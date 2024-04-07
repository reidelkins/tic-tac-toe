import { Component } from '@angular/core';
import { OnInit } from '@angular/core';
import { NgFor } from '@angular/common';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-game-board',
  standalone: true,
  imports: [
    NgFor
  ],
  templateUrl: './game-board.component.html',
  styleUrl: './game-board.component.scss'
})


export class GameBoardComponent implements OnInit {
  players: string[] = ['Player 1']; // Initially one player, add second player when joined
  viewersCount: number = 0;
  board: string[][] = [['', '', ''], ['', '', ''], ['', '', '']];
  gameID: string | undefined;

  constructor(private route: ActivatedRoute) {}

  shareGame() {
    console.log('Share/Invite functionality to be implemented');
    // Implement sharing/inviting functionality
  }

  makeMove(row: number, col: number) {
    if (!this.board[row][col]) {
      this.board[row][col] = 'X'; // Example move, modify as per game logic
      console.log(`Move at ${row},${col}`);
      // Add logic to handle the move, check for win, etc.
    }
  }

  getCellSymbol(row: number, col: number): string {
    return this.board[row][col];
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      this.gameID = params['gameID'];
      console.log('Game ID:', this.gameID);
      // You can now use this.gameID to load the specific game data
    });
  }
}
