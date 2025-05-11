package repositories

import (
	"context" // El contexto se sigue recibiendo en los métodos del repo
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mvialf/windraw/internal/app/window-api/models"
	"github.com/mvialf/windraw/internal/pkg/apiclient"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	// Duraciones de caché
	defaultCacheExpiration      = 30 * time.Minute
	shortCacheExpiration        = 5 * time.Minute
	defaultCacheCleanupInterval = 1 * time.Hour

	// Prefijos de Clave de Caché
	cachePrefixProfileByID           = "catalog:profile:id:"
	cachePrefixProfileBySKU          = "catalog:profile:sku:"
	cachePrefixProfileSystemByID     = "catalog:profilesystem:id:"
	cachePrefixProfileSystemsByType  = "catalog:profilesystems:type:"
	cachePrefixSystemProfileList     = "catalog:systemprofilelist:systemid:"
	cachePrefixColorByID             = "catalog:color:id:"
	cachePrefixSystemAvailableColors = "catalog:systemavailablecolors:systemid:"
	cachePrefixStockItem             = "catalog:stockitem:profileid_colorid:"
	cachePrefixStockItemsForProfile  = "catalog:stockitems:profileid:"
	cachePrefixMaterialByID          = "catalog:material:id:"
	cachePrefixSupplierByID          = "catalog:supplier:id:"
	cachePrefixProfileReinforcements = "catalog:profilereinforcements:mainprofileid:"
)

type supabaseProfileCatalogRepository struct {
	supabaseClient *apiclient.SupabaseClient
	cache          *cache.Cache
	logger         *logrus.Entry
}

func NewSupabaseProfileCatalogRepository(client *apiclient.SupabaseClient, logger *logrus.Logger) ProfileCatalogRepository {
	repoLogger := logger.WithField("repository", "supabase_profile_catalog")
	return &supabaseProfileCatalogRepository{
		supabaseClient: client,
		cache:          cache.New(defaultCacheExpiration, defaultCacheCleanupInterval),
		logger:         repoLogger,
	}
}

// --- Implementación de Métodos para Profile ---

