package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

func main() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&User{}, &Site{}, &Endpoint{})

	engine := django.New("./static", ".html")

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Static("/static", "./static")

	h := NewHandler(db)

	app.Get("/", h.getIndex).Name("index")
	app.Get("/login", h.getLogin).Name("login")
	app.Get("/site/:id", h.getSite).Name("site")
	app.Post("/add", h.postAddSite).Name("add")
	app.Post("/endpoint/add", h.postAddEndpoint).Name("add-endpoint")

	log.Fatal(app.Listen(":3000"))

}
