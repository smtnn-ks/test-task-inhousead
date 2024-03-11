package router

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ise(err error) error {
	log.Println("[ERROR] ", err)
	return fiber.ErrInternalServerError
}

func registerHandler(c *fiber.Ctx) error {
	var dto registerDto_t
	if err := c.BodyParser(&dto); err != nil {
		return fiber.ErrBadRequest
	}

	if err := validator.New().Struct(dto); err != nil {
		c.SendString(err.Error())
		return c.SendStatus(fiber.StatusBadRequest)
	}
	log.Printf("Register dto: %v\n", dto)

	if err := register(dto.Name, dto.Password); err != nil {
		if err == fiber.ErrBadRequest {
			c.SendString("Такой пользователь уже зарегистрирован")
			return c.SendStatus(fiber.StatusBadRequest)
		}
		return fiber.ErrInternalServerError
	}

	return nil
}

func getCategoriesHandler(c *fiber.Ctx) error {
	res, err := getCategories()
	if err != nil {
		return err
	}
	payload, err := json.Marshal(res)
	if err != nil {
		return ise(err)
	}
	return c.Send(payload)
}

func getItemsHandler(c *fiber.Ctx) error {
	categoryId, err := strconv.ParseInt(c.Params("categoryid"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	res, err := getItems(int(categoryId))
	if err != nil {
		return ise(err)
	}
	payload, err := json.Marshal(res)
	if err != nil {
		return ise(err)
	}
	return c.Send(payload)
}

func createCategoryHandler(c *fiber.Ctx) error {
	var dto crudDto_t
	if err := c.BodyParser(&dto); err != nil {
		return ise(err)
	}
	if dto.Title == "" {
		return fiber.ErrBadRequest
	}
	if err := createCategory(dto.Title); err != nil {
		if err == fiber.ErrBadRequest {
			return err
		}
		return ise(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func updateCategoryHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("categoryid"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	var dto crudDto_t
	if err := c.BodyParser(&dto); err != nil {
		return ise(err)
	}
	if dto.Title == "" {
		return fiber.ErrBadRequest
	}

	if err := updateCategory(int(id), dto.Title); err != nil {
		return ise(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func removeCategoryHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("categoryid"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}
	if err := removeCategory(int(id)); err != nil {
		return ise(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func createItemHandler(c *fiber.Ctx) error {
	var dto crudDto_t
	if err := c.BodyParser(&dto); err != nil {
		return ise(err)
	}
	if dto.Title == "" || len(dto.Categories) == 0 {
		return fiber.ErrBadRequest
	}
	if err := createItem(dto.Title, dto.Categories); err != nil {
		if err == fiber.ErrBadRequest {
			return err
		}
		return ise(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func updateItemHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("itemid"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	var dto crudDto_t
	if err := c.BodyParser(&dto); err != nil {
		return ise(err)
	}
	if dto.Title == "" || len(dto.Categories) == 0 {
		return fiber.ErrBadRequest
	}

	if err := updateItem(int(id), dto.Title, dto.Categories); err != nil {
		return ise(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func removeItemHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("itemid"), 10, 64)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := removeItem(int(id)); err != nil {
		return ise(err)
	}
	return c.SendStatus(fiber.StatusOK)
}
