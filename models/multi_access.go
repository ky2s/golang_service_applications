package models

import (
	"errors"
	"snapin-form/objects"

	"gorm.io/gorm"
)

type MultiAccessModels interface {
	GetSenderRole(SenderID int) (objects.UserOrganizationRoles, error)
	InsertFormToUserInvites(dataInput objects.SelectOrganization) (objects.SelectOrganization, error)
	CheckFormAlreadyHave(FormID int, ReceiverID int, OrganizationID int) ([]objects.UserOrganizationRoles, error)
	CheckUserAlreadyConnect(dataInput objects.InputFormUserOrganizations) ([]objects.InputFormUserOrganizations, error)
}

type multiAccessConnection struct {
	db *gorm.DB
}

func NewMultiAccessModels(dbg *gorm.DB) MultiAccessModels {
	return &multiAccessConnection{
		db: dbg,
	}
}

func (con *multiAccessConnection) GetSenderRole(SenderID int) (objects.UserOrganizationRoles, error) {
	var data objects.UserOrganizationRoles
	err := con.db.Table("usr.user_organization_roles").
		Select("user_organization_roles.id, user_organization_roles.role_id").
		Joins("left join usr.user_organizations uo on uo.id = user_organization_roles.user_organization_id").
		Where("uo.user_id = ?", SenderID).
		Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return objects.UserOrganizationRoles{}, nil
	}

	if err != nil {
		return objects.UserOrganizationRoles{}, err
	}

	return data, nil
}

func (con *multiAccessConnection) InsertFormToUserInvites(dataInput objects.SelectOrganization) (objects.SelectOrganization, error) {

	err := con.db.Scopes(SchemaFrm("form_to_user_invites")).Create(&dataInput).Error
	if err != nil {
		return objects.SelectOrganization{}, err
	}
	return dataInput, err
}

func (con *multiAccessConnection) CheckFormAlreadyHave(FormID int, ReceiverID int, OrganizationID int) ([]objects.UserOrganizationRoles, error) {
	var data []objects.UserOrganizationRoles
	err := con.db.Table("frm.form_to_user_invites").
		Select("form_to_user_invites.id").
		Where("form_to_user_invites.form_id = ?", FormID).
		Where("form_to_user_invites.user_receiver_id = ?", ReceiverID).
		Where("form_to_user_invites.organization_receiver_id = ?", OrganizationID).
		Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return data, nil
	}

	if err != nil {
		return data, err
	}

	return data, nil
}

func (con *multiAccessConnection) CheckUserAlreadyConnect(dataInput objects.InputFormUserOrganizations) ([]objects.InputFormUserOrganizations, error) {
	var result []objects.InputFormUserOrganizations
	err := con.db.Table("frm.form_users").
		Select("form_to_user_invites.*").
		Where(dataInput).
		Find(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return result, nil
	}

	if err != nil {
		return result, err
	}

	return result, nil
}
