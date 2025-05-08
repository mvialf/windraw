package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// APIConfig almacena la configuración para interactuar con la API de Supabase.
type APIConfig struct {
	BaseURL        string // Ej: https://<project-id>.supabase.co
	ServiceRoleKey string // Tu service_role secret key
}

// Config almacena toda la configuración de la aplicación.
type Config struct {
	SupabaseAPI APIConfig
	// APIPort     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env. Error:", err)
	}

	cfg := &Config{}

	cfg.SupabaseAPI.BaseURL = getEnv("SUPABASE_API_URL", "")
	cfg.SupabaseAPI.ServiceRoleKey = getEnv("SUPABASE_SERVICE_KEY", "")

	if cfg.SupabaseAPI.BaseURL == "" {
		return nil, fmt.Errorf("la variable de entorno SUPABASE_API_URL no está configurada")
	}
	if cfg.SupabaseAPI.ServiceRoleKey == "" {
		return nil, fmt.Errorf("la variable de entorno SUPABASE_SERVICE_KEY no está configurada")
	}

	// cfg.APIPort = getEnv("API_PORT", "8080")

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}