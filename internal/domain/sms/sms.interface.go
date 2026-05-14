package sms

import "context"

type SMSInterface interface {
	SendOTP(ctx context.Context, phone string, code string) error
}
