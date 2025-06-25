package database

import (
	"database/sql"
	"fmt"
	"time"

	"go-proxy-rotator/models"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS proxies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		username TEXT DEFAULT '',
		password TEXT DEFAULT '',
		protocol TEXT DEFAULT 'http',
		is_active BOOLEAN DEFAULT 1,
		last_checked DATETIME DEFAULT CURRENT_TIMESTAMP,
		response_time INTEGER DEFAULT 0,
		fail_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(host, port)
	);
	
	CREATE INDEX IF NOT EXISTS idx_proxies_active ON proxies(is_active);
	CREATE INDEX IF NOT EXISTS idx_proxies_health ON proxies(is_active, fail_count);
	`

	_, err := db.conn.Exec(query)
	return err
}

// AddProxy adds a new proxy to the database
func (db *DB) AddProxy(proxy *models.Proxy) error {
	query := `
	INSERT INTO proxies (host, port, username, password, protocol, is_active, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := db.conn.Exec(query, proxy.Host, proxy.Port, proxy.Username, 
		proxy.Password, proxy.Protocol, proxy.IsActive, now, now)
	if err != nil {
		return fmt.Errorf("failed to add proxy: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	proxy.ID = int(id)
	proxy.CreatedAt = now
	proxy.UpdatedAt = now
	return nil
}

// GetActiveProxies returns all active proxies
func (db *DB) GetActiveProxies() ([]*models.Proxy, error) {
	query := `
	SELECT id, host, port, username, password, protocol, is_active, 
		   last_checked, response_time, fail_count, created_at, updated_at
	FROM proxies 
	WHERE is_active = 1 AND fail_count < 5
	ORDER BY response_time ASC, fail_count ASC
	`
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active proxies: %w", err)
	}
	defer rows.Close()

	var proxies []*models.Proxy
	for rows.Next() {
		proxy := &models.Proxy{}
		err := rows.Scan(&proxy.ID, &proxy.Host, &proxy.Port, &proxy.Username,
			&proxy.Password, &proxy.Protocol, &proxy.IsActive, &proxy.LastChecked,
			&proxy.ResponseTime, &proxy.FailCount, &proxy.CreatedAt, &proxy.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// GetAllProxies returns all proxies
func (db *DB) GetAllProxies() ([]*models.Proxy, error) {
	query := `
	SELECT id, host, port, username, password, protocol, is_active, 
		   last_checked, response_time, fail_count, created_at, updated_at
	FROM proxies 
	ORDER BY created_at DESC
	`
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all proxies: %w", err)
	}
	defer rows.Close()

	var proxies []*models.Proxy
	for rows.Next() {
		proxy := &models.Proxy{}
		err := rows.Scan(&proxy.ID, &proxy.Host, &proxy.Port, &proxy.Username,
			&proxy.Password, &proxy.Protocol, &proxy.IsActive, &proxy.LastChecked,
			&proxy.ResponseTime, &proxy.FailCount, &proxy.CreatedAt, &proxy.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proxy: %w", err)
		}
		proxies = append(proxies, proxy)
	}

	return proxies, nil
}

// UpdateProxyHealth updates proxy health information
func (db *DB) UpdateProxyHealth(id int, responseTime int, success bool) error {
	var query string
	var args []interface{}
	
	if success {
		query = `
		UPDATE proxies 
		SET last_checked = ?, response_time = ?, fail_count = 0, updated_at = ?
		WHERE id = ?
		`
		now := time.Now()
		args = []interface{}{now, responseTime, now, id}
	} else {
		query = `
		UPDATE proxies 
		SET last_checked = ?, fail_count = fail_count + 1, updated_at = ?,
			is_active = CASE WHEN fail_count + 1 >= 5 THEN 0 ELSE is_active END
		WHERE id = ?
		`
		now := time.Now()
		args = []interface{}{now, now, id}
	}

	_, err := db.conn.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update proxy health: %w", err)
	}

	return nil
}

// DeleteProxy deletes a proxy by ID
func (db *DB) DeleteProxy(id int) error {
	query := "DELETE FROM proxies WHERE id = ?"
	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete proxy: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("proxy with id %d not found", id)
	}

	return nil
}

// GetProxyStats returns statistics about proxies
func (db *DB) GetProxyStats() (*models.ProxyStats, error) {
	query := `
	SELECT 
		COUNT(*) as total,
		COALESCE(SUM(CASE WHEN is_active = 1 THEN 1 ELSE 0 END), 0) as active,
		COALESCE(SUM(CASE WHEN is_active = 1 AND fail_count < 5 AND response_time < 10000 THEN 1 ELSE 0 END), 0) as healthy,
		COALESCE(SUM(CASE WHEN is_active = 0 OR fail_count >= 5 THEN 1 ELSE 0 END), 0) as failed
	FROM proxies
	`

	row := db.conn.QueryRow(query)
	stats := &models.ProxyStats{}
	err := row.Scan(&stats.TotalProxies, &stats.ActiveProxies, &stats.HealthyProxies, &stats.FailedProxies)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy stats: %w", err)
	}

	return stats, nil
}

// ClearAllProxies removes all proxies from the database
func (db *DB) ClearAllProxies() error {
	query := "DELETE FROM proxies"
	_, err := db.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to clear all proxies: %w", err)
	}
	return nil
}