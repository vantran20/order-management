package users

import (
	"omg/api/internal/controller/users"
)

type Handler struct {
	controller users.Controller
}

func New(controller users.Controller) Handler {
	return Handler{
		controller: controller,
	}
}
