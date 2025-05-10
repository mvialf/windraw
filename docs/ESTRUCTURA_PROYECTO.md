# Documentación de la Estructura del Proyecto Windraw

## Introducción

Este documento detalla la estructura de carpetas y archivos del proyecto Windraw. El objetivo es proporcionar una comprensión clara de la organización del código fuente, los archivos de configuración, la documentación y otros artefactos relevantes. Esta documentación es crucial para facilitar el mantenimiento, la colaboración y la incorporación de nuevos desarrolladores al proyecto.

Windraw es una aplicación Go diseñada para [Breve descripción del propósito principal de Windraw, por ejemplo: 'gestionar la fabricación y el diseño de ventanas y puertas']. Interactúa con un backend de Supabase para la persistencia de datos.

## Estructura General del Proyecto

A continuación, se muestra una representación en árbol de la estructura de directorios del proyecto Windraw. Los directorios y archivos sensibles o irrelevantes para la comprensión general (como `.git`, `.vscode`, archivos temporales) se han omitido.

```
Windraw/
├── .env (No versionado - Ejemplo de contenido si es relevante)
├── .gitignore
├── README.md
├── README2.md
├── cmd/
│   └── window-api/
│       └── main.go
├── configs/
│   └── config.dev.yaml
├── docs/
│   ├── ESTRUCTURA_PROYECTO.md (Este archivo)
│   ├── ejemplos_de_uso.md
│   └── resumen07_05.md
├── go.mod
├── go.sum
├── internal/
│   ├── app/
│   │   └── window-api/
│   │       ├── handlers/
│   │       │   └── http_handlers.go
│   │       ├── models/
│   │       │   ├── component_models.go
│   │       │   ├── element_models.go
│   │       │   └── project_models.go
│   │       ├── repositories/
│   │       │   └── (Vacío o placeholder para futuros repositorios)
│   │       └── services/
│   │           └── (Vacío o placeholder para futuros servicios)
│   └── pkg/
│       ├── apiclient/
│       │   └── supabase_client.go
│       ├── config/
│       │   └── config.go
│       └── constants/
│           └── constants.go
└── (Otros archivos o directorios principales si existen)
```

## Descripción Detallada de Carpetas y Archivos

### Archivos en el Directorio Raíz

#### `.gitignore`

Este archivo especifica los archivos y directorios que Git debe ignorar. Es crucial para evitar que archivos generados, dependencias locales, configuraciones sensibles o artefactos de compilación se incluyan en el control de versiones.

```gitignore
.idea/
.vscode/
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
*~
gopath/
vendor/
dependencies/
bin/
pkg/
# Env files
.env
.env.*
!.env.example
# Log files
*.log
# OS generated files
.DS_Store
Thumbs.db
```

#### `README.md`

El archivo README principal del proyecto. Generalmente contiene una descripción general del proyecto, instrucciones de instalación, cómo ejecutarlo, y otra información relevante para empezar.

```markdown
# windraw
```

#### `README2.md`

Un segundo archivo README. Su propósito actual es incierto ya que se encuentra vacío. Podría estar destinado a notas adicionales o documentación en progreso.

```markdown
Este archivo está actualmente vacío.
```

#### `go.mod`

Defines el módulo del proyecto Go (`github.com/mvialf/windraw`), sus dependencias directas y las versiones de Go requeridas. Es fundamental para la gestión de dependencias en Go y asegura compilaciones reproducibles.

```go
module github.com/mvialf/windraw

go 1.20

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.9.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

#### `go.sum`

Registra las sumas de verificación (checksums criptográficos) de las dependencias directas e indirectas utilizadas por el proyecto, tal como se especifican en `go.mod` y sus propias dependencias. Esto asegura la integridad y reproducibilidad de la compilación al garantizar que se utilicen exactamente las mismas versiones de los módulos.

```
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+cfYផ្សu_9j9uA/M=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/google/uuid v1.6.0 h1:cv3h5RMR2H7K3vERqWqFITLDRnLzO1x7Y3bKiDevfBA=
github.com/google/uuid v1.6.0/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/joho/godotenv v1.5.1 h1:G6B+H8F208Pj2Sj9XzMrPNA5qJzzJ0F6wZgR1i2Q84E=
github.com/joho/godotenv v1.5.1/go.mod h1:q3/xH8MAj26K5yLqUvb4z+gPZLeKzBQCf0iX0qgZqfI=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NG8rlTnDofKTM7xIcIuD/nNl1Yy0hNEwBHbY=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1j/nOPvVMqS2xUuJ6DOKfGesZ6cK2s=
github.com/sirupsen/logrus v1.9.3 h1:peRkK+htD3uxllqNBUTuHIftP3/VzDqF1U9kKNB5YdI=
github.com/sirupsen/logrus v1.9.3/go.mod h1:N8fcDxBiAZ3YOpD2+l+94R6Zl4a0J+N0hKDQBtM44uQ=
github.com/stretchr/objx v0.5.2/go.mod h1:CHkGj5zUz3KFs3T45xTjL5QYeCEQDU6D0RaKzLAXmY4=
github.com/stretchr/testify v1.9.0 h1:4JTRV9wEgiEwJ/BftRUEERqB5pjsiJwVojxVf4o6Pk4=
github.com/stretchr/testify v1.9.0/go.mod h1:t8hU9iJjXQQxRIIqk99pgV8n0zD0V9QVeCOG7LVnJ+M=
golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 h1:P29w7jL0Hl0V5OyP+sRnP3XW7WprrSjch1ZmuI3vWlM=
golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8/go.mod h1:oPkhp1MJrh7nUepCBck5+mAzfO9JrbApNNgaTdGDITg=
gopkg.in/yaml.v2 v2.4.0 h1:D8xgwECY7SRAWxzOqk1h7gaA+8tE2TzTjMfmPpcA1/Q=
gopkg.in/yaml.v2 v2.4.0/go.mod h1:KlV9M+X7E8Z2k3JJr7HlWHpGnYnZ+xR3mrhZrAXXrF8=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzJygEWbeMflLKCMXLqg8N+G+w2V3P4oQILKHA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZDabD6gA2SM=

