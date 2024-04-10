import { Component, OnInit, OnDestroy } from '@angular/core';
import { NgFor } from '@angular/common';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import confetti from 'canvas-confetti';
import Swal from 'sweetalert2';
import { Game } from '../game.model';
import { GameWebsocketService } from '../core/services/web-socket.service';
import { GameService } from '../core/services/game.service';


@Component({
  selector: 'app-game-board',
  standalone: true,
  imports: [
    NgFor
  ],
  templateUrl: './game-board.component.html',
  styleUrl: './game-board.component.scss'
})


export class GameBoardComponent implements OnInit, OnDestroy {
  gameState: Game = {
    ID: '0',
    Player1ID: 0,
    Player1Username: '',
    Player2ID: 0,
    Player2Username: '',
    Board: [['', '', ''], ['', '', ''], ['', '', '']],
    CurrentPlayer: 'X',
    Winner: '',
    Over: false,
    CreatedAt: '',
    UpdatedAt: ''
  };
  gameID: string = '';
  private gameStateSubscription: Subscription | undefined;

  constructor(
    private route: ActivatedRoute,
    private websocketService: GameWebsocketService,
    private gameService: GameService
    ) {}

  shareGame() {
    console.log('Share/Invite functionality to be implemented');
    // Implement sharing/inviting functionality
  }

  makeMove(row: number, col: number) {
    if (!this.gameState.Over && this.gameState.Board[row][col] === '' && this.isCurrentPlayer()) {
      const move = {
        playerId: this.websocketService.getPlayerID(),
        x: row,
        y: col
      };
      this.websocketService.sendMessage(move);
    } else {
      if (this.gameState.Over) {
        alert('Game over! Please start a new game.');
      } else if (!this.isCurrentPlayer()) {
        alert('It is not your turn!');
      }      
    }

  }

  getCellSymbol(row: number, col: number): string {
    return this.gameState.Board[row][col];
  }

  isCurrentPlayer(): boolean {
    const playerID = this.websocketService.getPlayerID();        
    return (
      (playerID === this.gameState.Player1ID && this.gameState.CurrentPlayer === 'X') ||
      (playerID === this.gameState.Player2ID && this.gameState.CurrentPlayer === 'O')
    );
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.gameID = params['gameID'];
      this.gameService.getGame(this.gameID).subscribe({
        next: (gameState) => {
          this.gameState = gameState;
          if (this.gameState.Player2ID === 0) {
            this.websocketService.playerID = this.gameState.Player1ID;
          } else {
            this.websocketService.playerID = this.gameState.Player2ID;
          }
          this.websocketService.connect(this.gameID);
        },
        error: (err) => {
          console.error('Error fetching game state:', err);
        }
      });
      // Subscribe to game state updates from the WebSocket service
      this.gameStateSubscription = this.websocketService.getGameUpdates().subscribe(
        (gameState: Game) => {          
          this.gameState = gameState;
          if (this.gameState.Over) {
            if (this.websocketService.getPlayerID() === this.gameState.Player1ID && this.gameState.CurrentPlayer === 'O') {            
              this.celebrateWin();
            } else {
              this.celebrateLoss();            
            }
          }
        }
      );
    });
  }

  ngOnDestroy() {
    this.websocketService.disconnect();
    if (this.gameStateSubscription) {
      this.gameStateSubscription.unsubscribe();
    }
  }

  celebrateWin() {
    // Create confetti effect
    confetti({
        particleCount: 100,
        spread: 70,
        origin: { y: 0.6 }
    });

    // Display a celebratory message
    Swal.fire({
        title: 'Congratulations!',
        text: 'You won!',
        icon: 'success',
        confirmButtonText: 'Cool'
    });
  }

  celebrateLoss() {
      Swal.fire({
          title: 'Oops!',
          text: 'You lost!',
          icon: 'error',
          confirmButtonText: 'Try Again'
      });
  }
  
}
