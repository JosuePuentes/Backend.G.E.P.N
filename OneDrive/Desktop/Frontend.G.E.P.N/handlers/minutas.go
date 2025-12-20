package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	minuta.OficialID = usuario.ID

	if err := database.CrearMinuta(&minuta); err != nil {
		http.Error(w, "Error al crear minuta", http.StatusInternalServerError)
		return
	}

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

	minutas, err := database.ListarMinutas()
	if err != nil {
		http.Error(w, "Error al listar minutas", http.StatusInternalServerError)
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
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	minuta, err := database.ObtenerMinutaPorID(id)
	if err != nil {
		http.Error(w, "Minuta no encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(minuta)
}