```

### Directorio `cmd`

El directorio `cmd` contiene los puntos de entrada principales de la aplicación o aplicaciones del proyecto. Cada subdirectorio dentro de `cmd` típicamente corresponde a un ejecutable diferente.

#### `cmd/window-api/`

Este subdirectorio aloja el código fuente para el ejecutable de la API de Windraw.

##### `cmd/window-api/main.go`

Este es el archivo principal que inicia la aplicación `window-api`. Sus responsabilidades incluyen:
*   Cargar la configuración de la aplicación.
*   Configurar e inicializar el cliente de Supabase.
*   (Eventualmente) Configurar el enrutador HTTP, los servicios, repositorios y arrancar el servidor.
*   Actualmente, también contiene código de ejemplo para probar la conexión y consulta a Supabase.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mvialf/windraw/internal/pkg/apiclient"
	"github.com/mvialf/windraw/internal/pkg/config"
	// Asegúrate de tener el paquete de modelos y constantes si son necesarios aquí
	// "github.com/mvialf/windraw/internal/app/window-api/models"
	// "github.com/mvialf/windraw/internal/pkg/constants"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load() // Carga desde .env en el directorio actual por defecto
	if err != nil {
		log.Println("Advertencia: No se pudo cargar el archivo .env. Asegúrate de que exista si no estás usando variables de entorno del sistema.")
	}

	// Cargar configuración de la aplicación
	cfg, err := config.LoadConfig(".") // Asume que config.dev.yaml o similar está en el directorio raíz o especificado
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// Crear cliente de Supabase
	supaClient, err := apiclient.NewSupabaseClient(cfg.Supabase.BaseURL, cfg.Supabase.ServiceRoleKey)
	if err != nil {
		log.Fatalf("Error al crear el cliente de Supabase: %v", err)
	}

	fmt.Println("Cliente de Supabase creado exitosamente.")
	fmt.Printf("Intentando consultar la tabla 'projects' en: %s\n", cfg.Supabase.BaseURL)

	// Ejemplo de cómo usar el cliente para obtener datos de una tabla
	// Reemplaza "projects" con el nombre real de tu tabla y asegúrate de que el tipo de destino sea apropiado.
	// El tipo `[]map[string]interface{}` es genérico para inspeccionar datos desconocidos.
	var projectsResponse []map[string]interface{}

	// Crear un contexto (puede ser útil para timeouts o cancelación)
	ctx := context.Background()

	// Realizar la petición GET a la tabla "projects"
	// Nota: Supabase usa PostgREST, por lo que las tablas son expuestas como endpoints REST.
	// Asegúrate de que la tabla 'projects' exista en tu base de datos Supabase.
	/*
		SQL para crear la tabla 'projects' (ejemplo básico):

		CREATE TABLE projects (
		    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    name TEXT NOT NULL,
		    created_at TIMESTAMPTZ DEFAULT now(),
		    client_name TEXT,
		    client_email TEXT,
		    iva_rate NUMERIC(5,2) DEFAULT 0.19 -- Ejemplo de IVA, ajustar según necesidad
		);

		-- También necesitarás configurar las políticas de RLS (Row Level Security) en Supabase
		-- para permitir el acceso a esta tabla, por ejemplo, para el rol 'service_role'.
		-- Ejemplo para permitir lectura a todos (¡NO RECOMENDADO PARA PRODUCCIÓN SIN AUTENTICACIÓN!):
		-- ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
		-- CREATE POLICY "Allow all read access to projects" ON projects FOR SELECT USING (true);

		-- Para permitir al service_role (usado por la ServiceRoleKey) leer y escribir:
		-- CREATE POLICY "Allow service_role full access to projects" ON projects
		-- FOR ALL
		-- USING (true) -- O (auth.role() = 'service_role') si quieres ser más explícito
		-- WITH CHECK (true);
	*/

	err = supaClient.Get(ctx, "projects", &projectsResponse, nil) // Sin parámetros de consulta adicionales
	if err != nil {
		log.Fatalf("Error al obtener datos de la tabla 'projects': %v", err)
		os.Exit(1) // Salir si hay un error para que se vea claramente en la ejecución
	}

	if len(projectsResponse) == 0 {
		fmt.Println("No se encontraron proyectos o la tabla está vacía.")
	} else {
		fmt.Printf("Proyectos encontrados (%d):\n", len(projectsResponse))
		for i, project := range projectsResponse {
			fmt.Printf("  Proyecto %d: ID=%v, Nombre=%v\n", i+1, project["id"], project["name"])
			// Puedes imprimir más campos según la estructura de tu tabla 'projects'
		}
	}

	fmt.Println("Ejecución de prueba de Supabase finalizada.")

	// Aquí iría la inicialización del servidor HTTP, servicios, etc.
	// Ejemplo: router := mux.NewRouter()
	// api.RegisterRoutes(router, projectService, elementService)
	// log.Println("Servidor iniciado en el puerto 8080")
	// log.Fatal(http.ListenAndServe(":8080", router))
}

---

### Directorio `configs`

Este directorio está destinado a almacenar los archivos de configuración de la aplicación. Es una práctica común tener diferentes archivos para distintos entornos (ej. `config.dev.yaml`, `config.prod.yaml`, `config.test.yaml`).

#### `configs/config.dev.yaml`

Archivo de configuración específico para el entorno de desarrollo. Actualmente, este archivo está vacío, pero podría contener configuraciones como niveles de log, parámetros específicos de desarrollo para bases de datos, o flags para funcionalidades en prueba.

```yaml
# Este archivo está actualmente vacío.
# Ejemplo de lo que podría contener:
# server:
#   port: 8081
#   debug_mode: true
# log_level: debug

