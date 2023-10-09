package models

import (
	"errors"
	"fmt"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"

	"github.com/jackc/pgconn"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/gorm"
)

type InputFormModels interface {
	InsertFormData(data objects.FormData) (bool, error)
	InsertFormDataWithOrganization(data objects.FormData, organizationID int) (bool, error)
	UpdateFormData(submissionID int, data objects.FormData) (bool, error)
	DeleteFormData(submissionID int, data objects.FormData) (bool, error)
	GetInputFormRows(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]tables.InputForms, error)
	GetInputFormUnscopedRows(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]tables.InputForms, error) // query menampilkan data yg sdh di softdelete
	GetInputFormRow(formID int, fields tables.InputForms, whreStr string) (tables.InputForms, error)
	GetInputDataRows(formID int, fstring string, fields tables.InputForms, whreStr string) ([][]string, error)
	GetInputDataUnscopedRows(formID int, fstring string, fields tables.InputForms, whreStr string) ([][]string, error)
	GetDates(strWhere string) ([]tables.Date, error)
	GetDatesNew(strWhere string) ([]tables.Date, error)
	GetDatesWithFilter(month string, year string, strWhere string) ([]tables.Date, error)
	GetMonths(strWhere string) ([]tables.Months, error)
	GetTotalDate(userID int, formID int, strWhre string) ([]tables.TotalDate, error)
	GetTotalDateMonthly(userID int, formID int, strWhre string) ([]tables.TotalDate, error)
	GetTotalMonth(userID int, formID int, strWhre string) ([]tables.TotalDate, error)
	GetDataHours(formID int, whreString string) ([]objects.GraficDataHours, error)
	GetDataPeriodeDays(formID int, whreString string) ([]objects.GraficDataPeriod, error)
	GetDataPeriodeMonthly(formID int, whreString string) ([]objects.GraficDataPeriod, error)
	GetDataPeriodeYearly(formID int, whreString string) ([]objects.GraficDataPeriod, error)
	GetActiveUserInputForm(formID int, fields tables.InputForms, whreStr string) ([]tables.JoinFormUsers, error)
	GetDataHoursUserResp(formID int, whreString string) ([]objects.GraficDataHours, error)
	GetDataPeriodeDaysResp(formID int, whreString string) ([]objects.GraficDataPeriod, error)
	GetDataPeriodeMonthlyResp(formID int, whreString string) ([]objects.GraficDataPeriod, error)
	GetDataPeriodeYearlyResp(formID int, whreString string) ([]objects.GraficDataPeriod, error)
	GetReportFormResponden(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]objects.ReportResponden, error)
	GetReportFormRespondenUnionTeam(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]objects.ReportResponden, error)
	GetInputFormCustomAnswerRow(formID int, form_field_id int, input_form_id int) (tables.InputFormCustomAnswers, error)

	GetInputFormOrganizationRows(formID int, fields tables.InputFormJoinOrganizations, whreStr string, paging objects.Paging) ([]tables.InputFormJoinOrganizations, error)
	GetInputDataOrganizationRows(formID int, fstring string, fields tables.InputForms, whreStr string) ([][]string, error)
	GetOrganizationInputForm(fields tables.InputFormOrganizations, whreStr string) ([]objects.InputFormOrganizations, error)

	InsertInputFormOrgData(data tables.InputFormOrganizations) (tables.InputFormOrganizations, error)

	GetUpdatedCount(formID int, submissionID int) (tables.InputForms, error)
	UpdatedCount(formID int, submissionID int, fields objects.UpdCnt) (bool, objects.UpdCnt, error)
	GetDeletedData(formID int, fields tables.InputForms, whreStr string) ([]tables.InputForms, error)
	GetUpdatedData(formID int, userID int, whreStr string) ([]tables.InputForms, error)
	GetUpdatedDataForOrg(formID int) ([]tables.InputForms, error)
	GetUpdatedDataNoDate(formID int, userID int) ([]tables.InputForms, error)
	GetUpdatedDataWithDate(formID int, whreStr string) ([]tables.InputForms, error)

	GetActiveUserInputFormOld(formID int, fields tables.InputForms, whreStr string) ([]tables.JoinFormUsers, error)
}

type inputFormConnection struct {
	db    *gorm.DB
	pgErr *pgconn.PgError
}

func NewInputFormModels(dbg *gorm.DB) InputFormModels {
	return &inputFormConnection{
		db: dbg,
	}
}

