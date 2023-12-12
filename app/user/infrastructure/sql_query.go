package infrastructure

const (
	queryInsertAuditTrail = `INSERT INTO audit_trail (menu,method,request,response,created_at,created_by) VALUES
	 (?,?,?,?,UTC_TIMESTAMP(),?)`

	queryGetUserById = `SELECT u.id, u.email,u.password,u.role_id, r.role_name FROM users u INNER JOIN roles r ON r.role_id=u.role_id WHERE u.id=?`

	queryGetUserByEmail = `SELECT u.id,u.email,u.password,u.role_id, r.role_name FROM users u INNER JOIN roles r ON r.role_id=u.role_id WHERE u.email=?`

	queryInsertUser = `INSERT INTO users (id,email,password,role_id,created_by) VALUES 
	 (?,?,?,?,?)`

	queryUpdateUserById = `UPDATE users
	SET password = CASE WHEN NULLIF(?,'') IS NOT NULL THEN ? ELSE password END,
	    email = CASE WHEN NULLIF(?, '') IS NOT NULL THEN ? ELSE email END,
	    role_id = CASE WHEN NULLIF(?, 0) IS NOT NULL THEN ? ELSE role_id END WHERE id=?`

	queryGetAllMenuByRoleId = `SELECT menu_id,parent_id,menu_name,menu_path,icon,role_view,role_add,role_edit,role_delete 
	FROM role_menu rm
	INNER JOIN (SELECT * FROM menus WHERE published = 1 AND deleted_at is null) m ON m.menu_id = rm.menu_id
	WHERE rm.role_id = ? AND rm.deleted_at is null ORDER BY m.ordering ASC`

	queryInsertCourse = `CALL sp_InsertCourse(?,?,?,?,?,?,?,?,?,?,?,?,?)`

	queryGetCourseContent = `SELECT`
)
