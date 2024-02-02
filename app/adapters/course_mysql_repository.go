package adapters

import (
	"context"
	"github.com/forhomme/app-base/infrastructure/baselogger"
	"github.com/forhomme/app-base/infrastructure/telemetry"
	db "github.com/forhomme/app-base/usecase/database"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
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
	cfg    *config.Config
	log    *baselogger.Logger
	tracer *telemetry.OtelSdk
	db.SqlHandler
}

func NewCourseMysqlRepository(cfg *config.Config, log *baselogger.Logger, sqlHandler db.SqlHandler, tracer *telemetry.OtelSdk) *CourseMysqlRepository {
	return &CourseMysqlRepository{
		cfg:        cfg,
		log:        log,
		tracer:     tracer,
		SqlHandler: sqlHandler,
	}
}

func (c *CourseMysqlRepository) AddCategory(ctx context.Context, categoryName string) (err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.add_category")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	_, err = c.SqlHandler.Exec(queryAddCategory, categoryName)
	if err != nil {
		return err
	}
	return nil
}

func (c *CourseMysqlRepository) GetCategories(ctx context.Context) (out []*course2.Category, err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.get_category")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	rows, err := c.SqlHandler.Query(queryGetAllCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out = make([]*course2.Category, 0)
	for rows.Next() {
		var data = &course2.Category{}
		dataRows := []interface{}{
			&data.CategoryId,
			&data.CategoryName,
		}
		err = rows.Scan(dataRows...)
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

func (c *CourseMysqlRepository) GetUserById(ctx context.Context, id string) (out *user.User, err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.getuserid")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	row, err := c.SqlHandler.Query(queryGetUserById, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	out = new(user.User)
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

func (c *CourseMysqlRepository) GetUserByEmail(ctx context.Context, email string) (out *user.User, err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.getuseremail")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	row, err := c.SqlHandler.Query(queryGetUserByEmail, email)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	out = new(user.User)
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

func (c *CourseMysqlRepository) InsertUser(ctx context.Context, user *user.User) error {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.insertuser")
	defer span.End()

	_, err := c.SqlHandler.Exec(queryInsertUser, user.UserId, user.Email, user.Password, user.RoleId, user.UserId)
	if err != nil {
		c.log.Error(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

func (c *CourseMysqlRepository) UpdateUser(ctx context.Context, id string, updateFn func(u *user.User) (*user.User, error)) (err error) {
	ctx, span := c.tracer.Tracer.Start(ctx, "db.updateuser")
	defer span.End()

	defer func() {
		if err != nil {
			c.log.Error(err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	existingUser, err := c.GetUserById(ctx, id)
	if err != nil {
		return err
	}

	updateUser, err := updateFn(existingUser)
	if err != nil {
		return err
	}

	_, err = c.SqlHandler.Exec(queryUpdateUserById, updateUser.Email, updateUser.Password, updateUser.RoleId, updateUser.UserId)
	if err != nil {
		return err
	}
	return nil
}
