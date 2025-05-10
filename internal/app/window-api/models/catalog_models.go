package models

import (
	"time" // <--- AÑADE ESTA LÍNEA
)

// Profile representa un tipo de perfil del catálogo.
type Profile struct {
	ID          string   `json:"id"` // O int, dependiendo de tu PK en Supabase
	SKU         string   `json:"sku"`
	Description string   `json:"description"`
	Material    string   `json:"material"` // e.g., constants.MATERIAL_PVC, constants.MATERIAL_ALUMINIO
	WeightPerM  float64  `json:"weight_per_meter"`
	Colors      []string `json:"available_colors"` // Si los colores son un array de texto en Supabase
	// ... otros campos relevantes del perfil: dimensiones, tipo, etc.

	// Campos de auditoría opcionales que Supabase podría gestionar
	CreatedAt time.Time `json:"created_at"` // Ahora 'time.Time' será reconocido
	UpdatedAt time.Time `json:"updated_at"` // Ahora 'time.Time' será reconocido
}

// Aquí podrías tener otros modelos de catálogo si los necesitas:
// type GlassType struct { ... }
// type HardwareItem struct { ... }
