package main

import (
	"context"
	"encoding/json" // Para imprimir structs de forma legible en el log
	"os"

	// "time" // No es estrictamente necesario en este main.go revisado, a menos que una prueba lo requiera

	// Paquetes del proyecto Windraw
	"github.com/mvialf/windraw/internal/app/window-api/models" // Necesario para crear el Element de prueba
	"github.com/mvialf/windraw/internal/app/window-api/repositories"
	"github.com/mvialf/windraw/internal/app/window-api/services" // Descomentado, ¡ahora tenemos servicios!
	"github.com/mvialf/windraw/internal/pkg/apiclient"
	"github.com/mvialf/windraw/internal/pkg/config"
	"github.com/mvialf/windraw/internal/pkg/constants" // Necesario para los roles, tipos, etc.

	// Paquetes de terceros
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// 1. Configurar Logger (logrus)
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel) // Cambiado a DebugLevel para ver más detalle durante el desarrollo

	logger.Info("Iniciando aplicación Windraw (window-api)...")

	// 2. Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		logger.Warn("Advertencia: No se pudo cargar el archivo .env. Usando variables de entorno del sistema si existen.")
	}

	// 3. Cargar configuración de la aplicación
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatalf("Error fatal al cargar la configuración: %v", err)
	}
	logger.Info("Configuración cargada exitosamente.")
	logger.Debugf("Supabase API URL desde config: %s", cfg.SupabaseAPI.BaseURL) // Asegúrate que cfg.SupabaseAPI.BaseURL sea el campo correcto

	// 4. Crear cliente API de Supabase
	supaClient := apiclient.NewSupabaseClient(&cfg.SupabaseAPI) // Asegúrate que &cfg.SupabaseAPI sea el campo correcto
	logger.Info("Cliente API de Supabase creado.")

	// --- Inicialización de Repositorios y Servicios ---
	// 5. Crear Repositorio de Catálogo
	profileRepo := repositories.NewSupabaseProfileCatalogRepository(supaClient, logger)
	logger.Info("Repositorio de Catálogo de Perfiles creado.")

	// 6. Crear Servicios
	profileSelectorSvc := services.NewProfileSelectorService(profileRepo, logger)
	logger.Info("Servicio ProfileSelector creado.")
	windowCalcSvc := services.NewWindowCalculationService(logger) // No tiene dependencias de repo directas por ahora
	logger.Info("Servicio WindowCalculation creado.")

	// --- Fase de Pruebas de la Nueva Lógica ---
	ctx := context.Background()
	testNewLogic(ctx, profileSelectorSvc, windowCalcSvc, logger) // Nueva función de prueba

	// testDirectProjectQuery(ctx, supaClient, logger) // Puedes mantenerla comentada si aún es útil para algo

	logger.Info("Fase de inicialización y pruebas completada.")
	logger.Info("La aplicación está lista para la lógica del servidor HTTP (pendiente de implementación).")

	// Aquí es donde configurarías y arrancarías tu servidor HTTP (ej. usando Gin, net/http)
	// router := gin.Default()
	// api.SetupRoutes(router, profileSelectorSvc, windowCalcSvc, projectManager) // Ejemplo de cómo podría ser
	// router.Run(":8080") // Por ejemplo
}

