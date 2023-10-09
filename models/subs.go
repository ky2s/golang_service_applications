package models

import (
	"errors"
	"snapin-form/objects"
	"snapin-form/tables"

	"gorm.io/gorm"
)

type SubsModels interface {
	GetPlanRow(fields objects.SubsPlan) (tables.SubsPlan, error)
	GetPlanPeriodRow(fields objects.SubsPlan) (objects.SubsPlanPeriod, error)
	InsertPlan(data tables.SubsPlan) (bool, tables.SubsPlan, error)
	UpdatePlan(orgID int, data tables.SubsPlan) (bool, tables.SubsPlan, error)
	UpdatePlanCurrent(orgID int, data tables.SubsPlan) (bool, tables.SubsPlan, error)
	InsertInjuryPlan(data tables.InjuryPlan) (bool, tables.InjuryPlan, error)
	GetInjuryPlanRows(fields objects.InjuryPlan) ([]tables.InjuryPlan, error)
}

type subsConnection struct {
	db *gorm.DB
}

func NewSubsModels(dbg *gorm.DB) SubsModels {
	return &subsConnection{
		db: dbg,
	}
}

func (con *subsConnection) GetPlanRow(fields objects.SubsPlan) (tables.SubsPlan, error) {
	var data tables.SubsPlan

	err := con.db.Scopes(SchemaUsr("organization_subscription_plans")).Select("organization_subscription_plans.* , sp.respondent_quota ").Joins("left join mstr.subscription_plans sp on sp.id = organization_subscription_plans.subscription_plan_id").Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *subsConnection) GetPlanPeriodRow(fields objects.SubsPlan) (objects.SubsPlanPeriod, error) {
	var data objects.SubsPlanPeriod

	err := con.db.Scopes(SchemaUsr("organization_subscription_plans as osp")).Select("osp.* , sp.respondent_quota, to_char((osp.created_at + (interval '1 day' * osp.total_period_days)),'yyyy-mm-dd') as period_end_date, DATE ((osp.created_at + (interval '1 day' * osp.total_period_days))) -  DATE (current_date) as period_remain ").Joins("left join mstr.subscription_plans sp on sp.id = osp.subscription_plan_id").Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *subsConnection) InsertPlan(data tables.SubsPlan) (bool, tables.SubsPlan, error) {

	err := con.db.Scopes(SchemaUsr("organization_subscription_plans")).Create(&data).Error
	if err != nil {
		return false, tables.SubsPlan{}, err
	}

	return true, data, nil
}

func (con *subsConnection) UpdatePlan(orgID int, data tables.SubsPlan) (bool, tables.SubsPlan, error) {

	err := con.db.Scopes(SchemaUsr("organization_subscription_plans")).Where("organization_id = ?", orgID).Updates(&data).Error
	if err != nil {
		return false, tables.SubsPlan{}, err
	}

	var dataSelect tables.SubsPlan
	var subPlan tables.SubsPlan
	subPlan.OrganizationID = orgID
	err2 := con.db.Scopes(SchemaUsr("organization_subscription_plans")).Where(subPlan).First(&dataSelect).Error
	if err2 != nil {
		return false, tables.SubsPlan{}, err2
	}

	return true, dataSelect, nil
}

func (con *subsConnection) UpdatePlanCurrent(orgID int, data tables.SubsPlan) (bool, tables.SubsPlan, error) {

	err := con.db.Exec("UPDATE usr.organization_subscription_plans SET respondent_current = ? , respon_updated_at=now() WHERE organization_id = ?", gorm.Expr("respondent_current + ?", data.RespondentCurrent), orgID).Error
	if err != nil {
		return false, tables.SubsPlan{}, err
	}

	err2 := con.db.Exec("UPDATE usr.organization_subscription_plans SET quota_current = ? WHERE organization_id = ?", gorm.Expr("quota_current - ?", data.RespondentCurrent), orgID).Error
	if err != nil {
		return false, tables.SubsPlan{}, err2
	}

	return true, data, nil
}

func (con *subsConnection) InsertInjuryPlan(data tables.InjuryPlan) (bool, tables.InjuryPlan, error) {

	err := con.db.Scopes(SchemaUsr("injury_plans")).Create(&data).Error
	if err != nil {
		return false, tables.InjuryPlan{}, err
	}

	return true, data, nil
}

func (con *subsConnection) GetInjuryPlanRows(fields objects.InjuryPlan) ([]tables.InjuryPlan, error) {
	var data []tables.InjuryPlan

	err := con.db.Scopes(SchemaUsr("injury_plans")).Where(fields).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}
