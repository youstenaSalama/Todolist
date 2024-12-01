package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)


type Task struct { 
    ID        string `json:"id,omitempty" bson:"_id,omitempty"`
    Title     string `json:"title" bson:"title"`
    Completed bool   `json:"completed" bson:"completed"`
}

var client *mongo.Client
var taskCollection *mongo.Collection


func initMongo() {
    var err error
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017") 
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
    taskCollection = client.Database("todoApp").Collection("tasks")
}


func getTasks(c *gin.Context) {
    var tasks []Task
    cursor, err := taskCollection.Find(context.TODO(), bson.D{{}})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    for cursor.Next(context.TODO()) {
        var task Task
        err := cursor.Decode(&task)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        tasks = append(tasks, task)
    }

    if err := cursor.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    }
    cursor.Close(context.TODO())

    c.JSON(http.StatusOK, tasks)
}

func createTask(c *gin.Context) {
    var newTask Task
    if err := c.BindJSON(&newTask); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := taskCollection.InsertOne(context.TODO(), newTask)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
 
func updateTask(c *gin.Context) {
    id := c.Param("id")
    var updatedTask Task
    if err := c.BindJSON(&updatedTask); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    filter := bson.M{"_id": id}
    update := bson.M{"$set": updatedTask}

    _, err := taskCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func deleteTask(c *gin.Context) {
    id := c.Param("id")
    _, err := taskCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func main() {
    
    initMongo()

   
    router := gin.Default()

 
    router.GET("/tasks", getTasks)
    router.POST("/tasks", createTask)
    router.PUT("/tasks/:id", updateTask)
    router.DELETE("/tasks/:id", deleteTask)

   
    router.Run(":8080")
}
