// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
	}
}

func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()
	}
}

func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ActivityID primitive.ObjectID `json:"id"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.Delete(body.ActivityID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func Participants() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ActivityID primitive.ObjectID `json:"id"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.Delete(body.ActivityID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
