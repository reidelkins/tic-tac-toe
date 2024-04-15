import { Component } from '@angular/core';
import { Router } from '@angular/router'; // Import Router
import { LoginComponent } from '../login/login.component';
import { AnalyticsService } from '../core/services/analytics.service'; // Import AnalyticsService

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    LoginComponent
  ],
  providers: [
    AnalyticsService
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent {
  constructor(    
    private router: Router,
    private analyticsService: AnalyticsService
    ) {} 

  ngOnInit() {
    this.analyticsService.trackEvent('PAGE_LOADED', 'Header loaded', 'HEADER');
  }

  goHome() {
    this.router.navigate(['/']);
  }
}
