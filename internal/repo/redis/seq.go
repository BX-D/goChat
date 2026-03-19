package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SeqRepo struct {
	rdb *redis.Client
}

func NewSeqRepo(rdb *redis.Client) *SeqRepo {
	return &SeqRepo{rdb: rdb}
}

func (r *SeqRepo) NextSeq(ctx context.Context, conversationID string) (int64, error) {
	// Return the next sequence number for the given conversation ID
	key := "conv:" + conversationID + ":seq"
	result, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}
