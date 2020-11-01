package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/Wilder60/KeyRing/internal/domain"
	"github.com/Wilder60/KeyRing/internal/security"
	"github.com/Wilder60/KeyRing/internal/sql"
	"github.com/gin-gonic/gin"
)

type KeyRing struct {
	logger *zap.Logger
	*sql.SQL
	*security.Security
}

func ProvideKeyRing(db *sql.SQL, sec *security.Security, log *zap.Logger) *KeyRing {
	return &KeyRing{
		log,
		db,
		sec,
	}
}

func (kr *KeyRing) getEvents(ctx *gin.Context) {
	queryValues := ctx.Request.URL.Query()
	limitStr := queryValues.Get("limit")
	offsetStr := queryValues.Get("offset")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var offset int64
	if offsetStr != "" {
		offset, err = strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}

	userID, err := kr.getUserID(ctx)
	if err != nil {
		kr.logger.Info(fmt.Sprintf("Encountered error deseralizing token %v", err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	entries, err := kr.GetKeyRing(userID, limit, offset)
	if err != nil {
		kr.logger.Warn(fmt.Sprintf("Error encountered while accessing database: %v", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, entries)
}

func (kr *KeyRing) addEvent(ctx *gin.Context) {
	var event domain.KeyEntry
	err := ctx.BindJSON(&event)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	userID, err := kr.getUserID(ctx)
	if err != nil {
		kr.logger.Info(fmt.Sprintf("Encountered error deseralizing token %v", err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	_, err = kr.AddKeyRing(event, userID)
	if err != nil {
		kr.logger.Warn(fmt.Sprintf("Error encountered while accessing database: %v", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (kr *KeyRing) updateEvent(ctx *gin.Context) {
	var event domain.KeyEntry
	err := ctx.BindJSON(&event)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := kr.getUserID(ctx)
	if err != nil {
		kr.logger.Info(fmt.Sprintf("Encountered error deseralizing token %v", err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	_, err = kr.UpdateKeyRing(event, userID)
	if err != nil {
		kr.logger.Warn(fmt.Sprintf("Error encountered while accessing database: %v", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (kr *KeyRing) deleteEvents(ctx *gin.Context) {
	var eventIds []string
	err := ctx.BindJSON(eventIds)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := kr.getUserID(ctx)
	if err != nil {
		kr.logger.Info(fmt.Sprintf("Encountered error deseralizing token %v", err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	rowsDel, err := kr.DeleteKeyRing(eventIds, userID)
	if err != nil {
		kr.logger.Warn(fmt.Sprintf("Error encountered while accessing database: %v", err))
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}
	ctx.JSON(http.StatusOK, rowsDel)
}

func (kr *KeyRing) decodeRequest(body io.ReadCloser, out interface{}) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(&out)
}

func (kr *KeyRing) getUserID(ctx *gin.Context) (string, error) {
	str := ctx.GetHeader("Authorization")
	splitString := strings.Split(str, " ")
	if len(splitString) != 2 {
		return "", errors.New("Malformed Token")
	}

	userID := kr.GetUserFromToken(splitString[1])
	return userID, nil
}

var KeyRingModule = fx.Option(
	fx.Provide(ProvideKeyRing),
)
