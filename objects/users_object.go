package objects

type UserLogin struct {
	UserID         string `json:"user_id"`
	RoleID         string `json:"role_id"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
}

type User struct {
	UserID    string `json:"user_id" form:"user_id"`
	Email     string `json:"email" form:"email"`
	Phone     string `json:"phone" form:"phone"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
}

type UserMember struct {
	ID          int                      `json:"id" form:"id"`
	Email       string                   `json:"email" form:"email"`
	Phone       string                   `json:"phone" form:"phone"`
	Name        string                   `json:"name" form:"name"`
	StatusID    int                      `json:"status_id" form:"status_id"`
	StatusName  string                   `json:"status_name" form:"status_name"`
	Permissions []FormUserPermissionJoin `json:"permissions,omitempty" form:"permissions"`
}

type FormUserPermissionJoin struct {
	ID             int    `json:"id" form:"id"`
	FormUserID     int    `json:"form_user_id,omitempty" form:"form_user_id"`
	PermissionID   int    `json:"permission_id" form:"permission_id"`
	PermissionName string `json:"permission_name" form:"permission_name"`
	Status         bool   `json:"status" form:"status"`
}

type Permissions struct {
	ID             int
	PermissionName string
	Status         bool
}

type Users struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Phone    string `json:"phone" form:"phone" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type Login struct {
	Email    string `form:"email" json:"email" binding:"required" `
	Phone    string `form:"phone" json:"phone" `
	Username string `form:"username" json:"username" `
	Password string `form:"password" json:"password" binding:"required"`
}

type UserOrganizationPermissions struct {
	Data []UOPermissionCheck `json:"data" form:"data"`
}

type UOPermissionCheck struct {
	ID        int  `json:"id" form:"id"`
	IsChecked bool `json:"is_checked" form:"is_checked"`
}

type AdminEks struct {
	ID               int    `json:"id,omitempty" form:"id"`
	UserID           int    `json:"user_id" form:"user_id"`
	OrganizationID   int    `json:"organization_id" form:"organization_id"`
	Name             string `json:"name" form:"name"`
	Email            string `json:"email" form:"email"`
	Phone            string `json:"phone" form:"phone"`
	OrganizationName string `json:"organization_name" form:"organization_name"`
	TotalForm        int    `json:"total_form" form:"total_form"`
}

type IDAdminEks struct {
	ID int `json:"id" form:"id"`
}

type TeamUsers struct {
	ID     int `json:"id,omitempty" form:"id"`
	TeamID int `json:"team_id"`
	UserID int `json:"user_id" form:"user_id"`
}

type FormTeams struct {
	ID             int `json:"id"`
	FormID         int `json:"form_id"`
	TeamID         int `json:"team_id"`
	OrganizationID int `json:"organization_id"`
}
