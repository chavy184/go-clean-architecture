// 浣滅敤锛氬疄鐜?domain.user.Repository锛岃礋璐ｅ叿浣撶殑鏁版嵁搴撹鍐欐搷浣滐紝骞跺皢 PO 涓庡疄浣撹繘琛岀浉浜掕浆鎹?
package postgres

import (
	"context"
	"go-clean-architecture/internal/domain/user"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet 灏嗗綋鍓嶅寘涓嬫墍鏈夌殑鍏蜂綋瀹炵幇浣滀负 Provider 鎻愪緵
var ProviderSet = wire.NewSet(NewDB, NewUserRepository, NewGormUoW)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, u *user.User) error {
	db := GetDB(ctx, r.db)
	
	po := &UserPO{
		ID:       u.ID,
		Username: u.Username,
	}
	
	return db.Save(po).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	db := GetDB(ctx, r.db)
	
	var po UserPO
	if err := db.Where("id = ?", id).First(&po).Error; err != nil {
		return nil, err
	}
	
	// Convert PO back to Domain Entity
	return &user.User{
		ID:       po.ID,
		Username: po.Username,
	}, nil
}
