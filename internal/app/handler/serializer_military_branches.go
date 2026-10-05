package handler

import "army/internal/app/ds"

type MilitaryBranchSerializer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ImageURL    string `json:"image"`
	VideoURL    string `json:"video"`
	SpeedPlain  int    `json:"speed_plain"`
	Food        int    `json:"food"`
	CreatorID   int64  `json:"creator_id"`
	CreatedAt   string `json:"created_at"`
	FormedAt    string `json:"formed_at"`
}

func NewMilitaryBranchSerializer(
	militaryBranch ds.MilitaryBranch,
) MilitaryBranchSerializer {
	return MilitaryBranchSerializer{
		ID:          militaryBranch.ID,
		Name:        militaryBranch.Name,
		Description: militaryBranch.Description,
		Status:      militaryBranch.Status,
		ImageURL:    militaryBranch.ImageURL,
		VideoURL:    militaryBranch.VideoURL,
		SpeedPlain:  militaryBranch.SpeedPlain,
		Food:        militaryBranch.Food,
		CreatorID:   militaryBranch.CreatorID,
		CreatedAt:   militaryBranch.CreatedAt.Format("2006-01-02 15:04:05"),
		FormedAt:    militaryBranch.FormedAt.Format("2006-01-02 15:04:05"),
	}
}
