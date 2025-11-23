# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2024-07-26

### Added
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

### Changed
- Updated `app.config.ts` to include `organization` parameter for B2B login support.
- Configured Frontend Dockerfile for multi-stage build (Node -> Nginx).

## [0.1.0] - 2025-11-19

### Added
- **Infrastructure**: Docker Compose setup with Azure SQL Edge, OpenLDAP, and phpLDAPadmin.
- **Backend**: Go (Gin) API service with JWT authentication middleware.
- **Frontend**: Angular 18 SPA configured with Auth0 SDK for authentication.
- **Security**: Basic JWT validation logic (mocked for initial connectivity).
- **Documentation**: `walkthrough.md` for setup and verification instructions.
- **CI/CD**: Initial Git repository initialization and push to GitHub.
