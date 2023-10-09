package objects

type AdminInvites struct {
	FormID int `json:"form_id"`
	// Email  string `json:"email"`
	// Type   string `json:"type"`
	AdminInvites []EmailInvites `json:"admin_invites"`
}

type EmailInvites struct {
	Email string `json:"email"`
	Type  string `json:"type"`
}

type SelectOrganization struct {
	UserSenderID           int    `json:"user_sender_id"`
	UserSenderTypeID       int    `json:"user_sender_type_id"`
	OrganizationSenderID   int    `json:"organization_sender_id"`
	UserReceiverID         int    `json:"user_receiver_id"`
	UserReceiverTypeID     int    `json:"user_receiver_type_id"`
	OrganizationReceiverID int    `json:"organization_receiver_id"`
	FormID                 int    `json:"form_id"`
	AccessType             string `json:"access_type"`
	IsQuotaSharing         bool   `json:"is_quota_sharing"`
}
