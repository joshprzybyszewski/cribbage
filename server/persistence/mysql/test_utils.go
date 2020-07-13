// +build !prod

package mysql

func GetTestConfig() Config {
	return Config{
		DSNUser:        `root`,                 // github actions use root user with a root password
		DSNPassword:    `githubactionpassword`, // defined as envvar in the go_tests.yaml
		DSNHost:        `localhost`,
		DSNPort:        3306,
		DatabaseName:   `testing_cribbage`,
		DSNParams:      `parseTime=true`,
		RunCreateStmts: true,
	}
}

func GetTestConfigForLocal() Config {
	return Config{
		DSNUser:        `root`, // locally, we just use "root" with no password
		DSNPassword:    ``,
		DSNHost:        `127.0.0.1`,
		DSNPort:        3306,
		DatabaseName:   `testing_cribbage`,
		DSNParams:      `parseTime=true`,
		RunCreateStmts: true,
	}
}
