import { Component, Inject } from '@angular/core';
import { CommonModule, DOCUMENT } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { AuthService } from '@auth0/auth0-angular';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet],
  template: `
    <div class="container">
      <header>
        <h1>Federation Auth MVP</h1>
        <div *ngIf="auth.user$ | async as user">
          <p>Welcome, {{ user.name }}</p>
          <button (click)="logout()">Log out</button>
        </div>
        <div *ngIf="(auth.isAuthenticated$ | async) === false">
          <button (click)="login()">Log in</button>
        </div>
      </header>

      <main>
        <div *ngIf="auth.user$ | async as user">
          <h2>User Profile</h2>
          <pre>{{ user | json }}</pre>
          
          <button (click)="callApi()">Call Protected API</button>
          <div *ngIf="apiResponse">
            <h3>API Response:</h3>
            <pre>{{ apiResponse | json }}</pre>
          </div>
        </div>
      </main>
    </div>
  `,
  styles: [`
    .container { padding: 2rem; font-family: sans-serif; }
    header { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #ccc; padding-bottom: 1rem; margin-bottom: 2rem; }
    button { padding: 0.5rem 1rem; cursor: pointer; background: #007bff; color: white; border: none; border-radius: 4px; }
    button:hover { background: #0056b3; }
  `]
})
export class AppComponent {
  apiResponse: any;

  constructor(
    @Inject(DOCUMENT) public document: Document,
    public auth: AuthService
  ) { }

  login() {
    this.auth.loginWithPopup().subscribe();
  }

  logout() {
    this.auth.logout({
      logoutParams: { returnTo: this.document.location.origin }
    }).subscribe();
  }

  callApi() {
    // TODO: Implement API call with token
    // For MVP, we will just log the token for now
    this.auth.getAccessTokenSilently().subscribe({
      next: (token) => {
        console.log('Access Token:', token);
        this.fetchData(token);
      },
      error: (err) => console.error(err)
    });
  }

  fetchData(token: string) {
    fetch('http://localhost:8080/api/messages', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
      .then(response => response.json())
      .then(data => this.apiResponse = data)
      .catch(err => console.error(err));
  }
}
