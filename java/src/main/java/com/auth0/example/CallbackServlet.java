package com.auth0.example;

import com.auth0.SessionUtils;
import com.auth0.client.auth.AuthAPI;
import com.auth0.exception.Auth0Exception;
import com.auth0.json.auth.TokenHolder;
import com.auth0.net.TokenRequest;

import javax.servlet.ServletConfig;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class CallbackServlet extends HttpServlet {

    @Override
    public void init(ServletConfig config) throws ServletException {
        super.init(config);
        // Ensure configuration is loaded
        Auth0Config.init(config.getServletContext());
    }

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse res) throws IOException {
        handle(req, res);
    }

    @Override
    protected void doPost(HttpServletRequest req, HttpServletResponse res) throws IOException {
        handle(req, res);
    }

    private void handle(HttpServletRequest req, HttpServletResponse res) throws IOException {
        // 1. State validation
        String expectedState = (String) SessionUtils.get(req, "auth0State");
        String returnedState = req.getParameter("state");

        if (expectedState == null || returnedState == null || !expectedState.equals(returnedState)) {
            System.out.println("State mismatch. Expected: " + expectedState + ", got: " + returnedState);
            res.sendError(HttpServletResponse.SC_BAD_REQUEST, "Invalid state parameter.");
            return;
        }

        // 2. Authorization code
        String code = req.getParameter("code");
        if (code == null || code.isEmpty()) {
            res.sendError(HttpServletResponse.SC_BAD_REQUEST, "Missing authorization code.");
            return;
        }

        String domain = Auth0Config.getDomain();
        String clientId = Auth0Config.getClientId();
        String clientSecret = Auth0Config.getClientSecret();
        String redirectUri = Auth0Config.getAppBaseUrl() + "/callback";

        try {
            // Use older 1.x style AuthAPI constructor (from mvc-auth-commons' transitive dependency)
            AuthAPI authAPI = new AuthAPI(domain, clientId, clientSecret);

            TokenRequest tokenRequest = authAPI.exchangeCode(code, redirectUri);

            // In auth0-java 1.x, execute() returns TokenHolder directly
            TokenHolder tokenHolder = tokenRequest.execute();

            String accessToken = tokenHolder.getAccessToken();
            String idToken = tokenHolder.getIdToken();

            System.out.println("Successfully exchanged code for tokens.");
            System.out.println("access_token (truncated): " +
                    (accessToken != null ? accessToken.substring(0, Math.min(30, accessToken.length())) + "..." : "null"));

            // Store tokens in session and redirect to protected area
            SessionUtils.set(req, "accessToken", accessToken);
            SessionUtils.set(req, "idToken", idToken);

            //String portalHomeUrl = req.getContextPath() + "/portal/home";
            String portalHomeUrl = Auth0Config.getAppBaseUrl() + "/portal/home";
            System.out.println("Redirecting to portalHomeUrl: " + portalHomeUrl);
            res.sendRedirect(portalHomeUrl);

        } catch (Auth0Exception e) {
            System.out.println("Error while exchanging authorization code for tokens: " + e.getMessage());
            e.printStackTrace();
            res.sendRedirect(req.getContextPath() + "/login");
        }
    }
}
