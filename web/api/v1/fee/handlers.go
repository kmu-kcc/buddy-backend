// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create handles the fee creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(fee.Fee)
		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := fee.New(body.Year, body.Semester, body.Amount).Create(); err != nil {
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
			fee.Fee
			MemberID string `json:"member_id"`
		})

		resp := new(struct {
			Error string `json:"error"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := body.Submit(body.MemberID); err != nil {
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

		resp.Data.Sum = sum
		c.JSON(http.StatusOK, resp)
	}
}

//Dones handles the inquiry of done submitted personel
func Dones() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				Dones []map[string]interface{} `json:"dones"`
			} `json:"data"`
			Error string `json:"error"`
		})
		body := new(fee.Fee)

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		res, err := body.Dones()
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		resp.Data.Dones = res.Memfilter()
		c.JSON(http.StatusOK, resp)
	}
}

//Yets handles the inquiry of yet submitted personel
func Yets() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				Yets []map[string]interface{} `json:"yets"`
			} `json:"data"`
			Error string `json:"error"`
		})
		body := new(fee.Fee)
		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		res, err := body.Yets()
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Data.Yets = res.Memfilter()
		c.JSON(http.StatusOK, resp)
	}
}

//All handles the inquiry of all fee logs
func All() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				Logs []map[string]interface{} `json:"logs"`
			} `json:"data"`
			Error string `json:"error"`
		})
		body := new(struct {
			Startdate int `json:"startdate"`
			Enddate   int `json:"enddate"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		res, err := fee.All(body.Startdate, body.Enddate)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		resp.Data.Logs = res.Logfilter()
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

func Reject() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
			} `json:"data"`
			Error string `json:"error"`
		})

		body := new(struct {
			IDs []primitive.ObjectID `json:"ids"`
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
