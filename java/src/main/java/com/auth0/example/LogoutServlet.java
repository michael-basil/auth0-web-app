package com.auth0.example;

import javax.servlet.ServletException;
import javax.servlet.http.HttpServlet;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class LogoutServlet extends HttpServlet {

    private String domain;
    private String clientId;

    @Override
    public void init() throws ServletException {
        // Ensure configuration is loaded
        Auth0Config.init(getServletContext());
        this.domain = Auth0Config.getDomain();
        this.clientId = Auth0Config.getClientId();
    }

    @Override
    protected void doGet(HttpServletRequest request, HttpServletResponse response)
            throws IOException {

        if (request.getSession(false) != null) {
            request.getSession().invalidate();
        }

        String returnUrl = buildReturnUrl();

        String logoutUrl = String.format(
                "https://%s/v2/logout?client_id=%s&returnTo=%s",
                domain,
                clientId,
                returnUrl
        );

        response.sendRedirect(logoutUrl);
    }

    private String buildReturnUrl() {
        String baseUrl = Auth0Config.getAppBaseUrl();
        return baseUrl + "/login";
    }
}
