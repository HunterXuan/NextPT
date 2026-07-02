package economy

import (
	"context"
	"fmt"

	"server/internal/dao"
	"server/internal/model/entity"
	"server/internal/model/out/economyout"
	"server/internal/service"
)

type sEconomyRewardDomain struct{}

func init() {
	service.RegisterEconomyRewardDomain(NewEconomyRewardDomain())
}

func NewEconomyRewardDomain() *sEconomyRewardDomain {
	return &sEconomyRewardDomain{}
}

func (s *sEconomyRewardDomain) InsertRewardRecord(ctx context.Context, record entity.EconomyRewardRecord) error {
	_, err := dao.EconomyRewardRecord.Ctx(ctx).Data(record).Insert()
	return err
}

func (s *sEconomyRewardDomain) QueryRewardSummaries(ctx context.Context, targetType string, targetId uint64, page, size int) ([]economyout.RewardSummary, int, error) {
	columns := dao.EconomyRewardRecord.Columns()

	totalValue, err := dao.EconomyRewardRecord.Ctx(ctx).
		Fields(fmt.Sprintf("COUNT(DISTINCT %s)", columns.FromUserId)).
		Where(columns.TargetType, targetType).
		Where(columns.TargetId, targetId).
		Value()
	if err != nil {
		return nil, 0, err
	}

	var summaries []economyout.RewardSummary
	err = dao.EconomyRewardRecord.Ctx(ctx).
		Fields(
			fmt.Sprintf("%s AS user_id", columns.FromUserId),
			fmt.Sprintf("SUM(%s) AS amount", columns.Amount),
			"COUNT(*) AS reward_count",
			fmt.Sprintf("MAX(%s) AS last_reward_at", columns.CreatedAt),
		).
		Where(columns.TargetType, targetType).
		Where(columns.TargetId, targetId).
		Group(columns.FromUserId).
		Page(page, size).
		OrderDesc("amount").
		OrderDesc("last_reward_at").
		Scan(&summaries)
	return summaries, totalValue.Int(), err
}

func (s *sEconomyRewardDomain) DeleteRewardRecordsByTarget(ctx context.Context, targetType string, targetId uint64) error {
	_, err := dao.EconomyRewardRecord.Ctx(ctx).
		Where(dao.EconomyRewardRecord.Columns().TargetType, targetType).
		Where(dao.EconomyRewardRecord.Columns().TargetId, targetId).
		Delete()
	return err
}

func (s *sEconomyRewardDomain) DeleteRewardRecordsByTargets(ctx context.Context, targetType string, targetIds []uint64) error {
	if len(targetIds) == 0 {
		return nil
	}
	_, err := dao.EconomyRewardRecord.Ctx(ctx).
		Where(dao.EconomyRewardRecord.Columns().TargetType, targetType).
		WhereIn(dao.EconomyRewardRecord.Columns().TargetId, targetIds).
		Delete()
	return err
}
