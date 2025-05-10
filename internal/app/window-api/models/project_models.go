package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"
)

//==============================================================================
// --- Tipos Auxiliares y Helpers (Relevantes para múltiples modelos) ---
//==============================================================================

type ContactType bool

const (
	PERSON  ContactType = false
	COMPANY ContactType = true
)

// generateID crea un ID string basado en crypto/rand.
// Nota: Si decides estandarizar a github.com/google/uuid, esta función y el tipo de Project.ID cambiarían.
func generateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// En caso de error con rand.Read, devuelve un ID basado en tiempo para evitar pánico,
		// aunque esto es altamente improbable. Considera un log aquí.
		return fmt.Sprintf("error-id-%d", time.Now().UnixNano())
	}
	// Formato similar a un UUID, pero como string directo.
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// IsValidOption verifica si un valor está dentro de una lista de opciones válidas.
// Esta función parece ser un helper general, lo cual está bien aquí si es usado por modelos en este paquete.
func IsValidOption(value string, validOptions []string) bool {
	for _, option := range validOptions {
		if value == option {
			return true
		}
	}
	return false
}

//==============================================================================
// --- Estructuras de Datos Principales (Proyecto y Contacto) ---
//==============================================================================

type Contact struct {
	Type     ContactType `json:"type"`               // false para Persona, true para Empresa
	Name     string      `json:"name"`               // Nombre del contacto o razón social
	Phone    string      `json:"phone,omitempty"`    // Teléfono de contacto
	Email    string      `json:"email,omitempty"`    // Email de contacto
	Address  string      `json:"address,omitempty"`  // Dirección
	District string      `json:"district,omitempty"` // Comuna o Distrito
	City     string      `json:"city,omitempty"`     // Ciudad
}

type ProjectCost struct {
	Name         string  `json:"name"`          // Nombre o descripción del costo
	IsPercentage bool    `json:"is_percentage"` // True si el valor es un porcentaje, false si es un monto fijo
	Value        float64 `json:"value"`         // Valor del costo (monto o porcentaje)
}

// Project define la estructura de un proyecto.
type Project struct {
	ID         string        `json:"id"`                   // ID único del proyecto, generado por generateID()
	Name       string        `json:"name"`                 // Nombre del proyecto
	CreatedAt  time.Time     `json:"created_at"`           // Fecha y hora de creación del proyecto
	Contact    Contact       `json:"contact"`              // Información de contacto del cliente
	Costs      []ProjectCost `json:"costs"`                // Lista de costos adicionales asociados al proyecto
	Components []Component   `json:"components,omitempty"` // Lista de componentes del proyecto (SUGERENCIA: añadido omitempty)
	IvaRate    float64       `json:"iva_rate"`             // Tasa de IVA aplicable al proyecto (ej: 0.19 para 19%)
}

// validateContact valida los campos requeridos de la estructura Contact.
func validateContact(contact Contact) error {
	if contact.Name == "" {
		return errors.New("el nombre del contacto (Contact.Name) no puede estar vacío")
	}
	// Aquí podrías añadir más validaciones para Contact si es necesario (ej. formato de email)
	return nil
}

// NewProject es el constructor para la estructura Project.
// Inicializa un nuevo proyecto con los datos proporcionados y genera un ID y CreatedAt.
func NewProject(name string, contact Contact, costs []ProjectCost, components []Component, ivaRate float64) (*Project, error) {
	if name == "" {
		return nil, errors.New("el nombre del proyecto (Project.Name) no puede estar vacío")
	}
	if err := validateContact(contact); err != nil {
		return nil, fmt.Errorf("información de contacto inválida: %w", err)
	}

	// Asegurar que las slices no sean nil para evitar problemas con marshalling JSON o lógica posterior
	if costs == nil {
		costs = []ProjectCost{}
	}
	if components == nil {
		components = []Component{}
	}

	project := &Project{
		ID:         generateID(), // Usa tu función para generar un ID string
		Name:       name,
		CreatedAt:  time.Now(),
		Contact:    contact,
		Costs:      costs,
		Components: components,
		IvaRate:    ivaRate,
	}
	return project, nil
}

// AddComponent añade un nuevo componente a la lista de componentes del proyecto.
func (p *Project) AddComponent(component Component) {
	p.Components = append(p.Components, component)
}
