package services

import (
	"context"
	"fmt"
	"strings" // Necesario para el strings.Join en el log de errores

	// "sort" // No es estrictamente necesario si el repositorio ya ordena por primacy

	"github.com/mvialf/windraw/internal/app/window-api/models"
	"github.com/mvialf/windraw/internal/app/window-api/repositories"
	"github.com/mvialf/windraw/internal/pkg/constants" // IMPORTANTE: Asegúrate que esta ruta sea correcta
	"github.com/sirupsen/logrus"
)

// SelectedProfileInfo agrupa la información de un perfil seleccionado para una ranura,
// incluyendo su información de stock item (que contiene el ColorID y el precio para ese color/perfil).
type SelectedProfileInfo struct {
	Profile   models.Profile
	StockItem models.StockItem
	// ElementPartRole string // Opcional: podríamos añadir el rol aquí si es útil para el llamador
}

// ProfileSelectorService encapsula la lógica para seleccionar perfiles
// basados en un sistema, tipo de elemento, y otros criterios.
type ProfileSelectorService struct {
	repo   repositories.ProfileCatalogRepository
	logger *logrus.Entry
}

// NewProfileSelectorService crea una nueva instancia de ProfileSelectorService.
func NewProfileSelectorService(repo repositories.ProfileCatalogRepository, logger *logrus.Logger) *ProfileSelectorService {
	return &ProfileSelectorService{
		repo:   repo,
		logger: logger.WithField("service", "profile_selector"),
	}
}

// SelectDefaultSystemAndProfilesInput define los parámetros de entrada para SelectDefaultSystemAndProfiles.
type SelectDefaultSystemAndProfilesInput struct {
	ElementType    string   // Ej. constants.TYPE_SLIDING (debe coincidir con ProfileSystem.Type)
	DesiredColorID int64    // El ID del color deseado por el usuario
	ElementRanuras []string // Lista de "element_part_role" que necesita este elemento. Ej: [constants.ROLE_FRAME_TOP_SLIDING, ...]
}

// SelectDefaultSystemAndProfilesOutput define el resultado de SelectDefaultSystemAndProfiles.
type SelectDefaultSystemAndProfilesOutput struct {
	SelectedSystem     models.ProfileSystem
	SelectedColor      models.Color
	ProfilesForRanuras map[string]SelectedProfileInfo // Mapea un element_part_role al perfil y stock item seleccionados
	Errors             []string                       // Lista de errores/advertencias encontradas durante la selección
}

// selectProfileForSingleRole es una función helper interna para seleccionar el mejor perfil "disponible"
// (con StockItem existente) para un rol específico, sistema y color.
func (s *ProfileSelectorService) selectProfileForSingleRole(
	ctx context.Context,
	systemID int64,
	targetRole string,
	colorID int64,
	systemName string, // Para mensajes de error más claros
) (*SelectedProfileInfo, string, error) { // Devuelve (SelectedProfileInfo, MensajeUsuarioSiNoHaySeleccion, ErrorDeSistema)
	log := s.logger.WithFields(logrus.Fields{
		"helperMethod": "selectProfileForSingleRole",
		"systemID":     systemID,
		"targetRole":   targetRole,
		"colorID":      colorID,
	})
	log.Debug("Intentando seleccionar perfil para rol específico")

	profileListItems, err := s.repo.GetSystemProfileListItems(ctx, systemID, &targetRole)
	if err != nil {
		errMsg := fmt.Sprintf("error obteniendo perfiles para rol '%s' en sistema ID %d: %v", targetRole, systemID, err)
		log.Error(errMsg)
		return nil, "", fmt.Errorf(errMsg) // Error de sistema
	}

	if len(profileListItems) == 0 {
		userMsg := fmt.Sprintf("no se encontraron perfiles candidatos para el rol '%s' en el sistema '%s'", targetRole, systemName)
		log.Warn(userMsg)
		return nil, userMsg, nil
	}

	for _, item := range profileListItems { // Ya están ordenados por primacy por el repo
		profileDetails, err := s.repo.GetProfileByID(ctx, item.ProfileID)
		if err != nil {
			log.WithError(err).Warnf("Error obteniendo detalles del perfil ID %d para rol '%s'. Intentando siguiente.", item.ProfileID, targetRole)
			continue
		}
		if profileDetails == nil {
			log.Warnf("No se encontraron detalles para el perfil ID %d listado para rol '%s'. Intentando siguiente.", item.ProfileID, targetRole)
			continue
		}

		stockItem, err := s.repo.GetStockItem(ctx, profileDetails.ProfileID, colorID)
		if err != nil {
			log.WithError(err).Warnf("Error verificando ítem de stock para perfil ID %d (SKU: %s), color ID %d (rol '%s'). Intentando siguiente.", profileDetails.ProfileID, profileDetails.ProfileSKU, colorID, targetRole)
			continue
		}

		if stockItem != nil { // Perfil "disponible" (existe en stock_items para este color)
			log.Infof("Perfil ID %d (SKU: %s) seleccionado para rol '%s' en color ID %d. Ítem de stock encontrado.",
				profileDetails.ProfileID, profileDetails.ProfileSKU, targetRole, colorID)

			return &SelectedProfileInfo{
				Profile:   *profileDetails,
				StockItem: *stockItem,
			}, "", nil // Éxito
		}
		log.Warnf("No se encontró ítem de stock para perfil ID %d (SKU: %s) en color ID %d (rol '%s'). Intentando siguiente.",
			profileDetails.ProfileID, profileDetails.ProfileSKU, colorID, targetRole)
	}

	userMsg := fmt.Sprintf("no se pudo seleccionar un perfil 'disponible' (con ítem de stock existente) para el rol '%s' en color ID %d para el sistema '%s'", targetRole, colorID, systemName)
	log.Warn(userMsg)
	return nil, userMsg, nil
}

