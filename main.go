package main

import (
	"bn-crud-ads/configuration"
	ds "bn-crud-ads/domain/datasources"
	repo "bn-crud-ads/domain/repositories"
	gw "bn-crud-ads/src/gateways"
	"bn-crud-ads/src/middlewares"
	sv "bn-crud-ads/src/services"
	"bn-crud-ads/utils/providers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/watchakorn-18k/scalar-go"
)

func main() {

	// // // remove this before deploy ###################
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// /// ############################################

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
	middlewares.Logger(app)
	app.Use(recover.New())
	app.Use(cors.New())

	mongodb := ds.NewMongoDB(10)
	redisdb := ds.NewRedisConnection()

	adsMongo := repo.NewAdsRepository(mongodb)
	ads_redis := repo.NewRedisRepository(redisdb)
	wrongMessageRepo := repo.NewWrongMessageRepository(mongodb)

	sv0 := sv.NewAdsService(adsMongo, ads_redis, providers.NewS3Provider(), wrongMessageRepo)

	gw.NewHTTPGateway(app, sv0)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "7300"
	}

	app.Listen(":" + PORT)
}
