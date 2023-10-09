package models

import (
	"errors"
	"snapin-form/tables"

	"gorm.io/gorm"
)

type FieldTypeModels interface {
	InsertFieldType(data tables.FieldTypes) (tables.FieldTypes, error)
	GetFieldTypeRows(data tables.FieldTypes) ([]tables.FieldTypeTrans, error)
	GetFieldTypeRow(data tables.FieldTypes) (tables.FieldTypeTrans, error)
	UpdateFieldType(id int, data tables.FieldTypes) (bool, error)
	DeleteFieldType(id int) (bool, error)
}

type fieldTypeConnection struct {
	db *gorm.DB
}

func NewFieldTypeModels(dbg *gorm.DB) FieldTypeModels {
	return &fieldTypeConnection{
		db: dbg,
	}
}

func (con *fieldTypeConnection) InsertFieldType(data tables.FieldTypes) (tables.FieldTypes, error) {
	err := con.db.Scopes(SchemaFrm("field_types")).Create(&data).Error

	return data, err
}

func (con *fieldTypeConnection) GetFieldTypeRows(fields tables.FieldTypes) ([]tables.FieldTypeTrans, error) {
	var data []tables.FieldTypeTrans
	err := con.db.Table("mstr.field_types").Select("field_types.id", "field_types.code", "t.translation as name", "field_types.real_var_type", "field_types.is_group", "field_types.is_media", "field_types.is_child_field", "field_types.name_textcontent_id", "'Pengguna bisa mengisi dengan teks apapaun' as info").Joins("join mstr.translations t on t.textcontent_id = field_types.name_textcontent_id").Where(fields).Where("t.language_id", 1).Where("field_types.deleted_at is null").Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *fieldTypeConnection) GetFieldTypeRow(fields tables.FieldTypes) (tables.FieldTypeTrans, error) {
	var data tables.FieldTypeTrans
	err := con.db.Scopes(SchemaFrm("field_types")).Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *fieldTypeConnection) UpdateFieldType(id int, fields tables.FieldTypes) (bool, error) {
	err := con.db.Scopes(SchemaFrm("field_types")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *fieldTypeConnection) DeleteFieldType(FieldTypeID int) (bool, error) {
	var FieldTypes tables.FieldTypes
	err := con.db.Scopes(SchemaFrm("field_types")).Delete(&FieldTypes, FieldTypeID).Error
	if err != nil {
		return false, err
	}

	return true, err
}
