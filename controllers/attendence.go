package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"snapin-form/helpers"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"

	"github.com/360EntSecGroup-Skylar/excelize"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"

	"github.com/gin-gonic/gin"
)

// interface
type AttendanceController interface {
	InsertAttendance(c *gin.Context)
	InsertWarningAttendance(c *gin.Context)
	InsertWarning2Attendance(c *gin.Context)
	FormAttendanceList(c *gin.Context)
	FormAttendanceListExport(c *gin.Context)
	FormAttendanceListExportCSV(c *gin.Context)
	FormAttendanceListExportPost(c *gin.Context)
	// InsertAttendanceOut(c *gin.Context)
	FormAttendanceMapsList(c *gin.Context)
	FormAttendanceMapsList2(c *gin.Context) // 1 user 2 latlong
	AutoInsertData(c *gin.Context)
	GetListLocationAttendance(c *gin.Context)
	GetListLocationAttendanceDash(c *gin.Context)
}

type attController struct {
	formMod models.FormModels
	userMod models.UserModels
	attMod  models.AttendanceModels
	helper  helpers.Helper
	shrtMod models.ShortenUrlModels
	pgErr   *pgconn.PgError
}

func NewAttController(formModel models.FormModels, userModel models.UserModels, attModel models.AttendanceModels, help helpers.Helper, shrtModel models.ShortenUrlModels) AttendanceController {
	return &attController{
		formMod: formModel,
		userMod: userModel,
		attMod:  attModel,
		helper:  help,
		shrtMod: shrtModel,
	}
}

func (ctr *attController) InsertAttendance(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	respondenID, _ := strconv.Atoi(claims["id"].(string))

	fmt.Println("my UserID :::", respondenID)
	// os.Exit(0)

	var reqData objects.Attendence
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	//set timezone,
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	//check
	currentDate := time.Now().Format("2006-01-02")

	// shorten url
	var shortenUrlPic objects.ShortURLResponse
	shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":   err,
			"status":  false,
			"message": "Shorten URL is failed",
		})
		return
	}

	// shortenUrlPic.Data.ShortURL = reqData.FacePic //disabled url shorten
	// shortenUrlPic.Data.ShortURL = shortenUrlPic.Data.ShortURL

	var whre1 tables.Attendances
	whre1.UserID = respondenID
	whre1.FormID = reqData.FormID
	strField := "TO_CHAR(attendances.attendance_in::date, 'yyyy-mm-dd') = '" + currentDate + "'"

	getToday, err := ctr.attMod.GetAttendanceRow(whre1, strField)
	if getToday.ID > 0 {
		if getToday.AttendanceOut.IsZero() == false {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Your attendance has check-out before",
				"data":    nil,
			})
			return
		}

		if reqData.OfflineTime != "" {
			now, _ = time.Parse("2006-01-02 15:04", reqData.OfflineTime)
		}

		var postData tables.Attendances
		postData.AttendanceOut = now
		postData.FacePicOut = shortenUrlPic.Data.ShortURL
		postData.AddressOut = reqData.Address

		var geo tables.Geometry
		geo.Latitude = reqData.Latitude
		geo.Longitude = reqData.Longitude
		udpateAtt, err := ctr.attMod.UpdateAttendance(getToday.ID, postData, geo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if udpateAtt {
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Successfuly check-out attendance",
				"data":    nil,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed check-out attendance",
				"data":    nil,
			})
			return
		}
	}

	// checkout overdate if available -----------------------------------

	var whrFrm tables.Forms
	whrFrm.ID = reqData.FormID
	getFormData, _ := ctr.formMod.GetFormRow(whrFrm)

	// cek chekout yg kosong sebelum hari ini (kurang dari sama dengan today) dari batas tanggal overdate At
	if getFormData.IsAttendanceRequired == true && getFormData.AttendanceOverdateAt.IsZero() == false {
		checkCheckoutBefore, err := ctr.attMod.GetLastAttendanceOverdate(getFormData.ID, respondenID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"error":   err,
				"message": "Data is available",
			})
			return
		}

		if checkCheckoutBefore.ID >= 1 && checkCheckoutBefore.AttendanceOut == "" {
			// jika absen sebelumnya belum checkout maka button akan terus checkout
			if reqData.OfflineTime != "" {
				now, _ = time.Parse("2006-01-02 15:04", reqData.OfflineTime)
			}

			var postData tables.Attendances
			postData.AttendanceOut = now
			postData.FacePicOut = shortenUrlPic.Data.ShortURL
			postData.AddressOut = reqData.Address

			var geo tables.Geometry
			geo.Latitude = reqData.Latitude
			geo.Longitude = reqData.Longitude
			udpateAtt, err := ctr.attMod.UpdateAttendance(checkCheckoutBefore.ID, postData, geo)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			if udpateAtt {
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Successfuly check-out attendance overdate",
					"data":    nil,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed check-out attendance overdate",
					"data":    nil,
				})
				return
			}
		}
	}

	// end overdate here -----------------------------------------------

	var postData tables.Attendances
	postData.FormID = reqData.FormID
	postData.UserID = respondenID
	postData.FacePicIn = shortenUrlPic.Data.ShortURL
	postData.AttendanceIn = now
	postData.AddressIn = reqData.Address

	var geo tables.Geometry
	geo.Latitude = reqData.Latitude
	geo.Longitude = reqData.Longitude
	save, err := ctr.attMod.InsertAttendance(postData, geo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if save.ID > 0 {
		getCompany, _ := ctr.attMod.GetCompanyID(reqData.FormID, respondenID)
		fmt.Println(getCompany)
		// os.Exit(0)
		var value objects.InsAttOrg
		value.AttendanceID = save.ID
		value.OrganizationID = getCompany.OrganizationID

		_, err := ctr.attMod.InsertAttendanceOrganization(value)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Successfuly check-in attendance",
			"data":    nil,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed check-in attendance",
			"data":    nil,
		})
		return
	}

}

