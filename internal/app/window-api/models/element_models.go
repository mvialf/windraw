package models

import (
	"errors"
	"fmt"
	// "github.com/tu-usuario/tu-proyecto-ventanas/internal/pkg/constants" // Descomentar cuando constants.go exista
)

type FrameDetail struct {
	Position       string  `json:"position"`
	ProfileSKU     string  `json:"profile_sku"`
	Color          string  `json:"color"`
	Dimension      int     `json:"dimension"`
	AngleLeft      float64 `json:"angle_left"`
	AngleRight     float64 `json:"angle_right"`
	ReinforcedUsed bool    `json:"reinforced_used"`
	ReinforcedSKU  string  `json:"reinforced_sku,omitempty"`
}

type WindDetail struct {
	Position       string  `json:"position"`
	ProfileSKU     string  `json:"profile_sku"`
	Color          string  `json:"color"`
	Dimension      int     `json:"dimension"`
	AngleLeft      float64 `json:"angle_left"`
	AngleRight     float64 `json:"angle_right"`
	ReinforcedUsed bool    `json:"reinforced_used"`
	ReinforcedSKU  string  `json:"reinforced_sku,omitempty"`
}

type Frame struct {
	Name      string                 `json:"name"`
	Inverted  bool                   `json:"inverted"`
	Geometry  string                 `json:"geometry"`
	Width     int                    `json:"width"`
	Height    int                    `json:"height"`
	Area      float64                `json:"area"`
	Perimeter float64                `json:"perimeter"`
	CutType   string                 `json:"cut_type"`
	Details   map[string]FrameDetail `json:"details"`
}

type Wind struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Kind             string                `json:"kind"`
	Status           string                `json:"status"`
	OpeningSide      string                `json:"opening_side,omitempty"`
	OpeningDirection string                `json:"opening_direction,omitempty"`
	Width            int                   `json:"width"`
	Height           int                   `json:"height"`
	Area             float64               `json:"area"`
	Perimeter        float64               `json:"perimeter"`
	CutType          string                `json:"cut_type"`
	Details          map[string]WindDetail `json:"details"`
}

