package service

import (
	"apartment-manager-backend/internal/domain/repository/redis"
	"apartment-manager-backend/internal/domain/sms"
	"fmt"
	"math/rand"
	"time"
)

type SendOtpService struct {
	otpRepo redis.OTPInterFace
	smsRepo sms.SMSInterface
}

func NewSendOtpService(otpRepo redis.OTPInterFace, smsRepo sms.SMSInterface) *SendOtpService {
	return &SendOtpService{
		otpRepo: otpRepo,
		smsRepo: smsRepo,
	}
}

func (u *SendOtpService) Generate() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (u *SendOtpService) Execute(phone string) error {

	code := u.Generate()

	err := u.otpRepo.Save(phone, code, 2*time.Minute)
	if err != nil {
		return err
	}

	message := "your code : " + code

	err = u.smsRepo.SendOTP(phone, message)
	if err != nil {
		return err
	}

	return nil
}
