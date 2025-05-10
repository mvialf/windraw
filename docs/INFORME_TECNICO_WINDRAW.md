# Informe Técnico del Proyecto Go: Windraw

## 1. Introducción y Visión General

### 1.1. Propósito del Proyecto
    - Basado en el código inicial y los comentarios en `cmd/window-api/main.go`, el proyecto tiene como objetivo desarrollar un backend en Go para una "aplicación de fabricación de ventanas".
    - Actualmente, la funcionalidad principal implementada es la capacidad de conectarse e interactuar con una base de datos Supabase para consultar datos (ej., una tabla `projects`). La intención parece ser construir una API para gestionar la información relacionada con la fabricación de ventanas.

### 1.2. Funcionalidades Principales Identificadas
    - **Conexión a Supabase:** El sistema puede cargar credenciales de Supabase (`BaseURL`, `ServiceRoleKey`) desde variables de entorno.
    - **Cliente Supabase:** Implementa un cliente (`internal/pkg/apiclient/SupabaseClient`) capaz de realizar consultas GET genéricas a la API REST de Supabase (PostgREST).
    - **Carga de Configuración:** Utiliza `internal/pkg/config/LoadConfig` para gestionar la configuración de la aplicación a partir de variables de entorno, con soporte para archivos `.env` mediante `github.com/joho/godotenv`.
    - **Consulta de Datos:** El `main.go` actual demuestra cómo consultar datos de tablas en Supabase.
    - (La funcionalidad de API HTTP y la lógica de negocio específica, como el cálculo de perfiles, aún no están implementadas o no son visibles en los archivos analizados hasta ahora).

## 2. Entorno de Desarrollo y Herramientas (Inferido del Código)

### 2.1. Versión de Go
    - Versión de Go: `1.24.2` (inferido de `go.mod`).

### 2.2. Supabase en el Proyecto
    - El proyecto utiliza Supabase como backend de base de datos. La interacción se realiza a través de su API REST (PostgREST).
    - Se ha implementado un cliente Go (`internal/pkg/apiclient/SupabaseClient`) que maneja la construcción de peticiones HTTP GET, la autenticación mediante `apikey` y `Bearer token` (usando la `ServiceRoleKey`), y la decodificación de respuestas JSON.
    - La configuración para Supabase (`BaseURL` y `ServiceRoleKey`) se carga desde variables de entorno (`SUPABASE_API_URL`, `SUPABASE_SERVICE_KEY`) gestionadas por `internal/pkg/config/config.go`.
    - El archivo `cmd/window-api/main.go` contiene ejemplos de cómo usar este cliente para consultar tablas (ej., una tabla `projects`).

### 2.3. Este Editor: Windsurf AI
    - Este análisis y la generación de este informe se realizan con Windsurf AI, una herramienta de asistencia de IA para la comprensión y documentación de proyectos de software.

## 3. Arquitectura del Proyecto (Analizada desde el Código)

