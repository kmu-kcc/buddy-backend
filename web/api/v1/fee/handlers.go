// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
)

// // Create handles the fee creation request.
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

// // Amount handles the submission amount request.
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

// //Dones handles the inquiry of done submitted personel
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

// //Yets handles the inquiry of yet submitted personel
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

// //All handles the inquiry of all fee logs
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

// Pay handles the payment request.
func Pay() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		body := new(struct {
			Year     int `json:"year"`
			Semester int `json:"semester"`
			Payments []struct {
				ID     string `json:"id"`
				Amount int    `json:"amount"`
			} `json:"payments"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		ids := make([]string, len(body.Payments))
		amounts := make([]int, len(body.Payments))

		for idx, payment := range body.Payments {
			ids[idx], amounts[idx] = payment.ID, payment.Amount
		}

		if err := fee.Pay(body.Year, body.Semester, ids, amounts); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Deposit handles the deposit request.
func Deposit() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		body := new(struct {
			Year     int `json:"year"`
			Semester int `json:"semester"`
			Amount   int `json:"amount"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := fee.Deposit(body.Year, body.Semester, body.Amount); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
