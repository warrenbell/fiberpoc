package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/sandstone2/fiberpoc/common/services"
	"go.uber.org/zap"
)

type FooHandler struct {
	fooService *services.FooServiceInterface
	logger     *zap.Logger
}

func NewFooHandler(fooService services.FooServiceInterface, logger *zap.Logger) *FooHandler {
	return &FooHandler{fooService: &fooService, logger: logger}
}

func (fooHandler *FooHandler) HandleGetFoos(c *fiber.Ctx) error {
	foos, err := (*fooHandler.fooService).GetFoos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error J5TSGF - Getting foos in handler."})
	}
	return c.JSON(foos)
}

func (fooHandler *FooHandler) HandleCreateFoo(c *fiber.Ctx) error {
	rowsAffected, err := (*fooHandler.fooService).CreateFoo()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error QONMRA - Creating foo in handler."})
	}
	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d foos created.", rowsAffected)})
}

func (fooHandler *FooHandler) HandleDeleteFoos(c *fiber.Ctx) error {
	rowsAffected, err := (*fooHandler.fooService).DeleteFoos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error 8HCIPG - Deleting foos in handler."})
	}
	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d foos deleted.", rowsAffected)})
}

func (fooHandler *FooHandler) HandleUpdateFoo(c *fiber.Ctx) error {
	fooID := c.QueryInt("fooId")

	rowsAffected, err := (*fooHandler.fooService).UpdateFoo(int64(fooID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error E4LP9X - Updating foo in handler."})
	}
	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d foos updated.", rowsAffected)})
}
