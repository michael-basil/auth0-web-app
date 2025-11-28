package com.auth0.example;

import com.auth0.SessionUtils;

import javax.servlet.*;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class Auth0Filter implements Filter {

    @Override
    public void init(FilterConfig filterConfig) throws ServletException {
        // Make sure config is loaded so we can use appBaseUrl
        Auth0Config.init(filterConfig.getServletContext());
    }

    @Override
    public void doFilter(ServletRequest req, ServletResponse res, FilterChain chain)
            throws IOException, ServletException {

        HttpServletRequest request = (HttpServletRequest) req;
        HttpServletResponse response = (HttpServletResponse) res;

        String accessToken = (String) SessionUtils.get(request, "accessToken");
        String idToken = (String) SessionUtils.get(request, "idToken");

        if (accessToken == null && idToken == null) {
            // *** This is the redirect you are seeing ***
            String loginUrl = Auth0Config.getAppBaseUrl() + "/login";
            System.out.println("Auth0Filter redirecting to: " + loginUrl);
            response.sendRedirect(loginUrl);
            return;
        }

        chain.doFilter(req, res);
    }

    @Override
    public void destroy() {
        // no-op
    }
}
