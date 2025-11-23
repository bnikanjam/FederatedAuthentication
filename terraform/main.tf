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

# 1. Create the Organization
resource "auth0_organization" "test_biz" {
  name         = "test-biz"
  display_name = "Test Business Inc."
  branding {
    logo_url = "https://example.com/logo.png"
  }
}

# 2. Create the Connection (Simulated LDAP via Custom DB or just Username-Password for now)
# Note: Real LDAP requires the Auth0 Connector to be installed on your infrastructure.
# For this MVP, we'll use a standard DB connection to simulate the "Own IdP".
resource "auth0_connection" "ldap_sim" {
  name     = "ldap-sim-connection"
  strategy = "auth0" # Using standard DB for simplicity in Terraform, swap to 'ad' or 'ldap' if connector is ready
}

# 3. Enable the Connection for the Organization
resource "auth0_organization_connection" "test_biz_conn" {
  organization_id = auth0_organization.test_biz.id
  connection_id   = auth0_connection.ldap_sim.id
  assign_membership_on_login = true
}

# Output the Org ID for Angular
output "organization_id" {
  value = auth0_organization.test_biz.id
}