func (r *supabaseProfileCatalogRepository) GetProfileByID(ctx context.Context, profileID int64) (*models.Profile, error) {
	cacheKey := cachePrefixProfileByID + strconv.FormatInt(profileID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetProfileByID", "profileID": profileID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if profile, ok := cachedData.(*models.Profile); ok {
			log.Debug("Cache HIT")
			return profile, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/profiles"
	queryParams := fmt.Sprintf("profile_id=eq.%d&select=profile_id,profile_sku,profile_name,profile_type,material_id,supplier_id,profile_weigth_meter,profile_h,profile_w,profile_h1,profile_h2,profile_h3,profile_h4,profile_w1,profile_cut_type,cut_margin_mm,track_count,uses_overlap,profile_position,profile_structure,weld_margin,created_at,updated_at", profileID)

	var profiles []models.Profile
	if err := r.supabaseClient.QueryData(path, queryParams, &profiles); err != nil { // No pasar ctx aquí
		log.WithError(err).Error("Error obteniendo perfil de Supabase")
		return nil, fmt.Errorf("error obteniendo perfil ID %d de Supabase: %w", profileID, err)
	}

	if len(profiles) == 0 {
		log.Info("Perfil no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	profileToCache := &profiles[0]
	r.cache.Set(cacheKey, profileToCache, defaultCacheExpiration)
	log.Info("Perfil obtenido de Supabase y guardado en caché")
	return profileToCache, nil
}

func (r *supabaseProfileCatalogRepository) GetProfileBySKU(ctx context.Context, sku string) (*models.Profile, error) {
	cacheKey := cachePrefixProfileBySKU + sku
	log := r.logger.WithFields(logrus.Fields{"method": "GetProfileBySKU", "sku": sku, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if profile, ok := cachedData.(*models.Profile); ok {
			log.Debug("Cache HIT")
			return profile, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/profiles"
	queryParams := fmt.Sprintf("profile_sku=eq.%s&select=profile_id,profile_sku,profile_name,profile_type,material_id,supplier_id,profile_weigth_meter,profile_h,profile_w,profile_h1,profile_h2,profile_h3,profile_h4,profile_w1,profile_cut_type,cut_margin_mm,track_count,uses_overlap,profile_position,profile_structure,weld_margin,created_at,updated_at", sku)

	var profiles []models.Profile
	if err := r.supabaseClient.QueryData(path, queryParams, &profiles); err != nil { // No pasar ctx aquí
		log.WithError(err).Error("Error obteniendo perfil por SKU de Supabase")
		return nil, fmt.Errorf("error obteniendo perfil SKU '%s' de Supabase: %w", sku, err)
	}

	if len(profiles) == 0 {
		log.Info("Perfil SKU no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	profileToCache := &profiles[0]
	r.cache.Set(cacheKey, profileToCache, defaultCacheExpiration)
	log.Info("Perfil SKU obtenido de Supabase y guardado en caché")
	return profileToCache, nil
}

// --- Implementación de Métodos para ProfileSystem ---

func (r *supabaseProfileCatalogRepository) GetProfileSystemByID(ctx context.Context, systemID int64) (*models.ProfileSystem, error) {
	cacheKey := cachePrefixProfileSystemByID + strconv.FormatInt(systemID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetProfileSystemByID", "systemID": systemID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if system, ok := cachedData.(*models.ProfileSystem); ok {
			log.Debug("Cache HIT")
			return system, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/profile_systems"
	queryParams := fmt.Sprintf("system_id=eq.%d&select=system_id,name,supplier_id,type,material_id,uses_glass_bead,glass_margin_mm,top_overlap_mm,bottom_overlap_mm,side_overlap_mm,prymacy,created_at,updated_at", systemID)

	var systems []models.ProfileSystem
	if err := r.supabaseClient.QueryData(path, queryParams, &systems); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo sistema de perfiles de Supabase")
		return nil, fmt.Errorf("error obteniendo sistema de perfiles ID %d: %w", systemID, err)
	}

	if len(systems) == 0 {
		log.Info("Sistema de perfiles no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	systemToCache := &systems[0]
	r.cache.Set(cacheKey, systemToCache, defaultCacheExpiration)
	log.Info("Sistema de perfiles obtenido de Supabase y guardado en caché")
	return systemToCache, nil
}

func (r *supabaseProfileCatalogRepository) GetProfileSystemsByType(ctx context.Context, systemType string) ([]models.ProfileSystem, error) {
	cacheKey := cachePrefixProfileSystemsByType + systemType
	log := r.logger.WithFields(logrus.Fields{"method": "GetProfileSystemsByType", "systemType": systemType, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if systems, ok := cachedData.([]models.ProfileSystem); ok {
			log.Debug("Cache HIT")
			return systems, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/profile_systems"
	queryParams := fmt.Sprintf("type=eq.%s&select=system_id,name,supplier_id,type,material_id,uses_glass_bead,glass_margin_mm,top_overlap_mm,bottom_overlap_mm,side_overlap_mm,prymacy,created_at,updated_at&order=prymacy.asc,system_id.asc", systemType)

	var systems []models.ProfileSystem
	if err := r.supabaseClient.QueryData(path, queryParams, &systems); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo sistemas de perfiles por tipo de Supabase")
		return nil, fmt.Errorf("error obteniendo sistemas de perfiles tipo '%s': %w", systemType, err)
	}

	r.cache.Set(cacheKey, systems, defaultCacheExpiration)
	log.Infof("Sistemas de perfiles tipo '%s' obtenidos. Total: %d", systemType, len(systems))
	return systems, nil
}

// --- Implementación de Métodos para SystemProfileListItem ---

func (r *supabaseProfileCatalogRepository) GetSystemProfileListItems(ctx context.Context, systemID int64, elementPartRole *string) ([]models.SystemProfileListItem, error) {
	var sb strings.Builder
	sb.WriteString(cachePrefixSystemProfileList)
	sb.WriteString(strconv.FormatInt(systemID, 10))
	if elementPartRole != nil && *elementPartRole != "" {
		sb.WriteString(":role:")
		sb.WriteString(*elementPartRole)
	}
	cacheKey := sb.String()

	logFields := logrus.Fields{"method": "GetSystemProfileListItems", "systemID": systemID, "cache_key": cacheKey}
	if elementPartRole != nil && *elementPartRole != "" {
		logFields["elementPartRole"] = *elementPartRole
	}
	log := r.logger.WithFields(logFields)

	if cachedData, found := r.cache.Get(cacheKey); found {
		if items, ok := cachedData.([]models.SystemProfileListItem); ok {
			log.Debug("Cache HIT")
			return items, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/system_profile_list"
	baseQuery := fmt.Sprintf("system_id=eq.%d&select=system_id,profile_id,primacy,element_part_role&order=primacy.desc,profile_id.asc", systemID)

	var finalQuery string
	if elementPartRole != nil && *elementPartRole != "" {
		finalQuery = fmt.Sprintf("%s&element_part_role=eq.%s", baseQuery, *elementPartRole)
	} else {
		finalQuery = baseQuery
	}

	var items []models.SystemProfileListItem
	if err := r.supabaseClient.QueryData(path, finalQuery, &items); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo lista de perfiles de sistema de Supabase")
		return nil, fmt.Errorf("error obteniendo lista de perfiles para sistema ID %d: %w", systemID, err)
	}

	r.cache.Set(cacheKey, items, defaultCacheExpiration)
	log.Infof("Lista de perfiles de sistema ID %d obtenida. Total: %d", systemID, len(items))
	return items, nil
}

// --- Implementación de Métodos para Color y SystemAvailableColor ---

func (r *supabaseProfileCatalogRepository) GetColorByID(ctx context.Context, colorID int64) (*models.Color, error) {
	cacheKey := cachePrefixColorByID + strconv.FormatInt(colorID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetColorByID", "colorID": colorID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if color, ok := cachedData.(*models.Color); ok {
			log.Debug("Cache HIT")
			return color, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/colors"
	queryParams := fmt.Sprintf("color_id=eq.%d&select=color_id,name,hex_code", colorID)

	var colors []models.Color
	if err := r.supabaseClient.QueryData(path, queryParams, &colors); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo color de Supabase")
		return nil, fmt.Errorf("error obteniendo color ID %d de Supabase: %w", colorID, err)
	}

	if len(colors) == 0 {
		log.Info("Color no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	colorToCache := &colors[0]
	r.cache.Set(cacheKey, colorToCache, defaultCacheExpiration)
	log.Info("Color obtenido de Supabase y guardado en caché")
	return colorToCache, nil
}

func (r *supabaseProfileCatalogRepository) GetSystemAvailableColors(ctx context.Context, systemID int64) ([]models.SystemAvailableColor, error) {
	cacheKey := cachePrefixSystemAvailableColors + strconv.FormatInt(systemID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetSystemAvailableColors", "systemID": systemID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if items, ok := cachedData.([]models.SystemAvailableColor); ok {
			log.Debug("Cache HIT")
			return items, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/system_available_colors"
	queryParams := fmt.Sprintf("system_id=eq.%d&select=system_id,color_id,color_code_suffix", systemID)

	var items []models.SystemAvailableColor
	if err := r.supabaseClient.QueryData(path, queryParams, &items); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo colores disponibles para sistema de Supabase")
		return nil, fmt.Errorf("error obteniendo colores disponibles para sistema ID %d: %w", systemID, err)
	}

	r.cache.Set(cacheKey, items, defaultCacheExpiration)
	log.Infof("Colores disponibles para sistema ID %d obtenidos. Total: %d", systemID, len(items))
	return items, nil
}

// --- Implementación de Métodos para StockItem ---

func (r *supabaseProfileCatalogRepository) GetStockItem(ctx context.Context, profileID int64, colorID int64) (*models.StockItem, error) {
	cacheKey := cachePrefixStockItem + strconv.FormatInt(profileID, 10) + "_c_" + strconv.FormatInt(colorID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetStockItem", "profileID": profileID, "colorID": colorID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if item, ok := cachedData.(*models.StockItem); ok {
			log.Debug("Cache HIT")
			return item, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/stock_items"
	queryParams := fmt.Sprintf("profile_id=eq.%d&color_id=eq.%d&select=stock_item_id,profile_id,color_id,supplier_id,item_sku,profile_price,profile_length", profileID, colorID)

	var items []models.StockItem
	if err := r.supabaseClient.QueryData(path, queryParams, &items); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo ítem de stock de Supabase")
		return nil, fmt.Errorf("error obteniendo ítem de stock para profileID %d, colorID %d: %w", profileID, colorID, err)
	}

	if len(items) == 0 {
		log.Info("Ítem de stock no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	itemToCache := &items[0]
	r.cache.Set(cacheKey, itemToCache, defaultCacheExpiration)
	log.Info("Ítem de stock obtenido de Supabase y guardado en caché")
	return itemToCache, nil
}

func (r *supabaseProfileCatalogRepository) GetStockItemsForProfile(ctx context.Context, profileID int64) ([]models.StockItem, error) {
	cacheKey := cachePrefixStockItemsForProfile + strconv.FormatInt(profileID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetStockItemsForProfile", "profileID": profileID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if items, ok := cachedData.([]models.StockItem); ok {
			log.Debug("Cache HIT")
			return items, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/stock_items"
	queryParams := fmt.Sprintf("profile_id=eq.%d&select=stock_item_id,profile_id,color_id,supplier_id,item_sku,profile_price,profile_length", profileID)

	var items []models.StockItem
	if err := r.supabaseClient.QueryData(path, queryParams, &items); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo ítems de stock para perfil de Supabase")
		return nil, fmt.Errorf("error obteniendo ítems de stock para profileID %d: %w", profileID, err)
	}

	r.cache.Set(cacheKey, items, defaultCacheExpiration)
	log.Infof("Ítems de stock para perfil ID %d obtenidos. Total: %d", profileID, len(items))
	return items, nil
}

// --- Implementación de Métodos para Material ---

func (r *supabaseProfileCatalogRepository) GetMaterialByID(ctx context.Context, materialID int64) (*models.Material, error) {
	cacheKey := cachePrefixMaterialByID + strconv.FormatInt(materialID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetMaterialByID", "materialID": materialID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if material, ok := cachedData.(*models.Material); ok {
			log.Debug("Cache HIT")
			return material, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/materials"
	queryParams := fmt.Sprintf("material_id=eq.%d&select=material_id,name", materialID)

	var materials []models.Material
	if err := r.supabaseClient.QueryData(path, queryParams, &materials); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo material de Supabase")
		return nil, fmt.Errorf("error obteniendo material ID %d de Supabase: %w", materialID, err)
	}

	if len(materials) == 0 {
		log.Info("Material no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	materialToCache := &materials[0]
	r.cache.Set(cacheKey, materialToCache, defaultCacheExpiration)
	log.Info("Material obtenido de Supabase y guardado en caché")
	return materialToCache, nil
}

// --- Implementación de Métodos para Supplier ---

func (r *supabaseProfileCatalogRepository) GetSupplierByID(ctx context.Context, supplierID int64) (*models.Supplier, error) {
	cacheKey := cachePrefixSupplierByID + strconv.FormatInt(supplierID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetSupplierByID", "supplierID": supplierID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if cachedData == nil {
			log.Debug("Cache HIT (not found)")
			return nil, nil
		}
		if supplier, ok := cachedData.(*models.Supplier); ok {
			log.Debug("Cache HIT")
			return supplier, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/suppliers"
	queryParams := fmt.Sprintf("supplier_id=eq.%d&select=supplier_id,name", supplierID)

	var suppliers []models.Supplier
	if err := r.supabaseClient.QueryData(path, queryParams, &suppliers); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo proveedor de Supabase")
		return nil, fmt.Errorf("error obteniendo proveedor ID %d de Supabase: %w", supplierID, err)
	}

	if len(suppliers) == 0 {
		log.Info("Proveedor no encontrado en Supabase")
		r.cache.Set(cacheKey, nil, shortCacheExpiration)
		return nil, nil
	}

	supplierToCache := &suppliers[0]
	r.cache.Set(cacheKey, supplierToCache, defaultCacheExpiration)
	log.Info("Proveedor obtenido de Supabase y guardado en caché")
	return supplierToCache, nil
}

// --- Implementación de Métodos para ProfileReinforcement ---

func (r *supabaseProfileCatalogRepository) GetProfileReinforcements(ctx context.Context, mainProfileID int64) ([]models.ProfileReinforcement, error) {
	cacheKey := cachePrefixProfileReinforcements + strconv.FormatInt(mainProfileID, 10)
	log := r.logger.WithFields(logrus.Fields{"method": "GetProfileReinforcements", "mainProfileID": mainProfileID, "cache_key": cacheKey})

	if cachedData, found := r.cache.Get(cacheKey); found {
		if items, ok := cachedData.([]models.ProfileReinforcement); ok {
			log.Debug("Cache HIT")
			return items, nil
		}
		log.Warn("Tipo de dato incorrecto en caché. Eliminando.")
		r.cache.Delete(cacheKey)
	}

	log.Debug("Cache MISS")
	path := "/rest/v1/profile_reinforcements"
	queryParams := fmt.Sprintf("main_profile_id=eq.%d&select=main_profile_id,reinforcement_profile_id,reinforcement_gap_mm", mainProfileID)

	var items []models.ProfileReinforcement
	if err := r.supabaseClient.QueryData(path, queryParams, &items); err != nil { // No pasar ctx
		log.WithError(err).Error("Error obteniendo refuerzos de perfil de Supabase")
		return nil, fmt.Errorf("error obteniendo refuerzos para perfil principal ID %d: %w", mainProfileID, err)
	}

	r.cache.Set(cacheKey, items, defaultCacheExpiration)
	log.Infof("Refuerzos para perfil principal ID %d obtenidos. Total: %d", mainProfileID, len(items))
	return items, nil
}
