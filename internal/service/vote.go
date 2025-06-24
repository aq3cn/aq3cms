package service

import (
	"fmt"
	"time"

	"aq3cms/config"
	"aq3cms/internal/model"
	"aq3cms/pkg/cache"
	"aq3cms/pkg/database"
	"aq3cms/pkg/logger"
)

// VoteService 投票服务
type VoteService struct {
	db             *database.DB
	cache          cache.Cache
	config         *config.Config
	voteModel      *model.VoteModel
	voteOptionModel *model.VoteOptionModel
	voteLogModel   *model.VoteLogModel
}

// NewVoteService 创建投票服务
func NewVoteService(db *database.DB, cache cache.Cache, config *config.Config) *VoteService {
	return &VoteService{
		db:             db,
		cache:          cache,
		config:         config,
		voteModel:      model.NewVoteModel(db),
		voteOptionModel: model.NewVoteOptionModel(db),
		voteLogModel:   model.NewVoteLogModel(db),
	}
}

// GetVote 获取投票
func (s *VoteService) GetVote(id int64) (*model.Vote, []*model.VoteOption, error) {
	// 获取投票
	vote, err := s.voteModel.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	// 获取投票选项
	options, err := s.voteOptionModel.GetByVoteID(id)
	if err != nil {
		return nil, nil, err
	}

	return vote, options, nil
}

// CreateVote 创建投票
func (s *VoteService) CreateVote(vote *model.Vote, options []*model.VoteOption) (int64, error) {
	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return 0, err
	}
	defer tx.Rollback()

	// 创建投票
	voteID, err := s.voteModel.Create(vote)
	if err != nil {
		return 0, err
	}

	// 创建投票选项
	for i, option := range options {
		option.VoteID = voteID
		option.OrderID = i + 1
		_, err = s.voteOptionModel.Create(option)
		if err != nil {
			return 0, err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return 0, err
	}

	return voteID, nil
}

// UpdateVote 更新投票
func (s *VoteService) UpdateVote(vote *model.Vote, options []*model.VoteOption) error {
	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 更新投票
	err = s.voteModel.Update(vote)
	if err != nil {
		return err
	}

	// 删除原有选项
	err = s.voteOptionModel.DeleteByVoteID(vote.ID)
	if err != nil {
		return err
	}

	// 创建新选项
	for i, option := range options {
		option.VoteID = vote.ID
		option.OrderID = i + 1
		_, err = s.voteOptionModel.Create(option)
		if err != nil {
			return err
		}
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// DeleteVote 删除投票
func (s *VoteService) DeleteVote(id int64) error {
	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 删除投票日志
	err = s.voteLogModel.DeleteByVoteID(id)
	if err != nil {
		return err
	}

	// 删除投票选项
	err = s.voteOptionModel.DeleteByVoteID(id)
	if err != nil {
		return err
	}

	// 删除投票
	err = s.voteModel.Delete(id)
	if err != nil {
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// DoVote 执行投票
func (s *VoteService) DoVote(voteID int64, optionIDs []int64, memberID int64, ip string) error {
	// 获取投票
	vote, err := s.voteModel.GetByID(voteID)
	if err != nil {
		return err
	}

	// 检查投票状态
	if vote.Status != 1 {
		return fmt.Errorf("vote is disabled")
	}

	// 检查投票时间
	now := time.Now()
	if now.Before(vote.StartTime) {
		return fmt.Errorf("vote has not started")
	}
	if now.After(vote.EndTime) {
		return fmt.Errorf("vote has ended")
	}

	// 检查是否已投票
	voted, err := s.voteLogModel.CheckVoted(voteID, memberID, ip)
	if err != nil {
		return err
	}
	if voted {
		return fmt.Errorf("already voted")
	}

	// 检查选项数量
	if len(optionIDs) == 0 {
		return fmt.Errorf("no option selected")
	}
	if vote.IsMulti == 0 && len(optionIDs) > 1 {
		return fmt.Errorf("multiple options not allowed")
	}
	if vote.IsMulti == 1 && len(optionIDs) > vote.MaxChoices {
		return fmt.Errorf("too many options selected")
	}

	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		logger.Error("开始事务失败", "error", err)
		return err
	}
	defer tx.Rollback()

	// 增加选项投票数
	for _, optionID := range optionIDs {
		// 检查选项是否存在
		option, err := s.voteOptionModel.GetByID(optionID)
		if err != nil {
			return err
		}
		if option.VoteID != voteID {
			return fmt.Errorf("option not belong to vote")
		}

		// 增加投票数
		err = s.voteOptionModel.IncrementCount(optionID)
		if err != nil {
			return err
		}

		// 创建投票日志
		log := &model.VoteLog{
			VoteID:   voteID,
			OptionID: optionID,
			MemberID: memberID,
			IP:       ip,
		}
		_, err = s.voteLogModel.Create(log)
		if err != nil {
			return err
		}
	}

	// 增加总投票数
	err = s.voteModel.IncrementTotalCount(voteID, 1)
	if err != nil {
		return err
	}

	// 提交事务
	err = tx.Commit()
	if err != nil {
		logger.Error("提交事务失败", "error", err)
		return err
	}

	return nil
}

// GetVoteResult 获取投票结果
func (s *VoteService) GetVoteResult(id int64) (*model.Vote, []*model.VoteOption, error) {
	// 获取投票
	vote, err := s.voteModel.GetByID(id)
	if err != nil {
		return nil, nil, err
	}

	// 获取投票选项
	options, err := s.voteOptionModel.GetByVoteID(id)
	if err != nil {
		return nil, nil, err
	}

	return vote, options, nil
}

// GetVoteLogs 获取投票日志
func (s *VoteService) GetVoteLogs(voteID int64, page, pageSize int) ([]*model.VoteLog, int, error) {
	// 获取投票日志
	logs, total, err := s.voteLogModel.GetByVoteID(voteID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