# supabase:
#   timeout_seconds: 10

```

### Directorio `docs`

Este directorio contiene toda la documentación relacionada con el proyecto Windraw. Es fundamental para entender la arquitectura, el uso y los detalles de implementación del sistema.

#### `docs/ESTRUCTURA_PROYECTO.md`

Este mismo archivo. Proporciona una descripción detallada de la estructura de carpetas y archivos del proyecto Windraw, incluyendo el propósito de cada componente y el contenido de los archivos de código relevantes.

#### `docs/ejemplos_de_uso.md`

Este documento muestra cómo instanciar y utilizar las estructuras y funciones principales del backend de fabricación de ventanas. Incluye requisitos previos y un ejemplo completo en Go que simula la creación de un proyecto, un elemento con marco y hoja, la asignación de perfiles y el cálculo de detalles.

```markdown
# Ejemplos de Uso del Backend de Fabricación de Ventanas

Este documento muestra cómo instanciar y utilizar las estructuras y funciones principales del backend de fabricación de ventanas, asumiendo que el código Go está organizado en los paquetes `models` y `constants` como se definió previamente.

**Nota Importante:** Estos ejemplos se ejecutarían típicamennte dentro de la función `main` (para pruebas iniciales), o más realisticamente, la lógica de creación y manipulación sería orquestada por *servicios* que a su vez son llamados por *handlers* de API.

## Requisitos Previos

Asegúrate de que tu `go.mod` esté inicializado correctamente (ej. `go mod init github.com/tu-usuario/tu-proyecto-ventanas`) y que los paquetes `models` y `constants` existan en las rutas:

*   `github.com/tu-usuario/tu-proyecto-ventanas/internal/app/window-api/models`
*   `github.com/tu-usuario/tu-proyecto-ventanas/internal/pkg/constants`

(Reemplaza `github.com/tu-usuario/tu-proyecto-ventanas` con el nombre real de tu módulo).

## Ejemplo Completo

Este ejemplo simula la creación de un proyecto, un elemento con un marco y una hoja, y la asignación de perfiles y cálculo de detalles.

