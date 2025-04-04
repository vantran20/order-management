package products

import "omg/api/internal/controller/products"

type Handler struct {
	controller products.Controller
}

func New(controller products.Controller) Handler {
	return Handler{
		controller: controller,
	}
}