### 3.1. Estructura de Carpetas y Archivos
    - A continuación, se presenta una representación del árbol de directorios del proyecto (excluyendo `docs/`):
      ```
/windraw
|-- .env
|-- .gitignore
|-- README.md
|-- README2.md
|-- api/
|-- cmd/
|   `-- window-api/
|       `-- main.go
|-- configs/
|   `-- config.dev.yaml
|-- go.mod
|-- go.sum
|-- internal/
|   |-- app/
|   |   `-- window-api/
|   |       |-- handlers/
|   |       |   `-- http_handlers.go
|   |       |-- models/
|   |       |   |-- catalog_models.go
|   |       |   |-- component_models.go
|   |       |   |-- element_models.go
|   |       |   `-- project_models.go
|   |       |-- repositories/
|   |       |   |-- profile_catalog_repository.go
|   |       |   |-- project_repository.go
|   |       |   `-- supabase_profile_catalog_repository.go
|   |       |-- services/  # (Contenido de services no listado por find_by_name, se asumirá vacío o se explorará después)
|   `-- pkg/
|       |-- apiclient/
|       |   `-- supabase_client.go
|       |-- config/
|       |   `-- config.go
|       |-- constants/
|       |   `-- constants.go
|       |-- database/    # (Contenido de database no listado por find_by_name)
|       |-- projectfile/
|       |   `-- manager.go
|       `-- utils/       # (Contenido de utils no listado por find_by_name)
|-- main.exe
|-- pkg/                 # (Contenido de pkg no listado por find_by_name)
|-- scripts/             # (Contenido de scripts no listado por find_by_name)
|-- tests/
|   |-- integration/
|   `-- unit/
`-- web/
    `-- static/
      ```
    - **Descripción de Carpetas y Archivos Principales (hasta ahora analizados):**
        - `cmd/window-api/main.go`: Punto de entrada de la aplicación "window-api". Actualmente, inicializa la configuración, el cliente Supabase y, actualmente, ejecuta código de ejemplo para consultar Supabase. No implementa un servidor HTTP completo todavía, pero indica la intención de ser una API.
        - `internal/pkg/config/config.go`: Define y carga la configuración de la aplicación (ej., credenciales de Supabase) desde variables de entorno, con soporte para archivos `.env`.
        - `configs/config.dev.yaml`: Archivo de configuración YAML. Actualmente está vacío y el código de `config.go` no lo utiliza para cargar la configuración (se basa en variables de entorno).
        - `internal/pkg/apiclient/supabase_client.go`: Contiene el cliente (`SupabaseClient`) para interactuar con la API REST de Supabase. Maneja la autenticación y las solicitudes GET.
        - `internal/app/window-api/handlers/http_handlers.go`: Previsto para contener los manejadores de las rutas HTTP de la API. Actualmente está vacío.
        - `internal/app/window-api/models/`: Paquete que contiene las estructuras de datos (modelos) del proyecto.
            - `catalog_models.go`: Actualmente vacío. Probablemente destinado a modelos relacionados con catálogos de perfiles, vidrios, herrajes, etc.
            - `component_models.go`: Define `AuxProfile`, `Module` y `Component`, que estructuran un proyecto en partes jerárquicas.
            - `element_models.go`: Contiene los modelos detallados para `FrameDetail`, `WindDetail`, `Frame`, `Wind` y `Element`. Estos son cruciales para describir una ventana o puerta y sus partes, incluyendo dimensiones y ángulos para el despiece.
            - `project_models.go`: Define la estructura `Project` principal, así como `Contact`, `ProjectCost` y funciones helper como `generateID` y `IsValidOption`.
        - `internal/pkg/constants/constants.go`: Define un conjunto de constantes globales que representan el vocabulario del dominio del proyecto (materiales, tipos de perfiles, geometrías, etc.), usadas para validación y lógica de negocio.
    - (Se describirá el propósito de otras carpetas y archivos principales a medida que se analicen).

### 3.2. Patrones de Diseño y Principios Arquitectónicos Identificados
    - **Carga de Configuración Centralizada:** El paquete `internal/pkg/config` centraliza la lógica de carga de configuración.
    - **Cliente API Dedicado:** El paquete `internal/pkg/apiclient` encapsula la lógica de comunicación con Supabase, lo que es una buena práctica para separar preocupaciones.
    - **Estructura Orientada a Funcionalidades (parcial):** La organización en `internal/app/window-api` con subcarpetas como `handlers`, `models`, `repositories`, `services` sugiere una intención de seguir patrones como Clean Architecture o similar, aunque aún está en una etapa temprana de desarrollo.
    - (Se identificarán más patrones a medida que se explore el código).

## 4. Dependencias Clave (Librerías y Frameworks)

### 4.1. Análisis de `go.mod`
    - **Dependencias Directas:**
        - `github.com/joho/godotenv v1.5.1`: Utilizada en `internal/pkg/config/config.go` para cargar variables de entorno desde un archivo `.env` durante el desarrollo, facilitando la gestión de la configuración sin necesidad de establecer variables de sistema.
    - **Dependencias Indirectas:**
        - `github.com/patrickmn/go-cache v2.1.0+incompatible`
        - `github.com/sirupsen/logrus v1.9.3`
        - `golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8`
    - (Se analizará el uso de estas dependencias, especialmente las indirectas que puedan ser relevantes, a medida que se explore el código).

