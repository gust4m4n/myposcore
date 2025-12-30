package dto

type CreatePINRequest struct {
	PIN        string `json:"pin" binding:"required,len=6,numeric"`
	ConfirmPIN string `json:"confirm_pin" binding:"required,len=6,numeric"`
}

type ChangePINRequest struct {
	OldPIN     string `json:"old_pin" binding:"required,len=6,numeric"`
	NewPIN     string `json:"new_pin" binding:"required,len=6,numeric"`
	ConfirmPIN string `json:"confirm_pin" binding:"required,len=6,numeric"`
}

type AdminChangePINRequest struct {
	Username   string `json:"username" binding:"required"`
	PIN        string `json:"pin" binding:"required,len=6,numeric"`
	ConfirmPIN string `json:"confirm_pin" binding:"required,len=6,numeric"`
}
