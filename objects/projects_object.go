package objects

type Projects struct {
	ID          int    `json:"id"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
	ParentID    int    `json:"parent_id,omitempty" form:"parent_id" gorm:"default:null"`
}

type ProjectRes struct {
	ID int `json:"id" form:"id"`
}

type ProjectListRes struct {
	ID          int     `json:"id" form:"id"`
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	FormCount   int     `json:"form_count,omitempty" form:"form_count"`
	FormList    []Forms `json:"form_list,omitempty" form:"form_list"`
}

type ProjectForm struct {
	ProjectID int `json:"project_id" form:"project_id" binding:"required"`
	FormID    int `json:"form_id" form:"form_id" binding:"required"`
}

type DataDetail struct {
	LastSubmissionDate string `json:"last_submission_date"`
}

type DataRows struct {
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
}

type DataRowsDetail struct {
	AllRows int `json:"all_rows"`
}
