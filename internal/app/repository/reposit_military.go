package repository

import (
	"army/internal/app/ds"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}

func (r *Repository) GetMilitaryBranches() ([]ds.MilitaryBranch, error) {
	var militaryBranches []ds.MilitaryBranch

	err := r.db.
		Where("status != ?", "удален").
		Find(&militaryBranches).Error

	if err != nil {
		return nil, err
	}

	for i := range militaryBranches {
		militaryBranches[i].Image = militaryBranches[i].ImageURL
		militaryBranches[i].Video = militaryBranches[i].VideoURL
	}

	return militaryBranches, nil
}

func (r *Repository) GetMilitaryBranch(id int64) (*ds.MilitaryBranch, error) {
	var militaryBranch ds.MilitaryBranch

	err := r.db.
		Where("id = ? AND status != ?", id, "удален").
		First(&militaryBranch).Error

	if err != nil {
		return nil, err
	}

	militaryBranch.Image = militaryBranch.ImageURL
	militaryBranch.Video = militaryBranch.VideoURL

	return &militaryBranch, nil
}

func (r *Repository) GetMilitaryBranchDraft(creatorID int64) (*ds.MilitaryBranch, error) {
	var draft ds.MilitaryBranch

	err := r.db.
		Where("status = ? AND creator_id = ?", "черновик", creatorID).
		First(&draft).Error

	if err != nil {
		return nil, err
	}

	return &draft, nil
}

func (r *Repository) GetMilitaryBranchesBySpeed(speed int) ([]ds.MilitaryBranch, error) {
	var militaryBranches []ds.MilitaryBranch

	err := r.db.
		Where("status = ? AND speed_plain = ?", "опубликован", speed).
		Find(&militaryBranches).Error

	if err != nil {
		return nil, err
	}

	for i := range militaryBranches {
		militaryBranches[i].Image = militaryBranches[i].ImageURL
		militaryBranches[i].Video = militaryBranches[i].VideoURL
	}

	return militaryBranches, nil
}

func (r *Repository) CreateMilitaryBranch(name string, creatorID int64) (*ds.MilitaryBranch, error) {
	var draft ds.MilitaryBranch

	result := r.db.
		Where("status = ? AND creator_id = ?", "черновик", creatorID).
		First(&draft)

	if result.Error == nil {
		return &draft, nil
	}

	draft = ds.MilitaryBranch{
		Name:      name,
		Status:    "черновик",
		CreatorID: creatorID,
	}

	err := r.db.Create(&draft).Error
	if err != nil {
		return nil, err
	}

	return &draft, nil
}

func (r *Repository) PublishMilitaryBranch(
	id int64,
	description string,
	speedPlain int,
	speedMountains int,
	speedForest int,
	speedRiver int,
	food int,
) error {
	return r.db.
		Model(&ds.MilitaryBranch{}).
		Where("id = ? AND status = ?", id, "черновик").
		Updates(map[string]interface{}{
			"description":     description,
			"speed_plain":     speedPlain,
			"speed_mountains": speedMountains,
			"speed_forest":    speedForest,
			"speed_river":     speedRiver,
			"food":            food,
			"status":          "опубликован",
			"formed_at":       time.Now(),
		}).Error
}

func (r *Repository) DeleteMilitaryBranch(id int64) error {
	query := `
        UPDATE military_branches
        SET status = 'удален'
        WHERE id = $1
        RETURNING id
    `

	row := r.db.Raw(query, id).Row()

	var deletedID int64

	err := row.Scan(&deletedID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLikesCount(militaryBranchID int64) (int, error) {
	var count int64

	err := r.db.
		Table("military_branch_likes").
		Where("military_branch_id = ?", militaryBranchID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) GetNextMilitaryBranch(id int64) (*ds.MilitaryBranch, error) {
	var militaryBranch ds.MilitaryBranch

	err := r.db.
		Where("status = ? AND id > ?", "опубликован", id).
		Order("id ASC").
		First(&militaryBranch).Error

	if err != nil {
		err = r.db.
			Where("status = ?", "опубликован").
			Order("id ASC").
			First(&militaryBranch).Error

		if err != nil {
			return nil, err
		}
	}

	militaryBranch.Image = militaryBranch.ImageURL
	militaryBranch.Video = militaryBranch.VideoURL

	return &militaryBranch, nil
}
