package postgre

import (
	"fmt"
	"nostr-relay/pkg/db"
	"nostr-relay/pkg/models"

	"gorm.io/gorm"
)

func GetMainDB() (*gorm.DB, error) {
	db := db.GetMainDB()
	if db == nil {
		return nil, fmt.Errorf("main db is nil")
	}

	return db, nil
}

func SaveEvent(data *models.RelayEvent) error {
	db, err := GetMainDB()
	if err != nil {
		return err
	}

	dbc := db.Create(data)
	return dbc.Error
}

func GetEvent(limit int) []models.RelayEvent {
	var ret []models.RelayEvent
	db, err := GetMainDB()
	if err != nil {
		return ret
	}

	db.Find(ret).Limit(limit)
	return ret
}

func GetLastEvent() *models.RelayEvent {
	var ret models.RelayEvent
	db, err := GetMainDB()
	if err != nil {
		return nil
	}

	db.Last(&ret)
	return &ret
}

func GetEventFrom(id int) []models.RelayEvent {
	var ret []models.RelayEvent
	db, err := GetMainDB()
	if err != nil {
		return ret
	}

	db.Where("ID >?", id).Find(&ret)
	return ret
}
