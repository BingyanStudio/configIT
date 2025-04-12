package client

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/BingyanStudio/configIT/internal/controller"
	"github.com/BingyanStudio/configIT/internal/model"
	"github.com/BingyanStudio/configIT/internal/utils"
)

type ConfigUri struct {
	AppName   string `uri:"app_name" binding:"required"`
	ConfigKey string `uri:"config_key" binding:"required"`
}

type GetConfigResp struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
	TS    int64  `json:"ts"`
}

func GetConfig(c *gin.Context) {
	var (
		uri ConfigUri
		v   any
		err error
	)
	if err = c.ShouldBindUri(&uri); err != nil {
		controller.ErrBadRequest(c)
		return
	}

	config, err := model.GetConfig(c, uri.AppName, uri.ConfigKey)
	if err != nil {
		controller.ErrNotFoundOrInternal(c, err)
		return
	}

	switch config.Type {
	case model.ConfigTypeString:
		v = config.Value
	case model.ConfigTypeNumber:
		v, err = strconv.Atoi(config.Value)
		if err != nil {
			v, err = strconv.ParseFloat(config.Value, 64)
		}
	case model.ConfigTypeBool:
		v, err = strconv.ParseBool(config.Value)
	default:
		err = json.Unmarshal([]byte(config.Value), &v)
	}

	if err != nil {
		controller.ErrBadRequest(c)
		return
	}

	resp := GetConfigResp{
		Type:  string(config.Type),
		Value: v,
		TS:    config.UpdatedAt.Unix(),
	}

	controller.OK(c, resp)
}

func UpdateConfig(c *gin.Context) {
	var (
		uri ConfigUri
		v   = c.PostForm("value")
		err error
	)
	if err = c.ShouldBindUri(&uri); err != nil {
		controller.ErrBadRequest(c)
		return
	}

	config, err := model.GetConfig(c, uri.AppName, uri.ConfigKey)
	if err != nil {
		controller.ErrNotFoundOrInternal(c, err)
		return
	}

	// Validate the value based on the config type
	switch config.Type {
	case model.ConfigTypeString:
		config.Value = c.PostForm("value")
	case model.ConfigTypeNumber:
		_, err = strconv.Atoi(v)
		if err != nil {
			_, err = strconv.ParseFloat(v, 64)
		}
	case model.ConfigTypeBool:
		_, err = strconv.ParseBool(v)
	default:
		if !json.Valid([]byte(v)) {
			err = json.Unmarshal([]byte(v), &struct{}{})
		}
	}
	if err != nil {
		controller.ErrBadRequest(c)
		return
	}
	err = model.UpdateConfig(c, config)
	if err != nil {
		controller.ErrNotFoundOrInternal(c, err)
		return
	}

	//Trigger the refresh hook
	switch config.App.AutoRefresh {
	case model.AutoRefreshConfigMap:
		ok, err := config.App.Namespace.InCluster(c)
		if err != nil {
			controller.ErrInternal(c, err)
			return
		}
		if !ok {
			controller.ErrNotFound(c)
			return
		}
		cm, err := utils.GetConfigMap(c, config.App.Namespace.Name, config.App.Name)
		if err != nil {
			controller.ErrInternal(c, err)
			return
		}
		cm[config.Key] = config.Value
		err = utils.UpdateConfigMap(c, config.App.Namespace.Name, config.App.Name, cm)
		if err != nil {
			controller.ErrInternal(c, err)
			return
		}
	case model.AutoRefreshHook:
		err = utils.SendHook(config.App.AutoRefreshParam, config.Key)
		if err != nil {
			controller.ErrInternal(c, err)
			return
		}
	case model.AutoRefreshPassive:
		// Do nothing
	}
	controller.OK(c, nil)
}
