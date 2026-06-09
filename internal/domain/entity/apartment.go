package entity

type Apartment struct {
	BaseModel

	Name       string
	Province   string
	City       string
	Address    string
	PostalCode string

	Users         []User
	Announcements []Announcement
	Rules         []Rule
	InviteCodes   []InviteCode
	Units         []Unit
	Tickets       []Ticket
}
