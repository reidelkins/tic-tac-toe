import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Game } from '../../game.model'; // Adjust the import according to your game model path
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class GameService {
  constructor(private http: HttpClient) {}

  listActiveGames(): Observable<Game[]> {
    return this.http.get<Game[]>(`${environment.backendUrl}/list-active-games`);
  }

  createGame(): Observable<string> {
    return this.http.post<string>(`${environment.backendUrl}/create-game`, { }, { withCredentials: true });
  }

  getGame(gameId: string): Observable<Game> {
    return this.http.get<Game>(`${environment.backendUrl}/get-game/${gameId}`);
  }

  joinGame(gameId: string, player2Username: string): Observable<string> {
    return this.http.post<string>(`${environment.backendUrl}/join-game`, { gameId, player2Username }, { withCredentials: true });
  }
}