# Auth0 Go Web App (Minimal Example for Codespaces)

This example demonstrates a minimal **Go (Gin)** web application using
**Auth0** for authentication.\
It runs cleanly inside **GitHub Codespaces** using `go run .` and
supports Codespaces' external HTTPS URLs for the Auth0 redirect flow.

The project is intentionally simple: a single router, a login handler, a
callback handler, and a protected `/user` page.

------------------------------------------------------------------------

## 1. Configure Auth0

This project uses a simple `.env` file for configuration (easy in
Codespaces).

Inside the project folder (e.g., `/go`):

Create a file named `.env`:

    AUTH0_DOMAIN=YOUR_DOMAIN
    AUTH0_CLIENT_ID=YOUR_CLIENT_ID
    AUTH0_CLIENT_SECRET=YOUR_CLIENT_SECRET
    APP_BASE_URL=https://your-codespace-id-3000.app.github.dev
    AUTH0_AUDIENCE=https://YOUR-API-IDENTIFIER   # optional, required if you want an API Access Token
    AUTH0_SCOPES=openid profile offline_access  # optional override (default already requests offline_access)

> Codespaces gives you a public URL that looks like\
> `https://<id>-3000.app.github.dev`\
> whenever you run a server on port **3000**.

No secrets will be committed if `.env` is in `.gitignore` (recommended).

------------------------------------------------------------------------

## 2. Update Auth0 Application Settings

Assume your Codespaces public URL is:

    https://your-codespace-id-3000.app.github.dev

### **Allowed Callback URLs**

    https://your-codespace-id-3000.app.github.dev/callback

### **Allowed Logout URLs**

    https://your-codespace-id-3000.app.github.dev/

### **Allowed Web Origins**

    https://your-codespace-id-3000.app.github.dev

> Note: Even though your Go app listens on port 3000, Codespaces exposes
> it **without** `:3000` in the external URL.

### Enable Refresh Tokens

If you want to test refresh tokens:

1. In the Auth0 Dashboard, open your Application → **Settings → Advanced Settings → Grant Types** and ensure **Refresh Token** is enabled.
2. Under **Token Settings**, enable **Refresh Token Rotation** (recommended) and configure the absolute expiry to taste.
3. If you set `AUTH0_AUDIENCE`, make sure that API allows the `offline_access` scope.

------------------------------------------------------------------------

## 3. How the Login Flow Works (Go Version)

-   Visit `/`\
    → simple home page with a "Sign In" link

-   `/login`\
    → generates a state value\
    → builds a PKCE challenge (S256)\
    → redirects to Auth0 Universal Login (includes `audience` if configured)

-   After successful login\
    Auth0 redirects back to:

        GET /callback?code=...&state=...

-   The app:

    1.  Confirms the `state` matches\
    2.  Exchanges the code for tokens\
    3.  Verifies the ID Token using the Auth0 tenant's JWKS\
    4.  Stores the user profile in the session

-   Redirects the user to the protected page: `/user`

-   `/logout`\
    → redirects to Auth0's logout endpoint\
    → returns the user to your Codespaces URL

-   `/userinfo`\
    → uses the stored Access Token\
    → calls Auth0's `/userinfo` endpoint and returns the JSON response

-   `/refresh`\
    → exchanges the stored refresh token for new tokens\
    → updates the session (Auth0 refresh token rotation supported)

------------------------------------------------------------------------

## 4. Run the Application (Go 1.21+)

From inside the `go` folder:

``` bash
go run .
```

GitHub Codespaces will automatically forward port **3000**.

Open the forwarded HTTPS URL in your browser, and you should be able to:

-   Load the home page\
-   Click **Sign In**\
-   Use Auth0 Universal Login\
-   Return authenticated to `/user`\
-   View your profile info
-   Click **Call /userinfo** to see the Access Token in action\
-   Click **Refresh Tokens** to perform a refresh grant

------------------------------------------------------------------------

## 5. Notes

This example keeps configuration and structure intentionally minimal to
make it easy to understand, extend, and use as a Codespaces demo.
