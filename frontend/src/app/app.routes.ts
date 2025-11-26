import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { AppComponent } from './app.component'; // Note: AppComponent is usually the root, we might need a separate HomeComponent or guard.

export const routes: Routes = [
    { path: '', component: LoginComponent },
    // We'll redirect back to home after login, which might need handling in AppComponent or a separate component
];

