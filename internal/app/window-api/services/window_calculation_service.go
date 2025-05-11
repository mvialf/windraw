package services

import (
	"context"
	"fmt"

	// "math" // Podría ser necesario para redondeos o funciones matemáticas avanzadas en el futuro

	"github.com/mvialf/windraw/internal/app/window-api/models"
	"github.com/mvialf/windraw/internal/pkg/constants" // Necesitaremos las constantes
	"github.com/sirupsen/logrus"
)

// WindowCalculationService se encarga de calcular las dimensiones de la hoja/marco terminados.
// Las dimensiones de corte exactas (con descuentos de soldadura/mecanizado) se manejarán
// en una etapa posterior o en un servicio de "despiece".
type WindowCalculationService struct {
	logger *logrus.Entry
}

// NewWindowCalculationService crea una nueva instancia de WindowCalculationService.
func NewWindowCalculationService(logger *logrus.Logger) *WindowCalculationService {
	return &WindowCalculationService{
		logger: logger.WithField("service", "window_calculation"),
	}
}

// CalculateDetailsInput agrupa todos los datos necesarios para los cálculos.
// Incluye el elemento (que se modificará), el sistema de perfiles y
// el mapa de perfiles seleccionados para cada rol/ranura.
type CalculateDetailsInput struct {
	Element            *models.Element // El elemento cuyos detalles se calcularán y rellenarán
	System             models.ProfileSystem
	ProfilesForRanuras map[string]SelectedProfileInfo // Resultado del ProfileSelectorService
}

// CalculateDetails calcula las dimensiones terminadas del marco y las hojas de un elemento.
// Modifica el input.Element directamente, rellenando los campos .Dimension de los Details.
// Los ángulos se asumen de 90 grados por ahora.
func (s *WindowCalculationService) CalculateDetails(ctx context.Context, input CalculateDetailsInput) error {
	log := s.logger.WithFields(logrus.Fields{
		"method":    "CalculateDetails",
		"elementID": input.Element.ID,
		"systemID":  input.System.SystemID,
	})
	log.Info("Iniciando cálculos de dimensiones terminadas para el elemento")

	// Validaciones básicas del input
	if input.Element == nil {
		return fmt.Errorf("el input Element no puede ser nulo")
	}
	if input.Element.Height <= 0 || input.Element.Width <= 0 {
		return fmt.Errorf("las dimensiones del Element (Height: %.2f, Width: %.2f) deben ser positivas", input.Element.Height, input.Element.Width)
	}
	if input.ProfilesForRanuras == nil || len(input.ProfilesForRanuras) == 0 {
		return fmt.Errorf("ProfilesForRanuras no puede estar vacío; se necesitan perfiles seleccionados")
	}

	// Asegurar que los maps de Details estén inicializados
	if input.Element.Frame.Details == nil {
		input.Element.Frame.Details = make(map[string]models.FrameDetail)
	}
	for i := range input.Element.Winds {
		if input.Element.Winds[i].Details == nil {
			input.Element.Winds[i].Details = make(map[string]models.WindDetail)
		}
		// Limpiar errores previos del elemento si este proceso es una recalculación
		// input.Element.Winds[i].Errors = []string{} // Opcional, depende de cómo se manejen errores acumulados
	}
	// input.Element.Errors = []string{} // Opcional

	// --- Cálculo para el Marco (Frame) ---
	// err := s.calculateFrameTerminatedDimensions(ctx, input) // TODO: Implementar
	// if err != nil {
	//     log.WithError(err).Error("Error calculando dimensiones del marco")
	//     return fmt.Errorf("error en dimensiones del marco: %w", err)
	// }
	log.Info("Cálculo de dimensiones del marco (TODO)")

	// --- Cálculo para las Hojas (Winds) ---
	// La lógica de qué tipo de hoja calcular vendrá del tipo de sistema
	if input.System.Type == constants.TYPE_SLIDING {
		err := s.calculateSlidingWindsTerminatedDimensions(ctx, input)
		if err != nil {
			log.WithError(err).Error("Error calculando dimensiones de hojas correderas")
			// Aquí podríamos decidir si añadir a input.Element.Errors o devolver el error directamente.
			// Por ahora, devolvemos error para detener el proceso si las hojas no se pueden calcular.
			return fmt.Errorf("error en dimensiones de hojas correderas: %w", err)
		}
	} else if input.System.Type == constants.TYPE_CASEMENT {
		log.Warnf("Cálculo para hojas tipo '%s' aún no implementado", input.System.Type)
		// err := s.calculateCasementWindsTerminatedDimensions(ctx, input)
		// if err != nil { return err }
	} else {
		log.Warnf("Tipo de sistema '%s' no reconocido para cálculo de hojas o lógica no implementada", input.System.Type)
	}

	log.Info("Cálculos de dimensiones terminadas para el elemento completados (parcialmente)")
	return nil
}

