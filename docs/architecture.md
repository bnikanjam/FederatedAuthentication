# System Architecture & Data Flow (Phase 3)

## High-Level Architecture
This diagram illustrates the production-ready architecture with Backend Discovery and Real IdP Integration.

```mermaid
graph TD
    subgraph "Nextpower Enterprise Customer"
        User[Users]
        Browser[Web Browser]
        AzureAD["Azure Entra (IdP)"]
    end

    subgraph "Federation Broker"
        Auth0["Auth0 Tenant"]
    end

    subgraph "Nextpower Platform"
        Angular[Angular App]
        Nginx[Nginx Container]
        GoAPI[Go Backend API]
        SQL[Azure SQL Edge]
    end

    %% Relationships
    User -->|Accesses| Browser
    Browser -->|Loads App| Nginx
    Browser -- Runs --> Angular
    
    Angular -->|"1. Lookup Org (Email Domain)"| GoAPI
    GoAPI -->|2. Query Domain| SQL
    SQL -- Returns Org ID --> GoAPI
    GoAPI -- Returns Org ID --> Angular

    Angular -->|"3. Login (OIDC + Org ID)"| Auth0
    Auth0 -->|"4. Federate (SAML/OIDC)"| AzureAD
    AzureAD -- Federated Response --> Auth0
    Auth0 -- Tokens --> Angular
    
    Angular -->|"5. API Request (Bearer Token)"| GoAPI
    GoAPI -->|"6. Validate Token (JWKS)"| Auth0
```

## Authentication & API Sequence Flow
## Authentication & API Sequence Flow
This sequence diagram details the **"Dynamic Discovery"** login flow with a **"No Prompt"** experience for Business Users.

### Workflow Description
1.  **Organization Discovery**: The user enters their email (e.g., `alice@azure-corp.com`). The Angular App calls the Nextpower API to lookup the `organization_id` associated with that domain.
2.  **Direct Federation (No Prompt)**: The Angular App initiates the Auth0 login, passing the `organization_id` explicitly.
    *   *Configuration*: The Auth0 Application is configured for **"Business Users"** (Team Members) only.
    *   *Experience*: Auth0 skips the Universal Login prompt and immediately redirects the user to their configured Identity Provider (Azure Entra).
3.  **Authentication**: The user authenticates with their corporate credentials at Azure Entra.
4.  **Token Issuance**: Azure Entra returns a SAML/OIDC response to Auth0. Auth0 issues an OIDC Access Token and ID Token to the Angular App.
5.  **Protected Access**: The Angular App uses the Access Token to call the protected Go Backend API. The Backend validates the token against Auth0's JWKS.

```mermaid
sequenceDiagram
    autonumber
    box "Nextpower Enterprise Customer"
        participant IdP as Azure Entra
        actor User as Users
    end

    box "Federation Broker"
        participant Auth0 as Auth0 Tenant
    end

    box "Nextpower"
        participant App as Angular App
        participant API as Go Backend
        participant DB as SQL Database
    end

    User->>App: Enter Email (alice@azure-corp.com)
    
    Note right of App: Phase 3: Backend Discovery
    App->>API: GET /api/directory/lookup?domain=azure-corp.com
    API->>DB: SELECT * FROM organizations WHERE domain = ...
    DB-->>API: Return Org ID (org_xyz...)
    API-->>App: Return { organization_id: "org_xyz..." }

    App->>Auth0: Redirect to /authorize<br/>(client_id, organization_id, login_hint)
    
    Note right of Auth0: Auth0 maps Org ID to Connection
    
    Auth0->>IdP: To the User's Company's Idp
    User->>IdP: Enter Credentials
    IdP-->>Auth0: Federated Response (SAML/OIDC)
    
    Auth0-->>App: Redirect with Authorization Code
    App->>Auth0: Exchange Code for Tokens
    Auth0-->>App: Return Access Token & ID Token
    
    Note over App: User is logged in

    Note right of App: Phase 4: Protected API Access
    App->>API: GET /api/messages<br/>(Authorization: Bearer <Access Token>)
    
    Note right of API: Middleware Validates JWT<br/>(Signature, Audience, Issuer)
    
    API-->>App: Return JSON Response
```

## Component Details

### 1. Frontend (Angular)
- **Role**: Single Page Application serving the UI.
- **Auth SDK**: `@auth0/auth0-angular`.
- **Logic**: `LoginComponent` extracts email domain and calls Backend API for discovery.
- **Hosting**: Served via Nginx in a Docker container.

### 2. Backend (Go)
- **Role**: Protected Resource Server & Directory Service.
- **Framework**: Gin Web Framework.
- **Database**: GORM with Azure SQL Edge.
- **Endpoints**:
  - `GET /api/directory/lookup`: Public endpoint for Org Discovery.
  - `GET /api/messages`: Protected endpoint requiring JWT.
- **Middleware**: CORS and Auth0 JWT Validation.

### 3. Database (Azure SQL Edge)
- **Role**: Stores Organization mappings.
- **Schema**: `Organizations` table (Domain -> Auth0 Org ID).
- **Management**: Auto-migrated and seeded by the Go Backend on startup.

### 4. Auth0 (Broker)
- **Role**: Central Federation Broker.
- **Application Settings**:
  - **Login Experience**: "Business Users" (Users must be a member of an organization).
  - **Prompt**: "No Prompt" (Application passes `organization` parameter to skip Auth0 selection screen).
- **Resources Managed**:
  - **Organizations**: One per B2B customer (e.g., "Real Azure Corp").
  - **Connections**: "waad" (Azure AD) and "auth0" (Simulated).
  - **Policies**: Signups disabled for B2B connections.

### 5. Infrastructure
- **Docker Compose**: Orchestrates Frontend, Backend, and SQL Server.
- **Terraform**: Provisions Auth0 resources (Orgs, Connections, Users) and manages secrets via `terraform.tfvars`.
