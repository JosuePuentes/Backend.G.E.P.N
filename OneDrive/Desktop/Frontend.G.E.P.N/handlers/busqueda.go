package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
)

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

	// Buscar en los más buscados
	masBuscado, err := database.BuscarMasBuscadoPorCedula(req.Cedula)
	var encontrado bool
	var resultado models.MasBuscado
	
	busqueda := models.BusquedaCedula{
		Cedula:    req.Cedula,
		OficialID: usuario.ID,
	}
	
	if err == nil && masBuscado != nil {
		encontrado = true
		resultado = *masBuscado
		busqueda.Resultado = "encontrado"
		busqueda.Nombre = masBuscado.Nombre
		busqueda.Apellido = masBuscado.Apellido
	} else {
		busqueda.Resultado = "no_encontrado"
	}

	// Registrar búsqueda en MongoDB
	database.CrearBusqueda(&busqueda)

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

	masBuscados, err := database.ListarMasBuscados()
	if err != nil {
		http.Error(w, "Error al listar más buscados", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masBuscados)
}

