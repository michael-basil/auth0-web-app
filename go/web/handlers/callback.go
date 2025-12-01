package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	authenticator "goAuth0/platform/auth"
)

// Callback handles the Auth0 redirect and stores the user session.
func Callback(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		queryState := ctx.Query("state")
		authorizationCode := ctx.Query("code")
		codeVerifier, _ := session.Get("code_verifier").(string)

		log.Printf("[callback] received redirect state=%s codeLen=%d", queryState, len(authorizationCode))

		if queryState != session.Get("state") {
			log.Printf("[callback] state mismatch session=%v query=%s", session.Get("state"), queryState)
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		// Exchange an authorization code for a token.
		exchangeCtx := withHTTPClient(ctx.Request.Context(), auth.HTTPClient)
		token, err := auth.Exchange(
			exchangeCtx,
			authorizationCode,
			oauth2.SetAuthURLParam("code_verifier", codeVerifier),
		)
		if err != nil {
			log.Printf("[callback] token exchange failed: %v", err)
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}
		log.Printf("[callback] token exchange succeeded expiresIn=%s", time.Until(token.Expiry).Round(time.Second))

		verifyCtx := withHTTPClient(ctx.Request.Context(), auth.HTTPClient)
		idToken, err := auth.VerifyIDToken(verifyCtx, token)
		if err != nil {
			log.Printf("[callback] failed to verify ID token: %v", err)
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		log.Printf("[callback] verified ID token for subject=%v", profile["sub"])

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		session.Delete("code_verifier")
		session.Delete("state")
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		log.Printf("[callback] session established for subject=%v profileEmail=%v", profile["sub"], profile["email"])

		// Redirect to logged-in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}

func withHTTPClient(ctx context.Context, client *http.Client) context.Context {
	if client == nil {
		return ctx
	}
	return context.WithValue(ctx, oauth2.HTTPClient, client)
}
