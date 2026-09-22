package repository

import "army/internal/app/ds"

type MilitaryBranchRepository interface {
	GetMilitaryBranches() ([]ds.MilitaryBranch, error)
	GetMilitaryBranch(id int64) (*ds.MilitaryBranch, error)
	GetMilitaryBranchDraft(creatorID int64) (*ds.MilitaryBranch, error)
	GetMilitaryBranchesBySpeed(speed int) ([]ds.MilitaryBranch, error)

	CreateMilitaryBranch(name string, creatorID int64) (*ds.MilitaryBranch, error)

	PublishMilitaryBranch(
		id int64,
		description string,
		speedPlain int,
		speedMountains int,
		speedForest int,
		speedRiver int,
		food int,
	) error

	DeleteMilitaryBranch(id int64) error
	GetLikesCount(militaryBranchID int64) (int, error)
	GetNextMilitaryBranch(id int64) (*ds.MilitaryBranch, error)
}
