package session

import (
	"Project/Doit/entity"
	"github.com/go-ozzo/ozzo-routing"
)

const (
	sessionKeyCO = "sess:co"
)

func SetUseression(c *routing.Context, user entity.User) {
	c.Set(sessionKeyCO, user)
}

func GetUserSession(c *routing.Context) entity.User {
	return c.Get(sessionKeyCO).(entity.User)
}

func GetSession(c *routing.Context) error {
	return c.Write(getSessionValue(c))
}

func getSessionValue(c *routing.Context) entity.User {
	return GetUserSession(c)
}
