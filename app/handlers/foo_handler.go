package handlers

import (
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
	"gitlab.com/sandstone2/fiberpoc/common/models"
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error J5TSGF - Getting foos in handler. Error: %v", err)})
	}
	return c.JSON(foos)
}

func (fooHandler *FooHandler) HandleCreateFoo(c *fiber.Ctx) error {
	newFoo := models.Foo{}
	if err := c.BodyParser(&newFoo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error O1WQ9B - Bad request body."})
	}

	resultFoo, err := (*fooHandler.fooService).CreateFoo(newFoo.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error QONMRA - Creating foo in handler. Error: %v", err)})
	}
	return c.JSON(resultFoo)
}

func (fooHandler *FooHandler) HandleDeleteFoos(c *fiber.Ctx) error {
	rowsAffected, err := (*fooHandler.fooService).DeleteFoos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error 8HCIPG - Deleting foos in handler."})
	}
	return c.JSON(fiber.Map{"message": fmt.Sprintf("%d foos deleted.", rowsAffected)})
}

func (fooHandler *FooHandler) HandleUpdateFoo(c *fiber.Ctx) error {
	fooId, err := c.ParamsInt("id", 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error L41Q1S - Foo id is not a number."})
	}
	if fooId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error ZR53ES - No foo id was provided."})
	}

	updatedFoo := models.Foo{}
	if err := c.BodyParser(&updatedFoo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error O1WQ9B - Bad request body."})
	}

	foo, err := (*fooHandler.fooService).UpdateFoo(int64(fooId), updatedFoo.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": fmt.Sprintf("Error FSYTGZ - Updating foo. Error: %v", err)})
	}
	// if rowsAffected == 0 {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": fmt.Sprintf("Error XLM18M - Foo was not found with id %d.", fooId)})
	// }
	return c.JSON(foo)
}
