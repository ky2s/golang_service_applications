package models

import (
	"errors"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

type RuleModels interface {
	ConditionRuleList(langID int, fields tables.ConditionRules) ([]tables.ConditionRules, error)
	InsertFormFieldRule(fields tables.FormFieldConditionRules) (tables.FormFieldConditionRules, error)
	UpdateFormFieldRule(id int, fields tables.FormFieldConditionRules) (bool, error)
	GetFormFieldRuleRow(fields tables.FormFieldConditionRules, fieldString string) (tables.FormFieldConditionRules, error)
	GetFormFieldRuleRows(fields tables.FormFieldConditionRules, stringFields string) ([]tables.FormFieldConditionRules, error)
	DeleteFormFieldRule(formFieldID int) (bool, error)
	DeleteFormFieldRulePrimary(formFieldID int) (bool, error)
}

type ruleConnection struct {
	db *gorm.DB
}

func NewRuleModels(dbg *gorm.DB) RuleModels {
	return &ruleConnection{
		db: dbg,
	}
}

func (con *ruleConnection) ConditionRuleList(langID int, fields tables.ConditionRules) ([]tables.ConditionRules, error) {
	var data []tables.ConditionRules

	err := con.db.Table("mstr.condition_rules").Select("condition_rules.id, t.translation as name, condition_rules.code").Joins("join mstr.translations t on condition_rules.name_textcontent_id = t.textcontent_id ").Where("t.language_id", langID).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *ruleConnection) InsertFormFieldRule(data tables.FormFieldConditionRules) (tables.FormFieldConditionRules, error) {

	err := con.db.Scopes(SchemaFrm("form_field_condition_rules")).Create(&data).Error
	if err != nil {
		return tables.FormFieldConditionRules{}, err
	}
	return data, err
}

func (con *ruleConnection) UpdateFormFieldRule(id int, data tables.FormFieldConditionRules) (bool, error) {

	err := con.db.Scopes(SchemaFrm("form_field_condition_rules")).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *ruleConnection) GetFormFieldRuleRow(fields tables.FormFieldConditionRules, fieldString string) (tables.FormFieldConditionRules, error) {

	var data tables.FormFieldConditionRules
	err := con.db.Scopes(SchemaFrm("form_field_condition_rules")).Where(fields).Where(fieldString).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *ruleConnection) GetFormFieldRuleRows(fields tables.FormFieldConditionRules, stringFields string) ([]tables.FormFieldConditionRules, error) {

	var data []tables.FormFieldConditionRules
	err := con.db.Scopes(SchemaFrm("form_field_condition_rules")).Where(fields).Where(stringFields).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *ruleConnection) DeleteFormFieldRule(formFieldID int) (bool, error) {

	err := con.db.Exec("DELETE FROM frm.form_field_condition_rules where form_field_id = " + strconv.Itoa(formFieldID) + " and condition_parent_field_id is not null").Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *ruleConnection) DeleteFormFieldRulePrimary(formFieldID int) (bool, error) {

	err := con.db.Exec("DELETE FROM frm.form_field_condition_rules where form_field_id = " + strconv.Itoa(formFieldID) + " and condition_parent_field_id is null").Error
	if err != nil {
		return false, err
	}

	return true, err
}
