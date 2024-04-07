import { Component } from '@angular/core';
import { Router } from '@angular/router'; // Import Router
import { FormsModule } from '@angular/forms';
import { LeaderboardComponent } from '../leaderboard/leaderboard.component';
import { GameListComponent } from '../game-list/game-list.component';

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

  constructor(private router: Router) {} // Inject Router here

  createGame() {
    if (!this.username) {
      alert('Please enter a username');
      return;
    }
    console.log('Creating game for', this.username);
    // Logic to create a game, likely involving calling a service method

    // Navigate to the game route with a static gameID for now
    this.router.navigate(['/game', '1']);
  }
}
