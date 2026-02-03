package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"gepn/database"
	"gepn/models"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Almacenamiento de tokens de patrullaje (en memoria)
var patrullajeTokens = make(map[string]*models.Oficial)

// LoginPatrullajeHandler maneja el login con credencial y PIN
func LoginPatrullajeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginPatrullajeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al decodificar la petición",
		})
		return
	}

	// Validar campos obligatorios
	if req.Credencial == "" || req.PIN == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Credencial y PIN son obligatorios",
		})
		return
	}

	// Validar formato de PIN (6 dígitos)
	if err := ValidarPIN(req.PIN); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Buscar funcionario por credencial
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("oficiales")

	var oficial models.Oficial
	err := collection.FindOne(ctx, bson.M{"credencial": req.Credencial}).Decode(&oficial)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Credencial o PIN incorrectos",
		})
		return
	}

	// Verificar que esté activo
	if !oficial.Activo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Funcionario inactivo",
		})
		return
	}

	// Verificar que tenga PIN configurado
	if oficial.PIN == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Funcionario no tiene PIN configurado para patrullaje",
		})
		return
	}

	// Verificar PIN con bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(oficial.PIN), []byte(req.PIN))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Credencial o PIN incorrectos",
		})
		return
	}

	// Generar token
	token := generatePatrullajeToken()
	patrullajeTokens[token] = &oficial

	// Respuesta exitosa
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Login exitoso",
		"data": map[string]interface{}{
			"id":         oficial.ID.Hex(),
			"nombre":     oficial.PrimerNombre,
			"apellido":   oficial.PrimerApellido,
			"credencial": oficial.Credencial,
			"rango":      oficial.Rango,
			"unidad":     "Patrullaje",
		},
		"token": token,
	})
}

// IniciarPatrullajeHandler inicia un nuevo patrullaje
func IniciarPatrullajeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener funcionario del token
	token := r.Header.Get("Authorization")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token requerido",
		})
		return
	}

	// Remover "Bearer " del token si existe
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	oficial, ok := patrullajeTokens[token]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token inválido",
		})
		return
	}

	var req models.IniciarPatrullajeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al decodificar la petición",
		})
		return
	}

	// Validar coordenadas
	if err := ValidarCoordenadas(req.Latitud, req.Longitud); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Verificar que no tenga un patrullaje activo
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("patrullajes")

	var patrullajeExistente models.Patrullaje
	err := collection.FindOne(ctx, bson.M{
		"funcionario_id": oficial.ID,
		"activo":         true,
	}).Decode(&patrullajeExistente)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Ya tienes un patrullaje activo",
		})
		return
	}

	// Asignar color
	color := AsignarColor(collection)

	// Crear patrullaje
	ahora := time.Now()
	patrullaje := models.Patrullaje{
		FuncionarioID:       oficial.ID,
		Credencial:          oficial.Credencial,
		Nombre:              oficial.PrimerNombre,
		Apellido:            oficial.PrimerApellido,
		Rango:               oficial.Rango,
		Unidad:              "Patrullaje",
		Latitud:             req.Latitud,
		Longitud:            req.Longitud,
		Color:               color,
		FechaInicio:         ahora,
		UltimaActualizacion: ahora,
		Activo:              true,
	}

	result, err := collection.InsertOne(ctx, patrullaje)
	if err != nil {
		log.Printf("Error al crear patrullaje: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al crear patrullaje",
		})
		return
	}

	patrullajeID := result.InsertedID.(primitive.ObjectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Patrullaje iniciado correctamente",
		"data": map[string]interface{}{
			"patrullajeId": patrullajeID.Hex(),
			"nombre":       oficial.PrimerNombre + " " + oficial.PrimerApellido,
			"credencial":   oficial.Credencial,
			"color":        color,
			"fecha_inicio": ahora,
		},
	})
}

