package auth

import (
	"github.com/ValeryVerkhoturov/chat/config"
	"github.com/ValeryVerkhoturov/chat/db"
	"github.com/ValeryVerkhoturov/chat/utils/requestUtils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SessionRequired is a simple middleware to check the session.
func SessionRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(config.SessionUserKey)
	if user == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, requestUtils.GetApplicationError("unauthorized"))
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

// CreateSessionIfNotExists is a middleware that automatically create user
// if session userId does not exist or userId was not found id bd,
// when chat is loaded.
func CreateSessionIfNotExists(c *gin.Context) {
	session := sessions.Default(c)
	userId := session.Get(config.SessionUserKey)
	if userId == nil {
		createSession(c)
	} else {
		// Check if user exists
		var objUserId, err = db.ConvertValueToObjectId(userId)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				requestUtils.GetApplicationError("Failed to get user id"),
			)
		}
		userQuery := db.User{ID: objUserId}
		if _, err = userQuery.FindOneById(); err != nil {
			// Create user and set session
			createSession(c)
		}
	}
	// Continue down the chain to handler etc
	c.Next()
}

// createSession creates user and sets session
func createSession(c *gin.Context) {
	session := sessions.Default(c)
	newUser := db.User{Source: db.EmbeddedChat, IsAgent: false}

	insertResult, err := newUser.InsertOne()
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			requestUtils.GetApplicationError("Failed to create user"),
		)
	}
	objectId, ok := db.ConvertInsertOneResultToId(insertResult)
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			requestUtils.GetApplicationError("Failed to get created user id"),
		)
	}

	session.Set(config.SessionUserKey, objectId)
	if err := session.Save(); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			requestUtils.GetApplicationError("Failed to save session"),
		)
		return
	}
}
