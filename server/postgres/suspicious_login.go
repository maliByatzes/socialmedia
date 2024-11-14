package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

type SuspiciousLoginService struct {
	db *DB
}

func NewSuspiciousLoginService(db *DB) *SuspiciousLoginService {
	return &SuspiciousLoginService{db: db}
}

func (s *SuspiciousLoginService) FindSLByID(ctx context.Context, id uint) (*sm.SuspiciousLogin, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	sl, err := findSLByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (s *SuspiciousLoginService) FindSLByUserID(ctx context.Context, userID uint) (*sm.SuspiciousLogin, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	sl, err := findSLByUserID(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	return sl, nil
}

func (s *SuspiciousLoginService) FindSLs(ctx context.Context, filter sm.SLFilter) ([]*sm.SuspiciousLogin, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findSLs(ctx, tx, filter)
}

func (s *SuspiciousLoginService) CreateSL(ctx context.Context, sl *sm.SuspiciousLogin) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createSL(ctx, tx, sl); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SuspiciousLoginService) UpdateSL(ctx context.Context, id uint, upd sm.SLUpdate) (*sm.SuspiciousLogin, error) {
	return nil, nil
}

func (s *SuspiciousLoginService) DeleteSL(ctx context.Context, id uint) error {
	return nil
}

func createSL(ctx context.Context, tx *Tx, sl *sm.SuspiciousLogin) error {
	sl.CreatedAt = tx.now
	sl.UpdatedAt = sl.CreatedAt

	query := `INSERT INTO "suspicious_logins" ("user_id", "email", "ip", "country", "city", "browser", "platform", "os", "device", "device_type", "unverified_attempts", "is_trusted", "is_blocked", "created_at", "updated_at") 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id`
	args := []interface{}{
		sl.UserID,
		sl.Email,
		sl.IP,
		sl.Country,
		sl.City,
		sl.Browser,
		sl.Platform,
		sl.OS,
		sl.Device,
		sl.DeviceType,
		sl.UnverifiedAttempts,
		sl.IsTrusted,
		sl.IsBlocked,
		(*NullTime)(&sl.CreatedAt),
		(*NullTime)(&sl.UpdatedAt),
	}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&sl.ID)
	if err != nil {
		return err
	}

	return nil
}

func findSLByID(ctx context.Context, tx *Tx, id uint) (*sm.SuspiciousLogin, error) {
	a, _, err := findSLs(ctx, tx, sm.SLFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Suspicious Login not found."}
	}
	return a[0], nil
}

func findSLByUserID(ctx context.Context, tx *Tx, userID uint) (*sm.SuspiciousLogin, error) {
	a, _, err := findSLs(ctx, tx, sm.SLFilter{UserID: &userID})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Suspicious Login not found."}
	}
	return a[0], nil
}

func findSLs(ctx context.Context, tx *Tx, filter sm.SLFilter) (_ []*sm.SuspiciousLogin, n int, err error) {
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
		argPos++
	}

	if v := filter.UnverifiedAttempts; v != nil {
		where, args = append(where, fmt.Sprintf(`"unverified_attempts" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.IsTrusted; v != nil {
		where, args = append(where, fmt.Sprintf(`"is_trusted" = $%d`, argPos)), append(args, *v)
		argPos++
	}

	if v := filter.IsBlocked; v != nil {
		where, args = append(where, fmt.Sprintf(`"is_blocked" = $%d`, argPos)), append(args, *v)
	}

	query := `SELECT "id", "user_id", "email", "ip", "country", "city", "browser", "platform", "os", "device", "device_type", "unverified_attempts", "is_trusted", "is_blocked", "created_at", "updated_at", COUNT(*) OVER()
	 FROM "suspicious_logins"` + formatWhereClause(where) + ` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	sls := make([]*sm.SuspiciousLogin, 0)

	for rows.Next() {
		var sl sm.SuspiciousLogin

		if err := rows.Scan(
			&sl.ID,
			&sl.UserID,
			(*NullString)(&sl.Email),
			(*NullString)(&sl.IP),
			(*NullString)(&sl.Country),
			(*NullString)(&sl.City),
			(*NullString)(&sl.Browser),
			(*NullString)(&sl.Platform),
			(*NullString)(&sl.OS),
			(*NullString)(&sl.Device),
			(*NullString)(&sl.DeviceType),
			&sl.UnverifiedAttempts,
			&sl.IsTrusted,
			&sl.IsBlocked,
			(*NullTime)(&sl.CreatedAt),
			(*NullTime)(&sl.UpdatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		sls = append(sls, &sl)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return sls, n, nil
}
