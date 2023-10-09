package models

import (
	"snapin-form/objects"
	"snapin-form/tables"

	"gorm.io/gorm"
)

type FormOtpModels interface {
	InsertFormOtp(data tables.FormOtps) (tables.FormOtps, error)
	GetFormOtp(fields tables.FormOtps) (tables.FormOtps, error)
	UpdateFormOtp(fielID int, data objects.FormOtp) (bool, error)
	DeleteWhrFormOtp(whre tables.FormOtps) (bool, error)
}

type formOtpConnection struct {
	db *gorm.DB
}

func NewFormOtpModels(dbg *gorm.DB) FormOtpModels {
	return &formOtpConnection{
		db: dbg,
	}
}

func (con *formOtpConnection) InsertFormOtp(data tables.FormOtps) (tables.FormOtps, error) {

	err := con.db.Scopes(SchemaFrm("form_otps")).Create(&data).Error
	if err != nil {
		return tables.FormOtps{}, err
	}
	return data, nil
}

func (con *formOtpConnection) GetFormOtp(fields tables.FormOtps) (tables.FormOtps, error) {

	var data tables.FormOtps
	err := con.db.Scopes(SchemaFrm("form_otps")).Where(fields).First(&data).Error
	if err != nil {
		return tables.FormOtps{}, err
	}
	return data, nil
}

func (con *formOtpConnection) UpdateFormOtp(id int, data objects.FormOtp) (bool, error) {

	err := con.db.Scopes(SchemaFrm("form_otps")).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return false, err
	}

	err = con.db.Scopes(SchemaFrm("form_otps")).Where("id = ?", id).Update("status", data.Status).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *formOtpConnection) DeleteWhrFormOtp(whre tables.FormOtps) (bool, error) {

	// err := con.db.Scopes(SchemaFrm("form_otps")).Where("id = ?", id).Updates(data).Error
	var tbl tables.FormOtps
	err := con.db.Scopes(SchemaFrm("form_otps")).Where(whre).Delete(&tbl).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
