package gateways

import (
	service "bn-crud-ads/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	AdsService service.IAdsService
}

func NewHTTPGateway(app *fiber.App, ads service.IAdsService) {
	gateway := &HTTPGateway{
		AdsService: ads,
	}
	GatewayAds(*gateway, app)
}
