package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	authenticator "goAuth0/platform/auth"
)

// Login starts the Auth0 login flow.
func Login(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		codeVerifier, err := generateCodeVerifier()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}
		codeChallenge := codeChallengeFromVerifier(codeVerifier)

		// Save the PKCE state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		session.Set("code_verifier", codeVerifier)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		options := []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("code_challenge", codeChallenge),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		}
		if audience := os.Getenv("AUTH0_AUDIENCE"); audience != "" {
			options = append(options, oauth2.SetAuthURLParam("audience", audience))
		}

		redirectURL := auth.AuthCodeURL(state, options...)
		log.Printf("[login] generated state %s pkce challenge %s redirecting to %s", state, codeChallenge, redirectURL)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b), nil
}

func generateCodeVerifier() (string, error) {
	// 32 random bytes -> 43 char base64url verifier (>= recommended length)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func codeChallengeFromVerifier(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
