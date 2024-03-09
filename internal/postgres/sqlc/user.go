package sqlc

import "github.com/tedmo/go-rest-template/internal/app"

func (u *User) DomainModel() *app.User {
	return &app.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
