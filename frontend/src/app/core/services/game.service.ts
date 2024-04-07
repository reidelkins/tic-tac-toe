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
}