// SelectDefaultSystemAndProfiles selecciona el sistema de perfiles por defecto, el color,
// y los perfiles por defecto para las ranuras especificadas de un elemento.
func (s *ProfileSelectorService) SelectDefaultSystemAndProfiles(
	ctx context.Context,
	input SelectDefaultSystemAndProfilesInput,
) (*SelectDefaultSystemAndProfilesOutput, error) {
	log := s.logger.WithFields(logrus.Fields{
		"method":      "SelectDefaultSystemAndProfiles",
		"elementType": input.ElementType,
		"colorID":     input.DesiredColorID,
	})
	log.Info("Iniciando selección de sistema y perfiles por defecto")

	output := &SelectDefaultSystemAndProfilesOutput{
		ProfilesForRanuras: make(map[string]SelectedProfileInfo),
		Errors:             []string{},
	}

	// 1. Seleccionar el ProfileSystem por defecto
	systems, err := s.repo.GetProfileSystemsByType(ctx, input.ElementType)
	if err != nil {
		log.WithError(err).Error("Error obteniendo sistemas por tipo")
		return nil, fmt.Errorf("error al buscar sistemas de tipo '%s': %w", input.ElementType, err)
	}
	if len(systems) == 0 {
		msg := fmt.Sprintf("no se encontraron sistemas de perfiles para el tipo '%s'", input.ElementType)
		log.Warn(msg)
		output.Errors = append(output.Errors, msg)
		return output, nil
	}
	output.SelectedSystem = systems[0] // El repo ya ordena por primacy.asc, system_id.asc
	log = log.WithField("selectedSystemID", output.SelectedSystem.SystemID)
	log.Infof("Sistema seleccionado: %s (ID: %d)", output.SelectedSystem.Name, output.SelectedSystem.SystemID)

	// 2. Validar y obtener el Color seleccionado
	availableSystemColors, err := s.repo.GetSystemAvailableColors(ctx, output.SelectedSystem.SystemID)
	if err != nil {
		log.WithError(err).Errorf("Error obteniendo colores disponibles para el sistema ID %d", output.SelectedSystem.SystemID)
		return nil, fmt.Errorf("error al obtener colores para el sistema %d: %w", output.SelectedSystem.SystemID, err)
	}
	isColorValidForSystem := false
	for _, asc := range availableSystemColors {
		if asc.ColorID == input.DesiredColorID {
			isColorValidForSystem = true
			break
		}
	}
	if !isColorValidForSystem {
		msg := fmt.Sprintf("el color ID %d no está disponible para el sistema '%s' (ID %d)", input.DesiredColorID, output.SelectedSystem.Name, output.SelectedSystem.SystemID)
		log.Warn(msg)
		output.Errors = append(output.Errors, msg)
		return output, nil
	}
	colorDetails, err := s.repo.GetColorByID(ctx, input.DesiredColorID)
	if err != nil || colorDetails == nil {
		msg := fmt.Sprintf("no se pudieron obtener los detalles para el color ID %d", input.DesiredColorID)
		log.WithError(err).Warn(msg)
		output.Errors = append(output.Errors, msg)
		return output, nil
	}
	output.SelectedColor = *colorDetails
	log.Infof("Color seleccionado: %s (ID: %d)", output.SelectedColor.Name, output.SelectedColor.ColorID)

	// 3. Seleccionar Perfiles para cada Ranura especificada
	if len(input.ElementRanuras) == 0 {
		log.Warn("No se especificaron ranuras de elemento para la selección de perfiles.")
		return output, nil
	}
	log.Infof("Intentando seleccionar perfiles para %d ranuras de input", len(input.ElementRanuras))

	processedRanuras := make(map[string]bool)

	for _, ranuraRole := range input.ElementRanuras {
		if _, ok := processedRanuras[ranuraRole]; ok {
			log.Debugf("Rol '%s' ya procesado o en proceso, omitiendo duplicado en input.", ranuraRole)
			continue
		}

		selectedProfile, userMsg, err := s.selectProfileForSingleRole(ctx, output.SelectedSystem.SystemID, ranuraRole, output.SelectedColor.ColorID, output.SelectedSystem.Name)
		if err != nil {
			output.Errors = append(output.Errors, fmt.Sprintf("Error de sistema procesando rol '%s': %v", ranuraRole, err))
			processedRanuras[ranuraRole] = true
			continue
		}
		if selectedProfile != nil {
			output.ProfilesForRanuras[ranuraRole] = *selectedProfile
			processedRanuras[ranuraRole] = true

			// LÓGICA CONDICIONAL PARA TRASLAPO ADICIONAL
			if ranuraRole == constants.ROLE_WIND_JAMB_MEETING_SLIDING &&
				selectedProfile.Profile.UsesOverlap != nil && *selectedProfile.Profile.UsesOverlap {

				overlapRole := constants.ROLE_WIND_VERTICAL_OVERLAP_SLIDING
				log.Infof("Perfil de encuentro (Rol: %s, SKU: %s) usa traslapo. Buscando perfil para rol de traslapo: %s",
					ranuraRole, selectedProfile.Profile.ProfileSKU, overlapRole)

				if _, overlapProcessed := processedRanuras[overlapRole]; !overlapProcessed { // Evitar buscar traslapo si ya se hizo
					overlapProfileInfo, overlapUserMsg, overlapErr := s.selectProfileForSingleRole(ctx, output.SelectedSystem.SystemID, overlapRole, output.SelectedColor.ColorID, output.SelectedSystem.Name)
					if overlapErr != nil {
						output.Errors = append(output.Errors, fmt.Sprintf("Error de sistema procesando rol de traslapo '%s': %v", overlapRole, overlapErr))
					} else if overlapProfileInfo != nil {
						output.ProfilesForRanuras[overlapRole] = *overlapProfileInfo
						log.Infof("Perfil de traslapo (Rol: %s, SKU: %s) seleccionado.", overlapRole, overlapProfileInfo.Profile.ProfileSKU)
					} else {
						output.Errors = append(output.Errors, overlapUserMsg)
					}
					processedRanuras[overlapRole] = true
				} else {
					log.Debugf("Rol de traslapo '%s' ya procesado.", overlapRole)
				}
			}
		} else {
			output.Errors = append(output.Errors, userMsg)
			processedRanuras[ranuraRole] = true
		}
	}

	log.Info("Selección de sistema y perfiles por defecto completada")
	if len(output.Errors) > 0 {
		log.Warnf("Se encontraron %d problemas durante la selección: %v", len(output.Errors), strings.Join(output.Errors, "; "))
	}
	return output, nil
}