func (ctr *attController) InsertWarningAttendance(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	respondenID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.WarnAttendence
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	//set timezone,
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	currentDate := time.Now().Format("2006-01-02")
	fmt.Println("postData :::", currentDate)

	if reqData.Type == "check_in" {
		var shortenUrlPic objects.ShortURLResponse
		// shorten url

		shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   err,
				"status":  false,
				"message": "Shorten URL is failed",
			})
			return
		}

		// shortenUrlPic.Data.ShortURL = reqData.FacePic //disabled url shorten
		// shortenUrlPic.Data.ShortURL = shortenUrlPic.Data.ShortURL

		//checking today
		var whre1 tables.Attendances
		whre1.UserID = respondenID
		whre1.FormID = reqData.FormID
		strField := "TO_CHAR(attendances.attendance_in::date, 'yyyy-mm-dd') = '" + currentDate + "'"

		getToday, _ := ctr.attMod.GetAttendanceRow(whre1, strField)
		if getToday.ID >= 1 {
			if getToday.AttendanceIn.IsZero() == false {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": "Your attendance has check-in before",
					"data":    nil,
				})
				return
			}
		}

		var postData tables.Attendances
		postData.FormID = reqData.FormID
		postData.UserID = respondenID
		postData.FacePicIn = shortenUrlPic.Data.ShortURL
		postData.AttendanceIn = now
		postData.AddressIn = reqData.Address
		// postData.FormAttendanceLocationIdIn = reqData.LocationID

		var geo tables.Geometry
		geo.Latitude = reqData.Latitude
		geo.Longitude = reqData.Longitude
		save, err := ctr.attMod.InsertAttendance(postData, geo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if save.ID > 0 {
			getCompany, _ := ctr.attMod.GetCompanyID(reqData.FormID, respondenID)

			if getCompany.OrganizationID >= 1 {
				var value objects.InsAttOrg
				value.AttendanceID = save.ID
				value.OrganizationID = getCompany.OrganizationID

				_, err := ctr.attMod.InsertAttendanceOrganization(value)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Successfuly check-in attendance",
					"data":    nil,
				})
				return
			} else if getCompany.OrganizationID <= 0 {
				getTeam, _ := ctr.attMod.GetTeamByRespondent(respondenID)

				if len(getTeam) > 0 {
					for i := 0; i < len(getTeam); i++ {
						getFormTeam, _ := ctr.attMod.GetFormTeam(reqData.FormID, getTeam[i].TeamID)

						if getFormTeam.OrganizationID >= 1 {
							var value objects.InsAttOrg
							value.AttendanceID = save.ID
							value.OrganizationID = getFormTeam.OrganizationID

							_, err := ctr.attMod.InsertAttendanceOrganization(value)
							if err != nil {
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							c.JSON(http.StatusOK, gin.H{
								"status":  true,
								"message": "Successfuly check-in attendance",
								"data":    nil,
							})
							return
						}
					}
				}

			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": "User responden belum masuk ke form tersebut",
					"data":    nil,
				})
				return
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed check-in attendance",
				"data":    nil,
			})
			return
		}
	} else if reqData.Type == "check_out" {
		var shortenUrlPic objects.ShortURLResponse
		// shorten url

		shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   err,
				"status":  false,
				"message": "Shorten URL is failed",
			})
			return
		}

		//check
		var whre1 tables.Attendances
		whre1.UserID = respondenID
		whre1.FormID = reqData.FormID
		strField := "TO_CHAR(attendances.attendance_in::date, 'yyyy-mm-dd') = '" + currentDate + "'"

		getToday, _ := ctr.attMod.GetAttendanceRow(whre1, strField)
		if getToday.ID >= 1 {
			if getToday.AttendanceOut.IsZero() == false {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": "Your attendance has check-out before",
					"data":    nil,
				})
				return
			}

			if reqData.OfflineTime != "" {
				now, _ = time.Parse("2006-01-02 15:04", reqData.OfflineTime)
			}

			var postData tables.Attendances
			postData.AttendanceOut = now
			postData.FacePicOut = shortenUrlPic.Data.ShortURL
			postData.AddressOut = reqData.Address
			// postData.FormAttendanceLocationIdOut = reqData.LocationID // disabled lebih dulu untuk next update

			var geo tables.Geometry
			geo.Latitude = reqData.Latitude
			geo.Longitude = reqData.Longitude
			udpateAtt, err := ctr.attMod.UpdateAttendance(getToday.ID, postData, geo)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			if udpateAtt {
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Successfuly check-out attendance",
					"data":    nil,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed check-out attendance",
					"data":    nil,
				})
				return
			}
		}

		// checkout overdate if available -----------------------------------

		var whrFrm tables.Forms
		whrFrm.ID = reqData.FormID
		getFormData, _ := ctr.formMod.GetFormRow(whrFrm)

		// cek chekout yg kosong sebelum hari ini (kurang dari sama dengan today) dari batas tanggal overdate At
		if getFormData.IsAttendanceRequired == true && getFormData.AttendanceOverdateAt.IsZero() == false {
			checkCheckoutBefore, err := ctr.attMod.GetLastAttendanceOverdate(getFormData.ID, respondenID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"error":   err,
					"message": "Data is available",
				})
				return
			}

			if checkCheckoutBefore.ID >= 1 && checkCheckoutBefore.AttendanceOut == "" {
				// jika absen sebelumnya belum checkout maka button akan terus checkout
				if reqData.OfflineTime != "" {
					now, _ = time.Parse("2006-01-02 15:04", reqData.OfflineTime)
				}
				var shortenUrlPic objects.ShortURLResponse
				// shorten url

				shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{
						"error":   err,
						"status":  false,
						"message": "Shorten URL is failed",
					})
					return
				}

				// shortenUrlPic.Data.ShortURL = reqData.FacePic //disabled url shorten
				shortenUrlPic.Data.ShortURL = shortenUrlPic.Data.ShortURL
				var postData tables.Attendances
				postData.AttendanceOut = now
				postData.FacePicOut = shortenUrlPic.Data.ShortURL
				postData.AddressOut = reqData.Address

				var geo tables.Geometry
				geo.Latitude = reqData.Latitude
				geo.Longitude = reqData.Longitude
				udpateAtt, err := ctr.attMod.UpdateAttendance(checkCheckoutBefore.ID, postData, geo)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if udpateAtt {
					c.JSON(http.StatusOK, gin.H{
						"status":  true,
						"message": "Successfuly check-out attendance overdate",
						"data":    nil,
					})
					return
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": "Failed check-out attendance overdate",
						"data":    nil,
					})
					return
				}
			}
		} else if getFormData.IsAttendanceRequired == true && getFormData.AttendanceOverdateAt.IsZero() == true {
			// disini next ditambahkan pengecekan data pending
			// if check pending time stamp dari Apps
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Berhasil absen dari data pending",
			})
			return
		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Anda sudah absen masuk hari ini",
			})
			return
		}

		// end overdate here -----------------------------------------------

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Attendance type is wrong",
		})
		return
	}

}

