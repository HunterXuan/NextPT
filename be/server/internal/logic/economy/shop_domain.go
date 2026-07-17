package economy

import (
	"context"
	"sort"

	"server/internal/consts"
	"server/internal/dao"
	"server/internal/model"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gtime"
)

type sEconomyShopDomain struct{}

func init() {
	service.RegisterEconomyShopDomain(NewEconomyShopDomain())
}

func NewEconomyShopDomain() *sEconomyShopDomain {
	return &sEconomyShopDomain{}
}

func (s *sEconomyShopDomain) LoadProducts(ctx context.Context) (model.EconomyShopProducts, error) {
	value := service.SiteConfigDomain().GetByPath(ctx, consts.SiteConfigEconomyShopProducts)
	var products model.EconomyShopProducts
	if err := value.Scan(&products); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.invalid_config"))
	}
	products = products.Normalized()
	if err := products.Validate(); err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.invalid_config"))
	}
	sort.SliceStable(products, func(i, j int) bool {
		if products[i].SortOrder == products[j].SortOrder {
			return products[i].Key < products[j].Key
		}
		return products[i].SortOrder < products[j].SortOrder
	})
	return products, nil
}

func (s *sEconomyShopDomain) InsertOrder(ctx context.Context, order entity.EconomyShopOrder) (uint64, error) {
	result, err := dao.EconomyShopOrder.Ctx(ctx).Data(do.EconomyShopOrder{
		UserId:          order.UserId,
		ProductKey:      order.ProductKey,
		ProductType:     order.ProductType,
		ProductSnapshot: order.ProductSnapshot,
		Price:           order.Price,
		Status:          order.Status,
		TargetType:      order.TargetType,
		TargetId:        order.TargetId,
		CreatedAt:       order.CreatedAt,
		CompletedAt:     order.CompletedAt,
	}).Insert()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	return uint64(id), err
}

func (s *sEconomyShopDomain) CompleteOrder(ctx context.Context, id uint64, targetType string, targetId uint64) error {
	columns := dao.EconomyShopOrder.Columns()
	_, err := dao.EconomyShopOrder.Ctx(ctx).
		Where(columns.Id, id).
		Data(do.EconomyShopOrder{
			Status:      consts.EconomyShopOrderStatusCompleted,
			TargetType:  targetType,
			TargetId:    targetId,
			CompletedAt: gtime.Now(),
		}).Update()
	return err
}

func (s *sEconomyShopDomain) QueryOrdersByUser(ctx context.Context, userId uint64, page int, size int) ([]entity.EconomyShopOrder, int, error) {
	columns := dao.EconomyShopOrder.Columns()
	m := dao.EconomyShopOrder.Ctx(ctx).Where(columns.UserId, userId)
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []entity.EconomyShopOrder
	err = m.Page(page, size).OrderDesc(columns.Id).Scan(&list)
	return list, total, err
}
