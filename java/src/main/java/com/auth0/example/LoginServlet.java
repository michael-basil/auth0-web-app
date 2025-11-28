package com.auth0.example;

import com.auth0.AuthenticationController;
import com.auth0.SessionUtils;

import javax.servlet.ServletConfig;
import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;

public class LoginServlet extends HttpServlet {

    private AuthenticationController authenticationController;

    @Override
    public void init(ServletConfig config) throws ServletException {
        super.init(config);
        // Ensure config is loaded
        Auth0Config.init(config.getServletContext());
        this.authenticationController = AuthenticationControllerProvider.getInstance(config);
    }

    @Override
    protected void doGet(HttpServletRequest req, HttpServletResponse res) throws IOException {
        String redirectUri = Auth0Config.getAppBaseUrl() + "/callback";

        // Build the /authorize URL with Auth0's MVC helper
        String authorizeUrl = authenticationController
                .buildAuthorizeUrl(req, res, redirectUri)
                .withScope("openid profile email")
                .build();

        // Extract the 'state' from the authorize URL and store it in session
        String state = extractQueryParam(authorizeUrl, "state");
        if (state != null) {
            SessionUtils.set(req, "auth0State", state);
        }

        System.out.println("Using redirectUri: " + redirectUri);
        System.out.println("Authorize URL: " + authorizeUrl);

        res.sendRedirect(authorizeUrl);
    }

    private String extractQueryParam(String url, String paramName) {
        try {
            URI uri = new URI(url);
            String query = uri.getQuery();
            if (query == null) {
                return null;
            }
            String[] pairs = query.split("&");
            for (String pair : pairs) {
                int idx = pair.indexOf('=');
                if (idx > 0) {
                    String name = URLDecoder.decode(pair.substring(0, idx), StandardCharsets.UTF_8);
                    if (paramName.equals(name)) {
                        return URLDecoder.decode(pair.substring(idx + 1), StandardCharsets.UTF_8);
                    }
                }
            }
        } catch (URISyntaxException e) {
            e.printStackTrace();
        }
        return null;
    }
}
