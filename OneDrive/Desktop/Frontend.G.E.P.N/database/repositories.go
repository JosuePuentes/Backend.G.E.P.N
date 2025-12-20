package database

import (
	"gepn/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Usuarios
func CrearUsuario(usuario *models.Usuario) error {
	usuario.FechaCreacion = time.Now()
	collection := GetCollection("usuarios")
	_, err := collection.InsertOne(Ctx, usuario)
	return err
}

func ObtenerUsuarioPorCredencial(credencial string) (*models.Usuario, error) {
	collection := GetCollection("usuarios")
	var usuario models.Usuario
	err := collection.FindOne(Ctx, bson.M{"credencial": credencial}).Decode(&usuario)
	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

func ObtenerUsuarioPorID(id primitive.ObjectID) (*models.Usuario, error) {
	collection := GetCollection("usuarios")
	var usuario models.Usuario
	err := collection.FindOne(Ctx, bson.M{"_id": id}).Decode(&usuario)
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
	
	// 4. Verificar que las demás colecciones estén listas (se crearán automáticamente al usar)
	// detenidos, minutas, busquedas, panico - se crearán cuando se inserten datos
	
	log.Println("✅ Inicialización de base de datos completada")
	log.Println("📋 Colecciones disponibles: usuarios, detenidos, minutas, busquedas, mas_buscados, panico")
	return nil
}

