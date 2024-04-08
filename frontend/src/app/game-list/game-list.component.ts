import { Component, OnInit } from '@angular/core';
import {NgFor, NgIf} from '@angular/common';
import { Router } from '@angular/router'; // Import Router
import { Game } from '../game.model';
import { GameService } from '../core/services/game.service';

@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.scss'],
  imports: [NgFor, NgIf], // Add this line to import NgFor
  standalone: true // Add this line to make it a standalone component
})
export class GameListComponent implements OnInit {
  games: Game[] = []; // This will be fetched from the backend

  
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

  viewGame(gameId: string) {
    console.log(`Viewing game: ${gameId}`);
    // Logic to view a game
    this.router.navigate(['/game', gameId]);
  }

  joinGame(gameId: string) {
    console.log(`Joining game: ${gameId}`);
    // Logic to join a game
  }
}
