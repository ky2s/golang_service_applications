package helpers

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"snapin-form/config"
	"snapin-form/tables"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	"gopkg.in/gomail.v2"

	"gorm.io/gorm"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"golang.org/x/crypto/bcrypt"
)

type Helper interface {
	UploadFileToOSS(file string, fileName string, fileType string) string
	UploadFileExtToOSS(file string, fileName string, fileGroup string, fileExtention string) (string, error)
	FormatDateOutput(date string) time.Time
	AddLogBook(UserID int, permissionsID int, formID int) error
	SendGoMail(to string, cc string, subject, message string) (bool, error)
}

type helperConfiguration struct {
	conf config.Configurations
	db   *gorm.DB
}

func NewHelper(c config.Configurations, dbg *gorm.DB) Helper {
	return &helperConfiguration{
		conf: c,
		db:   dbg,
	}
}

func SchemaUsr(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Table("usr" + "." + tableName)
	}
}

func (con *helperConfiguration) AddLogBook(UserID int, permissionsID int, formID int) error {

	var userData tables.AddLog
	userData.UserID = UserID
	userData.PermissionID = permissionsID
	userData.FormID = formID

	err := con.db.Scopes(SchemaUsr("form_user_activity_logs")).Create(&userData).Error
	if err != nil {
		return err
	}

	return err
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func Substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func BaseURLDash() string {

	// baseURL := "https://" + ctx.Request.Host + ":" + strconv.Itoa(h.conf.Server.Port)
	baseURL := ""

	var c *gin.Context
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL = scheme + "://" + c.Request.Host

	return baseURL
}

func DateNow() time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	return now
}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page == 0 {
			page = 1
		}

		limit, _ := strconv.Atoi(q.Get("limit"))
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

func (h *helperConfiguration) FormatDateOutput(stringDate string) time.Time {
	format := "2006-01-02"
	date, err := time.Parse(format, stringDate)
	if err != nil {
		fmt.Println("Error Parse date:", err)
		return time.Time{}
	}

	return date
}

func (h *helperConfiguration) UploadFileToOSS(file string, fileName string, fileType string) string {

	OSSAccessKeyID := h.conf.OSS_ACCESS_KEY_ID
	OSSAccessKeySecret := h.conf.OSS_SECRET_ACCESS_KEY
	OSSBucket := h.conf.OSS_BUCKET
	OSSRegion := h.conf.OSS_REGION
	OSSUrl := h.conf.OSS_URL
	OSSEndPoint := h.conf.OSS_ENDPOINT
	folder := h.conf.ENV_TYPE

	fmt.Println("OSS----", OSSEndPoint, OSSAccessKeyID, OSSAccessKeySecret)

	client, err := oss.New(OSSEndPoint, OSSAccessKeyID, OSSAccessKeySecret)
	if err != nil {
		fmt.Println("Error 1:", err)
		os.Exit(-1)
	}

	// Specify the name of the bucket. Example: examplebucket.
	bucket, err := client.Bucket(OSSBucket)
	if err != nil {
		fmt.Println("Error 2:", err)
		os.Exit(-1)
	}

	// // Set the storage class of the object to Infrequent Access (IA).
	// storageType := oss.ObjectStorageClass(oss.StorageIA)

	// // Set the access control list (ACL) of the object to private.
	// objectAcl := oss.ObjectACL(oss.ACLPrivate)

	objectFile := file
	idx := strings.Index(objectFile, ";base64,")

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(objectFile[idx+8:]))

	buff := bytes.Buffer{}
	_, err = buff.ReadFrom(reader)

	// fmt.Println("reader----", reader)

	_, fm, err := image.DecodeConfig(bytes.NewReader(buff.Bytes()))
	fileLastName := fileName + bson.NewObjectId().Hex()
	// Upload the "Hello OSS" string to the exampleobject.txt object in the exampledir directory.

	fmt.Println(OSSUrl, OSSRegion, folder, "fm ----------->", fm, " ----end")
	objectURL := folder + "/" + fileType + "/" + fileLastName + ".jpg"

	err = bucket.PutObject(objectURL, bytes.NewReader(buff.Bytes()))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	err = bucket.SetObjectACL(objectURL, oss.ACLPublicRead)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	aclRes, err := bucket.GetObjectACL(objectURL)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
	fmt.Println("Object ACL:", aclRes.ACL)

	return OSSUrl + "/" + objectURL
}

