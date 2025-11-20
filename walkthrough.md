# Federation Authentication MVP - Walkthrough

## Status
The application stack is currently **RUNNING**.

## Access Points
| Service | URL | Credentials / Notes |
| :--- | :--- | :--- |
| **Frontend** | [http://localhost:4200](http://localhost:4200) | Angular App |
| **Backend API** | [http://localhost:8080](http://localhost:8080) | Go API (Health: `/health`) |
| **LDAP Admin** | [https://localhost:6443](https://localhost:6443) | Login: `cn=admin,dc=mycompany,dc=com` / `adminpassword` |
| **SQL Server** | `localhost:1433` | User: `sa` / Pass: `YourStrong!Passw0rd` |

## Next Steps: Auth0 Configuration
To make the login functional, you must configure Auth0:

1.  **Create Auth0 Tenant**: If you haven't already.
2.  **Create SPA Application**:
    -   Allowed Callback URLs: `http://localhost:4200`
    -   Allowed Logout URLs: `http://localhost:4200`
    -   Allowed Web Origins: `http://localhost:4200`
3.  **Create API**:
    -   Identifier: `https://your-api-identifier` (Update `.env` and `app.config.ts` with this)
4.  **Update Configuration**:
    -   Edit `frontend/src/app/app.config.ts` with your `domain`, `clientId`, and `audience`.
    -   Edit `.env` with your `AUTH0_DOMAIN` and `AUTH0_AUDIENCE`.
    -   Restart containers: `docker-compose up -d`

## LDAP Federation (Optional for MVP Start)
To connect Auth0 to the local LDAP:
1.  Expose port `389` via `ngrok tcp 389`.
2.  In Auth0, create a "Custom Database" or use the "AD/LDAP Connector" pointing to the ngrok address.
