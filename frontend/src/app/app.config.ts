import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { provideAuth0 } from '@auth0/auth0-angular';

import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    provideHttpClient(),
    provideAuth0({
      domain: 'dev-bnik.us.auth0.com',
      clientId: '6DUcggvMzHN8HcWJ1JnlC9femCBeafhk',
      authorizationParams: {
        redirect_uri: window.location.origin,
        audience: 'https://fedauthoneapi/',
      },
    }),
  ]
};