func (ctr *attController) InsertWarning2Attendance(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	respondenID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.WarnAttendence
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	// set created at from mobile apps
	dateTime := reqData.CreatedAt
	if reqData.OfflineCreatedAt != "" {
		dateTime = reqData.OfflineCreatedAt
	}
	currentDate := dateTime[0:10]

	loc, _ := time.LoadLocation("Asia/Jakarta")
	const shortForm = "2006-01-02 15:04:05"

	if reqData.Type == "check_in" {
		var shortenUrlPic objects.ShortURLResponse
		// shorten url

		shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   err,
				"status":  false,
				"message": "Shorten URL is failed",
			})
			return
		}

		//checking today
		var whre1 tables.Attendances
		whre1.UserID = respondenID
		whre1.FormID = reqData.FormID
		strField := "TO_CHAR(attendances.attendance_in::date, 'yyyy-mm-dd') = '" + currentDate + "'"

		getToday, _ := ctr.attMod.GetAttendanceRow(whre1, strField)
		if getToday.ID >= 1 {
			if getToday.AttendanceIn.IsZero() == false {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": "Your attendance has check-in before",
					"data":    nil,
				})
				return
			}
		}

		// Parse string tanggal dan waktu menjadi time.Time
		// parsedAttendanceIn, err := time.Parse("2006-01-02 15:04:05", reqData.CreatedAt)
		parsedAttendanceIn, err := time.ParseInLocation(shortForm, reqData.CreatedAt, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Datetime is empty",
			})
			return
		}
		if reqData.OfflineCreatedAt != "" {
			// parsedAttendanceIn, err = time.Parse("2006-01-02 15:04:05", reqData.OfflineCreatedAt)
			parsedAttendanceIn, err = time.ParseInLocation(shortForm, reqData.OfflineCreatedAt, loc)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Datetime is empty",
				})
				return
			}
		}

		var postData tables.Attendances
		postData.FormID = reqData.FormID
		postData.UserID = respondenID
		postData.FacePicIn = shortenUrlPic.Data.ShortURL
		postData.AttendanceIn = parsedAttendanceIn
		postData.AddressIn = reqData.Address
		// postData.FormAttendanceLocationIdIn = reqData.LocationID

		var geo tables.Geometry
		geo.Latitude = reqData.Latitude
		geo.Longitude = reqData.Longitude
		save, err := ctr.attMod.InsertAttendance(postData, geo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if save.ID > 0 {
			getCompany, _ := ctr.attMod.GetCompanyID(reqData.FormID, respondenID)

			if getCompany.OrganizationID >= 1 {
				var value objects.InsAttOrg
				value.AttendanceID = save.ID
				value.OrganizationID = getCompany.OrganizationID

				_, err := ctr.attMod.InsertAttendanceOrganization(value)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Successfuly check-in attendance",
					"data":    nil,
				})
				return
			} else if getCompany.OrganizationID <= 0 {
				getTeam, _ := ctr.attMod.GetTeamByRespondent(respondenID)

				if len(getTeam) > 0 {
					for i := 0; i < len(getTeam); i++ {
						getFormTeam, _ := ctr.attMod.GetFormTeam(reqData.FormID, getTeam[i].TeamID)

						if getFormTeam.OrganizationID >= 1 {
							var value objects.InsAttOrg
							value.AttendanceID = save.ID
							value.OrganizationID = getFormTeam.OrganizationID

							_, err := ctr.attMod.InsertAttendanceOrganization(value)
							if err != nil {
								c.JSON(http.StatusBadRequest, gin.H{
									"error": err,
								})
								return
							}

							c.JSON(http.StatusOK, gin.H{
								"status":  true,
								"message": "Successfuly check-in attendance",
								"data":    nil,
							})
							return
						}
					}
				}

			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": "User responden belum masuk ke form tersebut",
					"data":    nil,
				})
				return
			}

		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Failed check-in attendance",
				"data":    nil,
			})
			return
		}
	} else if reqData.Type == "check_out" {
		var shortenUrlPic objects.ShortURLResponse
		// shorten url

		shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{
				"error":   err,
				"status":  false,
				"message": "Shorten URL is failed",
			})
			return
		}

		// Parse string tanggal dan waktu menjadi time.Time
		// parsedAttendanceOut, err := time.Parse("2006-01-02 15:04:05", reqData.CreatedAt)
		parsedAttendanceOut, err := time.ParseInLocation(shortForm, reqData.CreatedAt, loc)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Datetime is empty",
			})
			return
		}
		if reqData.OfflineCreatedAt != "" {
			// parsedAttendanceOut, err = time.Parse("2006-01-02 15:04:05", reqData.OfflineCreatedAt)
			parsedAttendanceOut, err = time.ParseInLocation(shortForm, reqData.OfflineCreatedAt, loc)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Datetime is empty",
				})
				return
			}
		}

		//check
		var whre1 tables.Attendances
		whre1.UserID = respondenID
		whre1.FormID = reqData.FormID
		strField := "TO_CHAR(attendances.attendance_in::date, 'yyyy-mm-dd') = '" + currentDate + "'"

		getToday, _ := ctr.attMod.GetAttendanceRow(whre1, strField)
		if getToday.ID >= 1 {
			if getToday.AttendanceOut.IsZero() == false {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": "Your attendance datetime check-out is empty",
					"data":    nil,
				})
				return
			}

			var postData tables.Attendances
			postData.AttendanceOut = parsedAttendanceOut
			postData.FacePicOut = shortenUrlPic.Data.ShortURL
			postData.AddressOut = reqData.Address
			// postData.FormAttendanceLocationIdOut = reqData.LocationID // disabled lebih dulu untuk next update

			var geo tables.Geometry
			geo.Latitude = reqData.Latitude
			geo.Longitude = reqData.Longitude
			udpateAtt, err := ctr.attMod.UpdateAttendance(getToday.ID, postData, geo)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}

			if udpateAtt {
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"message": "Successfuly check-out attendance",
					"data":    nil,
				})
				return
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "Failed check-out attendance",
					"data":    nil,
				})
				return
			}
		}

		// checkout overdate if available -----------------------------------

		var whrFrm tables.Forms
		whrFrm.ID = reqData.FormID
		getFormData, _ := ctr.formMod.GetFormRow(whrFrm)

		// cek chekout yg kosong sebelum hari ini (kurang dari sama dengan today) dari batas tanggal overdate At
		if getFormData.IsAttendanceRequired == true && getFormData.AttendanceOverdateAt.IsZero() == false {
			checkCheckoutBefore, err := ctr.attMod.GetLastAttendanceOverdate(getFormData.ID, respondenID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  false,
					"error":   err,
					"message": "Data is available",
				})
				return
			}

			if checkCheckoutBefore.ID >= 1 && checkCheckoutBefore.AttendanceOut == "" {
				// jika absen sebelumnya belum checkout maka button akan terus checkout
				// if reqData.OfflineTime != "" {
				// 	now, _ = time.Parse("2006-01-02 15:04", reqData.OfflineTime)
				// }
				var shortenUrlPic objects.ShortURLResponse
				// shorten url

				shortenUrlPic, err = ctr.shrtMod.ShortenImageURL(reqData.FacePic)
				if err != nil {
					c.JSON(http.StatusBadGateway, gin.H{
						"error":   err,
						"status":  false,
						"message": "Shorten URL is failed",
					})
					return
				}

				var postData tables.Attendances
				postData.AttendanceOut = parsedAttendanceOut
				postData.FacePicOut = shortenUrlPic.Data.ShortURL
				postData.AddressOut = reqData.Address

				var geo tables.Geometry
				geo.Latitude = reqData.Latitude
				geo.Longitude = reqData.Longitude
				udpateAtt, err := ctr.attMod.UpdateAttendance(checkCheckoutBefore.ID, postData, geo)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": err,
					})
					return
				}

				if udpateAtt {
					c.JSON(http.StatusOK, gin.H{
						"status":  true,
						"message": "Successfuly check-out attendance overdate",
						"data":    nil,
					})
					return
				} else {
					c.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": "Failed check-out attendance overdate",
						"data":    nil,
					})
					return
				}
			}

		} else {

			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Anda sudah absen masuk hari ini",
			})
			return
		}

		// end overdate here -----------------------------------------------

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Attendance type is wrong",
		})
		return
	}

}

