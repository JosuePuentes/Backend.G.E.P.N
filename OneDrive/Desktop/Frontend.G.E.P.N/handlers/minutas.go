package handlers

import (
	"encoding/json"
	"gepn/models"
	"net/http"
	"strconv"
	"time"
)

// Almacenamiento temporal en memoria
var minutas = []models.Minuta{}
var minutaIDCounter = 1

// CrearMinutaHandler crea una nueva minuta digital
func CrearMinutaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener usuario del token
	token := r.Header.Get("Authorization")
	usuario, ok := GetUsuarioFromToken(token)
	if !ok {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	var minuta models.Minuta
	if err := json.NewDecoder(r.Body).Decode(&minuta); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	minuta.ID = minutaIDCounter
	minutaIDCounter++
	minuta.OficialID = usuario.ID
	minuta.FechaCreacion = time.Now()

	minutas = append(minutas, minuta)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(minuta)
}

// ListarMinutasHandler lista todas las minutas
func ListarMinutasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener usuario del token
	token := r.Header.Get("Authorization")
	_, ok := GetUsuarioFromToken(token)
	if !ok {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(minutas)
}

// ObtenerMinutaHandler obtiene una minuta por ID
func ObtenerMinutaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener usuario del token
	token := r.Header.Get("Authorization")
	_, ok := GetUsuarioFromToken(token)
	if !ok {
		http.Error(w, "No autorizado", http.StatusUnauthorized)
		return
	}

	// Obtener ID de la query string
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, m := range minutas {
		if m.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(m)
			return
		}
	}

	http.Error(w, "Minuta no encontrada", http.StatusNotFound)
}

