package models

import (
	"errors"
	"snapin-form/objects"
	"snapin-form/tables"

	"gorm.io/gorm"
)

type RegionModels interface {
	GetListProvinces(fields tables.Province) ([]tables.Province, error)
	GetListCitiesByProvinceID(fields tables.Cities) ([]tables.Cities, error)
	GetListAllCities(fields tables.Cities) ([]tables.Cities, error)
	GetListDistrictsByProvinceID(fields tables.Districts) ([]tables.Districts, error)
	GetListAllDistricts(fields tables.Districts) ([]tables.Districts, error)
	GetListSubDistrictsByProvinceID(fields tables.SubDistricts) ([]tables.SubDistricts, error)
	GetListAllSubDistricts(fields tables.SubDistricts) ([]tables.SubDistricts, error)

	GetRadius(fields objects.Radius) (objects.Radius, error)
}

type regionConnection struct {
	db *gorm.DB
}

func NewRegionModels(dbg *gorm.DB) RegionModels {
	return &regionConnection{
		db: dbg,
	}
}

func (con *regionConnection) GetListProvinces(fields tables.Province) ([]tables.Province, error) {
	var data []tables.Province
	err := con.db.Scopes(SchemaMstr("provinces")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetListCitiesByProvinceID(fields tables.Cities) ([]tables.Cities, error) {
	var data []tables.Cities
	err := con.db.Scopes(SchemaMstr("cities")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetListAllCities(fields tables.Cities) ([]tables.Cities, error) {
	var data []tables.Cities
	err := con.db.Scopes(SchemaMstr("cities")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetListDistrictsByProvinceID(fields tables.Districts) ([]tables.Districts, error) {
	var data []tables.Districts
	err := con.db.Scopes(SchemaMstr("districts")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetListAllDistricts(fields tables.Districts) ([]tables.Districts, error) {
	var data []tables.Districts
	err := con.db.Scopes(SchemaMstr("districts")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetListSubDistrictsByProvinceID(fields tables.SubDistricts) ([]tables.SubDistricts, error) {
	var data []tables.SubDistricts
	err := con.db.Scopes(SchemaMstr("sub_districts")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetListAllSubDistricts(fields tables.SubDistricts) ([]tables.SubDistricts, error) {
	var data []tables.SubDistricts
	err := con.db.Scopes(SchemaMstr("sub_districts")).Where(fields).Find(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *regionConnection) GetRadius(fields objects.Radius) (objects.Radius, error) {
	err := con.db.Raw(`select fgr.distance, fgr.is_radius from mstr.f_get_radius(?,?,?) as fgr`, fields.Latitude, fields.Longitude, fields.LocationID).Find(&fields).Error
	return fields, err
}
