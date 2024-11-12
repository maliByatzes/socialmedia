package postgres

import (
	"context"
	"fmt"

	sm "github.com/maliByatzes/socialmedia"
)

var _ sm.PreferenceService = (*PreferenceService)(nil)

type PreferenceService struct {
	db *DB
}

func NewPreferenceService(db *DB) *PreferenceService {
	return &PreferenceService{db: db}
}
func (s *PreferenceService) FindPreferenceByID(ctx context.Context, id uint) (*sm.Preference, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	preference, err := findPreferenceByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return preference, nil
}

func (s *PreferenceService) FindPreferenceByUserID(ctx context.Context, userID uint) (*sm.Preference, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	preference, err := findPreferenceByUserID(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	return preference, nil
}

func (s *PreferenceService) FindPreferences(ctx context.Context, filter sm.PreferenceFilter) ([]*sm.Preference, int, error) {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	return findPreferences(ctx, tx, filter)
}

func (s *PreferenceService) CreatePreference(ctx context.Context, preference *sm.Preference) error {
	tx := s.db.BeginTx(ctx, nil)
	defer tx.Rollback()

	if err := createPreference(ctx, tx, preference); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *PreferenceService) UpdatePreference(ctx context.Context, id uint, upd sm.PreferenceUpdate) (*sm.Preference, error) {
	panic("not implemented")
}

func (s *PreferenceService) DeletePreference(ctx context.Context, id uint) error {
	panic("not implemented")
}

func findPreferenceByID(ctx context.Context, tx *Tx, id uint) (*sm.Preference, error) {
	a, _, err := findPreferences(ctx, tx, sm.PreferenceFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Preference not found."}
	}

	return a[0], nil
}

func findPreferenceByUserID(ctx context.Context, tx *Tx, userID uint) (*sm.Preference, error) {
	a, _, err := findPreferences(ctx, tx, sm.PreferenceFilter{UserID: &userID})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &sm.Error{Code: sm.ENOTFOUND, Message: "Preference not found."}
	}

	return a[0], nil
}

func findPreferences(ctx context.Context, tx *Tx, filter sm.PreferenceFilter) (_ []*sm.Preference, n int, err error) {
	where, args := []string{}, []interface{}{}
	argPosition := 0

	if v := filter.ID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf(`"id" = $%d`, argPosition)), append(args, *v)
	}

	if v := filter.UserID; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf(`"user_id" = $%d`, argPosition)), append(args, *v)
	}

	if v := filter.EnabledContextBasedAuth; v != nil {
		argPosition++
		where, args = append(where, fmt.Sprintf(`"enabled_context_auth_enabled" = $%d`, argPosition)), append(args, *v)
	}

	query := `SELECT "id", "user_id", "enable_context_based_auth", "created_at", COUNT(*) OVER() FROM "preferences"` + formatWhereClause(where) +
		` ORDER BY id ASC` + formatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	preferences := make([]*sm.Preference, 0)
	for rows.Next() {
		var preference sm.Preference
		if err := rows.Scan(
			&preference.ID,
			&preference.UserID,
			&preference.EnabledContextBasedAuth,
			(*NullTime)(&preference.CreatedAt),
			&n,
		); err != nil {
			return nil, n, err
		}

		preferences = append(preferences, &preference)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return preferences, n, nil
}

func createPreference(ctx context.Context, tx *Tx, preference *sm.Preference) error {
	preference.CreatedAt = tx.now

	if err := preference.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO "preferences" ("user_id", "enable_context_based_auth", "created_at")
			  VALUES ($1, $2, $3) RETURNING "id"`
	args := []interface{}{preference.UserID, preference.EnabledContextBasedAuth, (*NullTime)(&preference.CreatedAt)}

	err := tx.QueryRowxContext(ctx, query, args...).Scan(&preference.ID)
	if err != nil {
		return err
	}

	return nil
}
