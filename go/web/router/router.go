package router

import (
    "encoding/gob"
    "net/http"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"

    authenticator "goAuth0/platform/auth"
    "goAuth0/web/handlers"
)

// New registers the routes and returns the router.
func New(auth *authenticator.Authenticator) *gin.Engine {
    router := gin.Default()

    // To store custom types in our cookies, we must first register them using gob.Register
    gob.Register(map[string]interface{}{})

    store := cookie.NewStore([]byte("secret")) // TODO: replace with stronger key for real use
    router.Use(sessions.Sessions("auth-session", store))

    router.Static("/public", "web/static")
    router.LoadHTMLGlob("web/template/*")

    router.GET("/", func(ctx *gin.Context) {
        ctx.HTML(http.StatusOK, "home.html", nil)
    })

    router.GET("/login", handlers.Login(auth))
    router.GET("/callback", handlers.Callback(auth))
    router.GET("/user", handlers.IsAuthenticated, handlers.User)
    router.GET("/logout", handlers.Logout)

    return router
}