// calculateSlidingWindsTerminatedDimensions calcula el alto y ancho de las hojas correderas terminadas
// y asigna estas dimensiones y SKUs a los WindDetail correspondientes.
func (s *WindowCalculationService) calculateSlidingWindsTerminatedDimensions(ctx context.Context, input CalculateDetailsInput) error {
	log := s.logger.WithFields(logrus.Fields{
		"method":    "calculateSlidingWindsTerminatedDimensions",
		"elementID": input.Element.ID,
	})
	log.Info("Calculando dimensiones terminadas para hojas correderas")

	// --- 1. CÁLCULO DEL ALTO DE HOJA TERMINADA (calculatedWindHeight) ---
	var frameTopProfileH2 float64 = 0.0
	if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_FRAME_PERIMETER_TOP_SLIDING]; ok {
		if pInfo.Profile.ProfileH2 != nil {
			frameTopProfileH2 = *pInfo.Profile.ProfileH2
		} else {
			log.Warnf("Perfil para rol '%s' (SKU: %s) tiene ProfileH2 nulo. Asumiendo H2 (contribución a altura marco) = 0.", constants.ROLE_FRAME_PERIMETER_TOP_SLIDING, pInfo.Profile.ProfileSKU)
		}
	} else {
		return fmt.Errorf("perfil para rol '%s' no encontrado en ProfilesForRanuras, es necesario para calcular el alto de hoja", constants.ROLE_FRAME_PERIMETER_TOP_SLIDING)
	}

	var frameBottomProfileH2 float64 = 0.0
	if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_FRAME_PERIMETER_BOTTOM_SLIDING]; ok {
		if pInfo.Profile.ProfileH2 != nil {
			frameBottomProfileH2 = *pInfo.Profile.ProfileH2
		} else {
			log.Warnf("Perfil para rol '%s' (SKU: %s) tiene ProfileH2 nulo. Asumiendo H2 (contribución a altura marco) = 0.", constants.ROLE_FRAME_PERIMETER_BOTTOM_SLIDING, pInfo.Profile.ProfileSKU)
		}
	} else {
		return fmt.Errorf("perfil para rol '%s' no encontrado en ProfilesForRanuras, es necesario para calcular el alto de hoja", constants.ROLE_FRAME_PERIMETER_BOTTOM_SLIDING)
	}

	systemTopTraslape := input.System.TopOverlapMM
	systemBottomTraslape := input.System.BottomOverlapMM

	calculatedWindHeight := input.Element.Height - frameTopProfileH2 - systemTopTraslape - systemBottomTraslape - frameBottomProfileH2
	if calculatedWindHeight <= 0 {
		errorMsg := fmt.Sprintf("alto de hoja calculado (%.2f) es inválido. EH:%.2f, FTH2:%.2f, STT:%.2f, SBT:%.2f, FBH2:%.2f",
			calculatedWindHeight, input.Element.Height, frameTopProfileH2, systemTopTraslape, systemBottomTraslape, frameBottomProfileH2)
		log.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}
	log.Infof("CalculatedWindHeight (alto de hoja terminada): %.2f", calculatedWindHeight)

	// --- 2. CÁLCULO DEL ANCHO DE VANO INTERNO DEL MARCO (calculatedAnchoLuzTotal) ---
	var frameSideLeftEffectiveWidth float64 = 0.0
	if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING]; ok { // Asume un rol para ambos lados o el izquierdo
		if pInfo.Profile.ProfileH2 != nil { // ProfileH2 del perfil lateral es su "ancho efectivo"
			frameSideLeftEffectiveWidth = *pInfo.Profile.ProfileH2
		} else {
			log.Warnf("Perfil para rol '%s' (SKU: %s) tiene ProfileH2 nulo. Asumiendo ancho efectivo = 0 para marco lateral.", constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING, pInfo.Profile.ProfileSKU)
		}
	} else {
		return fmt.Errorf("perfil para rol '%s' no encontrado en ProfilesForRanuras, necesario para ancho de vano", constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING)
	}
	// Asumimos que el marco es simétrico y usa el mismo perfil lateral (y H2) en ambos lados.
	// Si el derecho fuera diferente, se obtendría con su rol específico de input.ProfilesForRanuras.
	frameSideRightEffectiveWidth := frameSideLeftEffectiveWidth

	calculatedAnchoLuzTotal := input.Element.Width - frameSideLeftEffectiveWidth - frameSideRightEffectiveWidth
	if calculatedAnchoLuzTotal <= 0 {
		errorMsg := fmt.Sprintf("AnchoLuzTotal calculado (%.2f) es inválido. EW:%.2f, FSLEW(H2):%.2f, FSREW(H2):%.2f",
			calculatedAnchoLuzTotal, input.Element.Width, frameSideLeftEffectiveWidth, frameSideRightEffectiveWidth)
		log.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}
	log.Infof("CalculatedAnchoLuzTotal (ancho de vano interno del marco): %.2f", calculatedAnchoLuzTotal)

	// --- 3. CÁLCULO DEL ANCHO DE HOJA TERMINADA (calculatedWindWidth) ---
	numActiveWinds := 0
	for _, w := range input.Element.Winds {
		if w.Kind == constants.KIND_SLIDING_WIND { // Usa la constante correcta para hoja corredera móvil
			numActiveWinds++
		}
	}

	if numActiveWinds == 0 {
		log.Info("No hay hojas correderas activas para calcular ancho de hoja.")
		return nil
	}
	log.Infof("Número de hojas correderas activas (numActiveWinds): %d", numActiveWinds)

	var meetingStyleTraslape float64 = 0.0
	if numActiveWinds > 1 {
		if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_JAMB_MEETING_SLIDING]; ok {
			if pInfo.Profile.ProfileH != nil { // ProfileH del perfil de encuentro es su "ancho visible frontal"
				meetingStyleTraslape = *pInfo.Profile.ProfileH
			} else {
				log.Warnf("Perfil para rol '%s' (SKU: %s) tiene ProfileH nulo. Asumiendo meetingStyleTraslape = 0.", constants.ROLE_WIND_JAMB_MEETING_SLIDING, pInfo.Profile.ProfileSKU)
			}
		} else {
			log.Warnf("Se esperaban %d hojas activas pero no se encontró perfil para rol '%s'. Asumiendo meetingStyleTraslape = 0.", numActiveWinds, constants.ROLE_WIND_JAMB_MEETING_SLIDING)
		}
	}
	log.Infof("MeetingStyleTraslape (basado en ProfileH de '%s'): %.2f", constants.ROLE_WIND_JAMB_MEETING_SLIDING, meetingStyleTraslape)

	systemSideTraslape := input.System.SideOverlapMM // Cuánto se mete la hoja en el marco lateral

	// Fórmula para N hojas iguales: (AnchoLuzTotal + (N-1)*TraslapeEncuentro + 2*TraslapeLateralHojaMarco) / N
	calculatedWindWidth := (calculatedAnchoLuzTotal + (float64(numActiveWinds-1) * meetingStyleTraslape) + (2 * systemSideTraslape)) / float64(numActiveWinds)

	if calculatedWindWidth <= 0 {
		errorMsg := fmt.Sprintf("ancho de hoja calculado (%.2f) es inválido. ALT:%.2f, N:%d, MST:%.2f, SST:%.2f",
			calculatedWindWidth, calculatedAnchoLuzTotal, numActiveWinds, meetingStyleTraslape, systemSideTraslape)
		log.Error(errorMsg)
		return fmt.Errorf(errorMsg)
	}
	log.Infof("CalculatedWindWidth (ancho de hoja terminada): %.2f", calculatedWindWidth)

	// --- 4. ASIGNAR DIMENSIONES Y SKUS A LOS WindDetail ---
	// Esto rellena los .Dimension con las medidas de la HOJA TERMINADA.
	// Los descuentos de corte/soldadura son una etapa posterior.
	// Los ángulos son 90 por ahora.

	for i := range input.Element.Winds {
		wind := &input.Element.Winds[i]
		if wind.Kind != constants.KIND_SLIDING_WIND {
			continue
		}
		logWind := log.WithFields(logrus.Fields{"windIndex": i, "windID": wind.ID})

		// Asignar SKUs y Dimensiones a los perfiles verticales de la hoja
		// La Dimensión de los perfiles verticales es calculatedWindHeight
		// La lógica para determinar qué rol (JAMB_SIDE vs JAMB_MEETING) va a LEFT vs RIGHT
		// de cada hoja puede ser compleja y depende de la configuración (ej. 2 hojas, 3 hojas, etc.)
		// Para V1.0, una simplificación:
		// - Si es la primera hoja (i=0) o la última (i=numActiveWinds-1) en un extremo, usa JAMB_SIDE en ese extremo.
		// - Los encuentros usan JAMB_MEETING.
		// ESTA LÓGICA DE ASIGNACIÓN DE ROLES A POSICIONES ES UNA SIMPLIFICACIÓN Y NECESITARÁ REFINARSE.

		// Detalle Izquierdo de la Hoja
		var leftProfileSKU string
		if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_JAMB_SIDE_SLIDING]; ok { // Default a jamba lateral
			leftProfileSKU = pInfo.Profile.ProfileSKU
		}
		if numActiveWinds > 1 && i > 0 { // Si no es la primera hoja, su lado izquierdo podría ser un encuentro
			if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_JAMB_MEETING_SLIDING]; ok {
				leftProfileSKU = pInfo.Profile.ProfileSKU
			}
		}
		if leftProfileSKU != "" {
			detail := wind.Details[constants.POSITION_LEFT] // Asume que el mapa Details ya tiene estas claves
			detail.ProfileSKU = leftProfileSKU
			detail.Dimension = calculatedWindHeight
			detail.AngleLeft = 90.0
			detail.AngleRight = 90.0
			wind.Details[constants.POSITION_LEFT] = detail
			logWind.Infof("  Detalle %s: SKU %s, Dim(H) %.2f", constants.POSITION_LEFT, detail.ProfileSKU, detail.Dimension)
		} else {
			logWind.Warnf("  No se pudo determinar el perfil para %s de la hoja %d", constants.POSITION_LEFT, i)
		}

		// Detalle Derecho de la Hoja
		var rightProfileSKU string
		if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_JAMB_SIDE_SLIDING]; ok { // Default a jamba lateral
			rightProfileSKU = pInfo.Profile.ProfileSKU
		}
		if numActiveWinds > 1 && i < numActiveWinds-1 { // Si no es la última hoja, su lado derecho podría ser un encuentro
			if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_JAMB_MEETING_SLIDING]; ok {
				rightProfileSKU = pInfo.Profile.ProfileSKU
			}
		}
		if rightProfileSKU != "" {
			detail := wind.Details[constants.POSITION_RIGHT]
			detail.ProfileSKU = rightProfileSKU
			detail.Dimension = calculatedWindHeight
			detail.AngleLeft = 90.0
			detail.AngleRight = 90.0
			wind.Details[constants.POSITION_RIGHT] = detail
			logWind.Infof("  Detalle %s: SKU %s, Dim(H) %.2f", constants.POSITION_RIGHT, detail.ProfileSKU, detail.Dimension)
		} else {
			logWind.Warnf("  No se pudo determinar el perfil para %s de la hoja %d", constants.POSITION_RIGHT, i)
		}

		// Perfiles Horizontales de la Hoja (Rieles)
		// Su Dimensión es calculatedWindWidth
		if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_RAIL_TOP_SLIDING]; ok {
			detail := wind.Details[constants.POSITION_TOP]
			detail.ProfileSKU = pInfo.Profile.ProfileSKU
			detail.Dimension = calculatedWindWidth
			detail.AngleLeft = 90.0
			detail.AngleRight = 90.0
			wind.Details[constants.POSITION_TOP] = detail
			logWind.Infof("  Detalle %s: SKU %s, Dim(W) %.2f", constants.POSITION_TOP, detail.ProfileSKU, detail.Dimension)
		} else {
			logWind.Warnf("  No se encontró perfil para %s para la hoja %d", constants.ROLE_WIND_RAIL_TOP_SLIDING, i)
			// Devolver error o añadir a element.Errors
		}

		if pInfo, ok := input.ProfilesForRanuras[constants.ROLE_WIND_RAIL_BOTTOM_SLIDING]; ok {
			detail := wind.Details[constants.POSITION_BOTTOM]
			detail.ProfileSKU = pInfo.Profile.ProfileSKU
			detail.Dimension = calculatedWindWidth
			detail.AngleLeft = 90.0
			detail.AngleRight = 90.0
			wind.Details[constants.POSITION_BOTTOM] = detail
			logWind.Infof("  Detalle %s: SKU %s, Dim(W) %.2f", constants.POSITION_BOTTOM, detail.ProfileSKU, detail.Dimension)
		} else {
			logWind.Warnf("  No se encontró perfil para %s para la hoja %d", constants.ROLE_WIND_RAIL_BOTTOM_SLIDING, i)
			// Devolver error o añadir a element.Errors
		}

		// Lógica de Traslapo Adicional (si el perfil de encuentro lo usa)
		if pInfoMeeting, okMeeting := input.ProfilesForRanuras[constants.ROLE_WIND_JAMB_MEETING_SLIDING]; okMeeting {
			if pInfoMeeting.Profile.UsesOverlap != nil && *pInfoMeeting.Profile.UsesOverlap {
				if pInfoOverlap, okOverlap := input.ProfilesForRanuras[constants.ROLE_WIND_VERTICAL_OVERLAP_SLIDING]; okOverlap {
					// Este es el perfil de traslapo adicional seleccionado por ProfileSelectorService.
					// Su dimensión vertical también sería calculatedWindHeight.
					logWind.Infof("  Hoja requiere traslapo adicional. Perfil de traslapo SKU: %s. Dim(H) %.2f",
						pInfoOverlap.Profile.ProfileSKU, calculatedWindHeight)
					// TODO: Decidir cómo se modela/almacena este perfil de traslapo adicional en la hoja.
					// ¿Es un "ExtraDetail" o una propiedad de la hoja?
					// Por ahora, está disponible en input.ProfilesForRanuras si fue seleccionado.
					// Aquí podríamos añadirlo a una lista de "piezas adicionales" de la hoja si el modelo lo soporta.
					// Ejemplo: wind.AdditionalParts = append(wind.AdditionalParts, models.WindDetail{ProfileSKU: pInfoOverlap.Profile.ProfileSKU, Dimension: calculatedWindHeight, ...})
				} else {
					// El ProfileSelectorService ya debería haber registrado un error si el traslapo era requerido pero no se encontró.
					// Aquí podemos añadir un error al elemento si es crítico.
					errMsg := fmt.Sprintf("Hoja %d: Perfil de encuentro SKU %s requiere traslapo, pero no se seleccionó perfil para rol '%s'", i, pInfoMeeting.Profile.ProfileSKU, constants.ROLE_WIND_VERTICAL_OVERLAP_SLIDING)
					logWind.Error(errMsg)
					// input.Element.Errors = append(input.Element.Errors, errMsg) // Si Element tiene un slice de Errors
				}
			}
		}
	}
	return nil
}

// TODO: Implementar calculateFrameTerminatedDimensions
// func (s *WindowCalculationService) calculateFrameTerminatedDimensions(ctx context.Context, input CalculateDetailsInput) error { ... }