func (ctr *attController) FormAttendanceList(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	organization_id, _ := strconv.Atoi(claims["organization_id"].(string))

	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	filterDate := c.Request.URL.Query().Get("date")
	// startDate := c.Request.URL.Query().Get("start_date")
	// endDate := c.Request.URL.Query().Get("end_date")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")
	company_id, _ := strconv.Atoi(c.Request.URL.Query().Get("company_id"))

	if formID > 0 {

		var whre1 tables.Attendances
		whre1.FormID = formID

		whrStr := ""
		whrTime := ""
		whrComp := ""
		if filterDate == "" {
			whrTime = " AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char(now()::date,'yyyy-mm-dd')"
		}

		if filterDate != "" && searchKeyWord == "" {
			whrTime = " AND to_char(attendances.created_at::date,  'yyyy-mm-dd') = to_char('" + filterDate + "'::date,'yyyy-mm-dd')"
		}

		if searchKeyWord != "" && filterDate == "" {
			whrStr = " AND u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%'"
		}

		if searchKeyWord != "" && filterDate != "" {
			whrStr = " AND (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char('" + filterDate + "'::date, 'yyyy-mm-dd')"
		}

		// if company_id != 0 {
		getCompanyID, _ := ctr.attMod.GetCompanyIDByFormID(formID)
		if organization_id != getCompanyID.OrganizationID && company_id <= 0 {
			// company ID by TOKEN
			whrComp = " AND o.id = " + claims["organization_id"].(string)
		} else if organization_id == getCompanyID.OrganizationID && company_id >= 1 {
			// select company option (form sharing only)
			whrComp = " AND o.id = " + strconv.Itoa(company_id)
		}
		// whrComp = " AND o.id =  " + strconv.Itoa(company_id) + ""
		// }

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		results, err := ctr.attMod.GetUserTeamAttendanceReports(whre1, whrStr, whrTime, whrComp, paging)
		// results, err := ctr.attMod.GetUserAttendanceReports(whre1, whrStr, whrTime, paging)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		fmt.Println(whrComp)
		resultAll, err := ctr.attMod.GetUserTeamAttendanceReports(whre1, whrStr, whrTime, whrComp, objects.Paging{})
		// resultAll, err := ctr.attMod.GetUserAttendanceReports(whre1, whrStr, whrTime, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			var res []objects.AttendenceReport
			for i := 0; i < len(results); i++ {

				var each objects.AttendenceReport
				each.ID = results[i].ID
				each.FormID = results[i].FormID
				each.UserID = results[i].UserID
				each.UserName = results[i].UserName
				each.UserPhone = results[i].UserPhone
				each.OrganizationID = results[i].OrganizationID
				each.OrganizationName = results[i].OrganizationName
				each.AttendanceIn = results[i].AttendanceIn
				each.AttendanceOut = results[i].AttendanceOut
				each.AddressIn = results[i].AddressIn
				each.AddressOut = results[i].AddressOut
				each.Duration = results[i].Duration
				each.LocationIn = results[i].LocationIn
				each.LocationOut = results[i].LocationOut
				each.CreatedAt = results[i].CreatedAt

				var facePicIn = results[i].FacePicIn
				var facePicOut = results[i].FacePicOut

				if facePicIn != "" && len(facePicIn) < 35 {
					words := strings.Split(facePicIn, "/")
					lastWord := words[len(words)-1]

					linkReal, err := ctr.shrtMod.GetLinkReal(lastWord)
					if err != nil {
						// Handle the error.
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					each.FacePicIn = linkReal.Data.URL
				}

				if facePicOut != "" && len(facePicOut) < 35 {
					words := strings.Split(facePicOut, "/")
					lastWord := words[len(words)-1]

					linkReal, err := ctr.shrtMod.GetLinkReal(lastWord)
					if err != nil {
						// Handle the error.
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}

					each.FacePicOut = linkReal.Data.URL
				}

				res = append(res, each)
			}

			totalPage := 0
			if limit > 0 {
				totalPage = len(resultAll) / limit
				if (len(resultAll) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var paging objects.DataRows
			paging.TotalRows = len(resultAll)
			paging.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *attController) FormAttendanceListExport(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err := strconv.Atoi(c.Param("formid"))

	searchKeyWord := c.Request.URL.Query().Get("search")
	// filterDate := c.Request.URL.Query().Get("date")
	startDate := c.Request.URL.Query().Get("start_date") + " 00:00"
	endDate := c.Request.URL.Query().Get("end_date") + " 23:59"

	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var whre1 tables.Attendances
		whre1.FormID = formID

		whrStr := ""
		// if searchKeyWord != "" && filterDate == "" {
		// 	whrStr = " u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%'"
		// }

		// if filterDate != "" && searchKeyWord == "" {
		// 	whrStr = " to_char(attendances.created_at::date,  'yyyy-mm-dd') = to_char('" + filterDate + "'::date,'yyyy-mm-dd')"
		// }

		// if searchKeyWord != "" && filterDate != "" {
		// 	whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char('" + filterDate + "'::date, 'yyyy-mm-dd')"
		// }

		if startDate != "" && searchKeyWord == "" {
			whrStr = " attendances.created_at >= '" + startDate + "'"
		}
		if endDate != "" && searchKeyWord == "" {
			whrStr = " attendances.created_at <= '" + endDate + "'"
		}
		if startDate != "" && endDate != "" && searchKeyWord == "" {
			whrStr = " attendances.created_at BETWEEN '" + startDate + "' AND '" + endDate + "'"
		}
		//
		if searchKeyWord != "" && startDate == "" && endDate == "" {
			whrStr = " u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%'"
		}
		if searchKeyWord != "" && startDate != "" && endDate == "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND attendances.created_at >= '" + startDate + "'"
		}
		if searchKeyWord != "" && startDate == "" && endDate != "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND attendances.created_at <= '" + startDate + "'"
		}

		if searchKeyWord != "" && startDate != "" && endDate != "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND attendances.created_at BETWEEN '" + startDate + "' and '" + endDate + "'"
		}

		results, err := ctr.attMod.GetAttendanceReports(whre1, whrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {
			xlsx := excelize.NewFile()
			sheetName := "Export-Absen-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheetName)

			style, err := xlsx.NewStyle(`{"fill":{"type":"pattern","color":["#e4e4e4"],"pattern":1}}`)
			if err != nil {
				fmt.Println(err)
				c.JSON(http.StatusBadGateway, gin.H{
					"status":  false,
					"message": err.Error(),
				})
				return
			}
			xlsx.SetCellStyle(sheetName, "A1", "M1", style)
			xlsx.SetColWidth(sheetName, "A", "K", 20)

			xlsx.SetCellValue(sheetName, "A1", "Nama")
			xlsx.SetCellValue(sheetName, "B1", "Kontak")
			xlsx.SetCellValue(sheetName, "C1", "Selfie Absen Masuk")
			xlsx.SetCellValue(sheetName, "D1", "Tanggal Sistem Absen Masuk (WIB)")
			xlsx.SetCellValue(sheetName, "E1", "Tanggal Absen Masuk")
			xlsx.SetCellValue(sheetName, "F1", "Jam Absen Masuk")
			xlsx.SetCellValue(sheetName, "G1", "Lokasi Absen Masuk")

			xlsx.SetCellValue(sheetName, "H1", "Selfie Absen Keluar")
			xlsx.SetCellValue(sheetName, "I1", "Tanggal Sistem Absen Keluar (WIB)")
			xlsx.SetCellValue(sheetName, "J1", "Tanggal Absen Keluar")
			xlsx.SetCellValue(sheetName, "K1", "Jam Absen Keluar")
			xlsx.SetCellValue(sheetName, "L1", "Lokasi Absen Keluar")
			xlsx.SetCellValue(sheetName, "M1", "Durasi Kerja")
			// xlsx.SetCellValue(sheetName, "L1", "Lokasi Masuk")
			// xlsx.SetCellValue(sheetName, "M1", "Lokasi Keluar")

			row := 2
			for i := 0; i < len(results); i++ {

				updatedAt := ""
				if results[i].AttendanceOut != "" {
					updatedAt = results[i].UpdatedAt
				}

				xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(row), results[i].UserName)
				xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(row), results[i].UserPhone)
				xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(row), results[i].FacePicIn)
				xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(row), results[i].CreatedAt)
				xlsx.SetCellValue(sheetName, "E"+strconv.Itoa(row), results[i].AttendanceDateIn)
				xlsx.SetCellValue(sheetName, "F"+strconv.Itoa(row), results[i].AttendanceTimeIn)
				xlsx.SetCellValue(sheetName, "G"+strconv.Itoa(row), results[i].AddressIn)

				xlsx.SetCellValue(sheetName, "H"+strconv.Itoa(row), results[i].FacePicOut)
				xlsx.SetCellValue(sheetName, "I"+strconv.Itoa(row), updatedAt) // timestime out
				xlsx.SetCellValue(sheetName, "J"+strconv.Itoa(row), results[i].AttendanceDateOut)
				xlsx.SetCellValue(sheetName, "K"+strconv.Itoa(row), results[i].AttendanceTimeOut)
				xlsx.SetCellValue(sheetName, "L"+strconv.Itoa(row), results[i].AddressOut)
				xlsx.SetCellValue(sheetName, "M"+strconv.Itoa(row), results[i].Duration)
				// xlsx.SetCellValue(sheetName, "N"+strconv.Itoa(row), results[i].LocationIn)
				// xlsx.SetCellValue(sheetName, "O"+strconv.Itoa(row), results[i].LocationOut)
				row++
			}

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-1504")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			formName = strings.Replace(formName, "#", "-", 100)
			formName = strings.Replace(formName, "/", "-", 100)
			fileName := "Absensi-Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "attendance_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(err2)
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data export is available",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data Form ID is required",
			"data":    nil,
		})
		return
	}
}

