package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"net/http"
	"strconv"
	"time"
)

type CarController struct {
	cars         service.CarService
	recognitions service.RecognitionService
	vignettes    service.VignetteService
}

func NewCarController(cars service.CarService, recognitions service.RecognitionService, vignettes service.VignetteService) CarController {
	return CarController{cars, recognitions, vignettes}
}

func (controller *CarController) Route(routes *gin.RouterGroup) {
	routes.GET("/:id", controller.GetCar)
	routes.GET("/:id/edit", controller.EditCar)
	routes.PATCH("/:id", controller.UpdateCar)
	routes.GET("/:id/vignette", controller.GetVignette)
}

func (controller *CarController) GetCar(ctx *gin.Context) {
	carId := ctx.Param("id")
	carUuid, _ := uuid.Parse(carId)

	pageQuery := ctx.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageQuery)

	pages, _ := controller.recognitions.CountPagesWithCarId(&carUuid)
	car, _ := controller.cars.GetById(carUuid)
	recognitions, _ := controller.recognitions.GetPageWithCarId(&carUuid, page)

	ctx.HTML(http.StatusOK, "car/overview", gin.H{
		"carId":        car.ID,
		"plate":        car.Plate,
		"description":  car.Description,
		"authorized":   car.IsAuthorized,
		"pages":        pages,
		"recognitions": recognitions,
	})
}

func (controller *CarController) EditCar(ctx *gin.Context) {
	carId := ctx.Param("id")
	carUuid, _ := uuid.Parse(carId)

	car, _ := controller.cars.GetById(carUuid)

	ctx.HTML(http.StatusOK, "car/edit", gin.H{
		"car": car,
	})
}

func (controller *CarController) UpdateCar(ctx *gin.Context) {
	carId := ctx.Param("id")
	carUuid, _ := uuid.Parse(carId)

	authorized := ctx.PostForm("authorized") == "on"
	description := ctx.PostForm("description")

	_ = controller.cars.Update(carUuid, authorized, description)

	ctx.Header("HX-Redirect", "/car/"+carUuid.String())
	ctx.String(http.StatusOK, "ok")
}

func (controller *CarController) GetVignette(ctx *gin.Context) {
	carId := ctx.Param("id")
	carUuid, _ := uuid.Parse(carId)

	car, _ := controller.cars.GetById(carUuid)
	vignette, err := controller.vignettes.ValidatePlate(car.Plate)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "failed retrieving vignette for %s: %w", car.Plate, err)
		return
	}

	now := time.Now()
	valid := false
	for _, charge := range vignette.Charges {
		if charge.IsValidFor(now) {
			valid = true
			break
		}
	}

	ctx.HTML(http.StatusOK, "car/vignette", gin.H{
		"carId":   car.ID,
		"valid":   valid,
		"charges": vignette.Charges,
	})
}
