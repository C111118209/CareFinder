package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"user_id"`
	Email        string    `gorm:"varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"varchar(255);not null" json:"-"` // Do not expose password hash
	Role         string    `gorm:"varchar(50);not null" json:"role"` // 'user', 'caregiver', 'admin'
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CaregiverProfile struct {
	ProfileID  uint          `gorm:"primaryKey" json:"profile_id"`
	UserID     uint          `gorm:"not null" json:"user_id"`
	User       User          `gorm:"foreignKey:UserID"`
	FullName   string        `gorm:"varchar(100);not null" json:"full_name"`
	Gender     string        `gorm:"varchar(10)" json:"gender"`
	Phone      string        `gorm:"varchar(20)" json:"phone"`
	Address    string        `gorm:"text" json:"address"`
	Bio        string        `gorm:"text" json:"bio"`
	AvgRating  float32       `gorm:"numeric(2,1);default:0.0" json:"avg_rating"`
	ServiceRate float64      `gorm:"decimal" json:"service_rate"`
	Licenses   []License     `gorm:"foreignKey:CaregiverID" json:"licenses"`
	Availabilities []Availability `gorm:"foreignKey:CaregiverID" json:"availabilities"`
}

type License struct {
	LicenseID   uint      `gorm:"primaryKey" json:"license_id"`
	CaregiverID uint      `gorm:"not null" json:"-"` // Foreign key to CaregiverProfile
	Name        string    `gorm:"varchar(100);not null" json:"name"`
	IssueDate   time.Time `json:"issue_date"`
	ExpiryDate  time.Time `json:"expiry_date"`
	Status      string    `gorm:"varchar(50);default:'pending'" json:"status"` // 'pending', 'approved', 'rejected'
	ProofURL    string    `gorm:"varchar(255)" json:"proof_url"`
}

type Availability struct {
	AvailabilityID uint    `gorm:"primaryKey" json:"availability_id"`
	CaregiverID    uint    `gorm:"not null" json:"-"` // Foreign key to CaregiverProfile
	DayOfWeek      int     `gorm:"not null" json:"day_of_week"` // 1=Mon, 7=Sun
	StartTime      string  `gorm:"type:time" json:"start_time"`
	EndTime        string  `gorm:"type:time" json:"end_time"`
}

type Contract struct {
	ContractID        uint      `gorm:"primaryKey" json:"contract_id"`
	UserID            uint      `gorm:"not null" json:"user_id"`
	User              User      `gorm:"foreignKey:UserID"`
	CaregiverID       uint      `gorm:"not null" json:"caregiver_id"`
	Caregiver         User      `gorm:"foreignKey:CaregiverID"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Status            string    `gorm:"varchar(50);default:'pending'" json:"status"` // 'pending', 'active', 'completed', 'canceled'
	IsRenewal         bool      `gorm:"default:false" json:"is_renewal"`
	OriginalContractID *uint    `json:"original_contract_id"`
	OriginalContract  *Contract `gorm:"foreignKey:OriginalContractID"`
}

type Review struct {
	ReviewID    uint      `gorm:"primaryKey" json:"review_id"`
	ContractID  uint      `gorm:"unique;not null" json:"contract_id"`
	Contract    Contract  `gorm:"foreignKey:ContractID"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID"`
	CaregiverID uint      `gorm:"not null" json:"caregiver_id"`
	Caregiver   User      `gorm:"foreignKey:CaregiverID"`
	Rating      int       `gorm:"not null" json:"rating"` // 1 - 5
	Comment     string    `gorm:"text" json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

// Note: In CaregiverProfile, I'm referencing User directly.
// The user's schema has CaregiverProfile.user_id -> User.user_id
// and License.caregiver_id -> CaregiverProfile.profile_id
// I've adjusted this slightly for GORM's conventions where the relationship
// is more direct (e.g. License links to CaregiverProfile). I'll also link
// contracts and reviews to the User ID directly, which simplifies queries.
// The user schema had caregiver_id pointing to user_id in contract/review, which is correct.
