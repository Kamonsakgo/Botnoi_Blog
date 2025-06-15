package gateways

import (
	service "go-fiber-unittest/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	UserService      service.IUsersService
	BlogService      service.IBlogsService
	HighlightService service.IHighlightsService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, blogs service.IBlogsService, Highlights service.IHighlightsService) {
	gateway := &HTTPGateway{
		UserService:      users,
		BlogService:      blogs,
		HighlightService: Highlights,
	}
	GatewayUsers(*gateway, app)
	//RouteBlog(*gateway, app)
}
