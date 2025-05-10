package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv" // Asegúrate de tener esta dependencia: go get github.com/joho/godotenv
)

// APIConfig almacena la configuración para interactuar con la API de Supabase.
type APIConfig struct {
	BaseURL        string // Ej: https://ajviyshznobobicabwbw.supabase.co
	ServiceRoleKey string // Tu service_role secret key
}

// Config almacena toda la configuración de la aplicación.
type Config struct {
	SupabaseAPI APIConfig
	// APIPort     string // Descomenta si necesitas configurar el puerto de tu API
}

// LoadConfig carga la configuración desde variables de entorno (o un archivo .env).
func LoadConfig() (*Config, error) {
	// Carga variables desde .env en la raíz del proyecto
	err := godotenv.Load()
	if err != nil {
		// No es un error fatal si .env no existe (podría estar en un entorno de producción donde se usan variables de sistema)
		log.Println("Advertencia: No se pudo cargar el archivo .env. Se usarán variables de entorno del sistema si están disponibles. Error:", err)
	}

	cfg := &Config{}

	// Cargar configuración de la API de Supabase
	cfg.SupabaseAPI.BaseURL = getEnv("SUPABASE_API_URL", "")
	cfg.SupabaseAPI.ServiceRoleKey = getEnv("SUPABASE_SERVICE_KEY", "")

	// Validar que las variables necesarias estén presentes
	if cfg.SupabaseAPI.BaseURL == "" {
		return nil, fmt.Errorf("la variable de entorno SUPABASE_API_URL no está configurada")
	}
	if cfg.SupabaseAPI.ServiceRoleKey == "" {
		return nil, fmt.Errorf("la variable de entorno SUPABASE_SERVICE_KEY no está configurada")
	}

	// Cargar otras configuraciones
	// cfg.APIPort = getEnv("API_PORT", "8080")

	return cfg, nil
}

// getEnv lee una variable de entorno o devuelve un valor por defecto.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
