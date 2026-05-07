// 浣滅敤锛氱敤鎴烽鍩熸湇鍔★紝澶勭悊涓嶅睘浜庡崟涓疄浣撶殑銆佽法瀹炰綋鐨勭敤鎴风浉鍏虫牳蹇冧笟鍔￠€昏緫
package user

import (
	"context"
	"errors"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewDomainService)

// DomainService 璐熻矗涓嶅睘浜庡崟涓€瀹炰綋鐨勬牳蹇冧笟鍔￠€昏緫锛堜緥濡傝法瀹炰綋鏍￠獙銆佸鏉傜姸鎬佹祦杞級
type DomainService struct {
	// 棰嗗煙鏈嶅姟鍙互渚濊禆 Repository 杩涜涓氬姟瑙勫垯鏍￠獙锛屼絾涓嶈鍦ㄨ繖閲岀洿鎺ヨ皟鐢ㄧ涓夋柟鍩虹璁炬柦 API
	repo Repository 
}

func NewDomainService(repo Repository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

// ValidateRegistration 绀轰緥锛氶鍩熸湇鍔′腑鐨勪笟鍔¤鍒欐牎楠?// 涓轰粈涔堜笉鏀惧湪 User 瀹炰綋涓紵鍥犱负闇€瑕佽皟鐢?repo 鏌ヨ鏄惁瀛樺湪閲嶅悕锛屽疄浣撲笉鑳戒緷璧?Repo銆?
func (s *DomainService) ValidateRegistration(ctx context.Context, username string) error {
	if len(username) < 3 {
		return errors.New("domain_error: username too short")
	}

	// 鍋囪瀛樺湪 FindByName 鏂规硶
	// existing, _ := s.repo.FindByName(ctx, username)
	// if existing != nil {
	// 
	return errors.New("domain_error: username already exists")
	// }

	return nil
}

// CalculateRiskScore 绀轰緥锛氬鏉傜殑鑱氬悎绾у埆璁＄畻閫昏緫
func (s *DomainService) CalculateRiskScore(ctx context.Context, u *User) int {
	// 妯℃嫙澶嶆潅鐨勯鍩熸墦鍒嗛€昏緫
	score := 0
	if u.Username == "admin" {
		score += 100
	}
	return score
}
