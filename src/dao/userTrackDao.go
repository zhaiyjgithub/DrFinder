package dao

import (
	"DrFinder/src/models"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

const UserActionColName = "user_action"
const UserViewTimeColName = "user_view_time"
const UserSearchDrsRecord = "user_search_drs_record"
const SearchDrResultRecord = "search_dr_result_record"

type UserTrackDao struct {
	engine *mongo.Database
}

func NewUserTrackDao(engine *mongo.Database) *UserTrackDao {
	return &UserTrackDao{engine:engine}
}

func (d *UserTrackDao) AddActionEvent(event *models.UserAction) error {
	_, err := d.engine.Collection(UserActionColName).InsertOne(context.TODO(), event)

	return err
}

func (d *UserTrackDao) AddManyActionEvent(events []models.UserAction) error {
	var tmps []interface{}
	for i := 0; i< len(events); i ++ {
		tmps = append(tmps, events[i])
	}
	_, err := d.engine.Collection(UserActionColName).InsertMany(context.TODO(), tmps)

	return err
}

func (d *UserTrackDao) AddViewEvent(event *models.UserView) error {
	_, err := d.engine.Collection(UserViewTimeColName).InsertOne(context.TODO(), event)

	return err
}

func (d *UserTrackDao) AddManyViewTimeEvent(events []models.UserView) error {
	var docs []interface{}
	for i := 0; i < len(events); i ++ {
		docs = append(docs, events[i])
	}
	_, err := d.engine.Collection(UserViewTimeColName).InsertMany(context.TODO(), docs)

	return err
}

func (d *UserTrackDao) AddSearchDrsRecord(record *models.UserSearchDrRecord) error {
	_, err := d.engine.Collection(UserSearchDrsRecord).InsertOne(context.TODO(), record)

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

func (d *UserTrackDao) AddSearchDrResultRecords(records []models.DrSearchResultRecord) error {
	var docs []interface{}
	for i := 0; i < len(records); i ++ {
		docs = append(docs, records[i])
	}

	_, err := d.engine.Collection(SearchDrResultRecord).InsertMany(context.TODO(), docs)
	return err
}