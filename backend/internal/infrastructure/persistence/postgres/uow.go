// 浣滅敤锛氥€愭牳蹇冨疄鐜般€戝疄鐜?application.common.UnitOfWork (鍖呭惈 txKey 涓?GORM 浜嬪姟鎺у埗)锛屼繚璇佸熀纭€璁炬柦灞傜殑鏁版嵁寮轰竴鑷存€?
package postgres

import (
	"context"
	"go-clean-architecture/internal/application/common"

	"gorm.io/gorm"
)

// 灏嗗師鏈叕寮€鐨?txKey 鏀逛负灏忓啓锛屼粎鍦ㄥ熀纭€璁炬柦灞傚唴閮ㄤ娇鐢?
type txKey struct{}

// gormUoW 瀹炵幇 application.common.UnitOfWork 鎺ュ彛
type gormUoW struct {
	db *gorm.DB
}

// NewGormUoW 閫氳繃渚濊禆娉ㄥ叆鎻愪緵瀹炰緥
func NewGormUoW(db *gorm.DB) common.UnitOfWork {
	return &gormUoW{db: db}
}

// Do 鎵ц浜嬪姟锛岀鍚?common.UnitOfWork 鎺ュ彛
func (u *gormUoW) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	// 浣跨敤娉ㄥ叆鐨勫疄渚?db锛岃€屼笉鏄叏灞€ db
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

// GetDB 鎻愪緵缁欏綋鍓嶅寘涓嬬殑鍚勪釜 Repository 瀹炵幇浣跨敤
func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	// 濡傛灉娌℃湁寮€鍚簨鍔★紝灏变娇鐢ㄩ粯璁ょ殑 DB 瀹炰緥
	return defaultDB.WithContext(ctx)
}
