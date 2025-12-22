package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Almacenamiento temporal de tokens en memoria (los tokens pueden estar en memoria)
var tokens = make(map[string]*models.Oficial)

// LoginPolicialHandler maneja el login de policiales
// Ahora usa la colección de oficiales con bcrypt
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

	// Buscar oficial por credencial en la colección oficiales
	oficial, err := database.ObtenerOficialPorCredencial(req.Credencial)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Credenciales inválidas",
		})
		return
	}

	// Verificar contraseña con bcrypt
	// Usar contraseña del request (puede venir como PIN o contraseña)
	// La contraseña registrada en RRHH es la que se usa para el login policial
	contraseña := req.PIN
	if contraseña == "" {
		contraseña = req.Contraseña
	}
	if contraseña == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "La contraseña es obligatoria",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(oficial.Contraseña), []byte(contraseña))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Credenciales inválidas",
		})
		return
	}

	// Verificar que esté activo
	if !oficial.Activo {
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
			OficialID:      oficial.ID,
			FechaInicio:    time.Now(),
			LatitudInicio:  req.Latitud,
			LongitudInicio: req.Longitud,
			Activa:         true,
		}
		if err := database.CrearGuardia(&guardia); err != nil {
			// Log error pero continuar con el login
			// En producción, considerar si esto debe fallar el login
		}
	}

	// Generar token simple (en producción usar JWT)
	token := generateToken()
	tokens[token] = oficial

	// Convertir oficial a usuario para la respuesta (compatibilidad)
	usuarioResp := models.Usuario{
		ID:           oficial.ID,
		Credencial:   oficial.Credencial,
		Nombre:       oficial.PrimerNombre + " " + oficial.SegundoNombre + " " + oficial.PrimerApellido + " " + oficial.SegundoApellido,
		Rango:        oficial.Rango,
		Activo:       oficial.Activo,
		FechaCreacion: oficial.FechaRegistro,
		EnGuardia:    true,
		Latitud:      req.Latitud,
		Longitud:     req.Longitud,
	}

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
// Ahora retorna un Usuario basado en el Oficial del token
func GetUsuarioFromToken(token string) (*models.Usuario, bool) {
	oficial, exists := tokens[token]
	if !exists {
		return nil, false
	}
	// Actualizar oficial desde la base de datos para asegurar datos actualizados
	oficialActualizado, err := database.ObtenerOficialPorID(oficial.ID)
	if err != nil {
		// Si falla, convertir el oficial del token a usuario
		usuario := &models.Usuario{
			ID:            oficial.ID,
			Credencial:    oficial.Credencial,
			Nombre:        oficial.PrimerNombre + " " + oficial.SegundoNombre + " " + oficial.PrimerApellido + " " + oficial.SegundoApellido,
			Rango:         oficial.Rango,
			Activo:        oficial.Activo,
			FechaCreacion: oficial.FechaRegistro,
		}
		return usuario, true
	}
	// Convertir oficial a usuario
	usuario := &models.Usuario{
		ID:            oficialActualizado.ID,
		Credencial:    oficialActualizado.Credencial,
		Nombre:        oficialActualizado.PrimerNombre + " " + oficialActualizado.SegundoNombre + " " + oficialActualizado.PrimerApellido + " " + oficialActualizado.SegundoApellido,
		Rango:         oficialActualizado.Rango,
		Activo:        oficialActualizado.Activo,
		FechaCreacion: oficialActualizado.FechaRegistro,
	}
	return usuario, true
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

	// La guardia se finaliza en la base de datos, no necesitamos actualizar usuario aquí
	// ya que ahora usamos oficiales

	response := map[string]interface{}{
		"mensaje": "Guardia finalizada correctamente",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

