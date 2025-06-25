package models

import (
	"fmt"
	"time"
)

// Proxy represents a proxy server with its credentials and health status
type Proxy struct {
	ID          int       `json:"id" db:"id"`
	Host        string    `json:"host" db:"host"`
	Port        int       `json:"port" db:"port"`
	Username    string    `json:"username" db:"username"`
	Password    string    `json:"password" db:"password"`
	Protocol    string    `json:"protocol" db:"protocol"` // http, https, socks5
	IsActive    bool      `json:"is_active" db:"is_active"`
	LastChecked time.Time `json:"last_checked" db:"last_checked"`
	ResponseTime int      `json:"response_time" db:"response_time"` // in milliseconds
	FailCount   int       `json:"fail_count" db:"fail_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// GetURL returns the full proxy URL
func (p *Proxy) GetURL() string {
	return fmt.Sprintf("%s://%s:%d", p.Protocol, p.Host, p.Port)
}

// IsHealthy returns true if the proxy is considered healthy
func (p *Proxy) IsHealthy() bool {
	return p.IsActive && p.FailCount < 5 && p.ResponseTime < 10000
}

// ProxyStats represents statistics for proxy usage
type ProxyStats struct {
	TotalProxies   int `json:"total_proxies"`
	ActiveProxies  int `json:"active_proxies"`
	HealthyProxies int `json:"healthy_proxies"`
	FailedProxies  int `json:"failed_proxies"`
}