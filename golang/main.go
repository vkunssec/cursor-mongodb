package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

// Estrutura que define as configurações de conexão com MongoDB
type MongoConfig struct {
	URL       string      // URL de conexão com MongoDB
	AppName   string      // Nome da aplicação
	DebugMode bool        // Modo de debug
	Log       slog.Logger // Logger para registrar eventos
}

// Função que implementa paginação usando cursor baseado em ObjectID
func paginateWithCursor(ctx context.Context, lastID *primitive.ObjectID, limit int) ([]bson.M, error) {
	collection := connect(ctx)

	// Cria filtro para buscar documentos após o último ID (se fornecido)
	filter := bson.D{}
	if lastID != nil {
		filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastID}}}}
	}

	// Configura opções de busca: limite e ordenação por _id
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "_id", Value: 1}})

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Função que estabelece conexão com MongoDB e configura monitoramento
func connect(ctx context.Context) *mongo.Collection {
	cfg := MongoConfig{
		URL:     os.Getenv("MONGODB_URI"),
		AppName: "sample_mflix",
	}

	options := options.Client().ApplyURI(cfg.URL)
	options.SetAppName(cfg.AppName)

	// Configura monitor para logging de comandos MongoDB
	monitor := &event.CommandMonitor{
		Started: func(_ context.Context, e *event.CommandStartedEvent) {
			// Registra comandos enviados ao MongoDB (exceto endSessions e ping)
			if e.CommandName != "endSessions" && e.CommandName != "ping" {
				var decoded map[string]interface{}
				if err := json.Unmarshal([]byte(e.Command.String()), &decoded); err != nil {
					log.Printf("Error decoding command: %v", err)
					return
				}

				formatted, err := json.MarshalIndent(decoded, "", "  ")
				if err != nil {
					log.Printf("Error formatting command: %v", err)
					return
				}

				log.Println(string(formatted))
			}
		},
		Succeeded: func(_ context.Context, e *event.CommandSucceededEvent) {
			// Registra respostas bem-sucedidas do MongoDB
			if e.CommandName != "endSessions" && e.CommandName != "ping" {
				var decoded map[string]interface{}
				if err := json.Unmarshal([]byte(e.Reply.String()), &decoded); err != nil {
					log.Printf("Error decoding reply: %v", err)
					return
				}

				formatted, err := json.MarshalIndent(decoded, "", "  ")
				if err != nil {
					log.Printf("Error formatting reply: %v", err)
					return
				}

				log.Println(string(formatted))
			}
		},
		Failed: func(context.Context, *event.CommandFailedEvent) {},
	}
	options.SetMonitor(monitor)

	clt, err := mongo.Connect(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	collection := clt.Database(os.Getenv("MONGODB_DATABASE")).Collection(os.Getenv("MONGODB_COLLECTION"))

	client = clt
	return collection
}

func main() {
	ctx := context.Background()

	// Carrega variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	// Primeira chamada: busca os primeiros 10 documentos
	results, err := paginateWithCursor(ctx, nil, 10)
	if err != nil {
		log.Fatal(err)
	}

	// Exibe resultados formatados em JSON
	r, _ := json.MarshalIndent(&results, "", "  ")
	log.Println(string(r))

	// Segunda chamada: busca próximos 10 documentos após o último ID
	lastID := results[len(results)-1]["_id"].(primitive.ObjectID)
	results, err = paginateWithCursor(ctx, &lastID, 10)
	if err != nil {
		log.Fatal(err)
	}

	r, _ = json.MarshalIndent(&results, "", "  ")
	log.Println(string(r))

	client.Disconnect(ctx)
	os.Exit(0)
}
