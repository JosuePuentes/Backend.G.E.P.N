package handlers

import (
	"encoding/json"
	"gepn/models"
	"net/http"
	"time"
)

// Almacenamiento temporal en memoria
var alertasPanico = []models.Panico{}
var panicoIDCounter = 1

// ActivarPanicoHandler activa el botón de pánico
func ActivarPanicoHandler(w http.ResponseWriter, r *http.Request) {
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

	var req struct {
		Latitud   float64 `json:"latitud"`
		Longitud  float64 `json:"longitud"`
		Ubicacion string  `json:"ubicacion"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Crear alerta de pánico
	alerta := models.Panico{
		ID:              panicoIDCounter,
		OficialID:       usuario.ID,
		Latitud:         req.Latitud,
		Longitud:        req.Longitud,
		Ubicacion:       req.Ubicacion,
		FechaActivacion: time.Now(),
		Estado:          "activo",
	}
	panicoIDCounter++

	alertasPanico = append(alertasPanico, alerta)

	// En producción, aquí se enviaría notificación a central, otros oficiales, etc.
	response := map[string]interface{}{
		"mensaje": "Alerta de pánico activada",
		"alerta":  alerta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListarAlertasPanicoHandler lista las alertas de pánico
func ListarAlertasPanicoHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(alertasPanico)
}

