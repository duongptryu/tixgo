package domain

import "time"

type CommandSendUserMailOTP struct {
	Email      string    `json:"email"`
	OTP        string    `json:"otp"`
	OccurredAt time.Time `json:"occurred_at"`
}

func NewCommandSendUserMailOTP(email string, otp string) *CommandSendUserMailOTP {
	return &CommandSendUserMailOTP{
		Email:      email,
		OTP:        otp,
		OccurredAt: time.Now(),
	}
}
