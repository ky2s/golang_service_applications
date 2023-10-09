package models

import (
	"fmt"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

type ProjectModels interface {
	InsertProject(data tables.Projects, organizationID int) (tables.Projects, error)
	GetProjectRows(data tables.Projects) ([]tables.Projects, error)
	GetProjectRow(data tables.Projects) (tables.Projects, error)
	UpdateProject(id int, data tables.Projects) (bool, error)
	DeleteProject(id int) (bool, error)
	InsertProjectForm(data tables.ProjectForms) (tables.ProjectForms, error)
	DeleteProjectForm(data tables.ProjectForms) (bool, error)
	GetProjectForms(fields tables.ProjectForms, whrStr string) ([]tables.ProjectFormsJoin, error)
	GetProjectInForms(userID int) ([]tables.Projects, error)
	CheckFormIn(formID int, organizationID int) ([]objects.FormOrganizations, error)
	CheckFormOut(formID int, organizationID int) ([]objects.FormOrganizations, error)
	GetFormInProject(whereString string, organizationID int, paging objects.Paging) ([]objects.FormAll, error)
}

type projectConnection struct {
	db *gorm.DB
}

func NewProjectModels(dbg *gorm.DB) ProjectModels {
	return &projectConnection{
		db: dbg,
	}
}

func (con *projectConnection) InsertProject(data tables.Projects, organizationID int) (tables.Projects, error) {
	err := con.db.Scopes(SchemaFrm("projects")).Create(&data).Error
	if err != nil {
		return tables.Projects{}, err
	}

	fmt.Println("data.ID ------", data.ID)
	var dataProject tables.ProjectOrganizations
	dataProject.ProjectID = data.ID
	dataProject.OrganizationID = organizationID
	err2 := con.db.Scopes(SchemaFrm("project_organizations")).Create(&dataProject).Error
	if err2 != nil {
		return tables.Projects{}, err2
	}

	return data, err
}

func (con *projectConnection) GetProjectRows(fields tables.Projects) ([]tables.Projects, error) {
	var data []tables.Projects
	err := con.db.Scopes(SchemaFrm("projects")).Where(fields).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *projectConnection) GetProjectRow(fields tables.Projects) (tables.Projects, error) {
	var data tables.Projects
	err := con.db.Scopes(SchemaFrm("projects")).Where(fields).First(&data).Error
	if err != nil {
		return tables.Projects{}, err
	}
	return data, nil
}

func (con *projectConnection) UpdateProject(id int, fields tables.Projects) (bool, error) {
	err := con.db.Scopes(SchemaFrm("projects")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *projectConnection) DeleteProject(projectID int) (bool, error) {
	var data tables.Projects
	err := con.db.Scopes(SchemaFrm("projects")).Delete(&data, projectID).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *projectConnection) InsertProjectForm(data tables.ProjectForms) (tables.ProjectForms, error) {
	err := con.db.Scopes(SchemaFrm("project_forms")).Create(&data).Error

	return data, err
}

func (con *projectConnection) DeleteProjectForm(fields tables.ProjectForms) (bool, error) {
	err := con.db.Exec("DELETE FROM frm.project_forms where project_id = " + strconv.Itoa(fields.ProjectID) + " and form_id = " + strconv.Itoa(fields.FormID)).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *projectConnection) GetProjectForms(fields tables.ProjectForms, whreString string) ([]tables.ProjectFormsJoin, error) {
	var data []tables.ProjectFormsJoin
	// err := con.db.Scopes(SchemaFrm("project_forms")).Where(fields).Find(&data).Error
	err := con.db.Table("frm.project_forms").Select("project_forms.project_id, project_forms.form_id, f.name, f.description, f.profile_pic, f.created_by, f.form_status_id, fs.name as form_status, u.name as created_by_name, u.email as created_by_email").Joins("left join frm.forms f on f.id = project_forms.form_id").Joins("left join usr.users u on u.id = f.created_by").Joins("left join mstr.form_statuses fs ON fs.id = f.form_status_id").Where(fields).Where(whreString).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *projectConnection) GetProjectInForms(userID int) ([]tables.Projects, error) {
	var data []tables.Projects

	err := con.db.Table("frm.projects").Where(" projects.id in (select pf.project_id from frm.project_forms pf where pf.form_id in (select fu.form_id from frm.form_users fu where fu.user_id = " + strconv.Itoa(userID) + ")) ").First(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *projectConnection) CheckFormIn(formID int, organizationID int) ([]objects.FormOrganizations, error) {
	var data []objects.FormOrganizations

	err := con.db.Table("frm.form_organizations").Select("form_organizations.id, form_organizations.form_id, form_organizations.organization_id").Where("form_organizations.form_id = ?", formID).Where("form_organizations.organization_id = ?", organizationID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *projectConnection) CheckFormOut(formID int, organizationID int) ([]objects.FormOrganizations, error) {
	var data []objects.FormOrganizations

	err := con.db.Table("frm.form_organization_invites").Select("form_organization_invites.id, form_organization_invites.form_id, form_organization_invites.organization_id").Where("form_organization_invites.form_id = ?", formID).Where("form_organization_invites.organization_id = ?", organizationID).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *projectConnection) GetFormInProject(whereString string, organizationID int, paging objects.Paging) ([]objects.FormAll, error) {
	var data []objects.FormAll

	var orderBy string
	if paging.SortBy != "" && paging.Sort != "" {
		orderBy = paging.SortBy + " " + paging.Sort + ", "
	}

	offset := (paging.Page - 1) * paging.Limit
	if paging.Limit == 0 {
		paging.Limit = 100 // default rows per page
	}

	err := con.db.Raw(`
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
								'internal' AS type,								
								pf.id as project_form_id,
								p2.id as project_id 
							FROM "frm"."form_organizations"
							left join frm.forms on forms.id = form_organizations.form_id
							left join mstr.form_statuses fs on fs.id = forms.form_status_id
							left join usr.users u on u.id = forms.created_by
							left join frm.project_forms pf on pf.form_id = forms.id
							left join frm.projects p2 on p2.id = pf.project_id 
							WHERE "form_organizations"."organization_id" = ` + strconv.Itoa(organizationID) + `
							` + whereString + `
							AND "form_organizations"."deleted_at" IS NULL
							ORDER BY
							` + orderBy + ` 
								forms.name asc
							OFFSET ` + strconv.Itoa(offset) + ` 
							limit ` + strconv.Itoa(paging.Limit) + ``).Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}
