package main

import "os"

type (
	Config struct {
		port    string
		sqlite3 Sqlite3
		fileStorage string
		certFile string
		keyFile string
	}
	Sqlite3 struct {
		dsn    string
		driver string
	}
)

func getConfig() *Config {
	return &Config{
		port: getEnv("PORT", "8082"),
		fileStorage: getEnv("FILE_STORAGE", "./uploads"),
		certFile: getEnv("CERT_FILE", "server.crt"),
		keyFile: getEnv("KEY_FILE", "server.key"),
		sqlite3: Sqlite3{
			dsn:    getEnv("SQLITE3_DSN", "storage.db"),
			driver: getEnv("SQLITE3_DRIVER", "sqlite3"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
