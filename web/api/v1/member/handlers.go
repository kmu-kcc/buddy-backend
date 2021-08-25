// Package member defines the router layer of the club member of the Buddy System.
package member

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"github.com/kmu-kcc/buddy-backend/pkg/oauth2"
)

// SignIn handles the signin request.
func SignIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(member.Member)
		resp := new(struct {
			Data struct {
				AccessToken oauth2.Token `json:"access_token"`
				ExpiredAt   int64        `json:"expired_at,string"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		err := body.SingIn()
		if err != nil {
			resp.Error = err.Error()
			if err == member.ErrIdentityMismatch {
				c.JSON(http.StatusConflict, resp)
			} else {
				c.JSON(http.StatusInternalServerError, resp)
			}
			return
		}

		resp.Data.AccessToken, resp.Data.ExpiredAt, err = oauth2.NewToken(body.ID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnprocessableEntity, resp)
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

// SignUp handles the signup request.
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		body := new(member.Member)
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err := member.New(body.ID, body.Name, body.Department, body.Phone, body.Email, body.Grade, body.Attendance).
			SignUp(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// SignUps handles the signup list request.
func SignUps() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		resp := new(struct {
			Data struct {
				SignUps member.Members `json:"signups"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		err := token.Valid()
		if err != nil {
			resp.Error = err.Error()
			resp.Data.SignUps = member.Members{}
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			resp.Data.SignUps = member.Members{}
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.MemberManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			resp.Data.SignUps = member.Members{}
			c.JSON(http.StatusForbidden, resp)
			return
		}

		resp.Data.SignUps, err = member.SignUps()
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Approve handles the signup approvement request.
func Approve() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			IDs []string `json:"ids"`
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
		} else if !role.MemberManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if err := member.Approve(body.IDs); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Delete handles the refusal and deletion request.
func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			IDs []string `json:"ids"`
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
		} else if !role.MemberManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if err := member.Delete(body.IDs); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Exit handles the exit request.
func Exit() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(member.Member)
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

		if err := body.Exit(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Exits handles the exit list request.
func Exits() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		resp := new(struct {
			Data struct {
				Exits member.Members `json:"exits"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			resp.Data.Exits = member.Members{}
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			resp.Data.Exits = member.Members{}
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.MemberManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			resp.Data.Exits = member.Members{}
			c.JSON(http.StatusForbidden, resp)
			return
		}

		var err error
		resp.Data.Exits, err = member.Exits()
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// My handles the personal information request.
func My() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			ID       string `json:"id"`
			Password string `json:"password"`
		})
		resp := new(struct {
			Data struct {
				Data map[string]interface{} `json:"data"`
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

		if resp.Data.Data, err = (&member.Member{ID: body.ID, Password: body.Password}).My(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Search handles the member search request.
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		query := c.Query("query")
		resp := new(struct {
			Data struct {
				Members []map[string]interface{} `json:"members"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		members, err := member.Search(query)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		resp.Data.Members = members.Public()
		c.JSON(http.StatusOK, resp)
	}
}

// Update handles the member update request.
func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			ID     string                 `json:"id"`
			Update map[string]interface{} `json:"update"`
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

		if err := (member.Member{ID: body.ID}).Update(body.Update); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Active handles the member signup activation status request.
func Active() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		resp := new(struct {
			Data struct {
				Active bool `json:"active"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		var err error
		if resp.Data.Active, err = member.Active(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Activate handles the member signup activation status update request.
func Activate() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			Activate bool `json:"activate"`
		})
		resp := new(struct {
			Data struct {
				Active bool `json:"active"`
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
		} else if !role.MemberManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if resp.Data.Active, err = member.Activate(body.Activate); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Graduates handles the graduate list request.
func Graduates() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		resp := new(struct {
			Data struct {
				Graduates member.Members `json:"graduates"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})

		if err := token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.MemberManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		var err error
		resp.Data.Graduates, err = member.Graduates()
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// UpdateRole handles the role update request.
func UpdateRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Request.Body.Close()

		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			ID   string      `json:"id"`
			Role member.Role `json:"role"`
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
		} else if !role.Master {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if err := member.UpdateRole(body.ID, body.Role); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}
