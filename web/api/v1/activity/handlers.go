// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create handles the activity creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(activity.Activity)
		resp := new(struct {
			Error string `json:"error"`
		})

		// if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.
			// New(body.Start, body.End, body.Place, body.Type, body.Description, body.Participants, body.Private).
			// Create(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Search handles the activity search request.
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			Query string `json:"query"`
		})

		resp := new(struct {
			Data struct {
				Activities []map[string]interface{} `json:"activities"`
			} `json:"data"`
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		activities, err := activity.Search(body.Query)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		resp.Data.Activities = activities.Actfilter()
		c.JSON(http.StatusOK, resp)
	}
}

// Update handles the activity update request.
func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			activity.Activity
			Update map[string]interface{} `json:"update"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if id, exists := body.Update["_id"]; exists {
			objectID, err := primitive.ObjectIDFromHex(id.(string))
			if err != nil {
				resp.Error = err.Error()
				c.JSON(http.StatusBadRequest, resp)
				return
			}
			body.Update["_id"] = objectID

		} else {
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := body.Activity.Update(body.Update); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Delete handles the activity deletion request.
func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ActivityID string `json:"_id"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		objectID, err := primitive.ObjectIDFromHex(body.ActivityID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := activity.Delete(objectID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
