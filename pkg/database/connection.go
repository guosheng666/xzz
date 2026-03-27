package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(host string, port int, user, password, database, charset string, timeout int) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=%ds&parseTime=true&loc=Local",
		user, password, host, port, database, charset, timeout)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec("SET GLOBAL max_allowed_packet=67108864")
	if err != nil {
		fmt.Printf("Warning: Could not set max_allowed_packet (may require SUPER privilege): %v\n", err)
	}

	_, err = db.Exec("SET SESSION foreign_key_checks=0")
	if err != nil {
		return nil, fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	_, err = db.Exec("SET SESSION unique_checks=0")
	if err != nil {
		return nil, fmt.Errorf("failed to disable unique checks: %w", err)
	}

	_, err = db.Exec("SET SESSION sql_log_bin=0")
	if err != nil {
		fmt.Printf("Warning: Could not disable binary logging (may require SUPER privilege): %v\n", err)
	}

	return db, nil
}

func ConnectOrCreate(host string, port int, user, password, database, charset string, timeout int) (*sql.DB, error) {
	if database == "" {
		return nil, fmt.Errorf("database name is required")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&timeout=%ds&parseTime=true&loc=Local",
		user, password, host, port, charset, timeout)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", database))
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	fmt.Printf("Database '%s' created or already exists\n", database)

	db.Close()

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&timeout=%ds&parseTime=true&loc=Local",
		user, password, host, port, database, charset, timeout)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec("SET GLOBAL max_allowed_packet=67108864")
	if err != nil {
		fmt.Printf("Warning: Could not set max_allowed_packet (may require SUPER privilege): %v\n", err)
	}

	_, err = db.Exec("SET SESSION foreign_key_checks=0")
	if err != nil {
		return nil, fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	_, err = db.Exec("SET SESSION unique_checks=0")
	if err != nil {
		return nil, fmt.Errorf("failed to disable unique checks: %w", err)
	}

	_, err = db.Exec("SET SESSION sql_log_bin=0")
	if err != nil {
		fmt.Printf("Warning: Could not disable binary logging (may require SUPER privilege): %v\n", err)
	}

	return db, nil
}
