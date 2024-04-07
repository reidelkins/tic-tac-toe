import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterOutlet } from '@angular/router';
// import { BrowserModule } from '@angular/platform-browser'

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    // BrowserModule
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  title = 'frontend';
}
