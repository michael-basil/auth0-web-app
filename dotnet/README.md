# Auth0 ASP.NET Core Web App (Minimal Example for Codespaces)

This example shows a minimal **ASP.NET Core** web application using **Auth0** for authentication.  
It runs cleanly inside **GitHub Codespaces** using `dotnet run` and uses the external Codespaces URL for correct redirect handling.

The project is intentionally small: one controller, one index page, one login link, and a protected page.

---

## 1. Configure Auth0

Set your Auth0 values using **dotnet user-secrets** (recommended – no secrets in git):

From the project folder:

```bash
dotnet user-secrets init
```

Then add:

```bash
dotnet user-secrets set "Auth0:Domain" "YOUR_DOMAIN"
dotnet user-secrets set "Auth0:ClientId" "YOUR_CLIENT_ID"
dotnet user-secrets set "Auth0:ClientSecret" "YOUR_CLIENT_SECRET"
```

---

## 2. Update Auth0 Application Settings

Let’s assume your Codespaces URL is something like:

```
https://your-codespaces-id-5000.app.github.dev
```

### **Allowed Callback URLs**

```
https://your-codespaces-id-5000.app.github.dev/callback
```

### **Allowed Logout URLs**

```
https://your-codespaces-id-5000.app.github.dev/
```

### **Allowed Web Origins**

```
https://your-codespaces-id-5000.app.github.dev
```

> Note: Codespaces publishes port **5000** externally without `:5000` in the URL.

---

## 3. How the Login Flow Works

- Visit `/`  
  → index page loads with a “Login” link

- `/Account/Login`  
  → triggers `Challenge("Auth0")` → redirects to Auth0 Universal Login

- After successful login  
  Auth0 redirects back to:

  ```
  GET /callback?code=...&state=...
  ```

  (using `response_mode=query`)

- The Auth0 middleware validates the tokens  
  → signs the user in  
  → issues an auth cookie

- You are redirected to the protected page: `/Home/Portal`

---

## 4. Build & Run (.NET 8)

Run the app:

```bash
dotnet run
```

Codespaces will automatically expose the public HTTPS URL for port **5000**.

Open the forwarded URL in your browser and test the Sign In flow.