package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"net/http"
	"strconv"
)

type CarController struct {
	cars         service.CarService
	recognitions service.RecognitionService
}

func NewCarController(cars service.CarService, recognitions service.RecognitionService) CarController {
	return CarController{cars, recognitions}
}

func (controller *CarController) Route(routes *gin.RouterGroup) {
	routes.GET("/:id", controller.GetCar)
}

func (controller *CarController) GetCar(ctx *gin.Context) {
	carId := ctx.Param("id")
	carUuid, _ := uuid.Parse(carId)

	pageQuery := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageQuery)

	pages, _ := controller.recognitions.CountPagesWithCarId(&carUuid)
	car, _ := controller.cars.GetById(carUuid)
	recognitions, _ := controller.recognitions.GetPageWithCarId(&carUuid, page)

	ctx.HTML(http.StatusOK, "car.html", gin.H{
		"carId":        car.ID,
		"plate":        car.Plate,
		"description":  car.Description,
		"authorized":   car.IsAuthorized,
		"pages":        pages,
		"recognitions": recognitions,
	})
}
