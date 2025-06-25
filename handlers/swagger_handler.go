package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed swagger_ui/*
var swaggerUIFiles embed.FS

// SwaggerHandler handles Swagger UI routes
type SwaggerHandler struct{}

// NewSwaggerHandler creates a new swagger handler
func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// SetupSwaggerRoutes sets up Swagger UI routes
func (h *SwaggerHandler) SetupSwaggerRoutes(app *fiber.App) {
	// Serve swagger.yaml file
	app.Get("/docs/swagger.yaml", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.yaml")
	})

	// Serve Swagger UI
	app.Get("/docs", h.SwaggerUI)
	app.Get("/docs/", h.SwaggerUI)
	
	// Serve Swagger UI static files
	swaggerUIFS, err := fs.Sub(swaggerUIFiles, "swagger_ui")
	if err == nil {
		app.Use("/docs/swagger-ui", filesystem.New(filesystem.Config{
			Root: http.FS(swaggerUIFS),
		}))
	}
}

// SwaggerUI serves the Swagger UI HTML page
func (h *SwaggerHandler) SwaggerUI(c *fiber.Ctx) error {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go Proxy Rotator API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui.css" />
    <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5.10.3/favicon-32x32.png" sizes="32x32" />
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin:0;
            background: #fafafa;
        }
        .swagger-ui .topbar {
            background-color: #667eea;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .swagger-ui .topbar .download-url-wrapper {
            display: none;
        }
        .swagger-ui .info .title {
            color: #3b4151;
        }
        .custom-header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            text-align: center;
            margin-bottom: 0;
        }
        .custom-header h1 {
            margin: 0;
            font-size: 2rem;
            font-weight: 600;
        }
        .custom-header p {
            margin: 10px 0 0 0;
            opacity: 0.9;
            font-size: 1.1rem;
        }
        .custom-links {
            background: white;
            padding: 15px 20px;
            border-bottom: 1px solid #e8e8e8;
            display: flex;
            justify-content: center;
            gap: 20px;
            flex-wrap: wrap;
        }
        .custom-links a {
            color: #667eea;
            text-decoration: none;
            font-weight: 500;
            padding: 8px 16px;
            border-radius: 6px;
            border: 1px solid #667eea;
            transition: all 0.3s ease;
        }
        .custom-links a:hover {
            background: #667eea;
            color: white;
        }
        #swagger-ui {
            max-width: 1200px;
            margin: 0 auto;
        }
    </style>
</head>
<body>
    <div class="custom-header">
        <h1>🔄 Go Proxy Rotator API</h1>
        <p>Advanced proxy rotation with health monitoring and management</p>
    </div>
    
    <div class="custom-links">
        <a href="/" target="_blank">🌐 Web Interface</a>
        <a href="/health" target="_blank">🏥 Health Check</a>
        <a href="/docs/swagger.yaml" target="_blank">📄 OpenAPI Spec</a>
        <a href="https://github.com/your-repo/go-proxy-rotator" target="_blank">📚 GitHub</a>
    </div>

    <div id="swagger-ui"></div>

    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.10.3/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/docs/swagger.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                tryItOutEnabled: true,
                requestInterceptor: function(request) {
                    // Add any request interceptors here
                    return request;
                },
                responseInterceptor: function(response) {
                    // Add any response interceptors here
                    return response;
                },
                onComplete: function() {
                    console.log("Swagger UI loaded successfully");
                },
                onFailure: function(error) {
                    console.error("Failed to load Swagger UI:", error);
                },
                docExpansion: "list",
                filter: true,
                showRequestHeaders: true,
                showCommonExtensions: true,
                defaultModelsExpandDepth: 2,
                defaultModelExpandDepth: 2,
                displayRequestDuration: true,
                operationsSorter: "alpha",
                tagsSorter: "alpha"
            });

            // Custom styling after load
            setTimeout(function() {
                // Hide the top bar
                const topbar = document.querySelector('.swagger-ui .topbar');
                if (topbar) {
                    topbar.style.display = 'none';
                }
            }, 1000);
        };
    </script>
</body>
</html>`

	c.Set("Content-Type", "text/html")
	return c.SendString(html)
}