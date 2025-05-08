package constants

const (
	IVA_RATE = 19.0
)

const (
	MATERIAL_PVC      = "PVC"
	MATERIAL_ALUMINIO = "Aluminio"
	MATERIAL_CRISTAL  = "Cristal"
	MATERIAL_MADERA   = "Madera"
	MATERIAL_ACERO    = "Acero"
)

const (
	TYPE_SLIDING  = "Sliding"
	TYPE_CASEMENT = "Casement"
)

const (
	PROFILE_TYPE_SLIDING_FRAME    = "Marco corredera"
	PROFILE_TYPE_SLIDING_WIND     = "Hoja corredera"
	PROFILE_TYPE_SLIDING_OVERLAP  = "Traslapo corredera"
	PROFILE_TYPE_SLIDING_ADAPTER  = "Adaptador de hoja corredera"
	PROFILE_TYPE_SLIDING_RAIL     = "Riel corredera"
	PROFILE_TYPE_CASEMENT_FRAME   = "Marco fijo"
	PROFILE_TYPE_CASEMENT_WIND_IN = "Hoja abatir interior"
	PROFILE_TYPE_CASEMENT_WIND_OUT= "Hoja abatir exterior"
	PROFILE_TYPE_CASEMENT_ADAPTER = "Adaptador de hoja abatir"
	PROFILE_TYPE_MULLION          = "Montante"
	PROFILE_TYPE_GLASSBEAD        = "Junquillo"
	PROFILE_TYPE_JOINT            = "Unión"
	PROFILE_TYPE_AUXILIAR         = "Perfil auxiliar"
	PROFILE_TYPE_REINFORCEMENT    = "Refuerzo"
)

const (
	STRUCTURE_VENTANA = "Ventana"
	STRUCTURE_PUERTA  = "Puerta"
)

const (
	GEOMETRY_RECTANGULAR = "Rectangular"
	GEOMETRY_TRIANGULAR  = "Triangular"
)

const (
	POSITION_LEFT   = "Izquierda"
	POSITION_BOTTOM = "Abajo"
	POSITION_RIGHT  = "Derecha"
	POSITION_TOP    = "Arriba"
)

const (
	CUT_SQUARE = "Cuadrado"
	CUT_ANGLE  = "Ángulo"
)

const (
	CUT_SQUARE_WIND             = "Cuadrado"
	CUT_HORIZONTAL_WIND         = "Horizontal"
	CUT_VERTICAL_WIND           = "Vertical"
	CUT_ANGLE_WIND              = "Ángulo"
	CUT_HORIZONTAL_OVERLAP_WIND = "Horizontal superpuesto"
	CUT_VERTICAL_OVERLAP_WIND   = "Vertical superpuesto"
	CUT_CUSTOM_WIND             = "Personalizado"
)

const (
	WIND_KIND_SLIDING_MOVIL   = "Hoja corredera móvil"
	WIND_KIND_SLIDING_FIXED   = "Hoja corredera fija"
	WIND_KIND_FIXED           = "Hoja fija"
	WIND_KIND_CASEMENT        = "Hoja abatir"
	WIND_KIND_PROJECTING      = "Hoja proyectante"
	WIND_KIND_TILT_TURN       = "Hoja oscilobatiente"
	WIND_KIND_TILT_ONLY       = "Hoja oscilante"
)

const (
	WIND_STATUS_ACTIVE   = "activa"
	WIND_STATUS_INACTIVE = "inactiva"
)

const (
	OPENING_SIDE_RIGHT  = "Derecha"
	OPENING_SIDE_LEFT   = "Izquierda"
	OPENING_SIDE_BOTTOM = "Abajo"
	OPENING_SIDE_TOP    = "Arriba"
)

const (
	OPENING_INTERIOR = "interior"
	OPENING_EXTERIOR = "exterior"
)