
### 3.2. Capas Lógicas

1.  **Capa de Modelos/Dominio (`internal/app/window-api/models/`)**:
    *   Contiene las definiciones de las estructuras de datos (`structs`) que representan las entidades centrales del negocio: `Project`, `Element`, `Frame`, `Wind`, `FrameDetail`, `WindDetail`, `Contact`, `Component`, `Module`, etc.
    *   Incluye funciones constructoras (`New*`) para estas entidades, asegurando la validación inicial de datos.
    *   Puede contener métodos asociados a estas estructuras para operaciones intrínsecas (ej. `AddComponent` a un `Project`).
    *   La lógica de cálculo más compleja (como `CalculateFrameDetails`) reside temporalmente aquí como métodos, pero se moverá progresivamente a la capa de servicios.

2.  **Capa de Constantes (`internal/pkg/constants/`)**:
    *   Define todas las constantes globales usadas en la aplicación para evitar "magic strings/numbers" y facilitar el mantenimiento. Esto incluye tipos de materiales, tipos de elementos, tipos de corte, posiciones, etc.

3.  **Capa de Repositorios (`internal/app/window-api/repositories/`)**:
    *   Define las **interfaces** (puertos en terminología hexagonal) que describen las operaciones de persistencia de datos (CRUD) para las entidades del dominio (ej. `ProjectRepository`, `ElementRepository`).
    *   Contiene las **implementaciones** (adaptadores) de estas interfaces, específicamente para interactuar con **Supabase (PostgreSQL)**. Esta capa abstrae los detalles de la base de datos del resto de la aplicación.

4.  **Capa de Servicios (`internal/app/window-api/services/`)**:
    *   Contiene la **lógica de negocio principal** y la orquestación de operaciones complejas.
    *   Los servicios utilizan los repositorios (a través de sus interfaces) para acceder a los datos y manipulan las entidades del dominio (modelos).
    *   Ejemplos: `ProjectService` para gestionar la creación de proyectos y sus componentes; `WindowCalculationService` para realizar los cálculos detallados de despiece de perfiles, holguras, y aplicación de reglas del catálogo.
    *   Esta capa es independiente de la capa de presentación (HTTP).

5.  **Capa de Manejadores/API (`internal/app/window-api/handlers/`)**:
    *   Responsable de manejar las peticiones HTTP entrantes de la API REST.
    *   Valida los datos de entrada de la API, decodifica JSON, etc.
    *   Llama a los métodos apropiados de la capa de servicios para ejecutar la lógica de negocio.
    *   Formatea las respuestas (generalmente JSON) y gestiona los códigos de estado HTTP.
    *   Se utiliza la biblioteca estándar `net/http` de Go, potencialmente complementada con un router ligero como `chi` para un manejo de rutas más avanzado.

6.  **Punto de Entrada (`cmd/window-api/main.go`)**:
    *   Inicializa la aplicación: carga la configuración, establece la conexión a la base de datos (Supabase), instancia los repositorios, servicios y handlers.
    *   Configura el router HTTP y arranca el servidor.

### 3.3. Flujo de Datos Típico (Ejemplo: Crear un Elemento)

1.  El cliente (Frontend React) envía una petición HTTP POST a un endpoint como `/api/v1/projects/{projectID}/elements` con los datos del nuevo elemento en formato JSON.
2.  El `ElementHandler` en la capa de `handlers` recibe la petición.
3.  El handler decodifica el JSON, valida los datos de la API.
4.  El handler llama a un método como `CreateElement(projectID, elementData)` en el `ElementService` (o `ProjectService`).
5.  El `ElementService`:
    *   Puede interactuar con `ProjectRepository` para verificar que el proyecto existe.
    *   Utiliza el constructor `models.NewElement()` para crear una instancia del elemento, aplicando validaciones de dominio.
    *   Si se proporcionan perfiles, podría invocar un `WindowCalculationService` para calcular los detalles del marco y las hojas (este servicio accedería a un catálogo de perfiles, que es una dependencia externa).
    *   Llama a `ElementRepository.Save(elemento)` para persistir el nuevo elemento en Supabase.
6.  El servicio devuelve el elemento creado (o un error) al handler.
7.  El handler formatea la respuesta (el elemento creado o un mensaje de error) como JSON y la envía de vuelta al cliente con el código HTTP apropiado (ej. 201 Created o 400 Bad Request).

## 4. Tecnologías Clave

*   **Lenguaje de Programación:** Go (Golang)
*   **Base de Datos:** Supabase (PostgreSQL)
    *   Se accede principalmente mediante conexión directa a PostgreSQL desde el backend Go utilizando el driver `database/sql` y `pgx`.
    *   Supabase Auth puede ser utilizado para la gestión de usuarios.
*   **API:** RESTful sobre HTTP, utilizando JSON como formato de intercambio de datos.
*   **Enrutamiento HTTP:** Inicialmente con la biblioteca estándar `net/http`, con posibilidad de integrar un router ligero como `chi` si es necesario.
*   **Gestión de Dependencias:** Go Modules.
*   **Frontend (Previsto):** React.

## 5. Funcionamiento Básico (Conceptos del Dominio)

El sistema se basa en una jerarquía de entidades:

*   **`Project`**: La entidad de más alto nivel. Contiene información del cliente, costos generales y una lista de `Component`s.
    *   Atributos: `ID`, `Name`, `CreatedAt`, `Contact`, `Costs`, `Components`, `IvaRate`.
*   **`Contact`**: Información del cliente (persona o empresa).
*   **`ProjectCost`**: Costos adicionales del proyecto (fijos o porcentuales).
*   **`Component`**: Una agrupación lógica de `Module`s dentro de un proyecto (ej. "Ventanas Planta Baja", "Puertas Terraza").
    *   Atributos: `ID`, `Modules`.
