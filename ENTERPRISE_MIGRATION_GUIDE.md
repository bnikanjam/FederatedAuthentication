# Enterprise Migration Guide: Federated Authentication

This guide outlines how to apply the Federated Authentication pattern (Angular + Go + Auth0 + Terraform) to your enterprise applications.

## 1. Architecture Overview

**Goal**: Enable B2B users to log in with their corporate credentials (Azure AD, Okta, etc.) via a single Angular SPA.

**Flow**:
1.  **User enters email** in Angular App.
2.  **Angular calls Backend API** (`/api/directory/lookup`) with the email domain.
3.  **Backend queries Database** to find the mapped Auth0 Organization ID.
4.  **Angular redirects to Auth0** using the `organization` ID and `login_hint`.
5.  **Auth0 routes user** to their corporate IdP (e.g., Azure AD) for authentication.
6.  **User returns** to Angular App with a valid token.

---

## 2. Frontend (Angular)

### Dependencies
- `@auth0/auth0-angular`

### Login Component Pattern
Create a component that captures the email and performs the lookup *before* triggering the Auth0 login.

```typescript
// login.component.ts
handleLogin() {
    const domain = this.email.split('@')[1];

    // 1. Lookup Organization ID
    this.http.get<any>(`/api/directory/lookup?domain=${domain}`).subscribe({
        next: (response) => {
            const orgId = response.organization_id;

            if (orgId) {
                // 2. Redirect to Auth0 with Org ID
                this.auth.loginWithRedirect({
                    authorizationParams: {
                        organization: orgId,
                        login_hint: this.email // Pre-fills email at IdP
                    }
                });
            } else {
                this.error = 'Organization not found';
            }
        }
    });
}
```

### Auth0 Configuration
Ensure your `AuthConfig` allows organization-based routing.

```typescript
// app.config.ts
provideAuth0({
    domain: 'YOUR_AUTH0_DOMAIN',
    clientId: 'YOUR_CLIENT_ID',
    authorizationParams: {
        redirect_uri: window.location.origin,
        // Do NOT hardcode organization here if you want dynamic routing
    }
})
```

---

## 3. Backend (Go)

### Database Schema
You need a table to map Email Domains to Auth0 Organization IDs.

```go
// models/organization.go
type Organization struct {
    gorm.Model
    Domain      string `gorm:"uniqueIndex;size:255;not null"` // e.g., "acme.com"
    Auth0OrgID  string `gorm:"not null"`                    // e.g., "org_xyz..."
    DisplayName string
}
```

### Lookup Endpoint
Expose a public (or rate-limited) endpoint for the frontend.

```go
// api/handlers.go
func GetOrganizationByDomain(c *gin.Context) {
    domain := c.Query("domain")
    var org models.Organization
    
    // Query DB
    if err := db.DB.Where("domain = ?", domain).First(&org).Error; err != nil {
        c.JSON(404, gin.H{"error": "Organization not found"})
        return
    }

    c.JSON(200, gin.H{
        "organization_id": org.Auth0OrgID,
        "display_name":    org.DisplayName,
    })
}
```

---

## 4. Infrastructure (Terraform)

Manage your Auth0 configuration as code to ensure consistency across environments.

### Resource: Organization
Create an Auth0 Organization for each B2B customer.

```hcl
resource "auth0_organization" "acme_corp" {
  name         = "acme-corp"
  display_name = "Acme Corp"
}
```

### Resource: Connection (Azure AD Example)
Configure the connection to the customer's IdP.

```hcl
resource "auth0_connection" "acme_azure" {
  name     = "acme-azure-conn"
  strategy = "waad" // Windows Azure AD

  options {
    client_id     = var.acme_client_id
    client_secret = var.acme_client_secret
    domain        = "acme.onmicrosoft.com"
    tenant_domain = "acme.onmicrosoft.com"
    waad_protocol = "openid-connect"
    identity_api  = "microsoft-identity-platform-v2.0"
  }
}
```

### Resource: Link Connection to Org
Crucial step: Enable the connection *only* for that organization.

```hcl
resource "auth0_organization_connection" "acme_link" {
  organization_id = auth0_organization.acme_corp.id
  connection_id   = auth0_connection.acme_azure.id
  
  assign_membership_on_login = true // Auto-add users to Org upon login
}
```

### Security Policy: Disable Public Signups
Prevent random users from signing up to these connections.

```hcl
resource "auth0_connection" "..." {
  // ...
  options {
    disable_signup = true
  }
}
```

---

## 5. Deployment Checklist

1.  **Secrets**: Store Client IDs and Secrets in a secure vault (e.g., Azure Key Vault, HashiCorp Vault) and inject them as Terraform variables.
2.  **CORS**: Ensure your Backend API allows requests from your Angular App's domain.
3.  **SSL/TLS**: Ensure all traffic (Frontend -> Backend, Frontend -> Auth0) is over HTTPS.
4.  **Rate Limiting**: Apply rate limiting to the `/api/directory/lookup` endpoint to prevent enumeration attacks.
