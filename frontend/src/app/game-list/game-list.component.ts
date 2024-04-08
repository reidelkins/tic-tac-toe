import { Component, OnInit } from '@angular/core';
import {NgFor, NgIf} from '@angular/common';
import { Router } from '@angular/router'; // Import Router
import { FormsModule } from '@angular/forms';
import { Game } from '../game.model';
import { GameService } from '../core/services/game.service';

@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.scss'],
  imports: [
    NgFor,
    NgIf,
    FormsModule
  ], // Add this line to import NgFor
  standalone: true // Add this line to make it a standalone component
})
export class GameListComponent implements OnInit {
  games: Game[] = []; // This will be fetched from the backend
  playerUsernames: { [key: string]: string } = {};

  
  constructor(
    private gameService: GameService, 
    private router: Router) {}

  ngOnInit(): void {
    // This would be replaced by a service call to fetch active games    
    this.gameService.listActiveGames().subscribe({
      next: (games) => {        
        if (games && games.length > 0) {
          this.games = games;             
        } else {
          // Handle the case where there are no games
          this.games = [];          
        }
      },
      error: (err) => {        
        console.error('Error fetching games:', err);
      },
    });
    
  }

  viewGame(gameId: number) {    
    // Logic to view a game
    this.router.navigate(['/game', gameId]);
  }

  joinGame(gameId: number) {
    const username = this.playerUsernames[gameId] || '';
    console.log(`Joining game: ${gameId} with username ${username}`);
    console.log(`Joining game: ${gameId}`);
    if (!username) {
      alert('Username is required to join a game');
      return;
    }
    // Logic to join a game
    this.gameService.joinGame(gameId, username).subscribe({
      next: (response) => {
        console.log('Joined game:', response);
        this.router.navigate(['/game', gameId]);
      },
      error: (err) => {
        console.error('Error joining game:', err);
      },
    });
  }
}
