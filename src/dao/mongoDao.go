package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func (d *MongoDao) GetHotDoctor()([]int, error) {
	type DoctorNpiMongo struct {
		Npi int `bson:"npi"`
		Count int `bson:"count"`
	}

	groupStage := bson.D{{"$group",
		bson.D{{"_id", "$npi"}, {"count", bson.D{{"$sum", 1}}}}}}

	addFieldsStage := bson.D{{"$addFields", bson.D{{"npi", "$_id"}}}}

	sortStage := bson.D{{"$sort", bson.D{{"count", -1}}}}
	limitStage := bson.D{{"$limit", 20}}

	cur, err := d.engine.Collection(SearchDrResultRecord).Aggregate(context.TODO(), mongo.Pipeline{
		groupStage, addFieldsStage, sortStage, limitStage})

	if err != nil {
		return nil, err
	}

	var resultList []bson.M
	if err = cur.All(context.TODO(), &resultList); err != nil {
		return nil, err
	}

	npiList := make([]int, 0)
	for _, doctor := range resultList {
		var d DoctorNpiMongo

		b, _ := bson.Marshal(doctor)
		_ = bson.Unmarshal(b, &d)

		npiList = append(npiList, d.Npi)
	}

	return npiList, nil
}
