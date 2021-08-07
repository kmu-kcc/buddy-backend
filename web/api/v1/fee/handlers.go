package fee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
