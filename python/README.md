# Auth0 Python Web App – Login Example

A minimal Flask application demonstrating how to authenticate users with **Auth0** using **OAuth 2.0 / OpenID Connect**.

## What’s inside

- **Flask** – web server, routing, and sessions  
- **Authlib** – OAuth2/OIDC client used to talk to Auth0  
- **python-dotenv** – loads configuration from `.env`  
- **ProxyFix** – ensures Codespaces external URL is used for redirects  

## Quick Start

### 1. Install dependencies

```bash
python3 -m venv .venv
source .venv/bin/activate       # Windows: .venv\Scripts\activate
pip install -r requirements.txt
```

### 2. Create a `.env` file

Fill in your Auth0 settings:

```env
AUTH0_DOMAIN=YOUR_TENANT.us.auth0.com
AUTH0_CLIENT_ID=YOUR_CLIENT_ID
AUTH0_CLIENT_SECRET=YOUR_CLIENT_SECRET
APP_SECRET_KEY=ANY_RANDOM_SECRET
SERVER_NAME=YOUR_SERVER_NAME
PREFERRED_URL_SCHEME=https
APPLICATION_ROOT=/
```

In your Auth0 Application settings, ensure **Allowed Callback URLs** includes:

```
https://YOUR-CODESPACES-URL/callback
```

In your Auth0 Application settings, ensure **Allowed Logout URLs** includes:

```
https://YOUR-CODESPACES-URL
```

### 3. Run the server

```bash
python server.py
```

Open the forwarded Codespaces URL in your browser, click **Log In**, and complete the Auth0 flow.  
You should be redirected back to the app with your profile information displayed.

---