func (ctr *attController) FormAttendanceListExportPost(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	var reqData objects.AttendenceForm
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		fmt.Println(err)
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
		})
		return
	}

	formID := reqData.FormID

	if formID > 0 {

		var whre1 tables.Attendances
		whre1.FormID = formID
		results, err := ctr.attMod.GetAttendanceReports(whre1, "")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {
			xlsx := excelize.NewFile()
			sheetName := "Export-Absen-Data"
			xlsx.SetSheetName(xlsx.GetSheetName(1), sheetName)

			xlsx.SetCellValue(sheetName, "A1", "Tanggal")
			xlsx.SetCellValue(sheetName, "B1", "Name")
			xlsx.SetCellValue(sheetName, "C1", "Kontak")
			xlsx.SetCellValue(sheetName, "D1", "Selfie Check-In")
			xlsx.SetCellValue(sheetName, "E1", "Waktu Check-In")
			xlsx.SetCellValue(sheetName, "F1", "Lokasi Check-In")
			xlsx.SetCellValue(sheetName, "G1", "Selfie Check-Out")
			xlsx.SetCellValue(sheetName, "H1", "Waktu Check-Out")
			xlsx.SetCellValue(sheetName, "I1", "Lokasi Check-Out")

			row := 2
			for i := 0; i < len(results); i++ {

				xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(row), results[i].CreatedAt)
				xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(row), results[i].UserName)
				xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(row), results[i].UserPhone)
				xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(row), results[i].FacePicIn)
				xlsx.SetCellValue(sheetName, "E"+strconv.Itoa(row), results[i].AttendanceIn)
				xlsx.SetCellValue(sheetName, "F"+strconv.Itoa(row), results[i].AddressIn)
				xlsx.SetCellValue(sheetName, "G"+strconv.Itoa(row), results[i].FacePicOut)
				xlsx.SetCellValue(sheetName, "H"+strconv.Itoa(row), results[i].AttendanceOut)
				xlsx.SetCellValue(sheetName, "I"+strconv.Itoa(row), results[i].AddressOut)
				row++
			}

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("200601020405")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			fileName := "Export-" + formName + "-" + dateFormat + "0" + strconv.Itoa(userID)
			fileGroup := "attendance_download"
			fileExtention := "xlsx"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			err2 := xlsx.SaveAs(fileLocation)
			if err2 != nil {
				fmt.Println(err2)
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data Form ID is required",
			"data":    nil,
		})
		return
	}
}

