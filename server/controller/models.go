package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"wechat-gptbot/config"
)

/*
* @Author: zouyx
* @Email:
* @Date:   2024/5/29 18:06
* @Package: 模型管理
 */
type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type CurrentModelResponse struct {
	BaseResponse
	Data ModelInfo `json:"data"`
}

// CurrentModel 获取当前模型
func CurrentModel(c *gin.Context) {
	c.JSON(http.StatusOK, CurrentModelResponse{
		BaseResponse: BaseResponse{Code: 200,
			Msg: "ok"},
		Data: ModelInfo{config.C.GetBaseModel(), openai.CreateImageModelDallE3}})
}

type ModelInfo struct {
	TextModel    string `json:"text_model"`
	DrawingModel string `json:"drawing_model"`
}

func ResetModel(c *gin.Context) {
	info := ModelInfo{}
	c.ShouldBindBodyWithJSON(&info)
	if info.TextModel != "" {
		config.C.SetBaseModel(info.TextModel)
	}
	c.JSON(http.StatusOK, BaseResponse{
		Code: 200,
		Msg:  "ok",
	})
}
