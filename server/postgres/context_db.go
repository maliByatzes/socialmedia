package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
	"github.com/maliByatzes/socialmedia/utils"
)

type ContextService struct {
	db *DB
}

func NewContextService(db *DB) *ContextService {
	return &ContextService{db: db}
}

func (s *ContextService) FindContextByID(ctx context.Context, id uint) (*sm.Context, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	context, err := findContextByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return context, nil
}

func (s *ContextService) FindContextByUserID(ctx context.Context, userID uint) (*sm.Context, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	context, err := findContextByUserID(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	return context, nil
}

func (s *ContextService) FindContexts(ctx context.Context, filter sm.ContextFilter) ([]*sm.Context, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findContexts(ctx, tx, filter)
}

func (s *ContextService) CreateContext(ctx context.Context, context *sm.Context) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createContext(ctx, tx, context); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ContextService) UpdateContext(ctx context.Context, id uint, upd sm.ContextUpdate) (*sm.Context, error) {
	panic("not implemented")
}

func (s *ContextService) DeleteContext(ctx context.Context, id uint) error {
	panic("not implemented")
}

func createContext(ctx context.Context, tx *Tx, context *sm.Context) error {
	context.CreatedAt = tx.now
	context.UpdatedAt = context.CreatedAt

	if err := context.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO "context" ("user_id", "email", "ip", "country", "city", "browser", "platform", "os", "device", "device_type", "is_trusted", "created_at", "updated_at") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`
	args := []interface{}{
		context.UserID,
		context.Email,
		utils.EncryptData([]byte(context.IP)),
		utils.EncryptData([]byte(context.Country)),
		utils.EncryptData([]byte(context.City)),
		utils.EncryptData([]byte(context.Browser)),
		utils.EncryptData([]byte(context.Platform)),
		utils.EncryptData([]byte(context.OS)),
		utils.EncryptData([]byte(context.Device)),
		utils.EncryptData([]byte(context.DeviceType)),
		context.IsTrusted,
		(*NullTime)(&context.CreatedAt),
		(*NullTime)(&context.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&context.ID)
	if err != nil {
		return err
	}

	return nil
}

func findContextByID(ctx context.Context, tx *Tx, id uint) (*sm.Context, error) {
	a, _, err := findContexts(ctx, tx, sm.ContextFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Context not found."}
	}
	return a[0], nil
}

func findContextByUserID(ctx context.Context, tx *Tx, userID uint) (*sm.Context, error) {
	a, _, err := findContexts(ctx, tx, sm.ContextFilter{UserID: &userID})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Context not found."}
	}
	return a[0], nil
}

func findContexts(ctx context.Context, tx *Tx, filter sm.ContextFilter) (_ []*sm.Context, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPos := 1

	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.UserID; v != nil {
		where, args = append(where, fmt.Sprintf(`"user_id" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.Email; v != nil {
		where, args = append(where, fmt.Sprintf(`"email" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.IP; v != nil {
		where, args = append(where, fmt.Sprintf(`"ip" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.Country; v != nil {
		where, args = append(where, fmt.Sprintf(`"country" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.City; v != nil {
		where, args = append(where, fmt.Sprintf(`"city" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.Browser; v != nil {
		where, args = append(where, fmt.Sprintf(`"browser" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.Platform; v != nil {
		where, args = append(where, fmt.Sprintf(`"platform" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.OS; v != nil {
		where, args = append(where, fmt.Sprintf(`"os" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.Device; v != nil {
		where, args = append(where, fmt.Sprintf(`"device" = $%d`, argPos)), append(args, *v)
		argPos++
	}
	if v := filter.DeviceType; v != nil {
		where, args = append(where, fmt.Sprintf(`"device_type" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "user_id", "email", "ip", "country", "city", "browser", "platform", "os", "device", "device_type", "is_trusted", "created_at", "updated_at", COUNT(*) OVER() 
	 FROM "context" WHERE` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	contexts := make([]*sm.Context, 0)
	for rows.Next() {
		var context sm.Context
		if err := rows.Scan(
			&context.ID,
			&context.UserID,
			(*NullString)(&context.Email),
			(*NullString)(&context.IP),
			(*NullString)(&context.Country),
			(*NullString)(&context.City),
			(*NullString)(&context.Browser),
			(*NullString)(&context.Platform),
			(*NullString)(&context.OS),
			(*NullString)(&context.Device),
			(*NullString)(&context.DeviceType),
			&context.IsTrusted,
			(*NullTime)(&context.CreatedAt),
			(*NullTime)(&context.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		// NOTE: Will this even work?? IDK we'll see
		context.IP = utils.DecryptData([]byte(context.IP))
		context.Country = utils.DecryptData([]byte(context.Country))
		context.City = utils.DecryptData([]byte(context.City))
		context.Browser = utils.DecryptData([]byte(context.Browser))
		context.Platform = utils.DecryptData([]byte(context.Platform))
		context.OS = utils.DecryptData([]byte(context.OS))
		context.Device = utils.DecryptData([]byte(context.Device))
		context.DeviceType = utils.DecryptData([]byte(context.DeviceType))

		contexts = append(contexts, &context)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return contexts, n, nil
}
