// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplyC() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ActivityID string `json:"activity_id"`
			MemberID   string `json:"member_id"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		activityID, err := primitive.ObjectIDFromHex(body.ActivityID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.ApplyC(activityID, body.MemberID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func CancelC() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			ActivityID primitive.ObjectID `json:"activity_id"`
			MemberID   string             `json:"member_id"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.CancelC(body.ActivityID, body.MemberID); err != nil {
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func Capplies() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			ActivityID primitive.ObjectID `json:"activity_id"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		res, err := activity.Capplies(body.ActivityID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, res)
	}
}

func ApproveC() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			ActivityID primitive.ObjectID `json:"activity_id"`
			MemberID   string             `json:"member_id"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.ApproveC(body.ActivityID, body.MemberID); err != nil {
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func RejectC() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			ActivityID primitive.ObjectID `json:"activity_id"`
			MemberID   string             `json:"member_id"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.RejectC(body.ActivityID, body.MemberID); err != nil {
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
