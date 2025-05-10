package main

import (
	"context"
	"os"
	"time" // Necesario para el tipo de ejemplo ProjectFromAPI y la demo de TTL

	// Paquetes del proyecto Windraw
	// "github.com/mvialf/windraw/internal/app/window-api/models" // Eliminado ya que no se usa directamente en main
	"github.com/mvialf/windraw/internal/app/window-api/repositories"
	// "github.com/mvialf/windraw/internal/app/window-api/services" // Descomentar cuando tengas servicios
	"github.com/mvialf/windraw/internal/pkg/apiclient"
	"github.com/mvialf/windraw/internal/pkg/config"

	// Paquetes de terceros
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus" // Asegúrate de ejecutar: go get github.com/sirupsen/logrus
)

// ProjectFromAPI fue una estructura de ejemplo de tu main.go original.
type ProjectFromAPI struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	// 1. Configurar Logger (logrus)
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Iniciando aplicación Windraw (window-api)...")

	// 2. Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		logger.Warn("Advertencia: No se pudo cargar el archivo .env. Usando variables de entorno del sistema si existen.")
	}

	// 3. Cargar configuración de la aplicación
	cfg, err := config.LoadConfig() // O config.LoadConfig(".") si tu función espera un path
	if err != nil {
		logger.Fatalf("Error fatal al cargar la configuración: %v", err)
	}
	logger.Info("Configuración cargada exitosamente.")
	// Usa cfg.SupabaseAPI si ese es el nombre del campo en tu struct config.Config
	// Si es diferente (ej. SupabaseConfig), ajústalo.
	logger.Debugf("Supabase API URL desde config: %s", cfg.SupabaseAPI.BaseURL)

	// 4. Crear cliente API de Supabase
	// Usa &cfg.SupabaseAPI si ese es el campo en tu struct config.Config.
	// Ajusta si el campo tiene otro nombre (ej. cfg.SupabaseConfig)
	supaClient := apiclient.NewSupabaseClient(&cfg.SupabaseAPI)
	logger.Info("Cliente API de Supabase creado.")

	// --- Inicialización de Repositorios ---
	// 5. Crear Repositorio de Catálogo de Perfiles (con caché)
	// Asegúrate que la función NewSupabaseProfileCatalogRepository exista y esté exportada
	// en tu paquete repositories.
	profileRepo := repositories.NewSupabaseProfileCatalogRepository(supaClient, logger)
	logger.Info("Repositorio de Catálogo de Perfiles (con caché) creado.")

	ctx := context.Background()
	testProfileCatalog(ctx, profileRepo, logger)
	// testDirectProjectQuery(ctx, supaClient, logger) // Descomenta si quieres correr esta prueba

	logger.Info("Fase de inicialización y pruebas completada.")
	logger.Info("La aplicación está lista para la lógica del servidor HTTP (pendiente de implementación).")

	// ... (resto del código para el servidor HTTP futuro) ...
}

// testProfileCatalog encapsula la lógica de prueba para el catálogo de perfiles.
func testProfileCatalog(ctx context.Context, profileRepo repositories.ProfileCatalogRepository, logger *logrus.Logger) {
	logger.Info("--- Iniciando prueba del Catálogo de Perfiles ---")

	logger.Info("Llamada 1: profileRepo.GetAllProfiles()")
	profiles, err := profileRepo.GetAllProfiles(ctx)
	if err != nil {
		logger.Errorf("Error en profileRepo.GetAllProfiles() (Llamada 1): %v", err)
	} else {
		logger.Infof("Perfiles obtenidos (Llamada 1 - Esperado MISS): %d perfiles.", len(profiles))
		if len(profiles) > 0 {
			// Asumiendo que tu models.Profile tiene un campo SKU.
			// Si no, esta parte de la prueba necesitará ajustarse.
			skuToTest := profiles[0].SKU
			logger.Infof("Llamada 2: profileRepo.GetProfileBySKU(%s)", skuToTest)
			profile, errSku := profileRepo.GetProfileBySKU(ctx, skuToTest)
			if errSku != nil {
				logger.Errorf("Error en profileRepo.GetProfileBySKU(%s) (Llamada 2): %v", skuToTest, errSku)
			} else if profile != nil {
				logger.Infof("Perfil SKU '%s' obtenido (Llamada 2 - Esperado MISS): %s", skuToTest, profile.Description)
			} else {
				logger.Warnf("Perfil SKU '%s' no encontrado (Llamada 2).", skuToTest)
			}

			logger.Infof("Llamada 3: profileRepo.GetProfileBySKU(%s) - Esperando HIT", skuToTest)
			profile, errSku = profileRepo.GetProfileBySKU(ctx, skuToTest)
			if errSku != nil {
				logger.Errorf("Error en profileRepo.GetProfileBySKU(%s) (Llamada 3): %v", skuToTest, errSku)
			} else if profile != nil {
				logger.Infof("Perfil SKU '%s' obtenido (Llamada 3 - Esperado HIT): %s", skuToTest, profile.Description)
			} else {
				logger.Warnf("Perfil SKU '%s' no encontrado (Llamada 3).", skuToTest)
			}
		} else {
			logger.Info("No hay perfiles en el catálogo para probar GetProfileBySKU.")
		}
	}

	logger.Info("Llamada 4: profileRepo.GetAllProfiles() - Esperando HIT")
	_, err = profileRepo.GetAllProfiles(ctx) // Re-asignar a 'profiles' si necesitas el valor
	if err != nil {
		logger.Errorf("Error en profileRepo.GetAllProfiles() (Llamada 4): %v", err)
	} else {
		logger.Infof("Perfiles obtenidos (Llamada 4 - Esperado HIT): %d perfiles.", len(profiles)) // len(profiles) usará el valor de la Llamada 1
	}
	logger.Info("--- Prueba del Catálogo de Perfiles finalizada ---")
}

// testDirectProjectQuery (sin cambios, asumiendo que la adaptarás o eliminarás)
func testDirectProjectQuery(ctx context.Context, supaClient *apiclient.SupabaseClient, logger *logrus.Logger) {
	logger.Info("--- Iniciando prueba de consulta directa a 'projects' (legado) ---")
	path := "/rest/v1/projects"
	queryParams := "select=id,name,created_at&order=created_at.desc&limit=2"
	var projects []ProjectFromAPI

	err := supaClient.QueryData(path, queryParams, &projects)
	if err != nil {
		logger.Errorf("Error al consultar datos de Supabase ('projects'): %v", err)
	} else {
		if len(projects) > 0 {
			logger.Infof("Consulta directa a 'projects' exitosa. %d proyectos encontrados:", len(projects))
			for _, p := range projects {
				logger.Infof("  ID: %d, Name: %s, CreatedAt: %s", p.ID, p.Name, p.CreatedAt.Format(time.RFC3339))
			}
		} else {
			logger.Info("Consulta directa a 'projects' exitosa, pero no se encontraron proyectos.")
		}
	}

	logger.Info("Intentando consulta directa a 'tabla_inexistente'...")
	var noData []interface{}
	err = supaClient.QueryData("/rest/v1/tabla_inexistente", "select=*", &noData)
	if err != nil {
		logger.Warnf("Error (esperado) al consultar 'tabla_inexistente': %v", err)
	} else {
		logger.Info("Consulta a 'tabla_inexistente' no dio error, lo cual es inesperado si no existe.")
	}
	logger.Info("--- Prueba de consulta directa a 'projects' (legado) finalizada ---")
}
