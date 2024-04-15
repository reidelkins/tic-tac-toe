import { Component } from '@angular/core';
import { Router } from '@angular/router'; // Import Router
import { FormsModule } from '@angular/forms';
import { NgIf } from '@angular/common';
import { SocialAuthService, SocialUser } from '@abacritt/angularx-social-login';
import { LeaderboardComponent } from '../leaderboard/leaderboard.component';
import { GameListComponent } from '../game-list/game-list.component';
import { GameService } from '../core/services/game.service';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [
    FormsModule,
    GameListComponent,
    LeaderboardComponent,
    NgIf    
  ],
  templateUrl: './home.component.html',
  styleUrl: './home.component.scss'
})
export class HomeComponent {  
  user: SocialUser | null = null;

  constructor(
    private gameService: GameService, 
    private router: Router,
    private authService: SocialAuthService
    ) {} 

  ngOnInit() {
    this.authService.authState.subscribe((user: SocialUser | null) => {
      this.user = user;
    });
  }

  createGame() {
    this.gameService.createGame().subscribe({
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
