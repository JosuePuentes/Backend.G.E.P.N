package handlers

import (
	"encoding/json"
	"net/http"
)

// HomeHandler maneja la ruta raíz
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"message": "Hola Mundo",
		"status":  "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthHandler verifica el estado del servidor
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status": "healthy",
		"service": "GEPN API",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CiudadanoHandler maneja la ruta para ciudadanos
func CiudadanoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"message": "Portal Ciudadano - GEPN",
		"status":  "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// FaviconHandler maneja las peticiones de favicon (evita 404 en logs)
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	// Responder con 204 No Content para evitar el 404
	w.WriteHeader(http.StatusNoContent)
}

