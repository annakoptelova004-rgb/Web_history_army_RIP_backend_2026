package ds

import (
	"time"
)

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Login    string `gorm:"size:50;not null"`
	Password string `gorm:"size:255;not null"`
}

type MilitaryBranch struct {
	ID             int64  `gorm:"primaryKey"`
	Name           string `gorm:"size:100;not null"`
	Description    string `gorm:"size:255"`
	Status         string `gorm:"size:20;not null"`
	ImageURL       string `gorm:"size:255"`
	VideoURL       string `gorm:"size:255"`
	SpeedPlain     int    `gorm:"not null"`
	SpeedMountains int    `gorm:"not null"`
	SpeedForest    int    `gorm:"not null"`
	SpeedRiver     int    `gorm:"not null"`
	Food           int    `gorm:"not null;default:0"`
	CreatedAt      time.Time
	CreatorID      int64 `gorm:"not null"`
	FormedAt       time.Time

	Creator User `gorm:"foreignKey:CreatorID"`

	Terrain string `gorm:"-"`
	Speed   int    `gorm:"-"`
	Image   string `gorm:"-"`
	Video   string `gorm:"-"`
	Likes   []int  `gorm:"-"`
}

type MilitaryBranchLike struct {
	ID               int64 `gorm:"primaryKey"`
	UserID           int64 `gorm:"not null"`
	MilitaryBranchID int64 `gorm:"not null"`

	User           User           `gorm:"foreignKey:UserID"`
	MilitaryBranch MilitaryBranch `gorm:"foreignKey:MilitaryBranchID"`
}
