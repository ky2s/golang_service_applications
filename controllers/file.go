package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"snapin-form/helpers"
	"snapin-form/objects"

	"github.com/globalsign/mgo/bson"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

// interface
type FileController interface {
	FieldSaveData(c *gin.Context)
	FieldSaveDataVideo(c *gin.Context)
}

type fileController struct {
	helper helpers.Helper
}

func NewFileController(h helpers.Helper) FileController {
	return &fileController{
		helper: h,
	}
}

func (ctr *fileController) FieldSaveData(c *gin.Context) {

	var reqData objects.File
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

	if reqData.File != "" {

		fileName := ctr.helper.UploadFileToOSS(reqData.File, reqData.FileName, reqData.FileType)

		var obj objects.FileRes
		obj.File = fileName

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success upload file",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed upload file",
			"data":    nil,
		})
		return
	}

}

func (ctr *fileController) FieldSaveDataVideo(c *gin.Context) {

	myFile, err := c.FormFile("myfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	formID := c.PostForm("form_id")
	fieldID := c.PostForm("field_id")

	now := time.Now() // current local time
	unixTimeStamp := now.Format("200602011504")

	//decode timestamp
	// i, err := strconv.ParseInt("1674712481", 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err,
	// 	})
	// 	return
	// }
	// tm := time.Unix(i, 0)
	// fmt.Println(tm)

	if myFile.Filename != "" {
		extension := filepath.Ext(myFile.Filename)
		originalName := helpers.Substr(myFile.Filename, 0, 100)
		split := strings.Split(originalName, ".")

		originalFileName := split[0]
		if utf8.RuneCountInString(originalFileName) > 20 {
			originalFileName = originalFileName[0:19]
		}

		formFileInit := ""
		if formID != "" {
			formFileInit = formID + "_" + fieldID + "_" + unixTimeStamp
		}

		newFileName := originalFileName + "_" + formFileInit + "_" + bson.NewObjectId().Hex()
		ext := helpers.Substr(extension, 1, 10)

		locationFile := "file/" + newFileName + extension

		err = c.SaveUploadedFile(myFile, locationFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		fmt.Println("myFile :::", myFile.Filename, newFileName)
		fmt.Println("myFile :::", locationFile)
		fmt.Println("myFile :::", newFileName)
		fmt.Println("myFile :::", "form_data_file")
		fmt.Println("myFile :::", ext)

		fileName, err := ctr.helper.UploadFileExtToOSS(locationFile, newFileName, "form_data_file", ext)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  false,
				"message": err,
				"data":    nil,
			})
			return
		}

		var obj objects.FileRes
		obj.File = fileName

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "Success upload file",
			"data":    obj,
		})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Failed : file upload is empty",
			"data":    nil,
		})
		return
	}
}
