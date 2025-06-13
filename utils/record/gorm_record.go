package record

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IDGetter interface {
	GetID() int64
}

type GormRecord[T IDGetter] struct {
	DB            *gorm.DB
	Redis         *redis.Client
	RedisKey      string
	StartRowID    int64                      // 初始 row_id
	SelectFields  []string                   // 指定查询列
	QueryModifier func(db *gorm.DB) *gorm.DB // 查询条件修改器
}

func (f *GormRecord[T]) Fetch(ctx context.Context, limit int) ([]T, error) {
	lastID, err := f.getLastID(ctx)
	if err != nil {
		return nil, err
	}

	// 构建基础查询
	tx := f.DB.WithContext(ctx).
		Where("id > ?", lastID).
		Order("id ASC").
		Limit(limit)

	// 添加指定列
	if len(f.SelectFields) > 0 {
		tx = tx.Select(f.SelectFields)
	}

	// 应用自定义查询条件
	if f.QueryModifier != nil {
		tx = f.QueryModifier(tx)
	}

	// 执行查询
	var results []T
	if err := tx.Find(&results).Error; err != nil {
		return nil, err
	}

	// 更新游标
	if len(results) > 0 {
		newID := results[len(results)-1].GetID()
		if newID > lastID {
			_ = f.Redis.Set(ctx, f.RedisKey, strconv.FormatInt(newID, 10), 144*time.Hour).Err()
		}
	}

	// 数据不足时安全重置游标（使用事务确保一致性）
	if len(results) < limit {
		err = f.Redis.Set(ctx, f.RedisKey, strconv.FormatInt(f.StartRowID, 10), 144*time.Hour).Err()
		if err != nil {
			return nil, fmt.Errorf("reset redis cursor failed: %w", err)
		}
	}

	return results, nil
}

func (f *GormRecord[T]) getLastID(ctx context.Context) (int64, error) {
	val, err := f.Redis.Get(ctx, f.RedisKey).Result()
	switch {
	case errors.Is(err, redis.Nil):
		if err := f.Redis.Set(ctx, f.RedisKey, strconv.FormatInt(f.StartRowID, 10), 144*time.Hour).Err(); err != nil {
			return 0, fmt.Errorf("init redis failed: %w", err)
		}
		return f.StartRowID, nil
	case err != nil:
		return 0, fmt.Errorf("redis error: %w", err)
	default:
		return strconv.ParseInt(val, 10, 64)
	}
}
