package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jinderamarak/alpr-dasboard/internal/model"
	"github.com/jinderamarak/alpr-dasboard/internal/service"
	"github.com/jinderamarak/alpr-dasboard/internal/util"
	"html/template"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type NotificationController struct {
	templ        *template.Template
	recognitions chan model.Recognition
	cars         service.CarService
	publisher    *util.Broker[*string]
}

func NewNotificationController(templ *template.Template, recognitions service.RecognitionService, cars service.CarService) NotificationController {
	return NotificationController{templ, recognitions.Notifications().Subscribe(), cars, util.NewBroker[*string]()}
}

func (controller *NotificationController) Route(routes *gin.RouterGroup) {
	go controller.start()

	routes.GET("/ws", controller.handler)
}

func (controller *NotificationController) start() {
	go controller.publisher.Start()
	for {
		next := <-controller.recognitions
		if next.CarID != nil {
			car, _ := controller.cars.GetById(uuid.MustParse(*next.CarID))
			next.Car = &car
		}

		var builder strings.Builder
		err := controller.templ.ExecuteTemplate(&builder, "notification/recognition", gin.H{
			"event": next,
		})
		if err != nil {
			//	TODO: handle err
			continue
		}

		result := builder.String()
		controller.publisher.Publish(&result)
	}
}

func (controller *NotificationController) handler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	messages := controller.publisher.Subscribe()
	defer controller.publisher.Unsubscribe(messages)

	for {
		message := <-messages
		if message == nil {
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, []byte(*message))
		if err != nil {
			return
		}
	}
}
