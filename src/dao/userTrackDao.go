package dao

import (
	"DrFinder/src/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserActionColName = "useraction"

type UserTrackDao struct {
	engine *mongo.Database
}

func NewUserTrackDao(engine *mongo.Database) *UserTrackDao {
	return &UserTrackDao{engine:engine}
}

func (d *UserTrackDao) AddActionEvent(actEvent *models.UserAction) error {
	_, err := d.engine.Collection(UserActionColName).InsertOne(context.TODO(), actEvent)

	return err
}

func (d *UserTrackDao) FindActionEvent(filter interface{}) []models.UserAction {
	cur, err := d.engine.Collection(UserActionColName).Find(context.TODO(), filter)
	if err != nil {
		return nil
	}

	var results []models.UserAction
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.UserAction
		err := cur.Decode(&elem)
		if err != nil {
			//
		}else {
			results = append(results, elem)
		}
	}

	return results
}