package handlers

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"go-proxy-rotator/models"
	"go-proxy-rotator/services"

	"github.com/gofiber/fiber/v2"
)

type ProxyHandler struct {
	proxyService *services.ProxyService
}

func NewProxyHandler(proxyService *services.ProxyService) *ProxyHandler {
	return &ProxyHandler{proxyService: proxyService}
}

// UploadProxyFile handles proxy file uploads
func (h *ProxyHandler) UploadProxyFile(c *fiber.Ctx) error {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	// Check file extension
	if file.Header.Get("Content-Type") != "text/plain" && 
	   !isTextFile(file.Filename) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only .txt files are allowed",
		})
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Parse proxies from file
	proxies, err := h.proxyService.ParseProxyFile(src)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to parse proxy file: %v", err),
		})
	}

	if len(proxies) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No valid proxies found in file",
		})
	}

	// Add proxies to database
	added, err := h.proxyService.AddProxiesFromFile(proxies)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add proxies: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message":        "Proxies uploaded successfully",
		"total_parsed":   len(proxies),
		"total_added":    added,
		"total_skipped":  len(proxies) - added,
	})
}

// GetAllProxies returns all proxies
func (h *ProxyHandler) GetAllProxies(c *fiber.Ctx) error {
	proxies, err := h.proxyService.DB.GetAllProxies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get proxies",
		})
	}

	return c.JSON(fiber.Map{
		"proxies": proxies,
		"count":   len(proxies),
	})
}

// GetActiveProxies returns only active proxies
func (h *ProxyHandler) GetActiveProxies(c *fiber.Ctx) error {
	proxies, err := h.proxyService.DB.GetActiveProxies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get active proxies",
		})
	}

	return c.JSON(fiber.Map{
		"proxies": proxies,
		"count":   len(proxies),
	})
}

// DeleteProxy deletes a proxy by ID
func (h *ProxyHandler) DeleteProxy(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid proxy ID",
		})
	}

	err = h.proxyService.DB.DeleteProxy(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Proxy deleted successfully",
	})
}

// AddProxy manually adds a single proxy
func (h *ProxyHandler) AddProxy(c *fiber.Ctx) error {
	var proxy models.Proxy
	if err := c.BodyParser(&proxy); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if proxy.Host == "" || proxy.Port == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Host and port are required",
		})
	}

	// Set defaults
	if proxy.Protocol == "" {
		proxy.Protocol = "http"
	}
	proxy.IsActive = true

	err := h.proxyService.DB.AddProxy(&proxy)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to add proxy: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Proxy added successfully",
		"proxy":   proxy,
	})
}

// GetProxyStats returns proxy statistics
func (h *ProxyHandler) GetProxyStats(c *fiber.Ctx) error {
	stats, err := h.proxyService.DB.GetProxyStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get proxy statistics",
		})
	}

	return c.JSON(stats)
}

// HealthCheckProxies performs health check on all proxies
func (h *ProxyHandler) HealthCheckProxies(c *fiber.Ctx) error {
	testURL := c.Query("url", "https://httpbin.org/ip")
	
	err := h.proxyService.HealthCheckAllProxies(testURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Health check failed: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Health check completed",
	})
}

// ClearAllProxies removes all proxies
func (h *ProxyHandler) ClearAllProxies(c *fiber.Ctx) error {
	err := h.proxyService.DB.ClearAllProxies()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear proxies",
		})
	}

	return c.JSON(fiber.Map{
		"message": "All proxies cleared successfully",
	})
}

// basicAuth returns the base64 encoding of username:password
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// isTextFile checks if the filename has a .txt extension
func isTextFile(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".txt")
}