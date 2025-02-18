package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinderamarak/alpr-dasboard/internal/controller"
	"github.com/jinderamarak/alpr-dasboard/internal/data"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"github.com/jinderamarak/alpr-dasboard/templates"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"path/filepath"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		errors.Unwrap(err)
	}

	db.AutoMigrate(&model.Car{}, &model.Recognition{})

	carRepo := data.NewCarRepository(db)
	recRepo := data.NewRecognitionRepository(db)
	vigData := data.NewEDalniceVignetteProvider()

	carService := service.NewCarService(carRepo)
	recognitionService := service.NewRecognitionService(recRepo, carService)
	vignetteService := service.NewVignetteService(vigData)

	indexController := controller.NewIndexController(recognitionService)
	carController := controller.NewCarController(carService, recognitionService, vignetteService)
	recognitionController := controller.NewRecognitionController(recognitionService)

	server := gin.Default()
	server.SetHTMLTemplate(loadTemplates(template.FuncMap{
		"seq":      templates.Sequence,
		"formatDT": templates.FormatDateTime,
	}))

	indexController.Route(server.Group("/"))
	carController.Route(server.Group("/car"))
	recognitionController.Route(server.Group("/recognition"))

	server.Run("localhost:8080")
}

func loadTemplates(funcMap template.FuncMap) *template.Template {
	files := make([]string, 0)
	err := includeSubfolders(&files, "templates/", "*.tmpl")
	if err != nil {
		errors.Unwrap(err)
	}

	return template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
}

func includeSubfolders(files *[]string, folder, ending string) error {
	level, err := filepath.Glob(folder + ending)
	if err != nil {
		return err
	}
	if len(level) == 0 {
		return nil
	}

	*files = append(*files, level...)
	return includeSubfolders(files, folder+"*/", ending)
}
