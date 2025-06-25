package services

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-proxy-rotator/database"
	"go-proxy-rotator/models"
)

type ProxyService struct {
	DB *database.DB
}

func NewProxyService(db *database.DB) *ProxyService {
	return &ProxyService{DB: db}
}

// ParseProxyFile parses a proxy file and returns a list of proxies
func (s *ProxyService) ParseProxyFile(reader io.Reader) ([]*models.Proxy, error) {
	var proxies []*models.Proxy
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		proxy, err := s.parseProxyLine(line)
		if err != nil {
			// Log the error but continue parsing other lines
			continue
		}

		proxies = append(proxies, proxy)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading proxy file: %w", err)
	}

	return proxies, nil
}

// parseProxyLine parses a single proxy line in various formats
func (s *ProxyService) parseProxyLine(line string) (*models.Proxy, error) {
	// Support multiple formats:
	// 1. host:port:username:password
	// 2. host:port
	// 3. protocol://host:port
	// 4. protocol://username:password@host:port

	proxy := &models.Proxy{
		Protocol:  "http",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Check if it's a URL format
	if strings.Contains(line, "://") {
		u, err := url.Parse(line)
		if err != nil {
			return nil, fmt.Errorf("invalid URL format: %s", line)
		}

		proxy.Protocol = u.Scheme
		proxy.Host = u.Hostname()
		
		if u.Port() != "" {
			port, err := strconv.Atoi(u.Port())
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", u.Port())
			}
			proxy.Port = port
		} else {
			// Default ports
			switch u.Scheme {
			case "http":
				proxy.Port = 80
			case "https":
				proxy.Port = 443
			case "socks5":
				proxy.Port = 1080
			default:
				return nil, fmt.Errorf("unknown protocol: %s", u.Scheme)
			}
		}

		if u.User != nil {
			proxy.Username = u.User.Username()
			if password, set := u.User.Password(); set {
				proxy.Password = password
			}
		}
	} else {
		// Parse colon-separated format
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid proxy format: %s", line)
		}

		proxy.Host = parts[0]
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", parts[1])
		}
		proxy.Port = port

		if len(parts) >= 4 {
			proxy.Username = parts[2]
			proxy.Password = parts[3]
		}
	}

	// Validate required fields
	if proxy.Host == "" || proxy.Port == 0 {
		return nil, fmt.Errorf("missing host or port")
	}

	return proxy, nil
}

// AddProxiesFromFile adds proxies from a parsed file to the database
func (s *ProxyService) AddProxiesFromFile(proxies []*models.Proxy) (int, error) {
	added := 0
	for _, proxy := range proxies {
		err := s.DB.AddProxy(proxy)
		if err != nil {
			// Skip duplicates and continue
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				continue
			}
			return added, fmt.Errorf("failed to add proxy %s:%d: %w", proxy.Host, proxy.Port, err)
		}
		added++
	}
	return added, nil
}

// GetRandomProxy returns a random healthy proxy
func (s *ProxyService) GetRandomProxy() (*models.Proxy, error) {
	proxies, err := s.DB.GetActiveProxies()
	if err != nil {
		return nil, fmt.Errorf("failed to get active proxies: %w", err)
	}

	if len(proxies) == 0 {
		return nil, fmt.Errorf("no active proxies available")
	}

	// Filter for healthy proxies first
	var healthyProxies []*models.Proxy
	for _, proxy := range proxies {
		if proxy.IsHealthy() {
			healthyProxies = append(healthyProxies, proxy)
		}
	}

	// If no healthy proxies, use all active proxies
	if len(healthyProxies) == 0 {
		healthyProxies = proxies
	}

	// Return random proxy
	return healthyProxies[rand.Intn(len(healthyProxies))], nil
}

// CheckProxyHealth checks if a proxy is working
func (s *ProxyService) CheckProxyHealth(proxy *models.Proxy, testURL string) (int, bool) {
	start := time.Now()
	
	// Create proxy URL
	proxyURL, err := url.Parse(proxy.GetURL())
	if err != nil {
		return 0, false
	}

	// Set up proxy authentication if needed
	if proxy.Username != "" && proxy.Password != "" {
		proxyURL.User = url.UserPassword(proxy.Username, proxy.Password)
	}

	// Create HTTP client with proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// Make test request
	resp, err := client.Get(testURL)
	if err != nil {
		return 0, false
	}
	defer resp.Body.Close()

	responseTime := int(time.Since(start).Milliseconds())
	
	// Check if response is successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return responseTime, true
	}

	return responseTime, false
}

// HealthCheckAllProxies performs health checks on all active proxies
func (s *ProxyService) HealthCheckAllProxies(testURL string) error {
	proxies, err := s.DB.GetActiveProxies()
	if err != nil {
		return fmt.Errorf("failed to get active proxies: %w", err)
	}

	for _, proxy := range proxies {
		responseTime, success := s.CheckProxyHealth(proxy, testURL)
		err := s.DB.UpdateProxyHealth(proxy.ID, responseTime, success)
		if err != nil {
			// Log error but continue with other proxies
			continue
		}
	}

	return nil
}