# go-starter

## TODO

- ✅ Create users with email
- ✅ log in with email
- ⬜ Handle sessions with middleware:
- access token (short lived and stateless)<br>
- refresh token (long lived and saved hashed in db)<br>
- Hash the Refresh Token<br>
- tokens saved in two seperate HTTP-only, Secure cookies (preferably SameSite=Strict or Lax)
- both cookies sent in each request
- 🚀 1. High-Level Flow (Server-Side Auto-Refresh)
  Step-by-Step Flow:<br>
  Client makes an API request (HTMX or standard fetch).<br>
  Server validates the access token:<br>
  If valid → processes the request as usual.<br>
  If expired:<br>
  Checks for a valid refresh token (sent via HTTP-only cookie).<br>
  If the refresh token is valid:<br>
  Issues a new access token (and optionally rotates the refresh token).<br>
  Fulfills the original request without sending an error to the client.<br>
  Sends the new access token in an HTTP-only cookie along with the response.<br>
  If the refresh token is invalid/expired:<br>
  Responds with 401 Unauthorized (or 403 Forbidden).<br>

- ⬜ Sign up with google and apple
- ⬜ Login with different providers

Access token and refresh tokens are now being created when log in.
Need to create them when signing up
Need to implement middleware that checks whether logged in or not
Revoke session if no cookie found

## Tech Stack

This project functions as a starter for any webapp. It uses the stack

- Gin Web Framework
- HTMX
- Typescript
- Tailwind CSS

For Bundeling, Minifying and Compiling TS to JS I use Webpack.
It also uses Docker. There is a make file which starts the whole stack in development and production mode.

## Backend architecture

We will follow an MVP architecture for better separation of concerns.
This is the request's life cycle.

```mermaid
flowchart TD
    Client --A middleware injects\n a db transaction here\n So the service can access the transaction \n using the context of the request--> Handler
    Handler --> Service
    Service --> Repository
    Repository --through gorm--> Gorm/[(postgresql)]

    Handler ~~~|"The handler:<br>Does gin specific request parsing<br> Uses a service to fulfill the request"| Handler
    Service ~~~|"The service:<br>Contains the business logic<br>Uses the domain models and repository interface"| Service
    Repository ~~~|"The repository is an interface <br> implemented by the files in pkg/db <br> It uses gorm to access the database"| Repository
```
