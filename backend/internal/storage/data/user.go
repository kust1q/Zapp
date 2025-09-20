package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/storage/cache"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

func (s *dataStorage) CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	query := fmt.Sprintf(`
        INSERT INTO %s (username, email, password, bio, gen, avatar_url, created_at, is_superuser) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`, postgres.UserTable)
	var id int
	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.Bio, user.Gen, user.CreatedAt, user.IsSuperuser).Scan(&id)
	if err != nil {
		return &entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	user.ID = id
	return user, nil
}

func (s *dataStorage) SetSecretQuestionTx(ctx context.Context, tx *sql.Tx, question *entity.SecretQuestion) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, question, answer)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET
        secret_question = EXCLUDED.secret_question,
        answer = EXCLUDED.answer`,
		postgres.SecretQuestionTable)
	_, err := tx.ExecContext(ctx, query, question.UserID, question.SecretQuestion, question.Answer)
	if err != nil {
		return fmt.Errorf("failed to set secret question: %w", err)
	}
	return nil
}

func (s *dataStorage) CreateAdmin(ctx context.Context, user *entity.User) error {
	query := fmt.Sprintf(`
	INSERT INTO %s (username, email, password, bio, gen, created_at, is_superuser) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, postgres.UserTable)
	_, err := s.db.ExecContext(ctx, query, user.Email, user.Password, user.Bio, user.Gen, user.CreatedAt, user.IsSuperuser)
	return err
}

func (s *dataStorage) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", postgres.UserTable)
	var user entity.User
	err := s.db.GetContext(ctx, &user, query, email)
	return &user, err
}

func (s *dataStorage) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", postgres.UserTable)
	var user entity.User
	err := s.db.GetContext(ctx, &user, query, username)
	return &user, err
}

func (s *dataStorage) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postgres.UserTable)
	var user entity.User
	err := s.db.GetContext(ctx, &user, query, userID)
	return &user, err
}

func (s *dataStorage) SetSecretQuestion(ctx context.Context, question *entity.SecretQuestion) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, question, answer) VALUES ($1, $2, $3)", postgres.SecretQuestionTable)
	_, err := s.db.ExecContext(ctx, query, question.UserID, question.SecretQuestion, question.Answer)
	return err
}

func (s *dataStorage) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE id = $2", postgres.UserTable)
	_, err := s.db.ExecContext(ctx, query, password, userID)
	return err
}

func (s *dataStorage) GetSecurityDataByUserID(ctx context.Context, userID int) (*entity.SecretQuestion, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postgres.SecretQuestionTable)
	var question entity.SecretQuestion
	err := s.db.GetContext(ctx, &question, query, userID)
	return &question, err
}

func (s *dataStorage) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := s.userCache.Exists(ctx, cache.EmailType, email)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", postgres.UserTable)
	var count int
	err = s.db.GetContext(ctx, &count, query, email)
	return count > 0, err
}

func (s *dataStorage) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := s.userCache.Exists(ctx, cache.UsernameType, username)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = $1", postgres.UserTable)
	var count int
	err = s.db.GetContext(ctx, &count, query, username)
	return count > 0, err
}

func (s *dataStorage) DeleteUser(ctx context.Context, userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", postgres.UserTable)
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}
