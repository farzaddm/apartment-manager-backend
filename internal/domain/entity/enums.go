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
	TicketOpen       TicketStatus = "open"
	TicketClosed     TicketStatus = "closed"
	TicketInProgress TicketStatus = "in-progress"
)

type TicketCategory string

const (
	TicketMaintenanceCategory TicketCategory = "maintenance"
	TicketPlumbingCategory    TicketCategory = "plumbing"
	TicketElectricityCategory TicketCategory = "electricity"
	TicketSecurityCategory    TicketCategory = "security"
	TicketCleaningCategory    TicketCategory = "cleaning"
	TicketParkingCategory     TicketCategory = "parking"
	TicketOtherCategory       TicketCategory = "other"
)

type RuleCategory string

const (
	RulePetPolicyCategory        RuleCategory = "pet_policy"
	RuleNoiseRegulationsCategory RuleCategory = "noise_regulations"
	RuleGymRulesCategory         RuleCategory = "gym_rules"
	RuleGarbageRecyclingCategory RuleCategory = "garbage_recycling"
	RuleParkingBylawsCategory    RuleCategory = "parking_bylaws"
	RulePoolPolicyCategory       RuleCategory = "pool_policy"
	RuleRuleOtherCategory        RuleCategory = "other"
)

type AnnouncementOrder string

const (
	WarningAnnouncementOrder       AnnouncementOrder = "warning"
	VeryImportantAnnouncementOrder AnnouncementOrder = "very_important"
	ImportantAnnouncementOrder     AnnouncementOrder = "important"
	OtherAnnouncementOrder         AnnouncementOrder = "other"
)

type TicketAccessability string

const (
	PublicTicket  TicketAccessability = "public"
	PrivateTicket TicketAccessability = "private"
)