// ListAlternativeProfilesForSlotInput define los parámetros de entrada para ListAlternativeProfilesForSlot.
type ListAlternativeProfilesForSlotInput struct {
	SystemID       int64
	RanuraRole     string // El element_part_role específico, ej: constants.ROLE_FRAME_TOP_SLIDING
	DesiredColorID int64
}

// ListAlternativeProfilesForSlotOutput define el resultado de ListAlternativeProfilesForSlot.
type ListAlternativeProfilesForSlotOutput struct {
	AlternativeProfiles []SelectedProfileInfo
}

// ListAlternativeProfilesForSlot devuelve una lista de todos los perfiles alternativos "disponibles"
// (con ítem de stock existente) para una ranura específica de un sistema y un color dados.
func (s *ProfileSelectorService) ListAlternativeProfilesForSlot(
	ctx context.Context,
	input ListAlternativeProfilesForSlotInput,
) (*ListAlternativeProfilesForSlotOutput, error) {
	log := s.logger.WithFields(logrus.Fields{
		"method":         "ListAlternativeProfilesForSlot",
		"systemID":       input.SystemID,
		"ranuraRole":     input.RanuraRole,
		"desiredColorID": input.DesiredColorID,
	})
	log.Info("Listando perfiles alternativos para la ranura")

	output := &ListAlternativeProfilesForSlotOutput{
		AlternativeProfiles: []SelectedProfileInfo{},
	}

	if input.RanuraRole == "" {
		log.Error("RanuraRole (element_part_role) no puede estar vacío")
		return nil, fmt.Errorf("RanuraRole (element_part_role) es requerido para listar alternativas")
	}

	profileListItems, err := s.repo.GetSystemProfileListItems(ctx, input.SystemID, &input.RanuraRole)
	if err != nil {
		log.WithError(err).Errorf("Error obteniendo perfiles para rol '%s' en sistema ID %d", input.RanuraRole, input.SystemID)
		return nil, fmt.Errorf("error obteniendo perfiles candidatos para rol '%s': %w", input.RanuraRole, err)
	}

	if len(profileListItems) == 0 {
		log.Warnf("No se encontraron perfiles candidatos para el rol '%s' en el sistema ID %d", input.RanuraRole, input.SystemID)
		return output, nil
	}
	log.Infof("Se encontraron %d perfiles candidatos para el rol '%s'", len(profileListItems), input.RanuraRole)

	for _, item := range profileListItems {
		logRanuraItem := log.WithField("candidateProfileID", item.ProfileID)

		profileDetails, err := s.repo.GetProfileByID(ctx, item.ProfileID)
		if err != nil {
			logRanuraItem.WithError(err).Warnf("Error obteniendo detalles del perfil ID %d. Omitiendo alternativa.", item.ProfileID)
			continue
		}
		if profileDetails == nil {
			logRanuraItem.Warnf("No se encontraron detalles para el perfil ID %d. Omitiendo alternativa.", item.ProfileID)
			continue
		}

		stockItem, err := s.repo.GetStockItem(ctx, profileDetails.ProfileID, input.DesiredColorID)
		if err != nil {
			logRanuraItem.WithError(err).Warnf("Error verificando ítem de stock para perfil ID %d, color ID %d. Omitiendo alternativa.", profileDetails.ProfileID, input.DesiredColorID)
			continue
		}

		if stockItem != nil {
			logRanuraItem.Infof("Perfil ID %d (SKU: %s) es una alternativa válida (ítem de stock encontrado) en color ID %d.",
				profileDetails.ProfileID, profileDetails.ProfileSKU, input.DesiredColorID)

			output.AlternativeProfiles = append(output.AlternativeProfiles, SelectedProfileInfo{
				Profile:   *profileDetails,
				StockItem: *stockItem,
			})
		} else {
			logRanuraItem.Debugf("No se encontró ítem de stock para perfil ID %d (SKU: %s) en color ID %d. Omitiendo como alternativa.",
				profileDetails.ProfileID, profileDetails.ProfileSKU, input.DesiredColorID)
		}
	}

	log.Infof("Se encontraron %d perfiles alternativos (con ítem de stock existente) para la ranura '%s' y color ID %d.", len(output.AlternativeProfiles), input.RanuraRole, input.DesiredColorID)
	return output, nil
}
