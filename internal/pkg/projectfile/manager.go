package projectfile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	// No necesitamos "github.com/google/uuid" aquí si los IDs de Project son strings
	"github.com/mvialf/windraw/internal/app/window-api/models" // Tus modelos
)

var (
	ErrProjectDataMissing = fmt.Errorf("projectfile: project data missing")
	ErrInvalidJSONFormat  = fmt.Errorf("projectfile: invalid JSON format")
)

func sanitizeFilename(name string) string {
	if name == "" {
		return "untitled"
	}
	re := regexp.MustCompile(`[^\w\s.-]`)
	name = re.ReplaceAllString(name, "")
	name = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(name), " ")
	return name
}

// GenerateProjectFilename genera un nombre de archivo para el proyecto.
// El formato es: "ProjectID - ClientName.json".
// Usa Project.ID (string) y Project.Contact.Name.
func GenerateProjectFilename(project *models.Project) (string, error) {
	// Ajustado para Project.ID como string
	if project == nil || project.ID == "" {
		return "", fmt.Errorf("%w: project or project ID is empty for filename generation", ErrProjectDataMissing)
	}
	clientName := "UnknownClient"
	if project.Contact.Name != "" {
		clientName = project.Contact.Name
	}

	sanitizedClientName := sanitizeFilename(clientName)
	// Project.ID ya es un string
	filename := fmt.Sprintf("%s - %s.json", project.ID, sanitizedClientName)
	return filename, nil
}

// SaveProject guarda la estructura del proyecto en un archivo JSON.
func SaveProject(project *models.Project, directoryPath string) (string, error) {
	if project == nil {
		return "", fmt.Errorf("%w: cannot save a nil project", ErrProjectDataMissing)
	}

	filename, err := GenerateProjectFilename(project)
	if err != nil {
		return "", fmt.Errorf("projectfile: failed to generate filename: %w", err)
	}

	filePath := filepath.Join(directoryPath, filename)

	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return "", fmt.Errorf("projectfile: failed to create directory %s: %w", directoryPath, err)
	}

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		// Si Project.ID es string, el error aquí no será por uuid.Nil
		return "", fmt.Errorf("projectfile: failed to marshal project %s to JSON: %w", project.ID, err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("projectfile: failed to write project file %s: %w", filePath, err)
	}

	return filePath, nil
}

// LoadProject carga un proyecto desde un archivo JSON.
func LoadProject(filePath string) (*models.Project, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("projectfile: project file %s not found: %w", filePath, err)
		}
		return nil, fmt.Errorf("projectfile: failed to read project file %s: %w", filePath, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("%w: project file %s is empty", ErrInvalidJSONFormat, filePath)
	}

	var project models.Project
	err = json.Unmarshal(data, &project)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to unmarshal JSON from %s: %w", ErrInvalidJSONFormat, filePath, err)
	}

	// Ajustado para Project.ID como string
	if project.ID == "" {
		// Podrías decidir si un proyecto cargado sin ID es un error crítico.
		// return nil, fmt.Errorf("%w: loaded project from %s has an empty ID", ErrInvalidJSONFormat, filePath)
	}

	return &project, nil
}
