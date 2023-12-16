// SQL语句
package common

var (
	CreateTable = `CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(255) NOT NULL,
		user_id INT NOT NULL,
		user_role VARCHAR(255) NOT NULL DEFAULT 'user',
		last_active_date DATE
	 );`
)