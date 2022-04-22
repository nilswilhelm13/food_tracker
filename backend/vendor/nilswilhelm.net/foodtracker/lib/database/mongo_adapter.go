package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

type DBAdapter interface {
	Insert(data interface{}, databaseName string, collectionName string) (string, error)
	Get(filter bson.D, databaseName string, collectionName string) *mongo.SingleResult
	GetMultiple(filter bson.D, databaseName string, collectionName string) (*mongo.Cursor, error)
	Update(filter bson.D, update bson.D, databaseName string, collectionName string) error
	Delete(filter bson.D, databaseName string, collectionName string) error
}

type MongoAdapter struct {
	client *mongo.Client
}

func NewMongoAdapter(uri string) (MongoAdapter, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)
	return MongoAdapter{client: client}, nil
}

func (m *MongoAdapter) Insert(data interface{}, databaseName string, collectionName string) (string, error) {

	collection := m.client.Database(databaseName).Collection(collectionName)
	insertResult, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return "", err
	}
	stringObjectID := insertResult.InsertedID.(primitive.ObjectID).Hex()
	log.Println(stringObjectID)
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return stringObjectID, nil
}

func (m *MongoAdapter) Get(filter bson.D, databaseName string, collectionName string) *mongo.SingleResult {

	collection := m.client.Database(databaseName).Collection(collectionName)
	return collection.FindOne(context.TODO(), filter)
}

func (m *MongoAdapter) GetMultiple(filter bson.D, databaseName string, collectionName string, options *options.FindOptions) (*mongo.Cursor, error) {

	collection := m.client.Database(databaseName).Collection(collectionName)
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, err
	}
	return cur, nil
}

func (m *MongoAdapter) Update(filter bson.D, update bson.D, databaseName string, collectionName string) error {

	collection := m.client.Database(databaseName).Collection(collectionName)
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoAdapter) Delete(filter bson.D, databaseName string, collectionName string) error {
	collection := m.client.Database(databaseName).Collection(collectionName)
	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoAdapter) Search(filter bson.A, databaseName string, collectionName string) (*mongo.Cursor, error) {

	collection := m.client.Database(databaseName).Collection(collectionName)
	cur, err := collection.Aggregate(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return cur, nil
}
