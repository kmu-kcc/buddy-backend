// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"
  "strconv"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

// Create handles the fee creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			Year     int `json:"year"`
			Semester int `json:"semester"`
			Amount   int `json:"amount"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := fee.Create(body.Year, body.Semester, body.Amount); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Submit handles the fee submission request.
func Submit() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			MemberID string `json:"member_id"`
			Year     int    `json:"year"`
			Semester int    `json:"semester"`
			Amount   int    `json:"amount"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := fee.Submit(body.MemberID, body.Year, body.Semester, body.Amount); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Amount handles the submission amount request.
func Amount() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(struct {
			MemberID string `json:"member_id"`
			Year     int    `json:"year"`
			Semester int    `json:"semester"`
		})

		resp := new(struct {
			Sum   int    `json:"sum"`
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		sum, err := fee.Amount(body.Year, body.Semester, body.MemberID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Sum = sum
		c.JSON(http.StatusOK, resp)
	}
}

// Dones handles the fee submittee list request.
func Dones() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Dones member.Members `json:"dones"`
			Error string         `json:"error"`
		})
		body := new(
			struct {
				Year     int `json:"year"`
				Semester int `json:"semester"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		res, err := fee.Dones(body.Year, body.Semester)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Dones = res
		c.JSON(http.StatusOK, resp)
	}
}

// Yets handles the unsubmittee list request.
func Yets() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Yets  member.Members `json:"yets"`
			Error string         `json:"error"`
		})
		body := new(
			struct {
				Year     int `json:"year"`
				Semester int `json:"semester"`
			})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		res, err := fee.Yets(body.Year, body.Semester)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Yets = res
		c.JSON(http.StatusOK, resp)
	}
}

// All handles the fee log list request.
func All() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Logs  fee.Logs `json:"logs"`
			Error string   `json:"error"`
		})
		body := new(struct {
			Year     int `json:"year"`
			Semester int `json:"semester"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		res, err := fee.All(body.Year, body.Semester)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Logs = res
		c.JSON(http.StatusOK, resp)
	}
}

func Approve() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			IDs []string `json:ids`
		})

		err := json.NewDecoder(c.Request.Body).Decode(body)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		ids := make([]primitive.ObjectID, 10)

		for idx, id := range body.IDs {
			ids[idx], err = primitive.ObjectIDFromHex(id)
			if err != nil {
				resp.Error = err.Error()
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
		}

		if err := fee.Approve(ids); err != nil {
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func Reject() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			IDs []primitive.ObjectID `json:ids`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := fee.Reject(body.IDs); err != nil {
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

func Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			Year     string `json: "year, string"`
			Semester string `json: "semester, string"`
			Amount   string `json: "amount, string"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		year, err := strconv.Atoi(body.Year)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		semester, err := strconv.Atoi(body.Semester)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		amount, err := strconv.Atoi(body.Amount)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := fee.Deposit(year, semester, amount); err != nil {
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
