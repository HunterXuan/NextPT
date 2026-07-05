package iam

import (
	"context"
	"fmt"
	"strings"

	"server/internal/dao"
	"server/internal/model/do"
	"server/internal/model/entity"
	"server/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sIamUserDomain struct{}

func init() {
	service.RegisterIamUserDomain(NewIamUserDomain())
}

func NewIamUserDomain() *sIamUserDomain {
	return &sIamUserDomain{}
}

func (s *sIamUserDomain) GetUserByLogin(ctx context.Context, login string) (*entity.IamUser, error) {
	var user *entity.IamUser
	err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Username, login).Scan(&user)
	if err != nil {
		return nil, err
	}
	if user == nil {
		err = dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Email, login).Scan(&user)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (s *sIamUserDomain) GetUserById(ctx context.Context, id uint64) (*entity.IamUser, error) {
	var user *entity.IamUser
	err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Id, id).Scan(&user)
	return user, err
}

func (s *sIamUserDomain) GetUserPasswordHash(ctx context.Context, userId uint64) (string, error) {
	columns := dao.IamUser.Columns()
	value, err := dao.IamUser.Ctx(ctx).Fields(columns.PasswordHash).Where(columns.Id, userId).Value()
	if err != nil || value == nil {
		return "", err
	}
	return value.String(), nil
}

func (s *sIamUserDomain) GetUserProfile(ctx context.Context, userId uint64) (*entity.IamUserProfile, error) {
	var profile *entity.IamUserProfile
	err := dao.IamUserProfile.Ctx(ctx).Where(dao.IamUserProfile.Columns().UserId, userId).Scan(&profile)
	return profile, err
}

func (s *sIamUserDomain) GetUserStat(ctx context.Context, userId uint64) (*entity.IamUserStat, error) {
	var stat *entity.IamUserStat
	err := dao.IamUserStat.Ctx(ctx).Where(dao.IamUserStat.Columns().UserId, userId).Scan(&stat)
	return stat, err
}

func (s *sIamUserDomain) CheckUsernameExists(ctx context.Context, username string) (bool, error) {
	count, err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Username, username).Count()
	return count > 0, err
}

func (s *sIamUserDomain) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	count, err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Email, email).Count()
	return count > 0, err
}

func (s *sIamUserDomain) InsertUser(ctx context.Context, data do.IamUser) (uint64, error) {
	id, err := dao.IamUser.Ctx(ctx).Data(data).InsertAndGetId()
	return uint64(id), err
}

func (s *sIamUserDomain) InsertUserProfile(ctx context.Context, data do.IamUserProfile) error {
	_, err := dao.IamUserProfile.Ctx(ctx).Data(data).Insert()
	return err
}

func (s *sIamUserDomain) InsertUserStat(ctx context.Context, data do.IamUserStat) error {
	_, err := dao.IamUserStat.Ctx(ctx).Data(data).Insert()
	return err
}

func (s *sIamUserDomain) UpdateUserProfile(ctx context.Context, userId uint64, avatar string, info string, signature string) error {
	columns := dao.IamUserProfile.Columns()
	_, err := dao.IamUserProfile.Ctx(ctx).Data(g.Map{
		columns.UserId:    userId,
		columns.Avatar:    avatar,
		columns.Info:      info,
		columns.Signature: signature,
	}).Save()
	return err
}

func (s *sIamUserDomain) UpdatePasswordHash(ctx context.Context, userId uint64, passwordHash string) error {
	columns := dao.IamUser.Columns()
	_, err := dao.IamUser.Ctx(ctx).
		Where(columns.Id, userId).
		Data(g.Map{
			columns.PasswordHash: passwordHash,
		}).
		Update()
	return err
}

func (s *sIamUserDomain) UpdatePasskey(ctx context.Context, userId uint64, passkey string) error {
	columns := dao.IamUser.Columns()
	_, err := dao.IamUser.Ctx(ctx).
		Where(columns.Id, userId).
		Data(g.Map{
			columns.Passkey: passkey,
		}).
		Update()
	return err
}

// GetUserByPasskey 通过 Passkey 获取用户（无缓存，纯领域逻辑）
func (s *sIamUserDomain) GetUserByPasskey(ctx context.Context, passkey string) (*entity.IamUser, error) {
	var user *entity.IamUser
	err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Passkey, passkey).Scan(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *sIamUserDomain) AdminListUsers(ctx context.Context, search string, order string, page int, size int) ([]*entity.IamUser, int, error) {
	m := dao.IamUser.Ctx(ctx)
	if search != "" {
		m = m.Where(m.Builder().
			WhereLike(dao.IamUser.Columns().Username, "%"+search+"%").
			WhereOrLike(dao.IamUser.Columns().Email, "%"+search+"%"))
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, err
	}
	var users []*entity.IamUser
	if total > 0 {
		m = m.Order(s.normalizeAdminUserOrder(order))
		err = m.Page(page, size).Scan(&users)
		if err != nil {
			return nil, 0, err
		}
	}
	return users, total, nil
}

