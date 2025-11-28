package handlers

import (
    "net/http"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
)

// User renders the logged-in user page.
func User(ctx *gin.Context) {
    session := sessions.Default(ctx)
    profile := session.Get("profile")
    ctx.HTML(http.StatusOK, "user.html", profile)
}
