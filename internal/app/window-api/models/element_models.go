package models

import (
	"errors"
	"fmt"
	// "github.com/mvialf/windraw/internal/pkg/constants" // Descomentar cuando constants.go exista y sea necesario y se use aquí
)

// FrameDetail describe una pieza individual de perfil para un marco.
type FrameDetail struct {
	Position       string  `json:"position"`                 // Posición del perfil (ej. "Top", "Bottom", "Left", "Right")
	ProfileSKU     string  `json:"profile_sku"`              // SKU del perfil utilizado
	Color          string  `json:"color"`                    // Color del perfil
	Dimension      float64 `json:"dimension"`                // MODIFICADO: Longitud de corte del perfil en mm (ahora float64)
	AngleLeft      float64 `json:"angle_left"`               // Ángulo de corte izquierdo en grados
	AngleRight     float64 `json:"angle_right"`              // Ángulo de corte derecho en grados
	ReinforcedUsed bool    `json:"reinforced_used"`          // Indica si se utiliza refuerzo
	ReinforcedSKU  string  `json:"reinforced_sku,omitempty"` // SKU del refuerzo, si se utiliza
	// ID             string  `json:"id"` // Considera si cada detalle necesita un ID único persistente
	// Notes          string  `json:"notes,omitempty"`
}

// WindDetail describe una pieza individual de perfil para una hoja.
type WindDetail struct {
	Position       string  `json:"position"`                 // Posición del perfil
	ProfileSKU     string  `json:"profile_sku"`              // SKU del perfil
	Color          string  `json:"color"`                    // Color del perfil
	Dimension      float64 `json:"dimension"`                // MODIFICADO: Longitud de corte en mm (ahora float64)
	AngleLeft      float64 `json:"angle_left"`               // Ángulo de corte izquierdo
	AngleRight     float64 `json:"angle_right"`              // Ángulo de corte derecho
	ReinforcedUsed bool    `json:"reinforced_used"`          // Si usa refuerzo
	ReinforcedSKU  string  `json:"reinforced_sku,omitempty"` // SKU del refuerzo
	// Drillings      []string `json:"drillings,omitempty"`
	// ID             string  `json:"id"` // Considera si cada detalle necesita un ID único persistente
	// Notes          string  `json:"notes,omitempty"`
}

// Frame representa el marco perimetral de un Element.
type Frame struct {
	Name      string                 `json:"name"`      // Nombre del marco (ej. "Marco Principal")
	Inverted  bool                   `json:"inverted"`  // Si el marco está invertido (puede afectar cálculos)
	Geometry  string                 `json:"geometry"`  // Geometría del marco (ej. "Rectangular")
	Width     int                    `json:"width"`     // Ancho exterior del marco en mm (generalmente igual al del Element)
	Height    int                    `json:"height"`    // Alto exterior del marco en mm (generalmente igual al del Element)
	Area      float64                `json:"area"`      // Área calculada del marco en m²
	Perimeter float64                `json:"perimeter"` // Perímetro calculado del marco en m
	CutType   string                 `json:"cut_type"`  // Tipo de corte para los perfiles del marco
	Details   map[string]FrameDetail `json:"details"`   // Mapa de detalles de perfiles por posición (constants.POSITION_*)
}

// Wind representa una hoja (panel móvil o fijo) dentro de un Element.
type Wind struct {
	ID               string                `json:"id"`                          // ID único de la hoja, generado por generateID()
	Name             string                `json:"name"`                        // Nombre de la hoja (ej. "Hoja Izquierda Móvil")
	Kind             string                `json:"kind"`                        // Tipo de hoja (ej. "Corredera", "Abatible", "Fijo" - Ver constants.KIND_*)
	Status           string                `json:"status"`                      // Estado de la hoja (ej. "Activa", "Inactiva")
	OpeningSide      string                `json:"opening_side,omitempty"`      // Lado de apertura (ej. "Izquierda", "Derecha")
	OpeningDirection string                `json:"opening_direction,omitempty"` // Dirección de apertura (ej. "Interior", "Exterior")
	Width            int                   `json:"width"`                       // Ancho de la hoja en mm (esto será calculado por WindowCalculationService)
	Height           int                   `json:"height"`                      // Alto de la hoja en mm (esto será calculado por WindowCalculationService)
	Area             float64               `json:"area"`                        // Área calculada de la hoja en m²
	Perimeter        float64               `json:"perimeter"`                   // Perímetro calculado de la hoja en m
	CutType          string                `json:"cut_type"`                    // Tipo de corte para los perfiles de la hoja
	Details          map[string]WindDetail `json:"details"`                     // Mapa de detalles de perfiles por posición (constants.POSITION_*)
	// Errors           []string           `json:"errors,omitempty"`      // Para registrar errores específicos de esta hoja
}

