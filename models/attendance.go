package models

import (
	"errors"
	"fmt"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

type AttendanceModels interface {
	InsertAttendance(data tables.Attendances, geo tables.Geometry) (tables.Attendances, error)
	UpdateAttendance(id int, data tables.Attendances, geo tables.Geometry) (bool, error)
	GetAttendanceRows(data tables.Attendances, strField string) ([]tables.Attendances, error)
	GetAttendanceRow(data tables.Attendances, strField string) (tables.Attendances, error)
	GetAttendanceObjRow(data tables.Attendances, strField string) (objects.Attendances, error)
	GetAttendanceReports(fields tables.Attendances, strField string) ([]tables.AttendenceReport, error)
	GetUserAttendanceReports(fields tables.Attendances, strField string, whreStr string, paging objects.Paging) ([]tables.AttendenceReport, error)
	GetUserTeamAttendanceReports(fields tables.Attendances, strField string, whreStr string, whrComp string, paging objects.Paging) ([]tables.AttendenceReport, error)
	GetUserAttendanceMapReports(fields tables.Attendances, strField string, whreTime string, paging objects.Paging) ([]tables.AttendenceReport, error)
	GetUserAttendanceMap2Reports(fields tables.Attendances, strField string, whreTime string, whrComp string) ([]tables.AttendenceMapReport, error)
	GetUserAttendanceMap3Reports(fields tables.Attendances, strField string, whreTime string, whrComp string) ([]tables.AttendenceMapReport, error) // new latest
	GetLastAttendanceOverdate(formID int, userID int) (objects.LastAttendance, error)
	GetMissingIDAtt(fields objects.MissingIDAtt) ([]objects.MissingIDAtt, error)
	GetCompanyID(form_id int, user_id int) (objects.FormUserOrg, error)
	InsertAttendanceOrganization(fields objects.InsAttOrg) (objects.InsAttOrg, error)
	GetCompanyIDByFormID(form_id int) (objects.GetOrgID, error)
	GetFormAttendanceLocationRows(fields objects.ObjectFormAttendanceLocations) ([]objects.ObjectFormAttendanceLocations, error)
	GetTeamByRespondent(respondenID int) ([]objects.TeamUsers, error)
	GetFormTeam(formID int, TeamID int) (objects.FormTeams, error)
}

type attendConnection struct {
	db *gorm.DB
}

func NewAttendanceModels(dbg *gorm.DB) AttendanceModels {
	return &attendConnection{
		db: dbg,
	}
}

func (con *attendConnection) InsertAttendance(data tables.Attendances, geo tables.Geometry) (tables.Attendances, error) {
	fmt.Println("lat, long models ====", geo)
	err := con.db.Scopes(SchemaUsr("attendances")).Create(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("ERROR :::: lat, long models ====", geo)
	}

	//geometry
	long := strconv.FormatFloat(geo.Longitude, 'f', 30, 64)
	lat := strconv.FormatFloat(geo.Latitude, 'f', 30, 64)

	fmt.Println("GEO LAT LONG ====", geo.Latitude, geo.Longitude, lat, long)
	err2 := con.db.Scopes(SchemaUsr("attendances")).Where("id = ?", data.ID).Update("geometry_in", gorm.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", long, lat)).Error
	if err2 != nil {
		return data, err2
	}
	return data, err
}

func (con *attendConnection) GetAttendanceRows(fields tables.Attendances, strField string) ([]tables.Attendances, error) {
	var data []tables.Attendances
	err := con.db.Scopes(SchemaUsr("attendances")).Where(fields).Where(strField).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}

	return data, nil
}

func (con *attendConnection) GetAttendanceReports(fields tables.Attendances, strField string) ([]tables.AttendenceReport, error) {
	var data []tables.AttendenceReport
	err := con.db.Table("usr.attendances").Select(`attendances.id, attendances.form_id, attendances.user_id, to_char(attendances.attendance_in::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_in, to_char(attendances.attendance_in::timestamp, 'yyyy-mm-dd') as attendance_date_in, to_char(attendances.attendance_in::timestamp, 'HH24:MI') as attendance_time_in, to_char(attendances.attendance_out::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_out, to_char(attendances.attendance_out::timestamp, 'yyyy-mm-dd') as attendance_date_out, to_char(attendances.attendance_out::timestamp, 'HH24:MI') as attendance_time_out, attendances.face_pic_in, attendances.face_pic_out, u.name as user_name, u.phone as user_phone, to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at, to_char(attendances.updated_at::timestamp, 'yyyy-mm-dd HH24:MI') as updated_at, attendances.address_in, attendances.address_out, mstr.get_duration_time(attendances.attendance_in, attendances.attendance_out) as duration, 
	(case when (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_in and  fal.form_id = attendances.form_id limit 1) is null 
	then '' else (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_in and  fal.form_id = attendances.form_id limit 1) 
	end) location_in,
	(case when (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_out and  fal.form_id = attendances.form_id limit 1) is null 
	then '' else (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_out and  fal.form_id = attendances.form_id  limit 1) 
	end) location_out`).Joins("left join usr.users u on u.id = attendances.user_id").Where(fields).Where(strField).Order("attendances.attendance_in").Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *attendConnection) GetUserAttendanceReports(fields tables.Attendances, strField string, whreTime string, paging objects.Paging) ([]tables.AttendenceReport, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.AttendenceReport
	err := con.db.Table("frm.form_users").Select("attendances.id, attendances.form_id, attendances.user_id, to_char(attendances.attendance_in::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_in, to_char(attendances.attendance_out::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_out, attendances.face_pic_in, attendances.face_pic_out, u.name as user_name, u.phone as user_phone, to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at, attendances.address_in, attendances.address_out,ST_Y(attendances.geometry_in::geometry) as latitude, ST_X(attendances.geometry_in::geometry) as longitude, mstr.get_duration_time(attendances.attendance_in, attendances.attendance_out) as duration").Joins("left join usr.attendances on attendances.user_id = form_users.user_id AND attendances.form_id = ? "+whreTime, fields.FormID).Joins("left join usr.users u on u.id = form_users.user_id").Where(fields).Where(strField).Where("form_users.type = 'respondent'").Limit(paging.Limit).Offset(offset).Order(orderBy).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *attendConnection) GetUserTeamAttendanceReports(fields tables.Attendances, strField string, whreTime string, whrComp string, paging objects.Paging) ([]tables.AttendenceReport, error) {

	offset := 0
	if paging.Page >= 1 {
		offset = (paging.Page - 1) * paging.Limit
	}

	orderBy := "user_name ASC"
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	limitOffset := ""
	if paging.Limit >= 1 {
		limitOffset = `LIMIT ` + strconv.Itoa(paging.Limit) + ` OFFSET ` + strconv.Itoa(offset)
	}

	formID := strconv.Itoa(fields.FormID)
	fmt.Println("fields ---------------------->>>::", paging.Limit, offset, paging.SortBy)
	var data []tables.AttendenceReport

	err := con.db.Raw(`select * FROM
						(
							SELECT attendances.id,
								attendances.form_id,
								attendances.user_id,
								ao.organization_id,
								o.name as organization_name,
								to_char(attendances.attendance_in::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_in,
								to_char(attendances.attendance_out::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_out,
								attendances.face_pic_in,
								attendances.face_pic_out,
								u.id as user_id,
								u.name as user_name,
								u.phone as user_phone,
								to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at,
								attendances.address_in,
								attendances.address_out,
								ST_Y(attendances.geometry_in::geometry) as latitude,
								ST_X(attendances.geometry_in::geometry) as longitude,
								mstr.get_duration_time(attendances.attendance_in, attendances.attendance_out) as duration,
								(case when (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_in and  fal.form_id = attendances.form_id limit 1) is null 
								then '' else (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_in and  fal.form_id = attendances.form_id limit 1) 
								end) location_in,
								(case when (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_out and  fal.form_id = attendances.form_id limit 1) is null 
								then '' else (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_out and  fal.form_id = attendances.form_id  limit 1) 
								end) location_out
							FROM "frm"."form_users"
							left join usr.attendances on attendances.user_id = form_users.user_id AND attendances.form_id = ? 
							
							left join usr.attendance_organizations ao on attendances.id = ao.attendance_id
							left join mstr.organizations o on ao.organization_id = o.id
							
							left join usr.users u on u.id = form_users.user_id
							WHERE "form_users"."form_id" = ? AND form_users.type = 'respondent' `+whrComp+`
							`+whreTime+` `+strField+`
						) as tb_1
						
						UNION
						
						(SELECT 
							attendances.id,
							attendances.form_id,
							attendances.user_id,
							ao.organization_id,
							o.name as organization_name,
							to_char(attendances.attendance_in::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_in,
							to_char(attendances.attendance_out::timestamp, 'yyyy-mm-dd HH24:MI') as attendance_out,
							attendances.face_pic_in,
							attendances.face_pic_out,
							u.id as user_id,
							u.name as user_name,
							u.phone as user_phone,
							to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at,
							attendances.address_in,
							attendances.address_out,
							ST_Y(attendances.geometry_in::geometry) as latitude,
							ST_X(attendances.geometry_in::geometry) as longitude,
							mstr.get_duration_time(attendances.attendance_in, attendances.attendance_out) as duration,
							(case when (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_in and  fal.form_id = attendances.form_id limit 1) is null 
								then '' else (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_in and  fal.form_id = attendances.form_id limit 1) 
								end) location_in,
								(case when (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_out and  fal.form_id = attendances.form_id limit 1) is null 
								then '' else (select fal."location" from frm.form_attendance_locations fal where fal.id = attendances.form_attendance_location_id_out and  fal.form_id = attendances.form_id  limit 1) 
								end) location_out
						FROM "usr"."team_users"
						left join usr.attendances on attendances.user_id = team_users.user_id AND attendances.form_id = ? 
						
						left join usr.attendance_organizations ao on attendances.id = ao.attendance_id
						left join mstr.organizations o on ao.organization_id = o.id
						
						left join usr.users u on u.id = team_users.user_id
						left join frm.form_teams ft on ft.team_id= team_users.team_id
						where ft.form_id = ? `+whrComp+`
						`+whreTime+` `+strField+`
						) 						
						ORDER BY `+orderBy+` `+limitOffset, formID, formID, formID, formID).Find(&data).Error

	if err != nil {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return data, nil
}

func (con *attendConnection) GetAttendanceRow(fields tables.Attendances, strField string) (tables.Attendances, error) {
	var data tables.Attendances
	err := con.db.Scopes(SchemaUsr("attendances")).Where(fields).Where(strField).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *attendConnection) GetAttendanceObjRow(fields tables.Attendances, strField string) (objects.Attendances, error) {
	var data objects.Attendances
	err := con.db.Table("usr.attendances").Select("attendances.id, attendances.form_id, attendances.user_id, attendances.attendance_in, attendances.attendance_out, ST_Y(attendances.geometry_in::geometry) as latitude, ST_X(attendances.geometry_in::geometry) as longitude").Where(fields).Where(strField).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *attendConnection) UpdateAttendance(id int, fields tables.Attendances, geo tables.Geometry) (bool, error) {
	err := con.db.Scopes(SchemaUsr("attendances")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		fmt.Println("ERROR :::: lat, long models update====", geo)
	}

	long := strconv.FormatFloat(geo.Longitude, 'f', 42, 64)
	lat := strconv.FormatFloat(geo.Latitude, 'f', 42, 64)

	err2 := con.db.Scopes(SchemaUsr("attendances")).Where("id = ?", id).Update("geometry_out", gorm.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", long, lat)).Error
	if err2 != nil {
		return false, err2
	}

	return true, nil
}

func (con *attendConnection) DeleteAttendance(AttendanceID int) (bool, error) {
	var data tables.Attendances
	err := con.db.Scopes(SchemaUsr("attendances")).Delete(&data, AttendanceID).Error
	if err != nil {
		return false, err
	}

	return true, err
}
func (con *attendConnection) GetUserAttendanceMapReports(fields tables.Attendances, strField string, whreTime string, paging objects.Paging) ([]tables.AttendenceReport, error) {

	return nil, nil
}

func (con *attendConnection) GetUserAttendanceMap2Reports(fields tables.Attendances, strField string, whreTime string, whrComp string) ([]tables.AttendenceMapReport, error) {

	// offset := (paging.Page - 1) * paging.Limit

	// orderBy := ""
	// if paging.SortBy != "" {
	// 	orderBy = paging.SortBy + " " + paging.Sort
	// }

	var data []tables.AttendenceMapReport
	err := con.db.Raw(`select * FROM
							(SELECT ROW_NUMBER() OVER(PARTITION BY form_users.id ORDER BY form_users.id DESC) as rownum,
								attendances.id,
								attendances.form_id,
								attendances.user_id,
								ao.organization_id,
								o.name as organization_name,
								u.name as user_name,
								u.phone as user_phone,
								u.avatar as user_avatar,
								to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at,
								to_char(attendances.attendance_in::timestamp, 'HH24:MI') as attendance,
								attendances.address_in as address,
								face_pic_in as face_pic,
								ST_Y(attendances.geometry_in::geometry) as latitude,
								ST_X(attendances.geometry_in::geometry) as longitude,
								true as is_checkin
							FROM frm.form_users
							left join usr.attendances on attendances.user_id = form_users.user_id 
								AND attendances.form_id = ?
								-- AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char('2023-01-12'::date, 'yyyy-mm-dd')
								`+strField+` `+whreTime+`
							left join usr.users u on u.id = form_users.user_id
							left join frm.forms f on f.id = form_users.form_id
							left join usr.attendance_organizations ao on attendances.id = ao.attendance_id
							left join mstr.organizations o on ao.organization_id = o.id
							WHERE form_users.form_id = ? 
							AND form_users.type = 'respondent' `+whrComp+`) as tb_1
							
							UNION 
							
							(SELECT ROW_NUMBER() OVER(PARTITION BY form_users.id ORDER BY form_users.id DESC)+1 as rownum,
								attendances.id,
								attendances.form_id,
								attendances.user_id,
								ao.organization_id,
								o.name as organization_name,
								u.name as user_name,
								u.phone as user_phone,
								u.avatar as user_avatar,
								to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at,
								to_char(attendances.attendance_out::timestamp, 'HH24:MI') as attendance,
								attendances.address_out as address,
								face_pic_out as face_pic,
								ST_Y(attendances.geometry_out::geometry) as latitude,
								ST_X(attendances.geometry_out::geometry) as longitude,
								false as is_checkin
							FROM frm.form_users
							left join usr.attendances on attendances.user_id = form_users.user_id 
								AND attendances.form_id = ?
								-- AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char('2023-01-12'::date, 'yyyy-mm-dd')
								`+strField+` `+whreTime+`
							left join usr.attendance_organizations ao on attendances.id = ao.attendance_id
							left join mstr.organizations o on ao.organization_id = o.id
							left join usr.users u on u.id = form_users.user_id
							left join frm.forms f on f.id = form_users.form_id
							WHERE form_users.form_id = ? 
							AND form_users.type = 'respondent' `+whrComp+`)
							order by user_name asc`, fields.FormID, fields.FormID, fields.FormID, fields.FormID).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *attendConnection) GetUserAttendanceMap3Reports(fields tables.Attendances, strField string, whreTime string, whrComp string) ([]tables.AttendenceMapReport, error) {

	var data []tables.AttendenceMapReport
	err := con.db.Raw(`select *
						FROM
						(SELECT ROW_NUMBER() OVER(PARTITION BY fu_merge.user_id ORDER BY fu_merge.user_id DESC) as rownum,
								attendances.id,
								attendances.form_id,
								attendances.user_id,
								u.name as user_name,
								u.phone as user_phone,
								u.avatar as user_avatar,
								to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at,
								to_char(attendances.attendance_in::timestamp, 'HH24:MI') as attendance,
								attendances.address_in as address,
								face_pic_in as face_pic,
								ST_Y(attendances.geometry_in::geometry) as latitude,
								ST_X(attendances.geometry_in::geometry) as longitude,
								true as is_checkin
						FROM
							(select *
							FROM
								(select attendances.user_id
								from usr.attendances
								where attendances.form_id = ?
								`+strField+` `+whreTime+` ) tb1
							UNION
								(select fu.user_id
								from frm.form_users fu
								where form_id= ?
								and fu.type = 'respondent' `+whrComp+`)) as fu_merge
						left join usr.attendances on attendances.user_id = fu_merge.user_id
						AND attendances.form_id = ?
						`+strField+` `+whreTime+`
						left join usr.users u on u.id = fu_merge.user_id 
						) as tbb1

						UNION

						(SELECT ROW_NUMBER() OVER(PARTITION BY fu_merge.user_id ORDER BY fu_merge.user_id DESC)+1 as rownum,
								attendances.id,
								attendances.form_id,
								attendances.user_id,
								u.name as user_name,
								u.phone as user_phone,
								u.avatar as user_avatar,
								to_char(attendances.created_at::timestamp, 'yyyy-mm-dd HH24:MI') as created_at,
								to_char(attendances.attendance_in::timestamp, 'HH24:MI') as attendance,
								attendances.address_in as address,
								face_pic_in as face_pic,
								ST_Y(attendances.geometry_in::geometry) as latitude,
								ST_X(attendances.geometry_in::geometry) as longitude,
								true as is_checkin
						FROM
							(select *
							FROM
								(select attendances.user_id
								from usr.attendances
								where attendances.form_id= ?
								`+strField+` `+whreTime+` ) tb1
							UNION
								(select fu.user_id
								from frm.form_users fu
								where form_id= ?
								and fu.type = 'respondent' `+whrComp+`)) as fu_merge
						left join usr.attendances on attendances.user_id = fu_merge.user_id
						AND attendances.form_id = ?
						`+strField+` `+whreTime+`
						left join usr.users u on u.id = fu_merge.user_id 
						)
						order by user_name asc`, fields.FormID, fields.FormID, fields.FormID, fields.FormID, fields.FormID, fields.FormID).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *attendConnection) GetLastAttendanceOverdate(formID int, userID int) (objects.LastAttendance, error) {

	var data objects.LastAttendance
	err := con.db.Raw(`SELECT t.id
						, t.user_id 
						, t.form_id 
						, t.attendance_in
						, t.attendance_out
					FROM usr.attendances t
					where t.user_id = ? and form_id = ?
					--and t.attendance_out is null 
					and to_char(t.attendance_in,'yyyy-mm-dd HH24:MI') BETWEEN (select to_char(f.attendance_overdate_at,'yyyy-mm-dd HH24:MI') 
																				from frm.forms f where f.id = t.form_id and f.id=?) 
																		AND to_char(now(),'yyyy-mm-dd HH24:MI')
					ORDER BY t.attendance_in desc limit 1`, userID, formID, formID).Find(&data).Error
	if err != nil {
		return objects.LastAttendance{}, err
	}

	return data, nil
}

func (con *attendConnection) GetMissingIDAtt(fields objects.MissingIDAtt) ([]objects.MissingIDAtt, error) {

	var data []objects.MissingIDAtt
	err := con.db.Table("usr.attendances").
		Select("attendances.id, attendances.form_id, attendances.user_id").
		Joins("left join usr.attendance_organizations ao ON attendances.id = ao.attendance_id").
		Where("ao.attendance_id IS null").
		Where(fields).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *attendConnection) GetCompanyID(form_id int, user_id int) (objects.FormUserOrg, error) {

	var data objects.FormUserOrg
	err := con.db.Table("frm.form_users").
		Select("form_users.id, form_users.user_id, form_users.form_id, fuo.organization_id").
		Joins("left join frm.form_user_organizations fuo on fuo.form_user_id = form_users.id").
		Where("form_users.form_id = ?", form_id).
		Where("form_users.user_id = ?", user_id).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return objects.FormUserOrg{}, err
	}
	return data, nil
}

func (con *attendConnection) GetCompanyIDByFormID(form_id int) (objects.GetOrgID, error) {

	var data objects.GetOrgID
	err := con.db.Table("frm.form_organizations").
		Select("form_organizations.organization_id").
		Where("form_organizations.form_id = ?", form_id).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return objects.GetOrgID{}, err
	}
	return data, nil
}

func (con *attendConnection) InsertAttendanceOrganization(fields objects.InsAttOrg) (objects.InsAttOrg, error) {

	err := con.db.Scopes(SchemaUsr("attendance_organizations")).Create(&fields).Error
	if err != nil {
		return objects.InsAttOrg{}, err
	}

	return fields, nil
}

func (con *attendConnection) GetFormAttendanceLocationRows(fields objects.ObjectFormAttendanceLocations) ([]objects.ObjectFormAttendanceLocations, error) {

	var data []objects.ObjectFormAttendanceLocations
	err := con.db.Table("frm.form_attendance_locations").
		Select(`form_attendance_locations.id,form_attendance_locations.form_id,form_attendance_locations.name,form_attendance_locations.location, 
		ST_Y(form_attendance_locations.geometry::geometry) latitude,
		ST_X(form_attendance_locations.geometry::geometry) longitude, form_attendance_locations.is_check_in,
		form_attendance_locations.is_check_out,form_attendance_locations.radius`).
		Where(fields).
		Find(&data).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *attendConnection) GetTeamByRespondent(respondenID int) ([]objects.TeamUsers, error) {

	var data []objects.TeamUsers
	err := con.db.Table("usr.team_users").
		Select("team_users.id, team_users.team_id, team_users.user_id").
		Where("team_users.user_id = ?", respondenID).
		Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *attendConnection) GetFormTeam(formID int, TeamID int) (objects.FormTeams, error) {

	var data objects.FormTeams
	err := con.db.Table("frm.form_teams").
		Select("form_teams.id, form_teams.team_id, form_teams.form_id, to2.organization_id").
		Joins("left join usr.team_organizations to2 on to2.team_id = form_teams.team_id").
		Where("form_teams.form_id = ?", formID).
		Where("form_teams.team_id = ?", TeamID).
		Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return objects.FormTeams{}, err
	}
	return data, nil
}
