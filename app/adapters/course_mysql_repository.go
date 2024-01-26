package adapters

import (
	"context"
	db "github.com/forhomme/app-base/usecase/database"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/pkg/errors"
	course2 "user-management/app/domain/course"
	"user-management/app/domain/user"
	"user-management/config"
)

const (
	queryAddCategory = `INSERT INTO category (category_name) VALUES (?)`

	queryGetAllCategory = `SELECT category_id, category_name FROM category WHERE deleted_at is null`

	queryGetUserById = `SELECT u.id, u.email,u.password,u.role_id, r.role_name FROM users u INNER JOIN roles r ON r.role_id=u.role_id WHERE u.id=?`

	queryGetUserByEmail = `SELECT u.id,u.email,u.password,u.role_id, r.role_name FROM users u INNER JOIN roles r ON r.role_id=u.role_id WHERE u.email=?`

	queryInsertUser = `INSERT INTO users (id,email,password,role_id,created_by) VALUES 
	 (?,?,?,?,?)`

	queryUpdateUserById = `UPDATE users SET email =  ?,password = ?,role_id = ? WHERE id=?`
)

type CourseMysqlRepository struct {
	cfg *config.Config
	log logger.Logger
	db.SqlHandler
}

func NewCourseMysqlRepository(cfg *config.Config, log logger.Logger, sqlHandler db.SqlHandler) *CourseMysqlRepository {
	return &CourseMysqlRepository{
		cfg:        cfg,
		log:        log,
		SqlHandler: sqlHandler,
	}
}

func (c *CourseMysqlRepository) AddCategory(ctx context.Context, categoryName string) error {
	_, err := c.SqlHandler.Exec(queryAddCategory, categoryName)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func (c *CourseMysqlRepository) GetCategories(ctx context.Context) ([]*course2.Category, error) {
	row, err := c.SqlHandler.Query(queryGetAllCategory)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	defer row.Close()

	out := make([]*course2.Category, 0)
	if row.Next() {
		var data = &course2.Category{}
		dataRows := []interface{}{
			&data.CategoryId,
			&data.CategoryName,
		}
		err = row.Scan(dataRows...)
		if err != nil {
			return nil, errors.Wrap(err, "Repo.GetCategories.Rows.Scan")
		}
		out = append(out, data)
	}
	return out, nil
}

func (c *CourseMysqlRepository) AddCourse(ctx context.Context, cr *course2.CoursePath) error {
	return errors.New("not implemented")
}

func (c *CourseMysqlRepository) GetCourses(ctx context.Context, in *course2.FilterCourse) ([]*course2.CoursePath, error) {
	return nil, errors.New("not implemented")
}

func (c *CourseMysqlRepository) UpdateCourse(ctx context.Context, id string, updateFn func(ctx context.Context, cm *course2.CoursePath) (*course2.CoursePath, error)) error {
	return errors.New("not implemented")
}

func (c *CourseMysqlRepository) GetUserById(id string) (*user.User, error) {
	row, err := c.SqlHandler.Query(queryGetUserById, id)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	defer row.Close()

	out := new(user.User)
	if row.Next() {
		dataRows := []interface{}{
			&out.UserId,
			&out.Email,
			&out.Password,
			&out.RoleId,
			&out.RoleName,
		}
		err = row.Scan(dataRows...)
		if err != nil {
			return nil, errors.Wrap(err, "Repo.GetUserById.Rows.Scan")
		}
	}
	return out, nil
}

func (c *CourseMysqlRepository) GetUserByEmail(email string) (*user.User, error) {
	row, err := c.SqlHandler.Query(queryGetUserByEmail, email)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}
	defer row.Close()

	out := new(user.User)
	if row.Next() {
		dataRows := []interface{}{
			&out.UserId,
			&out.Email,
			&out.Password,
			&out.RoleId,
			&out.RoleName,
		}
		err = row.Scan(dataRows...)
		if err != nil {
			return nil, errors.Wrap(err, "Repo.GetUserByEmail.Rows.Scan")
		}
	}
	return out, nil
}

func (c *CourseMysqlRepository) InsertUser(user *user.User) error {
	_, err := c.SqlHandler.Exec(queryInsertUser, user.UserId, user.Email, user.Password, user.RoleId, user.UserId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}

func (c *CourseMysqlRepository) UpdateUser(id string, updateFn func(u *user.User) (*user.User, error)) error {
	existingUser, err := c.GetUserById(id)
	if err != nil {
		c.log.Error(err)
		return err
	}

	updateUser, err := updateFn(existingUser)
	if err != nil {
		c.log.Error(err)
		return err
	}

	_, err = c.SqlHandler.Exec(queryUpdateUserById, updateUser.Email, updateUser.Password, updateUser.RoleId, updateUser.UserId)
	if err != nil {
		c.log.Error(err)
		return err
	}
	return nil
}
