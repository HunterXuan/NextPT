package economy

import (
	"context"
	"math"
	"time"

	"server/internal/consts"
	"server/internal/model"
	"server/internal/model/entity"
	"server/internal/model/in/economyin"
	"server/internal/model/out/economyout"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

type sEconomyShopUsecase struct{}

func init() {
	service.RegisterEconomyShopUsecase(NewEconomyShopUsecase())
}

func NewEconomyShopUsecase() *sEconomyShopUsecase {
	return &sEconomyShopUsecase{}
}

func (s *sEconomyShopUsecase) ListProducts(ctx context.Context, actor *model.Actor, in economyin.ShopProductListInp) (*economyout.ShopProductListOut, error) {
	products, err := s.loadProductsCache(ctx)
	if err != nil {
		return nil, err
	}
	return &economyout.ShopProductListOut{List: products.Enabled()}, nil
}

func (s *sEconomyShopUsecase) CreateOrder(ctx context.Context, actor *model.Actor, in economyin.ShopOrderCreateInp) (*economyout.ShopOrderCreateOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	products, err := s.loadProductsCache(ctx)
	if err != nil {
		return nil, err
	}
	product := products.Find(in.ProductKey)
	if product == nil {
		return nil, gerror.New(gi18n.T(ctx, "economy.shop.product_not_found"))
	}
	if !product.Enabled {
		return nil, gerror.New(gi18n.T(ctx, "economy.shop.product_disabled"))
	}

	price := math.Round(product.Price*100) / 100
	var out economyout.ShopOrderCreateOut
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		now := gtime.Now()
		orderId, err := service.EconomyShopDomain().InsertOrder(ctx, entity.EconomyShopOrder{
			UserId:          actor.Id,
			ProductKey:      product.Key,
			ProductType:     product.Type,
			ProductSnapshot: gjson.New(product),
			Price:           price,
			Status:          consts.EconomyShopOrderStatusPending,
			CreatedAt:       now,
		})
		if err != nil {
			return gerror.Wrap(err, gi18n.T(ctx, "economy.shop.create_order_failed"))
		}

		if err = service.EconomyBonusDomain().DebitBonusIfEnough(ctx, actor.Id, price); err != nil {
			return err
		}
		targetType, targetId, err := s.fulfillProduct(ctx, actor.Id, *product)
		if err != nil {
			return err
		}
		if err = service.EconomyShopDomain().CompleteOrder(ctx, orderId, targetType, targetId); err != nil {
			return gerror.Wrap(err, gi18n.T(ctx, "economy.shop.complete_order_failed"))
		}
		balanceAfter, err := service.EconomyBonusDomain().GetUserBonus(ctx, actor.Id)
		if err != nil {
			return err
		}
		if err = service.EconomyBonusDomain().InsertBonusLog(ctx, entity.EconomyBonusLog{
			UserId:       actor.Id,
			Amount:       -price,
			BalanceAfter: balanceAfter,
			Action:       consts.EconomyBonusActionShopPurchase,
			TargetType:   consts.EconomyBonusTargetTypeShopOrder,
			TargetId:     orderId,
		}); err != nil {
			return err
		}
		out = economyout.ShopOrderCreateOut{
			OrderId:      orderId,
			ProductKey:   product.Key,
			Price:        price,
			TargetType:   targetType,
			TargetId:     targetId,
			BalanceAfter: balanceAfter,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	service.IamUserUsecase().InvalidateUserCache(ctx, actor.Id)
	return &out, nil
}

func (s *sEconomyShopUsecase) loadProductsCache(ctx context.Context) (model.EconomyShopProducts, error) {
	cacheKey := service.SysCache().KeySiteConfigFullPath(ctx, consts.SiteConfigEconomyShopProducts)
	value, err := gcache.GetOrSetFunc(ctx, cacheKey, func(ctx context.Context) (any, error) {
		return service.EconomyShopDomain().LoadProducts(ctx)
	}, 10*time.Minute)
	if err != nil {
		return nil, err
	}
	if value.IsNil() {
		return service.EconomyShopDomain().LoadProducts(ctx)
	}
	if products, ok := value.Val().(model.EconomyShopProducts); ok {
		return products, nil
	}
	return service.EconomyShopDomain().LoadProducts(ctx)
}

func (s *sEconomyShopUsecase) ListMyOrders(ctx context.Context, actor *model.Actor, in economyin.ShopOrderListInp) (*economyout.ShopOrderListOut, error) {
	if actor == nil {
		return nil, gerror.New(gi18n.T(ctx, "iam.general.unauthorized"))
	}
	page, size := in.Page, in.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	orders, total, err := service.EconomyShopDomain().QueryOrdersByUser(ctx, actor.Id, page, size)
	if err != nil {
		return nil, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.query_orders_failed"))
	}
	list := make([]economyout.ShopOrderItem, 0, len(orders))
	for _, order := range orders {
		product := model.NewEconomyShopProductSnapshot(order.ProductKey, order.ProductType, order.Price)
		if order.ProductSnapshot != nil {
			var snapshot model.EconomyShopProductConfig
			if err = order.ProductSnapshot.Scan(&snapshot); err == nil && snapshot.Key != "" && snapshot.Type != "" {
				product = snapshot
			}
		}
		list = append(list, economyout.ShopOrderItem{
			Id:          order.Id,
			Product:     product,
			Price:       order.Price,
			Status:      order.Status,
			TargetType:  order.TargetType,
			TargetId:    order.TargetId,
			CreatedAt:   order.CreatedAt,
			CompletedAt: order.CompletedAt,
		})
	}
	return &economyout.ShopOrderListOut{Page: page, Size: size, Total: total, List: list}, nil
}

func (s *sEconomyShopUsecase) fulfillProduct(ctx context.Context, userId uint64, product model.EconomyShopProductConfig) (string, uint64, error) {
	switch product.Type {
	case consts.EconomyShopProductTypeInvite:
		return s.fulfillInvite(ctx, userId, product)
	case consts.EconomyShopProductTypeVip:
		return s.fulfillVip(ctx, userId, product)
	case consts.EconomyShopProductTypeUpload:
		return s.fulfillUpload(ctx, userId, product)
	case consts.EconomyShopProductTypeDownload:
		return s.fulfillDownload(ctx, userId, product)
	default:
		return "", 0, gerror.New(gi18n.T(ctx, "economy.shop.product_not_implemented"))
	}
}

func (s *sEconomyShopUsecase) fulfillInvite(ctx context.Context, userId uint64, product model.EconomyShopProductConfig) (string, uint64, error) {
	amount := product.InviteAmount()
	var firstInviteId uint64
	for i := 0; i < amount; i++ {
		inviteId, err := service.IamInviteDomain().CreateInvite(ctx, userId, false, nil)
		if err != nil {
			return "", 0, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.fulfill_failed"))
		}
		if firstInviteId == 0 {
			firstInviteId = inviteId
		}
	}
	return consts.EconomyShopOrderTargetTypeIamInvite, firstInviteId, nil
}

func (s *sEconomyShopUsecase) fulfillVip(ctx context.Context, userId uint64, product model.EconomyShopProductConfig) (string, uint64, error) {
	if err := service.IamUserDomain().ExtendUserVip(ctx, userId, product.VipDurationDays(), consts.EconomyShopVipRemark); err != nil {
		return "", 0, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.fulfill_failed"))
	}
	return consts.EconomyShopOrderTargetTypeIamUser, userId, nil
}

func (s *sEconomyShopUsecase) fulfillUpload(ctx context.Context, userId uint64, product model.EconomyShopProductConfig) (string, uint64, error) {
	if err := service.IamUserDomain().AddUserUploaded(ctx, userId, product.TrafficBytes()); err != nil {
		return "", 0, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.fulfill_failed"))
	}
	return consts.EconomyShopOrderTargetTypeIamUserStat, userId, nil
}

func (s *sEconomyShopUsecase) fulfillDownload(ctx context.Context, userId uint64, product model.EconomyShopProductConfig) (string, uint64, error) {
	stat, err := service.IamUserDomain().GetUserStat(ctx, userId)
	if err != nil {
		return "", 0, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.fulfill_failed"))
	}
	if stat == nil {
		return "", 0, gerror.New(gi18n.T(ctx, "economy.shop.fulfill_failed"))
	}
	if stat.Downloaded < product.TrafficBytes() {
		return "", 0, gerror.New(gi18n.T(ctx, "economy.shop.insufficient_download"))
	}
	if err := service.IamUserDomain().ReduceUserDownloaded(ctx, userId, product.TrafficBytes()); err != nil {
		return "", 0, gerror.Wrap(err, gi18n.T(ctx, "economy.shop.fulfill_failed"))
	}
	return consts.EconomyShopOrderTargetTypeIamUserStat, userId, nil
}
