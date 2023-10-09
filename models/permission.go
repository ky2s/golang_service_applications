package models

import (
	"errors"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

type PermissionModels interface {
	PermissionList(fields tables.Permissions) ([]tables.Permissions, error)
	GetPermissionRow(fields tables.Permissions, fieldString string) (tables.Permissions, error)
	GetPermissionRows(fields tables.Permissions, stringFields string) ([]tables.Permissions, error)
	InsertFormUserPermission(fields tables.FormUserPermission) (tables.FormUserPermission, error)
	UpdateFormUserPermission(id int, fields tables.FormUserPermission) (tables.FormUserPermission, error)
	DeleteFormUserPermission(id int, whreString string) (bool, error)
	GetFormUserPermissionRows(fields tables.FormUserPermissionJoin, stringFields string) ([]tables.FormUserPermissionJoin, error)
	GetFormUserPermissionRow(fields tables.FormUserPermissionJoin, stringFields string) (tables.FormUserPermissionJoin, error)
	// GetUserPermissionToFormRow(fields tables.FormUserPermissionJoin, stringFields string) ([]tables.FormUserPermissionJoin, error)
	UpdateUserOrganizationPermission(id int, fields tables.UserOrganizationPermission) (tables.UserOrganizationPermission, error)
	GetUserOrganizationPermissionRows(companyID int) ([]tables.UserOrganizationPermissionJoin, error)
	GetUserOrganizationPermissionRow(userID int) ([]tables.UserOrganizationPermissionJoin, error)
	InsertAttendanceOrganization(fields tables.UserOrganizationPermission) (tables.UserOrganizationPermission, error)
	InsertUserOrganizationPermission(fields tables.UserOrganizationPermission) (tables.UserOrganizationPermission, error)

	GetMissingID(fields objects.IDAdminEks) ([]objects.IDAdminEks, error)
}

type permissionConnection struct {
	db *gorm.DB
}

func NewPermissionModels(dbg *gorm.DB) PermissionModels {
	return &permissionConnection{
		db: dbg,
	}
}

func (con *permissionConnection) PermissionList(fields tables.Permissions) ([]tables.Permissions, error) {
	var data []tables.Permissions

	err := con.db.Table("usr.permissions").Select("permissions.id, t.translation as name, permissions.slug").Joins("join mstr.translations t on permissions.name_textcontent_id = t.textcontent_id ").Where("t.language_id= ?", 1).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *permissionConnection) GetPermissionRow(fields tables.Permissions, fieldString string) (tables.Permissions, error) {

	var data tables.Permissions
	err := con.db.Scopes(SchemaUsr("permissions")).Where(fields).Where(fieldString).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *permissionConnection) GetPermissionRows(fields tables.Permissions, stringFields string) ([]tables.Permissions, error) {

	var data []tables.Permissions
	err := con.db.Scopes(SchemaUsr("permissions")).Where(fields).Where(stringFields).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *permissionConnection) InsertFormUserPermission(fields tables.FormUserPermission) (tables.FormUserPermission, error) {

	err := con.db.Scopes(SchemaUsr("form_user_permissions")).Create(&fields).Error
	if err != nil {
		return tables.FormUserPermission{}, err
	}

	return fields, nil
}

func (con *permissionConnection) UpdateFormUserPermission(id int, fields tables.FormUserPermission) (tables.FormUserPermission, error) {

	err := con.db.Scopes(SchemaUsr("form_user_permissions")).Where("id = ?", id).Update("status", fields.Status).Error
	if err != nil {
		return tables.FormUserPermission{}, err
	}

	return fields, nil
}

func (con *permissionConnection) GetFormUserPermissionRows(fields tables.FormUserPermissionJoin, stringFields string) ([]tables.FormUserPermissionJoin, error) {

	var data []tables.FormUserPermissionJoin
	err := con.db.Table("usr.form_user_permissions").Select("form_user_permissions.id, form_user_permissions.form_user_id, form_user_permissions.permission_id, form_user_permissions.status, p.name as permission_name").Joins("left join  usr.permissions p on p.id = form_user_permissions.permission_id").Joins("left join frm.form_users fu ON fu.id = form_user_permissions.form_user_id").Where(fields).Where(stringFields).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *permissionConnection) GetFormUserPermissionRow(fields tables.FormUserPermissionJoin, stringFields string) (tables.FormUserPermissionJoin, error) {

	var data tables.FormUserPermissionJoin
	err := con.db.Table("usr.form_user_permissions").Select("form_user_permissions.id, form_user_permissions.form_user_id, form_user_permissions.permission_id, form_user_permissions.status, p.name as permission_name").Joins("left join  usr.permissions p on p.id = form_user_permissions.permission_id").Joins("left join frm.form_users fu ON fu.id = form_user_permissions.form_user_id").Where(fields).Where(stringFields).Find(&data).Error
	if err != nil {
		return tables.FormUserPermissionJoin{}, err
	}
	return data, nil
}

func (con *permissionConnection) DeleteFormUserPermission(id int, whreString string) (bool, error) {

	err := con.db.Exec("DELETE FROM usr.form_user_permissions where form_user_id = " + strconv.Itoa(id)).Error
	if err != nil {
		return false, err
	}

	return true, err
}

// func (con *permissionConnection) GetUserPermissionToFormRow(fields tables.FormUserPermissionJoin, stringFields string) ([]tables.FormUserPermissionJoin, error) {

//		var data []tables.FormUserPermissionJoin
//		err := con.db.Table("usr.form_user_permissions").Select("form_user_permissions.id, form_user_permissions.form_user_id, form_user_permissions.permission_id, form_user_permissions.status, p.name as permission_name").Joins("left join  usr.permissions p on p.id = form_user_permissions.permission_id").Joins("left join frm.form_users fu ON fu.id = form_user_permissions.form_user_id").Where(fields).Where(stringFields).Find(&data).Error
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, err
//		}
//		return data, nil
//	}
func (con *permissionConnection) UpdateUserOrganizationPermission(id int, fields tables.UserOrganizationPermission) (tables.UserOrganizationPermission, error) {

	err := con.db.Scopes(SchemaUsr("user_organization_permissions")).Where("id = ?", id).Update("is_checked", fields.IsChecked).Error
	if err != nil {
		return tables.UserOrganizationPermission{}, err
	}

	return fields, nil
}

func (con *permissionConnection) GetUserOrganizationPermissionRows(companyID int) ([]tables.UserOrganizationPermissionJoin, error) {

	var data []tables.UserOrganizationPermissionJoin
	err := con.db.Table("usr.user_organization_permissions").Select("user_organization_permissions.id, user_organization_permissions.user_organization_id, user_organization_permissions.permission_id, user_organization_permissions.is_checked, p.name as permission_name").Joins("left join  usr.permissions p on p.id = user_organization_permissions.permission_id").Joins("left join usr.user_organizations uo on uo.id = user_organization_permissions.user_organization_id").Where("uo.organization_id = ?", companyID).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *permissionConnection) GetUserOrganizationPermissionRow(userID int) ([]tables.UserOrganizationPermissionJoin, error) {

	var data []tables.UserOrganizationPermissionJoin
	err := con.db.Table("usr.user_organization_permissions").
		Select("user_organization_permissions.id, user_organization_permissions.user_organization_id, user_organization_permissions.permission_id, user_organization_permissions.is_checked, p.name as permission_name, uo.user_id").
		Joins("left join  usr.permissions p on p.id = user_organization_permissions.permission_id").
		Joins("left join usr.user_organizations uo on uo.id = user_organization_permissions.user_organization_id").
		Where("uo.user_id = ?", userID).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *permissionConnection) InsertAttendanceOrganization(fields tables.UserOrganizationPermission) (tables.UserOrganizationPermission, error) {

	err := con.db.Scopes(SchemaUsr("user_organization_permissions")).Create(&fields).Error
	if err != nil {
		return tables.UserOrganizationPermission{}, err
	}

	return fields, nil
}

func (con *permissionConnection) GetMissingID(fields objects.IDAdminEks) ([]objects.IDAdminEks, error) {

	var data []objects.IDAdminEks
	err := con.db.Table("usr.user_organizations_").
		Select("user_organizations.id").
		Joins("left join usr.user_organization_permissions uop ON user_organizations.id = uop.user_organization_id").
		Joins("left join usr.user_organization_roles uor ON user_organizations.id = uor.user_organization_id").Where("uop.user_organization_id IS null").Where("uor.role_id = ?", 2).Where("user_organization_invite_id is null").
		Where(fields).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *permissionConnection) InsertUserOrganizationPermission(fields tables.UserOrganizationPermission) (tables.UserOrganizationPermission, error) {

	err := con.db.Scopes(SchemaUsr("user_organization_permissions")).Create(&fields).Error
	if err != nil {
		return tables.UserOrganizationPermission{}, err
	}

	return fields, nil
}
