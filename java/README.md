# Auth0 Java Web App (Servlet + Maven + Jetty)

This example shows a minimal Java web application that uses **Auth0** for authentication.  
It runs inside **GitHub Codespaces** using Maven + Jetty and uses an external base URL for proper redirect handling.

---

## 1. Configure Auth0

Copy the example config file:

```bash
cp src/main/resources/auth0.properties.example src/main/resources/auth0.properties
```

Edit `auth0.properties`:

```properties
auth0.domain=YOUR_DOMAIN
auth0.clientId=YOUR_CLIENT_ID
auth0.clientSecret=YOUR_CLIENT_SECRET

# Codespaces external URL (no :3000)
app.baseUrl=https://YOUR-CODESPACES-URL
```

---

## 2. Update Auth0 Application Settings

### **Allowed Callback URLs**
```
https://YOUR-CODESPACES-URL/callback
```

### **Allowed Logout URLs**
```
https://YOUR-CODESPACES-URL/login
```

---

## 3. Login Flow

- Visit `/` → click “Go to Portal Home”
- `/portal/home` is protected and triggers redirect to `/login`
- `/login` redirects to Auth0 Universal Login
- After login, Auth0 redirects back to `/callback`
- Tokens are exchanged and stored in session
- User lands on `/portal/home`

---

## 4. Build & Run (Maven)

To build and run Jetty locally or in Codespaces:

```bash
mvn clean jetty:run
```

Jetty will start on port **3000**, and Codespaces will forward the public URL automatically.