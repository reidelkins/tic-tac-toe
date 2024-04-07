import { Routes } from '@angular/router';

import { HomeComponent } from './home/home.component';
import { GameBoardComponent } from './game-board/game-board.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'board', component: GameBoardComponent },
  { path: 'game/:gameID', component: GameBoardComponent },
  // other routes...
];
