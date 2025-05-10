package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/mvialf/windraw/internal/app/window-api/models" // Asegúrate que models.Profile esté definido
	"github.com/mvialf/windraw/internal/pkg/apiclient"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

const (
	// TTLs específicos para el catálogo de perfiles (podrían ser más largos)
	profilesCacheDefaultExpiration = 1 * time.Hour // Perfiles no cambian tan a menudo
	profilesCacheCleanupInterval   = 2 * time.Hour

	// Prefijos y claves para el caché de perfiles
	allProfilesCacheKey     = "catalog:all_profiles"
	profileBySKUCachePrefix = "catalog:profile_sku:"
)

// supabaseProfileCatalogRepository implementa ProfileCatalogRepository con Supabase y caché.
type supabaseProfileCatalogRepository struct {
	supabaseClient *apiclient.SupabaseClient
	cache          *cache.Cache
	logger         *logrus.Entry // Usar logrus.Entry para logging contextualizado
}

// NewSupabaseProfileCatalogRepository crea una nueva instancia del repositorio de catálogo de perfiles.
func NewSupabaseProfileCatalogRepository(client *apiclient.SupabaseClient, logger *logrus.Logger) ProfileCatalogRepository {
	// Creamos un logger contextualizado para este repositorio específico
	repoLogger := logger.WithField("repository", "profile_catalog")
	return &supabaseProfileCatalogRepository{
		supabaseClient: client,
		cache:          cache.New(profilesCacheDefaultExpiration, profilesCacheCleanupInterval),
		logger:         repoLogger,
	}
}

// GetAllProfiles obtiene todos los perfiles del catálogo, utilizando caché.
func (r *supabaseProfileCatalogRepository) GetAllProfiles(ctx context.Context) ([]models.Profile, error) {
	log := r.logger.WithField("method", "GetAllProfiles")

	// 1. Intentar obtener del caché
	if cachedData, found := r.cache.Get(allProfilesCacheKey); found {
		if profiles, ok := cachedData.([]models.Profile); ok {
			log.Info("Cache HIT")
			return profiles, nil
		}
		log.Warn("Tipo de dato incorrecto en caché para todos los perfiles. Se eliminará.")
		r.cache.Delete(allProfilesCacheKey)
	}

	log.Info("Cache MISS")

	// 2. Si no está en caché, obtener de Supabase
	var profiles []models.Profile
	// AJUSTA EL NOMBRE DE TU TABLA DE PERFILES EN SUPABASE SI ES NECESARIO
	supabasePath := "/rest/v1/profiles_catalog"
	// Es mejor seleccionar columnas explícitas en lugar de "select=*"
	queryParams := "select=id,sku,description,material,weight_per_meter,available_colors,created_at,updated_at" // Ejemplo

	if err := r.supabaseClient.QueryData(supabasePath, queryParams, &profiles); err != nil {
		log.WithError(err).Error("Error obteniendo perfiles de Supabase")
		return nil, fmt.Errorf("error obteniendo perfiles de Supabase: %w", err)
	}

	// 3. Guardar en caché
	r.cache.Set(allProfilesCacheKey, profiles, cache.DefaultExpiration) // Usa el default del caché para este repo
	log.Infof("Perfiles obtenidos de Supabase y guardados en caché. Total: %d", len(profiles))

	return profiles, nil
}

// GetProfileBySKU obtiene un perfil específico por su SKU, utilizando caché.
func (r *supabaseProfileCatalogRepository) GetProfileBySKU(ctx context.Context, sku string) (*models.Profile, error) {
	cacheKey := profileBySKUCachePrefix + sku
	log := r.logger.WithFields(logrus.Fields{"method": "GetProfileBySKU", "sku": sku, "cache_key": cacheKey})

	// 1. Intentar obtener del caché
	if cachedData, found := r.cache.Get(cacheKey); found {
		// Manejar el caso donde el caché podría tener un nil explícito para "no encontrado"
		if cachedData == nil {
			log.Info("Cache HIT (not found)")
			return nil, nil // O un error específico como models.ErrNotFound
		}
		if profile, ok := cachedData.(*models.Profile); ok {
			log.Info("Cache HIT")
			return profile, nil
		}
		log.Warn("Tipo de dato incorrecto en caché para perfil SKU. Se eliminará.")
		r.cache.Delete(cacheKey)
	}

	log.Info("Cache MISS")

	// 2. Si no está en caché, obtener de Supabase
	var profiles []models.Profile // Supabase generalmente devuelve un array
	// AJUSTA EL NOMBRE DE TU TABLA DE PERFILES EN SUPABASE SI ES NECESARIO
	supabasePath := "/rest/v1/profiles_catalog"
	// Es mejor seleccionar columnas explícitas. Asegúrate que los nombres de columna coincidan con tu struct models.Profile.
	supabaseQueryParams := fmt.Sprintf("sku=eq.%s&select=id,sku,description,material,weight_per_meter,available_colors,created_at,updated_at&limit=1", sku)

	if err := r.supabaseClient.QueryData(supabasePath, supabaseQueryParams, &profiles); err != nil {
		log.WithError(err).Error("Error obteniendo perfil por SKU de Supabase")
		return nil, fmt.Errorf("error obteniendo perfil SKU '%s' de Supabase: %w", sku, err)
	}

	if len(profiles) == 0 {
		log.Info("Perfil no encontrado en Supabase")
		// "Negative caching": cachear que no se encontró para evitar llamadas repetidas por un tiempo.
		r.cache.Set(cacheKey, nil, 5*time.Minute) // Cachea un nil explícito con un TTL corto
		return nil, nil                           // O un error models.ErrNotFound si prefieres manejarlo así
	}

	profileToCache := &profiles[0]

	// 3. Guardar en caché
	r.cache.Set(cacheKey, profileToCache, cache.DefaultExpiration) // Usa el default del caché para este repo
	log.Info("Perfil obtenido de Supabase y guardado en caché")

	return profileToCache, nil
}
