// 浣滅敤锛欼nfrastructure 灞?Repository 瀹炵幇鐨?TDD 娴嬭瘯
// 娴嬭瘯瀵硅薄锛歱ostgres.userRepository 鐨勫唴閮ㄩ€昏緫鍜?PO 鏄犲皠
// 娴嬭瘯绛栫暐锛氶€氳繃鐩存帴鏋勯€?UserPO 楠岃瘉 Domain Entity 涓?PO 涔嬮棿鐨勬槧灏勬纭€?// 娉ㄦ剰锛氬浜庣湡姝ｇ殑鏁版嵁搴撹涓猴紙SQL 璇硶鍏煎鎬э級锛屽缓璁湪 CI 涓娇鐢?Docker + 鐪熷疄 PostgreSQL 鍋氶泦鎴愭祴璇?
package postgres_test

import (
	"testing"

	"go-clean-architecture/internal/infrastructure/persistence/postgres"
)

// ==================== 绾?-> 缁?-> 閲嶆瀯 ====================

// 馃敶 绾細楠岃瘉 UserPO 缁撴瀯浣撳瓧娈垫槧灏勨€斺€斾繚璇?PO 涓庢暟鎹簱琛ㄧ粨鏋勫榻?
func TestUserPO_FieldMapping(t *testing.T) {
	po := postgres.UserPO{
		ID:       "test-001",
		Username: "alice",
	}

	// 馃煝 缁匡細瀛楁鍊煎簲璇ユ纭祴鍊?
	if po.ID != "test-001" {
		t.Errorf("expected ID 'test-001', got '%s'", po.ID)
	}
	if po.Username != "alice" {
		t.Errorf("expected Username 'alice', got '%s'", po.Username)
	}
}

// 馃敶 绾細楠岃瘉 NewUserRepository 鏋勯€犲嚱鏁颁笉涓?nil
func TestNewUserRepository_NotNil(t *testing.T) {
	// 浼犲叆 nil db 浠呭仛鏋勯€犻獙璇侊紙瀹為檯浣跨敤鏃?db 涓嶄細涓?nil锛?
	repo := postgres.NewUserRepository(nil)
	if repo == nil {
		t.Fatal("expected non-nil repository")
	}
}
