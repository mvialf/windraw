package main

import (
	"fmt"
	"log"

	// Reemplaza 'github.com/tu-usuario/tu-proyecto-ventanas' con tu nombre de módulo
	"github.com/tu-usuario/tu-proyecto-ventanas/internal/pkg/apiclient"
	"github.com/tu-usuario/tu-proyecto-ventanas/internal/pkg/config"
	// "github.com/tu-usuario/tu-proyecto-ventanas/internal/app/window-api/models"
)

// Define una estructura para el resultado de tu consulta de ejemplo
// Asegúrate de que los nombres de campo coincidan (o usa tags json `json:"nombre_columna"`)
// con las columnas que seleccionas de tu tabla.
type ProjectFromAPI struct {
	ID   int64  `json:"id"` // Asumiendo que 'id' es un número en la BD
	Name string `json:"name"`
	// Otros campos que quieras leer...
	// CreatedAt time.Time `json:"created_at"`
}

func main() {
	fmt.Println("Iniciando aplicación de fabricación de ventanas (modo API)...")

	// 1. Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}
	log.Println("Configuración cargada.")

	// 2. Crear cliente API de Supabase
	supaClient := apiclient.NewSupabaseClient(&cfg.SupabaseAPI)
	log.Println("Cliente API de Supabase creado.")

	// --- EJEMPLO DE CONSULTA DE DATOS ---
	log.Println("Intentando consultar datos de la tabla 'projects' (o la que definas)...")

	// Cambia "projects" por el nombre real de tu tabla
	// Cambia "id,name" por las columnas que quieres seleccionar
	// Puedes añadir filtros: "&name=eq.MiProyecto" o "&id=gte.10"
	// Ver documentación de PostgREST para filtros: https://postgrest.org/en/stable/api.html#filtering
	path := "/rest/v1/projects"              // Nombre de tu tabla
	queryParams := "select=id,name&limit=5"  // Selecciona columnas y limita resultados

	var projects []ProjectFromAPI // Un slice para almacenar los resultados

	err = supaClient.QueryData(path, queryParams, &projects)
	if err != nil {
		log.Printf("Error al consultar datos de Supabase: %v", err)
	} else {
		if len(projects) > 0 {
			log.Printf("Consulta exitosa. %d proyectos encontrados:\n", len(projects))
			for _, p := range projects {
				log.Printf("  ID: %d, Name: %s\n", p.ID, p.Name)
			}
		} else {
			log.Println("Consulta exitosa, pero no se encontraron proyectos (o la tabla está vacía / el filtro no coincide).")
		}
	}

	// Ejemplo con una tabla que podría no existir para ver el error
	log.Println("\nIntentando consultar datos de una tabla 'no_existe'...")
	var noData []interface{}
	err = supaClient.QueryData("/rest/v1/tabla_que_no_existe", "select=*", &noData)
	if err != nil {
		log.Printf("Error esperado al consultar 'tabla_que_no_existe': %v\n", err)
	}


	log.Println("La aplicación está lista para iniciar el servidor HTTP (lógica pendiente).")
	// ... lógica del servidor HTTP ...
}