package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	database    *mongo.Database
	collection  *mongo.Collection
)

type Task struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Completed bool               `json:"completed" bson:"completed"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func init() {
	// Inicializa o cliente MongoDB
	mongoURI := "mongodb+srv://demo1:demo1@demo1.f7x641l.mongodb.net/?retryWrites=true&w=majority&appName=demo1"
	if mongoURI == "" {
		log.Fatal("MONGODB_URI não foi especificado")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Erro ao criar cliente MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	mongoClient = client
	database = client.Database("todoapp")
	collection = database.Collection("tasks")
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CorsMiddleware())
	// Endpoint para adicionar uma nova tarefa
	r.POST("/tasks", func(c *gin.Context) {
		var task Task
		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		task.ID = primitive.NewObjectID()
		task.CreatedAt = time.Now()

		_, err := collection.InsertOne(context.Background(), task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao adicionar a tarefa"})
			return
		}

		c.JSON(http.StatusCreated, task)
	})

	// Endpoint para obter todas as tarefas
	r.GET("/tasks", func(c *gin.Context) {
		var tasks []Task

		cursor, err := collection.Find(context.Background(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar as tarefas"})
			return
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var task Task
			if err := cursor.Decode(&task); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao decodificar a tarefa"})
				return
			}
			tasks = append(tasks, task)
		}

		c.JSON(http.StatusOK, tasks)
	})

	// Endpoint para marcar uma tarefa como concluída
	r.PATCH("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da tarefa inválido"})
			return
		}

		result, err := collection.UpdateOne(
			context.Background(),
			bson.M{"_id": objID},
			bson.M{"$set": bson.M{"completed": true}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar a tarefa"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// Endpoint para deletar uma tarefa
	r.DELETE("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID da tarefa inválido"})
			return
		}

		result, err := collection.DeleteOne(context.Background(), bson.M{"_id": objID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar a tarefa"})
			return
		}

		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tarefa não encontrada"})
			return
		}

		c.Status(http.StatusNoContent)
	})

	// Inicia o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
