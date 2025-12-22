package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Almacenamiento temporal de tokens de master en memoria
var masterTokens = make(map[string]*models.UsuarioMaster)

// RegistroMasterHandler maneja el registro público de usuarios master
func RegistroMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegistroMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if req.Nombre == "" {
		http.Error(w, "El nombre es obligatorio", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "El email es obligatorio", http.StatusBadRequest)
		return
	}

	if req.Usuario == "" {
		http.Error(w, "El usuario es obligatorio", http.StatusBadRequest)
		return
	}

	if req.Contraseña == "" || len(req.Contraseña) < 6 {
		http.Error(w, "La contraseña debe tener al menos 6 caracteres", http.StatusBadRequest)
		return
	}

	// Verificar que el usuario no exista
	_, err := database.ObtenerUsuarioMasterPorUsuario(req.Usuario)
	if err == nil {
		http.Error(w, "El usuario ya está registrado", http.StatusConflict)
		return
	}

	// Verificar que el email no exista
	_, err = database.ObtenerUsuarioMasterPorEmail(req.Email)
	if err == nil {
		http.Error(w, "El email ya está registrado", http.StatusConflict)
		return
	}

	// Hashear contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contraseña), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al procesar la contraseña", http.StatusInternalServerError)
		return
	}

	// Crear usuario master
	master := models.UsuarioMaster{
		Nombre:     req.Nombre,
		Email:      req.Email,
		Usuario:    req.Usuario,
		Contraseña: string(hashedPassword),
		Rol:        "master",
		Activo:     true,
	}

	if err := database.CrearUsuarioMaster(&master); err != nil {
		http.Error(w, "Error al crear el usuario master: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// No retornar la contraseña
	master.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje": "Usuario master registrado correctamente",
		"master":  master,
	})
}

// LoginMasterHandler maneja el login de usuarios master
func LoginMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginMasterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Buscar master por usuario
	master, err := database.ObtenerUsuarioMasterPorUsuario(req.Usuario)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	// Verificar contraseña con bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(master.Contraseña), []byte(req.Contraseña))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario o contraseña incorrectos",
		})
		return
	}

	// Verificar que esté activo
	if !master.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuario inactivo",
		})
		return
	}

	// Actualizar último acceso
	database.ActualizarUltimoAccesoMaster(master.ID)

	// Generar token
	token := generateMasterToken()
	masterTokens[token] = master

	// No retornar la contraseña
	masterResp := *master
	masterResp.Contraseña = ""

	response := models.LoginMasterResponse{
		Token:   token,
		Master:  masterResp,
		Mensaje: "Login exitoso",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// generateMasterToken genera un token simple para masters (en producción usar JWT)
func generateMasterToken() string {
	return time.Now().Format("20060102150405") + "-master-token"
}

// GetMasterFromToken obtiene el master desde el token
func GetMasterFromToken(token string) (*models.UsuarioMaster, bool) {
	master, exists := masterTokens[token]
	if !exists {
		return nil, false
	}
	// Actualizar master desde la base de datos para asegurar datos actualizados
	masterActualizado, err := database.ObtenerUsuarioMasterPorID(master.ID)
	if err != nil {
		return master, true // Retornar el master del token si falla la actualización
	}
	return masterActualizado, true
}

// VerificarMasterHandler verifica si el token es válido y retorna información del master
func VerificarMasterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Token no proporcionado", http.StatusUnauthorized)
		return
	}

	master, ok := GetMasterFromToken(token)
	if !ok {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	// No retornar la contraseña
	masterResp := *master
	masterResp.Contraseña = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(masterResp)
}

