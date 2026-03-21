package mysql

import (
	"github.com/boxuanduan/gochat/pkg/model"
	"gorm.io/gorm"
)

type GroupRepo struct {
	db *gorm.DB
}

func NewGroupRepo(db *gorm.DB) *GroupRepo {
	return &GroupRepo{db: db}
}

// Create a new group
func (r *GroupRepo) CreateGroup(info *model.GroupInfo) error {
	tx := r.db.Create(info)
	return tx.Error
}

// Add a member to the group
func (r *GroupRepo) AddMember(member *model.GroupMember) error {
	tx := r.db.Create(member)
	return tx.Error
}

// List members' id of a group
func (r *GroupRepo) GetMemberIDs(groupID int64) ([]int64, error) {
	var memberIDs []int64
	err := r.db.Model(&model.GroupMember{}).
		Where("group_id = ?", groupID).
		Pluck("user_id", &memberIDs).
		Error

	if err != nil {
		return nil, err
	}

	return memberIDs, nil
}

func (r *GroupRepo) GetGroupByID(groupID int64) (*model.GroupInfo, error) {
	var group model.GroupInfo
	err := r.db.Where("id = ?", groupID).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepo) IsMember(groupID int64, userID int64) (bool, error) {
	var member model.GroupMember
	err := r.db.Model(&model.GroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&member).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil // User is not a member of the group
		}
		return false, err // An error occurred while checking membership
	}
	return true, nil // User is a member of the group
}
