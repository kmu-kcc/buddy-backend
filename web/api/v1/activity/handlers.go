// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Search handles the activity search request.
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			Filter map[string]interface{} `json:"filter"`
		})

		resp := new(struct {
			Activities activity.Activities `json:"activities"`
			Error      string              `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		activities, err := activity.Search(body.Filter)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Activities = activities
		c.JSON(http.StatusOK, resp)
	}
}

// Update handles the activity update request.
func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
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

		if err := activity.Update(body.Update); err != nil {
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

// Participants handles the participant list request.
func Participants() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			ActivityID string `json:"_id"`
		})

		resp := new(struct {
			Members member.Members `json:"members"`
			Error   string         `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		activityID, err := primitive.ObjectIDFromHex(body.ActivityID)
		if err != nil {
			return
		}

		members, err := activity.Participants(activityID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Members = members
		c.JSON(http.StatusOK, resp)
	}
}

//ApplyP handles the activity apply request
func ApplyP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})
		body := new(
			struct {
				ActivityID string `json:"_id"`
				ID         string `json:"member_id"`
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
		if err := activity.ApplyP(activityID, body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
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
			Data struct {
				Papplies []map[string]interface{} `json:"papplies"`
			} `json:"data"`
			Error string `json:"error"`
		})
		body := new(
			struct {
				ActivityID string `json:"_id"`
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
		members, err := activity.Papplies(activityID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Data.Papplies = members.Memfilter()
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
				ActivityID string   `json:"_id"`
				IDs        []string `json:"member_ids"`
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
		if err := activity.ApproveP(activityID, body.IDs); err != nil {
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
				ActivityID string   `json:"_id"`
				IDs        []string `json:"member_ids"`
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
		if err := activity.RejectP(activityID, body.IDs); err != nil {
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
				ActivityID string `json:"_id"`
				ID         string `json:"member_id"`
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
		if err := activity.CancelP(activityID, body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

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
