package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

// User represents a user in the new_users table.
type User struct {
	ID        int    `db:"id"`
	Email     string `db:"email"`
	EmailSent bool   `db:"email_sent"`
	InProcess bool   `db:"in_process"`
}

// App encapsulates dependencies for the application.
type App struct {
	DB     *sqlx.DB
	Config *Config
}

// Config holds the application configuration.
type Config struct {
	DBDSN       string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	SMTPFrom    string
	WorkerCount int
	BatchSize   int
}

// NewApp initializes the application with dependencies.
func NewApp(config *Config) (*App, error) {
	db, err := sqlx.Open("mysql", config.DBDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &App{
		DB:     db,
		Config: config,
	}, nil
}

// fetchNewUsersForUpdate retrieves a batch of new users from the database and marks them as in-process.
func (app *App) fetchNewUsersForUpdate(ctx context.Context, batchSize int) ([]string, error) {
	var emails = make([]string, 0)

	tx, err := app.DB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := "SELECT email FROM new_users WHERE email_sent = FALSE AND in_process = FALSE LIMIT ?"
	err = tx.SelectContext(ctx, &emails, query, batchSize)

	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if len(emails) == 0 {
		tx.Rollback()
		return nil, nil
	}

	updateQuery, args, err := sqlx.In("UPDATE new_users SET in_process = 1 WHERE email IN (?)", emails)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	updateQuery = app.DB.Rebind(updateQuery)
	_, err = tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update users: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return emails, nil
}

func (app *App) sendWelcomeEmail(email string) error {

	client, err := smtp.Dial("localhost:1025")
	if err != nil {
		log.Fatalf("Failed to connect to SMTP server: %v", err)
	}
	defer client.Close()

	sender := "no-reply@example.com"
	to := email
	subject := "Subject: Test email\r\n"
	body := "This is a test email sent from Go."

	// Set sender and recipient
	if err := client.Mail(sender); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}
	if err := client.Rcpt(to); err != nil {
		log.Fatalf("Failed to set recipient: %v", err)
	}

	// Send the email message
	w, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to start email data: %v", err)
	}
	defer w.Close()

	msg := []byte(subject + "\r\n" + body)
	_, err = w.Write(msg)
	if err != nil {
		log.Fatalf("Failed to write email message: %v", err)
	}

	log.Println("Email sent successfully")
	return nil
}

// processUsers processes a batch of users: sends emails and marks them as processed.
func (app *App) processUsers(ctx context.Context, emails []string) error {
	for _, email := range emails {
		log.Print("Sending Email", email, "... 💥")
		if err := app.sendWelcomeEmail(email); err != nil {
			log.Printf("Failed to send email to user %s: %v", email, err)
			return err
		}
		log.Print("Email Sent! 🚀")
	}

	updateQuery, args, err := sqlx.In("UPDATE new_users SET email_sent = TRUE, in_process = FALSE WHERE email IN (?)", emails)

	if err != nil {
		return err
	}
	updateQuery = app.DB.Rebind(updateQuery)
	_, err = app.DB.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update users: %v", err)
	}

	log.Print("Update Users Table! 🚀")
	return nil
}

// worker is a function that processes users in batches.
func (app *App) worker(ctx context.Context) {
	for {
		fetched_emails, err := app.fetchNewUsersForUpdate(ctx, app.Config.BatchSize)

		if err != nil || len(fetched_emails) == 0 {
			log.Print("No Users Fetched, Waiting... 😱")
			time.Sleep(5 * time.Second)
			continue
		}

		if err := app.processUsers(ctx, fetched_emails); err != nil {
			log.Printf("Failed to process users: %v", err)
		}
	}
}

func main() {
	seed()

	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // specify the config file type
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AutomaticEnv()          // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	config := &Config{
		DBDSN:       viper.GetString("DB_CONNECTION_STRING"),
		SMTPHost:    viper.GetString("SMTP_HOST"),
		SMTPPort:    viper.GetInt("SMTP_PORT"),
		SMTPUser:    viper.GetString("SMTP_USER"),
		SMTPPass:    viper.GetString("SMTP_PASS"),
		SMTPFrom:    viper.GetString("SMTP_FROM"),
		WorkerCount: viper.GetInt("WORKER_COUNT"),
		BatchSize:   viper.GetInt("BATCH_SIZE"),
	}

	app, err := NewApp(config)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	ctx := context.Background()
	quit := make(chan struct{})

	for w := 0; w < app.Config.WorkerCount; w++ {
		go app.worker(ctx)
	}

	// Keep the main function running.
	<-quit
}

// Config holds the application configuration.
type SeedConfig struct {
	DBDSN      string
	BatchSize  int
	NumBatches int
}

// NewConfig initializes the configuration using Viper.
func NewConfig() (*SeedConfig, error) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json")   // specify the config file type
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	viper.AutomaticEnv()          // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	config := &SeedConfig{
		DBDSN:      viper.GetString("DB_CONNECTION_STRING"),
		BatchSize:  viper.GetInt("BATCH_SIZE"),
		NumBatches: viper.GetInt("NUM_BATCHES"),
	}

	return config, nil
}

// generateRandomEmail generates a random email address.
func generateRandomEmail() string {
	domains := []string{"example.com", "test.com", "demo.com"}
	return fmt.Sprintf("user%d@%s", rand.Intn(1000000), domains[rand.Intn(len(domains))])
}

func insertDumpData(ctx context.Context, db *sqlx.DB) error {

	if _, err := db.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS new_users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            email VARCHAR(255) NOT NULL UNIQUE,
            email_sent BOOLEAN DEFAULT FALSE,
            in_process BOOLEAN DEFAULT FALSE
        )
    `); err != nil {
		return fmt.Errorf("failed to create new_users table: %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM new_users").Scan(&count); err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}

	if count == 0 {
		log.Print("Empty Table: Seeding it 🚀")
		for i := 0; i < 1000; i++ {
			var email = generateRandomEmail()
			query := "INSERT INTO new_users (email) VALUES (?)"
			if _, err := db.ExecContext(ctx, query, email); err != nil {
				return fmt.Errorf("failed to insert dump data: %v", err)
			}
		}
	}

	return nil
}

func seed() {
	config, err := NewConfig()

	if err != nil {
		log.Fatalf("Error initializing configuration: %v", err)
	}

	db, err := sqlx.Open("mysql", config.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	insertDumpData(context.Background(), db)

	log.Println("Data insertion completed.")
	defer db.Close()
}
