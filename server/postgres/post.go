package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

type PostService struct {
	db *DB
}

func NewPostService(db *DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) FindPostByID(ctx context.Context, id uint) (*sm.Post, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	post, err := findPostByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) FindPosts(ctx context.Context, filter sm.PostFilter) ([]*sm.Post, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findPosts(ctx, tx, filter)
}

func (s *PostService) CreatePost(ctx context.Context, post *sm.Post) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createPost(ctx, tx, post); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PostService) UpdatePost(ctx context.Context, id uint, upd sm.PostUpdate) (*sm.Post, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	post, err := updatePost(ctx, tx, id, upd)
	if err != nil {
		return post, err
	} else if err := tx.Commit(); err != nil {
		return post, err
	}

	return post, err
}

func (s *PostService) DeletePost(ctx context.Context, id uint) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := deletePost(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit()
}

func findPostByID(ctx context.Context, tx *Tx, id uint) (*sm.Post, error) {
	a, _, err := findPosts(ctx, tx, sm.PostFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Post not found."}
	}
	return a[0], nil
}

func findPosts(ctx context.Context, tx *Tx, filter sm.PostFilter) (_ []*sm.Post, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.CommunityID; v != nil {
		where, args = append(where, fmt.Sprintf(`"community_id" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.UserID; v != nil {
		where, args = append(where, fmt.Sprintf(`"user_id" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "content", "file_url", "community_id", "user_id", "created_at", "update_at", COUNT(*) OVER()
	FROM "posts"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	posts := make([]*sm.Post, 0)
	for rows.Next() {
		var post sm.Post
		if err := rows.Scan(
			&post.ID,
			(*NullString)(&post.Content),
			(*NullString)(&post.FileURL),
			&post.CommunityID,
			&post.UserID,
			(*NullTime)(&post.CreatedAt),
			(*NullTime)(&post.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return posts, n, nil
}

func createPost(ctx context.Context, tx *Tx, post *sm.Post) error {
	if user, err := findUserByID(ctx, tx, post.UserID); err != nil {
		return err
	} else if sm.UserIDFromContext(ctx) != user.ID {
		return sm.Errorf(sm.ENOTAUTHORIZED, "You are not allowed to create this post.")
	}

	post.CreatedAt = tx.now
	post.UpdatedAt = post.CreatedAt

	query := `INSERT INTO "posts" ("id", "content", "file_url", "community_id", "user_id", "created_at", "updated_at")
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	args := []interface{}{
		post.ID,
		post.FileURL,
		post.CommunityID,
		post.UserID,
		(*NullTime)(&post.CreatedAt),
		(*NullTime)(&post.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&post.ID)
	if err != nil {
		return err
	}

	return nil
}

func updatePost(ctx context.Context, tx *Tx, id uint, upd sm.PostUpdate) (*sm.Post, error) {
	post, err := findPostByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	user, err := findUserByID(ctx, tx, post.UserID)
	if err != nil {
		return post, err
	} else if sm.UserIDFromContext(ctx) != user.ID {
		return nil, sm.Errorf(sm.ENOTAUTHORIZED, "You are not allowed to update this post.")
	}

	if v := upd.Content; v != nil {
		post.Content = *v
	}

	if v := upd.FileURL; v != nil {
		post.FileURL = *v
	}

	post.UpdatedAt = tx.now

	args := []interface{}{
		post.Content,
		post.FileURL,
		post.UpdatedAt,
		post.ID,
		user.ID,
	}
	query := `UPDATE "posts" SET "content" = $1, "file_url" = $2, "updated_at" = $3 WHERE "id" = $4 and "user_id" = $5`

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return post, err
	}

	return post, nil
}

func deletePost(ctx context.Context, tx *Tx, id uint) error {
	post, err := findPostByID(ctx, tx, id)
	if err != nil {
		return err
	}

	user, err := findUserByID(ctx, tx, post.UserID)
	if err != nil {
		return err
	} else if sm.UserIDFromContext(ctx) != user.ID {
		return sm.Errorf(sm.ENOTAUTHORIZED, "You are not allowed to delete this post.")
	}

	query := `DELETE FROM "posts" WHERE "id" = $1 AND "user_id" = $2`
	args := []interface{}{
		post.ID,
		user.ID,
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}
