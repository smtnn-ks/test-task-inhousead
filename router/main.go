package router

import "github.com/gofiber/fiber/v2"

type registerDto_t struct {
	Name     string `json:"name" validate:"required,min=4,max=30"`
	Password string `json:"password" validate:"required,min=4,max=30"`
}

type crudDto_t struct {
	Title      string `json:"title"`
	Categories []int  `json:"categories"`
}

// category and item share the same structure
type entity_t struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func Init() (app *fiber.App) {
	app = fiber.New()
	app.Post("/register", registerHandler)
	app.Get("/categories", getCategoriesHandler)
	app.Get("/categories/:categoryid", getItemsHandler)

	app.Post("/categories", authMiddleware, createCategoryHandler)
	app.Put("/categories/:categoryid", authMiddleware, updateCategoryHandler)
	app.Delete("/categories/:categoryid", authMiddleware, removeCategoryHandler)

	app.Post("/items", authMiddleware, createItemHandler)
	app.Put("/items/:itemid", authMiddleware, updateItemHandler)
	app.Delete("/items/:itemid", authMiddleware, removeItemHandler)

	return
}
