// Package member defines the router layer of the club member of the Buddy System.
package member

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
)

// SignUp handles the signup request.
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(member.Member)
		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		memb := member.New(body.ID, body.Name, body.Department, body.Grade, body.Phone, body.Email, body.Attendance)
		if err := memb.SignUp(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
