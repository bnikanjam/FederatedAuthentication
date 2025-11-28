terraform {
  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = "~> 1.0"
    }
  }
}

provider "auth0" {
  domain        = var.auth0_domain
  client_id     = var.auth0_client_id
  client_secret = var.auth0_client_secret
}

variable "auth0_domain" {}
variable "auth0_client_id" {}
variable "auth0_client_secret" {}

variable "azure_client_id" {}
variable "azure_client_secret" {}

# Define the organizations and their simulated IdP types
locals {
  organizations = {
    "ldap-corp" = {
      domain       = "ldap-corp.com"
      display_name = "LDAP Corp (Simulated)"
      conn_name    = "sim-ldap-conn"
    },
    "azure-corp" = {
      domain       = "azure-corp.com"
      display_name = "Azure Corp (Simulated)"
      conn_name    = "sim-azure-conn"
    },
    "okta-corp" = {
      domain       = "okta-corp.com"
      display_name = "Okta Corp (Simulated)"
      conn_name    = "sim-okta-conn"
    },
    "google-corp" = {
      domain       = "google-corp.com"
      display_name = "Google Corp (Simulated)"
      conn_name    = "sim-google-conn"
    },
    "saml-corp" = {
      domain       = "saml-corp.com"
      display_name = "SAML Corp (Simulated)"
      conn_name    = "sim-saml-conn"
    },
    "real-azure-corp" = {
      domain       = "ba2kxoutlook.onmicrosoft.com"
      display_name = "Real Azure Corp"
      conn_name    = "real-azure-ad-conn"
    }
  }
}

# 1. Create Organizations
resource "auth0_organization" "orgs" {
  for_each     = local.organizations
  name         = each.key
  display_name = each.value.display_name
  branding {
    logo_url = "https://example.com/logo.png"
  }
}

# 2. Create Connections (Simulated & Real)
resource "auth0_connection" "conns" {
  for_each = local.organizations
  name     = each.value.conn_name
  
  # Use 'waad' for the real Azure connection, 'auth0' for simulations
  strategy = each.key == "real-azure-corp" ? "waad" : "auth0"

  # Options for Real Azure AD
  dynamic "options" {
    for_each = each.key == "real-azure-corp" ? [1] : []
    content {
      client_id     = var.azure_client_id
      client_secret = var.azure_client_secret
      domain        = "ba2kxoutlook.onmicrosoft.com" # Your Azure Tenant Domain
      tenant_domain = "ba2kxoutlook.onmicrosoft.com" # Your Azure Tenant Domain
      waad_protocol = "openid-connect"
      identity_api  = "microsoft-identity-platform-v2.0"
    }
  }

  # Options for Simulated Connections
  dynamic "options" {
    for_each = each.key != "real-azure-corp" ? [1] : []
    content {
      disable_signup = true
    }
  }
}

resource "auth0_connection_clients" "enable_clients" {
  for_each      = auth0_connection.conns
  connection_id = each.value.id
  enabled_clients = [
    var.auth0_client_id,               # M2M App (Terraform)
    "6DUcggvMzHN8HcWJ1JnlC9femCBeafhk" # Angular Frontend
  ]
}

# 3. Link Connection to Organization
resource "auth0_organization_connection" "org_conns" {
  for_each = local.organizations
  
  organization_id = auth0_organization.orgs[each.key].id
  connection_id   = auth0_connection.conns[each.key].id
  
  assign_membership_on_login = true
}

# 4. Create a Test User for Verification
resource "auth0_user" "alice_azure" {
  connection_name = "sim-azure-conn"
  email           = "alice@azure-corp.com"
  password        = "Password123!" # Simple password for testing
  email_verified  = true
  
  depends_on = [
    auth0_connection.conns,
    auth0_connection_clients.enable_clients
  ]
}

# Output the Map for Frontend
output "org_map" {
  value = {
    for k, v in local.organizations : v.domain => auth0_organization.orgs[k].id
  }
}