// func (ctr *attController) InsertAttendanceOut(c *gin.Context) {
// 	claims := jwt.ExtractClaims(c)
// 	respondenID, _ := strconv.Atoi(claims["id"].(string))

// 	fmt.Println(respondenID)

// 	var reqData objects.Attendence
// 	err := c.ShouldBindJSON(&reqData)
// 	if err != nil {
// 		fmt.Println(err)
// 		errorMessages := []string{}
// 		for _, e := range err.(validator.ValidationErrors) {
// 			errorMessage := fmt.Sprintf("Error validate %s, condition: %s", e.Field(), e.ActualTag())
// 			errorMessages = append(errorMessages, errorMessage)
// 		}

// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": errorMessages,
// 		})
// 		return
// 	}
// 	id := 0
// 	var postData tables.Attendances
// 	postData.FormID = reqData.FormID
// 	postData.UserId = respondenID
// 	postData.FacePicOut = reqData.FacePicOut
// 	// postData.AttendanceOut = time.Now().String()
// 	save, err := ctr.attMod.UpdateAttendance(id, postData)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err,
// 		})
// 		return
// 	}

// 	if save {
// 		c.JSON(http.StatusOK, gin.H{
// 			"status":  true,
// 			"message": "Successfuly submit data",
// 			"data":    nil,
// 		})
// 		return
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  false,
// 			"message": "Failed submit data",
// 			"data":    nil,
// 		})
// 		return
// 	}

// }

