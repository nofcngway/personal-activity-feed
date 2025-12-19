package pgstorage

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
)

const feedTable = "activity_feed"

type FeedItem struct {
	ActorID   int64
	Action    string
	TargetID  int64
	CreatedAt time.Time
}

func (s *PGStorage) InsertFeedItem(ctx context.Context, userID, actorID int64, action string, targetID int64, createdAt time.Time) error {
	q := squirrel.Insert(feedTable).
		Columns("user_id", "actor_id", "action", "target_id", "created_at").
		Values(userID, actorID, action, targetID, createdAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = s.db.Exec(ctx, sql, args...)
	return err
}

func (s *PGStorage) GetFeed(ctx context.Context, userID int64, limit, offset int32) ([]FeedItem, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	q := squirrel.Select("actor_id", "action", "target_id", "created_at").
		From(feedTable).
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]FeedItem, 0, limit)
	for rows.Next() {
		var it FeedItem
		if err := rows.Scan(&it.ActorID, &it.Action, &it.TargetID, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}


