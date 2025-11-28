package handlers

import (
    "net/http"
    "net/url"
    "os"

    "github.com/gin-gonic/gin"
)

// Logout clears the Auth0 session and redirects through Auth0.
func Logout(ctx *gin.Context) {
    logoutURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
    if err != nil {
        ctx.String(http.StatusInternalServerError, err.Error())
        return
    }

    baseURL := os.Getenv("APP_BASE_URL")
    if baseURL == "" {
        scheme := "http"
        if ctx.Request.TLS != nil {
            scheme = "https"
        }
        baseURL = scheme + "://" + ctx.Request.Host
    }

    returnTo, err := url.Parse(baseURL)
    if err != nil {
        ctx.String(http.StatusInternalServerError, err.Error())
        return
    }

    parameters := url.Values{}
    parameters.Add("returnTo", returnTo.String())
    parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))

    logoutURL.RawQuery = parameters.Encode()

    ctx.Redirect(http.StatusTemporaryRedirect, logoutURL.String())
}
