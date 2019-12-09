package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	err error
)

// MongoStore representation
type MongoStore struct {
	db                  *mongo.Client
	crawledCollection   *mongo.Collection
	uncrawledCollection *mongo.Collection
	mongoCtx            context.Context
}

// MongoDoc representation
type MongoDoc struct {
	ID       primitive.ObjectID `bson:"_id"`
	DocKey   string             `bson:"docKey"`
	DocValue []string           `bson:"docValue"`
}

// NewMongoStore constructor
func NewMongoStore() *MongoStore {
	return &MongoStore{
		db:                  nil,
		crawledCollection:   nil,
		uncrawledCollection: nil,
		mongoCtx:            context.Background(),
	}
}

// Init MongoStore with dbName
func (ms *MongoStore) Init(dbName string) error {
	ms.db, err = mongo.Connect(ms.mongoCtx,
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", dbName)))
	if err != nil {
		return err
	}

	log.Println("Pinging database ...")
	if err := ms.db.Ping(ms.mongoCtx, nil); err != nil {
		return err
	}

	ms.crawledCollection = ms.db.Database(dbName).Collection("crawled")

	log.Println("Dropping uncrawled collection ...")
	if err = ms.crawledCollection.Drop(ms.mongoCtx); err != nil {
		return err
	}

	ms.uncrawledCollection = ms.db.Database(dbName).Collection("uncrawled")

	log.Println("Dropping uncrawled collection ...")
	if err = ms.uncrawledCollection.Drop(ms.mongoCtx); err != nil {
		return err
	}

	log.Println("Successfully initialized the mongo database ...")

	return nil
}

/* --------------- Crawled Functions --------------- */

// PutCrawled v(s) into the k
func (ms *MongoStore) PutCrawled(k string, v []string) error {
	query := bson.M{
		"docKey": k,
	}

	var existingDoc MongoDoc

	check := ms.crawledCollection.FindOne(ms.mongoCtx, query)
	log.Println(check.Err())
	if check.Err() != nil {
		doc := MongoDoc{
			ID:       primitive.NewObjectID(),
			DocKey:   k,
			DocValue: v,
		}

		_, err := ms.crawledCollection.InsertOne(ms.mongoCtx, doc)
		if err != nil {
			return err
		}
	} else {
		err := check.Decode(&existingDoc)
		if err != nil {
			return err
		}

		for _, val := range v {
			existingDoc.DocValue = append(existingDoc.DocValue, val)
		}

		if err = ms.UpdateCrawled(k, existingDoc.DocValue); err != nil {
			return err
		}
	}

	return nil
}

// GetCrawled MongoDoc(s) from k
func (ms *MongoStore) GetCrawled(k string) ([]*MongoDoc, error) {
	query := bson.M{
		"docKey": k,
	}

	var docs []*MongoDoc

	cur, err := ms.crawledCollection.Find(ms.mongoCtx, query)

	for cur.Next(ms.mongoCtx) {
		var doc *MongoDoc

		err = cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// GetAllCrawled MongoDoc(s)
func (ms *MongoStore) GetAllCrawled() ([]*MongoDoc, error) {
	cur, err := ms.crawledCollection.Find(ms.mongoCtx, bson.D{})
	if err != nil {
		return nil, err
	}

	var vals []*MongoDoc

	for cur.Next(ms.mongoCtx) {
		var doc *MongoDoc

		err = cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		log.Printf("%v", doc)

		vals = append(vals, doc)
	}

	return vals, nil
}

// UpdateCrawled k with v(s)
func (ms *MongoStore) UpdateCrawled(k string, v []string) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"docKey", k}}
	update := bson.D{{"$set", bson.D{{"docValue", v}}}}

	result := ms.crawledCollection.FindOneAndUpdate(ms.mongoCtx, filter, update, opts)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}

		return err
	}

	return nil
}

/* --------------- Uncrawled Functions --------------- */

// PutUncrawled v(s) into k
func (ms *MongoStore) PutUncrawled(k string, v []string) error {
	query := bson.M{
		"docKey": k,
	}

	var existingDoc MongoDoc

	check := ms.uncrawledCollection.FindOne(ms.mongoCtx, query)
	log.Println(check.Err())
	if check.Err() != nil {
		doc := MongoDoc{
			ID:       primitive.NewObjectID(),
			DocKey:   k,
			DocValue: v,
		}

		_, err := ms.uncrawledCollection.InsertOne(ms.mongoCtx, doc)
		if err != nil {
			return err
		}
	} else {
		err := check.Decode(&existingDoc)
		if err != nil {
			return err
		}

		for _, val := range v {
			existingDoc.DocValue = append(existingDoc.DocValue, val)
		}

		if err = ms.UpdateUncrawled(k, existingDoc.DocValue); err != nil {
			return err
		}
	}

	return nil
}

// GetUncrawled MongoDoc(s) from k
func (ms *MongoStore) GetUncrawled(k string) ([]*MongoDoc, error) {
	query := bson.M{
		"docKey": k,
	}

	var docs []*MongoDoc

	cur, err := ms.uncrawledCollection.Find(ms.mongoCtx, query)

	for cur.Next(ms.mongoCtx) {
		var doc *MongoDoc

		err = cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// GetAllUncrawled MongoDoc(s)
func (ms *MongoStore) GetAllUncrawled() ([]*MongoDoc, error) {
	cur, err := ms.uncrawledCollection.Find(ms.mongoCtx, bson.D{})
	if err != nil {
		return nil, err
	}

	var vals []*MongoDoc

	for cur.Next(ms.mongoCtx) {
		var doc *MongoDoc

		err = cur.Decode(&doc)
		if err != nil {
			return nil, err
		}

		log.Printf("%v", doc)

		vals = append(vals, doc)
	}

	return vals, nil
}

// UpdateUncrawled k with v(s)
func (ms *MongoStore) UpdateUncrawled(k string, v []string) error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"docKey", k}}
	update := bson.D{{"$set", bson.D{{"docValue", v}}}}

	result := ms.uncrawledCollection.FindOneAndUpdate(ms.mongoCtx, filter, update, opts)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil
		}

		return err
	}

	return nil
}
