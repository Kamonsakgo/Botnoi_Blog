package gateways

import (
	"bn-crud-ads/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func GatewayAds(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/ads")

	apiNoJWT := api.Group("")
	apiNoJWT.Get("/set_ads", gateway.SetAdsRedis)
	apiNoJWT.Get("/get_ads_none_token", gateway.GetAdsData)

	apiJWT := api.Group("")
	apiJWT.Use(middlewares.SetBotnoiJWtHeaderHandler())
	apiJWT.Get("/get_ads", gateway.GetAdsData)
	apiJWT.Post("/update_ads", gateway.UpdateAds)
	apiJWT.Get("/get_marketplace_sound", gateway.GetMarketplaceSound)
	apiJWT.Post("/insert_ads", gateway.Insert_ads)
	apiJWT.Delete("/delete_ads", gateway.Delete_ads)
	apiJWT.Put("/update", gateway.Update)
	apiJWT.Get("/getall_ads", gateway.Getall_ads)
	apiJWT.Get("/find_one", gateway.Find_one_ads)
	apiJWT.Get("/find_one_random", gateway.Find_one_random_ads)

}
