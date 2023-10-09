package models

import (
	"errors"
	"snapin-form/tables"

	"gorm.io/gorm"
)

type SettingModels interface {
	GetSettingRow(fields tables.Settings) (tables.Settings, error)
}

type settingConnection struct {
	db *gorm.DB
}

func NewSettingModels(dbg *gorm.DB) SettingModels {
	return &settingConnection{
		db: dbg,
	}
}

func (con *settingConnection) GetSettingRow(fields tables.Settings) (tables.Settings, error) {

	var data tables.Settings
	err := con.db.Scopes(SchemaMstr("app_settings")).Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}