## 5. Código Fuente Relevante y Explicado
    - (Se completará progresivamente. Se excluirán `.gitignore` y `.env`. El archivo `configs/config.dev.yaml` se ha mencionado pero está vacío y no se usa actualmente para la carga de configuración).

    ### 5.1. Módulos/Paquetes Principales

        - **`cmd/window-api/`**
            - **Propósito e Importancia:** Contiene el ejecutable principal de la aplicación API.
            - **Archivos Clave:**
                - **`main.go`**
                    - Responsabilidad: Punto de entrada de la aplicación. Inicializa la configuración, el cliente Supabase y, actualmente, ejecuta código de ejemplo para consultar Supabase. No inicia un servidor HTTP por el momento.
                    - Fragmento de código representativo:
                    ```go
                    // en func main()
                    cfg, err := config.LoadConfig()
                    if err != nil {
                        log.Fatalf("Error al cargar la configuración: %v", err)
                    }
                    supaClient := apiclient.NewSupabaseClient(&cfg.SupabaseAPI)
                    log.Println("Cliente API de Supabase creado.")

                    // ... ejemplo de consulta ...
                    path := "/rest/v1/projects"
                    queryParams := "select=id,name,created_at&order=created_at.desc&limit=2"
                    var projects []ProjectFromAPI
                    err = supaClient.QueryData(path, queryParams, &projects)
                    // ... manejo de error y resultados ...
                    ```
                    - Explicación: Este fragmento muestra la inicialización de la configuración y del cliente Supabase, seguido de una consulta de ejemplo a una tabla `projects`.

        - **`internal/pkg/config/`**
            - **Propósito e Importancia:** Gestiona la carga y el acceso a la configuración de la aplicación.
            - **Archivos Clave:**
                - **`config.go`**
                    - Responsabilidad: Define las estructuras de configuración (`Config`, `APIConfig`) y proporciona la función `LoadConfig()` para cargar valores desde variables de entorno, con la ayuda de `godotenv`.
                    - Fragmento de código representativo:
                    ```go
                    // Config almacena toda la configuración de la aplicación.
                    type Config struct {
                        SupabaseAPI APIConfig
                    }

                    type APIConfig struct {
                        BaseURL        string
                        ServiceRoleKey string
                    }

                    func LoadConfig() (*Config, error) {
                        err := godotenv.Load() // Carga .env
                        // ... (manejo de error de godotenv)
                        cfg := &Config{}
                        cfg.SupabaseAPI.BaseURL = getEnv("SUPABASE_API_URL", "")
                        cfg.SupabaseAPI.ServiceRoleKey = getEnv("SUPABASE_SERVICE_KEY", "")
                        // ... (validaciones)
                        return cfg, nil
                    }
                    ```
                    - Explicación: Define cómo se estructura la configuración y cómo se carga desde el entorno.

        - **`internal/pkg/apiclient/`**
            - **Propósito e Importancia:** Encapsula la lógica para interactuar con la API de Supabase.
            - **Archivos Clave:**
                - **`supabase_client.go`**
                    - Responsabilidad: Define `SupabaseClient` y sus métodos para realizar peticiones a la API REST de Supabase, manejando la autenticación y la decodificación de respuestas.
                    - Fragmento de código representativo:
                    ```go
                    type SupabaseClient struct {
                        BaseURL        string
                        ServiceRoleKey string
                        HttpClient     *http.Client
                    }

                    func NewSupabaseClient(cfg *config.APIConfig) *SupabaseClient { /* ... */ }

                    func (c *SupabaseClient) QueryData(path string, queryParams string, target interface{}) error {
                        fullURL := c.BaseURL + path
                        if queryParams != "" {
                            fullURL += "?" + queryParams
                        }
                        req, err := http.NewRequest("GET", fullURL, nil)
                        // ... (manejo de error)
                        req.Header.Set("apikey", c.ServiceRoleKey)
                        req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
                        // ... (resto de la petición y decodificación)
                        return nil
                    }
                    ```
                    - Explicación: Muestra la estructura del cliente y la función clave para consultar datos, incluyendo cómo se establecen las cabeceras de autenticación.

        - **`internal/app/window-api/handlers/`**
            - **Propósito e Importancia:** Destinado a contener los manejadores de solicitudes HTTP para la API.
            - **Archivos Clave:**
                - **`http_handlers.go`**
                    - Responsabilidad: Se espera que implemente las funciones que manejan las rutas de la API.
                    - Estado Actual: El archivo está vacío. No hay manejadores definidos todavía.

        - **`internal/app/window-api/models/`**
            - **Propósito e Importancia:** Define las principales entidades de datos (structs) que representan el dominio de la aplicación de fabricación de ventanas. Estas estructuras son fundamentales para la lógica de negocio, el almacenamiento de datos y la comunicación dentro de la aplicación.
            - **Archivos Clave y Estructuras Relevantes:**
                - **`project_models.go`**:
                    - `Project`: Entidad raíz que agrupa todos los elementos de un pedido o trabajo. Contiene nombre, contacto, costos, componentes y tasa de IVA.
                    - `Contact`: Información del cliente (persona o empresa).
                    - `ProjectCost`: Costos adicionales (fijos o porcentuales) asociados al proyecto.
                    - `generateID()`: Función helper para crear identificadores únicos para los modelos.
                    - `IsValidOption()`: Función helper para validaciones.
                    ```go
                    // en project_models.go
                    type Project struct {
                        ID         string        `json:"id"`
                        Name       string        `json:"name"`
                        CreatedAt  time.Time     `json:"created_at"`
                        Contact    Contact       `json:"contact"`
                        Costs      []ProjectCost `json:"costs"`
                        Components []Component   `json:"components,omitempty"`
                        IvaRate    float64       `json:"iva_rate"`
                    }
                    ```
                - **`component_models.go`**:
                    - `Component`: Agrupación lógica de `Module`s (ej., "Ventanas Fachada Norte").
                    - `Module`: Agrupación de `Element`s que pueden estar unidos por un `AuxProfile`.
                    - `AuxProfile`: Perfil auxiliar para uniones entre módulos.
                    ```go
                    // en component_models.go
                    type Component struct {
                        ID      string   `json:"id"`
                        Modules []Module `json:"modules"`
                    }
                    type Module struct {
                        ID         string      `json:"id"`
                        Elements   []Element   `json:"elements"`
                        AuxProfile *AuxProfile `json:"aux_profile,omitempty"`
                    }
                    ```
                - **`element_models.go`**:
                    - `Element`: La unidad funcional principal (ventana, puerta). Contiene dimensiones generales, material, tipo, estructura, un `Frame` y una lista de `Wind`s.
                    - `Frame`: Marco perimetral del `Element`. Incluye dimensiones, geometría, tipo de corte y `FrameDetail`s.
                    - `Wind`: Hoja (panel móvil o fijo) dentro de un `Element`. Similar al marco, tiene tipo, dimensiones, tipo de corte y `WindDetail`s.
                    - `FrameDetail`: Pieza individual de perfil para un marco. Contiene `ProfileSKU`, `Dimension`, `AngleLeft`, `AngleRight`, esenciales para el despiece.
                    - `WindDetail`: Pieza individual de perfil para una hoja. Similar a `FrameDetail`.
                    - Constructores como `NewElement`, `NewFrame`, `NewWind` inicializan estas estructuras y realizan cálculos básicos.
                    ```go
                    // en element_models.go
                    type Element struct {
                        ID         string                 `json:"id"`
                        Width      int                    `json:"width"`
                        Height     int                    `json:"height"`
                        // ... otros campos ...
                        Frame      Frame                  `json:"frame"`
                        Winds      []Wind                 `json:"winds,omitempty"`
                    }

                    type FrameDetail struct {
                        Position       string  `json:"position"`
                        ProfileSKU     string  `json:"profile_sku"`
                        Dimension      int     `json:"dimension"` // Longitud de corte
                        AngleLeft      float64 `json:"angle_left"`
                        AngleRight     float64 `json:"angle_right"`
                        // ... otros campos ...
                    }
                    ```
                - **`catalog_models.go`**:
                    - Responsabilidad: Se espera que contenga modelos para catálogos de perfiles, vidrios, herrajes, etc.
                    - Estado Actual: El archivo está vacío.

        - **`internal/pkg/constants/`**
            - **Propósito e Importancia:** Centraliza todas las constantes literales usadas a través del proyecto, mejorando la mantenibilidad y reduciendo errores por strings mágicos. Define el vocabulario del dominio.
            - **Archivos Clave:**
                - **`constants.go`**
                    - Responsabilidad: Contiene definiciones de constantes para materiales, tipos de elementos, tipos de perfiles, estructuras, geometrías, posiciones, tipos de corte, tipos de hojas, estados, lados de apertura, etc.
                    - Agrupaciones de Constantes Clave y su Significado:
                        - `MATERIAL_*`: Define los materiales básicos (PVC, Aluminio, Cristal, etc.).
                        - `TYPE_*`: Tipologías generales de aberturas (Corredera, Abatible).
                        - `PROFILE_TYPE_*`: Clasificación detallada de los tipos de perfiles utilizados (marcos, hojas, traslapos, montantes, junquillos, refuerzos, etc.). Estos son fundamentales para la selección de perfiles en el catálogo y los cálculos de despiece.
                        - `STRUCTURE_*`: Define si el elemento es una Ventana o Puerta.
                        - `GEOMETRY_*`: Geometría de los marcos/elementos (Rectangular, Triangular).
                        - `POSITION_*`: Posiciones de los perfiles dentro de un marco u hoja (Izquierda, Derecha, Arriba, Abajo).
                        - `CUT_*` y `CUT_*_WIND`: Tipos de cortes aplicables a los perfiles (Cuadrado, Ángulo, etc.). Crucial para las instrucciones de corte.
                        - `WIND_KIND_*`: Variaciones de hojas (Corredera móvil, Fija, Abatible, Oscilobatiente, etc.).
                        - `WIND_STATUS_*`: Estado de una hoja (activa, inactiva).
                        - `OPENING_SIDE_*` y `OPENING_*`: Lados y dirección de apertura de las hojas.
                    ```go
                    // Ejemplo de agrupación en constants.go
                    const (
                        PROFILE_TYPE_SLIDING_FRAME     = "Marco corredera"
                        PROFILE_TYPE_SLIDING_WIND      = "Hoja corredera"
                        // ... más tipos de perfiles
                    )

                    const (
                        CUT_SQUARE = "Cuadrado"
                        CUT_ANGLE  = "Ángulo"
                        // ... más tipos de corte
                    )
                    ```
                    - Explicación: Estas constantes proporcionan un conjunto controlado de valores para propiedades de los modelos, facilitando la validación y la lógica de negocio específica (ej. qué perfiles son compatibles con qué tipos de elementos o qué tipo de corte aplicar).

        - **`internal/app/window-api/repositories/`**
            - **Propósito e Importancia:** Esta capa es responsable de abstraer las interacciones con la fuente de datos (ej. Supabase). Contendría la lógica para realizar operaciones CRUD (Crear, Leer, Actualizar, Eliminar) sobre las entidades del dominio.
            - **Archivos Clave y Responsabilidades Previstas:**
                - `project_repository.go`: Manejaría la persistencia y recuperación de los datos de `Project`.
                - `profile_catalog_repository.go`: Se encargaría de acceder a los datos del catálogo de perfiles (SKUs, propiedades, etc.).
                - `supabase_profile_catalog_repository.go`: Sería una implementación concreta de `profile_catalog_repository.go` utilizando el cliente Supabase.
            - **Estado Actual:** Todos los archivos (`project_repository.go`, `profile_catalog_repository.go`, `supabase_profile_catalog_repository.go`) en este paquete están actualmente vacíos. Esto indica que la lógica de acceso a datos aún no ha sido implementada en esta capa, y las interacciones directas con Supabase se realizan a través del `SupabaseClient` en `internal/pkg/apiclient/` como se vio en `main.go` para pruebas.

