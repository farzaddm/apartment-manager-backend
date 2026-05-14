package redis

import "time"

type OTPInterFace interface {
	Save(phone string, code string, ttl time.Duration) error
	Get(phone string) (string, error)
	Delete(phone string) error
}
