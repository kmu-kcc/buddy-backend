// Copyright 2021 KMU KCC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fee defines the router layer of the club fee of the Buddy System.
package fee

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/fee"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"github.com/kmu-kcc/buddy-backend/pkg/oauth2"
)

// // Create handles the fee creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(fee.Fee)
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.FeeManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
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

// // Amount handles the submission amount request.
func Amount() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
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

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
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
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(fee.Fee)
		resp := new(struct {
			Data struct {
				Payers member.Members `json:"payers"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		err := json.NewDecoder(c.Request.Body).Decode(body)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.FeeManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if resp.Data.Payers, err = body.Payers(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// Deptors handles deptor list request.
func Deptors() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(fee.Fee)
		resp := new(struct {
			Data struct {
				Deptors []struct {
					member.Member
					Dept int `json:"dept"`
				} `json:"deptors"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.FeeManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
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
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(fee.Fee)
		resp := new(struct {
			Data struct {
				CarryOver int                      `json:"carry_over"`
				Logs      []map[string]interface{} `json:"logs"`
				Total     int                      `json:"total"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		err := json.NewDecoder(c.Request.Body).Decode(body)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
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

// Pay handles the payment request.
func Pay() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			Year     int `json:"year"`
			Semester int `json:"semester"`
			Payments []struct {
				ID     string `json:"id"`
				Amount int    `json:"amount"`
			} `json:"payments"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.FeeManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
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
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			Year        int    `json:"year"`
			Semester    int    `json:"semester"`
			Amount      int    `json:"amount"`
			Description string `json:"description"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.FeeManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if err := fee.Deposit(body.Year, body.Semester, body.Amount, body.Description); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Exempt handles the exemption request.
func Exempt() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			fee.Fee
			ID string `json:"id"`
		})

		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.FeeManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if err := body.Exempt(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
