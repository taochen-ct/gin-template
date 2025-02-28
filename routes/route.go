package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

// ProviderSet is router providers.
var ProviderSet = wire.NewSet(CreateBaseRouter)

func CreateBaseRouter() *gin.Engine {

	router := gin.New()
	router.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	apiGroup := router.Group("/api")
	if err := RegisterApiRouter(apiGroup); err != nil {
		panic(err)
	}
	return router
}

func RegisterApiRouter(router *gin.RouterGroup) (err error) {
	routeFuncs := []func(*gin.RouterGroup) error{
		RegisterTestRoute,
	}

	for _, register := range routeFuncs {
		if err = register(router); err != nil {
			return err
		}
	}
	return nil
}
