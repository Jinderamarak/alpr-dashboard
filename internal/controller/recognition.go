package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"net/http"
)

type RecognitionController struct {
	recognitions service.RecognitionService
}

func NewRecognitionController(recognitions service.RecognitionService) RecognitionController {
	return RecognitionController{recognitions}
}

func (controller *RecognitionController) Route(routes *gin.RouterGroup) {
	routes.GET("/:id", controller.GetRecognition)
	routes.POST("/", controller.PostRecognition)
}

func (controller *RecognitionController) GetRecognition(ctx *gin.Context) {
	recognitionId := ctx.Param("id")
	recognitionUuid, _ := uuid.Parse(recognitionId)

	recognition, _ := controller.recognitions.GetByIdWithCar(recognitionUuid)

	ctx.HTML(http.StatusOK, "recognition.html", gin.H{
		"recognized": recognition.CreatedAt,
		"car":        recognition.Car,
	})
}

func (controller *RecognitionController) PostRecognition(ctx *gin.Context) {
	plate := ctx.PostForm("plate")
	if len(plate) < 3 {
		ctx.String(http.StatusBadRequest, "Invalid license plate, minimum length is 3.")
		return
	}

	_, _ = controller.recognitions.CreateByPlate(plate)
	ctx.Redirect(http.StatusSeeOther, "/")
}
