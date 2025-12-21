package database

import (
	"errors"
	"fmt"
	"gepn/models"
	"log"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Usuarios
func CrearUsuario(usuario *models.Usuario) error {
	ctx, cancel := GetContext()
	defer cancel()
	usuario.FechaCreacion = time.Now()
	collection := GetCollection("usuarios")
	_, err := collection.InsertOne(ctx, usuario)
	return err
}

func ObtenerUsuarioPorCredencial(credencial string) (*models.Usuario, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("usuarios")
	var usuario models.Usuario
	err := collection.FindOne(ctx, bson.M{"credencial": credencial}).Decode(&usuario)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func ObtenerUsuarioPorID(id primitive.ObjectID) (*models.Usuario, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("usuarios")
	var usuario models.Usuario
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&usuario)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

// Detenidos
func CrearDetenido(detenido *models.Detenido) error {
	detenido.FechaDetencion = time.Now()
	if detenido.Estado == "" {
		detenido.Estado = "detenido"
	}
	collection := GetCollection("detenidos")
	_, err := collection.InsertOne(Ctx, detenido)
	return err
}

func ListarDetenidos() ([]models.Detenido, error) {
	collection := GetCollection("detenidos")
	cursor, err := collection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var detenidos []models.Detenido
	if err = cursor.All(Ctx, &detenidos); err != nil {
		return nil, err
	}
	return detenidos, nil
}

func ObtenerDetenidoPorID(id primitive.ObjectID) (*models.Detenido, error) {
	collection := GetCollection("detenidos")
	var detenido models.Detenido
	err := collection.FindOne(Ctx, bson.M{"_id": id}).Decode(&detenido)
	if err != nil {
		return nil, err
	}
	return &detenido, nil
}

// Minutas
func CrearMinuta(minuta *models.Minuta) error {
	minuta.FechaCreacion = time.Now()
	collection := GetCollection("minutas")
	_, err := collection.InsertOne(Ctx, minuta)
	return err
}

func ListarMinutas() ([]models.Minuta, error) {
	collection := GetCollection("minutas")
	cursor, err := collection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var minutas []models.Minuta
	if err = cursor.All(Ctx, &minutas); err != nil {
		return nil, err
	}
	return minutas, nil
}

func ObtenerMinutaPorID(id primitive.ObjectID) (*models.Minuta, error) {
	collection := GetCollection("minutas")
	var minuta models.Minuta
	err := collection.FindOne(Ctx, bson.M{"_id": id}).Decode(&minuta)
	if err != nil {
		return nil, err
	}
	return &minuta, nil
}

// Búsquedas
func CrearBusqueda(busqueda *models.BusquedaCedula) error {
	busqueda.FechaBusqueda = time.Now()
	collection := GetCollection("busquedas")
	_, err := collection.InsertOne(Ctx, busqueda)
	return err
}

func BuscarMasBuscadoPorCedula(cedula string) (*models.MasBuscado, error) {
	collection := GetCollection("mas_buscados")
	var masBuscado models.MasBuscado
	err := collection.FindOne(Ctx, bson.M{"cedula": cedula}).Decode(&masBuscado)
	if err != nil {
		return nil, err
	}
	
	// Incrementar contador de veces buscado
	collection.UpdateOne(Ctx, 
		bson.M{"cedula": cedula},
		bson.M{"$inc": bson.M{"veces_buscado": 1}},
	)
	
	return &masBuscado, nil
}

func ListarMasBuscados() ([]models.MasBuscado, error) {
	collection := GetCollection("mas_buscados")
	
	// Ordenar por veces_buscado descendente
	opts := options.Find().SetSort(bson.D{{Key: "veces_buscado", Value: -1}})
	cursor, err := collection.Find(Ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var masBuscados []models.MasBuscado
	if err = cursor.All(Ctx, &masBuscados); err != nil {
		return nil, err
	}
	return masBuscados, nil
}

// Pánico
func CrearAlertaPanico(alerta *models.Panico) error {
	alerta.FechaActivacion = time.Now()
	if alerta.Estado == "" {
		alerta.Estado = "activo"
	}
	collection := GetCollection("panico")
	_, err := collection.InsertOne(Ctx, alerta)
	return err
}

func ListarAlertasPanico() ([]models.Panico, error) {
	collection := GetCollection("panico")
	cursor, err := collection.Find(Ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(Ctx)

	var alertas []models.Panico
	if err = cursor.All(Ctx, &alertas); err != nil {
		return nil, err
	}
	return alertas, nil
}

// Inicializar datos por defecto
func InicializarDatos() error {
	
	// 1. Inicializar colección de usuarios
	collection := GetCollection("usuarios")
	count, err := collection.CountDocuments(Ctx, bson.M{})
	if err != nil {
		// Si la colección no existe, MongoDB la creará automáticamente
		count = 0
	}
	
	if count == 0 {
		usuarios := []interface{}{
			models.Usuario{
				ID:           primitive.NewObjectID(),
				Credencial:   "POL001",
				PIN:          "123456",
				Nombre:       "Juan Pérez",
				Rango:        "Oficial",
				Activo:       true,
				FechaCreacion: time.Now(),
			},
			models.Usuario{
				ID:           primitive.NewObjectID(),
				Credencial:   "POL002",
				PIN:          "654321",
				Nombre:       "María González",
				Rango:        "Sargento",
				Activo:       true,
				FechaCreacion: time.Now(),
			},
		}
		
		_, err = collection.InsertMany(Ctx, usuarios)
		if err != nil {
			return err
		}
		log.Println("✅ Colección 'usuarios' inicializada con 2 usuarios por defecto")
	} else {
		log.Printf("ℹ️  Colección 'usuarios' ya existe con %d usuarios", count)
	}
	
	// 2. Inicializar colección de más buscados
	collection = GetCollection("mas_buscados")
	count, err = collection.CountDocuments(Ctx, bson.M{})
	if err != nil {
		count = 0
	}
	
	if count == 0 {
		masBuscados := []interface{}{
			models.MasBuscado{
				ID:           primitive.NewObjectID(),
				Cedula:       "1234567890",
				Nombre:       "Juan",
				Apellido:     "Delincuente",
				Motivo:       "Robo a mano armada",
				Prioridad:    "alta",
				VecesBuscado: 15,
			},
			models.MasBuscado{
				ID:           primitive.NewObjectID(),
				Cedula:       "0987654321",
				Nombre:       "María",
				Apellido:     "Fugitiva",
				Motivo:       "Homicidio",
				Prioridad:    "alta",
				VecesBuscado: 12,
			},
		}
		
		_, err = collection.InsertMany(Ctx, masBuscados)
		if err != nil {
			return err
		}
		log.Println("✅ Colección 'mas_buscados' inicializada con 2 registros")
	} else {
		log.Printf("ℹ️  Colección 'mas_buscados' ya existe con %d registros", count)
	}
	
	// 3. Crear índices para mejorar el rendimiento
	// Índice único en credencial de usuarios
	usuariosCollection := GetCollection("usuarios")
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "credencial", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	usuariosCollection.Indexes().CreateOne(Ctx, indexModel)
	
	// Índice en cedula de mas_buscados
	masBuscadosCollection := GetCollection("mas_buscados")
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "cedula", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	masBuscadosCollection.Indexes().CreateOne(Ctx, indexModel)
	
	// Índice único en cedula de ciudadanos
	ciudadanosCollection := GetCollection("ciudadanos")
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "cedula", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ciudadanosCollection.Indexes().CreateOne(Ctx, indexModel)
	
	// Índice en ciudadano_id de denuncias
	denunciasCollection := GetCollection("denuncias")
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "ciudadano_id", Value: 1}},
	}
	denunciasCollection.Indexes().CreateOne(Ctx, indexModel)
	
	// 4. Crear índices para oficiales (RRHH)
	oficialesCollection := GetCollection("oficiales")
	// Índice único en credencial
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "credencial", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	oficialesCollection.Indexes().CreateOne(Ctx, indexModel)
	// Índice único en cédula
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "cedula", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	oficialesCollection.Indexes().CreateOne(Ctx, indexModel)
	// Índice en rango
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "rango", Value: 1}},
	}
	oficialesCollection.Indexes().CreateOne(Ctx, indexModel)
	// Índice en estado
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "estado", Value: 1}},
	}
	oficialesCollection.Indexes().CreateOne(Ctx, indexModel)
	
	// 5. Verificar que las demás colecciones estén listas (se crearán automáticamente al usar)
	// detenidos, minutas, busquedas, panico - se crearán cuando se inserten datos
	
	log.Println("✅ Inicialización de base de datos completada")
	log.Println("📋 Colecciones disponibles: usuarios, detenidos, minutas, busquedas, mas_buscados, panico, ciudadanos, denuncias, oficiales")
	return nil
}

