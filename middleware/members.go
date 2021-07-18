// Package middleware defines the router layer of the Buddy System.
package middleware

import "github.com/gin-gonic/gin"

// Members handles the member query request.
func Members() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		// check auth level
	}
}
