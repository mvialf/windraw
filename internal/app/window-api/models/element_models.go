package models

import (
	"errors"
	"fmt"
	// "github.com/mvialf/windraw/internal/pkg/constants" // Descomentar cuando constants.go exista y sea necesario
)

// FrameDetail describe una pieza individual de perfil para un marco.
type FrameDetail struct {
	Position       string  `json:"position"`                 // Posición del perfil (ej. "Top", "Bottom", "Left", "Right")
	ProfileSKU     string  `json:"profile_sku"`              // SKU del perfil utilizado
	Color          string  `json:"color"`                    // Color del perfil
	Dimension      int     `json:"dimension"`                // Longitud de corte del perfil en mm
	AngleLeft      float64 `json:"angle_left"`               // Ángulo de corte izquierdo en grados
	AngleRight     float64 `json:"angle_right"`              // Ángulo de corte derecho en grados
	ReinforcedUsed bool    `json:"reinforced_used"`          // Indica si se utiliza refuerzo
	ReinforcedSKU  string  `json:"reinforced_sku,omitempty"` // SKU del refuerzo, si se utiliza
}

// WindDetail describe una pieza individual de perfil para una hoja.
type WindDetail struct {
	Position       string  `json:"position"`                 // Posición del perfil
	ProfileSKU     string  `json:"profile_sku"`              // SKU del perfil
	Color          string  `json:"color"`                    // Color del perfil
	Dimension      int     `json:"dimension"`                // Longitud de corte en mm
	AngleLeft      float64 `json:"angle_left"`               // Ángulo de corte izquierdo
	AngleRight     float64 `json:"angle_right"`              // Ángulo de corte derecho
	ReinforcedUsed bool    `json:"reinforced_used"`          // Si usa refuerzo
	ReinforcedSKU  string  `json:"reinforced_sku,omitempty"` // SKU del refuerzo
}

// Frame representa el marco perimetral de un Element.
type Frame struct {
	Name      string                 `json:"name"`      // Nombre del marco (ej. "Marco Principal")
	Inverted  bool                   `json:"inverted"`  // Si el marco está invertido (puede afectar cálculos de descuento)
	Geometry  string                 `json:"geometry"`  // Geometría del marco (ej. "Rectangular", "Trapezoidal")
	Width     int                    `json:"width"`     // Ancho exterior del marco en mm
	Height    int                    `json:"height"`    // Alto exterior del marco en mm
	Area      float64                `json:"area"`      // Área calculada del marco en m²
	Perimeter float64                `json:"perimeter"` // Perímetro calculado del marco en m
	CutType   string                 `json:"cut_type"`  // Tipo de corte para los perfiles del marco (ej. "Ángulo", "Cuadrado")
	Details   map[string]FrameDetail `json:"details"`   // Mapa de detalles de perfiles por posición
}

// Wind representa una hoja (panel móvil o fijo) dentro de un Element.
type Wind struct {
	ID               string                `json:"id"`                          // ID único de la hoja, generado por generateID()
	Name             string                `json:"name"`                        // Nombre de la hoja (ej. "Hoja Izquierda Móvil")
	Kind             string                `json:"kind"`                        // Tipo de hoja (ej. "Corredera", "Abatible", "Fija")
	Status           string                `json:"status"`                      // Estado de la hoja (ej. "Activa", "Inactiva")
	OpeningSide      string                `json:"opening_side,omitempty"`      // Lado de apertura (ej. "Izquierda", "Derecha")
	OpeningDirection string                `json:"opening_direction,omitempty"` // Dirección de apertura (ej. "Interior", "Exterior")
	Width            int                   `json:"width"`                       // Ancho de la hoja en mm
	Height           int                   `json:"height"`                      // Alto de la hoja en mm
	Area             float64               `json:"area"`                        // Área calculada de la hoja en m²
	Perimeter        float64               `json:"perimeter"`                   // Perímetro calculado de la hoja en m
	CutType          string                `json:"cut_type"`                    // Tipo de corte para los perfiles de la hoja
	Details          map[string]WindDetail `json:"details"`                     // Mapa de detalles de perfiles por posición
}

// Element representa la unidad funcional principal, una ventana o puerta individual.
type Element struct {
	ID         string                 `json:"id"`                   // ID único del elemento, generado por generateID()
	Width      int                    `json:"width"`                // Ancho total del elemento en mm
	Height     int                    `json:"height"`               // Alto total del elemento en mm
	Material   string                 `json:"material"`             // Material principal del elemento (ej. "PVC", "Aluminio")
	Type       string                 `json:"type"`                 // Tipología del elemento (ej. "Corredera", "Abatible")
	Structure  string                 `json:"structure"`            // Estructura (ej. "Ventana", "Puerta")
	Area       float64                `json:"area"`                 // Área calculada del elemento en m²
	Perimeter  float64                `json:"perimeter"`            // Perímetro calculado del elemento en m
	Frame      Frame                  `json:"frame"`                // El marco del elemento
	Winds      []Wind                 `json:"winds,omitempty"`      // Lista de hojas dentro del elemento
	Properties map[string]interface{} `json:"properties,omitempty"` // Propiedades adicionales (ej. color vidrio, tipo manilla)
}