// testNewLogic demuestra cómo usar el ProfileSelectorService y WindowCalculationService.
func testNewLogic(ctx context.Context, selectorSvc *services.ProfileSelectorService, calcSvc *services.WindowCalculationService, logger *logrus.Logger) {
	log := logger.WithField("test_function", "testNewLogic")
	log.Info("--- Iniciando prueba de la Nueva Lógica de Selección y Cálculo ---")

	// Parámetros de ejemplo para una ventana corredera
	elementType := constants.TYPE_SLIDING
	// Deberás elegir un DesiredColorID que sepas que existe en tu tabla 'colors'
	// y que esté asociado a un sistema de perfiles 'Sliding' en 'system_available_colors'.
	// Supongamos que el color con ID 1 es "Blanco".
	desiredColorID := int64(1) // CAMBIA ESTO por un ID de color válido en tu BD

	// Definir las ranuras (element_part_role) necesarias para una ventana corredera de 2 hojas.
	// Esta lista la generaría una capa superior (ej. ElementBuilderService) en una app real.
	elementRanuras := []string{
		constants.ROLE_FRAME_PERIMETER_TOP_SLIDING,
		constants.ROLE_FRAME_PERIMETER_BOTTOM_SLIDING,
		constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING, // Para el lado izquierdo del marco
		constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING, // Para el lado derecho del marco (mismo rol)

		// Hoja 1
		constants.ROLE_WIND_JAMB_SIDE_SLIDING,    // Jamba lateral de la hoja 1
		constants.ROLE_WIND_JAMB_MEETING_SLIDING, // Jamba de encuentro de la hoja 1
		constants.ROLE_WIND_RAIL_TOP_SLIDING,     // Riel superior de la hoja 1
		constants.ROLE_WIND_RAIL_BOTTOM_SLIDING,  // Riel inferior de la hoja 1

		// Hoja 2
		constants.ROLE_WIND_JAMB_MEETING_SLIDING, // Jamba de encuentro de la hoja 2
		constants.ROLE_WIND_JAMB_SIDE_SLIDING,    // Jamba lateral de la hoja 2
		constants.ROLE_WIND_RAIL_TOP_SLIDING,     // Riel superior de la hoja 2
		constants.ROLE_WIND_RAIL_BOTTOM_SLIDING,  // Riel inferior de la hoja 2
		// Nota: ROLE_WIND_VERTICAL_OVERLAP_SLIDING no se pide explícitamente, el selector lo añade si es necesario.
	}

	selectorInput := services.SelectDefaultSystemAndProfilesInput{
		ElementType:    elementType,
		DesiredColorID: desiredColorID,
		ElementRanuras: elementRanuras,
	}

	log.Infof("Llamando a ProfileSelectorService.SelectDefaultSystemAndProfiles con ElementType: %s, ColorID: %d", elementType, desiredColorID)
	selectionOutput, err := selectorSvc.SelectDefaultSystemAndProfiles(ctx, selectorInput)
	if err != nil {
		log.WithError(err).Error("Error en SelectDefaultSystemAndProfiles")
		return
	}

	if len(selectionOutput.Errors) > 0 {
		log.Warnf("Se encontraron problemas durante la selección de perfiles: %v", selectionOutput.Errors)
		// Podrías decidir no continuar si hay errores críticos.
		// Por ahora, continuamos para ver qué se pudo seleccionar.
	}

	if selectionOutput.SelectedSystem.SystemID == 0 {
		log.Error("No se pudo seleccionar un sistema de perfiles.")
		return
	}
	if selectionOutput.SelectedColor.ColorID == 0 {
		log.Error("No se pudo seleccionar un color.")
		return
	}

	log.Infof("Sistema Seleccionado: %s (ID: %d)", selectionOutput.SelectedSystem.Name, selectionOutput.SelectedSystem.SystemID)
	log.Infof("Color Seleccionado: %s (ID: %d)", selectionOutput.SelectedColor.Name, selectionOutput.SelectedColor.ColorID)
	log.Info("Perfiles seleccionados para las ranuras:")
	for ranura, profileInfo := range selectionOutput.ProfilesForRanuras {
		log.Infof("  Ranura '%s': Perfil SKU '%s', ColorID %d (ItemSKU: %s, Precio: %.2f)",
			ranura,
			profileInfo.Profile.ProfileSKU,
			profileInfo.StockItem.ColorID, // Confirmando que es el color deseado
			profileInfo.StockItem.ItemSKU,
			profileInfo.StockItem.ProfilePrice,
		)
		// Imprimir más detalles si es necesario
		// profileJSON, _ := json.MarshalIndent(profileInfo.Profile, "    ", "  ")
		// log.Debugf("    Detalles del Perfil: %s", string(profileJSON))
	}

	// Si no hubo errores fatales y tenemos perfiles, intentamos el cálculo
	if len(selectionOutput.Errors) == 0 && len(selectionOutput.ProfilesForRanuras) > 0 {
		log.Info("--- Iniciando prueba de WindowCalculationService ---")

		// Crear un elemento de ejemplo. En una app real, esto vendría de la UI o de una definición.
		// Usar constructores de models si están listos, o crear manualmente.
		// Las dimensiones del elemento Width/Height deben ser realistas.
		exampleElement := &models.Element{
			ID:        "test-elem-001", // Deberías usar generateID() de tus modelos
			Width:     2000,            // Ancho total del elemento en mm
			Height:    1500,            // Alto total del elemento en mm
			Type:      elementType,
			Structure: "Ventana",              // Ejemplo
			Material:  constants.MATERIAL_PVC, // Ejemplo
			Frame: models.Frame{ // Marco inicializado, sus Details se llenarán si implementamos calculateFrame...
				Details: make(map[string]models.FrameDetail),
			},
			Winds: []models.Wind{ // Dos hojas para este ejemplo
				{ID: "wind-001", Kind: constants.KIND_SLIDING_WIND, Details: make(map[string]models.WindDetail)}, // Hoja 1
				{ID: "wind-002", Kind: constants.KIND_SLIDING_WIND, Details: make(map[string]models.WindDetail)}, // Hoja 2
			},
		}

		// Pre-poblar los Details con las posiciones esperadas para que el servicio de cálculo las llene.
		// Esto es crucial y simula lo que un ElementBuilderService haría.
		expectedWindPositions := []string{constants.POSITION_TOP, constants.POSITION_BOTTOM, constants.POSITION_LEFT, constants.POSITION_RIGHT}
		for i := range exampleElement.Winds {
			for _, pos := range expectedWindPositions {
				exampleElement.Winds[i].Details[pos] = models.WindDetail{Position: pos}
			}
		}
		// Podríamos hacer lo mismo para exampleElement.Frame.Details si `calculateFrameTerminatedDimensions` estuviera listo.

		calcInput := services.CalculateDetailsInput{
			Element:            exampleElement,
			System:             selectionOutput.SelectedSystem,
			ProfilesForRanuras: selectionOutput.ProfilesForRanuras,
		}

		log.Infof("Llamando a WindowCalculationService.CalculateDetails para Elemento ID: %s", exampleElement.ID)
		err = calcSvc.CalculateDetails(ctx, calcInput)
		if err != nil {
			log.WithError(err).Error("Error en CalculateDetails")
		} else {
			log.Info("CalculateDetails completado. Mostrando dimensiones calculadas para las hojas:")
			for i, wind := range exampleElement.Winds {
				if wind.Kind == constants.KIND_SLIDING_WIND {
					log.Infof("  Hoja %d (ID: %s):", i, wind.ID)
					for pos, detail := range wind.Details {
						log.Infof("    Posición '%s': Perfil SKU '%s', Dimensión %.2f, Ángulos (%.1f, %.1f)",
							pos, detail.ProfileSKU, detail.Dimension, detail.AngleLeft, detail.AngleRight)
					}
				}
			}
			// Para ver el elemento completo:
			elementJSON, _ := json.MarshalIndent(exampleElement, "  ", "  ")
			log.Debugf("Elemento completo después del cálculo:\n%s", string(elementJSON))
		}
	} else {
		log.Warn("No se procederá con el cálculo de dimensiones debido a errores en la selección de perfiles o falta de perfiles.")
	}

	log.Info("--- Prueba de la Nueva Lógica finalizada ---")
}

// testDirectProjectQuery (sin cambios, puedes mantenerla o eliminarla)
// type ProjectFromAPI struct { ... } // Definición ya estaba en tu main.go original
// func testDirectProjectQuery(ctx context.Context, supaClient *apiclient.SupabaseClient, logger *logrus.Logger) { ... }
