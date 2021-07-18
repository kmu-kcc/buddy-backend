// Package middleware defines the router layer of the Buddy System.
package middleware

import "github.com/gin-gonic/gin"

// Delete handles the member deletion request.
func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO
		// check auth level
	}
}
