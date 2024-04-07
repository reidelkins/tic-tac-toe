import { Component, OnInit } from '@angular/core';
import {NgFor} from '@angular/common';

interface Player {
  name: string;
  wins: number;
  losses: number;
}

@Component({
  selector: 'app-leaderboard',
  templateUrl: './leaderboard.component.html',
  styleUrls: ['./leaderboard.component.scss'],
  imports: [NgFor], // Add this line to import NgFor
  standalone: true // Add this line to make it a standalone component
})
export class LeaderboardComponent implements OnInit {
  players: Player[] = [];

  constructor() { }

  ngOnInit(): void {
    // This would be replaced by a call to your service to fetch leaderboard data
    this.players = [
      { name: 'Player1', wins: 5, losses: 3 },
      { name: 'Player2', wins: 4, losses: 4 }
      // more players...
    ];
  }
}
