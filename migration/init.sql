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

CREATE TABLE IF NOT EXISTS `category` (
     category_id BIGINT AUTO_INCREMENT NOT NULL,
     category_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (category_id),
    UNIQUE KEY unique_category (category_name)
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