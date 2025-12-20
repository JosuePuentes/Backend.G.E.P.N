package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CrearDetenidoHandler crea un nuevo registro de detenido
func CrearDetenidoHandler(w http.ResponseWriter, r *http.Request) {
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

	var detenido models.Detenido
	if err := json.NewDecoder(r.Body).Decode(&detenido); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	detenido.OficialID = usuario.ID
	if detenido.Estado == "" {
		detenido.Estado = "detenido"
	}

	if err := database.CrearDetenido(&detenido); err != nil {
		http.Error(w, "Error al crear detenido", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(detenido)
}

// ListarDetenidosHandler lista todos los detenidos
func ListarDetenidosHandler(w http.ResponseWriter, r *http.Request) {
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

	detenidos, err := database.ListarDetenidos()
	if err != nil {
		http.Error(w, "Error al listar detenidos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detenidos)
}

// ObtenerDetenidoHandler obtiene un detenido por ID
func ObtenerDetenidoHandler(w http.ResponseWriter, r *http.Request) {
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

	detenido, err := database.ObtenerDetenidoPorID(id)
	if err != nil {
		http.Error(w, "Detenido no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detenido)
}

