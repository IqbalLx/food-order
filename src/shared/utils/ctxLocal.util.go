package utils

import "github.com/gofiber/fiber/v2"

// ref: https://betterprogramming.pub/sharing-a-database-connection-in-go-fiber-eedb305e9348

func SetLocal[T any](c *fiber.Ctx, key string, value T) {
	c.Locals(key, value)
}

func GetLocal[T any](c *fiber.Ctx, key string) T {
	return c.Locals(key).(T)
}