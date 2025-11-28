package com.auth0.example;

import com.auth0.AuthenticationController;
import com.auth0.jwk.JwkProvider;
import com.auth0.jwk.JwkProviderBuilder;

import javax.servlet.ServletConfig;

public class AuthenticationControllerProvider {

    private static AuthenticationController INSTANCE;

    private AuthenticationControllerProvider() {
    }

    public static synchronized AuthenticationController getInstance(ServletConfig config) {
        if (INSTANCE == null) {
            // ensure config is loaded
            Auth0Config.init(config.getServletContext());

            String domain = Auth0Config.getDomain();
            String clientId = Auth0Config.getClientId();
            String clientSecret = Auth0Config.getClientSecret();

            JwkProvider jwkProvider = new JwkProviderBuilder(domain).build();

            INSTANCE = AuthenticationController.newBuilder(domain, clientId, clientSecret)
                    .withJwkProvider(jwkProvider)
                    .build();
        }
        return INSTANCE;
    }
}
