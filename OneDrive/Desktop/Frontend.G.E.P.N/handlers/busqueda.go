package handlers

import (
	"encoding/json"
	"gepn/models"
	"net/http"
	"time"
)

// Almacenamiento temporal en memoria
var busquedas = []models.BusquedaCedula{}
var masBuscados = []models.MasBuscado{
	{
		ID:           1,
		Cedula:       "1234567890",
		Nombre:       "Juan",
		Apellido:     "Delincuente",
		Motivo:       "Robo a mano armada",
		Prioridad:    "alta",
		VecesBuscado: 15,
	},
	{
		ID:           2,
		Cedula:       "0987654321",
		Nombre:       "María",
		Apellido:     "Fugitiva",
		Motivo:       "Homicidio",
		Prioridad:    "alta",
		VecesBuscado: 12,
	},
}

// BuscarCedulaHandler busca una cédula
func BuscarCedulaHandler(w http.ResponseWriter, r *http.Request) {
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
		Cedula string `json:"cedula"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Registrar búsqueda
	busqueda := models.BusquedaCedula{
		Cedula:        req.Cedula,
		FechaBusqueda: time.Now(),
		OficialID:     usuario.ID,
	}

	// Buscar en los más buscados
	var encontrado bool
	var resultado models.MasBuscado
	for _, mb := range masBuscados {
		if mb.Cedula == req.Cedula {
			encontrado = true
			resultado = mb
			busqueda.Resultado = "encontrado"
			busqueda.Nombre = mb.Nombre
			busqueda.Apellido = mb.Apellido
			break
		}
	}

	if !encontrado {
		busqueda.Resultado = "no_encontrado"
	}

	busquedas = append(busquedas, busqueda)

	// Respuesta
	response := map[string]interface{}{
		"cedula":   req.Cedula,
		"encontrado": encontrado,
	}

	if encontrado {
		response["persona"] = resultado
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListarMasBuscadosHandler lista los más buscados
func ListarMasBuscadosHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(masBuscados)
}

