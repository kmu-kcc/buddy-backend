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

// Package activity defines the router layer of the club activity of the Buddy System.
package activity

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kmu-kcc/buddy-backend/pkg/activity"
	"github.com/kmu-kcc/buddy-backend/pkg/member"
	"github.com/kmu-kcc/buddy-backend/pkg/oauth2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Create handles the activity creation request.
func Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(activity.Activity)
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
		} else if !role.ActivityManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if err := activity.New(body.Title, body.Start, body.End, body.Place, body.Description, body.Type, body.Participants, body.Private).
			Create(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Search handles the public activity search request.
func Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		resp := new(struct {
			Data struct {
				Activities activity.Activities `json:"activities"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})
		var err error

		resp.Data.Activities, err = activity.Search(query, false)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Private handles the private activity search request.
func Private() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		query := c.Query("query")
		resp := new(struct {
			Data struct {
				Activities activity.Activities `json:"activities"`
			} `json:"data"`
			Error string `json:"error,omitempty"`
		})
		var err error

		if err = token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !(role.MemberManagement || role.ActivityManagement || role.FeeManagement) {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		resp.Data.Activities, err = activity.Search(query, true)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}
		c.JSON(http.StatusOK, resp)
	}
}

// Update handles the activity update request.
func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			ID     string            `json:"id"`
			Update activity.Activity `json:"update"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})
		var err error

		if err = json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		if err = token.Valid(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusUnauthorized, resp)
			return
		}

		if role, err := token.Role(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		} else if !role.ActivityManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		body.Update.ID, err = primitive.ObjectIDFromHex(body.ID)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		if err = body.Update.Update(); err != nil {
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
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
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
		} else if !role.ActivityManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = activity.Delete(objectID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}

// Upload handles the file upload request.
func Upload() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		id := c.Query("id")
		resp := new(struct {
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
		} else if !role.ActivityManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		filename := filepath.Base(file.Filename)

		if err = c.SaveUploadedFile(file, activity.NewFile(filename).Absolute()); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(id); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = (activity.Activity{ID: objectID}).Upload(filename); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}

// Download handles the file download request.
func Download() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := new(struct {
			FileName string `json:"filename"`
		})
		resp := new(struct {
			Error string `json:"error,omitempty"`
		})

		if err := json.NewDecoder(c.Request.Body).Decode(body); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
		c.File(activity.NewFile(body.FileName).Absolute())
	}
}

// DeleteFile handles the file deletion request.
func DeleteFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := oauth2.Token(c.Request.Header.Get("Authorization"))
		body := new(struct {
			ID       string `json:"id"`
			FileName string `json:"filename"`
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
		} else if !role.ActivityManagement {
			resp.Error = member.ErrPermissionDenied.Error()
			c.JSON(http.StatusForbidden, resp)
			return
		}

		if objectID, err := primitive.ObjectIDFromHex(body.ID); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else if err = (activity.Activity{ID: objectID}).DeleteFile(body.FileName); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusInternalServerError, resp)
		} else {
			c.JSON(http.StatusOK, resp)
		}
	}
}