func (con *inputFormConnection) InsertFormData(data objects.FormData) (bool, error) {

	long := strconv.FormatFloat(data.Longitude, 'f', 42, 64)
	lat := strconv.FormatFloat(data.Latitude, 'f', 42, 64)

	fieldData := map[string]interface{}{}
	fieldData["created_at"] = "now()"
	fieldData["updated_at"] = "now()"
	fieldData["created_by"] = data.UserID
	fieldData["user_id"] = data.UserID
	fieldData["address"] = data.Address
	// fieldData["is_pending"] = data.IsPending
	fieldData["geometry"] = gorm.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", long, lat)

	if len(data.FieldData) > 0 {

		// var formFields []tables.FormFields
		// err := con.db.Scopes(SchemaFrm("form_fields")).Where("form_id = ?", data.FormID).Find(&formFields).Error
		// if err != nil {
		// 	return false, err
		// }

		// for j := 0; j < len(formFields); j++ {
		// 	fmt.Println("ini fields :: ----", j, formFields[j].ID)
		// }

		shortenUrlModels := shortenurlConnection{}
		for i := 0; i < len(data.FieldData); i++ {

			//cek field mandatory here
			// var formField tables.FormFields
			// err := con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", data.FieldData[i].FieldID).Find(&formField).Error
			// if err != nil {
			// 	return false, err
			// }

			// if formField.IsRequired == true && data.FieldData[i].Answer == "" {
			// 	return false, errors.New("Terjadi kesalahan\n(Kode Error: 507)")
			// }

			// checking field URL for shorting
			fmt.Println(shortenUrlModels)
			/*
				var formField tables.FormFields
				err := con.db.Scopes(SchemaFrm("form_fields as ff")).Select("ff.field_type_id").Where("ff.id = ?", data.FieldData[i].FieldID).Find(&formField).Error
				if err != nil {
					return false, err
				}

				var urlShorten objects.ShortURLResponse
				if formField.FieldTypeID == 10 || formField.FieldTypeID == 18 || formField.FieldTypeID == 20 {
					urlShorten, err = shortenUrlModels.ShortenImageURL(data.FieldData[i].Answer)
					if err != nil {
						return false, err
					}

					sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
					fieldData[sFieldID] = urlShorten.Data.ShortURL
				} else {
					sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
					fieldData[sFieldID] = data.FieldData[i].Answer
				}
			*/

			sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
			fieldData[sFieldID] = data.FieldData[i].Answer
		}
	}

	err := con.db.Table("frm.input_forms_" + strconv.Itoa(data.FormID)).Create(fieldData).Error
	if err != nil {
		fmt.Println("tx-------------------", err)
		return false, err
	}
	var inptFrm tables.InputForms
	err = con.db.Table("frm.input_forms_" + strconv.Itoa(data.FormID)).Last(&inptFrm).Error
	if err != nil {
		fmt.Println("tx-------------------", err)
		return false, err
	}
	fmt.Println("last id-------------------", inptFrm)

	// save custom option
	var customField tables.InputFormCustomAnswers
	customField.FormID = data.FormID
	customField.InputFormID = inptFrm.ID

	if len(data.FieldData) > 0 {
		for j := 0; j < len(data.FieldData); j++ {

			if data.FieldData[j].CustomOption != "" {
				customField.FormFieldID = data.FieldData[j].FieldID
				customField.CustomAnswer = data.FieldData[j].CustomOption
				err := con.db.Table("frm.input_form_custom_answers").Create(&customField).Error
				if err != nil {
					fmt.Println("tx-------------------", err)
					return false, err
				}
			}
		}
	}

	return true, nil
}

func (con *inputFormConnection) InsertFormDataWithOrganization(data objects.FormData, organizationID int) (bool, error) {

	long := strconv.FormatFloat(data.Longitude, 'f', 42, 64)
	lat := strconv.FormatFloat(data.Latitude, 'f', 42, 64)

	fieldData := map[string]interface{}{}
	fieldData["created_at"] = "now()"
	fieldData["updated_at"] = "now()"
	fieldData["created_by"] = data.UserID
	fieldData["user_id"] = data.UserID
	fieldData["address"] = data.Address
	// fieldData["is_pending"] = data.IsPending
	fieldData["geometry"] = gorm.Expr("ST_SetSRID(ST_MakePoint(?, ?), 4326)", long, lat)

	if len(data.FieldData) > 0 {

		// var formFields []tables.FormFields
		// err := con.db.Scopes(SchemaFrm("form_fields")).Where("form_id = ?", data.FormID).Find(&formFields).Error
		// if err != nil {
		// 	return false, err
		// }

		// for j := 0; j < len(formFields); j++ {
		// 	fmt.Println("ini fields :: ----", j, formFields[j].ID)
		// }

		shortenUrlModels := shortenurlConnection{}
		for i := 0; i < len(data.FieldData); i++ {

			//cek field mandatory here
			// var formField tables.FormFields
			// err := con.db.Scopes(SchemaFrm("form_fields")).Where("id = ?", data.FieldData[i].FieldID).Find(&formField).Error
			// if err != nil {
			// 	return false, err
			// }

			// if formField.IsRequired == true && data.FieldData[i].Answer == "" {
			// 	return false, errors.New("Terjadi kesalahan\n(Kode Error: 507)")
			// }

			// checking field URL for shorting
			// fmt.Println(shortenUrlModels)

			var formField tables.FormFields
			err := con.db.Scopes(SchemaFrm("form_fields as ff")).Select("ff.field_type_id").Where("ff.id = ?", data.FieldData[i].FieldID).Find(&formField).Error
			if err != nil {
				return false, err
			}

			var urlShorten objects.ShortURLResponse
			if formField.FieldTypeID == 10 || formField.FieldTypeID == 18 || formField.FieldTypeID == 20 {
				urlShorten, err = shortenUrlModels.ShortenImageURL(data.FieldData[i].Answer)
				if err != nil {
					return false, err
				}
				// fmt.Println(urlShorten)
				// os.Exit(0)
				sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
				fieldData[sFieldID] = urlShorten.Data.ShortURL
				// fmt.Println(urlShorten.Data.ShortURL)
				// os.Exit(0)
			} else {
				sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
				fieldData[sFieldID] = data.FieldData[i].Answer
			}

			// sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
			// fieldData[sFieldID] = data.FieldData[i].Answer
		}
	}

	err := con.db.Table("frm.input_forms_" + strconv.Itoa(data.FormID)).Create(fieldData).Error
	if err != nil {
		fmt.Println("tx-------------------", err)
		return false, err
	}

	var inptFrm tables.InputForms
	err = con.db.Table("frm.input_forms_" + strconv.Itoa(data.FormID)).Last(&inptFrm).Error
	if err != nil {
		fmt.Println("tx-------------------", err)
		return false, err
	}
	fmt.Println("last id-------------------", inptFrm)

	// save input_form organization (pembeda submit input per company)
	if organizationID >= 1 {
		var inputFormOrg tables.InputFormOrganization
		inputFormOrg.FormID = data.FormID
		inputFormOrg.OrganizationID = organizationID
		inputFormOrg.InputFormID = inptFrm.ID
		err := con.db.Table("frm.input_form_organizations").Create(&inputFormOrg).Error
		if err != nil {
			fmt.Println("tx-------------------", err)
			return false, err
		}
	}

	// save custom option
	var customField tables.InputFormCustomAnswers
	customField.FormID = data.FormID
	customField.InputFormID = inptFrm.ID

	if len(data.FieldData) > 0 {
		for j := 0; j < len(data.FieldData); j++ {

			if data.FieldData[j].CustomOption != "" {
				customField.FormFieldID = data.FieldData[j].FieldID
				customField.CustomAnswer = data.FieldData[j].CustomOption
				err := con.db.Table("frm.input_form_custom_answers").Create(&customField).Error
				if err != nil {
					fmt.Println("tx-------------------", err)
					return false, err
				}
			}
		}
	}

	return true, nil
}

