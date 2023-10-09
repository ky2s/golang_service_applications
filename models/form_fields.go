package models

import (
	"errors"
	"fmt"
	"snapin-form/tables"
	"strconv"

	"gorm.io/gorm"
)

type FormFieldModels interface {
	InsertFormField(data tables.FormFields) (tables.FormFields, error)
	UpdateFormField(fielID int, data tables.FormFields) (bool, error)
	UpdateFormFieldSort(fielID int, data tables.FormFields) (bool, error)
	UpdateFormOnly(id int, data tables.FormFields) (bool, error)
	GetFormFieldRow(data tables.FormFields) (tables.SelectFormFieldConditionRules, error)
	GetFormFieldRows(data tables.FormFields) ([]tables.SelectFormFieldConditionRules, error)
	GetFormFieldWhrRows(data tables.FormFields, whrStr string) ([]tables.SelectFormFieldConditionRules, error)
	GetFormFieldNotParentRows(fields tables.FormFields, whrStr string) ([]tables.SelectFormFieldConditionRules, error)
	DeleteFormField(fieldID int) (bool, error)
	InsertFormFieldPic(data tables.FormFieldPics) (tables.FormFieldPics, error)
	UpdateFormFieldPic(id int, data tables.FormFieldPics) (bool, error)
	DeleteFormFieldPic(id int) (bool, error)
	GetFormFieldPicRow(data tables.FormFieldPics) (tables.FormFieldPics, error)
	InsertFormFieldSection(data tables.FormFieldSection) (tables.FormFieldSection, error)
	// GetProjectFormRows(data tables.FormFields) ([]tables.SelectFormFieldConditionRules, error)
	// GetFormFieldParentBySortOrderRow(data tables.FormFields) (tables.SelectFormFieldConditionRules, error)
}

type formFieldConnection struct {
	db *gorm.DB
}

func NewFormFieldModels(dbg *gorm.DB) FormFieldModels {
	return &formFieldConnection{
		db: dbg,
	}
}

func (con *formFieldConnection) InsertFormField(data tables.FormFields) (tables.FormFields, error) {

	// insert form field
	err := con.db.Scopes(SchemaFrm("form_fields")).Create(&data).Error
	if err != nil {
		return tables.FormFields{}, err
	}

	if data.ID > 0 {

		con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", data.ID).Update("is_required", data.IsRequired)

		con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", data.ID).Update("is_multiple", data.IsMultiple)

		con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", data.ID).Update("is_section", data.IsSection)

		var fieldType tables.FieldTypes
		var fields tables.FieldTypes

		con.db.Scopes(SchemaMstr("field_types")).Where(fields).First(&fieldType)

		fmt.Println("fieldType--------->", fieldType.RealVarType)

		con.db.Exec(`ALTER TABLE "frm"."input_forms_` + strconv.Itoa(data.FormID) + `" ADD COLUMN "f` + strconv.Itoa(data.ID) + `" ` + fieldType.RealVarType + `;`)

		return data, nil
	} else {
		err := errors.New("Error: failed saved form fields")
		return tables.FormFields{}, err
	}
}