## 6. Funcionamiento e Interacción del Código (Identificado por Windsurf AI)
    - (Se completará después de analizar el código).

    ### 6.1. Flujos de Datos y Lógica de Negocio Clave
        - **Flujo de Inicialización y Consulta (actual):**
            1. `cmd/window-api/main.go` se ejecuta.
            2. Llama a `config.LoadConfig()` para cargar la configuración (URL y clave de Supabase) desde variables de entorno (o `.env`).
            3. Se crea una instancia de `apiclient.NewSupabaseClient()` con la configuración cargada.
            4. El cliente Supabase se utiliza para llamar a `QueryData()` con una ruta y parámetros específicos (ej., `/rest/v1/projects?select=...`).
            5. `QueryData()` construye la URL completa, crea una petición HTTP GET, añade las cabeceras `apikey` y `Authorization`, y envía la petición.
            6. La respuesta de Supabase se decodifica en la estructura proporcionada.
        - (La lógica de negocio específica, como el cálculo de perfiles, aún no es visible en el código analizado).

    ### 6.2. Endpoints de la API (si aplica)
        - Actualmente, el proyecto no expone ningún endpoint de API HTTP.
        - El archivo `internal/app/window-api/handlers/http_handlers.go` está vacío, y `cmd/window-api/main.go` no configura ni inicia un servidor HTTP (ej., usando `net/http` o un router como `gorilla/mux` o `gin-gonic`).
        - La estructura del proyecto sugiere que la intención es desarrollar una API, pero esta funcionalidad aún no está implementada.

    ### 6.3. Interacción Detallada con Supabase
        - La interacción con Supabase se centraliza en `internal/pkg/apiclient/SupabaseClient`.
        - El método `QueryData` es el principal medio de interacción mostrado hasta ahora. Realiza peticiones GET a la API REST de Supabase (PostgREST).
        - **Autenticación:** Utiliza la `ServiceRoleKey` tanto para la cabecera `apikey` como para el `Bearer token` en la cabecera `Authorization`.
          ```go
          // En SupabaseClient.QueryData()
          req.Header.Set("apikey", c.ServiceRoleKey)
          req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey)
          ```
        - **Construcción de Consultas:** Las consultas se definen pasando una `path` (ej., `/rest/v1/nombre_tabla`) y `queryParams` (ej., `select=columna1,columna2&filtro=eq.valor`).
          ```go
          // En cmd/window-api/main.go
          path := "/rest/v1/projects"
          queryParams := "select=id,name,created_at&order=created_at.desc&limit=2"
          ```
        - **Manejo de Respuestas:** Las respuestas JSON son decodificadas en la estructura `target` proporcionada usando `json.NewDecoder()`.
        - El cliente actual solo implementa la funcionalidad de consulta (GET). Las operaciones de creación (POST), actualización (PATCH/PUT) o eliminación (DELETE) no están implementadas en `SupabaseClient`.

    ### 6.4. Modelos de Datos (Structs de Go) y su Relevancia
        - El archivo `cmd/window-api/main.go` define una estructura local `ProjectFromAPI` como ejemplo para decodificar los resultados de la tabla `projects`:
          ```go
          type ProjectFromAPI struct {
              ID        int64     `json:"id"`
              Name      string    `json:"name"`
              CreatedAt time.Time `json:"created_at"`
          }
          ```
        - Los modelos de datos principales del dominio residen en `internal/app/window-api/models/` y han sido detallados en la Sección 5.1. Son cruciales para definir la estructura de un proyecto de fabricación de ventanas/puertas, desde el `Project` general hasta los `Element`s individuales y sus `FrameDetail` y `WindDetail` (que contienen información para el despiece como `ProfileSKU`, `Dimension` y ángulos de corte).
        - **Jerarquía principal de modelos:**
            - `Project` (contiene `Contact`, `ProjectCost`, `IvaRate`)
                - `Component` (agrupación lógica)
                    - `Module` (agrupación de `Element`s, opcionalmente unidos por `AuxProfile`)
                        - `Element` (ventana/puerta individual)
                            - `Frame` (marco del elemento)
                                - `FrameDetail` (pieza de perfil del marco con SKU, dimensiones, ángulos)
                            - `Wind` (hoja del elemento)
                                - `WindDetail` (pieza de perfil de la hoja con SKU, dimensiones, ángulos)
        - El paquete `internal/pkg/constants/constants.go` provee los valores literales controlados para muchos campos de estos modelos (ej., `Element.Type`, `Frame.CutType`, `Wind.Kind`).
        - **Relevancia para Supabase:** Estas estructuras Go (o un subconjunto/variación de ellas) probablemente se mapearán a tablas en la base de datos Supabase. Por ejemplo, podría haber tablas para `projects`, `components`, `elements`, `profiles_catalog`, etc. Los campos como `ProfileSKU` en `FrameDetail` y `WindDetail` sugieren una relación con una tabla de catálogo de perfiles. La capa de repositorios (`internal/app/window-api/repositories/`), una vez implementada, se encargaría de este mapeo y de las operaciones de base de datos.

    ### 6.5. Estrategia de Manejo de Errores Observada
        - Las funciones que pueden fallar (ej., `config.LoadConfig()`, `apiclient.SupabaseClient.QueryData()`) devuelven un `error` como último valor.
        - Se utiliza `fmt.Errorf` para envolver errores y añadir contexto.
        - En `main.go`, los errores fatales durante la inicialización (como no poder cargar la configuración) usan `log.Fatalf()`.
        - El cliente Supabase (`QueryData`) verifica los códigos de estado HTTP y trata los códigos no exitosos (fuera del rango 200-299) como errores, intentando incluir el cuerpo de la respuesta de error de Supabase.
        - Si la decodificación JSON falla, también se intenta registrar el cuerpo de la respuesta para depuración.

    ### 6.6. Configuración y Variables de Entorno (Basado en uso, no en archivo `.env`)
        - A partir del análisis de `internal/pkg/config/config.go`, las siguientes variables de entorno son cruciales:
            - `SUPABASE_API_URL`: La URL base de la API de Supabase (ej., `https://<project-ref>.supabase.co`).
            - `SUPABASE_SERVICE_KEY`: La clave de servicio (rol `service_role`) para la autenticación con la API de Supabase.
        - (Otras variables como `API_PORT` están comentadas en el código de configuración, indicando que podrían usarse en el futuro).

