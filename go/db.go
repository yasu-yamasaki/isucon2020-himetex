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
		Host:     "10.161.12.102",
		Port:     "3306",
		User:     "isucon",
		DBName:   "isuumo",
		Password: "isucon",
	}
	res.noState = &MySQLConnectionEnvDetail{
		Host:     "localhost",
		Port:     "3306",
		User:     "isucon",
		DBName:   "isuumo",
		Password: "isucon",
	}
	return res
}
