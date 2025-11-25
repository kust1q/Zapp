package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kust1q/Zapp/backend/internal/domain/entity"
	conv "github.com/kust1q/Zapp/backend/internal/pkg/conv/db"
	"github.com/kust1q/Zapp/backend/internal/providers/db/models"
	"github.com/kust1q/Zapp/backend/internal/providers/db/redis/cache"
	"github.com/kust1q/Zapp/backend/internal/storage/postgres"
)

func (pg *PostgresDB) CreateUserTx(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	userModel := conv.FromDomainToUserModel(user)
	if userModel == nil {
		return nil, fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (username, email, password, bio, gen, created_at, is_active, is_superuser) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`, postgres.UserTable)

	var id int
	err := tx.QueryRowContext(ctx, query,
		userModel.Username, userModel.Email, userModel.Password,
		userModel.Bio, userModel.Gen, userModel.CreatedAt,
		userModel.IsActive, userModel.IsSuperuser).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	userModel.ID = id
	createduser := conv.FromUserModelToDomain(userModel)
	return createduser, nil
}

func (pg *PostgresDB) CreateAdmin(ctx context.Context, user *entity.User) error {
	userModel := conv.FromDomainToUserModel(user)
	if userModel == nil {
		return fmt.Errorf("cannot convert nil entity to DB model")
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (username, email, password, bio, gen, created_at, is_active, is_superuser) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, postgres.UserTable)

	_, err := pg.db.ExecContext(ctx, query,
		userModel.Username, userModel.Email, userModel.Password,
		userModel.Bio, userModel.Gen, userModel.CreatedAt,
		userModel.IsActive, userModel.IsSuperuser)
	return err
}

func (pg *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1", postgres.UserTable)
	var userModel models.User
	err := pg.db.GetContext(ctx, &userModel, query, email)
	if err != nil {
		return nil, err
	}
	user := conv.FromUserModelToDomain(&userModel)
	return user, nil
}

func (pg *PostgresDB) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", postgres.UserTable)
	var userModel models.User
	err := pg.db.GetContext(ctx, &userModel, query, username)
	if err != nil {
		return nil, err
	}
	user := conv.FromUserModelToDomain(&userModel)
	return user, nil
}

func (pg *PostgresDB) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", postgres.UserTable)
	var userModel models.User
	err := pg.db.GetContext(ctx, &userModel, query, userID)
	if err != nil {
		return nil, err
	}
	user := conv.FromUserModelToDomain(&userModel)
	return user, nil
}

func (pg *PostgresDB) UpdateUserPassword(ctx context.Context, userID int, password string) error {
	query := fmt.Sprintf("UPDATE %s SET password = $1 WHERE id = $2", postgres.UserTable)
	_, err := pg.db.ExecContext(ctx, query, password, userID)
	return err
}

// UpdateUserBio
func (pg *PostgresDB) UpdateUserBio(ctx context.Context, userID int, bio string) error {
	query := fmt.Sprintf("UPDATE %s SET bio = $1 WHERE id = $2", postgres.UserTable)
	_, err := pg.db.ExecContext(ctx, query, bio, userID)
	return err
}

func (pg *PostgresDB) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := pg.userCache.Exists(ctx, cache.UsernameType, email)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email = $1", postgres.UserTable)
	var count int
	err = pg.db.GetContext(ctx, &count, query, email)
	return count > 0, err
}

func (pg *PostgresDB) UserExistsByUsername(ctx context.Context, username string) (bool, error) {
	exists, err := pg.userCache.Exists(ctx, cache.UsernameType, username)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = $1", postgres.UserTable)
	var count int
	err = pg.db.GetContext(ctx, &count, query, username)
	return count > 0, err
}

func (pg *PostgresDB) FollowToUser(ctx context.Context, followerID, followingID int, createdAt time.Time) (*entity.Follow, error) {
	query := fmt.Sprintf("INSERT INTO %s (follower_id, following_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", postgres.FollowsTable)
	_, err := pg.db.ExecContext(ctx, query, followerID, followingID, createdAt)
	if err != nil {
		return nil, err
	}

	followEntity := &entity.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   createdAt,
	}
	return followEntity, nil
}

func (pg *PostgresDB) UnfollowUser(ctx context.Context, followerID, followingID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE follower_id = $1 AND following_id = $2", postgres.FollowsTable)
	result, err := pg.db.ExecContext(ctx, query, followerID, followingID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("follow not found")
	}
	return nil
}

func (pg *PostgresDB) GetFollowersIds(ctx context.Context, username string) ([]int, error) {
	query := fmt.Sprintf(`
        SELECT f.follower_id
        FROM %s f
        JOIN %s u ON f.following_id = u.id
        WHERE u.username = $1`,
		postgres.FollowsTable, postgres.UserTable)

	var res []int
	if err := pg.db.SelectContext(ctx, &res, query, username); err != nil {
		return nil, err
	}
	return res, nil
}

func (pg *PostgresDB) GetFollowingsIds(ctx context.Context, username string) ([]int, error) {
	query := fmt.Sprintf(`
        SELECT f.following_id
        FROM %s f
        JOIN %s u ON f.follower_id = u.id
        WHERE u.username = $1`,
		postgres.FollowsTable, postgres.UserTable)

	var res []int
	if err := pg.db.SelectContext(ctx, &res, query, username); err != nil {
		return nil, err
	}
	return res, nil
}

func (pg *PostgresDB) DeleteUser(ctx context.Context, userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", postgres.UserTable)
	result, err := pg.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
