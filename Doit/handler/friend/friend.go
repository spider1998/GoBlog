package friend

import (
	"Project/Doit/entity"
	"Project/Doit/service"
	"github.com/go-ozzo/ozzo-routing"
)

func QueryUsers(c *routing.Context) (err error) {
	var req entity.QueryUserRequest
	req.Name = c.Param("name")
	gender := c.Param("gender")
	req.Tag = c.Param("tags")
	if gender == "1" {
		req.Gender = entity.UserGenderMale
	} else {
		req.Gender = entity.UserGenderFemale
	}
	persons, err := service.Friend.QueryUsers(req)
	if err != nil {
		return
	}
	if len(persons) == 0 {
		persons = []entity.BaseUser{}
	}
	return c.Write(persons)
}