## 7. Análisis Estático: Problemas, Advertencias y Oportunidades de Mejora Detectadas
    - (Se completará después del análisis).

    ### 7.1. Errores de Compilación o Linter
        - (No se han detectado errores de compilación evidentes en los fragmentos analizados. Se requeriría una ejecución de `go build` o linters para una evaluación completa).

    ### 7.2. Código "Sospechoso" (Code Smells) o Anti-patrones
        - **`main.go` como Script de Prueba:** El archivo `cmd/window-api/main.go` actualmente funciona más como un script para probar la conexión y consultas a Supabase que como el punto de entrada de un servidor API persistente. Carece de la configuración de un router HTTP y el inicio de un servidor.
        - **Archivo de Configuración YAML no Utilizado:** La presencia del archivo `configs/config.dev.yaml` (actualmente vacío) sugiere una intención de usar configuración basada en archivos YAML, pero `internal/pkg/config/config.go` solo implementa la carga desde variables de entorno. Esto podría causar confusión o indicar una funcionalidad incompleta o abandonada.
        - **Manejo de Cuerpo de Respuesta en `supabase_client.go`:** La función `getResetBody` y su uso en `QueryData` para re-leer el cuerpo de la respuesta en caso de error de decodificación JSON es un poco compleja y podría simplificarse o hacerse más robusta. Leer el cuerpo una vez y luego intentar decodificarlo (y usar esos bytes para el mensaje de error si falla la decodificación) es una aproximación más común.

    ### 7.3. Oportunidades de Refactorización o Mejora
        - **Implementar Servidor HTTP:** El paso más obvio es desarrollar la funcionalidad del servidor API en `main.go`, utilizando un router (del paquete `net/http` o una librería externa) y conectándolo a los manejadores que se desarrollarán en `internal/app/window-api/handlers/`.
        - **Clarificar Estrategia de Configuración:** Decidir si se usará `config.dev.yaml` (y otros archivos YAML por entorno) e implementar su carga en `config.go`, o eliminar el archivo si la configuración será exclusivamente mediante variables de entorno.
        - **Expandir `SupabaseClient`:** Añadir métodos para otras operaciones CRUD (POST, PATCH, DELETE) si son necesarias para las implementaciones de los repositorios.
        - **Logging Estructurado:** Considerar el uso de una librería de logging más avanzada (como `logrus`, que ya está en `go.sum` como dependencia indirecta, o `zap`) para un logging estructurado y configurable, especialmente si la aplicación crece.
        - **Definición de Modelos de Catálogo:** Completar los modelos en `internal/app/window-api/models/catalog_models.go` para perfiles, vidrios, herrajes, etc., ya que son fundamentales para la lógica de cálculo.
        - **Implementar Lógica de Cálculo de Perfiles:** Esta es una funcionalidad clave ausente. Requerirá servicios que utilicen los modelos de `Element`, `FrameDetail`, `WindDetail` y los modelos de catálogo para determinar las dimensiones exactas y los cortes de cada perfil.
        - **Validaciones en Constructores:** Activar y completar las validaciones en los constructores de modelos (ej., `NewFrame`, `NewElement`) usando las constantes definidas en `internal/pkg/constants/constants.go`.

