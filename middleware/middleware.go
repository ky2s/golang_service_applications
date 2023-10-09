package middleware

import (
	"fmt"
	"net/mail"
	"snapin-form/models"
	"snapin-form/objects"
	"snapin-form/tables"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

var identityKey = "id"
var roleIDKey = "role_id"
var organizationIDKey = "organization_id"

func SetupMiddleware(db *gorm.DB) *jwt.GinJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "jwt",
		Key:         []byte("#snapin-new#"),
		Timeout:     time.Duration(24*365) * time.Hour,
		MaxRefresh:  time.Duration(24*365) * time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			// simpan data login (save token)
			fmt.Println("PayloadFunc ")

			if v, ok := data.(*objects.UserLogin); ok {

				tokenResult := jwt.MapClaims{
					identityKey:       v.UserID,
					roleIDKey:         v.RoleID,
					organizationIDKey: v.OrganizationID,
				}

				// save token
				// saveToken := db.Debug().Scopes(models.SchemaPublic("users"))
				// var userData models.Users
				// db.Model(&userData).Update("remember_token", tokenResult)

				fmt.Println("dataaaa payload----- ", v.UserID, v.Email, v.RoleID, tokenResult)

				return tokenResult
			}

			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			fmt.Println("IdentityHandler ----- ")
			claims := jwt.ExtractClaims(c)

			fmt.Println("identityKey ----", len(claims), identityKey, roleIDKey, "---", claims[identityKey].(string))

			claimsOrganizationIDKeyString := ""
			if len(claims) == 4 {
				if claims[identityKey] == nil || claims[roleIDKey] == nil {
					return &objects.UserLogin{}
				}

				claimsOrganizationIDKeyString = ""
			} else if len(claims) >= 5 {
				if claims[identityKey] == nil || claims[roleIDKey] == nil || claims[organizationIDKey] == nil {
					return &objects.UserLogin{}
				}

				claimsOrganizationIDKeyString = claims[organizationIDKey].(string)
			}

			return &objects.UserLogin{
				UserID:         claims[identityKey].(string),
				RoleID:         claims[roleIDKey].(string),
				OrganizationID: claimsOrganizationIDKeyString,
			}

		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//pengecekan token yg sudah disimpan di DB
			fmt.Println("Authorizator ----- ")
			fmt.Println("data Authorizator tables user------->>", data.(*objects.UserLogin).UserID, "---- ", data.(*objects.UserLogin).OrganizationID)

			// if data.(*objects.UserLogin).OrganizationID == "" {
			// 	return false
			// }

			if v, ok := data.(*objects.UserLogin); ok {

				fmt.Println("v.UserID------>>>>>>", v.UserID)
				var userData tables.Users

				errc := db.Debug().Scopes(models.SchemaUsr("users")).Where("status is true").First(&userData, "id = ? ", v.UserID).Error
				if errc != nil {
					fmt.Println(errc)
					return false
				}

				fmt.Println("return userData.ID------>>>>>>", userData.ID)
				if userData.ID > 0 {
					return true
				}
			}

			fmt.Println("---false---->>", data)

			return false
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// pengecekan akun login
			fmt.Println("Authenticator ----- ")

			var loginVals objects.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				fmt.Println("Error: JWT", jwt.ErrMissingLoginValues)
				return "", jwt.ErrMissingLoginValues
			}

			var userData tables.Users
			errc := db.Debug().Scopes(models.SchemaUsr("users")).First(&userData, "lower(email) = lower(?) ", loginVals.Email).Error
			if errc != nil {
				fmt.Println(errc)
			}

			fmt.Println("userData.ID---", userData.ID)
			organizationID := 0
			if userData.RoleID == 1 {
				// super admin/owner
				// userOrgData.ID = 1
				var userOrgData tables.Organizations

				errc := db.Debug().Scopes(models.SchemaMstr("organizations")).Where("is_default", "true").Where("created_by", userData.ID).First(&userOrgData).Error
				if errc != nil {
					fmt.Println("jwt.ErrForbidden", jwt.ErrForbidden)
					fmt.Println(errc)
				}

				organizationID = userOrgData.ID

			} else if userData.RoleID == 2 {
				var userOrgData tables.UserOrganizations
				// admin inviting
				errc := db.Debug().Scopes(models.SchemaUsr("user_organizations")).Where("user_id", userData.ID).First(&userOrgData).Error
				if errc != nil {
					fmt.Println("jwt.ErrForbidden", jwt.ErrForbidden)
					fmt.Println(errc)
				}

				organizationID = userOrgData.OrganizationID
			}

			if userData.ID > 0 && organizationID > 0 {

				checkPassword := VerifyPassword(loginVals.Password, userData.Password)
				if checkPassword {
					fmt.Println("getUserData---", userData)

					// save tokeN here
					return &objects.UserLogin{
						UserID:         strconv.Itoa(userData.ID),
						Email:          userData.Email,
						RoleID:         strconv.Itoa(userData.RoleID),
						OrganizationID: strconv.Itoa(organizationID),
					}, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			fmt.Println("Unauthorized ----- ")

			c.JSON(code, gin.H{
				"code":    code,
				"status":  false,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		fmt.Println("Err: ", err)
		return nil
	}

	return authMiddleware
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