*   **`Module`**: Una agrupación de `Element`s que pueden estar unidos por perfiles auxiliares (ej. una ventana fija unida a una puerta corredera).
    *   Atributos: `ID`, `Elements`, `AuxProfile`.
*   **`Element`**: La unidad funcional principal, una ventana o puerta individual.
    *   Atributos: `ID`, `Width`, `Height`, `Material`, `Type` (tipología), `Structure` (ventana/puerta), `Frame`, `Winds` (lista de hojas), `Properties`.
*   **`Frame`**: El marco perimetral del `Element`.
    *   Atributos: `Name`, `Geometry`, `Width`, `Height`, `CutType`, `Details` (mapa de `FrameDetail`).
*   **`Wind` (Hoja)**: Un panel móvil o fijo dentro de un `Element`.
    *   Atributos: `ID`, `Name`, `Kind` (tipo de hoja), `Status`, `OpeningSide`, `OpeningDirection`, `Width`, `Height`, `CutType`, `Details` (mapa de `WindDetail`).
*   **`FrameDetail` / `WindDetail`**: Describe una pieza individual de perfil para un marco o una hoja.
    *   Atributos: `Position` (Top, Bottom, Left, Right), `ProfileSKU`, `Color`, `Dimension` (longitud de corte), `AngleLeft`, `AngleRight`, `ReinforcedUsed`, `ReinforcedSKU`.

### 5.1. Proceso de Cálculo (Simplificado)

1.  Se crea un `Element` con sus dimensiones generales, material, tipo y estructura.
2.  Se crea un `Frame` para el elemento, con su geometría y tipo de corte.
3.  Se asignan perfiles (mediante `ProfileSKU` y `Color`) a cada `Position` del `Frame` usando `SetFrameProfile()`.
4.  Se invoca `CalculateFrameDetails()`:
    *   Esta función (eventualmente un servicio) consultaría un catálogo de perfiles (no implementado aún).
    *   Basándose en el `ProfileSKU`, las dimensiones del marco, el tipo de corte y las reglas del catálogo, calcula:
        *   `Dimension` (longitud de corte) para cada `FrameDetail`.
        *   `AngleLeft` y `AngleRight`.
        *   Si se necesita `ReinforcedUsed` y el `ReinforcedSKU`.
5.  Se crean `Wind`s (hojas) para el elemento, con sus dimensiones, tipo de hoja y tipo de corte.
6.  Se asignan perfiles a cada `Position` de cada `Wind` usando `SetWindProfile()`.
7.  Se invoca `CalculateWindDetails()` para cada hoja, siguiendo una lógica similar a la del marco pero con reglas específicas para hojas y sus tipos de corte (ej. solapes).
8.  Las hojas se añaden al `Element`.

## 6. Constantes Importantes

El sistema utiliza un conjunto de constantes definidas en `internal/pkg/constants/` para estandarizar valores y mejorar la legibilidad:

*   `MATERIAL_*` (ej. `MATERIAL_PVC`, `MATERIAL_ALUMINIO`)
*   `TYPE_*` (ej. `TYPE_SLIDING`, `TYPE_CASEMENT`)
*   `PROFILE_TYPE_*` (ej. `PROFILE_TYPE_SLIDING_FRAME`)
*   `STRUCTURE_*` (ej. `STRUCTURE_VENTANA`)
*   `GEOMETRY_*` (ej. `GEOMETRY_RECTANGULAR`)
*   `POSITION_*` (ej. `POSITION_LEFT`, `POSITION_TOP`)
*   `CUT_*` (ej. `CUT_SQUARE`, `CUT_ANGLE_WIND`)
*   `WIND_KIND_*` (ej. `WIND_KIND_SLIDING_MOVIL`)
*   `WIND_STATUS_*` (ej. `WIND_STATUS_ACTIVE`)
*   `OPENING_SIDE_*` / `OPENING_DIRECTION_*`

## 7. Cómo Contribuir / Próximos Pasos

(Esta sección es más para un proyecto real, pero es bueno pensar en ella)

*   **Implementación de la Capa de Repositorios:** Conectar con Supabase y escribir las implementaciones para las interfaces de repositorio.
*   **Desarrollo de la Capa de Servicios:** Mover la lógica de cálculo compleja y la orquestación a los servicios.
*   **Implementación del Catálogo de Perfiles:** Diseñar cómo se almacenará y accederá la información detallada de los perfiles (probablemente en Supabase).
*   **Desarrollo Completo de los Handlers de API:** Definir todos los endpoints REST necesarios.
*   **Autenticación y Autorización:** Integrar con Supabase Auth.
*   **Pruebas:** Escribir pruebas unitarias y de integración exhaustivas.
*   **Documentación de API:** Generar documentación OpenAPI/Swagger.

## 8. Configuración e Instalación

(Detalles a añadir cuando el proyecto sea desplegable)

1.  Clonar el repositorio.
2.  Configurar variables de entorno (ej. credenciales de Supabase).
3.  Ejecutar `go mod tidy` para descargar dependencias.
4.  Construir el binario: `go build -o windraw_api ./cmd/window-api/`
5.  Ejecutar la aplicación: `./windraw_api`

---

Este `README.md` es un punto de partida. Deberías actualizarlo continuamente a medida que tu proyecto evoluciona. Cuanto más detallado y preciso sea, mejor podrá entenderlo un modelo de IA (y otros desarrolladores).

Considera añadir diagramas de arquitectura o de flujo de datos si eso ayuda a clarificar conceptos complejos.