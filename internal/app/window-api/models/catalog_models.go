package models

import (
	"time"
)

// Color representa un color del catálogo.
type Color struct {
	ColorID int64  `json:"color_id"` // PK
	Name    string `json:"name"`
	HexCode string `json:"hex_code"`
}

// Material representa un material del catálogo.
type Material struct {
	MaterialID int64  `json:"material_id"` // PK
	Name       string `json:"name"`
}

// Supplier representa un proveedor.
type Supplier struct {
	SupplierID int64  `json:"supplier_id"` // PK
	Name       string `json:"name"`
}

// Profile representa un tipo de perfil del catálogo.
// Dimensiones H/W: Para perfiles horizontales, H es alto visible frontal, W es profundidad.
// Para perfiles verticales, H es ancho visible frontal, W es profundidad.
// H2 es una sub-dimensión de H para encajes/reducción de vano.
type Profile struct {
	ProfileID          int64     `json:"profile_id"`
	ProfileSKU         string    `json:"profile_sku"`
	ProfileName        string    `json:"profile_name"`
	ProfileType        string    `json:"profile_type"`
	MaterialID         int64     `json:"material_id"`
	SupplierID         int64     `json:"supplier_id"`
	ProfileWeightMeter *float64  `json:"profile_weigth_meter"`
	ProfileH           *float64  `json:"profile_h"` // Para horiz: alto visible; Para vert: ancho visible
	ProfileW           *float64  `json:"profile_w"` // Profundidad del perfil
	ProfileH1          *float64  `json:"profile_h1"`
	ProfileH2          *float64  `json:"profile_h2"` // Dimensión de encaje/reducción de vano
	ProfileH3          *float64  `json:"profile_h3"`
	ProfileH4          *float64  `json:"profile_h4"`
	ProfileW1          *float64  `json:"profile_w1"`
	ProfileCutType     *string   `json:"profile_cut_type"`
	CutMarginMM        *float64  `json:"cut_margin_mm"`
	TrackCount         *int      `json:"track_count"`
	UsesOverlap        *bool     `json:"uses_overlap"` // Si el perfil de encuentro de hoja necesita traslapo adicional
	ProfilePosition    string    `json:"profile_position"`
	ProfileStructure   *string   `json:"profile_structure"` // Ej. "Marco", "Hoja"
	WeldMargin         bool      `json:"weld_margin"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// ProfileSystem representa un sistema de perfiles.
type ProfileSystem struct {
	SystemID        int64     `json:"system_id"`
	Name            string    `json:"name"`
	SupplierID      int64     `json:"supplier_id"`
	Type            string    `json:"type"` // Ej. constants.TYPE_SLIDING
	MaterialID      int64     `json:"material_id"`
	UsesGlassBead   bool      `json:"uses_glass_bead"`
	GlassMarginMM   float64   `json:"glass_margin_mm"`
	TopOverlapMM    float64   `json:"top_overlap_mm"`    // Traslape vertical hoja sobre marco
	BottomOverlapMM float64   `json:"bottom_overlap_mm"` // Traslape vertical hoja sobre marco
	SideOverlapMM   float64   `json:"side_overlap_mm"`   // Cuánto se "mete" cada hoja en el marco lateral
	Primacy         *int64    `json:"prymacy"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
}

// StockItem representa un ítem de perfil en un color específico.
// No incluye cantidad de stock por ahora.
type StockItem struct {
	StockItemID   int64    `json:"stock_item_id"`
	ProfileID     int64    `json:"profile_id"`
	ColorID       int64    `json:"color_id"`
	SupplierID    int64    `json:"supplier_id"`
	ItemSKU       string   `json:"item_sku"`
	ProfilePrice  float64  `json:"profile_price"`
	ProfileLength *float64 `json:"profile_length"`
}

// ProfileReinforcement representa la relación de refuerzo para un perfil.
type ProfileReinforcement struct {
	MainProfileID          int64   `json:"main_profile_id"`
	ReinforcementProfileID int64   `json:"reinforcement_profile_id"`
	ReinforcementGapMM     float64 `json:"reinforcement_gap_mm"`
}

// SystemAvailableColor representa un color disponible para un sistema de perfiles.
type SystemAvailableColor struct {
	SystemID        int64  `json:"system_id"`
	ColorID         int64  `json:"color_id"`
	ColorCodeSuffix string `json:"color_code_suffix"`
}

// SystemProfileListItem representa un perfil dentro de un sistema de perfiles.
type SystemProfileListItem struct {
	SystemID        int64  `json:"system_id"`
	ProfileID       int64  `json:"profile_id"`
	Primacy         *int16 `json:"primacy"`
	ElementPartRole string `json:"element_part_role,omitempty"` // Rol funcional (constants.ROLE_*)
}
