package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	"github.com/kust1q/Zapp/backend/internal/storage/cache"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

type UserCache interface {
	Exists(ctx context.Context, dataType, data string) (bool, error)
}

type userStorage struct {
	db    *sqlx.DB
	cache UserCache
}

func NewUserStorage(db *sqlx.DB, cache UserCache) *userStorage {
	return &userStorage{
		db:    db,
		cache: cache,
	}
}

func (s *userStorage) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

func (s *userStorage) CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (entity.User, error) {
	query := fmt.Sprintf(`
        INSERT INTO %s (
            username, 
            email, 
            password, 
            bio, 
            gen, 
            avatar_url, 
            created_at, 
            is_superuser
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `, postgres.UserTable)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := tx.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Bio,
		user.Gen,
		user.AvatarURL,
		user.CreatedAt,
		user.IsSuperuser,
	).Scan(&user.ID)

	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return *user, nil
}

func (s *userStorage) SetSecretQuestionTx(ctx context.Context, tx *sql.Tx, question *entity.SecretQuestion) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (user_id, question, answer)
		VALUES ($1, $2, $3) ON CONFLICT (user_id) DO UPDATE SET
        secret_question = EXCLUDED.secret_question,
        answer = EXCLUDED.answer`,
		postgres.SecretQuestionTable)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := tx.ExecContext(ctx, query,
		question.UserID,
		question.SecretQuestion,
		question.Answer,
	)

	if err != nil {
		return fmt.Errorf("failed to set secret question: %w", err)
	}

	return nil
}

func (s *userStorage) CreateAdmin(user *entity.User) error {
	query := fmt.Sprintf("INSERT INTO %s (username, email, password, bio, gen, avatar_url, created_at, is_superuser) VALUES (:username, :email, :password, :bio, :gen, :avatar_url, :created_at, :is_superuser) RETURNING id", postgres.UserTable)

	stmt, err := s.db.PrepareNamed(query)
	if err != nil {
		return err
	}

	err = stmt.QueryRow(user).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *userStorage) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", postgres.UserTable)
	var user entity.User
	err := s.db.GetContext(ctx, &user, query, email)
	return user, err
}

func (s *userStorage) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", postgres.UserTable)
	var user entity.User
	err := s.db.GetContext(ctx, &user, query, username)
	return user, err
}

func (s *userStorage) GetUserByID(ctx context.Context, userID int) (entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postgres.UserTable)
	var user entity.User
	err := s.db.GetContext(ctx, &user, query, userID)
	return user, err
}

func (s *userStorage) SetSecretQuestion(ctx context.Context, question *entity.SecretQuestion) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, question, answer) VALUES (:user_id, :question, :answer)", postgres.SecretQuestionTable)
	_, err := s.db.NamedExecContext(ctx, query, question)
	return err
}

func (s *userStorage) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE id = $2", postgres.UserTable)
	_, err := s.db.ExecContext(ctx, query, password, userID)
	return err
}

func (s *userStorage) GetSecurityDataByUserID(ctx context.Context, userID int) (entity.SecretQuestion, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postgres.SecretQuestionTable)
	var question entity.SecretQuestion
	err := s.db.GetContext(ctx, &question, query, userID)
	return question, err
}

func (s *userStorage) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := s.cache.Exists(ctx, cache.EmailType, email)
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

func (s *userStorage) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := s.cache.Exists(ctx, cache.UsernameType, username)
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

func (s *userStorage) DeleteUser(ctx context.Context, userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", postgres.UserTable)
	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}