func (ctr *attController) FormAttendanceMapsList(c *gin.Context) {

	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	filterDate := c.Request.URL.Query().Get("date")
	page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	sortBy := c.Request.URL.Query().Get("sortby")
	sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {

		var whre1 tables.Attendances
		whre1.FormID = formID

		whrStr := ""
		whrTime := ""
		if filterDate == "" {
			whrTime = " to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char(now()::date,'yyyy-mm-dd')"
		}

		if filterDate != "" && searchKeyWord == "" {
			whrTime = " to_char(attendances.created_at::date,  'yyyy-mm-dd') = to_char('" + filterDate + "'::date,'yyyy-mm-dd')"
		}

		if searchKeyWord != "" && filterDate == "" {
			whrStr = " u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%'"
		}

		if searchKeyWord != "" && filterDate != "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char('" + filterDate + "'::date, 'yyyy-mm-dd')"
		}

		var paging objects.Paging
		paging.Page = page
		paging.Limit = limit
		paging.SortBy = sortBy
		paging.Sort = sort

		results, err := ctr.attMod.GetUserAttendanceMapReports(whre1, whrStr, whrTime, paging)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		resultAll, err := ctr.attMod.GetUserAttendanceMapReports(whre1, whrStr, whrTime, objects.Paging{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			var res []objects.AttendenceMapReport
			for i := 0; i < len(results); i++ {

				var each objects.AttendenceMapReport
				each.FormID = results[i].FormID
				each.FormName = results[i].FormName
				each.UserID = results[i].UserID
				each.UserName = results[i].UserName
				each.UserPhone = results[i].UserPhone
				each.UserAvatar = results[i].UserAvatar
				each.Latitude = results[i].Latitude
				each.Longitude = results[i].Longitude

				res = append(res, each)
			}

			totalPage := 0
			if limit > 0 {
				totalPage = len(resultAll) / limit
				if (len(resultAll) % limit) > 0 {
					totalPage = totalPage + 1
				}
			}

			var paging objects.DataRows
			paging.TotalRows = len(resultAll)
			paging.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    res,
				"paging":  paging,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}
}

func (ctr *attController) FormAttendanceListExportCSV(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	searchKeyWord := c.Request.URL.Query().Get("search")
	filterStartDate := c.Request.URL.Query().Get("start_date") + " 00:00"
	filterEndDate := c.Request.URL.Query().Get("end_date") + " 23:59"

	formID, err := strconv.Atoi(c.Param("formid"))

	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusNoContent, gin.H{
			"error": err,
		})
		return
	}

	if formID > 0 {

		var whre1 tables.Attendances
		whre1.FormID = formID

		whrStr := ""
		if searchKeyWord != "" && filterStartDate == "" && filterEndDate == "" {
			whrStr = " u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%'"
		}

		//---
		if filterStartDate != "" && filterEndDate == "" && searchKeyWord == "" {
			whrStr = " to_char(attendances.created_at,'yyyy-mm-dd') >= '" + filterStartDate + "' "
		}

		if filterStartDate == "" && filterEndDate != "" && searchKeyWord == "" {
			whrStr = " to_char(attendances.created_at,'yyyy-mm-dd') <= '" + filterEndDate + "'  "
		}

		if filterStartDate != "" && filterEndDate != "" && searchKeyWord == "" {
			whrStr = " to_char(attendances.created_at,'yyyy-mm-dd') BETWEEN '" + filterStartDate + "' AND '" + filterEndDate + "' "
		}

		//---

		// if searchKeyWord == "" && filterDate != "" {
		// 	whrStr = " to_char(attendances.created_at::date,  'yyyy-mm-dd') = to_char('" + filterDate + "'::date,'yyyy-mm-dd')"
		// }

		if searchKeyWord != "" && filterStartDate != "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at::date, 'yyyy-mm-dd') >= to_char('" + filterStartDate + "'::date, 'yyyy-mm-dd')"
		}

		if searchKeyWord != "" && filterEndDate != "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at::date, 'yyyy-mm-dd') <= to_char('" + filterEndDate + "'::date, 'yyyy-mm-dd')"
		}

		if searchKeyWord != "" && filterStartDate != "" && filterEndDate != "" {
			whrStr = " (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at,'yyyy-mm-dd') BETWEEN '" + filterStartDate + "' AND '" + filterEndDate + "' "
		}

		results, err := ctr.attMod.GetAttendanceReports(whre1, whrStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		if len(results) > 0 {

			// row header --------------------------
			var header = []string{"Nama", "Kontak", "Selfie Absen Masuk", "Tanggal Sistem Absen Masuk", "Waktu Absen Masuk", "Lokasi Absen Masuk", "Selfie Absen Keluar", "Tanggal Sistem Absen Masuk", "Waktu Absen Keluar", "Lokasi Absen Keluar", "Durasi Kerja"}

			var exportData = make([][]string, len(results)+1)
			exportData[0] = header

			// rows data
			no := 1
			for i := 0; i < len(results); i++ {
				updatedAt := ""
				if results[i].AttendanceOut != "" {
					updatedAt = results[i].UpdatedAt
				}
				var rowData = []string{results[i].UserName, results[i].UserPhone, results[i].FacePicIn, results[i].CreatedAt, results[i].AttendanceIn, results[i].AddressIn, results[i].FacePicOut, updatedAt, results[i].AttendanceOut, results[i].AddressOut, results[i].Duration}

				exportData[no] = rowData
				no++
			}

			// CONFIG file --------------------------------------------------
			var fieldForm tables.Forms
			fieldForm.ID = formID
			getForm, _ := ctr.formMod.GetFormRow(fieldForm)

			today := time.Now()
			dateFormat := today.Format("02012006-150405")

			formName := strings.Replace(getForm.Name, " ", "-", 100)
			fileName := "Absensi-Snapin-" + formName + "-" + dateFormat + "-" + strconv.Itoa(userID)
			fileGroup := "attendance_download"
			fileExtention := "csv"
			fileLocation := "file/" + fileName + "." + fileExtention // local path file location

			f, err := os.Create(fileLocation)
			// defer f.Close()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			w := csv.NewWriter(f)
			// w.UseCRLF = false
			err = w.WriteAll(exportData) // calls Flush internally
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			}

			var obj objects.FileRes
			obj.File, err = ctr.helper.UploadFileExtToOSS(fileLocation, fileName, fileGroup, fileExtention)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  false,
					"message": err,
					"data":    nil,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data export is available",
				"data":    obj,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	}
}

func (ctr *attController) FormAttendanceMapsList2(c *gin.Context) {

	ID := c.Param("formid")
	formID, err := strconv.Atoi(ID)
	if err != nil {
		fmt.Println("InsertFormFieldRule", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	searchKeyWord := c.Request.URL.Query().Get("search")
	filterDate := c.Request.URL.Query().Get("date")
	company_id, _ := strconv.Atoi(c.Request.URL.Query().Get("company_id"))

	// page, _ := strconv.Atoi(c.Request.URL.Query().Get("page"))
	// limit, _ := strconv.Atoi(c.Request.URL.Query().Get("limit"))
	// sortBy := c.Request.URL.Query().Get("sortby")
	// sort := c.Request.URL.Query().Get("sort")

	if formID > 0 {

		var whre1 tables.Attendances
		whre1.FormID = formID

		whrStr := ""
		whrTime := ""
		whrComp := ""

		if filterDate == "" {
			whrTime = " AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char(now()::date,'yyyy-mm-dd')"
		}

		if filterDate != "" && searchKeyWord == "" {
			whrTime = " AND to_char(attendances.created_at::date,  'yyyy-mm-dd') = to_char('" + filterDate + "'::date,'yyyy-mm-dd')"
		}

		if searchKeyWord != "" && filterDate == "" {
			whrStr = " AND u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%'"
		}

		if searchKeyWord != "" && filterDate != "" {
			whrStr = " AND (u.name ilike '%" + searchKeyWord + "%' OR u.phone ilike '%" + searchKeyWord + "%') AND to_char(attendances.created_at::date, 'yyyy-mm-dd') = to_char('" + filterDate + "'::date, 'yyyy-mm-dd')"
		}
		if company_id != 0 {
			whrComp = " AND o.id =  " + strconv.Itoa(company_id) + ""
		}
		// var paging objects.Paging
		// paging.Page = page
		// paging.Limit = limit
		// paging.SortBy = sortBy
		// paging.Sort = sort

		// fmt.Println("whre1, whrStr, whrTime ::", whre1, whrStr, whrTime)
		results, err := ctr.attMod.GetUserAttendanceMap3Reports(whre1, whrStr, whrTime, whrComp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		// resultAll, err := ctr.attMod.GetUserAttendanceMapReports(whre1, whrStr, whrTime, objects.Paging{})
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"error": err,
		// 	})
		// 	return
		// }

		if len(results) > 0 {

			// var res []objects.AttendenceMapReport
			// for i := 0; i < len(results); i++ {

			// 	var each objects.AttendenceMapReport
			// 	each.ID = results[i].ID
			// 	each.FormID = results[i].FormID
			// 	each.FormName = results[i].FormName
			// 	each.UserID = results[i].UserID
			// 	each.UserName = results[i].UserName
			// 	each.UserPhone = results[i].UserPhone
			// 	each.UserAvatar = results[i].UserAvatar
			// 	each.Latitude = results[i].Latitude
			// 	each.Longitude = results[i].Longitude

			// 	res = append(res, each)
			// }

			// totalPage := 0
			// if limit > 0 {
			// 	totalPage = len(resultAll) / limit
			// 	if (len(resultAll) % limit) > 0 {
			// 		totalPage = totalPage + 1
			// 	}
			// }

			// var paging objects.DataRows
			// paging.TotalRows = len(resultAll)
			// paging.TotalPages = totalPage

			c.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Data is available",
				"data":    results,
				// "paging":  paging,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Data is not available",
				"data":    nil,
			})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Form ID is not available",
			"data":    nil,
		})
		return
	}
}