func (s *sIamUserDomain) normalizeAdminUserOrder(order string) string {
	columns := dao.IamUser.Columns()
	defaultOrder := columns.Id + " DESC"
	fields := strings.Fields(order)
	if len(fields) == 0 {
		return defaultOrder
	}
	if len(fields) > 2 {
		return defaultOrder
	}

	column := s.adminUserOrderColumn(fields[0])
	if column == "" {
		return defaultOrder
	}

	direction := "DESC"
	if len(fields) == 2 {
		switch strings.ToLower(fields[1]) {
		case "asc":
			direction = "ASC"
		case "desc":
			direction = "DESC"
		default:
			return defaultOrder
		}
	}

	return column + " " + direction
}

func (s *sIamUserDomain) adminUserOrderColumn(field string) string {
	columns := dao.IamUser.Columns()
	switch strings.ToLower(strings.ReplaceAll(field, "_", "")) {
	case "id":
		return columns.Id
	case "username":
		return columns.Username
	case "email":
		return columns.Email
	case "status":
		return columns.Status
	case "role":
		return columns.Role
	case "createdat":
		return columns.CreatedAt
	case "updatedat":
		return columns.UpdatedAt
	case "lastlogin":
		return columns.LastLogin
	case "vipuntil":
		return columns.VipUntil
	default:
		return ""
	}
}

func (s *sIamUserDomain) AdminUpdateUser(ctx context.Context, id uint64, status *int, role *uint, passkey *string) error {
	data := do.IamUser{}
	if status != nil {
		data.Status = *status
	}
	if role != nil {
		data.Role = *role
	}
	if passkey != nil {
		data.Passkey = *passkey
	}
	_, err := dao.IamUser.Ctx(ctx).Where(dao.IamUser.Columns().Id, id).Data(data).Update()
	return err
}

func (s *sIamUserDomain) AdminGetUserStat(ctx context.Context, id uint64) (*entity.IamUserStat, error) {
	var stat *entity.IamUserStat
	err := dao.IamUserStat.Ctx(ctx).Where(dao.IamUserStat.Columns().UserId, id).Scan(&stat)
	return stat, err
}

func (s *sIamUserDomain) AdminUpdateUserStat(ctx context.Context, id uint64, uploadedDiff *int64, downloadedDiff *int64) (int64, error) {
	data := g.Map{}
	if uploadedDiff != nil && *uploadedDiff != 0 {
		diff := *uploadedDiff
		if diff > 0 {
			data[dao.IamUserStat.Columns().Uploaded] = gdb.Raw(fmt.Sprintf("uploaded + %d", diff))
		} else {
			data[dao.IamUserStat.Columns().Uploaded] = gdb.Raw(fmt.Sprintf("GREATEST(0, CAST(uploaded AS SIGNED) + %d)", diff))
		}
	}
	if downloadedDiff != nil && *downloadedDiff != 0 {
		diff := *downloadedDiff
		if diff > 0 {
			data[dao.IamUserStat.Columns().Downloaded] = gdb.Raw(fmt.Sprintf("downloaded + %d", diff))
		} else {
			data[dao.IamUserStat.Columns().Downloaded] = gdb.Raw(fmt.Sprintf("GREATEST(0, CAST(downloaded AS SIGNED) + %d)", diff))
		}
	}
	if len(data) == 0 {
		return 0, nil
	}
	res, err := dao.IamUserStat.Ctx(ctx).Where(dao.IamUserStat.Columns().UserId, id).Data(data).Update()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *sIamUserDomain) GetUsersByIds(ctx context.Context, ids []uint64) ([]entity.IamUser, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var users []entity.IamUser
	err := dao.IamUser.Ctx(ctx).WhereIn(dao.IamUser.Columns().Id, ids).Scan(&users)
	return users, err
}

func (s *sIamUserDomain) GetUserIdsByRoles(ctx context.Context, roleIds []uint) ([]uint64, error) {
	if len(roleIds) == 0 {
		return nil, nil
	}
	var users []entity.IamUser
	err := dao.IamUser.Ctx(ctx).
		Fields(dao.IamUser.Columns().Id).
		WhereIn(dao.IamUser.Columns().Role, roleIds).
		OrderAsc(dao.IamUser.Columns().Id).
		Scan(&users)
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(users))
	for _, user := range users {
		if user.Id == 0 {
			continue
		}
		ids = append(ids, user.Id)
	}
	return ids, nil
}

func (s *sIamUserDomain) GetUserProfilesByUserIds(ctx context.Context, userIds []uint64) ([]entity.IamUserProfile, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	var profiles []entity.IamUserProfile
	err := dao.IamUserProfile.Ctx(ctx).WhereIn(dao.IamUserProfile.Columns().UserId, userIds).Scan(&profiles)
	return profiles, err
}
