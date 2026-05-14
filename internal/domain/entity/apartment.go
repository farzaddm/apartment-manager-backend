package entity

type Apartment struct {
	BaseModel

	Name       string
	Province   string
	City       string
	Address    string
	PostalCode string

	UnitCount int

	Users         []User
	Announcements []Announcement
	Rules         []Rule
	InviteCodes   []InviteCode
}
