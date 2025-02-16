package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinderamarak/alpr-dasboard/internal/controller"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/repository"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"github.com/jinderamarak/alpr-dasboard/templates"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Car{}, &model.Recognition{})

	carRepo := repository.NewCarRepository(db)
	recRepo := repository.NewRecognitionRepository(db)

	carService := service.NewCarService(carRepo)
	recognitionService := service.NewRecognitionService(recRepo, carService)

	indexController := controller.NewIndexController(recognitionService)
	carController := controller.NewCarController(carService, recognitionService)
	recognitionController := controller.NewRecognitionController(recognitionService)

	server := gin.Default()
	server.SetFuncMap(template.FuncMap{
		"seq":      templates.Sequence,
		"formatDT": templates.FormatDateTime,
	})
	server.LoadHTMLGlob("templates/*")

	indexController.Route(server.Group("/"))
	carController.Route(server.Group("/car"))
	recognitionController.Route(server.Group("/recognition"))

	server.Run("localhost:8080")
}