// Ciudadanos
func CrearCiudadano(ciudadano *models.Ciudadano) error {
	ciudadano.FechaRegistro = time.Now()
	if !ciudadano.Activo {
		ciudadano.Activo = true
	}
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("ciudadanos")
	_, err := collection.InsertOne(ctx, ciudadano)
	return err
}

func ObtenerCiudadanoPorCedula(cedula string) (*models.Ciudadano, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("ciudadanos")
	var ciudadano models.Ciudadano
	err := collection.FindOne(ctx, bson.M{"cedula": cedula, "activo": true}).Decode(&ciudadano)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &ciudadano, nil
}

func ObtenerCiudadanoPorID(id primitive.ObjectID) (*models.Ciudadano, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("ciudadanos")
	var ciudadano models.Ciudadano
	err := collection.FindOne(ctx, bson.M{"_id": id, "activo": true}).Decode(&ciudadano)
	if err != nil {
		return nil, err
	}
	return &ciudadano, nil
}

// Denuncias
func CrearDenuncia(denuncia *models.Denuncia) error {
	ctx, cancel := GetContext()
	defer cancel()
	denuncia.FechaDenuncia = time.Now()
	if denuncia.Estado == "" {
		denuncia.Estado = "Pendiente"
	}
	collection := GetCollection("denuncias")
	_, err := collection.InsertOne(ctx, denuncia)
	return err
}

