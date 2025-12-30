package handlers

import (
	"encoding/json"
	"gepn/database"
	"gepn/models"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// esAdmin verifica si el usuario es admin
func esAdmin(master *models.UsuarioMaster) bool {
	if master == nil {
		return false
	}
	// Admin tiene todos los permisos o es el usuario "admin"
	return master.Usuario == "admin" || len(master.Permisos) == len(ModulosDisponibles)
}

// tieneAccesoEstado verifica si el usuario tiene acceso a un estado/región
func tieneAccesoEstado(master *models.UsuarioMaster, estado string) bool {
	if esAdmin(master) {
		return true // Admin ve todo
	}
	return master.Estado == estado // RRHH regional solo su estado
}

// CentrosHandler maneja GET y POST para centros
func CentrosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CrearCentroHandler(w, r)
	} else if r.Method == http.MethodGet {
		ListarCentrosHandler(w, r)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// CrearCentroHandler crea un nuevo centro de coordinación (solo admin)
func CrearCentroHandler(w http.ResponseWriter, r *http.Request) {

	// Verificar autenticación y que sea admin
	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	if !esAdmin(master) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Solo el administrador puede crear centros",
		})
		return
	}

	var centro models.Centro
	if err := json.NewDecoder(r.Body).Decode(&centro); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if centro.Nombre == "" || centro.Codigo == "" || centro.Estado == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Nombre, código y estado son obligatorios",
		})
		return
	}

	centro.ID = primitive.NewObjectID()
	centro.FechaCreacion = time.Now()
	centro.Activo = true

	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("centros")

	// Verificar código único
	var existente models.Centro
	err := collection.FindOne(ctx, bson.M{"codigo": centro.Codigo}).Decode(&existente)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ya existe un centro con ese código",
		})
		return
	}

	_, err = collection.InsertOne(ctx, centro)
	if err != nil {
		log.Printf("Error al crear centro: %v", err)
		http.Error(w, "Error al crear centro", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(centro)
}

// ListarCentrosHandler lista centros (admin ve todos, RRHH regional solo su estado)
func ListarCentrosHandler(w http.ResponseWriter, r *http.Request) {

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("centros")

	filter := bson.M{}
	// Si no es admin, filtrar por estado
	if !esAdmin(master) && master.Estado != "" {
		filter["estado"] = master.Estado
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error al listar centros: %v", err)
		http.Error(w, "Error al listar centros", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var centros []models.Centro
	if err := cursor.All(ctx, &centros); err != nil {
		log.Printf("Error al decodificar centros: %v", err)
		http.Error(w, "Error al procesar centros", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(centros)
}

// EstacionesHandler maneja GET y POST para estaciones
func EstacionesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CrearEstacionHandler(w, r)
	} else if r.Method == http.MethodGet {
		ListarEstacionesHandler(w, r)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// CrearEstacionHandler crea una nueva estación (solo admin)
func CrearEstacionHandler(w http.ResponseWriter, r *http.Request) {

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	if !esAdmin(master) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Solo el administrador puede crear estaciones",
		})
		return
	}

	var estacion models.Estacion
	if err := json.NewDecoder(r.Body).Decode(&estacion); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if estacion.Nombre == "" || estacion.Codigo == "" || estacion.Estado == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Nombre, código y estado son obligatorios",
		})
		return
	}

	if estacion.CentroID.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Debe especificar un centro",
		})
		return
	}

	estacion.ID = primitive.NewObjectID()
	estacion.FechaCreacion = time.Now()
	estacion.Activa = true
	estacion.FuncionariosAsignados = []models.FuncionarioAsignado{}

	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("estaciones")

	// Verificar código único
	var existente models.Estacion
	err := collection.FindOne(ctx, bson.M{"codigo": estacion.Codigo}).Decode(&existente)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Ya existe una estación con ese código",
		})
		return
	}

	_, err = collection.InsertOne(ctx, estacion)
	if err != nil {
		log.Printf("Error al crear estación: %v", err)
		http.Error(w, "Error al crear estación", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(estacion)
}

// ListarEstacionesHandler lista estaciones (admin ve todas, RRHH regional solo su estado)
func ListarEstacionesHandler(w http.ResponseWriter, r *http.Request) {

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("estaciones")

	filter := bson.M{}
	// Si no es admin, filtrar por estado
	if !esAdmin(master) && master.Estado != "" {
		filter["estado"] = master.Estado
	}

	// Filtro opcional por centro_id
	centroID := r.URL.Query().Get("centro_id")
	if centroID != "" {
		centroObjID, err := primitive.ObjectIDFromHex(centroID)
		if err == nil {
			filter["centro_id"] = centroObjID
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error al listar estaciones: %v", err)
		http.Error(w, "Error al listar estaciones", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var estaciones []models.Estacion
	if err := cursor.All(ctx, &estaciones); err != nil {
		log.Printf("Error al decodificar estaciones: %v", err)
		http.Error(w, "Error al procesar estaciones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estaciones)
}

// AsignarFuncionarioHandler asigna un funcionario a una estación (solo admin)
func AsignarFuncionarioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	if !esAdmin(master) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Solo el administrador puede asignar funcionarios",
		})
		return
	}

	// Obtener estacion_id de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		http.Error(w, "ID de estación requerido", http.StatusBadRequest)
		return
	}

	estacionID, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-2])
	if err != nil {
		http.Error(w, "ID de estación inválido", http.StatusBadRequest)
		return
	}

	var req struct {
		FuncionarioID primitive.ObjectID `json:"funcionario_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Verificar que el funcionario existe
	oficial, err := database.ObtenerOficialPorID(req.FuncionarioID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Funcionario no encontrado",
		})
		return
	}

	// Obtener estación
	ctx, cancel := database.GetContext()
	defer cancel()
	estacionesCollection := database.GetCollection("estaciones")

	var estacion models.Estacion
	err = estacionesCollection.FindOne(ctx, bson.M{"_id": estacionID}).Decode(&estacion)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Estación no encontrada",
		})
		return
	}

	// Crear asignación
	asignacion := models.FuncionarioAsignado{
		FuncionarioID:   oficial.ID,
		Nombre:          oficial.PrimerNombre + " " + oficial.PrimerApellido,
		Credencial:      oficial.Credencial,
		Rango:           oficial.Rango,
		FechaAsignacion: time.Now(),
		Activo:          true,
	}

	// Agregar a la lista de funcionarios asignados
	estacion.FuncionariosAsignados = append(estacion.FuncionariosAsignados, asignacion)

	_, err = estacionesCollection.UpdateOne(
		ctx,
		bson.M{"_id": estacionID},
		bson.M{"$set": bson.M{"funcionarios_asignados": estacion.FuncionariosAsignados}},
	)
	if err != nil {
		log.Printf("Error al asignar funcionario: %v", err)
		http.Error(w, "Error al asignar funcionario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"mensaje":      "Funcionario asignado correctamente",
		"asignacion":   asignacion,
		"estacion_id":  estacionID,
	})
}

// ListarFuncionariosEstacionHandler lista funcionarios de una estación
func ListarFuncionariosEstacionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	// Obtener estacion_id de la URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 {
		http.Error(w, "ID de estación requerido", http.StatusBadRequest)
		return
	}

	estacionID, err := primitive.ObjectIDFromHex(pathParts[len(pathParts)-1])
	if err != nil {
		http.Error(w, "ID de estación inválido", http.StatusBadRequest)
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()
	collection := database.GetCollection("estaciones")

	var estacion models.Estacion
	err = collection.FindOne(ctx, bson.M{"_id": estacionID}).Decode(&estacion)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Estación no encontrada",
		})
		return
	}

	// Verificar acceso al estado
	if !tieneAccesoEstado(master, estacion.Estado) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No tienes acceso a esta estación",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estacion.FuncionariosAsignados)
}

// PartesHandler maneja GET y POST para partes
func PartesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CrearParteHandler(w, r)
	} else if r.Method == http.MethodGet {
		ListarPartesHandler(w, r)
	} else {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// CrearParteHandler crea un parte de servicio
func CrearParteHandler(w http.ResponseWriter, r *http.Request) {

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	var parte models.Parte
	if err := json.NewDecoder(r.Body).Decode(&parte); err != nil {
		http.Error(w, "Error al decodificar la petición", http.StatusBadRequest)
		return
	}

	// Validaciones
	if parte.EstacionID.IsZero() || parte.FuncionarioID.IsZero() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Estación y funcionario son obligatorios",
		})
		return
	}

	// Verificar que la estación existe y el usuario tiene acceso
	ctx, cancel := database.GetContext()
	defer cancel()
	estacionesCollection := database.GetCollection("estaciones")

	var estacion models.Estacion
	err := estacionesCollection.FindOne(ctx, bson.M{"_id": parte.EstacionID}).Decode(&estacion)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Estación no encontrada",
		})
		return
	}

	// Verificar acceso al estado
	if !tieneAccesoEstado(master, estacion.Estado) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No tienes acceso a esta estación",
		})
		return
	}

	// Verificar que el funcionario está asignado a la estación
	funcionarioAsignado := false
	for _, funcAsig := range estacion.FuncionariosAsignados {
		if funcAsig.FuncionarioID == parte.FuncionarioID && funcAsig.Activo {
			funcionarioAsignado = true
			break
		}
	}

	if !funcionarioAsignado {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "El funcionario no está asignado a esta estación",
		})
		return
	}

	// Generar número de parte
	parte.ID = primitive.NewObjectID()
	parte.NumeroParte = "PART-" + parte.ID.Hex()[:8]
	parte.FechaCreacion = time.Now()
	if parte.FechaHora.IsZero() {
		parte.FechaHora = time.Now()
	}
	if parte.Estado == "" {
		parte.Estado = "activo"
	}

	partesCollection := database.GetCollection("partes")
	_, err = partesCollection.InsertOne(ctx, parte)
	if err != nil {
		log.Printf("Error al crear parte: %v", err)
		http.Error(w, "Error al crear parte", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(parte)
}

// ListarPartesHandler lista partes (admin ve todos, RRHH regional solo su estado)
func ListarPartesHandler(w http.ResponseWriter, r *http.Request) {

	master, ok := GetMasterFromRequest(r)
	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Se requiere autenticación",
		})
		return
	}

	ctx, cancel := database.GetContext()
	defer cancel()
	partesCollection := database.GetCollection("partes")
	estacionesCollection := database.GetCollection("estaciones")

	// Si no es admin, necesitamos filtrar por estado
	var estadosPermitidos []string
	if !esAdmin(master) && master.Estado != "" {
		estadosPermitidos = []string{master.Estado}
	} else {
		// Admin: obtener todos los estados de las estaciones
		cursor, _ := estacionesCollection.Find(ctx, bson.M{})
		defer cursor.Close(ctx)
		var estaciones []models.Estacion
		cursor.All(ctx, &estaciones)
		estadosMap := make(map[string]bool)
		for _, est := range estaciones {
			estadosMap[est.Estado] = true
		}
		for estado := range estadosMap {
			estadosPermitidos = append(estadosPermitidos, estado)
		}
	}

	// Obtener IDs de estaciones permitidas
	var estacionesPermitidas []primitive.ObjectID
	if len(estadosPermitidos) > 0 {
		cursor, _ := estacionesCollection.Find(ctx, bson.M{"estado": bson.M{"$in": estadosPermitidos}})
		defer cursor.Close(ctx)
		var estaciones []models.Estacion
		cursor.All(ctx, &estaciones)
		for _, est := range estaciones {
			estacionesPermitidas = append(estacionesPermitidas, est.ID)
		}
	}

	filter := bson.M{}
	if len(estacionesPermitidas) > 0 {
		filter["estacion_id"] = bson.M{"$in": estacionesPermitidas}
	}

	// Filtros opcionales
	estacionID := r.URL.Query().Get("estacion_id")
	if estacionID != "" {
		estObjID, err := primitive.ObjectIDFromHex(estacionID)
		if err == nil {
			filter["estacion_id"] = estObjID
		}
	}

	funcionarioID := r.URL.Query().Get("funcionario_id")
	if funcionarioID != "" {
		funcObjID, err := primitive.ObjectIDFromHex(funcionarioID)
		if err == nil {
			filter["funcionario_id"] = funcObjID
		}
	}

	tipoParte := r.URL.Query().Get("tipo_parte")
	if tipoParte != "" {
		filter["tipo_parte"] = tipoParte
	}

	opts := options.Find().SetSort(bson.D{{Key: "fecha_hora", Value: -1}})
	cursor, err := partesCollection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("Error al listar partes: %v", err)
		http.Error(w, "Error al listar partes", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var partes []models.Parte
	if err := cursor.All(ctx, &partes); err != nil {
		log.Printf("Error al decodificar partes: %v", err)
		http.Error(w, "Error al procesar partes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(partes)
}

