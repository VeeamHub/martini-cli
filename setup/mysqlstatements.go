package setup

type NamedQuery struct {
	q string
	n string
}

func GetCreateStatements() []NamedQuery {
	return []NamedQuery{NamedQuery{`
	CREATE TABLE martini_tenant (
		id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		email VARCHAR(50),
		registered BIGINT,
		instancefqdn VARCHAR(100) NOT NULL,
		instanceusername VARCHAR(100) NOT NULL,
		instancepassword VARCHAR(100) NOT NULL
	);`, "Tenant"}, NamedQuery{`
	CREATE TABLE martini_user (
		id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(30) NOT NULL,
		email VARCHAR(50),
		hashpassword VARCHAR(150) NOT NULL
	);`, "User"}, NamedQuery{`	
	CREATE TABLE martini_token (
		id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		userid INT(6) NOT NULL,
		token VARCHAR(128),
		renew VARCHAR(128),
		validuntil BIGINT
	);`, "Token"}, NamedQuery{`	
	CREATE TABLE martini_securestore (
		id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		keyval VARCHAR(100),
		encryptedpassword VARCHAR(250)
	);`, "SecureStore"}, NamedQuery{`	
	CREATE TABLE martini_endpoint (
		id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		port int(7),
		validuntil BIGINT,
		pid BIGINT,
		tenantid INT(6)
	);`, "SecureStore"},
	}
}
