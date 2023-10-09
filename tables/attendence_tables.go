package tables

import (
	"time"

	"gorm.io/gorm"
)

type Attendances struct {
	ID                          int
	FormID                      int
	UserID                      int
	AttendanceIn                time.Time `gorm:"default:null"`
	AttendanceOut               time.Time `gorm:"default:null"`
	FacePicIn                   string
	FacePicOut                  string
	GeometryIn                  string `gorm:"default:null"`
	GeometryOut                 string `gorm:"default:null"`
	AddressIn                   string
	AddressOut                  string
	FormAttendanceLocationIdIn  int `gorm:"default:null"`
	FormAttendanceLocationIdOut int `gorm:"default:null"`
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
	DeletedAt                   gorm.DeletedAt `gorm:"index"`
}

type AttendenceReport struct {
	ID                int     `json:"id"`
	FormID            int     `json:"form_id"`
	FormName          string  `json:"form_name"`
	UserID            int     `json:"user_id"`
	UserName          string  `json:"user_name"`
	UserPhone         string  `json:"user_phone"`
	OrganizationID    int     `json:"organization_id"`
	OrganizationName  string  `json:"organization_name"`
	UserAvatar        string  `json:"user_avatar"`
	FacePicIn         string  `json:"face_pic_in"`
	FacePicOut        string  `json:"face_pic_out"`
	AttendanceIn      string  `json:"attendance_in"`
	AttendanceDateIn  string  `json:"attendance_date_in"`
	AttendanceTimeIn  string  `json:"attendance_time_in"`
	AttendanceOut     string  `json:"attendance_out"`
	AttendanceDateOut string  `json:"attendance_date_out"`
	AttendanceTimeOut string  `json:"attendance_time_out"`
	AddressIn         string  `json:"address_in"`
	AddressOut        string  `json:"address_out"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Duration          string  `json:"duration"`
	LocationIn        string  `json:"location_in"`
	LocationOut       string  `json:"location_out"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type Geometry struct {
	Latitude  float64
	Longitude float64
}

type AttendenceMapReport struct {
	ID               int     `json:"id"`
	FormID           int     `json:"form_id"`
	UserID           int     `json:"user_id"`
	OrganizationName string  `json:"organization_name"`
	UserName         string  `json:"user_name"`
	UserPhone        string  `json:"user_phone"`
	UserAvatar       string  `json:"user_avatar"`
	FacePic          string  `json:"face_pic"`
	Attendance       string  `json:"attendance"`
	Address          string  `json:"address"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Duration         string  `json:"duration"`
	CreatedAt        string  `json:"created_at"`
	IsCheckin        bool    `json:"is_checkin"`
}
