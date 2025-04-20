package test

import (
	"github.com/gin-gonic/gin"
)

// Custom logger to reduce verbosity
type testLogger struct{}

func (l *testLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// SetupTestRouter creates a new Gin router with reduced logging for tests
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.LoggerWithWriter(&testLogger{}))
	return router
}
