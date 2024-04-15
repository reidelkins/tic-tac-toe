import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  constructor(private http: HttpClient) {}

  loginWithGoogle(token: string) {
    return this.http.post<string>(`${environment.backendUrl}/login/google`, { token });
  }
}
