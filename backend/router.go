package main

import "github.com/gin-gonic/gin"

// newRouter builds the HTTP routes served by the backend.
func newRouter() *gin.Engine {
	return gin.New()
}
