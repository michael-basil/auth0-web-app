package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// User renders the logged-in user page.
func User(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile, ok := session.Get("profile").(map[string]interface{})
	if !ok || profile == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	prettyBytes, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.HTML(http.StatusOK, "user.html", gin.H{
		"Profile": profile,
		"Pretty":  string(prettyBytes),
	})
}
