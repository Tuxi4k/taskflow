package main

import (
	"log"
	"taskflow/internal/database"
	"taskflow/internal/modules/task"

	"github.com/Tuxi4k/swaggen"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

var swaggerJSON = `{"swagger":"2.0","info":{"title":"Loading..."}}`

// @title TaskFlow API
// @version 1.0
// @description API для управления задачами
func main() {
	go generateSwagger()

	app := fiber.New()

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	registerSwagger(app)

	tasksRoutes := app.Group("/tasks")
	taskRepo := task.NewRepository(db)
	taskService := task.NewService(taskRepo)
	taskHandler := task.NewHandler(taskService)

	taskHandler.RegisterRoutes(tasksRoutes)

	log.Fatalf("Server error: %v", app.Listen(":3000"))
}

func generateSwagger() {
	json, _ := swaggen.Generate(
		swaggen.WithMainAPIFile("cmd/main.go"),
	)
	swaggerJSON = json
}

func registerSwagger(app fiber.Router) {
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(swaggerJSON)
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
}
