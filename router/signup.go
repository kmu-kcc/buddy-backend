// Package router defines the router layer of the Buddy System.
package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/model/member"
)

// SignUp handles the signup request.
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(member.Member)
		var (
			res struct {
				Error string `json:"error"`
			}
			err error
		)

		if err = json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			res.Error = err.Error()
			c.JSON(http.StatusBadRequest, res)
			return
		}

		guest := member.New(body.ID, body.Name, body.Department, body.Grade, body.Phone, body.Email, body.Enrollment)

		if err = guest.SignUp(); err == nil {
			err = errors.New("signup request success")
		}

		res.Error = err.Error()
		c.JSON(http.StatusOK, res)
	}
}