func (ctr *attController) AutoInsertData(c *gin.Context) {
	var fields objects.MissingIDAtt
	res, _ := ctr.attMod.GetMissingIDAtt(fields)
	if len(res) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data Kosong brader",
			"data":    nil,
		})
		return
	}
	// fmt.Println(res)

	for i := 0; i < len(res); i++ {

		resCompanyID, _ := ctr.attMod.GetCompanyID(res[i].FormID, res[i].UserID)
		// fmt.Println(resCompanyID)
		if resCompanyID.OrganizationID > 0 {
			var postData objects.InsAttOrg
			postData.AttendanceID = res[i].ID
			postData.OrganizationID = resCompanyID.OrganizationID

			_, err := ctr.attMod.InsertAttendanceOrganization(postData)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Success Insert Data",
	})
	return
}

func (ctr *attController) GetListLocationAttendance(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err1 := strconv.Atoi(c.Param("formid"))
	if err1 != nil {
		fmt.Println("InsertFormFieldRule", err1)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err1,
		})
		return
	}

	typeAtt := c.Request.URL.Query().Get("type")

	if typeAtt == "check_in" {
		var data objects.ObjectFormAttendanceLocations
		data.FormID = formID
		data.IsCheckIn = true
		res, _ := ctr.attMod.GetFormAttendanceLocationRows(data)
		if len(res) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true,
				"message": "Data Kosong brader",
				"data":    []string{},
			})
			return
		}

		fmt.Println(userID, formID, len(res))

		var dataResultV []objects.ObjectFormAttendanceLocations
		for i := 0; i < len(res); i++ {
			var dataResult objects.ObjectFormAttendanceLocations
			dataResult.ID = res[i].ID
			dataResult.FormID = res[i].FormID
			dataResult.Name = res[i].Name
			dataResult.Location = res[i].Location
			dataResult.IsCheckIn = res[i].IsCheckIn
			dataResult.IsCheckOut = res[i].IsCheckOut
			dataResult.Latitude = res[i].Latitude
			dataResult.Longitude = res[i].Longitude
			dataResult.Radius = res[i].Radius
			dataResultV = append(dataResultV, dataResult)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data Ready",
			"data":    dataResultV,
		})
		return
	} else if typeAtt == "check_out" {

		var data objects.ObjectFormAttendanceLocations
		data.FormID = formID
		data.IsCheckOut = true
		res, _ := ctr.attMod.GetFormAttendanceLocationRows(data)
		if len(res) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true,
				"message": "Data Kosong brader",
				"data":    []string{},
			})
			return
		}

		fmt.Println(userID, formID, len(res))

		var dataResultV []objects.ObjectFormAttendanceLocations
		for i := 0; i < len(res); i++ {
			var dataResult objects.ObjectFormAttendanceLocations
			dataResult.ID = res[i].ID
			dataResult.FormID = res[i].FormID
			dataResult.Name = res[i].Name
			dataResult.Location = res[i].Location
			dataResult.IsCheckIn = res[i].IsCheckIn
			dataResult.IsCheckOut = res[i].IsCheckOut
			dataResult.Latitude = res[i].Latitude
			dataResult.Longitude = res[i].Longitude
			dataResult.Radius = res[i].Radius
			dataResultV = append(dataResultV, dataResult)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data Ready",
			"data":    dataResultV,
		})
		return

	} else {

		var data objects.ObjectFormAttendanceLocations
		data.FormID = formID
		res, _ := ctr.attMod.GetFormAttendanceLocationRows(data)
		if len(res) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  true,
				"message": "Data Kosong brader",
				"data":    []string{},
			})
			return
		}

		fmt.Println(userID, formID, len(res))

		var dataResultV []objects.ObjectFormAttendanceLocations
		for i := 0; i < len(res); i++ {
			var dataResult objects.ObjectFormAttendanceLocations
			dataResult.ID = res[i].ID
			dataResult.FormID = res[i].FormID
			dataResult.Name = res[i].Name
			dataResult.Location = res[i].Location
			dataResult.IsCheckIn = res[i].IsCheckIn
			dataResult.IsCheckOut = res[i].IsCheckOut
			dataResult.Latitude = res[i].Latitude
			dataResult.Longitude = res[i].Longitude
			dataResult.Radius = res[i].Radius
			dataResultV = append(dataResultV, dataResult)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Data Ready",
			"data":    dataResultV,
		})
		return
	}

}

func (ctr *attController) GetListLocationAttendanceDash(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, _ := strconv.Atoi(claims["id"].(string))

	formID, err1 := strconv.Atoi(c.Param("formid"))
	if err1 != nil {
		fmt.Println("InsertFormFieldRule", err1)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err1,
		})
		return
	}

	var data objects.ObjectFormAttendanceLocations
	data.FormID = formID
	data.IsCheckIn = true
	res, _ := ctr.attMod.GetFormAttendanceLocationRows(data)
	if len(res) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Data Kosong brader",
			"data":    nil,
		})
		return
	}

	fmt.Println(userID, formID, len(res))

	var dataResultV []objects.ObjectFormAttendanceLocations
	for i := 0; i < len(res); i++ {
		var dataResult objects.ObjectFormAttendanceLocations
		dataResult.ID = res[i].ID
		dataResult.FormID = res[i].FormID
		dataResult.Name = res[i].Name
		dataResult.Location = res[i].Location
		dataResult.IsCheckIn = res[i].IsCheckIn
		dataResult.IsCheckOut = res[i].IsCheckOut
		dataResult.Latitude = res[i].Latitude
		dataResult.Longitude = res[i].Longitude
		dataResult.Radius = res[i].Radius
		dataResultV = append(dataResultV, dataResult)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data Ready",
		"data":    dataResultV,
	})
	return

}
