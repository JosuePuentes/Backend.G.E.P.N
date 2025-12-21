package handlers

import (
	"encoding/json"
	"fmt"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"
)

// CrearDenunciaHandler crea una nueva denuncia
func CrearDenunciaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ciudadano del token
	authHeader := r.Header.Get("Authorization")
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}
	
	ciudadano, ok := GetCiudadanoFromToken(token)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No autorizado. Token inválido o expirado",
		})
		return
	}

	var req models.CrearDenunciaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if req.Denunciante.Nombre == "" || req.Denunciante.Cedula == "" ||
		req.Denunciante.Telefono == "" || req.Denuncia.Motivo == "" ||
		req.Denuncia.Hechos == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Todos los campos obligatorios deben ser completados",
		})
		return
	}

	// Generar número de denuncia
	numeroDenuncia, err := database.GenerarNumeroDenuncia()
	if err != nil {
		http.Error(w, "Error al generar número de denuncia", http.StatusInternalServerError)
		return
	}

	// Convertir fecha de nacimiento si viene en formato DD/MM/AAAA
	var fechaNacimiento string
	if req.Denunciante.FechaNacimiento != "" {
		// Intentar parsear formato DD/MM/AAAA
		var dia, mes, año string
		_, err := fmt.Sscanf(req.Denunciante.FechaNacimiento, "%2s/%2s/%4s", &dia, &mes, &año)
		if err == nil {
			fechaNacimiento = fmt.Sprintf("%s-%s-%s", año, mes, dia)
		} else {
			fechaNacimiento = req.Denunciante.FechaNacimiento
		}
	}

	// Crear denuncia
	denuncia := &models.Denuncia{
		CiudadanoID:               ciudadano.ID,
		NumeroDenuncia:            numeroDenuncia,
		NombreDenunciante:         req.Denunciante.Nombre,
		CedulaDenunciante:         req.Denunciante.Cedula,
		TelefonoDenunciante:       req.Denunciante.Telefono,
		FechaNacimientoDenunciante: fechaNacimiento,
		ParroquiaDenunciante:      req.Denunciante.Parroquia,
		Motivo:                    req.Denuncia.Motivo,
		Hechos:                    req.Denuncia.Hechos,
		NombreDenunciado:          req.Denunciado.Nombre,
		DireccionDenunciado:       req.Denunciado.Direccion,
		EstadoDenunciado:          req.Denunciado.Estado,
		MunicipioDenunciado:       req.Denunciado.Municipio,
		ParroquiaDenunciado:       req.Denunciado.Parroquia,
		FechaDenuncia:             time.Now(),
		Estado:                    "Pendiente",
	}

	if err := database.CrearDenuncia(denuncia); err != nil {
		http.Error(w, "Error al registrar la denuncia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Denuncia registrada correctamente",
		"data": map[string]interface{}{
			"id":             denuncia.ID.Hex(),
			"numero_denuncia": denuncia.NumeroDenuncia,
			"fecha":          denuncia.FechaDenuncia.Format(time.RFC3339),
		},
	})
}

// MisDenunciasHandler obtiene las denuncias del usuario autenticado
func MisDenunciasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener ciudadano del token
	authHeader := r.Header.Get("Authorization")
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		token = authHeader
	}
	
	ciudadano, ok := GetCiudadanoFromToken(token)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No autorizado. Token inválido o expirado",
		})
		return
	}

	denuncias, err := database.ObtenerDenunciasPorCiudadano(ciudadano.ID)
	if err != nil {
		http.Error(w, "Error al obtener denuncias", http.StatusInternalServerError)
		return
	}

	// Formatear denuncias para la respuesta
	denunciasResp := make([]map[string]interface{}, len(denuncias))
	for i, d := range denuncias {
		denunciasResp[i] = map[string]interface{}{
			"id":             d.ID.Hex(),
			"numero_denuncia": d.NumeroDenuncia,
			"motivo":         d.Motivo,
			"fecha_denuncia": d.FechaDenuncia.Format(time.RFC3339),
			"estado":         d.Estado,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    denunciasResp,
	})
}

