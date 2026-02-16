package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/yourusername/golang-auth-api-boilerplate/config"
	"github.com/yourusername/golang-auth-api-boilerplate/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establishes database connection
func ConnectDatabase() {
	cfg := config.AppConfig

	// Create database if not exists
	createDatabase(cfg)

	// DSN for MySQL/MariaDB
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate models
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")
}

// GetDB returns database instance
func GetDB() *gorm.DB {
	return DB
}

// createDatabase creates the database if it doesn't exist (like Eloquent ORM)
func createDatabase(cfg *config.Config) {
	// Connect without database name to create it
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Warning: Could not connect to MySQL server: %v", err)
		return
	}
	defer db.Close()

	// Create database if not exists
	createDBQuery := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	_, err = db.Exec(createDBQuery)
	if err != nil {
		log.Printf("Warning: Could not create database: %v", err)
		return
	}

	log.Printf("Database '%s' ensured to exist", cfg.DBName)
}
