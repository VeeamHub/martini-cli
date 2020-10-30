package setup

type NamedQuery struct {
	q string
	n string
}

func GetCreateStatements() []NamedQuery {
	return []NamedQuery{NamedQuery{`
	CREATE TABLE ` + "`" + `martini_tenant` + "`" + ` (
		` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL,
		` + "`" + `name` + "`" + ` varchar(30) NOT NULL,
		` + "`" + `email` + "`" + ` varchar(50) DEFAULT NULL,
		` + "`" + `password` + "`" + ` varchar(250) DEFAULT NULL,
		` + "`" + `registered` + "`" + ` bigint(20) DEFAULT NULL
	  );`, "Tenant"}, NamedQuery{`
	  CREATE TABLE ` + "`" + `martini_tenant_instances` + "`" + ` (
		` + "`" + `id` + "`" + ` int(10) NOT NULL,
		` + "`" + `name` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `tenant_id` + "`" + ` int(10) NOT NULL,
		` + "`" + `json` + "`" + ` text NOT NULL,
		` + "`" + `type` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `status` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `location` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `hostname` + "`" + ` varchar(250) DEFAULT NULL,
		` + "`" + `port` + "`" + ` int(10) NOT NULL DEFAULT '4443',
		` + "`" + `username` + "`" + ` varchar(250) DEFAULT NULL,
		` + "`" + `password` + "`" + ` varchar(250) DEFAULT NULL
	  );`, "martini_tenant_instances"}, NamedQuery{`
	  CREATE TABLE ` + "`" + `martini_user` + "`" + ` (
		` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL,
		` + "`" + `name` + "`" + ` varchar(30) NOT NULL,
		` + "`" + `email` + "`" + ` varchar(50) DEFAULT NULL,
		` + "`" + `hashpassword` + "`" + ` varchar(150) NOT NULL
	  );`, "User"}, NamedQuery{`	
	CREATE TABLE ` + "`" + `martini_token` + "`" + ` (
		` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL,
		` + "`" + `userid` + "`" + ` int(6) NOT NULL,
		` + "`" + `token` + "`" + ` varchar(128) DEFAULT NULL,
		` + "`" + `renew` + "`" + ` varchar(128) DEFAULT NULL,
		` + "`" + `validuntil` + "`" + ` bigint(20) DEFAULT NULL
	  );`, "Token"}, NamedQuery{`	
	CREATE TABLE ` + "`" + `martini_securestore` + "`" + ` (
		` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL,
		` + "`" + `keyval` + "`" + ` varchar(100) DEFAULT NULL,
		` + "`" + `encryptedpassword` + "`" + ` varchar(250) DEFAULT NULL,
		` + "`" + `encryptedpassword_aes` + "`" + ` BLOB DEFAULT NULL
	  );`, "SecureStore"}, NamedQuery{`
	CREATE TABLE ` + "`" + `martini_endpoint` + "`" + ` (
		` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL,
		` + "`" + `port` + "`" + ` int(7) DEFAULT NULL,
		` + "`" + `validuntil` + "`" + ` bigint(20) DEFAULT NULL,
		` + "`" + `pid` + "`" + ` bigint(20) DEFAULT NULL,
		` + "`" + `tenantid` + "`" + ` int(6) DEFAULT NULL
	  );`, "Endpoint"}, NamedQuery{`
	  CREATE TABLE ` + "`" + `martini_general_aws_config` + "`" + ` (
		` + "`" + `provider` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `region` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `accesskey` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `secretkey` + "`" + ` varchar(250) NOT NULL
	  );`, "AWS Config"}, NamedQuery{`
	  CREATE TABLE ` + "`" + `martini_provider_aws_region` + "`" + ` (
		` + "`" + `region` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `vpc` + "`" + ` varchar(250) NOT NULL,
		` + "`" + `privatekey` + "`" + ` text NOT NULL
	  );`, "Provider AWS Region"}, NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_endpoint` + "`" + `
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `);
	  `, "primary key endpoint"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_general_aws_config` + "`" + `
		ADD UNIQUE KEY ` + "`" + `provider` + "`" + ` (` + "`" + `provider` + "`" + `);
	  `, "unique key general aws config"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_provider_aws_region` + "`" + `
		ADD UNIQUE KEY ` + "`" + `region` + "`" + ` (` + "`" + `region` + "`" + `);
	  `, "unique key provider aws region"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_securestore` + "`" + `
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `);
	  `, "primary key secure store"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_tenant` + "`" + `
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `),
		ADD UNIQUE KEY ` + "`" + `email` + "`" + ` (` + "`" + `email` + "`" + `);
	  `, "primary key tenant"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_tenant_instances` + "`" + `
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `);
	  `, "primary key tenant instances"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_token` + "`" + `
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `);
	  `, "primary key token"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_user` + "`" + `
		ADD PRIMARY KEY (` + "`" + `id` + "`" + `),
		ADD UNIQUE KEY ` + "`" + `name` + "`" + ` (` + "`" + `name` + "`" + `);
	  `, "primary key user"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_endpoint` + "`" + `
		MODIFY ` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
	  `, "primary key auto increment endpoint"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_securestore` + "`" + `
		MODIFY ` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
	  `, "primary key auto increment securestore"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_tenant` + "`" + `
		MODIFY ` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
	  `, "primary key auto increment tenant"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_tenant_instances` + "`" + `
		MODIFY ` + "`" + `id` + "`" + ` int(10) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
	  `, "primary key auto increment tenant instances"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_token` + "`" + `
		MODIFY ` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
	  `, "primary key auto increment token"},
		NamedQuery{`
	  ALTER TABLE ` + "`" + `martini_user` + "`" + `
		MODIFY ` + "`" + `id` + "`" + ` int(6) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;
	  `, "primary key auto increment user"},
	}
}