## 8. Áreas de Complejidad o Interés Especial (Identificadas por Windsurf AI)
    - **Modelado Detallado del Dominio:** El proyecto muestra un esfuerzo significativo en modelar el dominio de la fabricación de ventanas y puertas con un buen nivel de detalle (jerarquía de `Project` -> `Component` -> `Module` -> `Element` -> `Frame/Wind` -> `FrameDetail/WindDetail`). Esto es fundamental para la precisión de cualquier cálculo posterior.
    - **Importancia de `constants.go`:** El archivo `internal/pkg/constants/constants.go` es crucial. Define un vocabulario controlado para muchos atributos, lo que es vital para la lógica de negocio y las validaciones. La correcta gestión y uso de estas constantes será importante para la robustez del sistema.
    - **Lógica de Cálculo de Perfiles (Pendiente):** La funcionalidad más compleja y crítica, el cálculo de las dimensiones exactas de corte de los perfiles (`FrameDetail.Dimension`, `WindDetail.Dimension`) y sus ángulos, aún no está implementada en ninguna de las capas analizadas (`handlers`, `services`, `repositories`). Esta lógica probablemente residirá en el paquete de `services` (una vez desarrollado) y requerirá una interacción cuidadosa entre los `Element` (con sus dimensiones y tipos), y un catálogo de perfiles (actualmente no definido en `catalog_models.go` ni implementado en `repositories`) que especifique las propiedades de cada SKU de perfil (ej., anchos, alturas, descuentos de marco/hoja, etc.). La identificación de cómo se descuentan las medidas de los perfiles según el tipo de sistema (corredera, abatible), el tipo de unión (soldado a 45º, atornillado a 90º), y las propiedades específicas de cada perfil del catálogo será el núcleo de esta complejidad. Su ausencia actual significa que una parte fundamental del propósito de la aplicación no está operativa.
    - (Se destacarán más áreas a medida que se analice la implementación de la lógica de negocio y los servicios).

