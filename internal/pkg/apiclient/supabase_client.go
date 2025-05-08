package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	// Reemplaza 'github.com/tu-usuario/tu-proyecto-ventanas' con tu nombre de módulo
	"github.com/tu-usuario/tu-proyecto-ventanas/internal/pkg/config"
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
	req.Header.Set("apikey", c.ServiceRoleKey) // Usamos la service_role key como apikey
	req.Header.Set("Authorization", "Bearer "+c.ServiceRoleKey) // Y también como Bearer token
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "return=representation") // Para que devuelva los datos en la respuesta

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error haciendo petición a Supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error de Supabase API (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Decodificar la respuesta JSON en la estructura 'target'
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		// Leer el cuerpo para depuración si falla el JSON
		// Esto es útil porque a veces el error no es por el JSON sino por la respuesta vacía
		// o un error previo que no se manejó.
		bodyBytes, readErr := io.ReadAll(io.MultiReader(bytes.NewReader(getResetBody(resp)), resp.Body))
		if readErr == nil && len(bodyBytes) > 0 {
			fmt.Printf("Cuerpo de la respuesta (error JSON): %s\n", string(bodyBytes))
		}
		return fmt.Errorf("error decodificando respuesta JSON: %w", err)
	}

	return nil
}

// getResetBody es una función helper para poder leer el cuerpo de la respuesta múltiples veces
// si es necesario para depuración, ya que resp.Body solo se puede leer una vez.
// Esta función no es estrictamente necesaria para el funcionamiento básico.
func getResetBody(resp *http.Response) []byte {
	if resp == nil || resp.Body == nil {
		return nil
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // "Resetea" el cuerpo
	return bodyBytes
}