// Element representa la unidad funcional principal, una ventana o puerta individual.
type Element struct {
	ID         string                 `json:"id"`                   // ID único del elemento, generado por generateID()
	Width      int                    `json:"width"`                // Ancho total del elemento en mm
	Height     int                    `json:"height"`               // Alto total del elemento en mm
	Material   string                 `json:"material"`             // Material principal del elemento (ej. "PVC", "Aluminio")
	Type       string                 `json:"type"`                 // Tipología del elemento (ej. "Corredera", "Abatible" - Ver constants.TYPE_*)
	Structure  string                 `json:"structure"`            // Estructura (ej. "Ventana", "Puerta")
	Area       float64                `json:"area"`                 // Área calculada del elemento en m²
	Perimeter  float64                `json:"perimeter"`            // Perímetro calculado del elemento en m
	Frame      Frame                  `json:"frame"`                // El marco del elemento
	Winds      []Wind                 `json:"winds,omitempty"`      // Lista de hojas dentro del elemento
	Properties map[string]interface{} `json:"properties,omitempty"` // Propiedades adicionales (ej. color vidrio, tipo manilla)
	// Errors     []string               `json:"errors,omitempty"` // Para registrar errores de cálculo o selección a nivel de elemento
}

// NewFrame es el constructor para la estructura Frame.
func NewFrame(width, height int, geometry, cutType string, positions []string, defaultGeometry, defaultCutType string) (*Frame, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) del marco deben ser mayores a 0")
	}

	// Validaciones de geometry y cutType (actualmente comentadas, requieren constants)
	// if !IsValidOption(geometry, []string{constants.GEOMETRY_RECTANGULAR, ...}) { geometry = defaultGeometry }
	// if !IsValidOption(cutType, []string{constants.CUT_SQUARE, ...}) { cutType = defaultCutType }
	if geometry == "" { // Si no se valida y viene vacío, usar default
		geometry = defaultGeometry
	}
	if cutType == "" { // Si no se valida y viene vacío, usar default
		cutType = defaultCutType
	}

	area := float64(width*height) / 1000000.0         // conversión a m²
	perimeter := 2.0 * float64(width+height) / 1000.0 // conversión a m

	details := make(map[string]FrameDetail)
	if positions != nil {
		for _, pos := range positions {
			// Inicializa con la posición, otros campos se llenarán después
			details[pos] = FrameDetail{Position: pos, AngleLeft: 90.0, AngleRight: 90.0} // Default angles
		}
	}

	frame := &Frame{
		Name:      "Marco Principal",
		Geometry:  geometry,
		Width:     width,
		Height:    height,
		Area:      area,
		Perimeter: perimeter,
		CutType:   cutType,
		Details:   details,
		Inverted:  false,
	}
	return frame, nil
}

// NewWind es el constructor para la estructura Wind.
func NewWind(name, kind string /*width, height int,*/, cutType string, status string, positions []string,
	validKinds []string, validCutTypes []string) (*Wind, error) { // Width y Height se calcularán

	if name == "" {
		return nil, errors.New("el nombre de la hoja (Wind.Name) no puede estar vacío")
	}
	// Width y Height ya no son input directo, se calcularán.
	// if width <= 0 || height <= 0 {
	// 	return nil, errors.New("las dimensiones (width, height) de la hoja deben ser mayores a 0")
	// }

	// Validaciones (requieren que IsValidOption y las constantes estén definidas)
	// if !IsValidOption(kind, validKinds) {
	// 	return nil, fmt.Errorf("tipo de hoja (kind) inválido: '%s'. Válidos: %v", kind, validKinds)
	// }
	// if !IsValidOption(cutType, validCutTypes) {
	// 	return nil, fmt.Errorf("tipo de corte de hoja (cutType) inválido: '%s'. Válidos: %v", cutType, validCutTypes)
	// }

	// Area y Perímetro se calcularán una vez que Width y Height de la hoja sean conocidos.
	// area := float64(width*height) / 1000000.0
	// perimeter := 2.0 * float64(width+height) / 1000.0

	details := make(map[string]WindDetail)
	if positions != nil {
		for _, pos := range positions {
			// Inicializa con la posición, otros campos se llenarán después
			details[pos] = WindDetail{Position: pos, AngleLeft: 90.0, AngleRight: 90.0} // Default angles
		}
	}

	// Asumiendo que generateID() está en el mismo paquete 'models' (definido en project_models.go)
	windID := generateID() // Asegúrate que esta función esté accesible

	wind := &Wind{
		ID:     windID,
		Name:   name,
		Kind:   kind,
		Status: status,
		// Width y Height se calcularán por WindowCalculationService
		// Area y Perimeter también
		CutType: cutType,
		Details: details,
	}
	return wind, nil
}

