// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
)

//Dones handles the inquiry of done submitted personel
func Dones() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
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
		if _, err := fee.Dones(body.Year, body.Semester); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

//Yets handles the inquiry of yet submitted personel
func Yets() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
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
		if _, err := fee.Yets(body.Year, body.Semester); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

//All handles the inquiry of all fee logs
func All() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Error string `json:"error"`
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
		if _, err := fee.All(body.Year, body.Semester); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
