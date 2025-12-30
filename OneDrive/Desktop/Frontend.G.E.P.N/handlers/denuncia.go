package handlers

import (
	"encoding/json"
	"fmt"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

// ListarTodasDenunciasHandler lista todas las denuncias (para usuarios del sistema con permiso "denuncias")
func ListarTodasDenunciasHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autenticación de master
	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No autorizado. Token inválido o expirado",
		})
		return
	}

	// Verificar permiso de denuncias
	tienePermiso := false
	for _, permiso := range master.Permisos {
		if permiso == "denuncias" {
			tienePermiso = true
			break
		}
	}

	if !tienePermiso {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No tiene permisos para acceder a denuncias",
		})
		return
	}

	// Obtener parámetros de paginación y filtros
	page := 1
	limit := 20
	estado := ""

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		fmt.Sscanf(pageStr, "%d", &page)
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}
	if estadoStr := r.URL.Query().Get("estado"); estadoStr != "" {
		estado = estadoStr
	}

	denuncias, total, err := database.ListarTodasDenuncias(page, limit, estado)
	if err != nil {
		http.Error(w, "Error al obtener denuncias", http.StatusInternalServerError)
		return
	}

	// Formatear denuncias para la respuesta
	denunciasResp := make([]map[string]interface{}, len(denuncias))
	for i, d := range denuncias {
		denunciasResp[i] = map[string]interface{}{
			"id":                        d.ID.Hex(),
			"ciudadano_id":              d.CiudadanoID.Hex(),
			"numero_denuncia":            d.NumeroDenuncia,
			"nombre_denunciante":         d.NombreDenunciante,
			"cedula_denunciante":         d.CedulaDenunciante,
			"telefono_denunciante":       d.TelefonoDenunciante,
			"fecha_nacimiento_denunciante": d.FechaNacimientoDenunciante,
			"parroquia_denunciante":      d.ParroquiaDenunciante,
			"motivo":                     d.Motivo,
			"hechos":                     d.Hechos,
			"nombre_denunciado":          d.NombreDenunciado,
			"direccion_denunciado":       d.DireccionDenunciado,
			"estado_denunciado":          d.EstadoDenunciado,
			"municipio_denunciado":        d.MunicipioDenunciado,
			"parroquia_denunciado":       d.ParroquiaDenunciado,
			"fecha_denuncia":             d.FechaDenuncia.Format(time.RFC3339),
			"estado":                     d.Estado,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    denunciasResp,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// ObtenerDenunciaHandler obtiene una denuncia específica por ID
func ObtenerDenunciaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autenticación de master
	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No autorizado. Token inválido o expirado",
		})
		return
	}

	// Verificar permiso de denuncias
	tienePermiso := false
	for _, permiso := range master.Permisos {
		if permiso == "denuncias" {
			tienePermiso = true
			break
		}
	}

	if !tienePermiso {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No tiene permisos para acceder a denuncias",
		})
		return
	}

	// Obtener ID de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID de denuncia requerido",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	denuncia, err := database.ObtenerDenunciaPorID(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Denuncia no encontrada",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":                        denuncia.ID.Hex(),
			"ciudadano_id":              denuncia.CiudadanoID.Hex(),
			"numero_denuncia":           denuncia.NumeroDenuncia,
			"nombre_denunciante":        denuncia.NombreDenunciante,
			"cedula_denunciante":        denuncia.CedulaDenunciante,
			"telefono_denunciante":       denuncia.TelefonoDenunciante,
			"fecha_nacimiento_denunciante": denuncia.FechaNacimientoDenunciante,
			"parroquia_denunciante":     denuncia.ParroquiaDenunciante,
			"motivo":                    denuncia.Motivo,
			"hechos":                    denuncia.Hechos,
			"nombre_denunciado":         denuncia.NombreDenunciado,
			"direccion_denunciado":      denuncia.DireccionDenunciado,
			"estado_denunciado":         denuncia.EstadoDenunciado,
			"municipio_denunciado":      denuncia.MunicipioDenunciado,
			"parroquia_denunciado":      denuncia.ParroquiaDenunciado,
			"fecha_denuncia":            denuncia.FechaDenuncia.Format(time.RFC3339),
			"estado":                    denuncia.Estado,
		},
	})
}

// ActualizarEstadoDenunciaHandler actualiza el estado de una denuncia
func ActualizarEstadoDenunciaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar autenticación de master
	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No autorizado. Token inválido o expirado",
		})
		return
	}

	// Verificar permiso de denuncias
	tienePermiso := false
	for _, permiso := range master.Permisos {
		if permiso == "denuncias" {
			tienePermiso = true
			break
		}
	}

	if !tienePermiso {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "No tiene permisos para acceder a denuncias",
		})
		return
	}

	var req struct {
		ID     string `json:"id"`
		Estado string `json:"estado"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al decodificar la petición",
		})
		return
	}

	// Validar estado
	estadosValidos := []string{"Pendiente", "En Proceso", "Resuelta", "Archivada"}
	estadoValido := false
	for _, e := range estadosValidos {
		if req.Estado == e {
			estadoValido = true
			break
		}
	}

	if !estadoValido {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Estado inválido. Estados válidos: Pendiente, En Proceso, Resuelta, Archivada",
		})
		return
	}

	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID inválido",
		})
		return
	}

	if err := database.ActualizarEstadoDenuncia(id, req.Estado); err != nil {
		http.Error(w, "Error al actualizar el estado de la denuncia", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Estado de denuncia actualizado correctamente",
	})
}

