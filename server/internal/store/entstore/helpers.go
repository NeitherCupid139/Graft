package entstore

import (
	"math"

	"graft/server/internal/store"
)

// toEntID 把上层稳定的 uint64 标识转换为 Ent 当前使用的 int 主键。
//
// 超出当前主键范围时统一返回上层约定的“未命中”语义，避免插件感知底层实现限制。
func toEntID(id uint64) (int, error) {
	if id == 0 || id > math.MaxInt {
		return 0, store.ErrUserNotFound
	}

	return int(id), nil
}
