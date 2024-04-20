package domain

import "time"

type User struct {
	Id       int64      `json:"id,omitempty"`
	Email    string     `json:"email,omitempty"`
	Name     string     `json:"name,omitempty"`
	Password string     `json:"password,omitempty"`
	Phone    string     `json:"phone,omitempty"`
	Status   string     `json:"status,omitempty"`
	CreateAt *time.Time `json:"create_at,omitempty"`
	UpdateAt *time.Time `json:"update_at,omitempty"`
}
