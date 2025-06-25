
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-proxy-rotator/config"
	"go-proxy-rotator/database"
	"go-proxy-rotator/handlers"
	"go-proxy-rotator/models"
	"go-proxy-rotator/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize services
	proxyService := services.NewProxyService(db)

	// Initialize handlers
	proxyHandler := handlers.NewProxyHandler(proxyService)

	// Load initial proxies if database is empty
	if err := loadInitialProxies(proxyService, db); err != nil {
		log.Printf("Warning: Failed to load initial proxies: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: int(cfg.MaxFileSize),
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// API routes
	api := app.Group("/api/v1")
	
	// Proxy management routes
	api.Post("/proxies/upload", proxyHandler.UploadProxyFile)
	api.Get("/proxies", proxyHandler.GetAllProxies)
	api.Get("/proxies/active", proxyHandler.GetActiveProxies)
	api.Post("/proxies", proxyHandler.AddProxy)
	api.Delete("/proxies/:id", proxyHandler.DeleteProxy)
	api.Delete("/proxies", proxyHandler.ClearAllProxies)
	api.Get("/proxies/stats", proxyHandler.GetProxyStats)
	api.Post("/proxies/health-check", proxyHandler.HealthCheckProxies)

	// Static files
	app.Static("/", "./static")

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now(),
		})
	})

	// Main proxy middleware (for actual proxy usage)
	app.Use(func(c *fiber.Ctx) error {
		// Skip API routes
		if c.Path() == "/health" || c.Path() == "/api/v1" || 
		   len(c.Path()) > 7 && c.Path()[:7] == "/api/v1" {
			return c.Next()
		}

		// Get a random proxy
		selectedProxy, err := proxyService.GetRandomProxy()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "No proxies available",
			})
		}

		// Create proxy URL
		proxyURL := selectedProxy.GetURL()
		if selectedProxy.Username != "" && selectedProxy.Password != "" {
			u, _ := url.Parse(proxyURL)
			u.User = url.UserPassword(selectedProxy.Username, selectedProxy.Password)
			proxyURL = u.String()
		}

		// Forward the request to the proxy
		if err := proxy.Do(c, proxyURL); err != nil {
			// Update proxy health on failure
			go func() {
				db.UpdateProxyHealth(selectedProxy.ID, 0, false)
			}()
			return err
		}

		// Update proxy health on success
		go func() {
			db.UpdateProxyHealth(selectedProxy.ID, 100, true)
		}()

		return nil
	})

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}

