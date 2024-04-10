import { Injectable } from '@angular/core';
import { WebSocketSubject, webSocket } from 'rxjs/webSocket';
import { Observable, Subject } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Game } from '../../game.model';

@Injectable({
  providedIn: 'root'
})
export class GameWebsocketService {
  private socket$: WebSocketSubject<any> | null = null;
  public playerID: number = 0;
  private gameUpdates = new Subject<any>();

  constructor() {}

  public connect(gameId: string): void {
    if (!this.socket$ || this.socket$.closed) {
      this.socket$ = webSocket(`${environment.wsURL}/ws?gameId=${gameId}`);
      this.socket$.subscribe(
        message => {
          this.gameUpdates.next(message);          
        },
        error => console.log('WebSocket error:', error),
        () => console.log('WebSocket connection closed')
      );
    }
  }

  public isConnected(): boolean {
    return this.socket$ !== null && !this.socket$.closed;
  }

  public sendMessage(message: any): void {
    if (this.isConnected()) {
      this.socket$?.next(message);
    } else {
      console.log('WebSocket connection is not established');
    }
  }

  public getPlayerID(): number {
    return this.playerID;
  }

  public getGameUpdates(): Observable<Game> {     
    return this.gameUpdates.asObservable();
  }

  public disconnect(): void {
    if (this.isConnected()) {
      this.socket$?.complete();
      this.socket$ = null;
    }
  }
}