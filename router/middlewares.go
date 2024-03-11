package router

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func authMiddleware(c *fiber.Ctx) error {
	authString := c.Get("Authorization")
	if authString == "" {
		return fiber.ErrUnauthorized
	}

	strs := strings.Split(authString, " ")
	if len(strs) != 2 {
		return fiber.ErrUnauthorized
	}
	if strs[0] != "Basic" {
		return fiber.ErrUnauthorized
	}

	credsByte, err := base64.URLEncoding.DecodeString(strs[1])
	if err != nil {
		return fiber.ErrUnauthorized
	}

	creds := string(credsByte)
	creds = strings.Trim(creds, "\n")

	strs = strings.Split(creds, ":")

	if len(strs) != 2 {
		return fiber.ErrUnauthorized
	}

	if err := checkAuth(strs[0], strs[1]); err != nil {
		if err == fiber.ErrUnauthorized {
			return err
		} else {
			log.Println("[ERROR] ", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	return c.Next()
}
