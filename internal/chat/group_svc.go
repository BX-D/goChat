package chat

import (
	"github.com/boxuanduan/gochat/internal/repo/mysql"
	"github.com/boxuanduan/gochat/pkg/model"
)

type GroupService struct {
	groupRepo *mysql.GroupRepo
}

func NewGroupService(groupRepo *mysql.GroupRepo) *GroupService {
	return &GroupService{groupRepo: groupRepo}
}

func (s *GroupService) CreateGroup(ownerID int64, groupName string) (*model.GroupInfo, error) {
	group := &model.GroupInfo{
		OwnerID: ownerID,
		Name:    groupName,
	}
	// Create the group
	err := s.groupRepo.CreateGroup(group)
	if err != nil {
		return nil, err
	}

	// Add the owner as a member of the group
	owner := &model.GroupMember{UserID: ownerID, GroupID: group.ID, Role: 2} // Role 2 for owner, 1 for admin, 0 for regular member
	err = s.groupRepo.AddMember(owner)

	return group, err
}

func (s *GroupService) JoinGroup(groupID, userID int64) error {
	// Check if the user is already a member of the group
	isMember, err := s.groupRepo.IsMember(groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return nil // User is already a member, no need to add again
	}

	// Add the user as a member of the group // Role 1 for regular member
	err = s.groupRepo.AddMember(&model.GroupMember{GroupID: groupID, UserID: userID, Role: 0})
	return err
}

func (s GroupService) GetMemberIDs(groupID int64) ([]int64, error) {
	return s.groupRepo.GetMemberIDs(groupID)
}
