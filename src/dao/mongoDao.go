package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoDao struct {
	engine *mongo.Database
}

func NewMongoDao(engine *mongo.Database) *MongoDao {
	return &MongoDao{engine:engine}
}

func (d *MongoDao) InsertOne(collection string, val interface{}) error {
	_, err := d.engine.Collection(collection).InsertOne(context.TODO(), val)

	return err
}

func (d *MongoDao) FindOne(collection string, filter interface{},) (interface{}, error) {

	var val struct{
		Name string
		Time time.Time
	}

	err := d.engine.Collection(collection).FindOne(context.TODO(), filter).Decode(&val)
	if err != nil {
		return err, nil
	}

	return val, nil
}
