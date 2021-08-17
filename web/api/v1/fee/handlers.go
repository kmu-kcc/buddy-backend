// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create handles the fee creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(fee.Fee)
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		if err := body.Create(); err != nil {
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
			Data struct {
				Sum int `json:"sum"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
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

		resp.Data.Sum = sum
		c.JSON(http.StatusOK, resp)
	}
}

// Payers handles the payer list request.
func Payers() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				Payers member.Members `json:"payers"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})
		body := new(fee.Fee)

		err := json.NewDecoder(c.Request.Body).Decode(body)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if resp.Data.Payers, err = body.Payers(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		resp.Data.Dones = res.Public()
		c.JSON(http.StatusOK, resp)
	}
}

// Deptors handles deptor list request.
func Deptors() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				Deptors []struct {
					member.Member
					Dept int `json:"dept"`
				} `json:"deptors"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		body := new(fee.Fee)

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		deptors, depts, err := body.Deptors()
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		resp.Data.Deptors = make([]struct {
			member.Member
			Dept int `json:"dept"`
		}, len(deptors))

		for idx, deptor := range deptors {
			resp.Data.Deptors[idx].Member = deptor
			resp.Data.Deptors[idx].Dept = depts[idx]
		}

    c.JSON(http.StatusOK, resp)
	}
}

// Search handles the fee search request.
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				CarryOver int      `json:"carry_over"`
				Logs      fee.Logs `json:"logs"`
				Total     int      `json:"total"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		body := new(fee.Fee)

		err := json.NewDecoder(c.Request.Body).Decode(body)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if resp.Data.CarryOver, resp.Data.Logs, resp.Data.Total, err = body.Search(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
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
			IDs []string `json:"ids"`
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

func Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
		})

		body := new(struct {
			Year     string `json:"year"`
			Semester string `json:"semester"`
			Amount   string `json:"amount"`
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
