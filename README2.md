
# interacción con la base de datos
- el usuario decide que tipo de elemento quiere crear (sliding o casement) tabla 'profile_system' 'type
- se selecciona el 'profile_system' 'profile_system_id' que se usará, si hay mas de una aternativa se seleccionará por importancia 'profile_system' 'primacy'(la primacy es un numero del que parte desde el 1 y ordena la prioridad de los sistema de perfiles)
- una vez seleccionado el 'profile_system' 'profile_system_id', se usan los perfiles que se usan para el sistema seleccionado, para ello se usa la tabla 'system_profile_list' en la que solo se podran usar los perfiles que el campo system_id sea igual al 'profile_system' 'profile_system_id'
- si existe un perfil que se repite su 'profile' 'profile_type' en la tabla 'system_profile_list' seleccionada se usara el perfil que tenga la mayor importancia 'profile' 'primacy'
- se debe seleccionar el color del 'profile_system' en la tabla 'system_available_colors'
- al seleccionar el color se hará la consulta a la tabla 'stock_items' donde se relacionara el campo 'profile_id y 'color_id'.
- en esta tabla se encuentra el 'item_sku', profile_price: precio del perfil, profile_length: dimension del perfil en mm, stock: stock disponible. el 'profileColorName' es igual a la concatenación de 'profile_name' + ' ' + 'color_name'

