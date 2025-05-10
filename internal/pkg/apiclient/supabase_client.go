package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	// Usa el nombre de tu módulo Go aquí
	"github.com/mvialf/windraw/internal/pkg/config"
)

// SupabaseClient encapsula la configuración para hacer llamadas a la API de Supabase.
type SupabaseClient struct {
	BaseURL        string
	ServiceRoleKey string
	HttpClient     *http.Client
}

// NewSupabaseClient crea una nueva instancia de SupabaseClient.
func NewSupabaseClient(cfg *config.APIConfig) *SupabaseClient {
	return &SupabaseClient{
		BaseURL:        cfg.BaseURL,
		ServiceRoleKey: cfg.ServiceRoleKey,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second, // Timeout para las peticiones
		},
	}
}

// QueryData hace una petición GET a un endpoint de Supabase (PostgREST)
// path: ej. "/rest/v1/nombre_tabla"
// queryParams: ej. "select=columna1,columna2&otra_columna=eq.valor"
// target: un puntero a la estructura o slice de estructuras donde decodificar la respuesta JSON.
func (c *SupabaseClient) QueryData(path string, queryParams string, target interface{}) error {
	fullURL := c.BaseURL + path
	if queryParams != "" {
		fullURL += "?" + queryParams
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("error creando petición GET: %w", err)
	}

	// Cabeceras importantes para Supabase
	req.Header.Set("apikey", c.ServiceRoleKey)                  // Usamos la service_role key como apikey
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey) // Y también como Bearer token
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Prefer", "return=representation") // Opcional para GET, más útil para POST/PATCH/PUT

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error haciendo petición a Supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body) // Intenta leer el cuerpo del error
		return fmt.Errorf("error de Supabase API (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Decodificar la respuesta JSON en la estructura 'target'
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		// Intenta leer el cuerpo si la decodificación JSON falla, para ver qué devolvió Supabase
		// Necesitamos "resetear" el lector del cuerpo porque ya fue leído (o parcialmente leído) por NewDecoder
		// Esta parte es para una mejor depuración.
		// Primero, guardamos el cuerpo original para no perderlo.
		originalBodyBytes, readErr := io.ReadAll(io.MultiReader(bytes.NewReader(getResetBody(resp)), resp.Body))
		// MultiReader y getResetBody es un truco para intentar leer un cuerpo que ya pudo haber sido leído.

		// Si la lectura del cuerpo original tuvo éxito y la decodificación JSON falló, muestra el cuerpo.
		if readErr == nil && len(originalBodyBytes) > 0 {
			fmt.Printf("Cuerpo de la respuesta (error JSON - status %d): %s\n", resp.StatusCode, string(originalBodyBytes))
		}
		return fmt.Errorf("error decodificando respuesta JSON de Supabase: %w. Status: %d", err, resp.StatusCode)
	}

	return nil
}

// getResetBody es una función auxiliar para intentar leer un http.Response.Body
// que ya podría haber sido leído, guardándolo y reemplazándolo con un nuevo lector.
// Esto es útil para depuración si json.NewDecoder falla.
func getResetBody(resp *http.Response) []byte {
	if resp == nil || resp.Body == nil {
		return nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		// No podemos hacer mucho si hay un error leyendo el cuerpo aquí
		return nil
	}
	// "Reemplaza" el cuerpo original con un nuevo lector del mismo contenido
	// para que pueda ser leído de nuevo si es necesario.
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}