func (con *formFieldConnection) UpdateFormField(id int, fields tables.FormFields) (bool, error) {

	fmt.Println("fields--->", fields)
	err := con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	err = con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Update("is_multiple", fields.IsMultiple).Error
	if err != nil {
		return false, err
	}
	err = con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Update("is_required", fields.IsRequired).Error
	if err != nil {
		return false, err
	}
	err = con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Update("is_section", fields.IsSection).Error
	if err != nil {
		return false, err
	}

	if fields.Description == "" {
		err = con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Update("description", nil).Error
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (con *formFieldConnection) UpdateFormFieldSort(id int, fields tables.FormFields) (bool, error) {

	fmt.Println("fields--->", fields.SortOrder)
	err := con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *formFieldConnection) UpdateFormOnly(id int, fields tables.FormFields) (bool, error) {

	fmt.Println("fields--->", fields)
	err := con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *formFieldConnection) GetFormFieldRows(fields tables.FormFields) ([]tables.SelectFormFieldConditionRules, error) {

	var data []tables.SelectFormFieldConditionRules

	where := ``
	whereParent := ``
	whereFieldType := ``

	if fields.ID > 0 {
		where = `AND form_fields.id = ` + strconv.Itoa(fields.ID)
	}

	if fields.FormID > 0 {
		where = `AND form_fields.form_id = ` + strconv.Itoa(fields.FormID)
	}

	if fields.ParentID > 0 {
		whereParent = `AND form_fields.parent_id = ` + strconv.Itoa(fields.ParentID)
	}

	if fields.ParentID == 0 {
		whereParent = `AND form_fields.parent_id is null`
	}

	if fields.FieldTypeID > 0 {
		whereFieldType = `AND form_fields.field_type_id = ` + strconv.Itoa(fields.FieldTypeID)
	} else if fields.FieldTypeID == -1 {
		whereFieldType = `AND form_fields.field_type_id is null`
	} else if fields.FieldTypeID == -2 {
		whereFieldType = `AND form_fields.field_type_id is not null`
	}

	// if fields.IsSection == true {
	// 	whereParent = `AND form_fields.is_section is true`
	// } else {
	// 	whereParent = `AND form_fields.is_section is false`
	// }

	err := con.db.Raw(`SELECT form_fields.id, form_fields.parent_id, form_fields.form_id, form_fields.field_type_id, t.translation as field_type_name, form_fields.label, form_fields.description,
				form_fields.option, form_fields.condition_type, form_fields.upperlower_case_type, form_fields.is_multiple, form_fields.is_required, form_fields.is_section, form_fields.section_color
				, ffcr.condition_rule_id, ffcr.value1, ffcr.value2, ffcr.err_msg, ffcr.condition_parent_field_id, ffcr.condition_all_right, ffcr.tab_max_one_per_line, ffcr.tab_each_line_require
				, form_fields.is_condition, form_fields.is_section, form_fields.section_color, form_fields.sort_order, form_fields.tag_loc_color, form_fields.tag_loc_icon
				, ffp.pic as image, form_fields.address_type, form_fields.province_id, form_fields.city_id, form_fields.district_id, form_fields.sub_district_id, form_fields.currency_type, form_fields.currency 
				
			FROM frm.form_fields
			
			left join (SELECT ROW_NUMBER() OVER(PARTITION BY form_field_id ORDER BY id DESC) as rownum
				, form_field_id, condition_rule_id, value1, value2, err_msg, condition_parent_field_id, condition_all_right,tab_max_one_per_line, tab_each_line_require
				
				FROM frm.form_field_condition_rules f1 where condition_parent_field_id is null
				
				) as ffcr on ffcr.form_field_id = form_fields.id and ffcr.rownum = 1
			
			LEFT JOIN mstr.field_types ft on ft.id = form_fields.field_type_id
			LEFT JOIN mstr.translations t on t.textcontent_id = ft.name_textcontent_id AND t.language_id = 1
			LEFT join frm.form_field_pics ffp on ffp.form_field_id = form_fields.id
			WHERE form_fields.deleted_at IS NULL ` + where + ` ` + whereParent + ` ` + whereFieldType + `
			
			order by sort_order asc, id asc
			`).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *formFieldConnection) GetFormFieldNotParentRows(fields tables.FormFields, whrStr string) ([]tables.SelectFormFieldConditionRules, error) {

	var data []tables.SelectFormFieldConditionRules

	where := ``
	whereParent := ``
	whereFieldType := ``

	if fields.ID > 0 {
		where = `AND form_fields.id = ` + strconv.Itoa(fields.ID)
	}

	if fields.FormID > 0 {
		where = `AND form_fields.form_id = ` + strconv.Itoa(fields.FormID)
	}

	if fields.ParentID > 0 {
		whereParent = `AND form_fields.parent_id = ` + strconv.Itoa(fields.ParentID)
	}

	if fields.FieldTypeID > 0 {
		whereFieldType = `AND form_fields.field_type_id = ` + strconv.Itoa(fields.FieldTypeID)
	} else if fields.FieldTypeID == -1 {
		whereFieldType = `AND form_fields.field_type_id is null`
	} else if fields.FieldTypeID == -2 {
		whereFieldType = `AND form_fields.field_type_id is not null`
	}

	// if fields.IsSection == true {
	// 	whereParent = `AND form_fields.is_section is true`
	// } else {
	// 	whereParent = `AND form_fields.is_section is false`
	// }

	err := con.db.Raw(`SELECT form_fields.id, form_fields.parent_id, form_fields.form_id, form_fields.field_type_id, t.translation as field_type_name, form_fields.label, form_fields.description,
				form_fields.option, form_fields.condition_type, form_fields.upperlower_case_type, form_fields.is_multiple, form_fields.is_required, form_fields.is_section, form_fields.section_color
				, form_fields.is_country_phone_code, form_fields.province_id, form_fields.city_id, form_fields.district_id, form_fields.sub_district_id , form_fields.currency_type, form_fields.currency, form_fields.address_type, ffcr.condition_rule_id, ffcr.value1, ffcr.value2, ffcr.err_msg, ffcr.condition_parent_field_id, ffcr.condition_all_right, ffcr.tab_max_one_per_line, ffcr.tab_each_line_require
				, form_fields.is_condition, form_fields.is_section, form_fields.section_color, form_fields.sort_order
				, ffp.pic as image
				
			FROM frm.form_fields
			
			left join (SELECT ROW_NUMBER() OVER(PARTITION BY form_field_id ORDER BY id DESC) as rownum
				, form_field_id, condition_rule_id, value1, value2, err_msg, condition_parent_field_id, condition_all_right,tab_max_one_per_line, tab_each_line_require
				
				FROM frm.form_field_condition_rules f1 where condition_parent_field_id is null
				
				) as ffcr on ffcr.form_field_id = form_fields.id and ffcr.rownum = 1
			
			LEFT JOIN mstr.field_types ft on ft.id = form_fields.field_type_id
			LEFT JOIN mstr.translations t on t.textcontent_id = ft.name_textcontent_id AND t.language_id = 1
			LEFT join frm.form_field_pics ffp on ffp.form_field_id = form_fields.id
			WHERE form_fields.deleted_at IS NULL ` + where + ` ` + whereParent + ` ` + whereFieldType + ` ` + whrStr + `

			order by sort_order asc, id asc
			`).Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (con *formFieldConnection) GetFormFieldRow(fields tables.FormFields) (tables.SelectFormFieldConditionRules, error) {

	var data tables.SelectFormFieldConditionRules

	where := ``
	whereParent := ``
	whereSortOrder := ``

	if fields.ID > 0 {
		where = `AND form_fields.id = ` + strconv.Itoa(fields.ID)
	}

	if fields.FormID > 0 {
		where = `AND form_fields.form_id = ` + strconv.Itoa(fields.FormID)
	}

	if fields.ParentID > 0 {
		whereParent = `AND form_fields.parent_id = ` + strconv.Itoa(fields.ParentID)
	}

	if fields.ParentID == 0 {
		whereParent = `AND form_fields.parent_id is null`
	}

	if fields.SortOrder > 0 {
		whereSortOrder = `AND form_fields.sort_order = ` + strconv.Itoa(fields.SortOrder)
	}

	// if fields.IsSection == true {
	// 	whereParent = `AND form_fields.is_section is true`
	// } else {
	// 	whereParent = `AND form_fields.is_section is false`
	// }

	err := con.db.Raw(`SELECT form_fields.id,form_fields.parent_id, form_fields.form_id, form_fields.field_type_id, form_fields.label, 
				form_fields.description,
				form_fields.option,form_fields.condition_type,form_fields.upperlower_case_type,form_fields.is_multiple,form_fields.is_required,form_fields.is_section, form_fields.section_color
				, ffcr.condition_rule_id, ffcr.value1, ffcr.value2, ffcr.err_msg, ffcr.condition_parent_field_id
				, form_fields.is_condition, form_fields.is_section, form_fields.section_color, form_fields.sort_order
				, ffp.pic as image
				
			FROM frm.form_fields
			
			left join (SELECT ROW_NUMBER() OVER(PARTITION BY form_field_id ORDER BY id DESC) as rownum
				, form_field_id, condition_rule_id, value1, value2, err_msg, condition_parent_field_id
				
				FROM frm.form_field_condition_rules f1 where condition_parent_field_id is null
				
				) as ffcr on ffcr.form_field_id = form_fields.id and ffcr.rownum = 1
			left join frm.form_field_pics ffp on ffp.form_field_id = form_fields.id
			WHERE form_fields.deleted_at IS NULL ` + where + ` ` + whereParent + ` ` + whereSortOrder + `
			
			order by sort_order asc, id asc
			`).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.SelectFormFieldConditionRules{}, err
	}
	return data, nil
}

func (con *formFieldConnection) GetFormFieldRow___(fields tables.FormFields) (tables.FormFields, error) {
	var data tables.FormFields
	err := con.db.Scopes(SchemaFrm("form_fields")).Where(fields).First(&data).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *formFieldConnection) DeleteFormField(fieldID int) (bool, error) {

	var formField tables.FormFields
	// delete child
	err := con.db.Scopes(SchemaFrm("form_fields")).Where("parent_id = ?", fieldID).Delete(&formField).Error
	if err != nil {
		return false, err
	}

	// delete parent
	err = con.db.Scopes(SchemaFrm("form_fields")).Delete(&formField, fieldID).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formFieldConnection) InsertFormFieldPic(data tables.FormFieldPics) (tables.FormFieldPics, error) {

	err := con.db.Scopes(SchemaFrm("form_field_pics")).Create(&data).Error
	if err != nil {
		return tables.FormFieldPics{}, err
	}

	return data, err
}

func (con *formFieldConnection) UpdateFormFieldPic(id int, data tables.FormFieldPics) (bool, error) {

	err := con.db.Scopes(SchemaFrm("form_field_pics")).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formFieldConnection) DeleteFormFieldPic(fieldID int) (bool, error) {

	err := con.db.Exec("DELETE FROM frm.form_field_pics where form_field_id = " + strconv.Itoa(fieldID)).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *formFieldConnection) GetFormFieldPicRow(fields tables.FormFieldPics) (tables.FormFieldPics, error) {

	var data tables.FormFieldPics
	err := con.db.Scopes(SchemaFrm("form_field_pics")).Where(fields).First(&data).Error
	if err != nil {
		return tables.FormFieldPics{}, err
	}
	return data, nil
}

func (con *formFieldConnection) InsertFormFieldSection(data tables.FormFieldSection) (tables.FormFieldSection, error) {

	err := con.db.Scopes(SchemaFrm("form_field_sections")).Create(&data).Error
	if err != nil {
		return tables.FormFieldSection{}, err
	}

	return data, err
}

func (con *formFieldConnection) GetFormFieldWhrRows(fields tables.FormFields, whreStr string) ([]tables.SelectFormFieldConditionRules, error) {

	var data []tables.SelectFormFieldConditionRules

	where := ``
	whereParent := ``
	whereFieldType := ``

	if fields.ID > 0 {
		where = `AND form_fields.id = ` + strconv.Itoa(fields.ID)
	}

	if fields.FormID > 0 {
		where = `AND form_fields.form_id = ` + strconv.Itoa(fields.FormID)
	}

	if fields.ParentID > 0 {
		whereParent = `AND form_fields.parent_id = ` + strconv.Itoa(fields.ParentID)
	}

	if fields.ParentID == 0 {
		whereParent = `AND form_fields.parent_id is null`
	}

	if fields.FieldTypeID > 0 {
		whereFieldType = `AND form_fields.field_type_id = ` + strconv.Itoa(fields.FieldTypeID)
	} else if fields.FieldTypeID == -1 {
		whereFieldType = `AND form_fields.field_type_id is null`
	} else if fields.FieldTypeID == -2 {
		whereFieldType = `AND form_fields.field_type_id is not null`
	}

	// if fields.IsSection == true {
	// 	whereParent = `AND form_fields.is_section is true`
	// } else {
	// 	whereParent = `AND form_fields.is_section is false`
	// }

	err := con.db.Raw(`SELECT form_fields.id, form_fields.parent_id, form_fields.form_id, form_fields.field_type_id, t.translation as field_type_name, form_fields.label, form_fields.description,
				form_fields.option, form_fields.condition_type, form_fields.upperlower_case_type, form_fields.is_multiple, form_fields.is_required, form_fields.is_section, form_fields.section_color
				, ffcr.condition_rule_id, ffcr.value1, ffcr.value2, ffcr.err_msg, ffcr.condition_parent_field_id, ffcr.condition_all_right, ffcr.tab_max_one_per_line, ffcr.tab_each_line_require
				, form_fields.is_condition, form_fields.is_section, form_fields.section_color, form_fields.sort_order, form_fields.tag_loc_color, form_fields.tag_loc_icon
				, ffp.pic as image
				
			FROM frm.form_fields
			
			left join (SELECT ROW_NUMBER() OVER(PARTITION BY form_field_id ORDER BY id DESC) as rownum
				, form_field_id, condition_rule_id, value1, value2, err_msg, condition_parent_field_id, condition_all_right,tab_max_one_per_line, tab_each_line_require
				
				FROM frm.form_field_condition_rules f1 where condition_parent_field_id is null
				
				) as ffcr on ffcr.form_field_id = form_fields.id and ffcr.rownum = 1
			
			LEFT JOIN mstr.field_types ft on ft.id = form_fields.field_type_id
			LEFT JOIN mstr.translations t on t.textcontent_id = ft.name_textcontent_id AND t.language_id = 1
			LEFT join frm.form_field_pics ffp on ffp.form_field_id = form_fields.id
			WHERE form_fields.deleted_at IS NULL ` + where + ` ` + whereParent + ` ` + whereFieldType + ` ` + whreStr + `
			
			order by sort_order asc, id asc
			`).Find(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}
