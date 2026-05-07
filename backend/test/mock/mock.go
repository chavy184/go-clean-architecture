// 浣滅敤锛氭彁渚?domain 涓?application 灞傛帴鍙ｇ殑鎵嬪啓 Mock 瀹炵幇锛屼緵 TDD 鍗曞厓娴嬭瘯浣跨敤
// 鍦ㄥ疄闄呴」鐩腑锛屼篃鍙互浣跨敤 mockgen 鑷姩鐢熸垚
package mock

import (
	"context"

	"go-clean-architecture/internal/domain/user"
)

// ---- Mock Repository ----

type MockUserRepository struct {
	// 瀛樺偍妯℃嫙鏁版嵁
	Users map[string]*user.User

	// 鐢ㄤ簬鏂█璋冪敤琛屼负
	SaveCalled    bool
	SaveCalledWith *user.User
	SaveErr       error

	FindByIDCalled bool
	FindByIDErr    error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[string]*user.User),
	}
}

func (m *MockUserRepository) Save(ctx context.Context, u *user.User) error {
	m.SaveCalled = true
	m.SaveCalledWith = u
	if m.SaveErr != nil {
		return m.SaveErr
	}
	m.Users[u.ID] = u
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	m.FindByIDCalled = true
	if m.FindByIDErr != nil {
		return nil, m.FindByIDErr
	}
	u, ok := m.Users[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

// ---- Mock UnitOfWork ----

type MockUnitOfWork struct {
	DoErr error
}

func NewMockUnitOfWork() *MockUnitOfWork {
	return &MockUnitOfWork{}
}

// Do 鐩存帴鎵ц闂寘锛屾ā鎷熸棤浜嬪姟琛屼负锛涙垨娉ㄥ叆 DoErr 妯℃嫙浜嬪姟澶辫触
func (m *MockUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if m.DoErr != nil {
		return m.DoErr
	}
	return fn(ctx)
}