```go
package main

import (
	"fmt"
	// Asegúrate de que la ruta de importación coincida con tu nombre de módulo
	"github.com/mvialf/windraw/internal/app/window-api/models"
	"github.com/mvialf/windraw/internal/pkg/constants"
)

func main() {
	fmt.Println("--- Iniciando Ejemplo de Uso ---")

	// 1. Crear Información de Contacto
	contactInfo := models.Contact{
		Type:     models.PERSON, // Usando constante de models.ContactType
		Name:     "Ana López",
		Email:    "ana.lopez@example.com",
		Phone:    "555-1234",
		Address:  "Calle Falsa 123",
		District: "Centro",
		City:     "Ciudad Ejemplo",
	}
	fmt.Printf("Contacto creado: %s\n", contactInfo.Name)

	// 2. Crear un Proyecto
	// No se añaden costos ni componentes inicialmente
	proyecto, err := models.NewProject(
		"Remodelación Cocina Principal",
		contactInfo,
		nil, // costs
		nil, // components
		constants.IVA_RATE,
	)
	if err != nil {
		fmt.Printf("Error creando proyecto: %v\n", err)
		return
	}
	fmt.Printf("Proyecto creado: ID %s, Nombre: '%s', IVA: %.2f%%\n", proyecto.ID, proyecto.Name, proyecto.IvaRate)

	// 3. Definir parámetros comunes para Elementos, Marcos y Hojas
	//    Estos vendrían de la lógica de la aplicación o del input del usuario.

	// Posiciones estándar para un marco o hoja rectangular
	defaultPositions := []string{
		constants.POSITION_LEFT,
		constants.POSITION_RIGHT,
		constants.POSITION_TOP,
		constants.POSITION_BOTTOM,
	}

	// Listas de opciones válidas (que los constructores usarán para validar)
	validMaterials := []string{constants.MATERIAL_PVC, constants.MATERIAL_ALUMINIO}
	validElementTypes := []string{constants.TYPE_SLIDING, constants.TYPE_CASEMENT}
	validStructures := []string{constants.STRUCTURE_VENTANA, constants.STRUCTURE_PUERTA}
	validWindKinds := []string{
		constants.WIND_KIND_SLIDING_MOVIL, constants.WIND_KIND_SLIDING_FIXED,
		constants.WIND_KIND_FIXED, constants.WIND_KIND_CASEMENT,
	}
	validWindCutTypes := []string{
		constants.CUT_SQUARE_WIND, constants.CUT_ANGLE_WIND,
		constants.CUT_VERTICAL_OVERLAP_WIND, constants.CUT_HORIZONTAL_OVERLAP_WIND,
	}


	// 4. Crear un Elemento (Ventana Corredera de PVC)
	elementoVentana, err := models.NewElement(
		1500, // width
		1200, // height
		constants.MATERIAL_PVC,
		constants.TYPE_SLIDING,
		constants.STRUCTURE_VENTANA,
		constants.GEOMETRY_RECTANGULAR, // defaultFrameGeometry
		constants.CUT_ANGLE,            // defaultFrameCutType
		defaultPositions,               // defaultFramePositions
		validMaterials,
		validElementTypes,
		validStructures,
	)
	if err != nil {
		fmt.Printf("Error creando elementoVentana: %v\n", err)
		return
	}
	fmt.Printf("Elemento Ventana Creado: ID %s, Material: %s, Tipo: %s\n",
		elementoVentana.ID, elementoVentana.Material, elementoVentana.Type)
	fmt.Printf("  Dimensiones Marco: %dx%d, Geometría: %s, Tipo Corte Marco: %s\n",
		elementoVentana.Frame.Width, elementoVentana.Frame.Height, elementoVentana.Frame.Geometry, elementoVentana.Frame.CutType)

	// 5. Asignar Perfiles al Marco del elementoVentana
	//    (Los SKUs y colores son ejemplos, vendrían de un catálogo/configuración)
	fmt.Println("\n--- Asignando perfiles al marco de elementoVentana ---")
	err = elementoVentana.Frame.SetFrameProfile(constants.POSITION_TOP, "PVC-M-TOP-001", "Blanco Nieve")
	if err != nil { fmt.Printf("Error SetFrameProfile TOP: %v\n", err); return }
	err = elementoVentana.Frame.SetFrameProfile(constants.POSITION_BOTTOM, "PVC-M-BOT-002", "Blanco Nieve")
	if err != nil { fmt.Printf("Error SetFrameProfile BOTTOM: %v\n", err); return }
	err = elementoVentana.Frame.SetFrameProfile(constants.POSITION_LEFT, "PVC-M-SIDE-003", "Blanco Nieve")
	if err != nil { fmt.Printf("Error SetFrameProfile LEFT: %v\n", err); return }
	err = elementoVentana.Frame.SetFrameProfile(constants.POSITION_RIGHT, "PVC-M-SIDE-003", "Blanco Nieve") // Mismo perfil para ambos lados
	if err != nil { fmt.Printf("Error SetFrameProfile RIGHT: %v\n", err); return }

	// 6. Calcular Detalles del Marco del elementoVentana
	//    El '70.0' es un ancho de perfil de ejemplo para el cálculo.
	//    En una app real, este valor se obtendría del catálogo de perfiles basado en el SKU.
	fmt.Println("\n--- Calculando detalles del marco de elementoVentana ---")
	err = elementoVentana.Frame.CalculateFrameDetails(
		70.0, // anchoPerfilEjemplo
		constants.POSITION_LEFT,
		constants.POSITION_RIGHT,
		constants.POSITION_TOP,
		constants.POSITION_BOTTOM,
		constants.CUT_ANGLE, // valor para cutTypeAngle
		constants.CUT_SQUARE, // valor para cutTypeSquare
	)
	if err != nil {
		fmt.Printf("Error calculando detalles del marco: %v\n", err)
		return
	}
	fmt.Println("Detalles del marco calculados. Ejemplo (Top):")
	if detail, ok := elementoVentana.Frame.Details[constants.POSITION_TOP]; ok {
		fmt.Printf("  Pos: %s, SKU: %s, Dim: %dmm, Ángulos: %.1f°/%.1f°\n",
			detail.Position, detail.ProfileSKU, detail.Dimension, detail.AngleLeft, detail.AngleRight)
	}

	// 7. Crear una Hoja Corredera para elementoVentana
	//    Las dimensiones de la hoja son menores que las del marco.
	//    El tipo de corte podría ser diferente (ej. para solapes en correderas).
	fmt.Println("\n--- Creando Hoja Corredera ---")
	hojaCorredera1, err := models.NewWind(
		"Hoja Móvil Izquierda",             // name
		constants.WIND_KIND_SLIDING_MOVIL, // kind
		730,                               // width (ej: (AnchoMarco - Holguras) / 2)
		1120,                              // height (ej: AltoMarco - HolgurasMarcoSupInf - AltoPerfilHojaSupInf)
		constants.CUT_VERTICAL_OVERLAP_WIND, // cutType (para solape vertical)
		constants.WIND_STATUS_ACTIVE,      // status
		defaultPositions,                  // positions
		validWindKinds,
		validWindCutTypes,
	)
	if err != nil {
		fmt.Printf("Error creando hojaCorredera1: %v\n", err)
		return
	}
	fmt.Printf("Hoja Corredera Creada: ID %s, Nombre: '%s', Tipo: %s, Corte: %s\n",
		hojaCorredera1.ID, hojaCorredera1.Name, hojaCorredera1.Kind, hojaCorredera1.CutType)

	// 8. Asignar Perfiles a la Hoja Corredera
	fmt.Println("\n--- Asignando perfiles a hojaCorredera1 ---")
	err = hojaCorredera1.SetWindProfile(constants.POSITION_TOP, "PVC-H-TOP-SLD-00A", "Blanco Nieve")
	if err != nil { fmt.Printf("Error SetWindProfile Hoja TOP: %v\n", err); return }
	err = hojaCorredera1.SetWindProfile(constants.POSITION_BOTTOM, "PVC-H-BOT-SLD-00B", "Blanco Nieve")
	if err != nil { fmt.Printf("Error SetWindProfile Hoja BOTTOM: %v\n", err); return }
	// Para solape vertical, los perfiles verticales son clave
	err = hojaCorredera1.SetWindProfile(constants.POSITION_LEFT, "PVC-H-SIDE-SLD-INT-00C", "Blanco Nieve") // Perfil interior/de encuentro
	if err != nil { fmt.Printf("Error SetWindProfile Hoja LEFT: %v\n", err); return }
	err = hojaCorredera1.SetWindProfile(constants.POSITION_RIGHT, "PVC-H-SIDE-SLD-EXT-00D", "Blanco Nieve") // Perfil exterior/lateral
	if err != nil { fmt.Printf("Error SetWindProfile Hoja RIGHT: %v\n", err); return }


	// 9. Calcular Detalles de la Hoja Corredera
	//    El '60.0' es un ancho de perfil de hoja de ejemplo.
	fmt.Println("\n--- Calculando detalles de hojaCorredera1 ---")
	err = hojaCorredera1.CalculateWindDetails(
		60.0, // anchoPerfilHojaEjemplo
		constants.POSITION_LEFT,
		constants.POSITION_RIGHT,
		constants.POSITION_TOP,
		constants.POSITION_BOTTOM,
		constants.CUT_ANGLE_WIND,
		constants.CUT_SQUARE_WIND,
		constants.CUT_VERTICAL_OVERLAP_WIND,
	)
	if err != nil {
		fmt.Printf("Error calculando detalles de la hoja: %v\n", err)
		return
	}
	fmt.Println("Detalles de la hoja calculados. Ejemplo (Left - Vertical con solape):")
	if detail, ok := hojaCorredera1.Details[constants.POSITION_LEFT]; ok {
		fmt.Printf("  Pos: %s, SKU: %s, Dim: %dmm, Ángulos: %.1f°/%.1f°\n",
			detail.Position, detail.ProfileSKU, detail.Dimension, detail.AngleLeft, detail.AngleRight)
	}


	// 10. Añadir la Hoja al Elemento
	err = elementoVentana.AddWind(*hojaCorredera1)
	if err != nil {
		fmt.Printf("Error añadiendo hojaCorredera1 a elementoVentana: %v\n", err)
		return
	}
	fmt.Printf("\nElemento Ventana ahora tiene %d hoja(s).\n", len(elementoVentana.Winds))

	// 11. Crear un Módulo y un Componente para organizar los elementos
	//     (Aunque aquí solo tenemos un elemento)
	moduloCocina := models.Module{
		ID:       "MOD-COCINA-01",
		Elements: []models.Element{*elementoVentana},
		// AuxProfile: &models.AuxProfile{ID: "UNION-90DEG", Angle: 90.0}, // Ejemplo si se uniera a otro módulo
	}

	componenteVentanas := models.Component{
		ID:      "COMP-VENTANAS-COCINA",
		Modules: []models.Module{moduloCocina},
	}

	// 12. Añadir el Componente al Proyecto
	proyecto.AddComponent(componenteVentanas)
	fmt.Printf("Proyecto '%s' ahora tiene %d componente(s).\n", proyecto.Name, len(proyecto.Components))
	if len(proyecto.Components) > 0 && len(proyecto.Components[0].Modules) > 0 {
		fmt.Printf("  El primer componente tiene %d módulo(s), y el primer módulo tiene %d elemento(s).\n",
			len(proyecto.Components[0].Modules), len(proyecto.Components[0].Modules[0].Elements))
	}

	fmt.Println("\n--- Ejemplo de Uso Finalizado ---")
}
```

