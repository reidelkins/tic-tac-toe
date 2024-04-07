import { Component, OnInit } from '@angular/core';
import {NgFor} from '@angular/common';
import { Router } from '@angular/router'; // Import Router

interface Game {
  id: string;
  players: string[];
  viewers: number;
}

@Component({
  selector: 'app-game-list',
  templateUrl: './game-list.component.html',
  styleUrls: ['./game-list.component.scss'],
  imports: [NgFor], // Add this line to import NgFor
  standalone: true // Add this line to make it a standalone component
})
export class GameListComponent implements OnInit {
  games: Game[] = []; // This will be fetched from the backend

  constructor(private router: Router) {} // Inject Router here

  ngOnInit(): void {
    // This would be replaced by a service call to fetch active games
    this.games = [
      { id: '1', players: ['Alice'], viewers: 5 },
      { id: '2', players: ['Bob', 'Charlie'], viewers: 3 }
      // more games...
    ];
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
