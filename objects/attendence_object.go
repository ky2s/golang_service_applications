package objects

import (
	"time"

	"gorm.io/gorm"
)

type Attendence struct {
	ID          int     `json:"id"`
	FormID      int     `json:"form_id" binding:"required"`
	FacePic     string  `json:"face_pic" binding:"required"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Address     string  `json:"address"`
	OfflineTime string  `json:"offline_time"`
}

type WarnAttendence struct {
	Type             string  `json:"type"`
	FormID           int     `json:"form_id" binding:"required"`
	FacePic          string  `json:"face_pic" binding:"required"`
	Latitude         float64 `json:"latitude" binding:"required"`
	Longitude        float64 `json:"longitude" binding:"required"`
	Address          string  `json:"address"`
	OfflineTime      string  `json:"offline_time"`
	CreatedAt        string  `json:"created_at"`
	OfflineCreatedAt string  `json:"offline_created_at"`
}

type AttendenceReport struct {
	ID               int     `json:"id"`
	FormID           int     `json:"form_id"`
	UserID           int     `json:"user_id"`
	UserName         string  `json:"user_name"`
	OrganizationID   int     `json:"organization_id"`
	OrganizationName string  `json:"organization_name"`
	UserPhone        string  `json:"user_phone"`
	FacePicIn        string  `json:"face_pic_in"`
	FacePicOut       string  `json:"face_pic_out"`
	AttendanceIn     string  `json:"attendance_in"`
	AttendanceOut    string  `json:"attendance_out"`
	AddressIn        string  `json:"address_in"`
	AddressOut       string  `json:"address_out"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	LocationIn       string  `json:"location_in"`
	LocationOut      string  `json:"location_out"`
	Duration         string  `json:"duration"`
	CreatedAt        string  `json:"created_at"`
}

type AttendenceForm struct {
	FormID int `json:"form_id" binding:"required"`
}

type Attendances struct {
	ID            int
	FormID        int
	UserID        int
	AttendanceIn  time.Time `gorm:"default:null"`
	AttendanceOut time.Time `gorm:"default:null"`
	FacePicIn     string
	FacePicOut    string
	Latitude      float64
	Longitude     float64
	AddressIn     string
	AddressOut    string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type AttendenceMapReport struct {
	ID         int
	FormID     int       `json:"form_id"`
	FormName   string    `json:"form_name"`
	UserID     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	UserPhone  string    `json:"user_phone"`
	UserAvatar string    `json:"user_avatar"`
	Attendance time.Time `json:"attendance"`
	Address    string    `json:"address"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	IsCheckin  bool      `json:"is_checkin"`
}

type LastAttendance struct {
	ID            int    `json:"ID"`
	FormID        int    `json:"form_id"`
	UserID        int    `json:"user_id"`
	AttendanceIn  string `json:"attendance_in"`
	AttendanceOut string `json:"attendance_out"`
}

type MissingIDAtt struct {
	ID     int `json:"id" form:"id"`
	FormID int `json:"form_id" form:"form_id"`
	UserID int `json:"user_id" form:"user_id"`
}

type FormUserOrg struct {
	ID             int `json:"id" form:"id"`
	FormID         int `json:"form_id" form:"form_id"`
	UserID         int `json:"user_id" form:"user_id"`
	OrganizationID int `json:"organization_id" form:"organization_id"`
}

type InsAttOrg struct {
	AttendanceID   int `json:"attendance_id" form:"attendance_id"`
	OrganizationID int `json:"organization_id" form:"organization_id"`
	CreatedAt      time.Time
}

type GetOrgID struct {
	OrganizationID int `json:"organization_id" form:"organization_id"`
}