#### `docs/resumen07_05.md`

Este archivo parece ser un resumen o borrador de la arquitectura y diseño del sistema, posiblemente creado en una fecha específica (07/05). Detalla las capas lógicas (Modelos, Constantes, Repositorios, Servicios, Handlers), el flujo de datos típico para crear un elemento, las tecnologías clave, conceptos del dominio (Project, Element, Frame, Wind, etc.), el proceso de cálculo simplificado, constantes importantes, y posibles próximos pasos o áreas de contribución.

```markdown
### 3.2. Capas Lógicas

1.  **Capa de Modelos/Dominio (`internal/app/window-api/models/`)**:
    *   Contiene las definiciones de las estructuras de datos (`structs`) que representan las entidades centrales del negocio: `Project`, `Element`, `Frame`, `Wind`, `FrameDetail`, `WindDetail`, `Contact`, `Component`, `Module`, etc.
    *   Incluye funciones constructoras (`New*`) para estas entidades, asegurando la validación inicial de datos.
    *   Puede contener métodos asociados a estas estructuras para operaciones intrínsecas (ej. `AddComponent` a un `Project`).
    *   La lógica de cálculo más compleja (como `CalculateFrameDetails`) reside temporalmente aquí como métodos, pero se moverá progresivamente a la capa de servicios.

2.  **Capa de Constantes (`internal/pkg/constants/`)**:
    *   Define todas las constantes globales usadas en la aplicación para evitar "magic strings/numbers" y facilitar el mantenimiento. Esto incluye tipos de materiales, tipos de elementos, tipos de corte, posiciones, etc.

3.  **Capa de Repositorios (`internal/app/window-api/repositories/`)**:
    *   Define las **interfaces** (puertos en terminología hexagonal) que describen las operaciones de persistencia de datos (CRUD) para las entidades del dominio (ej. `ProjectRepository`, `ElementRepository`).
    *   Contiene las **implementaciones** (adaptadores) de estas interfaces, específicamente para interactuar con **Supabase (PostgreSQL)**. Esta capa abstrae los detalles de la base de datos del resto de la aplicación.

4.  **Capa de Servicios (`internal/app/window-api/services/`)**:
    *   Contiene la **lógica de negocio principal** y la orquestación de operaciones complejas.
    *   Los servicios utilizan los repositorios (a través de sus interfaces) para acceder a los datos y manipulan las entidades del dominio (modelos).
    *   Ejemplos: `ProjectService` para gestionar la creación de proyectos y sus componentes; `WindowCalculationService` para realizar los cálculos detallados de despiece de perfiles, holguras, y aplicación de reglas del catálogo.
    *   Esta capa es independiente de la capa de presentación (HTTP).

5.  **Capa de Manejadores/API (`internal/app/window-api/handlers/`)**:
    *   Responsable de manejar las peticiones HTTP entrantes de la API REST.
    *   Valida los datos de entrada de la API, decodifica JSON, etc.
    *   Llama a los métodos apropiados de la capa de servicios para ejecutar la lógica de negocio.
    *   Formatea las respuestas (generalmente JSON) y gestiona los códigos de estado HTTP.
    *   Se utiliza la biblioteca estándar `net/http` de Go, potencialmente complementada con un router ligero como `chi` para un manejo de rutas más avanzado.

6.  **Punto de Entrada (`cmd/window-api/main.go`)**:
    *   Inicializa la aplicación: carga la configuración, establece la conexión a la base de datos (Supabase), instancia los repositorios, servicios y handlers.
    *   Configura el router HTTP y arranca el servidor.

### 3.3. Flujo de Datos Típico (Ejemplo: Crear un Elemento)

1.  El cliente (Frontend React) envía una petición HTTP POST a un endpoint como `/api/v1/projects/{projectID}/elements` con los datos del nuevo elemento en formato JSON.
2.  El `ElementHandler` en la capa de `handlers` recibe la petición.
3.  El handler decodifica el JSON, valida los datos de la API.
4.  El handler llama a un método como `CreateElement(projectID, elementData)` en el `ElementService` (o `ProjectService`).
5.  El `ElementService`:
    *   Puede interactuar con `ProjectRepository` para verificar que el proyecto existe.
    *   Utiliza el constructor `models.NewElement()` para crear una instancia del elemento, aplicando validaciones de dominio.
    *   Si se proporcionan perfiles, podría invocar un `WindowCalculationService` para calcular los detalles del marco y las hojas (este servicio accedería a un catálogo de perfiles, que es una dependencia externa).
    *   Llama a `ElementRepository.Save(elemento)` para persistir el nuevo elemento en Supabase.
6.  El servicio devuelve el elemento creado (o un error) al handler.
7.  El handler formatea la respuesta (el elemento creado o un mensaje de error) como JSON y la envía de vuelta al cliente con el código HTTP apropiado (ej. 201 Created o 400 Bad Request).

## 4. Tecnologías Clave

*   **Lenguaje de Programación:** Go (Golang)
*   **Base de Datos:** Supabase (PostgreSQL)
    *   Se accede principalmente mediante conexión directa a PostgreSQL desde el backend Go utilizando el driver `database/sql` y `pgx`.
    *   Supabase Auth puede ser utilizado para la gestión de usuarios.
*   **API:** RESTful sobre HTTP, utilizando JSON como formato de intercambio de datos.
*   **Enrutamiento HTTP:** Inicialmente con la biblioteca estándar `net/http`, con posibilidad de integrar un router ligero como `chi` si es necesario.
*   **Gestión de Dependencias:** Go Modules.
*   **Frontend (Previsto):** React.

## 5. Funcionamiento Básico (Conceptos del Dominio)

El sistema se basa en una jerarquía de entidades:

*   **`Project`**: La entidad de más alto nivel. Contiene información del cliente, costos generales y una lista de `Component`s.
    *   Atributos: `ID`, `Name`, `CreatedAt`, `Contact`, `Costs`, `Components`, `IvaRate`.
*   **`Contact`**: Información del cliente (persona o empresa).
*   **`ProjectCost`**: Costos adicionales del proyecto (fijos o porcentuales).
*   **`Component`**: Una agrupación lógica de `Module`s dentro de un proyecto (ej. "Ventanas Planta Baja", "Puertas Terraza").
    *   Atributos: `ID`, `Modules`.
*   **`Module`**: Una agrupación de `Element`s que pueden estar unidos por perfiles auxiliares (ej. una ventana fija unida a una puerta corredera).
    *   Atributos: `ID`, `Elements`, `AuxProfile`.
*   **`Element`**: La unidad funcional principal, una ventana o puerta individual.
    *   Atributos: `ID`, `Width`, `Height`, `Material`, `Type` (tipología), `Structure` (ventana/puerta), `Frame`, `Winds` (lista de hojas), `Properties`.
*   **`Frame`**: El marco perimetral del `Element`.
    *   Atributos: `Name`, `Geometry`, `Width`, `Height`, `CutType`, `Details` (mapa de `FrameDetail`).
*   **`Wind` (Hoja)**: Un panel móvil o fijo dentro de un `Element`.
    *   Atributos: `ID`, `Name`, `Kind` (tipo de hoja), `Status`, `OpeningSide`, `OpeningDirection`, `Width`, `Height`, `CutType`, `Details` (mapa de `WindDetail`).
*   **`FrameDetail` / `WindDetail`**: Describe una pieza individual de perfil para un marco o una hoja.
    *   Atributos: `Position` (Top, Bottom, Left, Right), `ProfileSKU`, `Color`, `Dimension` (longitud de corte), `AngleLeft`, `AngleRight`, `ReinforcedUsed`, `ReinforcedSKU`.

### 5.1. Proceso de Cálculo (Simplificado)

1.  Se crea un `Element` con sus dimensiones generales, material, tipo y estructura.
2.  Se crea un `Frame` para el elemento, con su geometría y tipo de corte.
3.  Se asignan perfiles (mediante `ProfileSKU` y `Color`) a cada `Position` del `Frame` usando `SetFrameProfile()`.
4.  Se invoca `CalculateFrameDetails()`:
    *   Esta función (eventualmente un servicio) consultaría un catálogo de perfiles (no implementado aún).
    *   Basándose en el `ProfileSKU`, las dimensiones del marco, el tipo de corte y las reglas del catálogo, calcula:
        *   `Dimension` (longitud de corte) para cada `FrameDetail`.
        *   `AngleLeft` y `AngleRight`.
        *   Si se necesita `ReinforcedUsed` y el `ReinforcedSKU`.
5.  Se crean `Wind`s (hojas) para el elemento, con sus dimensiones, tipo de hoja y tipo de corte.
6.  Se asignan perfiles a cada `Position` de cada `Wind` usando `SetWindProfile()`.
7.  Se invoca `CalculateWindDetails()` para cada hoja, siguiendo una lógica similar a la del marco pero con reglas específicas para hojas y sus tipos de corte (ej. solapes).
8.  Las hojas se añaden al `Element`.

## 6. Constantes Importantes

El sistema utiliza un conjunto de constantes definidas en `internal/pkg/constants/` para estandarizar valores y mejorar la legibilidad:

*   `MATERIAL_*` (ej. `MATERIAL_PVC`, `MATERIAL_ALUMINIO`)
*   `TYPE_*` (ej. `TYPE_SLIDING`, `TYPE_CASEMENT`)
*   `PROFILE_TYPE_*` (ej. `PROFILE_TYPE_SLIDING_FRAME`)
*   `STRUCTURE_*` (ej. `STRUCTURE_VENTANA`)
*   `GEOMETRY_*` (ej. `GEOMETRY_RECTANGULAR`)
*   `POSITION_*` (ej. `POSITION_LEFT`, `POSITION_TOP`)
*   `CUT_*` (ej. `CUT_SQUARE`, `CUT_ANGLE_WIND`)
*   `WIND_KIND_*` (ej. `WIND_KIND_SLIDING_MOVIL`)
*   `WIND_STATUS_*` (ej. `WIND_STATUS_ACTIVE`)
*   `OPENING_SIDE_*` / `OPENING_DIRECTION_*`

## 7. Cómo Contribuir / Próximos Pasos

(Esta sección es más para un proyecto real, pero es bueno pensar en ella)

*   **Implementación de la Capa de Repositorios:** Conectar con Supabase y escribir las implementaciones para las interfaces de repositorio.
*   **Desarrollo de la Capa de Servicios:** Mover la lógica de cálculo compleja y la orquestación a los servicios.
*   **Implementación del Catálogo de Perfiles:** Diseñar cómo se almacenará y accederá la información detallada de los perfiles (probablemente en Supabase).
*   **Desarrollo Completo de los Handlers de API:** Definir todos los endpoints REST necesarios.
*   **Autenticación y Autorización:** Integrar con Supabase Auth.
*   **Pruebas:** Escribir pruebas unitarias y de integración exhaustivas.
*   **Documentación de API:** Generar documentación OpenAPI/Swagger.

## 8. Configuración e Instalación

(Detalles a añadir cuando el proyecto sea desplegable)

1.  Clonar el repositorio.
2.  Configurar variables de entorno (ej. credenciales de Supabase).
3.  Ejecutar `go mod tidy` para descargar dependencias.
4.  Construir el binario: `go build -o windraw_api ./cmd/window-api/`
5.  Ejecutar la aplicación: `./windraw_api`

---

Este `README.md` es un punto de partida. Deberías actualizarlo continuamente a medida que tu proyecto evoluciona. Cuanto más detallado y preciso sea, mejor podrá entenderlo un modelo de IA (y otros desarrolladores).

Considera añadir diagramas de arquitectura o de flujo de datos si eso ayuda a clarificar conceptos complejos.

---

## Directorio `internal`

El directorio `internal` es una convención en los proyectos de Go para alojar código que no se pretende que sea importado por otros proyectos externos. Es el lugar ideal para la lógica principal de tu aplicación o librería privada.

### Directorio `internal/app`

Dentro de `internal`, el subdirectorio `app` suele contener el código específico de la aplicación. Si tu proyecto define múltiples aplicaciones (por ejemplo, una API y un worker), cada una podría tener su propio subdirectorio aquí.

#### Directorio `internal/app/window-api`

Este directorio contiene todo el código específico de la aplicación `window-api`. Esto incluye los modelos de datos, la lógica de negocio (servicios), los manejadores de API (handlers), y los repositorios para la interacción con la base de datos.

##### Directorio `internal/app/window-api/models`

Aquí se definen las estructuras de datos (structs) que representan las entidades del dominio de la aplicación `window-api`. Estos modelos son la base sobre la cual operan los servicios y se almacenan en la base de datos. A continuación, se detallarán los archivos contenidos en este directorio.

---
