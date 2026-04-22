\c hot-chip

DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE TABLE_NAME = 'users'
		AND COLUMN_NAME = 'password'
	) THEN
		ALTER TABLE "users"
		DROP COLUMN password;
	END IF;
END $$;

\c hot-chip_test

DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE TABLE_NAME = 'users'
		AND COLUMN_NAME = 'password'
	) THEN
		ALTER TABLE "users"
		DROP COLUMN password;
	END IF;
END $$;

/* Run using 'psql -U <username> -d <db_name> -f path_to_file/project_3_change_password_column.sql' 

If you use 'hot-chip' as the db name, you don't need to run it a second time for the test db. \c handles switching from 1 db to another.

Run 'psql -c "SELECT current_user;"' in the terminal to find your psql username
*/