// NewElement es el constructor para la estructura Element.
func NewElement(width int, height int, material string, elementType string, structure string,
	defaultFrameGeometry string, defaultFrameCutType string, defaultFramePositions []string,
	validMaterials []string, validTypes []string, validStructures []string) (*Element, error) {

	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) del elemento deben ser mayores a 0")
	}

	// Validaciones (requieren IsValidOption y constantes)
	// if !IsValidOption(material, validMaterials) { ... }
	// if !IsValidOption(elementType, validTypes) { ... }
	// if !IsValidOption(structure, validStructures) { ... }

	frame, err := NewFrame(width, height, defaultFrameGeometry, defaultFrameCutType, defaultFramePositions, defaultFrameGeometry, defaultFrameCutType)
	if err != nil {
		return nil, fmt.Errorf("error al crear el marco interno para el elemento: %w", err)
	}

	area := float64(width*height) / 1000000.0
	perimeter := 2.0 * float64(width+height) / 1000.0

	elementID := generateID() // Asegúrate que esta función esté accesible

	element := &Element{
		ID:         elementID,
		Width:      width,
		Height:     height,
		Material:   material,
		Type:       elementType,
		Structure:  structure,
		Area:       area,
		Perimeter:  perimeter,
		Frame:      *frame,
		Winds:      []Wind{},
		Properties: make(map[string]interface{}),
	}
	return element, nil
}

// AddWind añade una hoja (Wind) al Element.
func (e *Element) AddWind(wind Wind) error {
	for _, existingWind := range e.Winds {
		if existingWind.Name == wind.Name { // Podría ser mejor verificar por ID si son únicos globales
			return fmt.Errorf("ya existe una hoja con el nombre '%s' en el elemento ID %s", wind.Name, e.ID)
		}
	}
	e.Winds = append(e.Winds, wind)
	return nil
}

// SetFrameProfile establece el SKU del perfil y el color para una posición específica del marco.
// ESTA FUNCIÓN PODRÍA SER REDUNDANTE si ProfileSelectorService y WindowCalculationService
// se encargan de rellenar los ProfileSKU directamente en los Details.
// Si se mantiene, el color podría venir del StockItem asociado al perfil.
func (f *Frame) SetFrameProfile(position string, profileSKU string, color string) error {
	detail, ok := f.Details[position]
	if !ok {
		// Si la posición no existe, podría ser deseable crearla aquí.
		// f.Details[position] = FrameDetail{Position: position}
		// detail = f.Details[position]
		return fmt.Errorf("posición de marco inválida o no inicializada: '%s'", position)
	}
	if profileSKU == "" {
		return fmt.Errorf("el SKU del perfil (profileSKU) no puede estar vacío para la posición '%s'", position)
	}
	detail.ProfileSKU = profileSKU
	detail.Color = color // El color vendrá del models.StockItem asociado al perfil y color seleccionado
	f.Details[position] = detail
	return nil
}

// SetWindProfile establece el SKU del perfil y el color para una posición específica de la hoja.
// MISMA NOTA QUE SetFrameProfile: podría ser redundante.
func (w *Wind) SetWindProfile(position string, profileSKU string, color string) error {
	detail, ok := w.Details[position]
	if !ok {
		return fmt.Errorf("posición de hoja inválida o no inicializada: '%s' en hoja '%s'", position, w.Name)
	}
	if profileSKU == "" {
		return fmt.Errorf("el SKU del perfil (profileSKU) no puede estar vacío para la posición '%s' de la hoja '%s'", position, w.Name)
	}
	detail.ProfileSKU = profileSKU
	detail.Color = color // El color vendrá del models.StockItem
	w.Details[position] = detail
	return nil
}

// Las funciones CalculateFrameDetails y CalculateWindDetails HAN SIDO ELIMINADAS
// ya que esta lógica ahora reside y se está desarrollando en WindowCalculationService.
