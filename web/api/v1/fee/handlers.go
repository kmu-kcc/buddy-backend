// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
)

//Dones handles the inquiry of done submitted personel
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

//Yets handles the inquiry of yet submitted personel
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

//All handles the inquiry of all fee logs
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
