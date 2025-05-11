
# interacción con la base de datos
// el usuario crea un componente y se crea automaticamente un modulo y un elemento.
- el usuario decide que tipo de elemento quiere crear (sliding o casement) tabla 'profile_system' 'type
- se selecciona el 'profile_system' 'profile_system_id' que se usará, si hay mas de una aternativa se seleccionará por importancia 'profile_system' 'primacy'(la primacy es un numero del que parte desde el 1 y ordena la prioridad de los sistema de perfiles) # manejo de error. si existen 2 o mas 'profile_system' con la misma importancia se seleccionara el que tenga la menor 'profile_system_id'
- una vez seleccionado el 'profile_system' 'profile_system_id', se usan los perfiles que se usan para el sistema seleccionado, para ello se usa la tabla 'system_profile_list' en la que solo se podran usar los perfiles que el campo system_id sea igual al 'profile_system' 'profile_system_id'
- si existe un perfil que se repite su 'profile' 'profile_type' en la tabla 'system_profile_list' seleccionada se usara el perfil que tenga la mayor importancia 'profile' 'primacy'
- se debe seleccionar el color del 'profile_system' en la tabla 'system_available_colors'
- al seleccionar el color se hará la consulta a la tabla 'stock_items' donde se relacionara el campo 'profile_id y 'color_id'.
- en esta tabla se encuentra el 'item_sku', profile_price: precio del perfil, profile_length: dimension del perfil en mm, stock: stock disponible. el 'profileColorName' es igual a la concatenación de 'profile_name' + ' ' + 'color_name'

# calculo hojas
- topOverlap = profile_systems.top_overlap_mm
- bottomOverlap = profile_systems.bottom_overlap_mm
- sideOverlap = profile_systems.side_overlap_mm
- slidingFrameTopH = profiles.profile_h
- slidingFrameBottomH = profiles.profile_h


## alto hoja
alto hoja = alto ventana - slidingFrameTopH - topOverlap - bottomOverlap - slidingFrameBottomH

## reglas
- Las 'ventanas correderas' deben tener mínimo 2 hojas
- Wind 1 solo puede ser fija o abrir a la derecha
- Wind última solo puede ser fija o abrir a la izquierda 
- Debe existirá al menos un wind con apertura y activa
- Si wind i and wind i+1 estan en el mismo truk, wind i  solo puede ser fija o apertura  (por defecto) izquierda. Wind i+1 solo puede ser fija o abrir hacia la derecha (por defecto)
- Si wind i es = wind1, wind i+1 = a wind2. Wind1 solo puede ser fija, wind 2 puede ser fija o abrir a la derecha (por defecto)
- Si la ventana corredera tiene 3 trucks, win i y wind i+1 solo pueden estar a en el mismo truk o 1 trucks de distancia (por defecto) . // Ejemplo si wind i está en el truck 1, wind i+1 puede estar en el truck 1 o 2. Si wind i = truck 2, wind i+1 puede estar en truck 1 2 o 3
- Si tiene 2 trucks and wind i abre hacia la derecha wind i+1 puede ser fija o abrir a la izquierda 
- Si tiene 3 trucks wind i-1 está en truck 1, wind i+1 truck 3. Wind puede ser fija o inactiva ( defecto) 
- Si tiene 3 trucks wind i-1 está en truck 3, wind i+1 truck 1. Wind puede ser fija o inactiva ( defecto)

## ancho hoja
Para calcular las medidas de las hojas de una ventana corredera necesitamos respetar las reglas del armado, se me ocurre la siguiente forma, tu quizás lo puedas mejorar.

1. El ancho de las hojas dependerá del número de hojas que la ventana tenga y en que pista se ubiquen. Pero lo primero que debemos calcular es el ancho luz total
Fórmula General para el Ancho de una Hoja (simplificada para empezar):

Una fórmula base para el ancho de una hoja podría ser algo como:

AnchoLuzTotal = element.Width - (slidingFrameSideLeftH2) - sideOverlap - sideOverlap(slidingFrameSideRightH2) 

para ver el numero de hojas se me ocurre hacer un array

meetingStyleTraslape = overlapStyleTraslape
overlapStyleTraslape = profile_h #de la hoja

### ejemplo
# EJEMPLO DE CÁLCULO PARA VENTANA CORREDERA

## 1. Dimensiones Generales del Elemento (Ventana Completa)
element_width: 2000  // Ancho total de la ventana terminada (mm)
element_height: 1500 // Alto total de la ventana terminada (mm)
number_of_active_sliding_winds: 2 // Número de hojas correderas móviles

## 2. Propiedades del Sistema de Perfiles (`models.ProfileSystem` seleccionado)
system_top_overlap_mm: 8    // `TopOverlapMM` (Traslape vertical superior de la hoja sobre el marco)
system_bottom_overlap_mm: 8 // `BottomOverlapMM` (Traslape vertical inferior de la hoja sobre el marco)
system_side_overlap_mm: 8    // `SideOverlapMM` (Cuánto se "mete" cada hoja en el marco lateral)

## 3. Propiedades de los Perfiles del Marco (`models.Profile` para los roles del marco)
# Perfil del Marco Superior (ej. rol `constants.ROLE_FRAME_PERIMETER_TOP_SLIDING`)
frame_top_profile_h2: 54    // `ProfileH2` del perfil del marco superior (su contribución a la altura del marco)

# Perfil del Marco Inferior (ej. rol `constants.ROLE_FRAME_PERIMETER_BOTTOM_SLIDING`)
frame_bottom_profile_h2: 54 // `ProfileH2` del perfil del marco inferior (su contribución a la altura del marco)

# Perfil del Marco Lateral Izquierdo (ej. rol `constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING`)
# (Usaremos ProfileH2 de este perfil como su "ancho efectivo" que reduce el vano)
frame_side_left_profile_h2: 54 // `ProfileH2` del perfil del marco lateral izquierdo

# Perfil del Marco Lateral Derecho (ej. rol `constants.ROLE_FRAME_PERIMETER_SIDE_SLIDING`)
# Si es el mismo perfil que el izquierdo, tendrá el mismo ProfileH2. Si es diferente, proporciona su ProfileH2.
frame_side_right_profile_h2: 54 // `ProfileH2` del perfil del marco lateral derecho

## 4. Propiedades del Perfil de Encuentro de Hoja (`models.Profile` para el rol de encuentro)
# (Usado para `meetingStyleTraslape` entre hojas)
# Perfil de Encuentro de Hoja (ej. rol `constants.ROLE_WIND_JAMB_MEETING_SLIDING`)
wind_meeting_jamb_profile_h: 80 // `ProfileH` de este perfil (su "ancho visible frontal" que contribuye al traslape horizontal)

