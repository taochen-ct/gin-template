package routes

import (
	"awesomeProject/internal/middleware"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(CreateBaseRouter)

func CreateBaseRouter(
	recoveryMiddleware *middleware.Recovery,
	corsMiddleware *middleware.Cors,
	limiterMiddleware *middleware.Limiter,
) *gin.Engine {

	// create new gin engine
	router := gin.New()
	// add middleware
	router.Use(
		gin.Logger(),                 // default logger
		recoveryMiddleware.Handler(), // logger
		corsMiddleware.Handler(),     // cors
		limiterMiddleware.Handler(),  // request rate limiter
	)
	// website server
	router.Use(static.Serve("/", static.LocalFile("web", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("web/index.html")
	})

	// test service health
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	// create api group
	apiPrefix := "/api/v1"
	apiGroup := router.Group(apiPrefix)
	// register customer api
	if err := RegisterApiRouter(apiGroup); err != nil {
		panic(err)
	}
	return router
}

func RegisterApiRouter(router *gin.RouterGroup) (err error) {
	routeFuncs := []func(*gin.RouterGroup) error{
		// add register func here
		RegisterTestRoute,
	}

	for _, register := range routeFuncs {
		if err = register(router); err != nil {
			return err
		}
	}
	return nil
}
