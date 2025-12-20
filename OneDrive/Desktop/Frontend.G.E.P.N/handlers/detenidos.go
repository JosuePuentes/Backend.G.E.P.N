package handlers

import (
	"encoding/json"
	"gepn/models"
	"net/http"
	"strconv"
	"time"
)

// Almacenamiento temporal en memoria
var detenidos = []models.Detenido{}
var detenidoIDCounter = 1

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

	detenido.ID = detenidoIDCounter
	detenidoIDCounter++
	detenido.OficialID = usuario.ID
	detenido.FechaDetencion = time.Now()
	if detenido.Estado == "" {
		detenido.Estado = "detenido"
	}

	detenidos = append(detenidos, detenido)

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
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	for _, d := range detenidos {
		if d.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(d)
			return
		}
	}

	http.Error(w, "Detenido no encontrado", http.StatusNotFound)
}