func (con *inputFormConnection) UpdateFormData(submissionID int, data objects.FormData) (bool, error) {

	fieldData := map[string]interface{}{}
	fieldData["updated_at"] = "now()"
	fieldData["updated_by"] = data.UserID

	if len(data.FieldData) > 0 {

		shortenUrlModels := shortenurlConnection{}
		for i := 0; i < len(data.FieldData); i++ {

			sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
			fieldData[sFieldID] = data.FieldData[i].Answer

			var formField tables.FormFields
			err := con.db.Scopes(SchemaFrm("form_fields as ff")).Select("ff.field_type_id").Where("ff.id = ?", data.FieldData[i].FieldID).Find(&formField).Error
			if err != nil {
				return false, err
			}

			var urlShorten objects.ShortURLResponse
			if formField.FieldTypeID == 10 || formField.FieldTypeID == 18 || formField.FieldTypeID == 20 {
				urlShorten, err = shortenUrlModels.ShortenImageURL(data.FieldData[i].Answer)
				if err != nil {
					return false, err
				}
				// fmt.Println(urlShorten)
				// os.Exit(0)
				sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
				fieldData[sFieldID] = urlShorten.Data.ShortURL
				// fmt.Println(urlShorten.Data.ShortURL)
				// os.Exit(0)
			} else {
				sFieldID := "f" + strconv.Itoa(data.FieldData[i].FieldID)
				fieldData[sFieldID] = data.FieldData[i].Answer
			}
		}
	}

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(data.FormID)).Where("id = ?", submissionID).Updates(fieldData).Error
	if err != nil {
		fmt.Println("tx-------------------", err)
		return false, err
	}

	// save custom option
	var customField tables.InputFormCustomAnswers
	if len(data.FieldData) > 0 {
		for j := 0; j < len(data.FieldData); j++ {

			if data.FieldData[j].CustomOption != "" {
				customField.CustomAnswer = data.FieldData[j].CustomOption
				err := con.db.Table("frm.input_form_custom_answers").Where("input_form_id = ? AND form_id = ?", submissionID, data.FormID).Updates(&customField).Error
				if err != nil {
					fmt.Println("tx-------------------", err)
					return false, err
				}
			}
		}
	}

	return true, nil
}

func (con *inputFormConnection) DeleteFormData(submissionID int, data objects.FormData) (bool, error) {

	fieldData := map[string]interface{}{}
	fieldData["deleted_at"] = "now()"
	fieldData["deleted_by"] = data.UserID

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(data.FormID)).Where("id = ?", submissionID).Updates(fieldData).Error
	if err != nil {
		fmt.Println("tx-------DELETE data------------", err)
		return false, err
	}

	return true, nil
}

func (con *inputFormConnection) InsertFormData___(data objects.FormData) (bool, error) {

	con.db.Table("frm.input_forms_106").Create(map[string]interface{}{
		"created_at": "now()",
		"updated_at": "now()",
		"created_by": 41,
		"user_id":    41,
		"f121":       "kyky sukiawan",
	})

	return true, nil
}

