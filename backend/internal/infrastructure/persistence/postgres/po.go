// 浣滅敤锛氭寔涔呭寲瀵硅薄 Persistent Object锛堝甫 GORM tag锛屼笌 Domain Entity 闅旂锛夛紝鐢ㄤ簬鐩存帴鍜屾暟鎹簱琛ㄦ槧灏?
package postgres

type UserPO struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
}
