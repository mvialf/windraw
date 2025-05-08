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

func generateID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("error-id-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

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
	Type     ContactType `json:"type"`
	Name     string      `json:"name"`
	Phone    string      `json:"phone,omitempty"`
	Email    string      `json:"email,omitempty"`
	Address  string      `json:"address,omitempty"`
	District string      `json:"district,omitempty"`
	City     string      `json:"city,omitempty"`
}

type ProjectCost struct {
	Name         string  `json:"name"`
	IsPercentage bool    `json:"is_percentage"`
	Value        float64 `json:"value"`
}

type Project struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	CreatedAt  time.Time     `json:"created_at"`
	Contact    Contact       `json:"contact"`
	Costs      []ProjectCost `json:"costs"`
	Components []Component
	IvaRate    float64       `json:"iva_rate"`
}

func validateContact(contact Contact) error {
	if contact.Name == "" {
		return errors.New("el nombre del contacto (Contact.Name) no puede estar vacío")
	}
	return nil
}

func NewProject(name string, contact Contact, costs []ProjectCost, components []Component, ivaRate float64) (*Project, error) {
	if name == "" {
		return nil, errors.New("el nombre del proyecto (Project.Name) no puede estar vacío")
	}
	if err := validateContact(contact); err != nil {
		return nil, fmt.Errorf("información de contacto inválida: %w", err)
	}
	if costs == nil {
		costs = []ProjectCost{}
	}
	if components == nil {
		components = []Component{}
	}

	project := &Project{
		ID:         generateID(),
		Name:       name,
		CreatedAt:  time.Now(),
		Contact:    contact,
		Costs:      costs,
		Components: components,
		IvaRate:    ivaRate,
	}
	return project, nil
}

func (p *Project) AddComponent(component Component) {
	p.Components = append(p.Components, component)
}