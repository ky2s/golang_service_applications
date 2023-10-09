package models

import (
	"errors"
	"fmt"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type FormModels interface {
	InsertForm(data tables.Forms) (tables.Forms, error)
	InsertFormOrganization(data tables.FormOrganizations) (tables.FormOrganizations, error)
	InsertFormAttendanceLocation(data objects.ObjectFormAttendanceLocations) (objects.ObjectFormAttendanceLocations, error)
	GetFormOrganization(fields tables.FormOrganizations) (objects.FormOrganizations, error)
	GetFormsOrganization(fields objects.HistoryBalanceSaldo) ([]objects.HistoryBalanceSaldo, error)
	GetFormsOrganizationWithFilter(fields objects.HistoryBalanceSaldo, whereString string) ([]objects.HistoryBalanceSaldo, error)
	GetFormRows(data tables.Forms) ([]tables.FormAll, error)
	GetFormWhreRows(data tables.Forms, whreStr string) ([]tables.FormAll, error)
	GetFormRow(fields tables.Forms) (tables.FormOut, error)
	GetFormPeriodeRangeRow(fields tables.Forms) (tables.FormPeriodRange, error)
	GetFormNotInProjectRows(fields tables.Forms, whereString string, paging objects.Paging) ([]tables.FormAll, error)
	GetFormOwnerRows(fields tables.FormOrganizationsJoin, whereString string, paging objects.Paging) ([]tables.FormAll, error)
	GetFormUnionProjectRows(fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) // form & project UNION
	GetFormMergeSuperAdmin(fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeSuperAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeSuperAdminNew1(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeSuperAdminApps(fields tables.FormOrganizationsJoin, whereName string, whereString string, paging objects.Paging) ([]tables.FormAll, error)
	GetProjectSuperAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeAdmin(fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeAdminNew1(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error)
	GetFormMergeAdminApps(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error)
	GetProjectAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error)
	GetFormUnionProjectAndExternalRows(userID int, fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error)
	UpdateForm(id int, data tables.Forms) (bool, error)
	UpdateFormAttendanceLocation(id int, data objects.ObjectFormAttendanceLocations) (bool, error)
	UpdateFormStatus(id int, fields tables.Forms) (bool, error)
	DeleteForm(id int, authorID int) (bool, error)
	DeleteFormLocation(id int) (bool, error)
	InsertFormUser(data objects.Forms) (tables.Forms, error)
	ConnectFormUser(data tables.FormUsers) (tables.FormUsers, error)
	ConnectFormUserOrg(dataInput objects.InputFormUserOrganizations) (tables.FormUsers, error)
	DeleteFormUser(userID int, formID int, userType string) (bool, error)
	DeleteFormUserOrg(userID int, formID int) (bool, error)
	UpdateFormUser(userID int, formID int, data tables.FormUsers) (bool, error)
	GetFormUserRow(data tables.FormUsers) (tables.JoinFormUsers, error)
	GetFormUserRows(data tables.JoinFormUsers, whreStr string) ([]tables.JoinFormUsers, error)
	GetFormUserToOrganizationRows(fields tables.JoinFormUsers, whereStr string) ([]tables.JoinFormUsers, error)
	GetFormUserToOrganizationRow(fields tables.JoinFormUsers, whereStr string) (tables.JoinFormUsers, error)
	GetFormUserRespondenRows(respondeID int) ([]tables.JoinFormUsers, error)
	GetFormUserUnionTeamRows(data tables.FormUsers) ([]tables.JoinFormUsers, error)
	GetFormUserAdminRows(data tables.FormUsers) ([]tables.JoinFormUsers, error)
	InsertFieldFile(data tables.FormFieldTempAssets) (bool, tables.FormFieldTempAssets, error)
	GetFormUserUniqRows(authorID int, projectID int) ([]tables.FormUsers, error)
	GetDetailFormUserRow(fields tables.FormUsers, whrString string) (tables.FormOut, error)
	GetUserFormOrganization(fields tables.UserFormOrganizations, whrStr string) (tables.FormOrganizations, error)

	// form company invite
	GetUserCompaniesListInvitedRows(fields tables.UserOrganizationInvites, whrStr string) ([]tables.UserOrganizationInviteDetail, error)
	InsertFormCompanyInvites(data tables.FormOrganizationInvites) (bool, error)
	DeleteFormCompanyInvites(formID int, orgID int) (bool, error)
	GetFormCompanyInviteRows(fields tables.JoinFormCompanies, whereStr string) ([]tables.JoinFormCompanies, error)
	GetFormCompanyInviteRow(fields tables.JoinFormCompanies, whereStr string) (tables.JoinFormCompanies, error)
	UpdateCompanyInviteForm(id int, fields tables.FormOrganizationInvites) (bool, error)
	GetAllIDForDelete(ID int) (objects.DeleteAdminEksObj, error)
	DeleteAdminEks(data objects.DeleteAdminEksObj) (objects.DeleteAdminEksObj, error)
	GetFormOrganizationInvite(formID int) (tables.JoinFormCompanies, error)
	GetFormCompanyInviteNew(fields tables.JoinFormCompanies, whereStr string) ([]tables.JoinFormCompanies, error)

	// tab list form other company
	GetFormOtherCompanyRows(fields tables.FormOrganizationsJoin, whereString string, paging objects.Paging) ([]tables.FormAll, error)

	GenerateFormUserOrg(tables.FormUserOrganizations) (tables.FormUserOrganizations, error)

	GetFormUserToFormOrgRows(fields tables.JoinFormUsers, whereStr string) ([]tables.JoinFormUsers, error)
	GetListFormEksternal(fields tables.Forms, whereString string, paging objects.Paging, user_id int, organization_id int) ([]tables.FormAll, error)
	GetListFormEksternalSuperAdmin(fields tables.Forms, whereString string, paging objects.Paging, organization_id int) ([]tables.FormAll, error)
	GetFormEksternalOwnerRows(fields tables.FormOrganizationsJoin, whereString string, paging objects.Paging) ([]tables.FormAll, error)
	GetBlastInfoData(formID int, whereString string) ([]objects.BlastInfoData, error)
	GetBlastInfoDataUsers(formID int, userID int, whereString string) ([]objects.BlastInfoData, error)
	GetFillingType(fields objects.FillingType) ([]objects.FillingType, error)

	GetDate() ([]objects.Date, error)
	GetAllForm() ([]objects.FormData, error)
	GetFormByOrganization(fields objects.HistoryBalanceSaldo) ([]objects.HistoryBalanceSaldo, error)

	// GetFormExport() ([]objects.Date, error)
	GetHistoryTopupByDate(organizationID int, whreDate string) ([]objects.TopupHistory, error)

	GetAllFormRows(fields tables.Forms) ([]tables.FormAll, error)

	// Form Template
	GetFormTemplate(UserID int) ([]objects.FormTemplateNew, error)
	GetFormTemplateByProjectID(UserID int, ProjectID int) ([]objects.FormTemplateNew, error)
	GetProject(UserID int) ([]objects.Projects, error)
}

type formConnection struct {
	db  *gorm.DB
	ctx *gin.Context
}

func NewFormModels(dbg *gorm.DB) FormModels {

	// fmt.Println("---ooo---", *gin.Context)
	return &formConnection{
		db: dbg,
	}
}

func SchemaUsr(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("usr" + "." + tableName)
	}
}

func SchemaFrm(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("frm" + "." + tableName)
	}
}

func SchemaMstr(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("mstr" + "." + tableName)
	}
}

func (con *formConnection) InsertForm(data tables.Forms) (tables.Forms, error) {
	err := con.db.Scopes(SchemaFrm("forms")).Create(&data).Error

	con.db.Exec(`CREATE TABLE "frm"."input_forms_` + strconv.Itoa(data.ID) + `" ("id" serial,"user_id" integer NOT NULL,"address" character varying(300), "geometry" geography DEFAULT null,"created_by" integer DEFAULT null,"updated_by" integer DEFAULT null,"deleted_by" integer DEFAULT null,"created_at" timestamptz,"updated_at" timestamptz,"updated_count" integer DEFAULT 0, "deleted_at" timestamptz,"is_pending" bool DEFAULT null, PRIMARY KEY ("id"))`)

	con.db.Exec(`ALTER TABLE "frm"."input_forms_` + strconv.Itoa(data.ID) + `" ADD CONSTRAINT input_forms_` + strconv.Itoa(data.ID) + `_user_id_foreign FOREIGN KEY (user_id) REFERENCES usr.users(id) MATCH FULL`)

	return data, err
}

func (con *formConnection) InsertFormOrganization(data tables.FormOrganizations) (tables.FormOrganizations, error) {
	err := con.db.Scopes(SchemaFrm("form_organizations")).Create(&data).Error
	if err != nil {
		return tables.FormOrganizations{}, err
	}

	return data, err
}

func (con *formConnection) InsertFormAttendanceLocation(data objects.ObjectFormAttendanceLocations) (objects.ObjectFormAttendanceLocations, error) {
	err := con.db.Raw(`insert into frm.form_attendance_locations (
		form_id, 
		"name", 
		"location", 
		geometry, 
		is_check_in, 
		is_check_out, 
		radius,
		created_at, 
		updated_at
		) 
		values (?, ?, ?, ST_SetSRID(ST_MakePoint(?, ?), 4326), ?, ?, ?, now(), now()) RETURNING id`,
		data.FormID,
		data.Name,
		data.Location,
		data.Longitude,
		data.Latitude,
		data.IsCheckIn,
		data.IsCheckOut,
		data.Radius).Scan(&data).Error
	return data, err
}

func (con *formConnection) GetFormOrganization(fields tables.FormOrganizations) (objects.FormOrganizations, error) {

	var data objects.FormOrganizations
	err := con.db.Scopes(SchemaFrm("form_organizations")).Select("form_organizations.*, o.name as organization_name, f.name as form_name").Joins("join mstr.organizations o on o.id=form_organizations.organization_id").Joins("left join frm.forms f on f.id = form_organizations.form_id").Where(fields).First(&data).Error
	if err != nil {
		return objects.FormOrganizations{}, err
	}

	return data, err
}

func (con *formConnection) GetFormsOrganization(fields objects.HistoryBalanceSaldo) ([]objects.HistoryBalanceSaldo, error) {

	var data []objects.HistoryBalanceSaldo
	err := con.db.Scopes(SchemaFrm("form_organizations")).
		Select("form_organizations.*, o.name as organization_name, f.name as form_name").
		Joins("join frm.forms f on f.id=form_organizations.form_id").
		Joins("join mstr.organizations o on o.id=form_organizations.organization_id").
		Where(fields).Find(&data).Error
	if err != nil {
		return []objects.HistoryBalanceSaldo{}, err
	}

	return data, err
}

func (con *formConnection) GetFormsOrganizationWithFilter(fields objects.HistoryBalanceSaldo, whereString string) ([]objects.HistoryBalanceSaldo, error) {

	var data []objects.HistoryBalanceSaldo
	err := con.db.Scopes(SchemaFrm("form_organizations")).
		Select("form_organizations.*, o.name as organization_name, f.name as form_name, f.profile_pic as form_image").
		Joins("join frm.forms f on f.id=form_organizations.form_id").
		Joins("join mstr.organizations o on o.id=form_organizations.organization_id").
		Where(fields).Where(whereString).Find(&data).Error
	if err != nil {
		return []objects.HistoryBalanceSaldo{}, err
	}

	return data, err
}

// func (tables.InputForms) TableName(increment string) string {
// 	return "frm.input_my_forms_" + increment
// }

