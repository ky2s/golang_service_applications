package objects

type Home struct {
	UserID              int    `json:"user_id"`
	UserName            string `json:"user_name"`
	UserAvatar          string `json:"user_avatar"`
	CompanyName         string `json:"company_name"`
	IsUserCompanyActive bool   `json:"is_user_company_active"`
	IsProfileComplete   bool   `json:"is_profile_complete"`
	TotalForm           int    `json:"total_form"`
	TotalResponden      int    `json:"total_responden"`
	TotalRespon         int    `json:"total_respon"`
}

type HomeApps struct {
	UserID               int     `json:"user_id"`
	UserName             string  `json:"user_name"`
	UserAvatar           string  `json:"user_avatar"`
	CompanyName          string  `json:"company_name"`
	TotalForm            int     `json:"total_form"`
	TotalRespon          int     `json:"total_respon"`
	TotalAttendance      int     `json:"total_attendance"`
	TotalAttendanceFloat float64 `json:"total_attendance_float"`
}

type HomeAdminApps struct {
	UserID                int     `json:"user_id"`
	UserName              string  `json:"user_name"`
	UserAvatar            string  `json:"user_avatar"`
	CompanyName           string  `json:"company_name"`
	TotalPerformance      int     `json:"total_performance"`
	TotalPerformanceFloat float64 `json:"total_performance_float"`
	TotalRespon           int     `json:"total_respon"`
}
