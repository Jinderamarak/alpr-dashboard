package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"net/http"
	"strconv"
)

type IndexController struct {
	recognitions service.RecognitionService
}

func NewIndexController(recognitions service.RecognitionService) IndexController {
	return IndexController{recognitions}
}

func (controller *IndexController) Route(routes *gin.RouterGroup) {
	routes.GET("/", controller.GetIndex)
}

func (controller *IndexController) GetIndex(ctx *gin.Context) {
	pageQuery := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageQuery)

	pages, _ := controller.recognitions.CountPages()
	recognitions, _ := controller.recognitions.GetPage(page)

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"recognitions": recognitions,
		"pages":        pages,
	})
}
