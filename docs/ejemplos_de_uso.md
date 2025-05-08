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
	"github.com/tu-usuario/tu-proyecto-ventanas/internal/app/window-api/models"
	"github.com/tu-usuario/tu-proyecto-ventanas/internal/pkg/constants"
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


	// 10. Añadir la Hoja al Elemento Ventana
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