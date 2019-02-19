package friend

import (
	"github.com/go-ozzo/ozzo-routing"
	"Project/doit/entity"
	"Project/doit/service"
	"Project/doit/util"

)

func QueryUsers(c *routing.Context) (err error) {
	var req entity.QueryUserRequest
	req.Name = c.Param("name")
	gender := c.Param("gender")
	req.Tag = c.Param("tags")
	if gender == "1"{
		req.Gender = entity.UserGenderMale
	}else{
		req.Gender = entity.UserGenderFemale
	}
	persons,err := service.Friend.QueryUsers(req)
	if err != nil{
		return
	}
	if len(persons) == 0{
		persons = []entity.BaseUser{}
	}
	pager := util.GetPaginatedListFromRequest(c, len(persons))
	if pager.Offset()+pager.Limit() <= pager.TotalCount {
		return c.Write(persons[pager.Offset() : pager.Offset()+pager.Limit()])
	} else {
		return c.Write(persons[pager.Offset():pager.TotalCount])
	}
}
