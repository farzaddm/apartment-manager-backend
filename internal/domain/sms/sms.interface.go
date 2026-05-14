package sms

type SMSInterface interface {
	SendOTP(phone string, code string) error
}
