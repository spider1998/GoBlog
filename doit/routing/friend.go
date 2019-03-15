package routing

import (
	"Project/doit/handler/friend"
	"github.com/go-ozzo/ozzo-routing"
)

func FriendRegisterRoutes(router *routing.RouteGroup) {
	router.Get("/users", friend.QueryUsers) //查询用户
	router.Post("/add", friend.AddFriends)  //添加好友申请
	router.Post("/authorization",friend.AddAuthorization)

}
