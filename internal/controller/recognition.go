package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"net/http"
	"strconv"
)

type RecognitionController struct {
	recognitions service.RecognitionService
}

func NewRecognitionController(recognitions service.RecognitionService) RecognitionController {
	return RecognitionController{recognitions}
}

func (controller *RecognitionController) Route(routes *gin.RouterGroup) {
	routes.GET("/", controller.GetList)
	routes.GET("/:id", controller.GetRecognition)
	routes.GET("/:id/upload", controller.GetUploadForm)
	routes.POST("/", controller.PostRecognition)
}

func (controller *RecognitionController) GetList(ctx *gin.Context) {
	pageQuery := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageQuery)

	pages, _ := controller.recognitions.CountPages()
	recognitions, _ := controller.recognitions.GetPage(page)

	ctx.HTML(http.StatusOK, "recognition/list", gin.H{
		"recognitions": recognitions,
		"pages":        pages,
	})
}

func (controller *RecognitionController) GetRecognition(ctx *gin.Context) {
	recognitionId := ctx.Param("id")
	recognitionUuid, _ := uuid.Parse(recognitionId)

	recognition, _ := controller.recognitions.GetByIdWithCar(recognitionUuid)
	photos, _ := controller.recognitions.ImagesByRecognitionId(&recognitionUuid)

	ctx.HTML(http.StatusOK, "recognition/event", gin.H{
		"id":         recognition.ID,
		"recognized": recognition.CreatedAt,
		"car":        recognition.Car,
		"photos":     photos,
	})
}

func (controller *RecognitionController) GetUploadForm(ctx *gin.Context) {
	recognitionId := ctx.Param("id")
	recognitionUuid, _ := uuid.Parse(recognitionId)

	_, _ = controller.recognitions.GetByIdWithCar(recognitionUuid)
	presigned, form, _ := controller.recognitions.CreateImageUpload(recognitionUuid)

	ctx.HTML(http.StatusOK, "recognition/upload", gin.H{
		"url":  presigned,
		"form": form,
	})
}

func (controller *RecognitionController) PostRecognition(ctx *gin.Context) {
	plate := ctx.PostForm("plate")

	_, err := controller.recognitions.CreateByPlate(plate)
	if err != nil {
		if errors.Is(err, service.ErrPlateTooShort) {
			ctx.HTML(http.StatusOK, "recognition/creation", gin.H{
				"error": "Plate is too short, minimum length is 3 characters.",
			})
			return
		}
		ctx.HTML(http.StatusOK, "recognition/creation", gin.H{
			"error": "Unknown error",
		})
		return
	}

	ctx.Header("HX-Trigger", "recognition-event-created")
	ctx.HTML(http.StatusOK, "recognition/creation", gin.H{})
}
