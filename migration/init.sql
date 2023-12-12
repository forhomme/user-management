CREATE TABLE IF NOT EXISTS `users` (
    id VARCHAR(64) NOT NULL DEFAULT (UUID()),
    email VARCHAR(64) NOT NULL,
    password TEXT NOT NULL,
    role_id INT NOT NULL,
    is_teacher smallint default 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(64) DEFAULT NULL,
    deleted_at VARCHAR(64) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX (role_id)
);

CREATE TABLE IF NOT EXISTS `roles` (
    role_id BIGINT AUTO_INCREMENT NOT NULL,
    role_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(64) DEFAULT NULL,
    deleted_at VARCHAR(64) DEFAULT NULL,
    PRIMARY KEY (role_id)
);

CREATE TABLE IF NOT EXISTS `menus` (
    menu_id BIGINT AUTO_INCREMENT NOT NULL,
    parent_id INT NOT NULL,
    menu_name VARCHAR(50) NOT NULL,
    menu_path VARCHAR(255) NOT NULL,
    icon VARCHAR(255),
    ordering INT,
    published smallint default 0 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(64) DEFAULT NULL,
    deleted_at VARCHAR(64) DEFAULT NULL,
    PRIMARY KEY (menu_id),
    INDEX (parent_id)
    );

CREATE TABLE IF NOT EXISTS `role_menu` (
    id BIGINT AUTO_INCREMENT NOT NULL,
    menu_id INT NOT NULL,
    role_id INT NOT NULL,
    role_view smallint default 0,
    role_add smallint default 0,
    role_edit smallint default 0,
    role_delete smallint default 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(64) DEFAULT NULL,
    deleted_at VARCHAR(64) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX (menu_id),
    INDEX (role_id)
);

CREATE TABLE IF NOT EXISTS `audit_trail` (
   id BIGINT AUTO_INCREMENT NOT NULL,
   menu_name VARCHAR(1000) NOT NULL,
   menu VARCHAR(1000) NOT NULL,
   request VARCHAR(1000) NOT NULL,
   response VARCHAR(1000) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   created_by VARCHAR(64) NOT NULL,
    PRIMARY KEY (id),
    INDEX (created_by)
);

CREATE TABLE IF NOT EXISTS `course_category` (
    category_id BIGINT AUTO_INCREMENT NOT NULL,
    category VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(64) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS `courses` (
    id VARCHAR(64) NOT NULL DEFAULT (UUID()),
    parent_id VARCHAR(64) NOT NULL,
    category_id int NOT NULL,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    is_publish smallint default 0,
    ordering int,
    tags JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(64) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by VARCHAR(64) DEFAULT NULL,
    deleted_at VARCHAR(64) DEFAULT NULL,
    PRIMARY KEY (id),
    INDEX(parent_id),
    INDEX(category_id)
);

/*https://dba.stackexchange.com/questions/24531/mysql-create-index-if-not-exists*/
DROP PROCEDURE IF EXISTS `sp_CreateIndex`;
CREATE PROCEDURE `sp_CreateIndex`
(
    given_database VARCHAR(64),
    given_table    VARCHAR(64),
    given_index    VARCHAR(64),
    given_columns  VARCHAR(64)
)
BEGIN

    DECLARE IndexIsThere INTEGER;

SELECT COUNT(1) INTO IndexIsThere
FROM INFORMATION_SCHEMA.STATISTICS
WHERE table_schema = given_database
  AND   table_name   = given_table
  AND   index_name   = given_index;

IF IndexIsThere = 0 THEN
        SET @sqlstmt = CONCAT('CREATE INDEX ',given_index,' ON ',
        given_database,'.',given_table,' (',given_columns,')');
PREPARE st FROM @sqlstmt;
EXECUTE st;
DEALLOCATE PREPARE st;
ELSE
SELECT CONCAT('Index ',given_index,' already exists on Table ',
              given_database,'.',given_table) CreateindexErrorMessage;
END IF;

END;

DROP PROCEDURE IF EXISTS `sp_CreateUniqueIndex`;
CREATE PROCEDURE `sp_CreateUniqueIndex`
(
    given_database VARCHAR(64),
    given_table    VARCHAR(64),
    given_index    VARCHAR(64),
    given_columns  VARCHAR(64)
)
BEGIN

    DECLARE IndexIsThere INTEGER;

SELECT COUNT(1) INTO IndexIsThere
FROM INFORMATION_SCHEMA.STATISTICS
WHERE table_schema = given_database
  AND   table_name   = given_table
  AND   index_name   = given_index;

IF IndexIsThere = 0 THEN
        SET @sqlstmt = CONCAT('CREATE UNIQUE INDEX ',given_index,' ON ',
        given_database,'.',given_table,' (',given_columns,')');
PREPARE st FROM @sqlstmt;
EXECUTE st;
DEALLOCATE PREPARE st;
ELSE
SELECT CONCAT('Index ',given_index,' already exists on Table ',
              given_database,'.',given_table) CreateindexErrorMessage;
END IF;

END;

DROP PROCEDURE IF EXISTS `sp_AlterTable`;
CREATE PROCEDURE sp_AlterTable()
BEGIN
END;
CALL sp_AlterTable();
CALL sp_CreateIndex('puspeknubika', 'courses', 'tags_courses_idx', '(CAST(tags AS UNSIGNED ARRAY))');