func (con *inputFormConnection) GetInputFormRows(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]tables.InputForms, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.InputForms

	// err := con.db.Scopes(SchemaFrm("input_forms_" + strconv.Itoa(formID))).Where(fields).Find(&data).Error

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(formID)+" as if").
		Select("if.id", "if.user_id", "u.name as user_name", "u.phone", "u.avatar", "if.address", "ST_Y(if.geometry::geometry) as latitude", "ST_X(if.geometry::geometry) as longitude", "if.created_at", "if.updated_at, if.updated_count").
		Joins("join usr.users u on u.id = if.user_id").
		Joins("left join frm.form_users fu on if.user_id = fu.user_id AND fu.form_id =" + strconv.Itoa(formID)).
		Joins("left join frm.form_user_organizations fuo on fu.id = fuo.form_user_id").
		Joins("left join frm.input_form_organizations ifo on ifo.input_form_id=if.id AND ifo.form_id=" + strconv.Itoa(formID)).
		Where(fields).Where(whreStr).Order(orderBy).Order("id desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetDeletedData(formID int, fields tables.InputForms, whreStr string) ([]tables.InputForms, error) {
	var data []tables.InputForms
	err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").
		Select("if.id").
		Where(fields).Where(whreStr).
		Where("if.deleted_at is not null").Unscoped().Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetUpdatedData(formID int, userID int, whreStr string) ([]tables.InputForms, error) {
	fID := strconv.Itoa(formID)
	uID := strconv.Itoa(userID)

	var data []tables.InputForms
	err := con.db.Raw(`SELECT sum(updated_count) as updated_count
						from frm.input_forms_` + fID + ` as if
						where ` + whreStr + ` AND if.user_id = ` + uID + ``).Unscoped().Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetUpdatedDataNoDate(formID int, userID int) ([]tables.InputForms, error) {
	fID := strconv.Itoa(formID)
	uID := strconv.Itoa(userID)

	var data []tables.InputForms
	err := con.db.Raw(`SELECT sum(updated_count) as updated_count
						from frm.input_forms_` + fID + ` as if
						where if.user_id = ` + uID + ``).Unscoped().Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetUpdatedDataForOrg(formID int) ([]tables.InputForms, error) {
	fID := strconv.Itoa(formID)

	var data []tables.InputForms
	err := con.db.Raw(`SELECT sum(updated_count) as updated_count
						from frm.input_forms_` + fID + ` as if`).Unscoped().Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetUpdatedDataWithDate(formID int, whreStr string) ([]tables.InputForms, error) {
	fID := strconv.Itoa(formID)

	var data []tables.InputForms
	err := con.db.Raw(`SELECT sum(updated_count) as updated_count
						from frm.input_forms_` + fID + ` as if
						where ` + whreStr + ``).Unscoped().Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetUpdatedCount(formID int, submissionID int) (tables.InputForms, error) {

	var data tables.InputForms

	// err := con.db.Scopes(SchemaFrm("input_forms_" + strconv.Itoa(formID))).Where(fields).Find(&data).Error

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(formID)+" as if").
		Select("if.id, if.updated_count").
		Where("if.id = ?", submissionID).
		Find(&data).Error

	if err != nil {
		return tables.InputForms{}, err
	}
	return data, nil
}

func (con *inputFormConnection) UpdatedCount(formID int, submissionID int, data objects.UpdCnt) (bool, objects.UpdCnt, error) {

	query := "UPDATE frm.input_forms_" + strconv.Itoa(formID) + " SET updated_count = ?  WHERE id = ?"
	err := con.db.Exec(query, gorm.Expr("updated_count + ?", data.UpdatedCount), submissionID).Error
	if err != nil {
		return false, objects.UpdCnt{}, err
	}

	return true, data, nil
}

func (con *inputFormConnection) GetInputFormUnscopeRows(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]tables.InputForms, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.InputForms

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(formID)+" as if").Select("if.id", "if.user_id", "u.name as user_name", "u.phone", "u.avatar", "if.address", "ST_Y(if.geometry::geometry) as latitude", "ST_X(if.geometry::geometry) as longitude", "if.created_at", "if.updated_at").Joins("join usr.users u on u.id = if.user_id").Where(fields).Where(whreStr).Unscoped().Order(orderBy).Order("id desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetInputFormUnscopedRows(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]tables.InputForms, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.InputForms

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(formID)+" as if").Select("if.id", "if.user_id", "if.updated_count", "u.name as user_name", "u.phone", "u.avatar", "if.address", "ST_Y(if.geometry::geometry) as latitude", "ST_X(if.geometry::geometry) as longitude", "if.created_at", "if.updated_at", "if.deleted_at, if.deleted_by").Joins("join usr.users u on u.id = if.user_id").Joins("left join frm.input_form_organizations ifo on ifo.input_form_id=if.id AND ifo.form_id=" + strconv.Itoa(formID)).Where(fields).Where(whreStr).Order(orderBy).Order("id desc").Limit(paging.Limit).Offset(offset).Unscoped().Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("err err record not")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetInputFormObjectRows(formID int, fields tables.InputForms, whreStr string) ([]objects.InputForms, error) {

	var data []objects.InputForms

	// err := con.db.Scopes(SchemaFrm("input_forms_" + strconv.Itoa(formID))).Where(fields).Find(&data).Error

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(formID)+" as if").Select("if.id", "if.user_id", "u.name as user_name", "u.phone", "if.created_at").Joins("join usr.users u on u.id = if.user_id").Where(fields).Where(whreStr).Order("id desc").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetInputDataRows(formID int, fieldString string, fields tables.InputForms, whreStr string) ([][]string, error) {

	// var result []tables.InputForms

	// err := con.db.Scopes(SchemaFrm("input_forms_" + strconv.Itoa(formID))).Where(fields).Find(&data).Error

	rows, err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").Select(fieldString).Where(fields).Where(whreStr).Rows()
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()

	length := len(cols)

	//----
	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]

		}

		err = rows.Scan(pointers...)
		if err != nil {
			fmt.Println("err pointers-----------", err)
			panic(err.Error())
		}

		result = append(result, container)
	}

	return result, nil
}

func (con *inputFormConnection) GetInputDataUnscopedRows(formID int, fieldString string, fields tables.InputForms, whreStr string) ([][]string, error) {

	rows, err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").Select(fieldString).Unscoped().Where(fields).Where(whreStr).Rows()
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()

	length := len(cols)

	//----
	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]

		}

		err = rows.Scan(pointers...)
		if err != nil {
			fmt.Println("err pointers-----------", err)
			panic(err.Error())
		}

		result = append(result, container)
	}

	return result, nil
}

func (con *inputFormConnection) GetInputFormRow(formID int, fields tables.InputForms, whreStr string) (tables.InputForms, error) {

	var data tables.InputForms

	err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").Select("if.id, if.user_id, u.name as user_name, to_char(if.created_at::timestamp, 'yyyy-mm-dd HH24:MI')::timestamp as created_at").Joins("join usr.users u on u.id = if.user_id").Where(fields).Where(whreStr).Order("id desc").First(&data).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *inputFormConnection) UpdateGetInputForm(id int, fields tables.Projects) (bool, error) {
	err := con.db.Scopes(SchemaFrm("projects")).Where("id = ?", id).Updates(fields).Error
	if err != nil {
		return false, err
	}

	return true, nil
}

func (con *inputFormConnection) DeleteGetInputForm(projectID int) (bool, error) {
	var data tables.Projects
	err := con.db.Scopes(SchemaFrm("projects")).Delete(&data, projectID).Error
	if err != nil {
		return false, err
	}

	return true, err
}

