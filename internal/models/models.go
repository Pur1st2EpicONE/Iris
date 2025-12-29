package models

const (
	StatusPending            = "pending"
	StatusCanceled           = "canceled"
	StatusFailedToSendInTime = "failed to send in time"
	StatusFailed             = "failed to send"
	StatusLate               = "running late"
	StatusSent               = "sent"
)

const (
	Email    = "email"
	Stdout   = "stdout"
	Telegram = "telegram"
)

const (
	MaxEmailLength   = 254
	MaxSubjectLength = 254
	MaxMessageLength = 254
)