// loadInitialProxies loads initial proxies if database is empty
func loadInitialProxies(proxyService *services.ProxyService, db *database.DB) error {
	// Check if we already have proxies
	stats, err := db.GetProxyStats()
	if err != nil {
		return err
	}
	
	if stats.TotalProxies > 0 {
		log.Printf("Database already contains %d proxies", stats.TotalProxies)
		return nil
	}

	log.Println("Loading initial proxy list...")
	
	proxyList := []string{
		"31.58.16.86:6053:nufcfuwi:wtohyhtey9hp",
		"198.37.116.237:6196:nufcfuwi:wtohyhtey9hp",
		"216.173.111.14:6724:nufcfuwi:wtohyhtey9hp",
		"45.43.70.68:6355:nufcfuwi:wtohyhtey9hp",
		"89.249.193.146:5884:nufcfuwi:wtohyhtey9hp",
		"154.30.244.127:9568:nufcfuwi:wtohyhtey9hp",
		"50.114.15.66:6051:nufcfuwi:wtohyhtey9hp",
		"82.22.230.178:7516:nufcfuwi:wtohyhtey9hp",
		"104.245.244.67:6507:nufcfuwi:wtohyhtey9hp",
		"107.181.152.104:5141:nufcfuwi:wtohyhtey9hp",
		"166.88.235.80:5708:nufcfuwi:wtohyhtey9hp",
		"181.41.218.183:5769:nufcfuwi:wtohyhtey9hp",
		"45.41.169.116:6777:nufcfuwi:wtohyhtey9hp",
		"45.61.125.159:6170:nufcfuwi:wtohyhtey9hp",
		"82.25.216.37:6879:nufcfuwi:wtohyhtey9hp",
		"67.227.37.102:5644:nufcfuwi:wtohyhtey9hp",
		"38.154.217.209:7400:nufcfuwi:wtohyhtey9hp",
		"45.39.50.10:6428:nufcfuwi:wtohyhtey9hp",
		"45.61.127.65:6004:nufcfuwi:wtohyhtey9hp",
		"154.29.239.97:6136:nufcfuwi:wtohyhtey9hp",
		"45.61.122.177:6469:nufcfuwi:wtohyhtey9hp",
		"45.147.186.165:7038:nufcfuwi:wtohyhtey9hp",
		"194.113.112.134:6029:nufcfuwi:wtohyhtey9hp",
		"92.113.236.49:6634:nufcfuwi:wtohyhtey9hp",
		"64.137.93.234:6691:nufcfuwi:wtohyhtey9hp",
		"82.22.228.183:7523:nufcfuwi:wtohyhtey9hp",
		"38.154.197.219:6885:nufcfuwi:wtohyhtey9hp",
		"45.41.177.216:5866:nufcfuwi:wtohyhtey9hp",
		"23.27.93.231:5810:nufcfuwi:wtohyhtey9hp",
		"64.43.91.69:6840:nufcfuwi:wtohyhtey9hp",
		"162.220.246.240:6524:nufcfuwi:wtohyhtey9hp",
		"206.206.124.50:6631:nufcfuwi:wtohyhtey9hp",
		"154.92.114.134:5829:nufcfuwi:wtohyhtey9hp",
		"31.58.26.44:6627:nufcfuwi:wtohyhtey9hp",
		"38.170.161.240:9291:nufcfuwi:wtohyhtey9hp",
		"45.61.98.210:5894:nufcfuwi:wtohyhtey9hp",
		"31.59.18.60:6641:nufcfuwi:wtohyhtey9hp",
		"31.57.83.211:5785:nufcfuwi:wtohyhtey9hp",
		"38.153.137.40:5348:nufcfuwi:wtohyhtey9hp",
		"31.57.43.105:6179:nufcfuwi:wtohyhtey9hp",
		"142.111.44.26:5738:nufcfuwi:wtohyhtey9hp",
		"173.211.30.252:6686:nufcfuwi:wtohyhtey9hp",
		"161.123.101.238:6864:nufcfuwi:wtohyhtey9hp",
		"82.22.228.79:7419:nufcfuwi:wtohyhtey9hp",
		"82.23.235.24:5348:nufcfuwi:wtohyhtey9hp",
		"192.3.48.17:6010:nufcfuwi:wtohyhtey9hp",
		"45.138.119.199:5948:nufcfuwi:wtohyhtey9hp",
		"82.29.243.183:8004:nufcfuwi:wtohyhtey9hp",
		"142.147.128.160:6660:nufcfuwi:wtohyhtey9hp",
		"155.254.39.166:6124:nufcfuwi:wtohyhtey9hp",
		"206.232.71.122:6693:nufcfuwi:wtohyhtey9hp",
		"23.27.208.191:5901:nufcfuwi:wtohyhtey9hp",
		"45.38.78.134:6071:nufcfuwi:wtohyhtey9hp",
		"45.61.98.55:5739:nufcfuwi:wtohyhtey9hp",
		"104.252.20.92:6024:nufcfuwi:wtohyhtey9hp",
		"161.123.5.80:5129:nufcfuwi:wtohyhtey9hp",
		"173.211.8.33:6145:nufcfuwi:wtohyhtey9hp",
		"172.98.168.234:6881:nufcfuwi:wtohyhtey9hp",
		"217.69.127.116:6737:nufcfuwi:wtohyhtey9hp",
		"46.202.227.221:6215:nufcfuwi:wtohyhtey9hp",
		"45.56.173.118:6101:nufcfuwi:wtohyhtey9hp",
		"192.186.172.144:9144:nufcfuwi:wtohyhtey9hp",
		"206.232.13.213:5879:nufcfuwi:wtohyhtey9hp",
		"82.23.202.122:6974:nufcfuwi:wtohyhtey9hp",
		"147.124.198.57:5916:nufcfuwi:wtohyhtey9hp",
		"23.27.93.195:5774:nufcfuwi:wtohyhtey9hp",
		"104.239.39.220:6149:nufcfuwi:wtohyhtey9hp",
		"173.211.0.30:6523:nufcfuwi:wtohyhtey9hp",
		"38.154.195.140:9228:nufcfuwi:wtohyhtey9hp",
		"45.41.162.217:6854:nufcfuwi:wtohyhtey9hp",
		"173.214.176.243:6214:nufcfuwi:wtohyhtey9hp",
		"185.216.107.119:5696:nufcfuwi:wtohyhtey9hp",
		"45.56.175.15:5689:nufcfuwi:wtohyhtey9hp",
		"50.114.93.187:6171:nufcfuwi:wtohyhtey9hp",
		"154.29.235.146:6487:nufcfuwi:wtohyhtey9hp",
	}

	// Parse and add proxies
	var proxies []*models.Proxy
	for _, line := range proxyList {
		proxy, err := parseProxyLine(line)
		if err != nil {
			continue
		}
		proxies = append(proxies, proxy)
	}

	added, err := proxyService.AddProxiesFromFile(proxies)
	if err != nil {
		return err
	}

	log.Printf("Loaded %d initial proxies", added)
	return nil
}

// parseProxyLine parses a proxy line in format host:port:username:password
func parseProxyLine(line string) (*models.Proxy, error) {
	parts := strings.Split(line, ":")
	if len(parts) != 4 {
		return nil, fmt.Errorf("invalid proxy format")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port")
	}

	return &models.Proxy{
		Host:      parts[0],
		Port:      port,
		Username:  parts[2],
		Password:  parts[3],
		Protocol:  "http",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// basicAuth returns the base64 encoding of username:password.
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
