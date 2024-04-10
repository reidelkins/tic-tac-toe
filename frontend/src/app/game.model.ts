export interface Game {
  ID: string;
  Player1ID: number;
  Player1Username: string;
  Player2ID: number;
  Player2Username: string;
  Board: string[][];
  CurrentPlayer: string;
  Winner: string;
  Over: boolean;
  CreatedAt: string;
  UpdatedAt: string;
}