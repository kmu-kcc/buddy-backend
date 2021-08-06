// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ApplyP handles the activity apply request
func ApplyP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})
		body := new(
			struct {
				ActivityID primitive.ObjectID `json:"_id"`
				MemberID   string             `json:"member_id"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		if err := activity.ApplyP(body.ActivityID, body.MemberID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

//Papplies handles the inquire of applicants list
func Papplies() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})
		body := new(
			struct {
				ID primitive.ObjectID `json:"_id"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		if _, err := activity.Papplies(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

//ApplyP handles the activity apply request
func ApproveP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})
		body := new(
			struct {
				ID  primitive.ObjectID `json:"_id"`
				IDs []string           `json:"member_ids"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		if err := activity.ApproveP(body.ID, body.IDs); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

//ApplyP handles the activity apply request
func RejectP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})
		body := new(
			struct {
				ID  primitive.ObjectID `json:"_id"`
				IDs []string           `json:"member_ids"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		if err := activity.RejectP(body.ID, body.IDs); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func CancelP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})
		body := new(
			struct {
				ActivityID primitive.ObjectID `json:"_id"`
				MemberID   string             `json:"member_id"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		if err := activity.CancelP(body.ActivityID, body.MemberID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