// NewFrame es el constructor para la estructura Frame.
// Los parámetros defaultFrameGeometry y defaultFrameCutType se usarían si las validaciones con constants estuvieran activas.
func NewFrame(width, height int, geometry, cutType string, positions []string, defaultGeometry, defaultCutType string) (*Frame, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) del marco deben ser mayores a 0")
	}

	// Las validaciones con IsValidOption y constants están comentadas.
	// Cuando el paquete constants esté disponible y poblado, se podrán descomentar.
	// if !IsValidOption(geometry, []string{constants.GEOMETRY_RECTANGULAR, constants.GEOMETRY_TRIANGULAR /*, ...otras geometrías */}) {
	//     fmt.Printf("Advertencia: Geometría de marco '%s' no reconocida, usando '%s'\n", geometry, defaultGeometry)
	//     geometry = defaultGeometry
	// }
	// if !IsValidOption(cutType, []string{constants.CUT_SQUARE, constants.CUT_ANGLE /*, ...otros tipos de corte */}) {
	//    fmt.Printf("Advertencia: Tipo de corte de marco '%s' no reconocido, usando '%s'\n", cutType, defaultCutType)
	//    cutType = defaultCutType
	// }

	if geometry == "" {
		geometry = defaultGeometry
	}
	if cutType == "" {
		cutType = defaultCutType
	}

	area := float64(width) * float64(height) / 1000000.0           // conversión a m²
	perimeter := 2.0 * (float64(width) + float64(height)) / 1000.0 // conversión a m

	details := make(map[string]FrameDetail)
	if positions != nil {
		for _, pos := range positions {
			details[pos] = FrameDetail{Position: pos}
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
func NewWind(name, kind string, width, height int, cutType string, status string, positions []string,
	validKinds []string, validCutTypes []string) (*Wind, error) {

	if name == "" {
		return nil, errors.New("el nombre de la hoja (Wind.Name) no puede estar vacío")
	}
	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) de la hoja deben ser mayores a 0")
	}

	if !IsValidOption(kind, validKinds) {
		return nil, fmt.Errorf("tipo de hoja (kind) inválido: '%s'. Válidos: %v", kind, validKinds)
	}
	if !IsValidOption(cutType, validCutTypes) {
		return nil, fmt.Errorf("tipo de corte de hoja (cutType) inválido: '%s'. Válidos: %v", cutType, validCutTypes)
	}

	area := float64(width) * float64(height) / 1000000.0
	perimeter := 2.0 * (float64(width) + float64(height)) / 1000.0

	details := make(map[string]WindDetail)
	if positions != nil {
		for _, pos := range positions {
			details[pos] = WindDetail{Position: pos}
		}
	}

	// Asumiendo que generateID() está en el mismo paquete 'models' (definido en project_models.go)
	windID := generateID()

	wind := &Wind{
		ID:        windID,
		Name:      name,
		Kind:      kind,
		Status:    status,
		Width:     width,
		Height:    height,
		Area:      area,
		Perimeter: perimeter,
		CutType:   cutType, // profile_cut_type tabla profile
		Details:   details,
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

	if !IsValidOption(material, validMaterials) {
		return nil, fmt.Errorf("material principal inválido: '%s'. Válidos: %v", material, validMaterials)
	}
	if !IsValidOption(elementType, validTypes) {
		return nil, fmt.Errorf("tipo de elemento inválido: '%s'. Válidos: %v", elementType, validTypes)
	}
	if !IsValidOption(structure, validStructures) {
		return nil, fmt.Errorf("estructura inválida: '%s'. Válidos: %v", structure, validStructures)
	}

	frame, err := NewFrame(width, height, defaultFrameGeometry, defaultFrameCutType, defaultFramePositions, defaultFrameGeometry, defaultFrameCutType)
	if err != nil {
		return nil, fmt.Errorf("error al crear el marco interno para el elemento: %w", err)
	}

	area := float64(width) * float64(height) / 1000000.0
	perimeter := 2.0 * (float64(width) + float64(height)) / 1000.0

	// Asumiendo que generateID() está en el mismo paquete 'models'
	elementID := generateID()

	element := &Element{
		ID:       elementID,
		Width:    width,
		Height:   height,
		Material: material,    //tabla material
		Type:     elementType, // sliding or casement
		//Profyle_system:	Systema de perfiles	# tabla profyle_systems
		Structure:  structure, // profyle_structure tabla profiles
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
		if existingWind.Name == wind.Name {
			return fmt.Errorf("ya existe una hoja con el nombre '%s' en el elemento ID %s", wind.Name, e.ID)
		}
	}
	e.Winds = append(e.Winds, wind)
	return nil
}

// SetFrameProfile establece el SKU del perfil y el color para una posición específica del marco.
func (f *Frame) SetFrameProfile(position string, profileSKU string, color string) error {
	detail, ok := f.Details[position]
	if !ok {
		return fmt.Errorf("posición de marco inválida o no inicializada: '%s'", position)
	}
	if profileSKU == "" {
		return fmt.Errorf("el SKU del perfil (profileSKU) no puede estar vacío para la posición '%s'", position)
	}
	detail.ProfileSKU = profileSKU
	detail.Color = color
	f.Details[position] = detail
	return nil
}

// CalculateFrameDetails calcula las dimensiones y ángulos de corte para los perfiles del marco.
func (f *Frame) CalculateFrameDetails(anchoPerfilEjemplo float64,
	positionLeft, positionRight, positionTop, positionBottom,
	cutTypeAngle, cutTypeSquare string) error { // Estos strings deberían ser constantes

	calculatedDetails := make(map[string]FrameDetail)
	for pos, detail := range f.Details {
		if detail.ProfileSKU == "" {
			calculatedDetails[pos] = detail
			continue
		}

		calculatedDim := 0
		calculatedAngleL := 0.0
		calculatedAngleR := 0.0
		needsReinforcement := false
		reinforcementSKU := ""

		if f.CutType == cutTypeAngle {
			descuento := int(anchoPerfilEjemplo)
			if pos == positionLeft || pos == positionRight {
				calculatedDim = f.Height - descuento
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = f.Width - descuento
			}
			calculatedAngleL = 45.0
			calculatedAngleR = 45.0
		} else if f.CutType == cutTypeSquare {
			if pos == positionLeft || pos == positionRight {
				calculatedDim = f.Height
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = f.Width
			}
			calculatedAngleL = 90.0
			calculatedAngleR = 90.0
		}

		detail.Dimension = calculatedDim
		detail.AngleLeft = calculatedAngleL
		detail.AngleRight = calculatedAngleR
		detail.ReinforcedUsed = needsReinforcement
		detail.ReinforcedSKU = reinforcementSKU
		calculatedDetails[pos] = detail
	}
	f.Details = calculatedDetails
	return nil
}

// SetWindProfile establece el SKU del perfil y el color para una posición específica de la hoja.
func (w *Wind) SetWindProfile(position string, profileSKU string, color string) error {
	detail, ok := w.Details[position]
	if !ok {
		return fmt.Errorf("posición de hoja inválida o no inicializada: '%s' en hoja '%s'", position, w.Name)
	}
	if profileSKU == "" {
		return fmt.Errorf("el SKU del perfil (profileSKU) no puede estar vacío para la posición '%s' de la hoja '%s'", position, w.Name)
	}
	detail.ProfileSKU = profileSKU
	detail.Color = color
	w.Details[position] = detail
	return nil
}

// CalculateWindDetails calcula las dimensiones y ángulos de corte para los perfiles de la hoja.
func (w *Wind) CalculateWindDetails(anchoPerfilHojaEjemplo float64,
	positionLeft, positionRight, positionTop, positionBottom,
	cutAngleWind, cutSquareWind, cutVerticalOverlapWind string) error { // Estos strings deberían ser constantes

	calculatedDetails := make(map[string]WindDetail)
	for pos, detail := range w.Details {
		if detail.ProfileSKU == "" {
			calculatedDetails[pos] = detail
			continue
		}

		calculatedDim := 0
		calculatedAngleL := 0.0
		calculatedAngleR := 0.0
		needsReinforcement := false
		reinforcementSKU := ""

		if w.CutType == cutAngleWind {
			descuento := int(anchoPerfilHojaEjemplo)
			if pos == positionLeft || pos == positionRight {
				calculatedDim = w.Height - descuento
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = w.Width - descuento
			}
			calculatedAngleL = 45.0
			calculatedAngleR = 45.0
		} else if w.CutType == cutSquareWind {
			if pos == positionLeft || pos == positionRight {
				calculatedDim = w.Height
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = w.Width
			}
			calculatedAngleL = 90.0
			calculatedAngleR = 90.0
		} else if w.CutType == cutVerticalOverlapWind {
			if pos == positionLeft || pos == positionRight {
				calculatedDim = w.Height
				calculatedAngleL = 90.0
				calculatedAngleR = 90.0
			} else if pos == positionTop || pos == positionBottom {
				descuento := int(anchoPerfilHojaEjemplo)
				calculatedDim = w.Width - descuento
				calculatedAngleL = 45.0
				calculatedAngleR = 45.0
			}
		}

		detail.Dimension = calculatedDim
		detail.AngleLeft = calculatedAngleL
		detail.AngleRight = calculatedAngleR
		detail.ReinforcedUsed = needsReinforcement
		detail.ReinforcedSKU = reinforcementSKU
		calculatedDetails[pos] = detail
	}
	w.Details = calculatedDetails
	return nil
}
