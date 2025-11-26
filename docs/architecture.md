# System Architecture & Data Flow (Phase 1)

## High-Level Architecture
This diagram illustrates the components of the Federated Authentication MVP and their interactions.

```mermaid
graph TD
    subgraph "Client Workstation"
        User[User]
        Browser[Web Browser]
        Angular[Angular SPA]
    end

    subgraph "Docker Host"
        Nginx[Nginx Container]
        GoAPI[Go Backend API]
    end

    subgraph "External Services"
        Auth0[Auth0 Tenant]
        Terraform[Terraform Cloud/Local]
    end

    subgraph "Identity Providers"
        UpstreamIdP[Test Business IdP\n(Simulated LDAP)]
    end

    %% Relationships
    User -->|Accesses| Browser
    Browser -->|Loads App| Nginx
    Browser -- Runs --> Angular
    
    Angular -->|1. Login (OIDC)| Auth0
    Auth0 -->|2. Federate| UpstreamIdP
    
    Angular -->|3. API Request (Bearer Token)| GoAPI
    GoAPI -->|4. Validate Token (JWKS)| Auth0
    
    Terraform -->|Configures| Auth0
```

## Authentication & API Sequence Flow
This sequence diagram details the "Business Users Only" login flow and subsequent protected API access.

```mermaid
sequenceDiagram
    autonumber
    actor User
    participant App as Angular SPA
    participant Auth0 as Auth0 Tenant
    participant IdP as Upstream IdP (Org)
    participant API as Go Backend

    Note over App, Auth0: Configured with Organization ID

    User->>App: Click "Log In"
    App->>Auth0: Redirect to /authorize<br/>(client_id, audience, scope, organization_id)
    
    Note right of Auth0: Auth0 detects Organization Context
    
    Auth0->>IdP: Redirect to Organization's Login Page
    User->>IdP: Enter Credentials
    IdP-->>Auth0: Authentication Success
    
    Auth0-->>App: Redirect with Authorization Code
    App->>Auth0: Exchange Code for Tokens
    Auth0-->>App: Return Access Token & ID Token
    
    Note over App: User is logged in

    User->>App: Click "Call API"
    App->>API: GET /api/messages<br/>Authorization: Bearer <Access Token>
    
    Note right of API: Middleware Validation
    API->>Auth0: Fetch JWKS (Public Keys)
    API->>API: Validate Token Signature & Claims
    
    alt Token Valid
        API-->>App: 200 OK (Protected Data)
        App-->>User: Display Data
    else Token Invalid
        API-->>App: 401 Unauthorized
    end
```

## Component Details

### 1. Frontend (Angular)
- **Role**: Single Page Application serving the UI.
- **Auth SDK**: `@auth0/auth0-angular`.
- **Configuration**:
  - `domain`: `dev-bnik.us.auth0.com`
  - `audience`: `https://fedauthoneapi/`
  - `organization`: `org_...` (Hardcoded for Phase 1)
- **Hosting**: Served via Nginx in a Docker container.

### 2. Backend (Go)
- **Role**: Protected Resource Server.
- **Framework**: Gin Web Framework.
- **Middleware**: Custom middleware using `github.com/auth0/go-jwt-middleware/v2`.
- **Validation**: Validates JWTs against Auth0's JSON Web Key Set (JWKS).

### 3. Auth0 (Broker)
- **Role**: Central Federation Broker.
- **Resources Managed**:
  - **API**: Represents the Go Backend.
  - **Application**: Represents the Angular Frontend.
  - **Organization**: Represents the customer ("Test Business").
  - **Connection**: The link to the Upstream IdP.

### 4. Infrastructure
- **Docker Compose**: Orchestrates the local runtime environment.
- **Terraform**: Provisions and manages the Auth0 Organization and Connection resources to ensure reproducibility.
