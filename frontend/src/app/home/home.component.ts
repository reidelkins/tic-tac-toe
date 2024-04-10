import { Component } from '@angular/core';
import { Router } from '@angular/router'; // Import Router
import { FormsModule } from '@angular/forms';
import { LeaderboardComponent } from '../leaderboard/leaderboard.component';
import { GameListComponent } from '../game-list/game-list.component';
import { GameService } from '../core/services/game.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    FormsModule,
    GameListComponent,
    LeaderboardComponent
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss'
})
export class HomeComponent {
  username: string = '';

  constructor(
    private gameService: GameService, 
    private router: Router
    ) {} 

  createGame() {
    if (!this.username) {
      alert('Please enter a username');
      return;
    }
    console.log('Creating game for', this.username);
    this.gameService.createGame(this.username).subscribe({
      next: (gameId) => {
        // Navigate to the game page
        this.router.navigate(['/game', gameId]);
      },
      error: (err) => {
        console.error('Error creating game:', err);
      },
    });    
        
  }
}
