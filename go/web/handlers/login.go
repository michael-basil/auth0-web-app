package handlers

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"

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

        // Save the state inside the session.
        session := sessions.Default(ctx)
        session.Set("state", state)
        if err := session.Save(); err != nil {
            ctx.String(http.StatusInternalServerError, err.Error())
            return
        }

        ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
    }
}

func generateRandomState() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(b), nil
}
