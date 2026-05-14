package entity

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleManager  UserRole = "manager"
	RoleResident UserRole = "resident"
)

type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
)

type TicketStatus string

const (
	TicketOpen   TicketStatus = "open"
	TicketClosed TicketStatus = "closed"
)