func (con *formConnection) GetFormRows(fields tables.Forms) ([]tables.FormAll, error) {

	var data []tables.FormAll

	// err := con.db.Scopes(SchemaFrm("forms")).Where(fields).Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Find(&data).Error
	err := con.db.Table("frm.forms").Select("forms.id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, u.name as created_by_name, u.email as created_by_email,forms.is_attendance_required, forms.submission_target_user").Where(fields).Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u on u.id = forms.created_by").Order("forms.name asc, forms.created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetAllFormRows(fields tables.Forms) ([]tables.FormAll, error) {

	var data []tables.FormAll

	// err := con.db.Scopes(SchemaFrm("forms")).Find(&data).Error
	err := con.db.Table("frm.forms").Select("forms.id").Order("forms.id asc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormWhreRows(fields tables.Forms, whreStr string) ([]tables.FormAll, error) {

	var data []tables.FormAll

	err := con.db.Table("frm.forms").Select("forms.id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, u.name as created_by_name, u.email as created_by_email").Where(fields).Where(whreStr).Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u on u.id = forms.created_by").Joins("left join frm.form_organizations fo on fo.form_id = forms.id").Order("forms.name asc, forms.created_at desc").Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormNotInProjectRows(fields tables.Forms, whereString string, paging objects.Paging) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.FormAll

	err := con.db.Table("frm.forms").Select("forms.id, fo.organization_id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, forms.is_attendance_required, u.name as created_by_name, u.email as created_by_email,forms.created_by, forms.share_url, forms.created_at, forms.updated_at, forms.archived_at, forms.submission_target_user").Where(fields).Where(whereString).
		Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").
		Joins("left join usr.users u on u.id = forms.created_by").
		Joins("left join frm.form_organizations fo on fo.form_id = forms.id").
		Order(orderBy).Order("forms.name asc, forms.created_at desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) __GetFormNotInProjectRows(fields tables.Forms, whereString string, paging objects.Paging) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	var data []tables.FormAll
	fmt.Println(offset)
	// err := con.db.Scopes(SchemaFrm("forms")).Where(fields).Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Find(&data).Error
	err := con.db.Table("frm.forms").Select("forms.id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, forms.is_attendance_required, u.name as created_by_name, u.email as created_by_email, forms.created_at, forms.updated_at, forms.archived_at").Where(fields).Where(whereString).Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u on u.id = forms.created_by").Order("forms.name asc, forms.created_at desc").Find(&data).Error

	if err != nil {
		fmt.Println("-----err", data)
		return nil, err
	}
	fmt.Println("-----gooo", data)
	return data, nil
}

func (con *formConnection) GetFormOwnerRows(fields tables.FormOrganizationsJoin, whereString string, paging objects.Paging) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.FormAll

	err := con.db.Table("frm.form_organizations").Select("forms.id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, forms.is_attendance_required, u.name as created_by_name, u.email as created_by_email,forms.share_url, forms.created_by, forms.created_at, forms.updated_at, coalesce(forms.archived_at, null) as archived_at, (case when u.id=(select org.created_by from mstr.organizations org where org.id = form_organizations.organization_id) then 'Owner' else 'Member' end) as author, forms.submission_target_user, form_organizations.organization_id").
		Joins("left join frm.forms on forms.id = form_organizations.form_id").
		Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u on u.id = forms.created_by").
		Where(fields).Where(whereString).Order(orderBy).Order("forms.name asc, forms.created_at desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) GetFormRow(fields tables.Forms) (tables.FormOut, error) {
	var data tables.FormOut
	err := con.db.Scopes(SchemaFrm("forms")).Select("forms.*, fs.name as form_status, u.phone as user_phone, fo.organization_id, o.name as organization_name").Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u ON u.id = forms.created_by").Joins("join frm.form_organizations fo ON fo.form_id=forms.id").Joins("join mstr.organizations o on o.id=fo.organization_id").Where(fields).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.FormOut{}, err
	}

	if err != nil {
		return tables.FormOut{}, err
	}
	return data, nil
}

func (con *formConnection) GetFormPeriodeRangeRow(fields tables.Forms) (tables.FormPeriodRange, error) {
	var data tables.FormPeriodRange
	err := con.db.Table("frm.forms").Select("extract(day from forms.period_end_date::timestamp - forms.period_start_date::timestamp) as period_range").Where(fields).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.FormPeriodRange{}, err
	}

	if err != nil {
		return tables.FormPeriodRange{}, err
	}

	return data, nil
}

func (con *formConnection) UpdateForm(id int, fields tables.Forms) (bool, error) {

	fmt.Println("fields ProfilePic", fields.ProfilePic)
	err := con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Update("is_attendance_required", fields.IsAttendanceRequired).Error
	if err != nil {
		return false, err
	}

	err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Update("is_attendance_radius", fields.IsAttendanceRadius).Error
	if err != nil {
		return false, err
	}

	if fields.PeriodEndDate == "" {
		err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Update("period_end_date", nil).Error
		if err != nil {
			return false, err
		}
	}

	if fields.ProfilePic == "" {
		err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Update("profile_pic", nil).Error
		if err != nil {
			return false, err
		}
	}

	if fields.AttendanceOverdateAt.IsZero() == false {
		err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Update("attendance_overdate_at", fields.AttendanceOverdateAt).Error
		if err != nil {
			return false, err
		}
	} else {
		err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", id).Update("attendance_overdate_at", nil).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (con *formConnection) UpdateFormStatus(id int, fields tables.Forms) (bool, error) {

	fmt.Println("fields ProfilePic", fields.ProfilePic)
	err := con.db.Scopes(SchemaFrm("forms")).Where("id", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *formConnection) UpdateFormAttendanceLocation(id int, fields objects.ObjectFormAttendanceLocations) (bool, error) {

	err := con.db.Raw(`update frm.form_attendance_locations set 
	name=?, 
	location=?,
	geometry=ST_SetSRID(ST_MakePoint(?, ?), 4326),
	is_check_in=?,
	is_check_out=?,
	radius=?
	where id=?`,
		fields.Name,
		fields.Location,
		fields.Longitude,
		fields.Latitude,
		fields.IsCheckIn,
		fields.IsCheckOut,
		fields.Radius,
		fields.ID).Scan(&fields).Error
	return true, err
}

func (con *formConnection) DeleteForm(formID int, authorID int) (bool, error) {
	var forms tables.Forms
	err := con.db.Scopes(SchemaFrm("forms")).Delete(&forms, formID).Error
	if err != nil {
		return false, err
	}

	err = con.db.Scopes(SchemaFrm("forms")).Where("id = ?", formID).Update("deleted_by", authorID).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formConnection) DeleteFormLocation(ID int) (bool, error) {
	var formslocation objects.ObjectFormAttendanceLocations
	err := con.db.Scopes(SchemaFrm("form_attendance_locations")).Delete(&formslocation, ID).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formConnection) InsertFormUser(data objects.Forms) (tables.Forms, error) {

	var fieldForms tables.Forms
	err := con.db.Scopes(SchemaFrm("forms")).Create(&fieldForms).Error
	fmt.Println("data.ID---", data.ID)

	con.db.Exec(`CREATE TABLE "frm"."input_forms_` + strconv.Itoa(fieldForms.ID) + `" ("id" serial,"user_id" integer NOT NULL,"created_by" integer DEFAULT null,"updated_by" integer DEFAULT null,"deleted_by" integer DEFAULT null,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,PRIMARY KEY ("id"))`)
	con.db.Exec(`ALTER TABLE "frm"."input_forms_` + strconv.Itoa(fieldForms.ID) + `" ADD CONSTRAINT "frm"."input_forms_` + strconv.Itoa(fieldForms.ID) + `"_user_id_foreign FOREIGN KEY (user_id) REFERENCES users(id) MATCH FULL;`)

	var fieldFormUser tables.FormUsers
	fieldFormUser.FormID = fieldForms.ID
	fieldFormUser.UserID = data.UserID
	err = con.db.Scopes(SchemaFrm("form_users")).Create(&fieldFormUser).Error

	return fieldForms, err
}

func (con *formConnection) ConnectFormUser(data tables.FormUsers) (tables.FormUsers, error) {
	err := con.db.Scopes(SchemaFrm("form_users")).Create(&data).Error
	if err != nil {
		return tables.FormUsers{}, err
	}

	return data, err
}

func (con *formConnection) ConnectFormUserOrg(dataInput objects.InputFormUserOrganizations) (tables.FormUsers, error) {

	var data tables.FormUsers
	data.FormID = dataInput.FormID
	data.UserID = dataInput.UserID
	data.Type = dataInput.Type
	data.FormUserStatusID = 1
	err := con.db.Scopes(SchemaFrm("form_users")).Create(&data).Error
	if err != nil {
		return tables.FormUsers{}, err
	}

	if dataInput.OrganizationID >= 1 {
		var dataFormUserOrg tables.FormUserOrganizations
		dataFormUserOrg.FormUserID = data.ID
		dataFormUserOrg.OrganizationID = dataInput.OrganizationID
		err = con.db.Scopes(SchemaFrm("form_user_organizations")).Create(&dataFormUserOrg).Error
		if err != nil {
			return tables.FormUsers{}, err
		}
	}

	return data, err
}

func (con *formConnection) DeleteFormUser(userID int, formID int, userType string) (bool, error) {
	// var forms tables.FormUsers
	// err := con.db.Unscoped().Delete(&forms, userID).Error
	//err := con.db.Scopes(SchemaFrm("projects")).Delete(&data, projectID).Error

	// err := con.db.Scopes(SchemaFrm("form_users")).Delete(&forms, ID).Error
	// err := con.db.Exec("DELETE FROM frm.form_users where id = " + strconv.Itoa(formUserID)).Error
	err := con.db.Exec("DELETE FROM frm.form_users where user_id = " + strconv.Itoa(userID) + " and form_id = " + strconv.Itoa(formID) + " and type ='" + userType + "'").Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formConnection) DeleteFormUserOrg(userID int, formID int) (bool, error) {

	//get formuser id
	var data tables.FormUsers
	err := con.db.Table("frm.form_users").Where(tables.FormUsers{UserID: userID, FormID: formID}).First(&data).Error
	if err != nil {
		fmt.Println("----form_users", err.Error())
		return false, err
	}

	err = con.db.Exec("DELETE FROM usr.form_user_permissions where form_user_id = " + strconv.Itoa(data.ID)).Error
	if err != nil {
		fmt.Println("----form_users_permission", err.Error())
		return false, err
	}

	err = con.db.Exec("DELETE FROM frm.form_user_organizations where form_user_id = (select fu.id from frm.form_users fu where fu.form_id=" + strconv.Itoa(formID) + " AND fu.user_id = " + strconv.Itoa(userID) + ")").Error
	if err != nil {
		fmt.Println("----form_users_organization", err.Error())
		return false, err
	}

	// get organization id
	var formOrg tables.FormOrganizations
	err = con.db.Table("frm.form_organizations").Where(tables.FormOrganizations{FormID: formID}).First(&formOrg).Error
	if err != nil {
		fmt.Println("----form_organization", err.Error())
		return false, err
	}

	// err = con.db.Exec("DELETE FROM usr.user_organization_roles where user_organization_id = () " + strconv.Itoa(formOrg.ID)).Error
	// if err != nil {
	// 	fmt.Println("----user_organization_roles", err.Error())
	// 	return false, err
	// }

	// err = con.db.Exec("DELETE FROM usr.user_organizations where user_id = " + strconv.Itoa(userID) + " and organization_id = " + strconv.Itoa(formOrg.OrganizationID)).Error
	// if err != nil {
	// 	fmt.Println("----user_organizations", err.Error())
	// 	return false, err
	// }

	err = con.db.Exec("DELETE FROM frm.form_users where user_id = " + strconv.Itoa(userID) + " and form_id = " + strconv.Itoa(formID)).Error
	if err != nil {
		fmt.Println("----form_users", err.Error())
		return false, err
	}

	return true, err
}

func (con *formConnection) UpdateFormUser(userID int, formID int, data tables.FormUsers) (bool, error) {
	err := con.db.Scopes(SchemaFrm("form_users")).Where("user_id = ? and form_id = ?", userID, formID).Updates(&data).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formConnection) GetFormUserRows(fields tables.JoinFormUsers, whereStr string) ([]tables.JoinFormUsers, error) {

	fmt.Println("------ctx----", con.ctx)
	// claims := jwt.ExtractClaims(con.ctx)
	// userID := claims["id"].(string)
	// fmt.Println(userID)

	var data []tables.JoinFormUsers
	err := con.db.Table("frm.form_users").
		Select("form_users.id, form_users.form_id, fus.id as form_user_status_id, t.translation as form_user_status_name, u.id as user_id, u.name as user_name, u.email, u.phone, u.avatar as user_image, f.name, f.description, f.notes, f.profile_pic, f.period_start_date, f.period_end_date, s.id as form_status_id, s.name as form_status, f.created_at").
		Joins("join frm.forms f on f.id=form_users.form_id").
		Joins("left join usr.users u on u.id = form_users.user_id").
		Joins("left join frm.form_user_organizations fuo on form_users.id = fuo.form_user_id").
		Joins("left join mstr.form_statuses s on s.id = f.form_status_id").
		Joins("left join mstr.form_user_statuses fus on fus.id = form_users.form_user_status_id").
		Where(fields).Where(whereStr).Order("form_users.id desc").
		Joins("left join mstr.translations t on t.textcontent_id = fus.name_textcontent_id AND t.language_id = (select ul.language_id from usr.user_languages ul where ul.user_id = 5)").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormUserToOrganizationRows(fields tables.JoinFormUsers, whereStr string) ([]tables.JoinFormUsers, error) {

	var data []tables.JoinFormUsers
	err := con.db.Table("frm.form_users").Select("form_users.id, form_users.form_id, fus.id as form_user_status_id, t.translation as form_user_status_name, u.id as user_id, u.name as user_name, u.email, u.phone, f.name, f.description, f.notes, f.profile_pic, f.period_start_date, f.period_end_date, s.id as form_status_id, s.name as form_status, f.created_at").Joins("join frm.forms f on f.id=form_users.form_id").Joins("left join usr.users u on u.id = form_users.user_id").Joins("left join mstr.form_statuses s on s.id = f.form_status_id").Joins("left join mstr.form_user_statuses fus on fus.id = form_users.form_user_status_id").Where(fields).Where(whereStr).Order("form_users.id desc").Joins("left join mstr.translations t on t.textcontent_id = fus.name_textcontent_id AND t.language_id = (select ul.language_id from usr.user_languages ul where ul.user_id = 5)").Joins("left join frm.form_user_organizations fuo on fuo.form_user_id=form_users.id").Joins("join frm.form_organizations fo on fo.form_id=form_users.form_id").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormUserToOrganizationRow(fields tables.JoinFormUsers, whereStr string) (tables.JoinFormUsers, error) {

	var data tables.JoinFormUsers
	err := con.db.Table("frm.form_users").
		Select("form_users.id, form_users.form_id, fus.id as form_user_status_id, t.translation as form_user_status_name, u.id as user_id, u.name as user_name, u.email, u.phone, f.name, f.description, f.notes, f.profile_pic, f.period_start_date, f.period_end_date, s.id as form_status_id, s.name as form_status, f.created_at, fuo.organization_id").
		Joins("join frm.forms f on f.id=form_users.form_id").Joins("left join usr.users u on u.id = form_users.user_id").
		Joins("left join mstr.form_statuses s on s.id = f.form_status_id").
		Joins("left join mstr.form_user_statuses fus on fus.id = form_users.form_user_status_id").
		Joins("left join mstr.translations t on t.textcontent_id = fus.name_textcontent_id AND t.language_id = (select ul.language_id from usr.user_languages ul where ul.user_id = 5)").
		Joins("left join frm.form_user_organizations fuo on fuo.form_user_id=form_users.id").
		Joins("join frm.form_organizations fo on fo.form_id=form_users.form_id").
		Where(fields).Where(whereStr).
		Order("form_users.id desc").
		First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.JoinFormUsers{}, nil
	}

	if err != nil {
		return tables.JoinFormUsers{}, err
	}
	return data, nil
}

func (con *formConnection) GetFormUserRespondenRows(respondenID int) ([]tables.JoinFormUsers, error) {

	var data []tables.JoinFormUsers
	err := con.db.Raw(`select f.id as form_id , f.is_attendance_required
						from frm.forms f 
						where f.id in (SELECT fu.form_id from frm.form_users fu
										where fu.form_user_status_id = 1
										AND (fu.user_id = ? or fu.form_id in (select ft.form_id FROM frm.form_teams ft
																				where ft.team_id in ( select tu.team_id from usr.team_users tu 
																									where tu.user_id = ?))
										)
										and fu.type='respondent'
										)
						and f.deleted_at is null
						and f.form_status_id = 1`, respondenID, respondenID).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormUserRows___(fields tables.FormUsers) ([]tables.JoinFormUsers, error) {
	var data []tables.JoinFormUsers
	err := con.db.Table("frm.form_users").Select("form_users.id, form_users.form_id, u.id as user_id, u.name, u.email, u.phone, f.name as form_name, f.description, f.notes, f.profile_pic, f.period_start_date, f.period_end_date, f.created_at").Joins("join frm.forms f on f.id=form_users.form_id").Joins("left join usr.users u on u.id = form_users.user_id").Where(fields).Order("form_users.id desc").Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *formConnection) GetFormUserUnionTeamRows(fields tables.FormUsers) ([]tables.JoinFormUsers, error) {
	var data []tables.JoinFormUsers

	err := con.db.Raw(`SELECT * from
    (SELECT form_users.form_id, u.id as user_id, u.name as user_name, u.email, u.phone, f.name, f.description, f.notes, f.profile_pic, 
    f.period_start_date, f.period_end_date, f.is_attendance_radius, s.id as form_status_id, s.name as form_status, f.created_at,f.updated_at, a.attendance_in, a.attendance_out, f.is_attendance_required
	, (case when f.attendance_overdate_at is not null AND f.is_attendance_required is true
		then true else false end) as is_attendance_overdate , fo.organization_id
    FROM "frm"."form_users" 
	join frm.forms f on f.id=form_users.form_id 
	left join usr.users u on u.id = form_users.user_id
    left join mstr.form_statuses s on s.id = f.form_status_id 
    left join usr.attendances a on a.form_id = f.id AND a.user_id = u.id  AND to_char(a.created_at::date, 'yyyy-mm-dd') = to_char(now(), 'yyyy-mm-dd')
	join frm.form_organizations fo on fo.form_id = f.id
	left join usr.organization_subscription_plans usp on usp.organization_id = fo.organization_id 

    where "form_users"."user_id" = ` + strconv.Itoa(fields.UserID) + `
	AND to_char(now()::date, 'yyyy-mm-dd') BETWEEN to_char(f.period_start_date::date, 'yyyy-mm-dd') AND to_char(COALESCE(f.period_end_date::date, '2050-01-01'::date), 'yyyy-mm-dd')
	AND "f"."form_status_id" = 1
	AND form_users.form_user_status_id = 1
	AND form_users.type = 'respondent'
    AND "form_users"."deleted_at" IS NULL 
	AND usp.is_blocked is false) as tbl_1
	    
    UNION
    
    (SELECT form_teams.form_id, u.id as user_id, u.name as user_name, u.email, u.phone, f.name, f.description, f.notes, f.profile_pic, 
    f.period_start_date, f.period_end_date,f.is_attendance_radius, s.id as form_status_id, s.name as form_status, f.created_at,f.updated_at, a.attendance_in, a.attendance_out, f.is_attendance_required
	, (case when f.attendance_overdate_at is not null AND f.is_attendance_required is true
		then true else false end) as is_attendance_overdate, fo.organization_id

    FROM "frm"."form_teams" 
    left join frm.forms f on f.id=form_teams.form_id 
    left join usr.teams t on t.id = form_teams.team_id
    left join usr.team_users tu on tu.team_id = t.id
    left join usr.users u on u.id = tu.user_id
    left join mstr.form_statuses s on s.id = f.form_status_id 
    left join usr.attendances a on a.form_id = f.id AND a.user_id = u.id AND to_char(a.created_at::date, 'yyyy-mm-dd') = to_char(now(), 'yyyy-mm-dd')
    join frm.form_organizations fo on fo.form_id = f.id
	left join usr.organization_subscription_plans usp on usp.organization_id = fo.organization_id

    where "tu"."user_id" = ` + strconv.Itoa(fields.UserID) + `
	AND to_char(now()::date, 'yyyy-mm-dd') BETWEEN to_char(f.period_start_date::date, 'yyyy-mm-dd') AND to_char(COALESCE(f.period_end_date::date, '2050-01-01'::date), 'yyyy-mm-dd')
	AND "f"."form_status_id" = 1
    AND "form_teams"."deleted_at" IS NULL 
	AND usp.is_blocked is false)

	ORDER BY name asc
    `).Find(&data).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *formConnection) GetFormUserAdminRows(fields tables.FormUsers) ([]tables.JoinFormUsers, error) {
	var data []tables.JoinFormUsers
	return data, nil
}

func (con *formConnection) GetFormUserRow(fields tables.FormUsers) (tables.JoinFormUsers, error) {
	var data tables.JoinFormUsers
	err := con.db.Table("frm.form_users").Select("form_users.id, form_users.form_id, 1 as form_user_status_id, 'Active' as form_user_status_name,  f.name, f.description, f.notes, f.profile_pic, f.period_start_date, f.period_end_date, f.created_at").Joins("join frm.forms f on f.id=form_users.form_id").Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *formConnection) InsertFieldFile(data tables.FormFieldTempAssets) (bool, tables.FormFieldTempAssets, error) {

	err := con.db.Scopes(SchemaFrm("form_field_temp_assets")).Create(&data).Error
	if err != nil {
		return false, tables.FormFieldTempAssets{}, err
	}

	return true, data, nil
}

func (con *formConnection) GetFormUserUniqRows(authorID int, projectID int) ([]tables.FormUsers, error) {
	var data []tables.FormUsers

	whreProject := ``
	if projectID > 0 {
		whreProject = "and f.id in (select pf.form_id from frm.project_forms pf where pf.project_id = " + strconv.Itoa(projectID) + ")"
	}
	err := con.db.Raw(`select fu.user_id 
						from frm.form_users fu
						where fu.form_id in (select f.id from frm.forms f
							where f.created_by = ?
							`+whreProject+`
							and f.deleted_at is null
							and f.form_status_id = 1)
						and type = 'respondent'
						group by fu.user_id
						
						UNION
						
						(select tu.user_id from usr.team_users tu
						where tu.team_id in (select ft.team_id from frm.form_teams ft 
						                      where ft.form_id in (select f.id from frm.forms f
							                                        where f.created_by = ?
																	`+whreProject+`
							                                        and f.deleted_at is null
							                                        and f.form_status_id = 1))
						) `, authorID, authorID).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *formConnection) GetFormUnionProjectRows(fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`SELECT * from 
							(SELECT t.id as project_id
								, 0 as id,
								0 as form_status_id,
								'' as form_status,
								t.name,
								t.description,
								'' as notes,
								now() as period_start_date,
								now() as period_end_date,
								'' as profile_pic,
								false as is_publish,
								false as is_attendance_required,
								'' as created_by_name,
								'' as created_by_email,
								'' as share_url,
								t.created_by,
								t.created_at,
								t.updated_at,
								now() as archived_at,
								'' as author,
								0 as submission_target_user,
								po.organization_id
							FROM frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + `
							ORDER BY t.id) as tb_2
							
							UNION (SELECT 
									0 as project_id,
									forms.id,
								forms.form_status_id,
								fs.name as form_status,
								forms.name,
								forms.description,
								forms.notes,
								forms.period_start_date,
								forms.period_end_date,
								forms.profile_pic,
								forms.is_publish,
								forms.is_attendance_required,
								u.name as created_by_name,
								u.email as created_by_email,
								forms.share_url,
								forms.created_by,
								forms.created_at,
								forms.updated_at,
								coalesce(forms.archived_at, null) as archived_at,
								(case
										when u.id=
											(select org.created_by
												from mstr.organizations org
												where org.id = form_organizations.organization_id) then 'Owner'
										else 'Member'
									end) as author,
								forms.submission_target_user,
								form_organizations.organization_id
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id not in (select pf.form_id from frm.project_forms pf)
							) 
							
							order by  ` + orderBy + ` project_id desc, name asc
							OFFSET ` + strconv.Itoa(offset) + ` limit ` + strconv.Itoa(paging.Limit) + `	
    						`).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeSuperAdmin(fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`SELECT * from 
							(SELECT
								0 as form_share_count, 
								t.id as project_id, 
								0 as id,
								0 as form_status_id,
								'' as form_status,
								t.name,
								t.description,
								'' as notes,
								now() as period_start_date,
								now() as period_end_date,
								'' as profile_pic,
								false as is_publish,
								false as is_attendance_required,
								'' as created_by_name,
								'' as created_by_email,
								'' as share_url,
								t.created_by,
								t.created_at,
								t.updated_at,
								now() as archived_at,
								'' as author,
								0 as submission_target_user,
								po.organization_id,
								'group' as form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'' as type
							FROM frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + ` ) as tb_2
							
							UNION (SELECT 
								(SELECT count(t.id) as count_share
								FROM frm.form_organization_invites t
								where form_id =forms.id) as form_share_count,
								0 as project_id,
								forms.id,
								forms.form_status_id,
								fs.name as form_status,
								forms.name,
								forms.description,
								forms.notes,
								forms.period_start_date,
								forms.period_end_date,
								forms.profile_pic,
								forms.is_publish,
								forms.is_attendance_required,
								u.name as created_by_name,
								u.email as created_by_email,
								forms.share_url,
								forms.created_by,
								forms.created_at,
								forms.updated_at,
								coalesce(forms.archived_at, null) as archived_at,
								(case
										when u.id=
											(select org.created_by
												from mstr.organizations org
												where org.id = form_organizations.organization_id) then 'Owner'
										else 'Member'
									end) as author,
								forms.submission_target_user,
								form_organizations.organization_id,
								(CASE 
									WHEN (
									SELECT 
										COUNT(t.id) AS count_share 
									FROM 
										frm.form_organization_invites t 
									WHERE 
										t.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END) AS form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'internal' as type
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id not in (select pf.form_id from frm.project_forms pf)
							) 
							UNION
							(
								SELECT 
									0 as form_share_count,
									0 as project_id, 
									forms.id, 
									forms.form_status_id, 
									fs.name as form_status, 
									forms.name, 
									forms.description, 
									forms.notes, 
									forms.period_start_date, 
									forms.period_end_date, 
									forms.profile_pic, 
									forms.is_publish, 
									forms.is_attendance_required, 
									u.name as created_by_name, 
									u.email as created_by_email, 
									forms.share_url, 
									forms.created_by, 
									forms.created_at, 
									forms.updated_at, 
									coalesce(forms.archived_at, null) as archived_at, 
									(
										case when u.id =(
										select 
											org.created_by 
										from 
											mstr.organizations org 
										where 
											org.id = form_organization_invites.organization_id
										) then 'Owner' else 'Member' end
									) as author, 
									forms.submission_target_user, 
									form_organization_invites.organization_id,
									'in' as form_shared,
									o.name as form_external_company_name,
									o.profile_pic as form_external_company_image,
									'external' as type									
								FROM 
									"frm"."form_organization_invites" 
									left join frm.forms on forms.id = form_organization_invites.form_id 
									left join frm.form_organizations fo on fo.form_id = form_organization_invites.form_id
									left join mstr.organizations o on o.id = fo.organization_id
									left join mstr.form_statuses fs on fs.id = forms.form_status_id 
									left join usr.users u on u.id = forms.created_by 
								WHERE "form_organization_invites"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								` + whereString + `	
							AND "form_organization_invites"."deleted_at" IS NULL
							)
							order by ` + orderBy + ` project_id desc, name asc
							OFFSET ` + strconv.Itoa(offset) + ` limit ` + strconv.Itoa(paging.Limit) + `	
    						`).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeSuperAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`WITH subquery1 AS (
							SELECT 
								0 AS form_share_count, 
								t.id AS project_id, 
								0 AS id, 
								0 AS form_status_id, 
								'' AS form_status, 
								t.name, 
								t.description, 
								'' AS notes, 
								NOW() AS period_start_date, 
								NOW() AS period_end_date, 
								'' AS profile_pic, 
								FALSE AS is_publish, 
								FALSE AS is_attendance_required, 
								'' AS created_by_name, 
								'' AS created_by_email, 
								'' AS share_url, 
								t.created_by, 
								t.created_at, 
								t.updated_at, 
								NOW() AS archived_at, 
								'' AS author, 
								0 AS submission_target_user, 
								po.organization_id, 
								'group' AS form_shared, 
								'' AS form_external_company_name, 
								'' AS form_external_company_image, 
								'' AS type 
							FROM 
								frm.projects t 
							LEFT JOIN frm.project_organizations po ON po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + ` 
						),
						subquery2 AS (
							SELECT 
								(
								SELECT 
									COUNT(t.id) AS count_share 
								FROM 
									frm.form_organization_invites t 
								WHERE 
									t.form_id = forms.id
								) AS form_share_count, 
								0 AS project_id, 
								forms.id, 
								forms.form_status_id, 
								fs.name AS form_status, 
								forms.name, 
								forms.description, 
								forms.notes, 
								forms.period_start_date, 
								forms.period_end_date, 
								forms.profile_pic, 
								forms.is_publish, 
								forms.is_attendance_required, 
								u.name AS created_by_name, 
								u.email AS created_by_email, 
								forms.share_url, 
								forms.created_by, 
								forms.created_at, 
								forms.updated_at, 
								COALESCE(forms.archived_at, NULL) AS archived_at, 
								(
								CASE 
									WHEN u.id = (
									SELECT 
										org.created_by 
									FROM 
										mstr.organizations org 
									WHERE 
										org.id = form_organizations.organization_id
									) THEN 'Owner' 
									ELSE 'Member' 
								END
								) AS author, 
								forms.submission_target_user, 
								form_organizations.organization_id, 
								(
								CASE 
									WHEN (
									SELECT 
										COUNT(t.id) AS count_share 
									FROM 
										frm.form_organization_invites t 
									WHERE 
										t.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END
								) AS form_shared, 
								'' AS form_external_company_name, 
								'' AS form_external_company_image, 
								'internal' AS type 
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id NOT IN (
							SELECT 
								pf.form_id 
							FROM 
								frm.project_forms pf
							)
						), 
						subquery3 AS (
							SELECT 
								0 AS form_share_count, 
								0 AS project_id, 
								forms.id, 
								forms.form_status_id, 
								fs.name AS form_status, 
								forms.name, 
								forms.description, 
								forms.notes, 
								forms.period_start_date, 
								forms.period_end_date, 
								forms.profile_pic, 
								forms.is_publish, 
								forms.is_attendance_required, 
								u.name AS created_by_name, 
								u.email AS created_by_email, 
								forms.share_url, 
								forms.created_by, 
								forms.created_at, 
								forms.updated_at, 
								COALESCE(forms.archived_at, NULL) AS archived_at, 
								(
								CASE 
									WHEN u.id =(
									SELECT 
										org.created_by 
									FROM 
										mstr.organizations org 
									WHERE 
										org.id = form_organization_invites.organization_id
									) THEN 'Owner' 
									ELSE 'Member' 
								END
								) AS author, 
								forms.submission_target_user, 
								form_organization_invites.organization_id, 
								'in' AS form_shared, 
								o.name AS form_external_company_name, 
								o.profile_pic AS form_external_company_image, 
								'external' AS type 								
								FROM 
									"frm"."form_organization_invites" 
									left join frm.forms on forms.id = form_organization_invites.form_id 
									left join frm.form_organizations fo on fo.form_id = form_organization_invites.form_id
									left join mstr.organizations o on o.id = fo.organization_id
									left join mstr.form_statuses fs on fs.id = forms.form_status_id 
									left join usr.users u on u.id = forms.created_by 
								WHERE "form_organization_invites"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								` + whereString + `	
								AND "form_organization_invites"."deleted_at" IS NULL
							)
							SELECT 
								form_share_count,
								project_id,
								id,
								form_status_id,
								form_status,
								name,
								description,
								notes,
								period_start_date,
								period_end_date,
								profile_pic,
								is_publish,
								is_attendance_required,
								created_by_name,
								created_by_email,
								share_url,
								created_by,
								created_at,
								updated_at,
								archived_at,
								author,
								submission_target_user,
								organization_id,
								form_shared,
								form_external_company_name,
								form_external_company_image,
								type
							FROM (
								SELECT * FROM subquery1
								UNION ALL
								SELECT * FROM subquery2
								UNION ALL
								SELECT * FROM subquery3
							) AS tb_2							
							` + whereName + `
							ORDER BY 
								(CASE WHEN form_shared = 'group' THEN 1 ELSE 2 END), 
								` + orderBy + ` 
								project_id desc,
								name asc
							OFFSET ` + strconv.Itoa(offset) + ` 
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeSuperAdminNew1(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`WITH subquery1 AS (
							SELECT 
								0 AS form_share_count, 
								t.id AS project_id, 
								0 AS id, 
								0 AS form_status_id, 
								'' AS form_status, 
								t.name, 
								t.description, 
								'' AS notes, 
								NOW() AS period_start_date, 
								NOW() AS period_end_date, 
								'' AS profile_pic, 
								FALSE AS is_publish, 
								FALSE AS is_attendance_required, 
								'' AS created_by_name, 
								'' AS created_by_email, 
								'' AS share_url, 
								t.created_by, 
								t.created_at, 
								t.updated_at, 
								NOW() AS archived_at, 
								'' AS author, 
								0 AS submission_target_user, 
								po.organization_id, 
								'group' AS form_shared, 
								'' AS form_external_company_name, 
								'' AS form_external_company_image, 
								'' AS type, 								
								'' AS access_type,
								false AS is_quota_sharing
							FROM 
								frm.projects t 
							LEFT JOIN frm.project_organizations po ON po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + ` 
						),
						subquery2 AS (
							SELECT 
								(
								SELECT 
									COUNT(ftoi.id) AS count_share 
								FROM 
									frm.form_to_user_invites ftoi
								WHERE 
									ftoi.form_id = forms.id
								) AS form_share_count, 
								0 AS project_id, 
								forms.id, 
								forms.form_status_id, 
								fs.name AS form_status, 
								forms.name, 
								forms.description, 
								forms.notes, 
								forms.period_start_date, 
								forms.period_end_date, 
								forms.profile_pic, 
								forms.is_publish, 
								forms.is_attendance_required, 
								u.name AS created_by_name, 
								u.email AS created_by_email, 
								forms.share_url, 
								forms.created_by, 
								forms.created_at, 
								forms.updated_at, 
								COALESCE(forms.archived_at, NULL) AS archived_at, 
								(
								CASE 
									WHEN u.id = (
									SELECT 
										org.created_by 
									FROM 
										mstr.organizations org 
									WHERE 
										org.id = form_organizations.organization_id
									) THEN 'Owner' 
									ELSE 'Member' 
								END
								) AS author, 
								forms.submission_target_user, 
								form_organizations.organization_id, 
								(
								CASE 
									WHEN (
									SELECT 
										COUNT(ftoi.id) AS count_share 
									FROM 
										frm.form_to_user_invites ftoi
									WHERE 
										ftoi.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END
								) AS form_shared, 
								'' AS form_external_company_name, 
								'' AS form_external_company_image, 
								'internal' AS type,
								'' AS access_type,
								false AS is_quota_sharing
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id NOT IN (
							SELECT 
								pf.form_id 
							FROM 
								frm.project_forms pf
							)
						), 
						subquery3 AS (
							SELECT 
								0 AS form_share_count, 
								0 AS project_id, 
								forms.id, 
								forms.form_status_id, 
								fs.name AS form_status, 
								forms.name, 
								forms.description, 
								forms.notes, 
								forms.period_start_date, 
								forms.period_end_date, 
								forms.profile_pic, 
								forms.is_publish, 
								forms.is_attendance_required, 
								u.name AS created_by_name, 
								u.email AS created_by_email, 
								forms.share_url, 
								forms.created_by, 
								forms.created_at, 
								forms.updated_at, 
								COALESCE(forms.archived_at, NULL) AS archived_at, 
								(
								CASE 
									WHEN u.id =(
									SELECT 
										org.created_by 
									FROM 
										mstr.organizations org 
									WHERE 
										org.id = form_to_user_invites.organization_receiver_id
									) THEN 'Owner' 
									ELSE 'Member' 
								END
								) AS author, 
								forms.submission_target_user, 
								form_to_user_invites.organization_receiver_id, 
								'in' AS form_shared, 
								o.name AS form_external_company_name, 
								o.profile_pic AS form_external_company_image, 
								'external' AS type,
								form_to_user_invites.access_type,
								form_to_user_invites.is_quota_sharing 								
								FROM 
									"frm"."form_to_user_invites" 
									left join frm.forms on forms.id = form_to_user_invites.form_id 
									left join frm.form_organizations fo on fo.form_id = form_to_user_invites.form_id
									left join mstr.organizations o on o.id = fo.organization_id
									left join mstr.form_statuses fs on fs.id = forms.form_status_id 
									left join usr.users u on u.id = forms.created_by 
								WHERE "form_to_user_invites"."organization_receiver_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								` + whereString + `	
								AND "form_to_user_invites"."deleted_at" IS NULL
								GROUP by
									form_share_count,
									project_id,
									forms.id, 
									forms.form_status_id, 
									fs.name, 
									forms.name, 
									forms.description, 
									forms.notes, 
									forms.period_start_date, 
									forms.period_end_date, 
									forms.profile_pic, 
									forms.is_publish, 
									forms.is_attendance_required, 
									u.name, 
									u.email, 
									forms.share_url, 
									forms.created_by, 
									forms.created_at, 
									forms.updated_at, 
									author, 
									forms.submission_target_user, 
									form_to_user_invites.organization_receiver_id, 
									form_shared,
									o.name, 
									o.profile_pic,
									type,
									form_to_user_invites.access_type,
									form_to_user_invites.is_quota_sharing 									
							)
							SELECT 
								form_share_count,
								project_id,
								id,
								form_status_id,
								form_status,
								name,
								description,
								notes,
								period_start_date,
								period_end_date,
								profile_pic,
								is_publish,
								is_attendance_required,
								created_by_name,
								created_by_email,
								share_url,
								created_by,
								created_at,
								updated_at,
								archived_at,
								author,
								submission_target_user,
								organization_id,
								form_shared,
								form_external_company_name,
								form_external_company_image,
								type,								
								access_type,
								is_quota_sharing 
							FROM (
								SELECT * FROM subquery1
								UNION ALL
								SELECT * FROM subquery2
								UNION ALL
								SELECT * FROM subquery3
							) AS tb_2							
							` + whereName + `
							ORDER BY 
								(CASE WHEN form_shared = 'group' THEN 1 ELSE 2 END), 
								` + orderBy + ` 
								project_id desc,
								name asc
							OFFSET ` + strconv.Itoa(offset) + ` 
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeSuperAdminApps(fields tables.FormOrganizationsJoin, whereName string, whereString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`WITH subquery1 AS (
							SELECT 
								(
								SELECT 
									COUNT(ftoi.id) AS count_share 
								FROM 
									frm.form_to_user_invites ftoi
								WHERE 
									ftoi.form_id = forms.id
								) AS form_share_count, 
								0 AS project_id, 
								forms.id, 
								forms.form_status_id, 
								fs.name AS form_status, 
								forms.name, 
								forms.description, 
								forms.notes, 
								forms.period_start_date, 
								forms.period_end_date, 
								forms.profile_pic, 
								forms.is_publish, 
								forms.is_attendance_required, 
								u.name AS created_by_name, 
								u.email AS created_by_email, 
								forms.share_url, 
								forms.created_by, 
								forms.created_at, 
								forms.updated_at, 
								COALESCE(forms.archived_at, NULL) AS archived_at, 
								(
								CASE 
									WHEN u.id = (
									SELECT 
										org.created_by 
									FROM 
										mstr.organizations org 
									WHERE 
										org.id = form_organizations.organization_id
									) THEN 'Owner' 
									ELSE 'Member' 
								END
								) AS author, 
								forms.submission_target_user, 
								form_organizations.organization_id, 
								(
								CASE 
									WHEN (
									SELECT 
										COUNT(ftoi.id) AS count_share 
									FROM 
										frm.form_to_user_invites ftoi
									WHERE 
										ftoi.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END
								) AS form_shared, 
								'' AS form_external_company_name, 
								'' AS form_external_company_image, 
								'internal' AS type,
								'' AS access_type,
								false AS is_quota_sharing
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
						), 
						subquery2 AS (
							SELECT 
								0 AS form_share_count, 
								0 AS project_id, 
								forms.id, 
								forms.form_status_id, 
								fs.name AS form_status, 
								forms.name, 
								forms.description, 
								forms.notes, 
								forms.period_start_date, 
								forms.period_end_date, 
								forms.profile_pic, 
								forms.is_publish, 
								forms.is_attendance_required, 
								u.name AS created_by_name, 
								u.email AS created_by_email, 
								forms.share_url, 
								forms.created_by, 
								forms.created_at, 
								forms.updated_at, 
								COALESCE(forms.archived_at, NULL) AS archived_at, 
								(
								CASE 
									WHEN u.id =(
									SELECT 
										org.created_by 
									FROM 
										mstr.organizations org 
									WHERE 
										org.id = form_to_user_invites.organization_receiver_id
									) THEN 'Owner' 
									ELSE 'Member' 
								END
								) AS author, 
								forms.submission_target_user, 
								form_to_user_invites.organization_receiver_id, 
								'in' AS form_shared, 
								o.name AS form_external_company_name, 
								o.profile_pic AS form_external_company_image, 
								'external' AS type,
								form_to_user_invites.access_type,
								form_to_user_invites.is_quota_sharing 								
								FROM 
									"frm"."form_to_user_invites" 
									left join frm.forms on forms.id = form_to_user_invites.form_id 
									left join frm.form_organizations fo on fo.form_id = form_to_user_invites.form_id
									left join mstr.organizations o on o.id = fo.organization_id
									left join mstr.form_statuses fs on fs.id = forms.form_status_id 
									left join usr.users u on u.id = forms.created_by 
								WHERE "form_to_user_invites"."organization_receiver_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								` + whereString + `	
								AND "form_to_user_invites"."deleted_at" IS NULL
								GROUP by
									form_share_count,
									project_id,
									forms.id, 
									forms.form_status_id, 
									fs.name, 
									forms.name, 
									forms.description, 
									forms.notes, 
									forms.period_start_date, 
									forms.period_end_date, 
									forms.profile_pic, 
									forms.is_publish, 
									forms.is_attendance_required, 
									u.name, 
									u.email, 
									forms.share_url, 
									forms.created_by, 
									forms.created_at, 
									forms.updated_at, 
									author, 
									forms.submission_target_user, 
									form_to_user_invites.organization_receiver_id, 
									form_shared,
									o.name, 
									o.profile_pic,
									type,
									form_to_user_invites.access_type,
									form_to_user_invites.is_quota_sharing 									
							)
							SELECT 
								form_share_count,
								project_id,
								id,
								form_status_id,
								form_status,
								name,
								description,
								notes,
								period_start_date,
								period_end_date,
								profile_pic,
								is_publish,
								is_attendance_required,
								created_by_name,
								created_by_email,
								share_url,
								created_by,
								created_at,
								updated_at,
								archived_at,
								author,
								submission_target_user,
								organization_id,
								form_shared,
								form_external_company_name,
								form_external_company_image,
								type,								
								access_type,
								is_quota_sharing 
							FROM (
								SELECT * FROM subquery1
								UNION ALL
								SELECT * FROM subquery2
							) AS tb_2							
							` + whereName + `
							ORDER BY 
								(CASE WHEN form_shared = 'group' THEN 1 ELSE 2 END), 
								` + orderBy + ` 
								project_id desc,
								name asc
							OFFSET ` + strconv.Itoa(offset) + ` 
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeAdmin(fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`SELECT * from 
							(SELECT
								0 as form_share_count, 
								t.id as project_id, 
								0 as id,
								0 as form_status_id,
								'' as form_status,
								t.name,
								t.description,
								'' as notes,
								now() as period_start_date,
								now() as period_end_date,
								'' as profile_pic,
								false as is_publish,
								false as is_attendance_required,
								'' as created_by_name,
								'' as created_by_email,
								'' as share_url,
								t.created_by,
								t.created_at,
								t.updated_at,
								now() as archived_at,
								'' as author,
								0 as submission_target_user,
								po.organization_id,
								'group' as form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'' as type
							FROM frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + `) as tb_2
							
							UNION (SELECT 
								(SELECT count(t.id) as count_share
								FROM frm.form_organization_invites t
								where form_id =forms.id) as form_share_count,
								0 as project_id,
								forms.id,
								forms.form_status_id,
								fs.name as form_status,
								forms.name,
								forms.description,
								forms.notes,
								forms.period_start_date,
								forms.period_end_date,
								forms.profile_pic,
								forms.is_publish,
								forms.is_attendance_required,
								u.name as created_by_name,
								u.email as created_by_email,
								forms.share_url,
								forms.created_by,
								forms.created_at,
								forms.updated_at,
								coalesce(forms.archived_at, null) as archived_at,
								(case
										when u.id=
											(select org.created_by
												from mstr.organizations org
												where org.id = form_organizations.organization_id) then 'Owner'
										else 'Member'
									end) as author,
								forms.submission_target_user,
								form_organizations.organization_id,
								'' as form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'internal' as type
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id not in (select pf.form_id from frm.project_forms pf)
							) 
							UNION
							(
								SELECT 
									0 as form_share_count,
									0 as project_id, 
									f.id, 
									f.form_status_id, 
									fs.name as form_status, 
									f.name, 
									f.description, 
									f.notes, 
									f.period_start_date, 
									f.period_end_date, 
									f.profile_pic, 
									f.is_publish, 
									f.is_attendance_required, 
									u.name as created_by_name, 
									u.email as created_by_email, 
									f.share_url, 
									f.created_by, 
									f.created_at, 
									f.updated_at, 
									coalesce(f.archived_at, null) as archived_at, 
									'Member' as author, 
									f.submission_target_user, 
									foi.organization_id,
									'in' as form_shared,
									o.name as form_external_company_name,
									o.profile_pic as form_external_company_image,
									'external' as type
								FROM 
								    "frm"."form_users" 
									left join frm.forms f on f.id = form_users.form_id 
									left join frm.form_organization_invites foi on form_users.form_id = foi.form_id 
									left join frm.form_organizations fo on fo.form_id = form_users.form_id
									left join mstr.organizations o on o.id = fo.organization_id
									left join mstr.form_statuses fs on fs.id = f.form_status_id 
									left join usr.users u on u.id = f.created_by  
								WHERE 
								` + whereStr + `
							AND "form_users"."user_id" = ` + strconv.Itoa(userID) + `
							AND foi.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							AND "form_users"."deleted_at" IS NULL
							)
							order by  ` + orderBy + ` name asc, created_at desc
							OFFSET ` + strconv.Itoa(offset) + ` limit ` + strconv.Itoa(paging.Limit) + `	
    						`).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`WITH subquery1 AS(
							SELECT
								0 AS form_share_count, 
								t.id AS project_id, 
								0 AS id,
								0 AS form_status_id,
								'' AS form_status,
								t.name,
								t.description,
								'' AS notes,
								NOW() AS period_start_date,
								NOW() AS period_end_date,
								'' AS profile_pic,
								false AS is_publish,
								false AS is_attendance_required,
								'' AS created_by_name,
								'' AS created_by_email,
								'' AS share_url,
								t.created_by,
								t.created_at,
								t.updated_at,
								NOW() AS archived_at,
								'' AS author,
								0 AS submission_target_user,
								po.organization_id,
								'group' AS form_shared,
								'' AS form_external_company_name,
								'' AS form_external_company_image,
								'' AS type
							FROM 
								frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + `
						),
						subquery2 AS (
							SELECT
								(
								SELECT 
									COUNT(t.id) as count_share
								FROM 
									frm.form_organization_invites t
								WHERE 
									t.form_id =forms.id
								) AS form_share_count,
								0 as project_id,
								forms.id,
								forms.form_status_id,
								fs.name as form_status,
								forms.name,
								forms.description,
								forms.notes,
								forms.period_start_date,
								forms.period_end_date,
								forms.profile_pic,
								forms.is_publish,
								forms.is_attendance_required,
								u.name as created_by_name,
								u.email as created_by_email,
								forms.share_url,
								forms.created_by,
								forms.created_at,
								forms.updated_at,
								COALESCE(forms.archived_at, null) as archived_at,
								(
									CASE
									WHEN u.id = (
									SELECT 
										org.created_by
									FROM 
										mstr.organizations org
									WHERE 
										org.id = form_organizations.organization_id
									) THEN 'Owner'
									ELSE 'Member'
								END
								) AS author,								
								forms.submission_target_user,
								form_organizations.organization_id,
								(CASE 
									WHEN (
									SELECT 
										COUNT(t.id) AS count_share 
									FROM 
										frm.form_organization_invites t 
									WHERE 
										t.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END
								) AS form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'internal' as type
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id NOT IN (
							SELECT 
								pf.form_id 
							FROM 
								frm.project_forms pf
							)
						),
						subquery3 AS (
							SELECT 
								0 as form_share_count,
								0 as project_id, 
								f.id, 
								f.form_status_id, 
								fs.name as form_status, 
								f.name, 
								f.description, 
								f.notes, 
								f.period_start_date, 
								f.period_end_date, 
								f.profile_pic, 
								f.is_publish, 
								f.is_attendance_required, 
								u.name as created_by_name, 
								u.email as created_by_email, 
								f.share_url, 
								f.created_by, 
								f.created_at, 
								f.updated_at, 
								coalesce(f.archived_at, null) as archived_at, 
								'Member' as author, 
								f.submission_target_user, 
								foi.organization_id,
								'in' as form_shared,
								o.name as form_external_company_name,
								o.profile_pic as form_external_company_image,
								'external' as type
							FROM 
								"frm"."form_users" 
								left join frm.forms f on f.id = form_users.form_id 
								left join frm.form_organization_invites foi on form_users.form_id = foi.form_id 
								left join frm.form_organizations fo on fo.form_id = form_users.form_id
								left join mstr.organizations o on o.id = fo.organization_id
								left join mstr.form_statuses fs on fs.id = f.form_status_id 
								left join usr.users u on u.id = f.created_by
							WHERE ` + whereStr + `
								AND "form_users"."user_id" = ` + strconv.Itoa(userID) + `
								AND foi.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
								AND "form_users"."deleted_at" IS NULL
						)	
						SELECT 
								form_share_count,
								project_id,
								id,
								form_status_id,
								form_status,
								name,
								description,
								notes,
								period_start_date,
								period_end_date,
								profile_pic,
								is_publish,
								is_attendance_required,
								created_by_name,
								created_by_email,
								share_url,
								created_by,
								created_at,
								updated_at,
								archived_at,
								author,
								submission_target_user,
								organization_id,
								form_shared,
								form_external_company_name,
								form_external_company_image,
								type
							FROM (
								SELECT * FROM subquery1
								UNION ALL
								SELECT * FROM subquery2
								UNION ALL
								SELECT * FROM subquery3
							) AS tb_2
							` + whereName + `
							ORDER BY 	
								(CASE WHEN form_shared = 'group' THEN 1 ELSE 2 END),
								` + orderBy + ` 
								name asc, 
								created_at desc
							OFFSET ` + strconv.Itoa(offset) + `
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeAdminNew1(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`WITH subquery1 AS(
							SELECT
								0 AS form_share_count, 
								t.id AS project_id, 
								0 AS id,
								0 AS form_status_id,
								'' AS form_status,
								t.name,
								t.description,
								'' AS notes,
								NOW() AS period_start_date,
								NOW() AS period_end_date,
								'' AS profile_pic,
								false AS is_publish,
								false AS is_attendance_required,
								'' AS created_by_name,
								'' AS created_by_email,
								'' AS share_url,
								t.created_by,
								t.created_at,
								t.updated_at,
								NOW() AS archived_at,
								'' AS author,
								0 AS submission_target_user,
								po.organization_id,
								'group' AS form_shared,
								'' AS form_external_company_name,
								'' AS form_external_company_image,
								'' AS type,						
								'' AS access_type,
								false AS is_quota_sharing
							FROM 
								frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + `
						),
						subquery2 AS (
							SELECT
								(
								SELECT 
									COUNT(t.id) as count_share
								FROM 
									frm.form_organization_invites t
								WHERE 
									t.form_id =forms.id
								) AS form_share_count,
								0 as project_id,
								forms.id,
								forms.form_status_id,
								fs.name as form_status,
								forms.name,
								forms.description,
								forms.notes,
								forms.period_start_date,
								forms.period_end_date,
								forms.profile_pic,
								forms.is_publish,
								forms.is_attendance_required,
								u.name as created_by_name,
								u.email as created_by_email,
								forms.share_url,
								forms.created_by,
								forms.created_at,
								forms.updated_at,
								COALESCE(forms.archived_at, null) as archived_at,
								(
									CASE
									WHEN u.id = (
									SELECT 
										org.created_by
									FROM 
										mstr.organizations org
									WHERE 
										org.id = form_organizations.organization_id
									) THEN 'Owner'
									ELSE 'Member'
								END
								) AS author,								
								forms.submission_target_user,
								form_organizations.organization_id,
								(CASE 
									WHEN (
									SELECT 
										COUNT(t.id) AS count_share 
									FROM 
										frm.form_organization_invites t 
									WHERE 
										t.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END
								) AS form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'internal' as type,
								'' AS access_type,
								false AS is_quota_sharing
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							AND forms.id NOT IN (
							SELECT 
								pf.form_id 
							FROM 
								frm.project_forms pf
							)
						),
						subquery3 AS (
							SELECT 
								0 as form_share_count,
								0 as project_id, 
								f.id, 
								f.form_status_id, 
								fs.name as form_status, 
								f.name, 
								f.description, 
								f.notes, 
								f.period_start_date, 
								f.period_end_date, 
								f.profile_pic, 
								f.is_publish, 
								f.is_attendance_required, 
								u.name as created_by_name, 
								u.email as created_by_email, 
								f.share_url, 
								f.created_by, 
								f.created_at, 
								f.updated_at, 
								coalesce(f.archived_at, null) as archived_at, 
								'Member' as author, 
								f.submission_target_user, 
								form_to_user_invites.organization_receiver_id, 
								'in' as form_shared,
								o.name as form_external_company_name,
								o.profile_pic as form_external_company_image,
								'external' as type,
								form_to_user_invites.access_type,
								form_to_user_invites.is_quota_sharing
							FROM 
								"frm"."form_to_user_invites" 
								left join frm.forms f on f.id = form_to_user_invites.form_id 
								left join frm.form_organizations fo on fo.form_id = form_to_user_invites.form_id 
								left join mstr.organizations o on o.id = fo.organization_id 
								left join mstr.form_statuses fs on fs.id = f.form_status_id 
								left join usr.users u on u.id = f.created_by 
							WHERE ` + whereStr + `
								AND "form_to_user_invites"."user_receiver_id" = ` + strconv.Itoa(userID) + `
								AND "form_to_user_invites"."organization_receiver_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								AND "form_to_user_invites"."deleted_at" IS NULL
						)	
						SELECT 
								form_share_count,
								project_id,
								id,
								form_status_id,
								form_status,
								name,
								description,
								notes,
								period_start_date,
								period_end_date,
								profile_pic,
								is_publish,
								is_attendance_required,
								created_by_name,
								created_by_email,
								share_url,
								created_by,
								created_at,
								updated_at,
								archived_at,
								author,
								submission_target_user,
								organization_id,
								form_shared,
								form_external_company_name,
								form_external_company_image,
								type,
								access_type,
								is_quota_sharing 	
							FROM (
								SELECT * FROM subquery1
								UNION ALL
								SELECT * FROM subquery2
								UNION ALL
								SELECT * FROM subquery3
							) AS tb_2
							` + whereName + `
							ORDER BY 	
								(CASE WHEN form_shared = 'group' THEN 1 ELSE 2 END),
								` + orderBy + ` 
								name asc, 
								created_at desc
							OFFSET ` + strconv.Itoa(offset) + `
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormMergeAdminApps(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`WITH subquery1 AS (
							SELECT
								(
								SELECT 
									COUNT(t.id) as count_share
								FROM 
									frm.form_organization_invites t
								WHERE 
									t.form_id =forms.id
								) AS form_share_count,
								0 as project_id,
								forms.id,
								forms.form_status_id,
								fs.name as form_status,
								forms.name,
								forms.description,
								forms.notes,
								forms.period_start_date,
								forms.period_end_date,
								forms.profile_pic,
								forms.is_publish,
								forms.is_attendance_required,
								u.name as created_by_name,
								u.email as created_by_email,
								forms.share_url,
								forms.created_by,
								forms.created_at,
								forms.updated_at,
								COALESCE(forms.archived_at, null) as archived_at,
								(
									CASE
									WHEN u.id = (
									SELECT 
										org.created_by
									FROM 
										mstr.organizations org
									WHERE 
										org.id = form_organizations.organization_id
									) THEN 'Owner'
									ELSE 'Member'
								END
								) AS author,								
								forms.submission_target_user,
								form_organizations.organization_id,
								(CASE 
									WHEN (
									SELECT 
										COUNT(t.id) AS count_share 
									FROM 
										frm.form_organization_invites t 
									WHERE 
										t.form_id = forms.id
									) > 0 THEN 'out' 
									ELSE '' 
								END
								) AS form_shared,
								'' as form_external_company_name,
								'' as form_external_company_image,
								'internal' as type
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
						),
						subquery2 AS (
							SELECT 
								0 as form_share_count,
								0 as project_id, 
								f.id, 
								f.form_status_id, 
								fs.name as form_status, 
								f.name, 
								f.description, 
								f.notes, 
								f.period_start_date, 
								f.period_end_date, 
								f.profile_pic, 
								f.is_publish, 
								f.is_attendance_required, 
								u.name as created_by_name, 
								u.email as created_by_email, 
								f.share_url, 
								f.created_by, 
								f.created_at, 
								f.updated_at, 
								coalesce(f.archived_at, null) as archived_at, 
								'Member' as author, 
								f.submission_target_user, 
								form_to_user_invites.organization_receiver_id, 
								'in' as form_shared,
								o.name as form_external_company_name,
								o.profile_pic as form_external_company_image,
								'external' as type
							FROM 
								"frm"."form_to_user_invites" 
								left join frm.forms f on f.id = form_to_user_invites.form_id 
								left join frm.form_organizations fo on fo.form_id = form_to_user_invites.form_id 
								left join mstr.organizations o on o.id = fo.organization_id 
								left join mstr.form_statuses fs on fs.id = f.form_status_id 
								left join usr.users u on u.id = f.created_by 
							WHERE ` + whereStr + `
								AND "form_to_user_invites"."user_receiver_id" = ` + strconv.Itoa(userID) + `
								AND "form_to_user_invites"."organization_receiver_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								AND "form_to_user_invites"."deleted_at" IS NULL
						)	
						SELECT 
								form_share_count,
								project_id,
								id,
								form_status_id,
								form_status,
								name,
								description,
								notes,
								period_start_date,
								period_end_date,
								profile_pic,
								is_publish,
								is_attendance_required,
								created_by_name,
								created_by_email,
								share_url,
								created_by,
								created_at,
								updated_at,
								archived_at,
								author,
								submission_target_user,
								organization_id,
								form_shared,
								form_external_company_name,
								form_external_company_image,
								type
							FROM (
								SELECT * FROM subquery1
								UNION ALL
								SELECT * FROM subquery2
							) AS tb_2
							` + whereName + `
							ORDER BY 	
								(CASE WHEN form_shared = 'group' THEN 1 ELSE 2 END),
								` + orderBy + ` 
								name asc, 
								created_at desc
							OFFSET ` + strconv.Itoa(offset) + `
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetProjectSuperAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`SELECT 
								0 AS form_share_count, 
								t.id AS project_id, 
								0 AS id, 
								0 AS form_status_id, 
								'' AS form_status, 
								t.name, 
								t.description, 
								'' AS notes, 
								NOW() AS period_start_date, 
								NOW() AS period_end_date, 
								'' AS profile_pic, 
								FALSE AS is_publish, 
								FALSE AS is_attendance_required, 
								'' AS created_by_name, 
								'' AS created_by_email, 
								'' AS share_url, 
								t.created_by, 
								t.created_at, 
								t.updated_at, 
								NOW() AS archived_at, 
								'' AS author, 
								0 AS submission_target_user, 
								po.organization_id, 
								'group' AS form_shared, 
								'' AS form_external_company_name, 
								'' AS form_external_company_image, 
								'' AS type 
							FROM 
								frm.projects t 
							LEFT JOIN frm.project_organizations po ON po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + ` 
							ORDER BY 
							` + orderBy + ` 
								t.name asc
							OFFSET ` + strconv.Itoa(offset) + ` 
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetProjectAdminNew(fields tables.FormOrganizationsJoin, whereName string, whereString string, whereGroupString string, whereStr string, userID int, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	// whereGroup := ""
	// if whereString != "" {
	// 	whereGroup = strings.Replace(whereString, "forms.name", "t.name", 10)
	// }

	err := con.db.Raw(`SELECT
								0 AS form_share_count, 
								t.id AS project_id, 
								0 AS id,
								0 AS form_status_id,
								'' AS form_status,
								t.name,
								t.description,
								'' AS notes,
								NOW() AS period_start_date,
								NOW() AS period_end_date,
								'' AS profile_pic,
								false AS is_publish,
								false AS is_attendance_required,
								'' AS created_by_name,
								'' AS created_by_email,
								'' AS share_url,
								t.created_by,
								t.created_at,
								t.updated_at,
								NOW() AS archived_at,
								'' AS author,
								0 AS submission_target_user,
								po.organization_id,
								'group' AS form_shared,
								'' AS form_external_company_name,
								'' AS form_external_company_image,
								'' AS type
							FROM 
								frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							` + whereGroupString + `
							ORDER BY 	
							` + orderBy + ` 
								t.name asc
							OFFSET ` + strconv.Itoa(offset) + `
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetDetailFormUserRow(fields tables.FormUsers, whrString string) (tables.FormOut, error) {
	var data tables.FormOut
	err := con.db.Scopes(SchemaFrm("forms")).Select("forms.*, fs.name as form_status").Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("join frm.form_users fu on fu.form_id=forms.id").Where(fields).Where(whrString).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return data, nil
	}

	if err != nil {
		return tables.FormOut{}, err
	}

	return data, nil
}

func (con *formConnection) GetUserFormOrganization(fields tables.UserFormOrganizations, whrString string) (tables.FormOrganizations, error) {

	var data tables.FormOrganizations
	err := con.db.Scopes(SchemaFrm("form_organizations")).Joins("join mstr.organizations o on o.id=form_organizations.organization_id").Where(fields).Where(whrString).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.FormOrganizations{}, nil
	}

	if err != nil {
		return data, err
	}

	return data, nil
}

func (con *formConnection) GetUserCompaniesListInvitedRows(fields tables.UserOrganizationInvites, whrStr string) ([]tables.UserOrganizationInviteDetail, error) {
	var data []tables.UserOrganizationInviteDetail
	err := con.db.Scopes(SchemaUsr("user_organization_invites")).Select("user_organization_invites.organization_id, user_organization_invites.user_id, o.name as organization_name, o.profile_pic as organization_profile_pic, o.contact_name as organization_contact_name,  o.contact_phone as organization_contact_phone").Joins("left join mstr.organizations o ON user_organization_invites.organization_id=o.id").Where(fields).Where(whrStr).Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) InsertFormCompanyInvites(data tables.FormOrganizationInvites) (bool, error) {

	err := con.db.Scopes(SchemaFrm("form_organization_invites")).Create(&data).Error
	if err != nil {
		fmt.Println("error FormOrganizationInvites--- ", err)
		return false, err
	}

	return true, nil
}

func (con *formConnection) DeleteFormCompanyInvites(formID int, organizationID int) (bool, error) {

	err := con.db.Exec("DELETE FROM frm.form_organization_invites where form_id = " + strconv.Itoa(formID) + " and organization_id = " + strconv.Itoa(organizationID)).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formConnection) GetFormCompanyInviteRows(fields tables.JoinFormCompanies, whereStr string) ([]tables.JoinFormCompanies, error) {

	var data []tables.JoinFormCompanies
	err := con.db.Table("frm.form_organization_invites").Select("form_organization_invites.id, form_organization_invites.form_id, form_organization_invites.is_quota_sharing, o.id as organization_id, o.name as organization_name, o.code as organization_code, o.contact_name as organization_contact_name, o.contact_phone as organization_contact_phone, o.profile_pic as organization_profile_pic").Joins("join mstr.organizations o on o.id=form_organization_invites.organization_id").Where(fields).Where(whereStr).Order("form_organization_invites.created_at desc").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormCompanyInviteRow(fields tables.JoinFormCompanies, whereStr string) (tables.JoinFormCompanies, error) {

	var data tables.JoinFormCompanies
	err := con.db.Table("frm.form_organization_invites").Select("form_organization_invites.id, form_organization_invites.form_id , form_organization_invites.is_quota_sharing, o.id as organization_id, o.name as organization_name, o.code as organization_code, o.contact_name as organization_contact_name, o.contact_phone as organization_contact_phone, o.profile_pic as organization_profile_pic").Joins("join mstr.organizations o on o.id=form_organization_invites.organization_id").Where(fields).Where(whereStr).Order("form_organization_invites.created_at desc").First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.JoinFormCompanies{}, nil
	}

	if err != nil {
		return tables.JoinFormCompanies{}, err
	}
	return data, nil
}

func (con *formConnection) UpdateCompanyInviteForm(id int, fields tables.FormOrganizationInvites) (bool, error) {

	if fields.IsQuotaSharing == false {
		err := con.db.Scopes(SchemaFrm("form_organization_invites")).Where("id = ?", id).Update("is_quota_sharing", "false").Error
		if err != nil {
			return false, err
		}
	} else {
		err := con.db.Scopes(SchemaFrm("form_organization_invites")).Where("id = ?", id).Updates(fields).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (con *formConnection) GetAllIDForDelete(ID int) (objects.DeleteAdminEksObj, error) {

	var data objects.DeleteAdminEksObj

	err := con.db.Table("usr.user_organization_invites").
		Select("usr.user_organization_invites.id,usr.user_organization_invites.user_id,usr.user_organization_invites.organization_id,uo.id as user_organization_id,uor.id as user_organization_roles_id,foi.organization_id as form_organization_invites_id").
		Where("usr.user_organization_invites.id = ?", ID).
		Joins("left join usr.user_organizations uo on uo.user_organization_invite_id = usr.user_organization_invites.id").
		Joins("left join usr.user_organization_roles uor on uor.user_organization_id  = uo.id").
		Joins("left join frm.form_organization_invites foi on foi.organization_id = usr.user_organization_invites.organization_id").
		First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *formConnection) DeleteAdminEks(data objects.DeleteAdminEksObj) (objects.DeleteAdminEksObj, error) {
	var dataID objects.DeleteAdminEksObj
	dataID.ID = data.ID
	dataID.UserID = data.UserID
	dataID.OrganizationID = data.OrganizationID
	dataID.UserOrganizationID = data.UserOrganizationID
	dataID.UserOrganizationRolesID = data.UserOrganizationRolesID
	dataID.FormOrganizationInvitesID = data.FormOrganizationInvitesID

	err := con.db.Scopes(SchemaUsr("user_organization_roles")).Where("id = ?", dataID.UserOrganizationRolesID).Delete(&dataID.UserOrganizationRolesID).Error
	if err != nil {
		fmt.Println("error user_organization_roles --- ", err)
		return objects.DeleteAdminEksObj{}, err
	}

	err = con.db.Scopes(SchemaUsr("user_organizations")).Where("id = ?", dataID.UserOrganizationID).Delete(&dataID.UserOrganizationID).Error
	if err != nil {
		fmt.Println("error user_organizations --- ", err)
		return objects.DeleteAdminEksObj{}, err
	}

	err = con.db.Scopes(SchemaUsr("user_organization_invites")).Where("id = ?", dataID.ID).Delete(&dataID.ID).Error
	if err != nil {
		fmt.Println("error user_organization_invites --- ", err)
		return objects.DeleteAdminEksObj{}, err
	}

	err = con.db.Scopes(SchemaFrm("form_organization_invites")).Where("organization_id = ?", dataID.FormOrganizationInvitesID).Delete(&dataID.FormOrganizationInvitesID).Error
	if err != nil {
		fmt.Println("error form_organization_invites --- ", err)
		return objects.DeleteAdminEksObj{}, err
	}

	return dataID, err
}

func (con *formConnection) GetFormOtherCompanyRows(fields tables.FormOrganizationsJoin, whereString string, paging objects.Paging) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.FormAll

	err := con.db.Table("frm.form_organization_invites").Select("forms.id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, forms.is_attendance_required, u.name as created_by_name, u.email as created_by_email,forms.share_url, forms.created_by, forms.created_at, forms.updated_at, coalesce(forms.archived_at, null) as archived_at, (case when u.id=(select org.created_by from mstr.organizations org where org.id = form_organization_invites.organization_id) then 'Owner' else 'Member' end) as author, forms.submission_target_user, form_organization_invites.organization_id, form_organization_invites.is_quota_sharing").Joins("left join frm.forms on forms.id = form_organization_invites.form_id").Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u on u.id = forms.created_by").Where(fields).Where(whereString).Order(orderBy).Order("forms.name asc, forms.created_at desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) GenerateFormUserOrg(data tables.FormUserOrganizations) (tables.FormUserOrganizations, error) {
	err := con.db.Scopes(SchemaFrm("form_user_organizations")).Create(&data).Error
	if err != nil {
		return tables.FormUserOrganizations{}, err
	}

	return data, nil
}

func (con *formConnection) GetFormUserToFormOrgRows(fields tables.JoinFormUsers, whereStr string) ([]tables.JoinFormUsers, error) {

	var data []tables.JoinFormUsers
	err := con.db.Table("frm.form_users").Select("form_users.id, form_users.form_id, fus.id as form_user_status_id, t.translation as form_user_status_name, u.id as user_id, u.name as user_name, u.email, u.phone, f.name, f.description, f.notes, f.profile_pic, f.period_start_date, f.period_end_date, s.id as form_status_id, s.name as form_status, f.created_at, fo.organization_id").Joins("join frm.forms f on f.id=form_users.form_id").Joins("left join usr.users u on u.id = form_users.user_id").Joins("left join mstr.form_statuses s on s.id = f.form_status_id").Joins("left join mstr.form_user_statuses fus on fus.id = form_users.form_user_status_id").Where(fields).Where(whereStr).Order("form_users.id asc").Joins("left join mstr.translations t on t.textcontent_id = fus.name_textcontent_id AND t.language_id = (select ul.language_id from usr.user_languages ul where ul.user_id = 5)").Joins("join frm.form_organizations fo on fo.form_id=form_users.form_id").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetListFormEksternal(fields tables.Forms, whereString string, paging objects.Paging, user_id int, organization_id int) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.FormAll

	err := con.db.
		Table("frm.form_users").
		Select("f.id, form_users.form_id, form_users.user_id, foi.organization_id, f.form_status_id, fs.name as form_status, f.name, f.description, f.notes, f.period_start_date, f.period_end_date, f.profile_pic, f.is_publish, f.is_attendance_required, u.name as created_by_name, u.email as created_by_email, f.share_url, f.created_by, f.created_at, f.updated_at, coalesce(f.archived_at, null) as archived_at, f.submission_target_user").
		Joins("left join frm.forms f on f.id = form_users.form_id").
		Joins("left join frm.form_organization_invites foi on form_users.form_id = foi.form_id ").
		Joins("left join mstr.form_statuses fs on fs.id = f.form_status_id").
		Joins("left join usr.users u on u.id = f.created_by").
		Where(fields).
		Where(whereString).
		Where("form_users.user_id = ?", user_id).
		Where("foi.organization_id = ?", organization_id).
		Order(orderBy).
		Order("f.name asc, f.created_at desc").
		Limit(paging.Limit).
		Offset(offset).
		Find(&data).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) GetListFormEksternalSuperAdmin(fields tables.Forms, whereString string, paging objects.Paging, organization_id int) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.FormAll

	err := con.db.
		Table("frm.form_users").
		Select("f.id, form_users.form_id, form_users.user_id, foi.organization_id, f.form_status_id, fs.name as form_status, f.name, f.description, f.notes, f.period_start_date, f.period_end_date, f.profile_pic, f.is_publish, f.is_attendance_required, u.name as created_by_name, u.email as created_by_email, f.share_url, f.created_by, f.created_at, f.updated_at, coalesce(f.archived_at, null) as archived_at, f.submission_target_user").
		Joins("left join frm.forms f on f.id = form_users.form_id").
		Joins("left join frm.form_organization_invites foi on form_users.form_id = foi.form_id ").
		Joins("left join mstr.form_statuses fs on fs.id = f.form_status_id").
		Joins("left join usr.users u on u.id = f.created_by").
		Where(fields).
		Where(whereString).
		Where("foi.organization_id = ?", organization_id).
		Order(orderBy).
		Order("f.name asc, f.created_at desc").
		Limit(paging.Limit).
		Offset(offset).
		Find(&data).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}
func (con *formConnection) GetFormEksternalOwnerRows(fields tables.FormOrganizationsJoin, whereString string, paging objects.Paging) ([]tables.FormAll, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.FormAll

	err := con.db.Table("frm.form_organization_invites").Select("forms.id, forms.form_status_id, fs.name as form_status, forms.name, forms.description, forms.notes, forms.period_start_date, forms.period_end_date, forms.profile_pic, forms.is_publish, forms.is_attendance_required, u.name as created_by_name, u.email as created_by_email,forms.share_url, forms.created_by, forms.created_at, forms.updated_at, coalesce(forms.archived_at, null) as archived_at, (case when u.id=(select org.created_by from mstr.organizations org where org.id = form_organization_invites.organization_id) then 'Owner' else 'Member' end) as author, forms.submission_target_user, form_organization_invites.organization_id").
		Joins("left join frm.forms on forms.id = form_organization_invites.form_id").
		Joins("left join mstr.form_statuses fs on fs.id = forms.form_status_id").Joins("left join usr.users u on u.id = forms.created_by").
		Where(fields).Where(whereString).Order(orderBy).Order("forms.name asc, forms.created_at desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}
func (con *formConnection) GetBlastInfoData(formID int, whereString string) ([]objects.BlastInfoData, error) {
	// fID := strconv.Itoa(formID)

	var data []objects.BlastInfoData
	err := con.db.Table("usr.notification_histories").
		Select("notification_histories.user_id, notification_histories.created_at, if.destination_form_id as form_id").
		Joins("left join mstr.broadcast_messages if on if.id = notification_histories.broadcast_message_id").
		Where("if.destination_form_id = ?", formID).
		Where(whereString).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetBlastInfoDataUsers(formID int, userID int, whereString string) ([]objects.BlastInfoData, error) {
	// fID := strconv.Itoa(formID)

	var data []objects.BlastInfoData
	err := con.db.Table("usr.notification_histories").
		Select("notification_histories.user_id, notification_histories.created_at, if.destination_form_id as form_id").
		Joins("left join mstr.broadcast_messages if on if.id = notification_histories.broadcast_message_id").
		Where("if.destination_form_id = ?", formID).
		Where("notification_histories.user_id = ?", userID).
		Where(whereString).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFillingType(fields objects.FillingType) ([]objects.FillingType, error) {
	// fID := strconv.Itoa(formID)

	var data []objects.FillingType
	err := con.db.Table("mstr.filling_type").
		Select("filling_type.id, filling_type.name, filling_type.status").
		Where(fields).
		Where("filling_type.status", true).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormOrganizationInvite(formID int) (tables.JoinFormCompanies, error) {
	// fID := strconv.Itoa(formID)

	var data tables.JoinFormCompanies
	err := con.db.Table("frm.form_organization_invites").
		Select("form_organization_invites.organization_id").
		Where("form_organization_invites.form_id = ?", formID).Find(&data).Error

	if err != nil {
		return tables.JoinFormCompanies{}, err
	}
	return data, nil
}

func (con *formConnection) GetFormCompanyInviteNew(fields tables.JoinFormCompanies, whereStr string) ([]tables.JoinFormCompanies, error) {

	var data []tables.JoinFormCompanies
	err := con.db.Table("frm.form_organization_invites").Select("form_organization_invites.id, form_organization_invites.form_id, form_organization_invites.is_quota_sharing, o.id as organization_id, o.name as organization_name, o.code as organization_code, form_organization_invites.user_organization_invite_id, u.name as organization_contact_name, u.phone as organization_contact_phone, u.avatar as organization_profile_pic, uoi.invite_type as type").Joins("join mstr.organizations o on o.id = form_organization_invites.organization_id").Joins("left join usr.user_organizations uo on uo.user_organization_invite_id = form_organization_invites.user_organization_invite_id").Joins("left join usr.users u on u.id = uo.user_id").Joins("left join usr.user_organization_invites uoi on uoi.id = form_organization_invites.user_organization_invite_id").Where(fields).Where(whereStr).Order("form_organization_invites.created_at desc").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormUnionProjectAndExternalRows(userID int, fields tables.FormOrganizationsJoin, whereString string, whereGroupString string, paging objects.Paging) ([]tables.FormAll, error) {
	var data []tables.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	err := con.db.Raw(`SELECT *
							from
							(SELECT '' as form_shared,    
									0 as form_share_count,
									'' as form_external_company_name,
									'' as form_external_company_image,
									t.id as project_id,
									0 as id,
									0 as form_status_id,
									'' as form_status,
									t.name,
									t.description,
									'' as notes,
									now() as period_start_date,
									now() as period_end_date,
									'' as profile_pic, --10
									false as is_publish,
										false as is_attendance_required,
										'' as created_by_name,
										'' as created_by_email,
										'' as share_url,
										t.created_by,
										t.created_at,
										t.updated_at,
										now() as archived_at, --10
										'' as author,
										0 as submission_target_user,
										po.organization_id,
										false as is_quota_sharing
							FROM frm.projects t
							LEFT JOIN frm.project_organizations po on po.project_id = t.id
							WHERE t.deleted_at is null AND po.organization_id = ` + strconv.Itoa(fields.OrganizationID) + `
							ORDER BY t.id) as tb_2
						UNION
							(SELECT '' as form_shared,
									(SELECT count(t.id) as count_share
										FROM frm.form_organization_invites t
										where form_id =forms.id) as form_share_count,
									'' as form_external_company_name,
									'' as form_external_company_image,
									0 as project_id,
									forms.id,
									forms.form_status_id,
									fs.name as form_status,
									forms.name,
									forms.description,
									forms.notes,
									forms.period_start_date,
									forms.period_end_date,
									forms.profile_pic, --10
									forms.is_publish,
									forms.is_attendance_required,
									u.name as created_by_name,
									u.email as created_by_email,
									forms.share_url,
									forms.created_by,
									forms.created_at,
									forms.updated_at, --8
									coalesce(forms.archived_at, null) as archived_at,
									(case
										when u.id=
												(select org.created_by
												from mstr.organizations org
												where org.id = form_organizations.organization_id) then 'Owner'
										else 'Member'
									end) as author, --10
									forms.submission_target_user,
									form_organizations.organization_id,
									false as is_quota_sharing
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								AND forms.form_status_id not in (3)
								AND forms.id not in
								(select pf.form_id
									from frm.project_forms pf
									where pf.project_id in
										(select p.id
										from frm.projects p
										where p.created_by = ` + strconv.Itoa(userID) + `))
								AND "form_organizations"."deleted_at" IS NULL
								AND forms.id not in (select pf.form_id from frm.project_forms pf) )
						UNION
							(SELECT 'in' as form_shared,
									0 as form_share_count,
									(select o.name from frm.form_organizations fo join mstr.organizations o on o.id=fo.organization_id where fo.form_id=forms.id) as form_external_company_name,
									(select o.profile_pic from frm.form_organizations fo join mstr.organizations o on o.id=fo.organization_id where fo.form_id=forms.id) as  form_external_company_image,
									0 as project_id,
									forms.id,
									forms.form_status_id,
									fs.name as form_status,
									forms.name,
									forms.description,
									forms.notes,
									forms.period_start_date,
									forms.period_end_date,
									forms.profile_pic, --10
									forms.is_publish,
									forms.is_attendance_required,
									u.name as created_by_name,
									u.email as created_by_email,
									forms.share_url,
									forms.created_by,
									forms.created_at,
									forms.updated_at, --8
									coalesce(forms.archived_at, null) as archived_at,
									(case
										when u.id=
												(select org.created_by
												from mstr.organizations org
												where org.id = form_organization_invites.organization_id) then 'Owner'
										else 'Member'
									end) as author, --10
									forms.submission_target_user,
									form_organization_invites.organization_id,
									form_organization_invites.is_quota_sharing
							FROM "frm"."form_organization_invites"
							left join frm.forms on forms.id = form_organization_invites.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							WHERE "form_organization_invites"."organization_id" = ` + strconv.Itoa(fields.OrganizationID) + `
								AND (forms.form_status_id not in (3)
									AND forms.id not in
										(select pf.form_id
										from frm.project_forms pf
										where pf.project_id in
											(select p.id
											from frm.projects p
											where p.created_by = ` + strconv.Itoa(userID) + `)))
								AND "form_organization_invites"."deleted_at" IS NULL
							ORDER BY forms.name asc, forms.created_at desc)
							
						order by  ` + orderBy + ` project_id desc, name asc
						OFFSET ` + strconv.Itoa(offset) + ` limit ` + strconv.Itoa(paging.Limit) + `	
						`).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetDate() ([]objects.Date, error) {
	var data []objects.Date
	err := con.db.Raw(`select to_char(g_date::date,'yyyy-mm-dd') as date FROM generate_series(date '2023-06-01', date '2023-07-31', '1 day') as g_date`).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) GetAllForm() ([]objects.FormData, error) {
	var data []objects.FormData
	err := con.db.Raw(`select fo.form_id from frm.form_organizations fo where fo.organization_id in (36,35,30,29)`).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formConnection) GetFormByOrganization(fields objects.HistoryBalanceSaldo) ([]objects.HistoryBalanceSaldo, error) {

	var data []objects.HistoryBalanceSaldo
	err := con.db.Scopes(SchemaFrm("form_organizations")).Where(fields).Find(&data).Error
	if err != nil {
		return []objects.HistoryBalanceSaldo{}, err
	}

	return data, err
}
func (con *formConnection) GetHistoryTopupByDate(organizationID int, whreDate string) ([]objects.TopupHistory, error) {
	var data []objects.TopupHistory

	err := con.db.
		Table("usr.organization_topup_histories").
		Select("sum(organization_topup_histories.respondent_quota) as quota").
		Joins("left join usr.organization_subscription_plans osp ON osp.id = organization_topup_histories.organization_subscription_plan_id").
		Where(whreDate).
		Where("osp.organization_id = ?", organizationID).
		Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormTemplate(UserID int) ([]objects.FormTemplateNew, error) {

	var data []objects.FormTemplateNew

	err := con.db.
		Table("frm.forms").
		Select("forms.id as form_id, pf.project_id, forms.name, forms.profile_pic, forms.description").
		Joins("left join frm.project_forms pf on pf.form_id = forms.id").
		Where("forms.created_by = ?", UserID).
		Where("forms.form_status_id = ?", 1).
		Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetFormTemplateByProjectID(UserID int, ProjectID int) ([]objects.FormTemplateNew, error) {

	var data []objects.FormTemplateNew

	err := con.db.
		Table("frm.project_forms").
		Select("project_forms.project_id, f.id as form_id, f.name, f.profile_pic, f.description").
		Joins("left join frm.forms f on f.id = project_forms.form_id").
		Where("project_forms.project_id = ?", ProjectID).
		Where("f.created_by = ?", UserID).
		Where("f.form_status_id = ?", 1).
		Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *formConnection) GetProject(UserID int) ([]objects.Projects, error) {

	var data []objects.Projects

	err := con.db.
		Table("frm.projects").
		Select("projects.id, projects.name, projects.description").
		Where("projects.created_by = ?", UserID).
		Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}