func (h *helperConfiguration) UploadFileExtToOSS(fileLocation string, fileName string, fileGroup string, fileExtention string) (string, error) {

	OSSAccessKeyID := h.conf.OSS_ACCESS_KEY_ID
	OSSAccessKeySecret := h.conf.OSS_SECRET_ACCESS_KEY
	OSSBucket := h.conf.OSS_BUCKET
	// OSSRegion := h.conf.OSS_REGION
	OSSUrl := h.conf.OSS_URL
	OSSEndPoint := h.conf.OSS_ENDPOINT
	folder := h.conf.ENV_TYPE

	client, err := oss.New(OSSEndPoint, OSSAccessKeyID, OSSAccessKeySecret)
	if err != nil {
		fmt.Println("Error 1:", err)
		return "", err
	}

	// Specify the name of the bucket. Example: examplebucket.
	bucket, err := client.Bucket(OSSBucket)
	if err != nil {
		fmt.Println("Error 2:", err)
		return "", err
	}

	objectURL := folder + "/" + fileGroup + "/" + fileName + "." + fileExtention

	// err = bucket.PutObject(objectURL, bytes.NewReader(buff.Bytes()))
	err = bucket.PutObjectFromFile(objectURL, fileLocation)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	err = bucket.SetObjectACL(objectURL, oss.ACLPublicRead)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	aclRes, err := bucket.GetObjectACL(objectURL)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	fmt.Println("Object ACL:", aclRes.ACL)

	// clear file in project
	// time.Sleep(time.Second * 3)
	fmt.Println("fileName :::", fileName+"."+fileExtention)
	err = os.Remove("./file/" + fileName + "." + fileExtention)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	return OSSUrl + "/" + objectURL, nil
}

func (h *helperConfiguration) SendMail(to []string, cc []string, subject, message string) error {
	CONFIG_SMTP_HOST := h.conf.SMTP_HOST
	CONFIG_SMTP_PORT := h.conf.SMTP_PORT
	CONFIG_EMAIL := h.conf.EMAIL
	CONFIG_PASSWORD := h.conf.PASSWORD
	SENDER_NAME := h.conf.SENDER_NAME

	body := "From: " + SENDER_NAME + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	auth := smtp.PlainAuth("", CONFIG_EMAIL, CONFIG_PASSWORD, CONFIG_SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", CONFIG_SMTP_HOST, CONFIG_SMTP_PORT)

	err := smtp.SendMail(smtpAddr, auth, CONFIG_EMAIL, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}
	return nil
}

func SendWA(phone string, msg string) bool {

	url := "https://apivalen.waviro.com/api/sendwa"
	method := "POST"

	payload := strings.NewReader(`{"nohp": "` + phone + `","pesan": "` + msg + `","notifyurl": ""}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("SecretKey", "0zmrbqirmfhI7fKAbSHh")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("err string Body")
		fmt.Println(err)
		return false
	}
	fmt.Println("string(body) : ", string(body))

	if strings.Contains(string(body), "504 Gateway Time-out") == true {
		return false
	}

	return true
}

func (h *helperConfiguration) SendGoMail(to string, cc string, subject, message string) (bool, error) {

	CONFIG_SMTP_HOST := "mail.snap-in.co.id"
	CONFIG_SMTP_PORT := 465
	CONFIG_EMAIL := "no-reply@snap-in.co.id"
	CONFIG_PASSWORD := "yG8;kEM$((>ynmRKH<5%"
	CONFIG_SENDER_NAME := "Snap-in <no-reply@snap-in.co.id>"

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", to)
	if cc != "" {
		mailer.SetAddressHeader("Cc", cc, "")
	}
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)
	// mailer.Attach("./sample.png")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_EMAIL,
		CONFIG_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("ERROR here : ", err.Error())
		// log.Fatal(err.Error())
		fmt.Println("ERROR end ---- ")
		return false, err
	}

	log.Println("Mail sent!")

	return true, nil
}
