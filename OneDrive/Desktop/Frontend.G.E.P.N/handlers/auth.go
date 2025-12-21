package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"
)

// Almacenamiento temporal de tokens en memoria (los tokens pueden estar en memoria)
var tokens = make(map[string]*models.Usuario)

// LoginPolicialHandler maneja el login de policiales
func LoginPolicialHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validar credenciales desde MongoDB
	usuario, err := database.ObtenerUsuarioPorCredencial(req.Credencial)
	if err != nil || usuario.PIN != req.PIN {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Credenciales inválidas",
		})
		return
	}

	if !usuario.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario inactivo",
		})
		return
	}

	// Registrar inicio de guardia si se proporciona ubicación
	if req.Latitud != 0 && req.Longitud != 0 {
		guardia := models.Guardia{
			OficialID:     usuario.ID,
			FechaInicio:   time.Now(),
			LatitudInicio: req.Latitud,
			LongitudInicio: req.Longitud,
			Activa:        true,
		}
		if err := database.CrearGuardia(&guardia); err != nil {
			// Log error pero continuar con el login
			// En producción, considerar si esto debe fallar el login
		} else {
			// Actualizar usuario con estado de guardia y ubicación
			usuario.EnGuardia = true
			usuario.Latitud = req.Latitud
			usuario.Longitud = req.Longitud
			database.ActualizarUsuario(usuario)
		}
	}

	// Generar token simple (en producción usar JWT)
	token := generateToken()
	tokens[token] = usuario

	// Crear respuesta sin el PIN
	usuarioResp := *usuario
	usuarioResp.PIN = ""

	response := models.LoginResponse{
		Token:   token,
		Usuario: usuarioResp,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateToken genera un token simple (en producción usar JWT)
func generateToken() string {
	return time.Now().Format("20060102150405") + "-token"
}

// GetUsuarioFromToken obtiene el usuario desde el token
func GetUsuarioFromToken(token string) (*models.Usuario, bool) {
	usuario, exists := tokens[token]
	if !exists {
		return nil, false
	}
	// Actualizar usuario desde la base de datos para asegurar datos actualizados
	usuarioActualizado, err := database.ObtenerUsuarioPorID(usuario.ID)
	if err != nil {
		return usuario, true // Retornar el usuario del token si falla la actualización
	}
	return usuarioActualizado, true
}

// FinalizarGuardiaHandler finaliza la guardia de un oficial
func FinalizarGuardiaHandler(w http.ResponseWriter, r *http.Request) {
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

	// Finalizar guardia activa
	err := database.FinalizarGuardia(usuario.ID)
	if err != nil {
		http.Error(w, "Error al finalizar guardia", http.StatusInternalServerError)
		return
	}

	// Actualizar usuario
	usuario.EnGuardia = false
	usuario.Latitud = 0
	usuario.Longitud = 0
	database.ActualizarUsuario(usuario)

	response := map[string]interface{}{
		"mensaje": "Guardia finalizada correctamente",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

