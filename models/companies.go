package models

import (
	"errors"
	"fmt"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

type CompaniesModels interface {
	GetCompaniesRow(data tables.Organizations) (tables.Organizations, error)
	GetCompaniesRows(fields tables.Organizations) ([]tables.Organizations, error)
	ConnectedUserCompanies(fields objects.UserOrganizations) (tables.UserOrganizations, error)
	InsertUserCompaniyToRole(data tables.UserOrganizationRoles) (tables.UserOrganizationRoles, error)
	UpdateUserCompaniyToRole(userID int, userOrgID int, roleID int) (bool, error)
	GetUserCompaniyToRole(data tables.UserOrganizationRoles, whreStr string) (objects.UserOrganizationRoles, error)
	GetUserCompaniesRow(fields objects.UserOrganizations, whereString string) (tables.SelectUserOrganizations, error)
	GetFormOrganizationRow(fields tables.FormOrganizations) (objects.FormOrganizations, error)
}

func NewCompaniesModels(dbg *gorm.DB) CompaniesModels {
	return &connection{
		db: dbg,
	}
}

func (con *connection) GetFormOrganizationRow(fields tables.FormOrganizations) (objects.FormOrganizations, error) {
	var data objects.FormOrganizations
	err := con.db.Table("frm.form_organizations").Select("form_organizations.*, o.name as organization_name, o.profile_pic as organization_profile_pic").Joins("join mstr.organizations o on o.id=form_organizations.organization_id").Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) GetCompaniesRow(fields tables.Organizations) (tables.Organizations, error) {
	var data tables.Organizations
	err := con.db.Scopes(SchemaMstr("organizations")).Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) GetCompaniesRows(fields tables.Organizations) ([]tables.Organizations, error) {
	var data []tables.Organizations
	err := con.db.Scopes(SchemaMstr("organizations")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) ConnectedUserCompanies(data objects.UserOrganizations) (tables.UserOrganizations, error) {

	var dataUserComp tables.UserOrganizations
	dataUserComp.UserID = data.UserID
	dataUserComp.OrganizationID = data.OrganizationID

	err := con.db.Scopes(SchemaUsr("user_organizations")).Create(&dataUserComp).Error
	if err != nil {
		fmt.Println("error user_organizations --- ", err)
		return tables.UserOrganizations{}, err
	}

	return dataUserComp, nil
}

func (con *connection) InsertUserCompaniyToRole(data tables.UserOrganizationRoles) (tables.UserOrganizationRoles, error) {

	err := con.db.Scopes(SchemaUsr("user_organization_roles")).Create(&data).Error
	if err != nil {
		fmt.Println("error user_organization_roles --- ", err)
		return tables.UserOrganizationRoles{}, err
	}

	return data, nil
}

func (con *connection) UpdateUserCompaniyToRole(userID int, userOrgID int, roleID int) (bool, error) {

	var data tables.UserOrganizationRoles
	data.RoleID = roleID
	err := con.db.Scopes(SchemaUsr("user_organization_roles")).Where("user_organization_roles.user_organization_id = (select uo.id from usr.user_organizations uo where uo.user_id = " + strconv.Itoa(userID) + " AND uo.organization_id = " + strconv.Itoa(userOrgID) + ")").Updates(data).Error
	if err != nil {
		fmt.Println("error user_organization_roles --- ", err)
		return false, err
	}

	return true, nil
}

func (con *connection) GetUserCompaniyToRole(fields tables.UserOrganizationRoles, whreStr string) (objects.UserOrganizationRoles, error) {

	var result objects.UserOrganizationRoles
	err := con.db.Table("usr.user_organization_roles").Select("user_organization_roles.id, user_organization_roles.user_organization_id, uo.organization_id, o.name as organization_name, user_organization_roles.role_id, r.name as role_name, uo.user_id, u.name as user_name").Joins("left join usr.user_organizations uo on uo.id = user_organization_roles.user_organization_id").Joins("left join mstr.organizations o on o.id = uo.organization_id").Joins("left join usr.roles r on r.id = user_organization_roles.role_id").Joins("left join usr.users u on u.id = uo.user_id").Where(fields).Where(whreStr).Find(&result).Error
	if err != nil {
		fmt.Println("error user_organization_roles --- ", err)
		return objects.UserOrganizationRoles{}, err
	}

	return result, nil
}

func (con *connection) GetUserCompaniesRow(fields objects.UserOrganizations, whereString string) (tables.SelectUserOrganizations, error) {
	var data tables.SelectUserOrganizations

	err := con.db.Table("usr.user_organizations").Select("user_organizations.id, user_organizations.user_id, u.phone as user_phone, u.email as user_email, ogz.id as organization_id, ogz.code, ogz.name, ogz.description, ogz.address1, ogz.phone1, ogz.phone2, ogz.email, ogz.contact_name, ogz.contact_phone, ogz.contact_position").Joins("left join mstr.organizations ogz on ogz.id = user_organizations.organization_id").Joins("left join usr.users u on u.id = user_organizations.user_id").Joins("join usr.organization_subscription_plans osp ON osp.organization_id = ogz.id AND osp.is_blocked = false").Where(fields).Where(whereString).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}

	return data, nil
}
