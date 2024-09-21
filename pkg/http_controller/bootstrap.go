package http_controller

import (
	"net/http"
	"time"

	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func Bootstrap(handler *gin.Engine, l logging.Logger, appName string) {
	handler.UseRawPath = true

	l.Info("Bootstrapping gin engine :)")

	// Debug print all routes using zap
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		l.Infof("%-6s %-25s --> %s (%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// Datadog trace
	handler.Use(gintrace.Middleware(appName, gintrace.WithAnalytics(true)))

	// Cors
	handler.Use(cors.Default())

	// Prints incoming requests and results
	handler.Use(ginzap.Ginzap(l.Desugar(), time.RFC3339, true))

	// Generic Error handling :)
	handler.Use(GenericErrorHandler(l))

	// Middleware
	handler.Use(ginzap.RecoveryWithZap(l.Desugar(), true))

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	l.Info("Gin engine bootstrapped!")
}