// ActualizarUbicacionHandler actualiza la ubicación del patrullaje
func ActualizarUbicacionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener funcionario del token
	token := r.Header.Get("Authorization")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token requerido",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	oficial, ok := patrullajeTokens[token]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token inválido",
		})
		return
	}

	var req models.ActualizarUbicacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al decodificar la petición",
		})
		return
	}

	// Validar coordenadas
	if err := ValidarCoordenadas(req.Latitud, req.Longitud); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Convertir ID
	patrullajeID, err := primitive.ObjectIDFromHex(req.PatrullajeID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID de patrullaje inválido",
		})
		return
	}

	// Actualizar ubicación
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("patrullajes")

	ahora := time.Now()
	update := bson.M{
		"$set": bson.M{
			"latitud":              req.Latitud,
			"longitud":             req.Longitud,
			"ultima_actualizacion": ahora,
		},
	}

	result, err := collection.UpdateOne(ctx, bson.M{
		"_id":            patrullajeID,
		"funcionario_id": oficial.ID,
		"activo":         true,
	}, update)

	if err != nil || result.MatchedCount == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Patrullaje no encontrado o no autorizado",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Ubicación actualizada",
		"data": map[string]interface{}{
			"latitud":              req.Latitud,
			"longitud":             req.Longitud,
			"ultima_actualizacion": ahora,
		},
	})
}

// ObtenerPatrullajesActivosHandler obtiene todos los patrullajes activos
func ObtenerPatrullajesActivosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar token
	token := r.Header.Get("Authorization")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token requerido",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	_, ok := patrullajeTokens[token]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token inválido",
		})
		return
	}

	// Obtener patrullajes activos
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("patrullajes")

	opts := options.Find().SetSort(bson.D{{Key: "ultima_actualizacion", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"activo": true}, opts)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al obtener patrullajes",
		})
		return
	}
	defer cursor.Close(ctx)

	var patrullajes []map[string]interface{}
	for cursor.Next(ctx) {
		var patrullaje models.Patrullaje
		if err := cursor.Decode(&patrullaje); err != nil {
			continue
		}

		patrullajes = append(patrullajes, map[string]interface{}{
			"id":                   patrullaje.ID.Hex(),
			"credencial":           patrullaje.Credencial,
			"nombre":               patrullaje.Nombre,
			"apellido":             patrullaje.Apellido,
			"rango":                patrullaje.Rango,
			"unidad":               patrullaje.Unidad,
			"latitud":              patrullaje.Latitud,
			"longitud":             patrullaje.Longitud,
			"color":                patrullaje.Color,
			"ultima_actualizacion": patrullaje.UltimaActualizacion,
		})
	}

	if patrullajes == nil {
		patrullajes = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    patrullajes,
	})
}

// FinalizarPatrullajeHandler finaliza un patrullaje
func FinalizarPatrullajeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener funcionario del token
	token := r.Header.Get("Authorization")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token requerido",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	oficial, ok := patrullajeTokens[token]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token inválido",
		})
		return
	}

	var req models.FinalizarPatrullajeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al decodificar la petición",
		})
		return
	}

	// Convertir ID
	patrullajeID, err := primitive.ObjectIDFromHex(req.PatrullajeID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "ID de patrullaje inválido",
		})
		return
	}

	// Buscar patrullaje
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("patrullajes")

	var patrullaje models.Patrullaje
	err = collection.FindOne(ctx, bson.M{
		"_id":            patrullajeID,
		"funcionario_id": oficial.ID,
		"activo":         true,
	}).Decode(&patrullaje)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Patrullaje no encontrado o no autorizado",
		})
		return
	}

	// Finalizar patrullaje
	ahora := time.Now()
	duracionMinutos := int(ahora.Sub(patrullaje.FechaInicio).Minutes())

	update := bson.M{
		"$set": bson.M{
			"activo":    false,
			"fecha_fin": ahora,
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": patrullajeID}, update)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al finalizar patrullaje",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Patrullaje finalizado correctamente",
		"data": map[string]interface{}{
			"patrullajeId":      patrullajeID.Hex(),
			"fecha_inicio":      patrullaje.FechaInicio,
			"fecha_fin":         ahora,
			"duracion_minutos":  duracionMinutos,
		},
	})
}

// HistorialPatrullajesHandler obtiene el historial de patrullajes
func HistorialPatrullajesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Verificar token
	token := r.Header.Get("Authorization")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token requerido",
		})
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	oficial, ok := patrullajeTokens[token]
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Token inválido",
		})
		return
	}

	// Obtener historial del funcionario
	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("patrullajes")

	opts := options.Find().SetSort(bson.D{{Key: "fecha_inicio", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"funcionario_id": oficial.ID}, opts)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al obtener historial",
		})
		return
	}
	defer cursor.Close(ctx)

	var historial []map[string]interface{}
	for cursor.Next(ctx) {
		var patrullaje models.Patrullaje
		if err := cursor.Decode(&patrullaje); err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":           patrullaje.ID.Hex(),
			"credencial":   patrullaje.Credencial,
			"nombre":       patrullaje.Nombre + " " + patrullaje.Apellido,
			"fecha_inicio": patrullaje.FechaInicio,
			"activo":       patrullaje.Activo,
		}

		if patrullaje.FechaFin != nil {
			item["fecha_fin"] = *patrullaje.FechaFin
			duracionMinutos := int(patrullaje.FechaFin.Sub(patrullaje.FechaInicio).Minutes())
			item["duracion_minutos"] = duracionMinutos
		}

		historial = append(historial, item)
	}

	if historial == nil {
		historial = []map[string]interface{}{}
	}

	total := len(historial)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    historial,
		"total":   total,
	})
}

