package handlers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	authenticator "goAuth0/platform/auth"
)

// Refresh exchanges a refresh token for new tokens and updates the session.
func Refresh(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		refreshToken, _ := session.Get("refresh_token").(string)
		if refreshToken == "" {
			ctx.String(http.StatusBadRequest, "No refresh token available. Ensure offline_access scope is granted.")
			return
		}

		ts := auth.TokenSource(
			withHTTPClient(ctx.Request.Context(), auth.HTTPClient),
			&oauth2.Token{RefreshToken: refreshToken},
		)
		newToken, err := ts.Token()
		if err != nil {
			ctx.String(http.StatusBadGateway, "Failed to refresh the token: "+err.Error())
			return
		}

		session.Set("access_token", newToken.AccessToken)
		if newToken.RefreshToken != "" {
			session.Set("refresh_token", newToken.RefreshToken)
		}
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token_expires_in": time.Until(newToken.Expiry).Round(time.Second).String(),
			"refresh_token_rotated":   newToken.RefreshToken != "" && newToken.RefreshToken != refreshToken,
		})
	}
}
