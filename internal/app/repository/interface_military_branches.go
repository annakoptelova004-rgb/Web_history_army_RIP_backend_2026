package repository

import (
	"army/internal/app/ds"
	"io"
)

type MilitaryBranchRepository interface {
	GetMilitaryBranches() ([]ds.MilitaryBranch, error)
	GetMilitaryBranch(id int64) (*ds.MilitaryBranch, error)
	GetMilitaryBranchDraft(creatorID int64) (*ds.MilitaryBranch, error)
	GetMilitaryBranchesBySpeed(speed int) ([]ds.MilitaryBranch, error)
	UpdateMilitaryBranchImage(id int64, fileName string) error
	UpdateMilitaryBranchVideo(id int64, fileName string) error

	CreateMilitaryBranch(name string, creatorID int64) (*ds.MilitaryBranch, error)
	UploadFile(
		file io.Reader,
		size int64,
		contentType string,
		extension string,
	) (string, error)

	PublishMilitaryBranch(
		id int64,
		description string,
		speedPlain int,
		food int,
	) error

	DeleteMilitaryBranch(id int64) error
	GetLikesCount(militaryBranchID int64) (int, error)
	CreateLike(like *ds.MilitaryBranchLike) error
	DeleteLike(userID int64, militaryBranchID int64) error
	CreateUser(user *ds.User) error
	GetNextMilitaryBranch(id int64) (*ds.MilitaryBranch, error)
}