// Funciones auxiliares

// AsignarColor asigna un color (rojo o azul) basado en el número de patrullajes activos
func AsignarColor(collection *mongo.Collection) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"activo": true})
	if err != nil {
		// En caso de error, asignar azul por defecto
		return "azul"
	}

	if count%2 == 0 {
		return "rojo"
	}
	return "azul"
}

// ValidarCoordenadas valida que las coordenadas sean válidas
func ValidarCoordenadas(lat, lon float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("Latitud inválida (debe estar entre -90 y 90)")
	}
	if lon < -180 || lon > 180 {
		return errors.New("Longitud inválida (debe estar entre -180 y 180)")
	}
	return nil
}

// CrearUsuarioPruebaPatrullajeHandler crea un oficial de prueba para patrullaje si no existe.
// GET /api/patrullaje/crear-usuario-prueba - devuelve credencial y PIN para poder entrar.
func CrearUsuarioPruebaPatrullajeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	const credencialPrueba = "PATRULLA-TEST"
	const pinPrueba = "123456"

	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("oficiales")

	var existente models.Oficial
	err := collection.FindOne(ctx, bson.M{"credencial": credencialPrueba}).Decode(&existente)
	if err == nil {
		// Ya existe, devolver las credenciales (sin cambiar el PIN)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Usuario de prueba ya existe. Usa estas credenciales para entrar.",
			"credencial": credencialPrueba,
			"pin":        pinPrueba,
		})
		return
	}
	if err != mongo.ErrNoDocuments {
		log.Printf("Error al buscar oficial de prueba: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al verificar usuario de prueba",
		})
		return
	}

	// Crear oficial de prueba
	hashedPIN, _ := bcrypt.GenerateFromPassword([]byte(pinPrueba), bcrypt.DefaultCost)
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("Prueba123"), bcrypt.DefaultCost)

	oficial := models.Oficial{
		ID:               primitive.NewObjectID(),
		PrimerNombre:     "Oficial",
		SegundoNombre:    "Prueba",
		PrimerApellido:   "Patrullaje",
		SegundoApellido:  "GEPN",
		Cedula:           "V-00000000",
		Contraseña:       string(hashedPass),
		PIN:              string(hashedPIN),
		FechaNacimiento:  "01/01/1990",
		Estatura:         1.70,
		ColorPiel:        "Morena",
		TipoSangre:       "O+",
		CiudadNacimiento: "Caracas",
		Credencial:       credencialPrueba,
		Rango:            "Oficial",
		Destacado:        "Patrullaje",
		FechaGraduacion:  "01/01/2020",
		Antiguedad:       5,
		Estado:           "Distrito Capital",
		Municipio:        "Libertador",
		Parroquia:        "El Valle",
		FotoCara:         "",
		FechaRegistro:    time.Now(),
		Activo:           true,
	}

	if err := database.CrearOficial(&oficial); err != nil {
		log.Printf("Error al crear oficial de prueba: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Error al crear usuario de prueba: " + err.Error(),
		})
		return
	}

	log.Printf("Usuario de prueba patrullaje creado: %s PIN: %s", credencialPrueba, pinPrueba)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Usuario de prueba creado. Usa estas credenciales en el login de patrullaje.",
		"credencial": credencialPrueba,
		"pin":        pinPrueba,
	})
}

// ValidarPIN valida que el PIN tenga el formato correcto
func ValidarPIN(pin string) error {
	if len(pin) != 6 {
		return errors.New("PIN debe tener 6 dígitos")
	}
	matched, _ := regexp.MatchString(`^[0-9]{6}$`, pin)
	if !matched {
		return errors.New("PIN debe contener solo números")
	}
	return nil
}

// generatePatrullajeToken genera un token simple para patrullaje
func generatePatrullajeToken() string {
	return "patrullaje-" + time.Now().Format("20060102150405")
}

