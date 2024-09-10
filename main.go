package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&User{}, &Site{})

	engine := html.New("./static", ".html")
	app := fiber.New(fiber.Config{Views: engine})

	h := NewHandler()

	app.Static("/static", "./static")

	app.Get("/", h.getIndex)
	app.Get("/login", h.getLogin)
	app.Get("/site/:id", h.getSite)
	app.Post("/add", h.postAddSite)

	log.Fatal(app.Listen(":3000"))

}