type Element struct {
	ID         string                 `json:"id"`
	Width      int                    `json:"width"`
	Height     int                    `json:"height"`
	Material   string                 `json:"material"`
	Type       string                 `json:"type"`
	Structure  string                 `json:"structure"`
	Area       float64                `json:"area"`
	Perimeter  float64                `json:"perimeter"`
	Frame      Frame                  `json:"frame"`
	Winds      []Wind                 `json:"winds,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

func NewFrame(width, height int, geometry, cutType string, positions []string, defaultGeometry, defaultCutType string) (*Frame, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) del marco deben ser mayores a 0")
	}

	// Validación de geometry (usará constants.GEOMETRY_RECTANGULAR etc. y IsValidOption)
	// if !IsValidOption(geometry, []string{constants.GEOMETRY_RECTANGULAR, constants.GEOMETRY_TRIANGULAR}) {
	// 	 fmt.Printf("Advertencia: Geometría de marco '%s' no reconocida, usando '%s'\n", geometry, defaultGeometry)
	// 	 geometry = defaultGeometry
	// }
	// Validación de cutType (usará constants.CUT_SQUARE etc. y IsValidOption)
	// if !IsValidOption(cutType, []string{constants.CUT_SQUARE, constants.CUT_ANGLE}) {
	//  fmt.Printf("Advertencia: Tipo de corte de marco '%s' no reconocido, usando '%s'\n", cutType, defaultCutType)
	// 	cutType = defaultCutType
	// }
    // Por ahora, asumimos que los valores de entrada son válidos o se usa un valor por defecto pasado como argumento
    if geometry == "" { geometry = defaultGeometry}
    if cutType == "" { cutType = defaultCutType}


	area := float64(width) * float64(height) / 1000000.0
	perimeter := 2.0 * (float64(width) + float64(height)) / 1000.0

	details := make(map[string]FrameDetail)
	for _, pos := range positions {
		details[pos] = FrameDetail{Position: pos}
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

func NewWind(name, kind string, width, height int, cutType string, status string, positions []string, validKinds []string, validCutTypes []string) (*Wind, error) {
	if name == "" {
		return nil, errors.New("el nombre de la hoja (Wind.Name) no puede estar vacío")
	}
	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) de la hoja deben ser mayores a 0")
	}

	if !IsValidOption(kind, validKinds) { // Usará constants.WIND_KIND_*
		return nil, fmt.Errorf("tipo de hoja (kind) inválido: '%s'. Válidos: %v", kind, validKinds)
	}
	if !IsValidOption(cutType, validCutTypes) { // Usará constants.CUT_*_WIND
		return nil, fmt.Errorf("tipo de corte de hoja (cutType) inválido: '%s'. Válidos: %v", cutType, validCutTypes)
	}

	area := float64(width) * float64(height) / 1000000.0
	perimeter := 2.0 * (float64(width) + float64(height)) / 1000.0

	details := make(map[string]WindDetail)
	for _, pos := range positions {
		details[pos] = WindDetail{Position: pos}
	}

	wind := &Wind{
		ID:        generateID(),
		Name:      name,
		Kind:      kind,
		Status:    status, // Usará constants.WIND_STATUS_*
		Width:     width,
		Height:    height,
		Area:      area,
		Perimeter: perimeter,
		CutType:   cutType,
		Details:   details,
	}
	return wind, nil
}

func NewElement(width int, height int, material string, elementType string, structure string,
	defaultFrameGeometry string, defaultFrameCutType string, defaultFramePositions []string,
	validMaterials []string, validTypes []string, validStructures []string) (*Element, error) {

	if width <= 0 || height <= 0 {
		return nil, errors.New("las dimensiones (width, height) del elemento deben ser mayores a 0")
	}

	if !IsValidOption(material, validMaterials) { // Usará constants.MATERIAL_*
		return nil, fmt.Errorf("material principal inválido: '%s'. Válidos: %v", material, validMaterials)
	}
	if !IsValidOption(elementType, validTypes) { // Usará constants.TYPE_*
		return nil, fmt.Errorf("tipo de elemento inválido: '%s'. Válidos: %v", elementType, validTypes)
	}
	if !IsValidOption(structure, validStructures) { // Usará constants.STRUCTURE_*
		return nil, fmt.Errorf("estructura inválida: '%s'. Válidos: %v", structure, validStructures)
	}

	frame, err := NewFrame(width, height, defaultFrameGeometry, defaultFrameCutType, defaultFramePositions, defaultFrameGeometry, defaultFrameCutType)
	if err != nil {
		return nil, fmt.Errorf("error al crear el marco interno para el elemento: %w", err)
	}

	area := float64(width) * float64(height) / 1000000.0
	perimeter := 2.0 * (float64(width) + float64(height)) / 1000.0

	element := &Element{
		ID:         generateID(),
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

func (e *Element) AddWind(wind Wind) error {
	for _, existingWind := range e.Winds {
		if existingWind.Name == wind.Name {
			return fmt.Errorf("ya existe una hoja con el nombre '%s' en el elemento ID %s", wind.Name, e.ID)
		}
	}
	e.Winds = append(e.Winds, wind)
	return nil
}

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
	// fmt.Printf("  [DEBUG] Asignado Perfil Marco: Pos '%s', SKU '%s', Color '%s'\n", position, profileSKU, color) // Comentado para limpieza
	return nil
}

func (f *Frame) CalculateFrameDetails(anchoPerfilEjemplo float64, positionLeft, positionRight, positionTop, positionBottom, cutTypeAngle, cutTypeSquare string) error { // Parámetros de ejemplo para constantes
	// fmt.Printf("  [DEBUG] Iniciando cálculo de detalles para Marco %dx%d, CutType: %s\n", f.Width, f.Height, f.CutType) // Comentado
	calculatedDetails := make(map[string]FrameDetail)
	for pos, detail := range f.Details {
		if detail.ProfileSKU == "" {
			// fmt.Printf("    Advertencia: Perfil no asignado para la posición '%s' del marco. Saltando cálculo.\n", pos) // Comentado
			calculatedDetails[pos] = detail
			continue
		}
		// fmt.Printf("    Calculando para Posición: %s, Perfil SKU: %s\n", pos, detail.ProfileSKU) // Comentado

		calculatedDim := 0
		calculatedAngleL := 0.0
		calculatedAngleR := 0.0
		needsReinforcement := false // Lógica de refuerzo simplificada/pendiente
		reinforcementSKU := ""    // Lógica de refuerzo simplificada/pendiente

		if f.CutType == cutTypeAngle { // Asume que cutTypeAngle es "Ángulo"
			descuento := int(anchoPerfilEjemplo)
			if pos == positionLeft || pos == positionRight {
				calculatedDim = f.Height - descuento
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = f.Width - descuento
			}
			calculatedAngleL = 45.0
			calculatedAngleR = 45.0
		} else if f.CutType == cutTypeSquare { // Asume que cutTypeSquare es "Cuadrado"
			if pos == positionLeft || pos == positionRight {
				calculatedDim = f.Height
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = f.Width
			}
			calculatedAngleL = 90.0
			calculatedAngleR = 90.0
		}
		// Añadir más lógica de tipos de corte si es necesario

		detail.Dimension = calculatedDim
		detail.AngleLeft = calculatedAngleL
		detail.AngleRight = calculatedAngleR
		detail.ReinforcedUsed = needsReinforcement
		detail.ReinforcedSKU = reinforcementSKU
		calculatedDetails[pos] = detail
		// fmt.Printf("      -> Detalle Marco Calculado %s: SKU %s, Dim %d, Angles %.1f/%.1f, Refuerzo: %t (%s)\n", // Comentado
		//	pos, detail.ProfileSKU, detail.Dimension, detail.AngleLeft, detail.AngleRight, detail.ReinforcedUsed, detail.ReinforcedSKU)
	}
	f.Details = calculatedDetails
	// fmt.Println("  [DEBUG] Cálculo de detalles del Marco finalizado.") // Comentado
	return nil
}

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
	// fmt.Printf("    [DEBUG] Asignado Perfil Hoja '%s': Pos '%s', SKU '%s', Color '%s'\n", w.Name, position, profileSKU, color) // Comentado
	return nil
}

func (w *Wind) CalculateWindDetails(anchoPerfilHojaEjemplo float64, positionLeft, positionRight, positionTop, positionBottom, cutAngleWind, cutSquareWind, cutVerticalOverlapWind string) error { // Parámetros de ejemplo para constantes
	// fmt.Printf("    [DEBUG] Iniciando cálculo de detalles para Hoja '%s' %dx%d, CutType: %s\n", w.Name, w.Width, w.Height, w.CutType) // Comentado
	calculatedDetails := make(map[string]WindDetail)
	for pos, detail := range w.Details {
		if detail.ProfileSKU == "" {
			// fmt.Printf("      Advertencia: Perfil no asignado para la posición '%s' de la hoja '%s'. Saltando cálculo.\n", pos, w.Name) // Comentado
			calculatedDetails[pos] = detail
			continue
		}
		// fmt.Printf("      Calculando para Hoja '%s', Posición: %s, Perfil SKU: %s\n", w.Name, pos, detail.ProfileSKU) // Comentado

		calculatedDim := 0
		calculatedAngleL := 0.0
		calculatedAngleR := 0.0
		needsReinforcement := false // Lógica de refuerzo simplificada/pendiente
		reinforcementSKU := ""    // Lógica de refuerzo simplificada/pendiente

		if w.CutType == cutAngleWind { // Asume constante para "Ángulo"
			descuento := int(anchoPerfilHojaEjemplo)
			if pos == positionLeft || pos == positionRight {
				calculatedDim = w.Height - descuento
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = w.Width - descuento
			}
			calculatedAngleL = 45.0
			calculatedAngleR = 45.0
		} else if w.CutType == cutSquareWind { // Asume constante para "Cuadrado"
			if pos == positionLeft || pos == positionRight {
				calculatedDim = w.Height
			} else if pos == positionTop || pos == positionBottom {
				calculatedDim = w.Width
			}
			calculatedAngleL = 90.0
			calculatedAngleR = 90.0
		} else if w.CutType == cutVerticalOverlapWind { // Asume constante para "Vertical superpuesto"
			if pos == positionLeft || pos == positionRight {
				calculatedDim = w.Height
				calculatedAngleL = 0.0 // o 90.0, depende de la convención
				calculatedAngleR = 0.0 // o 90.0
			} else if pos == positionTop || pos == positionBottom { // Comportamiento para horizontales en este caso?
				// Asumimos corte a 45 para horizontales si verticales son superpuestos, como en el original
				descuento := int(anchoPerfilHojaEjemplo)
				calculatedDim = w.Width - descuento
				calculatedAngleL = 45.0
				calculatedAngleR = 45.0
			}
		}
		// Añadir más lógica de tipos de corte si es necesario

		detail.Dimension = calculatedDim
		detail.AngleLeft = calculatedAngleL
		detail.AngleRight = calculatedAngleR
		detail.ReinforcedUsed = needsReinforcement
		detail.ReinforcedSKU = reinforcementSKU
		calculatedDetails[pos] = detail
		// fmt.Printf("        -> Detalle Hoja Calculado %s: SKU %s, Dim %d, Angles %.1f/%.1f, Refuerzo: %t (%s)\n", // Comentado
		//	pos, detail.ProfileSKU, detail.Dimension, detail.AngleLeft, detail.AngleRight, detail.ReinforcedUsed, detail.ReinforcedSKU)
	}
	w.Details = calculatedDetails
	// fmt.Printf("    [DEBUG] Cálculo de detalles de la Hoja '%s' finalizado.\n", w.Name) // Comentado
	return nil
}