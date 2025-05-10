package repositories

import (
	"context"

	"github.com/mvialf/windraw/internal/app/window-api/models" // Asumiendo que tienes un models.Profile
)

// ProfileCatalogRepository define las operaciones de consulta para el catálogo de perfiles.
type ProfileCatalogRepository interface {
	GetAllProfiles(ctx context.Context) ([]models.Profile, error)
	GetProfileBySKU(ctx context.Context, sku string) (*models.Profile, error)
	// Podrías tener más métodos como GetProfilesByMaterial(ctx context.Context, material string) ([]models.Profile, error)
}
