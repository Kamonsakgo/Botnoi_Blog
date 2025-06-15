package gateways

import (
	"go-fiber-unittest/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func GatewayUsers(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/v1/users")

	api.Post("/add_user", gateway.CreateNewUserAccount)
	api.Get("/users", gateway.GetAllUserData)

	apiblog := app.Group("/api/blog")
	apiblog.Get("/get_all_blog", gateway.GetAllblog)
	apiblog.Get("/get_blog", gateway.GetOneBlog)

	blogJwt := app.Group("/api/blog")
	blogJwt.Use(middlewares.SetBotnoiJWtHeaderHandler())
	blogJwt.Post("/insert_blog", gateway.Insertblog)
	blogJwt.Put("/update_blog", gateway.Updateblog)
	blogJwt.Delete("/delete_blog", gateway.Deleteblog)
	blogJwt.Post("/upload_image", gateway.UploadImage)

	HL_event := app.Group("/api/highlightevent")
	HL_event.Get("/get_all", gateway.GetAllHL_event)
	HL_event.Get("/get_one", gateway.GetOneHL_event)

	HL_eventJwt := app.Group("/api/highlightevent")
	HL_eventJwt.Use(middlewares.SetBotnoiJWtHeaderHandler())
	HL_eventJwt.Post("/insert", gateway.InsertHL_event)
	HL_eventJwt.Put("/update", gateway.UpdateHL_event)
	HL_eventJwt.Delete("/delete", gateway.DeleteHL_event)

}
