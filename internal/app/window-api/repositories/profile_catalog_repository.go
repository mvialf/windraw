package repositories

import (
	"context"

	"github.com/mvialf/windraw/internal/app/window-api/models"
)

type ProfileCatalogRepository interface {
	GetProfileByID(ctx context.Context, profileID int64) (*models.Profile, error)
	GetProfileBySKU(ctx context.Context, sku string) (*models.Profile, error)

	GetProfileSystemByID(ctx context.Context, systemID int64) (*models.ProfileSystem, error)
	GetProfileSystemsByType(ctx context.Context, systemType string) ([]models.ProfileSystem, error)

	GetSystemProfileListItems(ctx context.Context, systemID int64, elementPartRole *string) ([]models.SystemProfileListItem, error)

	GetColorByID(ctx context.Context, colorID int64) (*models.Color, error)
	GetSystemAvailableColors(ctx context.Context, systemID int64) ([]models.SystemAvailableColor, error)

	GetStockItem(ctx context.Context, profileID int64, colorID int64) (*models.StockItem, error)
	GetStockItemsForProfile(ctx context.Context, profileID int64) ([]models.StockItem, error)

	GetMaterialByID(ctx context.Context, materialID int64) (*models.Material, error)
	GetSupplierByID(ctx context.Context, supplierID int64) (*models.Supplier, error)
	GetProfileReinforcements(ctx context.Context, mainProfileID int64) ([]models.ProfileReinforcement, error)
}
