import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '@auth0/auth0-angular';
import { HttpClient } from '@angular/common/http';

@Component({
    selector: 'app-login',
    standalone: true,
    imports: [CommonModule, FormsModule],
    template: `
    <div class="login-container">
      <h2>Sign In</h2>
      <p>Enter your work email to continue.</p>
      
      <div class="form-group">
        <label for="email">Email Address</label>
        <input type="email" id="email" [(ngModel)]="email" placeholder="user@company.com" (keyup.enter)="handleLogin()">
      </div>

      <button (click)="handleLogin()" [disabled]="!email">Continue</button>
      
      <p class="error" *ngIf="error">{{ error }}</p>
    </div>
  `,
    styles: [`
    .login-container { max-width: 400px; margin: 2rem auto; padding: 2rem; border: 1px solid #ccc; border-radius: 8px; text-align: center; }
    .form-group { margin-bottom: 1rem; text-align: left; }
    label { display: block; margin-bottom: 0.5rem; font-weight: bold; }
    input { width: 100%; padding: 0.5rem; box-sizing: border-box; }
    button { width: 100%; padding: 0.75rem; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; font-size: 1rem; }
    button:disabled { background: #ccc; cursor: not-allowed; }
    .error { color: red; margin-top: 1rem; }
  `]
})
export class LoginComponent {
    email: string = '';
    error: string = '';

    constructor(private auth: AuthService, private http: HttpClient) { }

    handleLogin() {
        if (!this.email || !this.email.includes('@')) {
            this.error = 'Please enter a valid email address.';
            return;
        }

        const domain = this.email.split('@')[1];

        // Call Backend API for Organization Lookup
        this.http.get<any>(`http://localhost:8080/api/directory/lookup?domain=${domain}`).subscribe({
            next: (response) => {
                const orgId = response.organization_id;

                if (orgId) {
                    // Dynamic Login with Organization ID
                    this.auth.loginWithRedirect({
                        authorizationParams: {
                            organization: orgId,
                            login_hint: this.email // Pre-fill email in Auth0
                        }
                    });
                } else {
                    this.error = `No organization found for domain: ${domain}`;
                }
            },
            error: (err) => {
                console.error('Failed to load org map', err);
                this.error = 'System error: Could not load organization directory.';
            }
        });
    }
}
