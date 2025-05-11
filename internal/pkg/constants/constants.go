package constants

// General Application Constants
const (
	IVA_RATE = 0.19 // Ejemplo de tasa de IVA
)

// Material Types (ejemplos, ajústalos)
const (
	MATERIAL_PVC      = "PVC"
	MATERIAL_ALUMINUM = "Aluminum"
	MATERIAL_WOOD     = "Wood"
)

// Element Types / System Types (para ProfileSystem.Type)
const (
	TYPE_SLIDING  = "Sliding"
	TYPE_CASEMENT = "Casement"
	// TYPE_FIXED    = "Fixed" // Si tienes sistemas específicos para fijos
)

// Wind Kind Types (para models.Wind.Kind)
const (
	KIND_SLIDING_WIND  = "Corredera" // Hoja corredera móvil
	KIND_CASEMENT_WIND = "Abatible"  // Hoja abatible móvil
	KIND_FIXED_PANE    = "Fijo"      // Paño fijo (puede estar modelado como un Wind)
	// ... otros kinds ...
)

// Profile Types (ejemplos para `profiles.profile_type`)
// Estos son más descriptivos del perfil en sí, no de su rol en un sistema.
const (
	PROFILE_TYPE_SLIDING_FRAME   = "Marco Corredera"
	PROFILE_TYPE_SLIDING_WIND    = "Hoja Corredera"
	PROFILE_TYPE_CASEMENT_FRAME  = "Marco Abatible"
	PROFILE_TYPE_CASEMENT_WIND   = "Hoja Abatible"
	PROFILE_TYPE_MULLION         = "Montante"
	PROFILE_TYPE_GLASSBEAD       = "Junquillo"
	PROFILE_TYPE_REINFORCEMENT   = "Refuerzo"
	PROFILE_TYPE_THRESHOLD       = "Umbral"
	PROFILE_TYPE_TRACK_RAIL      = "Riel"
	PROFILE_TYPE_ADAPTER         = "Adaptador"
	PROFILE_TYPE_OVERLAP_PROFILE = "Perfil de Traslapo"
)

// Element Part Roles (para `system_profile_list.element_part_role`)
// Define el rol funcional/posicional de un perfil en un sistema.
const (
	// --- Roles Genéricos ---
	ROLE_GLASSBEAD_WIND  = "GLASSBEAD_WIND"
	ROLE_GLASSBEAD_FRAME = "GLASSBEAD_FRAME"
	ROLE_MULLION         = "MULLION"

	// --- Roles para Sistemas Correderos (Sliding) ---
	ROLE_FRAME_PERIMETER_TOP_SLIDING    = "FRAME_PERIMETER_TOP_SLIDING"
	ROLE_FRAME_PERIMETER_BOTTOM_SLIDING = "FRAME_PERIMETER_BOTTOM_SLIDING"
	ROLE_FRAME_PERIMETER_SIDE_SLIDING   = "FRAME_PERIMETER_SIDE_SLIDING"

	ROLE_WIND_JAMB_SIDE_SLIDING    = "WIND_JAMB_SIDE_SLIDING"
	ROLE_WIND_JAMB_MEETING_SLIDING = "WIND_JAMB_MEETING_SLIDING"
	ROLE_WIND_RAIL_TOP_SLIDING     = "WIND_RAIL_TOP_SLIDING"
	ROLE_WIND_RAIL_BOTTOM_SLIDING  = "WIND_RAIL_BOTTOM_SLIDING"

	ROLE_WIND_VERTICAL_OVERLAP_SLIDING = "WIND_VERTICAL_OVERLAP_SLIDING" // Perfil de traslapo adicional

	// --- (Roles para Casement y otros sistemas se añadirán después) ---
)

// Positions (para claves en Frame.Details y Wind.Details)
const (
	POSITION_TOP    = "Top"
	POSITION_BOTTOM = "Bottom"
	POSITION_LEFT   = "Left"
	POSITION_RIGHT  = "Right"
	POSITION_MIDDLE = "Middle"
)