func (con *inputFormConnection) GetDates(strWhere string) ([]tables.Date, error) {

	var result []tables.Date

	where := ""
	if strWhere != "" {
		where = strWhere
	}

	err := con.db.Raw(`select to_char(d,'DD Mon YYYY') as date from generate_series(date_trunc('month', CURRENT_DATE),  date_trunc('month',CURRENT_DATE) + interval '1 month' - interval '1 day', '1 day'::interval) d ` + where).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDatesNew(strWhere string) ([]tables.Date, error) {

	var result []tables.Date

	where := ""
	if strWhere != "" {
		where = strWhere
	}

	err := con.db.Raw(`SELECT to_char(d, 'DD Mon YYYY') AS date FROM generate_series(date_trunc('month', CURRENT_DATE):: date, CURRENT_DATE, '1 day' :: interval) d ` + where + `order by date desc`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDatesWithFilter(month string, year string, strWhere string) ([]tables.Date, error) {

	var result []tables.Date

	where := ""
	if strWhere != "" {
		where = strWhere
	}

	err := con.db.Raw(`
	SELECT to_char(d, 'DD Mon YYYY') AS date 
	FROM generate_series(
		DATE_TRUNC('month', DATE('` + year + `' || '-' ||  '` + month + `' || '-01')), 
		(DATE_TRUNC('MONTH', DATE('` + year + `' || '-' ||  '` + month + `' || '-01')) + INTERVAL '1 MONTH' - INTERVAL '1 DAY'),
    '1 DAY') d  ` + where + `order by date desc`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetMonths(strWhere string) ([]tables.Months, error) {

	var result []tables.Months
	where := ""
	if strWhere != "" {
		where = strWhere
	}

	err := con.db.Raw(`SELECT TO_CHAR(months, 'mm')::int as id, TO_CHAR(months, 'Mon') AS month
						FROM generate_series(
							'2008-01-01' :: DATE,
							'2008-12-01' :: DATE ,
							'1 month'
						) AS months` + where).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetTotalDate(userID int, formID int, strWhere string) ([]tables.TotalDate, error) {

	uID := strconv.Itoa(userID)
	fID := strconv.Itoa(formID)

	whre := ""
	if strWhere != "" {
		whre = strWhere
	}
	var result []tables.TotalDate
	err := con.db.Raw(`select to_char(d, 'DD Mon YYYY') as date, count(f.id) as value
						from generate_series(date_trunc('month', CURRENT_DATE),  date_trunc('month',CURRENT_DATE) + interval '1 month' - interval '1 day', '1 day'::interval) d
						left join frm.input_forms_` + fID + ` f on d::date = f.created_at::date 
						join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + fID + `
						WHERE f.user_id = ` + uID + ` 
						` + whre + `
						group by 1 order by 1`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetTotalDateMonthly(userID int, formID int, strWhere string) ([]tables.TotalDate, error) {

	uID := strconv.Itoa(userID)
	fID := strconv.Itoa(formID)

	whre := ""
	if strWhere != "" {
		whre = strWhere
	}
	var result []tables.TotalDate
	err := con.db.Raw(`select to_char(d, 'DD Mon YYYY') as date, count(f.id) as value
						from generate_series(date_trunc('month', CURRENT_DATE),  date_trunc('month',CURRENT_DATE) + interval '1 month' - interval '1 day', '1 day'::interval) d
						left join frm.input_forms_` + fID + ` f on d::date = f.created_at::date AND f.user_id = ` + uID + ` 
						join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + fID + `

						` + whre + `
						group by to_char(d, 'DD Mon YYYY') order by 1`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetTotalMonth(userID int, formID int, strWhere string) ([]tables.TotalDate, error) {

	uID := strconv.Itoa(userID)
	fID := strconv.Itoa(formID)

	whre := ""
	if strWhere != "" {
		whre = strWhere
	}
	var result []tables.TotalDate
	err := con.db.Raw(`select to_char(d, 'Mon YYYY') as date, count(f.id) as value
						from generate_series(date_trunc('month', CURRENT_DATE),  date_trunc('month',CURRENT_DATE) + interval '1 month' - interval '1 day', '1 day'::interval) d
						left join frm.input_forms_` + fID + ` f on d::date = f.created_at::date
						join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + fID + `
						
						where f.user_id = ` + uID + ` 
						` + whre + `
						group by to_char(d, 'Mon YYYY')`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataHours(formID int, whreString string) ([]objects.GraficDataHours, error) {

	var result []objects.GraficDataHours
	err := con.db.Raw(`select tbl_hours.id, tbl_hours.field_hours, coalesce(f.total, 0) as value
							FROM
							(select 
							(row_number() OVER () - 1)::int as id, field_hours
								from (select to_char(generate_series(CURRENT_DATE::timestamp, to_char(CURRENT_DATE + 1 - INTERVAL '1 min', 'YYYY-MM-DD HH24:MI:SS')::timestamp, '1 hours'), 'HH24:MI') as field_hours) as f_date
							) tbl_hours

						left join (select to_char(f.created_at,'HH24')::int as hours, count(f.id) as total 
									from frm.input_forms_` + strconv.Itoa(formID) + ` f 
									join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `
									` + whreString + `
									group by to_char(f.created_at,'HH24')
							) as f on f.hours = tbl_hours.id`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataPeriodeDays(formID int, whreString string) ([]objects.GraficDataPeriod, error) {

	var result []objects.GraficDataPeriod
	err := con.db.Raw(`select (row_number() OVER ())::int as id, days.str_indo as field_period, t_data.total as value
						FROM (select to_char(tbl_days::date, 'Day') as str_days, tbl_days::date , to_char(tbl_days::date, 'dd')::int as day,
								(select d.name_id from mstr.days d where d.id = to_number(to_char(tbl_days::date, 'D'), '99G999D9S')::int) as str_indo
								FROM GENERATE_SERIES(DATE_TRUNC('week', CURRENT_DATE), DATE_TRUNC('week', CURRENT_DATE + INTERVAL '7 day') + INTERVAL '-1 day', '1 day'::interval) tbl_days) AS days
						
						left join (select to_char(f.created_at,'dd')::int as data_days, count(f.id) as total 
															from frm.input_forms_` + strconv.Itoa(formID) + ` f 
															join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `
															` + whreString + `
															group by to_char(f.created_at,'dd')
									) as t_data ON t_data.data_days = days.day`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataPeriodeMonthly(formID int, whreString string) ([]objects.GraficDataPeriod, error) {

	var result []objects.GraficDataPeriod
	err := con.db.Raw(`select tbl_days.id as id , tbl_days.days as field_period,  coalesce(f.total, 0) as value
							from (select to_char(d::date, 'DD')::int as id , to_char(d::date, 'DD') as days 
									from generate_series(date_trunc('month', CURRENT_DATE),  date_trunc('month',CURRENT_DATE) + interval '1 month' - interval '1 day', '1 day'::interval) as d) as tbl_days
						left join (select to_char(f.created_at,'dd')::int as id, to_char(f.created_at,'dd') as days, count(f.id) as total 
															from frm.input_forms_` + strconv.Itoa(formID) + ` f 
															join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `

															` + whreString + `
															group by to_char(f.created_at,'dd')
													) as f on f.id = tbl_days.id`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataPeriodeYearly(formID int, whreString string) ([]objects.GraficDataPeriod, error) {

	var result []objects.GraficDataPeriod
	err := con.db.Raw(`select tbl_month.id, tbl_month.str_month as field_period, tbl_data.total as value
						from (SELECT TO_CHAR(months, 'mm')::int as id, TO_CHAR(months, 'Mon') AS str_month
							FROM generate_series(
								'2008-01-01' :: DATE,
								'2008-12-01' :: DATE ,
								'1 month'
							) AS months) as tbl_month
						left join (select to_char(f.created_at,'mm')::int as id, to_char(f.created_at,'mm') as month, count(f.id) as total 
															from frm.input_forms_` + strconv.Itoa(formID) + ` f 
															join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `

														 	` + whreString + `
															group by to_char(f.created_at,'mm')
								) as tbl_data on tbl_data.id = tbl_month.id
															`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataHoursUserResp(formID int, whreString string) ([]objects.GraficDataHours, error) {

	var result []objects.GraficDataHours

	err := con.db.Raw(`select tbl_hours.id,
							tbl_hours.field_hours,
							f.value
						FROM
						(select (row_number() OVER () - 1)::int as id, field_hours from (select to_char(generate_series(CURRENT_DATE::timestamp, to_char(CURRENT_DATE + 1 - INTERVAL '1 min', 'YYYY-MM-DD HH24:MI:SS')::timestamp, '1 hours'), 'HH24:MI') as field_hours) as f_date) tbl_hours
						left join
						(select temp.hours, count(*) as value
						from (select f.user_id, to_char(f.created_at, 'HH24')::int as hours
							from frm.input_forms_` + strconv.Itoa(formID) + ` f
							join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `

							` + whreString + `
							group by 2, f.user_id) temp
						group by temp.hours) f on f.hours = tbl_hours.id`).Scan(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (con *inputFormConnection) GetActiveUserInputForm(formID int, fields tables.InputForms, whreStr string) ([]tables.JoinFormUsers, error) {
	var data []tables.JoinFormUsers

	err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").Select("if.user_id").
		Joins("join usr.users u on u.id = if.user_id").
		Joins("left join frm.form_users fu on if.user_id = fu.user_id AND fu.form_id =" + strconv.Itoa(formID)).
		Joins("left join frm.form_user_organizations fuo on fu.id = fuo.form_user_id").
		Joins("left join frm.input_form_organizations ifo on ifo.input_form_id=if.id AND ifo.form_id=" + strconv.Itoa(formID)).
		Where(fields).Where(whreStr).Group("if.user_id").Find(&data).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *inputFormConnection) GetActiveUserInputFormOld(formID int, fields tables.InputForms, whreStr string) ([]tables.JoinFormUsers, error) {
	var data []tables.JoinFormUsers

	err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").Select("if.user_id").
		Joins("join usr.users u on u.id = if.user_id").
		Joins("left join frm.form_users fu on if.user_id = fu.user_id AND fu.form_id =" + strconv.Itoa(formID)).
		// Joins("left join frm.form_user_organizations fuo on fu.id = fuo.form_user_id").
		// Joins("left join frm.input_form_organizations ifo on ifo.input_form_id=if.id AND ifo.form_id=" + strconv.Itoa(formID)).
		Where(fields).Where(whreStr).Group("if.user_id").Find(&data).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return data, err
	}
	return data, nil
}

func (con *inputFormConnection) GetDataPeriodeDaysResp(formID int, whreString string) ([]objects.GraficDataPeriod, error) {

	var result []objects.GraficDataPeriod
	err := con.db.Raw(`select (row_number() OVER ())::int as id, days.str_indo as field_period, data_user.total as value
						FROM (select to_char(tbl_days::date, 'Day') as str_days, tbl_days::date , to_char(tbl_days::date, 'dd')::int as day,
								(select d.name_id from mstr.days d where d.id = to_number(to_char(tbl_days::date, 'D'), '99G999D9S')::int) as str_indo

								FROM GENERATE_SERIES(DATE_TRUNC('week', CURRENT_DATE), DATE_TRUNC('week', CURRENT_DATE + INTERVAL '7 day') + INTERVAL '-1 day', '1 day'::interval) tbl_days) AS days
						
						left join (select temp.data_days, count(*) as total
										from (select f.user_id, to_char(f.created_at,'dd')::int as data_days
												from frm.input_forms_` + strconv.Itoa(formID) + ` f
												join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `

								 				` + whreString + `
												group by 2, f.user_id) temp
									group by temp.data_days) data_user on data_user.data_days = days.day`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataPeriodeMonthlyResp(formID int, whreString string) ([]objects.GraficDataPeriod, error) {

	var result []objects.GraficDataPeriod
	err := con.db.Raw(`select tbl_days.id as id , tbl_days.days as field_period,  coalesce(data_user.total, 0) as value
						from (select to_char(d::date, 'DD')::int as id , to_char(d::date, 'DD') as days 
								from generate_series(date_trunc('month', CURRENT_DATE),  date_trunc('month',CURRENT_DATE) + interval '1 month' - interval '1 day', '1 day'::interval) as d) as tbl_days
						
						left join (select temp.data_days, count(*) as total
									from (select f.user_id, to_char(f.created_at,'dd')::int as data_days
											from frm.input_forms_` + strconv.Itoa(formID) + ` f
											join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `

								 			` + whreString + `
											group by 2, f.user_id) temp
								group by temp.data_days) data_user ON data_user.data_days = tbl_days.id`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetDataPeriodeYearlyResp(formID int, whreString string) ([]objects.GraficDataPeriod, error) {

	var result []objects.GraficDataPeriod
	err := con.db.Raw(`select tbl_month.id, tbl_month.str_month as field_period, data_user.total as value
						from (SELECT TO_CHAR(months, 'mm')::int as id, TO_CHAR(months, 'Mon') AS str_month
							FROM generate_series(
								'2008-01-01' :: DATE,
								'2008-12-01' :: DATE ,
								'1 month'
							) AS months) as tbl_month
						left join (select temp.data_month, count(*) as total
										from (select f.user_id, to_char(f.created_at,'mm')::int as data_month
												from frm.input_forms_` + strconv.Itoa(formID) + ` f
												join frm.input_form_organizations ifo on ifo.input_form_id=f.id AND ifo.form_id=` + strconv.Itoa(formID) + `

												` + whreString + `
												group by 2, f.user_id) temp
									group by temp.data_month) data_user on data_user.data_month=  tbl_month.id`).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetReportFormResponden(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]objects.ReportResponden, error) {
	var result []objects.ReportResponden

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort + " ,"
	}
	fmt.Println("formID ::::", formID, orderBy, paging.SortBy)
	offsetLimit := ``
	if paging.Limit > 0 {
		offsetLimit = `offset ` + strconv.Itoa(offset) + ` limit ` + strconv.Itoa(paging.Limit)
	}
	err := con.db.Raw(`select fu.user_id, u.name as user_name, u.phone as user_phone, u.avatar
							, input_f.user_respons as submission
							, input_f.last_submission
							, input_f.last_submission_date
							, case when fr.submission_target_user > 0 THEN
								((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int
								else '0' end 
								as performance
							, case when fr.submission_target_user > 0 THEN
								to_char(float8 ((input_f.user_respons::float / fr.submission_target_user::float) * 100), 'FM999999999.0')
								else '0' end 
								as performance_float
							,input_f.address
							,input_f.latitude
							,input_f.longitude
							,case when fr.submission_target_user > 0 AND
									((input_f.user_respons::float / fr.submission_target_user::float) * 100) < 50 then '#F9B3B3'
									when input_f.user_respons = 0 then '#F9B3B3'
									when input_f.user_respons is null then '#F9B3B3'
									ELSE '#F9B3B3'
							end as submission_color
							,case when fr.submission_target_user > 0 AND
									((input_f.user_respons::float / fr.submission_target_user::float) * 100) < 50 then '#F9B3B3'
									when input_f.user_respons = 0 then '#F9B3B3'
									when input_f.user_respons is null then '#F9B3B3'
									ELSE '#F9B3B3'
							end as performance_color
							
						from frm.form_users fu
						left join (select f.user_id
										,count(f.id) as user_respons 
										,max(TO_CHAR(f.created_at, 'yyyy-mm-dd HH24:MI')) as last_submission_date
										,(select mstr.get_duration_time_ago(if_2.created_at) 
												from frm.input_forms_`+strconv.Itoa(formID)+` if_2 
												where if_2.user_id = f.user_id order by if_2.created_at desc limit 1) as last_submission
										,max(f.address) as address
										,max(ST_Y(f.geometry::geometry)) as latitude
										,max(ST_X(f.geometry::geometry)) as longitude
										
											from frm.input_forms_`+strconv.Itoa(formID)+` f 
											where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')
											group by f.user_id) input_f on input_f.user_id = fu.user_id
						left join usr.users u on u.id = fu.user_id
						left join frm.forms fr on fr.id = fu.form_id
						where fu.form_id = ?
						and fu.form_user_status_id = 1
						and fu.type = 'respondent'
						order by  `+orderBy+` user_name asc
						`+offsetLimit+``, formID).Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetReportFormRespondenUnionTeam(formID int, fields tables.InputForms, whreStr string, paging objects.Paging) ([]objects.ReportResponden, error) {
	var result []objects.ReportResponden

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort + " ,"
	}
	fmt.Println("formID ::::", formID, orderBy, paging.SortBy)
	offsetLimit := ``
	if paging.Limit > 0 {
		offsetLimit = `offset ` + strconv.Itoa(offset) + ` limit ` + strconv.Itoa(paging.Limit)
	}
	err := con.db.Raw(`select * from (
							select fu.user_id,
								u.name as user_name,
								u.phone as user_phone,
								u.avatar ,
								input_f.user_respons as submission ,
								input_f.last_submission ,
								input_f.last_submission_date ,
								fr.submission_target_user,
								case
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int <= 100 THEN ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int > 100 THEN 100
									else '0'
								end as performance ,
								case
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int <= 100 THEN to_char(float8 ((input_f.user_respons::float / fr.submission_target_user::float) * 100), 'FM999999999.0')
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int > 100.0 THEN to_char(float8 (1 * 100), 'FM999999999.0')
									else '0'
								end as performance_float ,
								input_f.address ,
								input_f.latitude ,
								input_f.longitude ,
								case
									when fr.submission_target_user > 0
											AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100) < 50 then '#F9B3B3'
									when input_f.user_respons = 0 then '#F9B3B3'
									when input_f.user_respons is null then '#F9B3B3'
									ELSE '#F9B3B3'
								end as submission_color ,
								case
									when fr.submission_target_user > 0
											AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100) < 50 then '#F9B3B3'
									when input_f.user_respons = 0 then '#F9B3B3'
									when input_f.user_respons is null then '#F9B3B3'
									ELSE '#F9B3B3'
								end as performance_color
							from frm.form_users fu
							left join
							(select f.user_id ,
									count(f.id) as user_respons ,
									max(TO_CHAR(f.created_at, 'yyyy-mm-dd HH24:MI')) as last_submission_date ,
							
								(select mstr.get_duration_time_ago(if_2.created_at)
								from frm.input_forms_`+strconv.Itoa(formID)+` if_2
								where if_2.user_id = f.user_id
								order by if_2.created_at desc
								limit 1) as last_submission ,

									max(f.address) as address ,
									max(ST_Y(f.geometry::geometry)) as latitude ,
									max(ST_X(f.geometry::geometry)) as longitude
							from frm.input_forms_`+strconv.Itoa(formID)+` f
							where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')
							group by f.user_id) input_f on input_f.user_id = fu.user_id
							left join usr.users u on u.id = fu.user_id
							left join frm.forms fr on fr.id = fu.form_id
							where fu.form_id = ?
							and fu.form_user_status_id = 1
							and fu.type = 'respondent'
							order by user_name asc
							
							) as tb_1
							
							UNION
							
							(select u.id as user_id,
								u.name as user_name,
								u.phone as user_phone,
								u.avatar ,
								input_f.user_respons as submission ,
								input_f.last_submission ,
								input_f.last_submission_date ,
								fr.submission_target_user,
								case
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int <= 100 THEN ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int > 100 THEN 100
									else '0'
								end as performance ,
								case
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int <= 100 THEN to_char(float8 ((input_f.user_respons::float / fr.submission_target_user::float) * 100), 'FM999999999.0')
									when fr.submission_target_user > 0 AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100)::int > 100.0 THEN to_char(float8 (1 * 100), 'FM999999999.0')
									else '0'
								end as performance_float ,
								input_f.address ,
								input_f.latitude ,
								input_f.longitude ,
								case
									when fr.submission_target_user > 0
											AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100) < 50 then '#F9B3B3'
									when input_f.user_respons = 0 then '#F9B3B3'
									when input_f.user_respons is null then '#F9B3B3'
									ELSE '#F9B3B3'
								end as submission_color ,
								case
									when fr.submission_target_user > 0
											AND ((input_f.user_respons::float / fr.submission_target_user::float) * 100) < 50 then '#F9B3B3'
									when input_f.user_respons = 0 then '#F9B3B3'
									when input_f.user_respons is null then '#F9B3B3'
									ELSE '#F9B3B3'
								end as performance_color
							from frm.form_teams ft
							join usr.team_users tu on tu.team_id= ft.team_id
							join usr.users u on u.id=tu.user_id
							left join frm.forms fr on fr.id = ft.form_id
							left join
							(select f.user_id ,
									count(f.id) as user_respons ,
									max(TO_CHAR(f.created_at, 'yyyy-mm-dd HH24:MI')) as last_submission_date ,
							
								(select mstr.get_duration_time_ago(if_2.created_at)
								from frm.input_forms_`+strconv.Itoa(formID)+` if_2
								where if_2.user_id = f.user_id
								order by if_2.created_at desc
								limit 1) as last_submission ,
									max(f.address) as address ,
									max(ST_Y(f.geometry::geometry)) as latitude ,
									max(ST_X(f.geometry::geometry)) as longitude
							from frm.input_forms_`+strconv.Itoa(formID)+` f
							where TO_CHAR(f.created_at::date, 'yyyy-mm-dd') = TO_CHAR(NOW()::date, 'yyyy-mm-dd')
							group by f.user_id) input_f on input_f.user_id = u.id
							where ft.form_id=?)
							
							order by  `+orderBy+` user_name asc
						`+offsetLimit+``, formID, formID).Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (con *inputFormConnection) GetInputFormCustomAnswerRow(formID int, formFieldID int, inputFormID int) (tables.InputFormCustomAnswers, error) {

	var data tables.InputFormCustomAnswers

	err := con.db.Table("frm.input_form_custom_answers").Where("form_id", formID).Where("form_field_id", formFieldID).Where("input_form_id", inputFormID).First(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return tables.InputFormCustomAnswers{}, err
	}

	if err != nil {
		return tables.InputFormCustomAnswers{}, err
	}

	return data, nil
}

func (con *inputFormConnection) GetInputDataOrganizationRows(formID int, selectFieldStr string, fields tables.InputForms, whreStr string) ([][]string, error) {

	// var result []tables.InputForms

	// err := con.db.Scopes(SchemaFrm("input_forms_" + strconv.Itoa(formID))).Where(fields).Find(&data).Error

	rows, err := con.db.Table("frm.input_forms_" + strconv.Itoa(formID) + " as if").Select(selectFieldStr).Joins("left join frm.input_form_organizations ifo ON if.id = ifo.input_form_id AND ifo.form_id=" + strconv.Itoa(formID)).Where(fields).Where(whreStr).Rows()
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()

	length := len(cols)

	//----
	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]

		}

		err = rows.Scan(pointers...)
		if err != nil {
			fmt.Println("err pointers-----------", err)
			panic(err.Error())
		}

		result = append(result, container)
	}

	return result, nil
}

func (con *inputFormConnection) GetInputFormOrganizationRows(formID int, fields tables.InputFormJoinOrganizations, whreStr string, paging objects.Paging) ([]tables.InputFormJoinOrganizations, error) {

	offset := (paging.Page - 1) * paging.Limit

	orderBy := ""
	if paging.SortBy != "" {
		orderBy = paging.SortBy + " " + paging.Sort
	}

	var data []tables.InputFormJoinOrganizations

	err := con.db.Table("frm.input_forms_"+strconv.Itoa(formID)+" as if").Select("if.id", "if.user_id", "u.name as user_name", "u.phone", "u.avatar", "if.address", "ST_Y(if.geometry::geometry) as latitude", "ST_X(if.geometry::geometry) as longitude", "if.created_at", "if.updated_at", "o.name as organization_name", "if.updated_count").
		Joins("join usr.users u on u.id = if.user_id").
		Joins("left join frm.input_form_organizations ifo ON ifo.input_form_id=if.id AND ifo.form_id=" + strconv.Itoa(formID)).
		Joins("left join mstr.organizations o ON o.id=ifo.organization_id").
		Where(fields).Where(whreStr).Order(orderBy).Order("id desc").Limit(paging.Limit).Offset(offset).Find(&data).Error

	if err != nil {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) GetOrganizationInputForm(fields tables.InputFormOrganizations, whreStr string) ([]objects.InputFormOrganizations, error) {
	var data []objects.InputFormOrganizations

	err := con.db.Table("frm.input_form_organizations").Select("input_form_organizations.organization_id, o.name as organization_name").Joins("left join mstr.organizations o ON o.id=input_form_organizations.organization_id").Where(fields).Where(whreStr).Group("1,2").Find(&data).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return data, nil
}

func (con *inputFormConnection) InsertInputFormOrgData(data tables.InputFormOrganizations) (tables.InputFormOrganizations, error) {
	err := con.db.Scopes(SchemaFrm("input_form_organizations")).Create(&data).Error
	if err != nil {
		return tables.InputFormOrganizations{}, nil
	}
	return data, err
}
