package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type handler struct{}

var (
	sites []Site
	mu    sync.Mutex
)

func NewHandler() *handler {
	return &handler{}
}

func (*handler) getIndex(c *fiber.Ctx) error {
	var sites []Site
	site := Site{
		Url:         "https://www.google.com",
		Status:      checkSites("https://www.google.com"),
		Frequency:   5 * time.Minute,
		LastChecked: time.Now(),
	}
	// NOTE: This is temporary
	site.ID = 1

	sites = append(sites, site)
	return c.Render("index", fiber.Map{"Sites": sites}, "layouts/app")
}

func (h *handler) getSite(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Render("site", fiber.Map{"errors": "Invalid Site ID"}, "layouts/app")
	}
	// find site from sites by ID
	for _, site := range sites {
		if int(site.ID) == id {
			return c.Render("site", fiber.Map{"Site": site}, "layouts/app")
		}
	}
	return c.Redirect("/", 302)
}

func (*handler) getLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{}, "layouts/guest")
}

func (h *handler) postAddSite(c *fiber.Ctx) error {
	var site Site
	site.Url = c.FormValue("url")

	site.Frequency, _ = time.ParseDuration(c.FormValue("frequency"))
	site.Status = checkSites(site.Url)
	site.LastChecked = time.Now()
	sites = append(sites, site)
	return c.Redirect("/", 302)
}

func checkSites(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "ERROR"
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 404:
		return "404 NOT FOUND"
	case 500:
		return "500 SERVER ERROR"
	case 200:
		return "200 OK"
	default:
		return "DOWN"
	}
}

func monitorSites(site *Site, stop chan struct{}) {
	ticker := time.NewTicker(site.Frequency)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mu.Lock()
			site.Status = checkSites(site.Url)
			site.LastChecked = time.Now()
			mu.Unlock()
		case <-stop:
			return
		}
	}
}