## 9. Conclusión (Inferida del Estado del Proyecto)
    - El proyecto Windraw sienta las bases para una aplicación de gestión de fabricación de ventanas, con una estructura de directorios que sugiere una arquitectura por capas (handlers, services, repositories, models) y modelos de datos bien definidos para el dominio principal. La conexión inicial a Supabase y la carga de configuración también están presentes.
    - Sin embargo, las capas funcionales clave como `handlers`, `services`, y `repositories` están actualmente vacías. Esto significa que la funcionalidad principal de la API (endpoints HTTP), la lógica de negocio central (incluyendo el crítico cálculo de despiece de perfiles), y la interacción estructurada con la base de datos para las entidades del dominio aún no han sido implementadas.
    - La atención a la complejidad de la lógica de cálculo de perfiles y la gestión de las constantes del dominio será crucial para la precisión y robustez del sistema, y su implementación probablemente involucre las capas de `models`, `repositories` y `services` (una vez desarrolladas).
    - La implementación de la capa de repositorios para el acceso a datos y la definición de modelos de catálogo son fundamentales para completar la arquitectura del proyecto y permitir la implementación de la lógica de negocio.

## 10. Próximos Pasos Recomendados (para el desarrollo del proyecto Windraw)

Basado en el análisis del estado actual del proyecto Windraw, se recomiendan los siguientes pasos para avanzar en su desarrollo y completar la funcionalidad prevista:

