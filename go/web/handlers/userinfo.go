package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	authenticator "goAuth0/platform/auth"
)

// UserInfo calls Auth0's /userinfo endpoint using the stored Access Token.
func UserInfo(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		accessToken, _ := session.Get("access_token").(string)
		if accessToken == "" {
			ctx.String(http.StatusUnauthorized, "No access token available in the session.")
			return
		}

		domain := os.Getenv("AUTH0_DOMAIN")
		if domain == "" {
			ctx.String(http.StatusInternalServerError, "AUTH0_DOMAIN is not configured.")
			return
		}

		userInfoURL := fmt.Sprintf("https://%s/userinfo", domain)
		req, err := http.NewRequestWithContext(ctx.Request.Context(), http.MethodGet, userInfoURL, nil)
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		req.Header.Set("Authorization", "Bearer "+accessToken)

		client := auth.HTTPClient
		if client == nil {
			client = http.DefaultClient
		}

		resp, err := client.Do(req)
		if err != nil {
			ctx.String(http.StatusBadGateway, fmt.Sprintf("Failed to call userinfo: %v", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			ctx.String(resp.StatusCode, fmt.Sprintf("Auth0 userinfo returned %s", resp.Status))
			return
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, payload)
	}
}