func ObtenerDenunciasPorCiudadano(ciudadanoID primitive.ObjectID) ([]models.Denuncia, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("denuncias")
	
	// Ordenar por fecha descendente
	opts := options.Find().SetSort(bson.D{{Key: "fecha_denuncia", Value: -1}})
	cursor, err := collection.Find(ctx, bson.M{"ciudadano_id": ciudadanoID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var denuncias []models.Denuncia
	if err = cursor.All(ctx, &denuncias); err != nil {
		return nil, err
	}
	return denuncias, nil
}

func GenerarNumeroDenuncia() (string, error) {
	ctx, cancel := GetContext()
	defer cancel()
	año := time.Now().Year()
	collection := GetCollection("denuncias")
	
	// Contar denuncias del año actual
	fechaInicio := time.Date(año, 1, 1, 0, 0, 0, 0, time.UTC)
	fechaFin := time.Date(año+1, 1, 1, 0, 0, 0, 0, time.UTC)
	
	count, err := collection.CountDocuments(ctx, bson.M{
		"fecha_denuncia": bson.M{
			"$gte": fechaInicio,
			"$lt":  fechaFin,
		},
	})
	if err != nil {
		return "", err
	}
	
	numero := count + 1
	return fmt.Sprintf("DEN-%d-%04d", año, numero), nil
}

// Guardias
func CrearGuardia(guardia *models.Guardia) error {
	ctx, cancel := GetContext()
	defer cancel()
	guardia.FechaInicio = time.Now()
	guardia.Activa = true
	collection := GetCollection("guardias")
	_, err := collection.InsertOne(ctx, guardia)
	return err
}

func FinalizarGuardia(oficialID primitive.ObjectID) error {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("guardias")
	fechaFin := time.Now()
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"oficial_id": oficialID, "activa": true},
		bson.M{"$set": bson.M{"fecha_fin": fechaFin, "activa": false}},
	)
	return err
}

func ActualizarUsuario(usuario *models.Usuario) error {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("usuarios")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": usuario.ID},
		bson.M{"$set": bson.M{
			"en_guardia": usuario.EnGuardia,
			"latitud":    usuario.Latitud,
			"longitud":   usuario.Longitud,
		}},
	)
	return err
}

