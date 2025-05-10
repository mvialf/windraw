package models

// AuxProfile define un perfil auxiliar que puede usarse para unir módulos, por ejemplo.
type AuxProfile struct {
	ID    string  `json:"id"`              // SKU o identificador del perfil auxiliar
	Angle float64 `json:"angle,omitempty"` // Ángulo de unión o del perfil, si aplica
	Notes string  `json:"notes,omitempty"` // Notas adicionales sobre el perfil auxiliar
}

// Module representa una agrupación de Elements que pueden estar unidos (opcionalmente) por un AuxProfile.
// Por ejemplo, una ventana fija unida a una puerta corredera.
type Module struct {
	ID         string      `json:"id"`                    // ID único del módulo
	Elements   []Element   `json:"elements"`              // Lista de elementos que componen este módulo (SUGERENCIA: añadido tag JSON)
	AuxProfile *AuxProfile `json:"aux_profile,omitempty"` // Perfil auxiliar para unir este módulo al siguiente, si aplica
}

// Component es una agrupación lógica de Modules dentro de un proyecto.
// Por ejemplo, "Ventanas Fachada Norte" o "Puertas Terraza".
type Component struct {
	ID      string   `json:"id"`      // ID único del componente
	Modules []Module `json:"modules"` // Lista de módulos que componen este componente
}

// Aquí podrías añadir constructores para Component, Module, y AuxProfile si lo necesitas,
// de forma similar a como los tienes para Project, Element, etc.
// Ejemplo (muy básico, necesitarías pasar los datos necesarios):

/*
func NewModule(id string, elements []Element, auxProfile *AuxProfile) (*Module, error) {
	if id == "" {
		id = generateID() // Asumiendo que generateID() está disponible en el paquete models
	}
	if elements == nil {
		elements = []Element{}
	}
	return &Module{
		ID:         id,
		Elements:   elements,
		AuxProfile: auxProfile,
	}, nil
}

func NewComponent(id string, modules []Module) (*Component, error) {
	if id == "" {
		id = generateID() // Asumiendo que generateID() está disponible
	}
	if modules == nil {
		modules = []Module{}
	}
	return &Component{
		ID:      id,
		Modules: modules,
	}, nil
}
*/
