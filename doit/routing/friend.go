package routing

import (
	"Project/doit/handler/friend"
	"github.com/go-ozzo/ozzo-routing"
)

func FriendRegisterRoutes(router *routing.RouteGroup) {
	router.Get("/users", friend.QueryUsers)           //查询用户
	router.Post("/friends", friend.AddFriends)        //添加好友申请
	router.Patch("/friends", friend.AddAuthorization) //好友申请授权
	router.Delete("/friends", friend.DeleteFriend)    //删除好友
	router.Get("/friends", friend.GetFriendList)      //获取好友列表
	router.Patch("/blacklists", friend.PullBlack)     //拉黑好友/解除拉黑

}