func ObtenerOficialesEnGuardiaCercanos(latitud, longitud float64, radioKm float64) ([]models.Usuario, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("usuarios")
	
	// Buscar oficiales en guardia
	var oficiales []models.Usuario
	cursor, err := collection.Find(ctx, bson.M{"en_guardia": true, "activo": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	if err = cursor.All(ctx, &oficiales); err != nil {
		return nil, err
	}
	
	// Filtrar por distancia (cálculo simple de distancia haversine)
	oficialesCercanos := []models.Usuario{}
	for _, oficial := range oficiales {
		if oficial.Latitud != 0 && oficial.Longitud != 0 {
			distancia := calcularDistancia(latitud, longitud, oficial.Latitud, oficial.Longitud)
			if distancia <= radioKm {
				oficialesCercanos = append(oficialesCercanos, oficial)
			}
		}
	}
	
	return oficialesCercanos, nil
}

// calcularDistancia calcula la distancia en km entre dos puntos usando la fórmula de Haversine
func calcularDistancia(lat1, lon1, lat2, lon2 float64) float64 {
	const radioTierra = 6371 // Radio de la Tierra en km
	dLat := (lat2 - lat1) * 3.141592653589793 / 180.0
	dLon := (lon2 - lon1) * 3.141592653589793 / 180.0
	a := 0.5 - math.Cos(dLat)/2 + math.Cos(lat1*3.141592653589793/180.0)*math.Cos(lat2*3.141592653589793/180.0)*(1-math.Cos(dLon))/2
	return radioTierra * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// Oficiales - Funciones para gestión de oficiales (RRHH)
func CrearOficial(oficial *models.Oficial) error {
	ctx, cancel := GetContext()
	defer cancel()
	oficial.FechaRegistro = time.Now()
	if !oficial.Activo {
		oficial.Activo = true
	}
	collection := GetCollection("oficiales")
	_, err := collection.InsertOne(ctx, oficial)
	return err
}

func ObtenerOficialPorCredencial(credencial string) (*models.Oficial, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("oficiales")
	var oficial models.Oficial
	err := collection.FindOne(ctx, bson.M{"credencial": credencial, "activo": true}).Decode(&oficial)
	if err != nil {
		return nil, err
	}
	return &oficial, nil
}

func ObtenerOficialPorID(id primitive.ObjectID) (*models.Oficial, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("oficiales")
	var oficial models.Oficial
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&oficial)
	if err != nil {
		return nil, err
	}
	return &oficial, nil
}

func ObtenerOficialPorCedula(cedula string) (*models.Oficial, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("oficiales")
	var oficial models.Oficial
	err := collection.FindOne(ctx, bson.M{"cedula": cedula}).Decode(&oficial)
	if err != nil {
		return nil, err
	}
	return &oficial, nil
}

func ActualizarOficial(oficial *models.Oficial) error {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("oficiales")
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": oficial.ID}, oficial)
	return err
}

func ListarOficiales(page, limit int, rango, estado string) ([]models.Oficial, int64, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("oficiales")
	
	// Construir filtro
	filter := bson.M{}
	if rango != "" {
		filter["rango"] = rango
	}
	if estado != "" {
		filter["estado"] = estado
	}
	
	// Contar total
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	// Paginación
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)).SetSort(bson.D{{Key: "fecha_registro", Value: -1}})
	
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var oficiales []models.Oficial
	if err = cursor.All(ctx, &oficiales); err != nil {
		return nil, 0, err
	}
	
	return oficiales, total, nil
}

func ObtenerOficialesConAscensosPendientes() ([]models.Oficial, error) {
	ctx, cancel := GetContext()
	defer cancel()
	collection := GetCollection("oficiales")
	
	// Buscar todos los oficiales activos
	cursor, err := collection.Find(ctx, bson.M{"activo": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var oficiales []models.Oficial
	if err = cursor.All(ctx, &oficiales); err != nil {
		return nil, err
	}
	
	// Filtrar por ascensos pendientes (cada 4 años)
	var oficialesConAscenso []models.Oficial
	ahora := time.Now()
	
	for _, oficial := range oficiales {
		if oficial.FechaGraduacion != "" {
			fechaGraduacion, err := time.Parse("2006-01-02", oficial.FechaGraduacion)
			if err == nil {
				antiguedad := ahora.Sub(fechaGraduacion).Hours() / 24 / 365.25
				// Verificar si tiene un ascenso pendiente (cada 4 años)
				ascensosCompletos := int(antiguedad / 4)
				antiguedadActual := oficial.Antiguedad
				if antiguedadActual < float64(ascensosCompletos*4) {
					oficialesConAscenso = append(oficialesConAscenso, oficial)
				}
			}
		}
	}
	
	return oficialesConAscenso, nil
}

// VerificarOficialPorQR ya no se usa, la verificación se hace en el handler
// Esta función se mantiene por compatibilidad pero puede eliminarse

