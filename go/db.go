package main

import "github.com/jmoiron/sqlx"

var db dbType

type dbType struct {
	withState *sqlx.DB
	noState   *sqlx.DB
}

type MySQLConnectionEnvDetail struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

type MySQLConnectionEnv struct {
	withState *MySQLConnectionEnvDetail
	noState   *MySQLConnectionEnvDetail
}

func NewMySQLConnectionEnv() (res MySQLConnectionEnv) {
	res.withState = &MySQLConnectionEnvDetail{
		Host:     getEnv("MYSQL_HOST", "10.161.12.102"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		User:     getEnv("MYSQL_USER", "isucon"),
		DBName:   getEnv("MYSQL_DBNAME", "isuumo"),
		Password: getEnv("MYSQL_PASS", "isucon"),
	}
	res.noState = &MySQLConnectionEnvDetail{
		Host:     getEnv("MYSQL_HOST", "10.161.12.103"),
		Port:     getEnv("MYSQL_PORT", "3306"),
		User:     getEnv("MYSQL_USER", "isucon"),
		DBName:   getEnv("MYSQL_DBNAME", "isuumo"),
		Password: getEnv("MYSQL_PASS", "isucon"),
	}
	return res
}