1.  **Implementar el Servidor HTTP y Routing Básico:**
    *   En `cmd/window-api/main.go`, configurar un servidor HTTP (utilizando el paquete `net/http` estándar de Go o una librería de routing como `gorilla/mux` o `chi`).
    *   Definir rutas iniciales que se conectarán a los futuros manejadores.

2.  **Desarrollar los Manejadores HTTP (`handlers`):**
    *   En `internal/app/window-api/handlers/http_handlers.go`, implementar las funciones que manejarán las solicitudes HTTP para cada ruta definida.
    *   Estos manejadores recibirán las solicitudes, realizarán validaciones básicas de entrada y llamarán a los servicios apropiados.
    *   Se encargarán de serializar las respuestas (ej., a JSON) y enviarlas de vuelta al cliente.

3.  **Poblar los Modelos de Catálogo (`catalog_models.go`):**
    *   Definir las estructuras de Go en `internal/app/window-api/models/catalog_models.go` para representar los perfiles, vidrios, herrajes y otros componentes de catálogo.
    *   Estas estructuras deben contener todos los atributos necesarios para la lógica de negocio, especialmente para el cálculo de despiece (ej., para perfiles: SKU, descripción, material, dimensiones estándar, descuentos aplicables para cortes, secciones transversales si fuera necesario para cálculos más avanzados, etc.).

4.  **Implementar la Capa de Repositorios (`repositories`):**
    *   Desarrollar las implementaciones concretas para `project_repository.go` y `profile_catalog_repository.go` (y su variante `supabase_profile_catalog_repository.go`).
    *   Estas implementaciones interactuarán con Supabase (utilizando el `SupabaseClient` existente o expandiéndolo) para realizar operaciones CRUD (Crear, Leer, Actualizar, Eliminar) sobre las tablas correspondientes a los `Project`s y a los datos del catálogo de perfiles.
    *   Definir interfaces para estos repositorios para facilitar el testing y la posible sustitución de la capa de persistencia en el futuro.

5.  **Desarrollar la Capa de Servicios (`services`):**
    *   En `internal/app/window-api/services/`, implementar la lógica de negocio principal.
    *   **Servicio de Cálculo de Perfiles (Crítico):** Crear un servicio que tome un `Element` (con sus `Frame`s y `Wind`s) como entrada, consulte el catálogo de perfiles (a través del `ProfileCatalogRepository`) para obtener las propiedades de los `ProfileSKU` especificados, y calcule las dimensiones de corte exactas (`Dimension`) y los ángulos (`AngleLeft`, `AngleRight`) para cada `FrameDetail` y `WindDetail`.
        *   Esta lógica deberá considerar los tipos de perfiles, tipos de corte (definidos en `constants.go`), descuentos necesarios según el sistema (corredera, abatible), y las dimensiones generales del `Element`.
    *   Otros servicios podrían incluir la gestión de proyectos (creación, actualización, obtención), validaciones de negocio complejas, etc.
    *   Los servicios orquestarán las llamadas a los repositorios.

6.  **Expandir y Refinar `SupabaseClient`:**
    *   Añadir métodos para operaciones POST, PATCH, DELETE si son necesarias para las implementaciones de los repositorios.
    *   Mejorar el manejo de errores y la configuración del cliente si es necesario.

7.  **Implementar Validaciones Exhaustivas:**
    *   Activar y completar las validaciones en los constructores de los modelos (`NewFrame`, `NewWind`, `NewElement`) utilizando las `constants.go`.
    *   Añadir validaciones en los manejadores (para datos de entrada) y en los servicios (para reglas de negocio).

8.  **Estrategia de Configuración:**
    *   Decidir si se utilizará exclusivamente variables de entorno (como actualmente) o si se implementará la carga desde archivos YAML como `configs/config.dev.yaml`. Si se opta por YAML, actualizar `internal/pkg/config/config.go`.

9.  **Testing:**
    *   Implementar tests unitarios para los modelos, constantes, servicios y repositorios (usando mocks para las dependencias externas como Supabase en los tests de servicio).
    *   Implementar tests de integración para verificar la correcta interacción entre componentes, especialmente con la base de datos.

10. **Logging y Monitorización:**
    *   Considerar la integración de una librería de logging estructurado (ej. `logrus` o `zap`) para mejorar la observabilidad de la aplicación, especialmente a medida que crezca en complejidad.

Estos pasos proporcionan una hoja de ruta para transformar el proyecto Windraw de una base bien estructurada a una aplicación completamente funcional.

## 11. Anexos (Opcional)
    - (Espacio para diagramas de arquitectura detallados, esquemas de base de datos, etc., si fueran necesarios y disponibles).
