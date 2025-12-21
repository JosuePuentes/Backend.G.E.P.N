package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
)

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
		OficialID: usuario.ID,
		Latitud:   req.Latitud,
		Longitud:  req.Longitud,
		Ubicacion: req.Ubicacion,
		Estado:    "activo",
	}

	if err := database.CrearAlertaPanico(&alerta); err != nil {
		http.Error(w, "Error al crear alerta de pánico", http.StatusInternalServerError)
		return
	}

	// Notificar a oficiales cercanos (dentro de 5 km)
	oficialesCercanos, err := database.ObtenerOficialesEnGuardiaCercanos(
		req.Latitud,
		req.Longitud,
		5.0, // Radio de 5 km
	)
	if err == nil {
		// En producción, aquí se enviarían notificaciones push, SMS, etc.
		// Por ahora, solo lo registramos en la respuesta
		response := map[string]interface{}{
			"mensaje":            "Alerta de pánico activada",
			"alerta":             alerta,
			"oficiales_notificados": len(oficialesCercanos),
			"oficiales":          oficialesCercanos,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else {
		response := map[string]interface{}{
			"mensaje": "Alerta de pánico activada",
			"alerta":  alerta,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
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

	alertas, err := database.ListarAlertasPanico()
	if err != nil {
		http.Error(w, "Error al listar alertas de pánico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alertas)
}

