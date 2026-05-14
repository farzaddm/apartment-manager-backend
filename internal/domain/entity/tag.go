package entity

type Tag struct {
	BaseModel

	Name string

	Tickets       []Ticket       `gorm:"many2many:ticket_tags;"`
	Announcements []Announcement `gorm:"many2many:announcement_tags;"`
}
