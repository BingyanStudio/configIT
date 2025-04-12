package controller

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-.]+$`)

type OKResp struct {
	Data interface{} `json:"data"`
}

type ErrorResp struct {
	Message string `json:"msg"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &OKResp{
		Data: data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, &ErrorResp{
		Message: msg,
	})
}

func ErrInternal(c *gin.Context, err error) {
	// log.Errorw("internal error", "err", err)
	Error(c, http.StatusInternalServerError, "坏啦，怎么炸了")
}

func ErrNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, &ErrorResp{
		Message: "坏啦，没找到呢",
	})
}

func ErrBadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, &ErrorResp{
		Message: "坏啦，你说的啥",
	})
}

func ErrUnauthorize(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, &ErrorResp{
		Message: "坏啦，你是谁呢",
	})
}

func ErrNotFoundOrInternal(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ErrNotFound(c)
		return
	}
	ErrInternal(c, err)
}
