package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"federation-auth/internal/models"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Connection string for SQL Server
	// Format: sqlserver://username:password@host:port?database=dbname
	// Note: In docker-compose, the host is 'sql-server'.
	// We need to ensure the password matches what's in docker-compose.yml

	server := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	port := 1433

	// 1. Connect to 'master' to create the database if it doesn't exist
	masterConnString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=master&encrypt=disable", user, password, server, port)

	var err error
	var masterDB *gorm.DB

	// Retry logic for SQL Server startup
	for i := 0; i < 15; i++ {
		masterDB, err = gorm.Open(sqlserver.Open(masterConnString), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Waiting for SQL Server... (attempt %d/15): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("Could not connect to SQL Server master database")
	}

	// Create Database if not exists
	// Note: GORM doesn't support CREATE DATABASE directly via AutoMigrate, so we use Raw SQL
	// We need to use Exec, but GORM's Exec might try to wrap in transaction which CREATE DATABASE doesn't allow in MSSQL?
	// Actually, usually it's fine if we just run it.
	// But safely: IF NOT EXISTS(SELECT * FROM sys.databases WHERE name = 'FederationDB') CREATE DATABASE [FederationDB]

	createDbSQL := fmt.Sprintf("IF NOT EXISTS(SELECT * FROM sys.databases WHERE name = '%s') CREATE DATABASE [%s]", database, database)
	if err := masterDB.Exec(createDbSQL).Error; err != nil {
		log.Printf("Warning: Could not create database (might already exist): %v", err)
	}

	// Close master connection (optional, GORM manages pool)
	sqlDB, _ := masterDB.DB()
	sqlDB.Close()

	// 2. Connect to the actual database
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable", user, password, server, port, database)
	DB, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not connect to target database:", err)
	}

	log.Println("Connected to SQL Server!")

	// Auto Migrate
	err = DB.AutoMigrate(&models.Organization{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed.")

	// Seed Data (Idempotent)
	seedData()
}

func seedData() {
	orgs := []models.Organization{
		{Domain: "azure-corp.com", Auth0OrgID: "org_YZurKvaoBfFlbQTr", DisplayName: "Azure Corp"},
		{Domain: "google-corp.com", Auth0OrgID: "org_8mcqihdxiRlYI9t8", DisplayName: "Google Corp"},
		{Domain: "ldap-corp.com", Auth0OrgID: "org_fDx7gecDJ4hAmKDu", DisplayName: "LDAP Corp"},
		{Domain: "okta-corp.com", Auth0OrgID: "org_5DnOZCUZyLy2f8sT", DisplayName: "Okta Corp"},
		{Domain: "saml-corp.com", Auth0OrgID: "org_wtg281Bn9Lb3oaXn", DisplayName: "SAML Corp"},
		{Domain: "ba2kxoutlook.onmicrosoft.com", Auth0OrgID: "org_OqugCGseq85xF0rm", DisplayName: "Real Azure Corp"},
	}

	for _, org := range orgs {
		var count int64
		DB.Model(&models.Organization{}).Where("domain = ?", org.Domain).Count(&count)
		if count == 0 {
			DB.Create(&org)
			log.Printf("Seeded organization: %s", org.Domain)
		}
	}
}
