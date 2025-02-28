package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinderamarak/alpr-dasboard/internal/controller"
	"github.com/jinderamarak/alpr-dasboard/internal/data"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"github.com/jinderamarak/alpr-dasboard/templates"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"path/filepath"
)

func main() {
	s3, err := minio.New("localhost:9000", &minio.Options{
		Creds: credentials.NewStaticV4("UqpPskfiWz8RWzMo7hcO", "gF9OKGEDDLEPASbDyrS6SpmDsxovuoFBPvC3RxVp", ""),
	})

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		errors.Unwrap(err)
	}

	db.AutoMigrate(&model.Car{}, &model.Recognition{})

	carRepo := data.NewCarRepository(db)
	recognitionRepo := data.NewRecognitionRepository(db)
	vignetteProvider := data.NewEDalniceVignetteProvider()
	photoRepo := data.NewPhotoRepository(db, s3)

	carService := service.NewCarService(carRepo)
	recognitionService := service.NewRecognitionService(recognitionRepo, carService)
	vignetteService := service.NewVignetteService(vignetteProvider)

	funcMap := template.FuncMap{
		"seq":      templates.Sequence,
		"formatDT": templates.FormatDateTime,
	}
	notificationTemplates := loadTemplates(funcMap, "templates/notification/", "*.tmpl")

	indexController := controller.NewIndexController(recognitionService)
	carController := controller.NewCarController(carService, recognitionService, vignetteService)
	recognitionController := controller.NewRecognitionController(recognitionService)
	notificationController := controller.NewNotificationController(notificationTemplates, recognitionService, carService)

	server := gin.Default()
	server.SetHTMLTemplate(loadTemplates(funcMap, "templates/", "*.tmpl"))

	indexController.Route(server.Group("/"))
	carController.Route(server.Group("/car"))
	recognitionController.Route(server.Group("/recognition"))
	notificationController.Route(server.Group("/"))

	server.Run("localhost:8080")
}

func loadTemplates(funcMap template.FuncMap, folder, fileEnding string) *template.Template {
	files := make([]string, 0)
	err := includeSubfolders(&files, folder, fileEnding)
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
