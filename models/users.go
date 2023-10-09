package models

import (
	"errors"
	"fmt"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

func SchemaPublic(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("public" + "." + tableName)
	}
}

type UserModels interface {
	InsertUser(user tables.UsersMA) (tables.UsersMA, error)
	CreateUser(user tables.Users) (tables.Users, error)
	InputFormUsers(formID int) ([]tables.InputFormUsers, error)
	InputFormUserLeftJoin(formID int, joinEnd string) ([]objects.ObjectInputFormUsers, error)
	InputFormUserPaging(formID int, whreStr string, paging objects.Paging) ([]tables.InputFormUsers, error)
	GetUserRow(fields tables.Users) (tables.UserData, error)
	GetUserRows(fields tables.Users) ([]tables.UserData, error)
	GetListAdminEks(companyID int, whreStr string, paging objects.Paging) ([]objects.AdminEks, error)
	GetTotalFormAdminEksternal(userID int, companyID int) ([]objects.AdminEks, error)
	GetUserWhereRow(user tables.Users, whreStr string) (tables.UserData, error)
}

type connection struct {
	db *gorm.DB
}

func NewUserModels(dbg *gorm.DB) UserModels {
	return &connection{
		db: dbg,
	}
}

func (con *connection) InsertUser(data tables.UsersMA) (tables.UsersMA, error) {
	err := con.db.Scopes(SchemaUsr("users")).Create(&data).Error
	if err != nil {
		fmt.Println(err)
		return tables.UsersMA{}, err
	}
	return data, nil
}

func (con *connection) ListUsers() ([]tables.Users, error) {
	// err := r.db.Table("users").Select("name", "email").Row()
	// err := r.db.Model(&data).Limit(10).Find(&APIUser{}).Error
	// return data, err

	var data []tables.Users
	err := con.db.Find(&data).Error
	return data, err
}

func (con *connection) CreateUser(data tables.Users) (tables.Users, error) {
	err := con.db.Create(&data).Error
	return data, err
}

