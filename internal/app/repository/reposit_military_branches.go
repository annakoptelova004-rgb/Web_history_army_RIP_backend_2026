package repository

import (
	"army/internal/app/ds"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db          *gorm.DB
	minioClient *minio.Client
	minioBucket string
}

func New(
	dsn string,
	minioEndpoint string,
	minioAccessKey string,
	minioSecretKey string,
	minioBucket string,
) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	minioClient, err := minio.New(
		minioEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
			Secure: false,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}

	return &Repository{
		db:          db,
		minioClient: minioClient,
		minioBucket: minioBucket,
	}, nil
}
func (r *Repository) GetMilitaryBranches() ([]ds.MilitaryBranch, error) {
	var militaryBranches []ds.MilitaryBranch

	err := r.db.
		Where("status = ?", "опубликован").
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
	food int,
) error {
	return r.db.Model(&ds.MilitaryBranch{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"description": description,
			"speed_plain": speedPlain,
			"food":        food,
			"status":      "опубликован",
			"formed_at":   time.Now(),
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

func (r *Repository) CreateLike(like *ds.MilitaryBranchLike) error {
	return r.db.Create(like).Error
}

func (r *Repository) DeleteLike(userID int64, militaryBranchID int64) error {
	return r.db.
		Where("user_id = ? AND military_branch_id = ?", userID, militaryBranchID).
		Delete(&ds.MilitaryBranchLike{}).Error
}

func (r *Repository) CreateUser(user *ds.User) error {
	return r.db.Create(user).Error
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

func (r *Repository) UpdateMilitaryBranchImage(id int64, fileName string) error {
	return r.db.Model(&ds.MilitaryBranch{}).
		Where("id = ?", id).
		Update("image_url", fileName).Error
}
func (r *Repository) UpdateMilitaryBranchVideo(id int64, fileName string) error {
	return r.db.Model(&ds.MilitaryBranch{}).
		Where("id = ?", id).
		Update("video_url", fileName).Error
}

func (r *Repository) UploadFile(
	file io.Reader,
	size int64,
	contentType string,
	extension string,
) (string, error) {
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)

	_, err := r.minioClient.PutObject(
		context.Background(),
		r.minioBucket,
		fileName,
		file,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return "", err
	}

	return fileName, nil
}
