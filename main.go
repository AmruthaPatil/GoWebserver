package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"

    "github.com/IBM/sarama"
    "github.com/gorilla/mux"
    "github.com/go-redis/redis/v8"
)

var kafkaProducer sarama.SyncProducer
var redisClient *redis.Client
var ctx = context.Background()

func initKafka() {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true

    var err error
    kafkaProducer, err = sarama.NewSyncProducer([]string{"kafka:9092"}, config)
    if err != nil {
        log.Fatalf("Failed to start Sarama producer: %v", err)
    }
}

func initRedis() {
    redisClient = redis.NewClient(&redis.Options{
        Addr:     "redis:6379", // Use the Docker service name as the address
        Password: "",          // No password set
        DB:       0,           // Default DB
    })

    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
}

func postDataHandler(w http.ResponseWriter, r *http.Request) {
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    message := &sarama.ProducerMessage{
        Topic: "topic-name",
        Value: sarama.StringEncoder(body),
    }
    _, _, err = kafkaProducer.SendMessage(message)
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to send message to Kafka: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Data sent to Kafka successfully"))
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
    data, err := redisClient.Get(ctx, "dataKey").Result()
    if err != nil {
        http.Error(w, fmt.Sprintf("Failed to get data from Redis: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(data))
}

func consumeMessages() {
    config := sarama.NewConfig()
    config.Consumer.Return.Errors = true

    consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
    if err != nil {
        log.Fatalf("Failed to start Kafka consumer: %v", err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition("topic-name", 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatalf("Failed to get partition consumer: %v", err)
    }
    defer partitionConsumer.Close()

    for msg := range partitionConsumer.Messages() {
        log.Printf("Message received from Kafka: %s", string(msg.Value))
        err := redisClient.Set(ctx, "dataKey", msg.Value, 0).Err()
        if err != nil {
            log.Printf("Failed to save data to Redis: %v", err)
        }
    }
}

func main() {
    initKafka()
    initRedis()

    go consumeMessages()  // Start Kafka consumer as a goroutine

    router := mux.NewRouter()
    router.HandleFunc("/data", postDataHandler).Methods("POST")
    router.HandleFunc("/retrieve", getDataHandler).Methods("GET")

    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}
