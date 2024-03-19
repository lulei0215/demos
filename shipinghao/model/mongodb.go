package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDB() {
	// 设置MongoDB连接字符串
	connStr := "mongodb://hehe:JaA5aXDie8myN64w@101.35.227.154:27017/hehe"

	// 设置MongoDB客户端选项
	clientOptions := options.Client().ApplyURI(connStr)

	// 连接到MongoDB
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// 检查连接
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// 指定要插入文档的集合
	collection := client.Database("hehe").Collection("chat")

	// 创建要插入的文档
	document := bson.M{"name": "John", "age": 30, "city": "New York"}

	// 插入文档
	insertResult, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a document: ", insertResult.InsertedID)

}
func InsertDocument(client *mongo.Client, database, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := client.Database(database).Collection(collection)
	result, err := coll.InsertOne(context.Background(), document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// func ConnectToDB(name string) (*mongo.Database, error) {
// 	connStr := "mongodb://hehe:JaA5aXDie8myN64w@101.35.227.154:27017/hehe?authSource=admin"

// 	timeout := 10 * time.Second
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	o := options.Client().ApplyURI(connStr)
// 	o.SetMaxPoolSize(100)
// 	client, err := mongo.Connect(ctx, o)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client.Database(name), nil
// }
