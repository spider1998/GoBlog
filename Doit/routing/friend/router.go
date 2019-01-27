package friend

import (
	"Project/Doit/handler/friend"
	"github.com/go-ozzo/ozzo-routing"
)

func RegisterRoutes(router *routing.RouteGroup) {
	router.Get("/users", friend.QueryUsers) //查询用户
}
