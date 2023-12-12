package infrastructure

import (
	db "github.com/forhomme/app-base/usecase/database"
	"github.com/forhomme/app-base/usecase/logger"
	"github.com/pkg/errors"
	domain2 "user-management/app/user/domain"
	"user-management/config"
)

type database struct {
	cfg    *config.Config
	logger logger.Logger
	db.SqlHandler
}

func NewDatabase(cfg *config.Config, logger logger.Logger, sqlHandler db.SqlHandler) DatabasePorts {
	return &database{
		cfg:        cfg,
		logger:     logger,
		SqlHandler: sqlHandler,
	}
}

func (d *database) InsertAuditTrail(trail *domain2.AuditTrail) error {
	_, err := d.SqlHandler.Exec(queryInsertAuditTrail, trail.Menu, trail.Method, trail.Request, trail.Response, trail.UserId)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	return nil
}

func (d *database) GetUserById(id string) (*domain2.UsersDatabase, error) {
	row, err := d.SqlHandler.Query(queryGetUserById, id)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	defer row.Close()

	out := new(domain2.UsersDatabase)
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

func (d *database) GetUserByEmail(email string) (*domain2.UsersDatabase, error) {
	row, err := d.SqlHandler.Query(queryGetUserByEmail, email)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	defer row.Close()

	out := new(domain2.UsersDatabase)
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

func (d *database) InsertUser(usersDatabase *domain2.UsersDatabase) (*domain2.UsersDatabase, error) {
	_, err := d.SqlHandler.Exec(queryInsertUser, usersDatabase.UserId, usersDatabase.Email, usersDatabase.Password, usersDatabase.RoleId, usersDatabase.UserId)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}

	out, err := d.GetUserById(usersDatabase.UserId)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	return out, nil
}

func (d *database) UpdateUserById(id string, dataUser *domain2.UsersDatabase) error {
	_, err := d.SqlHandler.Exec(queryUpdateUserById,
		dataUser.Password,
		dataUser.Password,
		dataUser.Email,
		dataUser.Email,
		dataUser.RoleId,
		dataUser.RoleId,
		id)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	return nil
}

func (d *database) GetUserMenu(roleId int) ([]*domain2.Menu, error) {
	rows, err := d.SqlHandler.Query(queryGetAllMenuByRoleId, roleId)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	out := make([]*domain2.Menu, 0)
	if rows.Next() {
		m := new(domain2.Menu)
		dataOutM := []interface{}{
			m.MenuId,
			m.MenuParentId,
			m.MenuName,
			m.MenuPath,
			m.MenuIcon,
			m.RoleView,
			m.RoleAdd,
			m.RoleEdit,
			m.RoleDelete,
		}
		err = rows.Scan(dataOutM...)
		if err != nil {
			return nil, errors.Wrap(err, "Repo.GetUserMenu.Rows.Scan")
		}
		out = append(out, m)
	}
	return out, nil
}

func (d *database) GetAllCourseCategory() ([]*domain2.CourseCategoryDatabase, error) {
	out := make([]*domain2.CourseCategoryDatabase, 0)
	rows, err := d.SqlHandler.Query(queryGetAllMenuByRoleId)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		c := new(domain2.CourseCategoryDatabase)
		dataOutM := []interface{}{
			c.CategoryId,
			c.CategoryName,
		}
		err = rows.Scan(dataOutM...)
		if err != nil {
			return nil, errors.Wrap(err, "Repo.GetAllCourseCategory.Rows.Scan")
		}
		out = append(out, c)
	}
	return out, nil
}

func (d *database) InsertCourseContent(in *domain2.CourseDatabase) ([]*domain2.CourseDatabase, error) {
	rows, err := d.SqlHandler.Query(queryInsertCourse,
		in.Id,
		in.ParentId,
		in.CategoryId,
		in.Title,
		in.Description,
		in.Content,
		in.Image,
		in.Video,
		in.Tags,
		in.Ordering,
		in.IsPublish,
		in.CreatedBy,
		in.UpdatedBy)
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}
	defer rows.Close()

	out := make([]*domain2.CourseDatabase, 0)
	for rows.Next() {
		c := new(domain2.CourseDatabase)
		dataOutM := []interface{}{
			c.Id,
			c.ParentId,
			c.CategoryId,
			c.Title,
			c.Description,
			c.Content,
			c.Image,
			c.Video,
			c.Tags,
			c.Ordering,
			c.IsPublish,
			c.CreatedBy,
			c.UpdatedBy,
		}
		err = rows.Scan(dataOutM...)
		if err != nil {
			return nil, errors.Wrap(err, "Repo.InsertCourseContent.Rows.Scan")
		}
		out = append(out, c)
	}
	return out, nil
}
