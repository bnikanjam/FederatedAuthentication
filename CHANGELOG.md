# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2024-07-26

### Added
- **Phase 3: Backend Discovery API Complete**:
    - **Backend API**: Implemented `/api/directory/lookup` endpoint using Gin and GORM.
    - **Database**: Configured Azure SQL Edge with `Organizations` table, auto-migration, and seeding.
    - **Frontend**: Updated `LoginComponent` to fetch organization ID from the backend API.
    - **Infrastructure**: Resolved SQL Server TLS and volume issues; added CORS middleware.
- **Phase 2: Dynamic Org & Multi-IdP Simulation Complete**:
    - **Dynamic Discovery**: Implemented `LoginComponent` to route users to the correct Auth0 Organization based on email domain.
    - **Multi-IdP Simulation**: Refactored Terraform to provision 5 distinct Organizations (LDAP, Azure, Okta, Google, SAML) with simulated connections.
    - **Policy Enforcement**: Enforced "Existing Users Only" by disabling signups on all simulated connections.
    - **Verification**: Validated end-to-end login flow for `alice@azure-corp.com` and confirmed absence of signup option.
- **Phase 1 MVP Complete**:
    - **Frontend**: Angular SPA with Auth0 SDK integration.
    - **Backend**: Go API with Auth0 JWT Middleware for protected routes.
    - **Infrastructure**: Docker Compose setup for full stack (Frontend, Backend, SQL, LDAP).
    - **Auth0**: Terraform configuration (`terraform/main.tf`) to provision Organization and Connection.
    - **Login Flow**: Validated "Business Users Only" flow with hardcoded Organization ID.
- **Documentation**:
    - `implementation_plan.md`: Updated with Phase 2 roadmap (Dynamic Org & Multi-IdP).
    - `troubleshooting.md`: Added guides for common Auth0 errors (Service Not Found, Org Required).
    - `walkthrough.md`: Verification steps for the current build.
    - `docs/architecture.md`: Added system architecture and sequence diagrams.

### Changed
- Refactored `frontend/src/app/app.component.ts` to support routing and dynamic login.
- Updated `frontend/src/app/app.config.ts` to remove hardcoded organization ID.
- Updated `frontend/angular.json` to include `src/assets` in build output.
- Updated `terraform/main.tf` to use `for_each` for dynamic resource creation and added a test user.

## [0.1.0] - 2025-11-19

### Added
- **Infrastructure**: Docker Compose setup with Azure SQL Edge, OpenLDAP, and phpLDAPadmin.
- **Backend**: Go (Gin) API service with JWT authentication middleware.
- **Frontend**: Angular 18 SPA configured with Auth0 SDK for authentication.
- **Security**: Basic JWT validation logic (mocked for initial connectivity).
- **Documentation**: `walkthrough.md` for setup and verification instructions.
- **CI/CD**: Initial Git repository initialization and push to GitHub.
