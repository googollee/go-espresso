package framework

import (
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

type User struct {
	AccessKey string
	SecretKey string
	Name      string
}

type ContextData struct {
	User *User
}

func TestAuth(t *testing.T) {
	users := []User{
		{"ak1", "sk1", "user1"},
		{"ak2", "sk2", "user2"},
	}

	userMap := make(map[string]*User)
	for _, user := range users {
		user := user
		userMap[user.AccessKey] = &user
	}

	auth := NewAuthWithAkSk("Bearer", func(ctx *Context[ContextData], ak, hash string) error {
		user, ok := userMap[ak]
		if !ok {
			return errors.New("no user")
		}

		cloneUser := *user
		ctx.Data.User = &cloneUser

		return nil
	})

	eng := gin.Default()
	userGroup := Group(eng, "/users", ContextData{})
	userGroup.GET("/", auth.Handle, func(ctx *Context[ContextData]) {
		// consume ctx.Data.User
	})
}
