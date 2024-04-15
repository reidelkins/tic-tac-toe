import { Component } from '@angular/core';
import { Router } from '@angular/router'; // Import Router
import { LoginComponent } from '../login/login.component';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    LoginComponent
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent {
  constructor(    
    private router: Router    
    ) {} 

  goHome() {
    this.router.navigate(['/']);
  }
}
