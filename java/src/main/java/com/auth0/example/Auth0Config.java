package com.auth0.example;

import javax.servlet.ServletContext;
import java.io.IOException;
import java.io.InputStream;
import java.util.Properties;

public class Auth0Config {

    private static boolean initialized = false;

    private static String domain;
    private static String clientId;
    private static String clientSecret;
    private static String appBaseUrl; // external base URL (localhost or Codespaces)

    private Auth0Config() {
    }

    public static synchronized void init(ServletContext context) {
        if (initialized) {
            return;
        }

        String configPath = context.getInitParameter("auth0.configFile");
        if (configPath == null || configPath.isEmpty()) {
            configPath = "/auth0.properties";
        }

        Properties props = new Properties();
        try (InputStream is = Auth0Config.class.getResourceAsStream(configPath)) {
            if (is == null) {
                throw new IllegalStateException("Could not find Auth0 config file on classpath: " + configPath);
            }
            props.load(is);
        } catch (IOException e) {
            throw new RuntimeException("Failed to load Auth0 config from " + configPath, e);
        }

        domain = props.getProperty("auth0.domain");
        clientId = props.getProperty("auth0.clientId");
        clientSecret = props.getProperty("auth0.clientSecret");
        appBaseUrl = props.getProperty("app.baseUrl");

        if (domain == null || clientId == null || clientSecret == null || appBaseUrl == null) {
            throw new IllegalStateException(
                "Missing auth0.domain, auth0.clientId, auth0.clientSecret, or app.baseUrl in properties."
            );
        }

        initialized = true;
    }

    public static String getDomain() {
        assertInitialized();
        return domain;
    }

    public static String getClientId() {
        assertInitialized();
        return clientId;
    }

    public static String getClientSecret() {
        assertInitialized();
        return clientSecret;
    }

    public static String getAppBaseUrl() {
        assertInitialized();
        return appBaseUrl;
    }

    private static void assertInitialized() {
        if (!initialized) {
            throw new IllegalStateException("Auth0Config not initialized. Call Auth0Config.init(ServletContext) first.");
        }
    }
}
