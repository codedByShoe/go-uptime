package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

var (
	mu sync.Mutex
)

func NewHandler(db *gorm.DB) *handler {
	return &handler{
		db: db,
	}
}

func (h *handler) getIndex(c *fiber.Ctx) error {
	var sites []Site
	if err := h.db.Model(&Site{}).Preload("Endpoints").Find(&sites).Error; err != nil {
		return c.RedirectToRoute("index", fiber.Map{"Error": err.Error()})
	}
	return c.Render("index", fiber.Map{"Sites": sites}, "layouts/app")
}

func (h *handler) getSite(c *fiber.Ctx) error {
	var site Site
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.RedirectToRoute("index", fiber.Map{"Error": err.Error()})
	}
	if err := h.db.Model(&Site{}).Preload("Endpoints", "site_id = ?", id).First(&site, id).Error; err != nil {
		return c.RedirectToRoute("index", fiber.Map{"Error": err.Error()})
	}
	return c.Render("site", fiber.Map{"Site": site}, "layouts/app")
}

func (*handler) getLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{}, "layouts/guest")
}

func (h *handler) postAddSite(c *fiber.Ctx) error {
	url := c.FormValue("url")
	site := Site{
		Url: url,
	}
	// add site to database
	if err := h.db.Create(&site).Error; err != nil {
		return c.RedirectToRoute("index", fiber.Map{"Error": err.Error()})
	}
	// default endpoint & frequency
	endpoint := Endpoint{
		Path:        "/",
		Status:      checkSites(site.Url),
		Frequency:   5 * time.Minute,
		LastChecked: time.Now(),
		Uptime:      0,
		SiteID:      site.ID,
	}
	// add endpoint to database
	if err := h.db.Create(&endpoint).Error; err != nil {
		return c.RedirectToRoute("index", fiber.Map{"Error": err.Error()})
	}
	return c.Redirect(fmt.Sprintf("/site/%d", site.ID), 302)
}

func (h *handler) postAddEndpoint(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.FormValue("id"))
	path := c.FormValue("path")
	freq, _ := strconv.Atoi(c.FormValue("frequency"))
	url := c.FormValue("url")

	endpoint := Endpoint{
		Path:        path,
		Status:      checkSites(url + path),
		Frequency:   time.Duration(freq) * time.Second,
		LastChecked: time.Now(),
		Uptime:      0,
		SiteID:      uint(id),
	}
	if err := h.db.Create(&endpoint).Error; err != nil {
		return c.RedirectToRoute("index", fiber.Map{"Error": err.Error()})
	}
	return c.RedirectBack(fmt.Sprintf("/site/%d", id))
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

//
// func monitorSites(site *Site, stop chan struct{}) {
// 	ticker := time.NewTicker(site.Frequency)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			mu.Lock()
// 			site.Status = checkSites(site.Url)
// 			site.LastChecked = time.Now()
// 			mu.Unlock()
// 		case <-stop:
// 			return
// 		}
// 	}
// }
