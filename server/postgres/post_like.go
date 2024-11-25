package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

type PostLikeService struct {
	db *DB
}

func NewPostLikeService(db *DB) *PostLikeService {
	return &PostLikeService{db: db}
}

func (s *PostLikeService) CreatePostLike(ctx context.Context, postLike *sm.PostLike) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createPostLike(ctx, tx, postLike); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostLikeService) FindPostLikes(ctx context.Context, filter sm.PostLikeFilter) ([]*sm.PostLike, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findPostLikes(ctx, tx, filter)
}

func (s *PostLikeService) DeletePostLike(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deletePostLike(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func createPostLike(ctx context.Context, tx *Tx, postLike *sm.PostLike) error {
	// NOTE: Unsure about this check ...
	user, err := findUserByID(ctx, tx, postLike.UserID)
	if err != nil {
		return err
	} else if sm.UserIDFromContext(ctx) != user.ID {
		return sm.Errorf(sm.ENOTAUTHORIZED, "You are not allowed to like this post.")
	}

	postLike.CreatedAt = tx.now

	query := `INSERT INTO "post_likes" ("post_id", "user_id", "created_at")
	VALUES ($1, $2, $3)`
	args := []interface{}{
		postLike.PostID,
		postLike.UserID,
		postLike.CreatedAt,
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func findPostLikes(ctx context.Context, tx *Tx, filter sm.PostLikeFilter) (_ []*sm.PostLike, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.PostID; v != nil {
		where, args = append(where, fmt.Sprintf(`"post_id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.UserID; v != nil {
		where, args = append(where, fmt.Sprintf(`"user_id" = $%d`, argPos)), append(args, *v)
	}

	if filter.Fill {
		// TODO: Return filled columns instead of `id`s
	}

	query := `SELECT "post_id", "user_id", "created_at", COUNT(*) OVER() 
		FROM "post_likes"` + formatWhereClause(where) + ` ORDER BY "post_id" ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	postLikes := make([]*sm.PostLike, 0)
	for rows.Next() {
		var postLike sm.PostLike
		if err := rows.Scan(
			&postLike.PostID,
			&postLike.UserID,
			(*NullTime)(&postLike.CreatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		postLikes = append(postLikes, &postLike)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return postLikes, n, nil
}

func deletePostLike(ctx context.Context, tx *Tx, id uint) error {
	query := `DELETE FROM "post_likes" WHERE "post_id" = $1 AND "user_id" = $2`
	args := []interface{}{
		id,
		sm.UserIDFromContext(ctx),
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
