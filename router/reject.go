// Package router defines the router layer of the Buddy System.
package router

import "github.com/gin-gonic/gin"

// Reject handles the signup reject request.
func Reject() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		// check auth level
	}
}
