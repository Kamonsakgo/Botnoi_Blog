package main

import (
	"go-fiber-unittest/configuration"
	ds "go-fiber-unittest/domain/datasources"
	repo "go-fiber-unittest/domain/repositories"
	gw "go-fiber-unittest/src/gateways"
	"go-fiber-unittest/src/middlewares"
	sv "go-fiber-unittest/src/services"
	"go-fiber-unittest/src/utils/providers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/watchakorn-18k/scalar-go"
)

func main() {

	// // remove this before deploy ###################
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	/// ############################################
	app := fiber.New(configuration.NewFiberConfiguration())

	app.Use("/api/admin/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "bn-crud-admin API",
			},
			Theme:    "purple",
			Layout:   "modern",
			DarkMode: true,
		})
		if err != nil {
			return err
		}
		c.Type("html")
		return c.SendString(htmlContent)
	})
	app.Use(middlewares.NewLogger())
	app.Use(recover.New())
	app.Use(cors.New())

	mongodb := ds.NewMongoDB(10)

	userMongo := repo.NewUsersRepository(mongodb)
	blogMongo := repo.NewBlogsRepository(mongodb)
	highlighteventMongo := repo.NewHighlightsRepository(mongodb)
	sv0 := sv.NewUsersService(userMongo)
	blog := sv.NewBlogsService(blogMongo, userMongo, providers.NewS3Provider())
	highlightevent := sv.NewHighlightsService(highlighteventMongo, userMongo, providers.NewS3Provider())

	gw.NewHTTPGateway(app, sv0, blog, highlightevent)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}
