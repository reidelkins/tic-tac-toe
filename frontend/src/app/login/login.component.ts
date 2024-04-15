import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SocialAuthService, GoogleSigninButtonModule, SocialUser } from '@abacritt/angularx-social-login';
import { LoginService } from '../core/services/login.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, GoogleSigninButtonModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  user: SocialUser | null = null;

  constructor(
    private authService: SocialAuthService,
    private loginService: LoginService
  ) {}

  ngOnInit(): void {
    this.authService.authState.subscribe((user: SocialUser | null) => {
      this.user = user;
      if (user) {
        this.loginService.loginWithGoogle(user.idToken).subscribe({
          next: (resp: { email: string, token: string }) => {
          // Save the email and token in your component if needed
            
            // Save the token to local storage or cookie
            document.cookie = `token=${resp.token}; path=/; secure; httpOnly`;
          },
          error: (err) => {
            console.error('Error logging in:', err);
          }
        });
      }
    });
  }

  signOut(): void {
    this.authService.signOut().then(() => {
      // Remove token from local storage
      document.cookie = 'token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
      this.user = null;
    });
  }
}