func (con *connection) GetUserRow(fields tables.Users) (tables.UserData, error) {
	var data tables.UserData

	err := con.db.Table("usr.users").Select("users.id, users.name, users.phone, users.email, users.avatar, users.date_of_birth, users.gender_id, t.translation as gender_name, o.name as company_name").Joins("left join mstr.genders g on g.id = users.gender_id").Joins("left join mstr.translations t on t.textcontent_id = g.name_textcontent_id").Joins("left join mstr.organizations o on o.created_by = users.id and o.is_default = true").Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) GetUserRows(fields tables.Users) ([]tables.UserData, error) {
	var data []tables.UserData

	err := con.db.Table("usr.users").Select("users.id, users.name, users.phone, users.email, users.avatar, users.date_of_birth, users.gender_id, t.translation as gender_name, o.name as company_name").Joins("left join mstr.genders g on g.id = users.gender_id").Joins("left join mstr.translations t on t.textcontent_id = g.name_textcontent_id and t.language_id = 1").Joins("left join mstr.organizations o on o.created_by = users.id AND o.is_default is true").Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) FindByEmail(email string) tables.Users {
	var user tables.Users
	con.db.Where("email = ?", email).Take(&user)
	return user
}
func (con *connection) InputFormUserLeftJoin(formID int, joinEnd string) ([]objects.ObjectInputFormUsers, error) {

	var user []objects.ObjectInputFormUsers
	err := con.db.Raw(`select fu.user_id , ur.name as user_name, ur.phone as user_phone , 
						(select o."name" organizations  
							from frm.form_user_organizations fuo  
							join mstr.organizations o on o.id = fuo.organization_id 
							where fuo.form_user_id = fu.id ) as organizations,
						(select to_char(created_at::TIMESTAMP, 'DD MON YYYY; HH24:MI')  from frm.input_forms_` + strconv.Itoa(formID) + ` where user_id = fu.user_id order by created_at desc limit 1)submit_date
					from frm.form_users fu
					left join (SELECT f1.user_id, u.name as user_name, u.phone as user_phone
									from (select f.user_id from frm.input_forms_` + strconv.Itoa(formID) + ` f group by f.user_id) f1
									LEFT JOIN usr.users u on u.id= f1.user_id) as inf ON inf.user_id = fu.user_id
					left join usr.users ur on ur.id = fu.user_id
					left join frm.form_user_organizations fuo ON fuo.form_user_id=fu.id 
					
					where fu.form_id = ` + strconv.Itoa(formID) + ` ` + joinEnd).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (con *connection) InputFormUsers(formID int) ([]tables.InputFormUsers, error) {

	var user []tables.InputFormUsers
	err := con.db.Raw(`SELECT f1.user_id, u.name as user_name, u.phone as user_phone
				from (select f.user_id from frm.input_forms_` + strconv.Itoa(formID) + ` f group by f.user_id) f1
				LEFT JOIN usr.users u on u.id= f1.user_id
				`).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (con *connection) InputFormUserPaging(formID int, whreStr string, paging objects.Paging) ([]tables.InputFormUsers, error) {

	limitOffset := ``
	if paging.Limit > 0 && paging.Page > 0 {
		offset := (paging.Page - 1) * paging.Limit
		limitOffset = `LIMIT ` + strconv.Itoa(paging.Limit) + ` OFFSET ` + strconv.Itoa(offset)
	}

	sortData := ``
	if paging.Sort != "" && paging.SortBy != "" {
		sortData = `ORDER BY ` + paging.SortBy + ` ` + paging.Sort
	}

	var user []tables.InputFormUsers
	err := con.db.Raw(`SELECT fi.user_id, u.name as user_name, u.phone as user_phone, u.avatar
				from (select f.user_id from frm.input_forms_` + strconv.Itoa(formID) + ` f group by f.user_id) fi
				LEFT JOIN usr.users u on u.id= fi.user_id
				` + whreStr + `
				` + limitOffset + `
				` + sortData + `
				`).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (con *connection) GetListAdminEks(companyID int, whreStr string, paging objects.Paging) ([]objects.AdminEks, error) {
	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}
	var data []objects.AdminEks

	err := con.db.Raw(`
	SELECT uo.user_id, uo.organization_id, u.name, u.email, u.phone, o.name as organization_name
	FROM usr.user_organization_invites uoi 
	LEFT JOIN usr.user_organizations uo on uo.user_organization_invite_id = uoi.id 
	LEFT JOIN usr.users u on u.id = uo.user_id
	LEFT JOIN mstr.organizations o on o.id = uoi.organization_id 
	WHERE uo.organization_id = ?
	UNION
	(SELECT	fu.user_id , fo.organization_id, u2.name, u2.email, null as phone, null as organization_name  
	FROM frm.form_users fu
	LEFT JOIN usr.users u2 on u2.id = fu.user_id  
	LEFT JOIN frm.form_organizations fo on fo.form_id = fu.form_id 
	WHERE fu.type = 'guest' AND fo.organization_id = ?) 
	`, companyID, companyID).Where(whreStr).Order(orderBy).Limit(paging.Limit).Offset(offset).Group("uo.user_id, uo.organization_id, u.name, u.email, u.phone, o.name").Find(&data).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) GetTotalFormAdminEksternal(userID int, companyID int) ([]objects.AdminEks, error) {
	var data []objects.AdminEks

	err := con.db.Table("frm.form_users").
		Select("form_users.id").
		Joins("left join frm.form_organizations fo on fo.form_id = form_users.form_id").
		Where("form_users.user_id = ?", userID).
		Where("fo.organization_id = ?", companyID).
		Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *connection) GetUserWhereRow(fields tables.Users, whreStr string) (tables.UserData, error) {
	var data tables.UserData
	// err := con.db.Scopes(SchemaUsr("users")).Where(fields).First(&data).Error

	err := con.db.Table("usr.users").Select("users.id, users.name, users.phone, users.email, users.avatar, users.date_of_birth, users.gender_id, t.translation as gender_name, users.remember_token, users.created_at").Joins("left join mstr.genders g on g.id = users.gender_id").Joins("left join mstr.translations t on t.textcontent_id = g.name_textcontent_id").Where(fields).Where(whreStr).First(&data).Error
	if err != nil {
		return tables.UserData{}, err
	}
	return data, nil
}
