// Package router defines the router layer of the Buddy System.
package router

import "github.com/gin-gonic/gin"

// Approve handles the signup approval request.
func Approve() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		// check auth level
	}
}
