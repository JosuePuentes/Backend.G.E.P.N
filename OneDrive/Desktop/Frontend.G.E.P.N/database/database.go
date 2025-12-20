package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	DB       *mongo.Database
	Ctx      context.Context
	Cancel   context.CancelFunc
)

// Connect establece la conexión con MongoDB
func Connect() error {
	// Obtener la URI de conexión desde variable de entorno
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		// Si no está configurada, usar la URI proporcionada (reemplazar <db_password> con la contraseña real)
		// NOTA: En producción, siempre usa variables de entorno para mayor seguridad
		mongoURI = "mongodb+srv://Drocolven2019:Drocolven2019@drocolven2019.eof9ilx.mongodb.net/?appName=Drocolven2019"
		log.Println("⚠️  Usando URI por defecto. Configura MONGODB_URI en variables de entorno para mayor seguridad")
	}

	// Crear contexto con timeout
	Ctx, Cancel = context.WithTimeout(context.Background(), 10*time.Second)

	// Configurar opciones del cliente
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Conectar a MongoDB
	var err error
	Client, err = mongo.Connect(Ctx, clientOptions)
	if err != nil {
		return err
	}

	// Verificar la conexión
	err = Client.Ping(Ctx, nil)
	if err != nil {
		return err
	}

	// Obtener nombre de la base de datos desde variable de entorno
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "gepn" // Nombre por defecto
	}

	DB = Client.Database(dbName)

	log.Println("✅ Conectado a MongoDB exitosamente")
	return nil
}

// Disconnect cierra la conexión con MongoDB
func Disconnect() {
	if Cancel != nil {
		Cancel()
	}
	if Client != nil {
		if err := Client.Disconnect(Ctx); err != nil {
			log.Printf("Error al desconectar MongoDB: %v", err)
		} else {
			log.Println("✅ Desconectado de MongoDB")
		}
	}
}

// GetCollection retorna una colección de MongoDB
func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}

