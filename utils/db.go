package utils

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
)

var db *sql.DB

func Open(userName, password, address, databaseName string) error {
	var err error
	db, err = sql.Open("mysql", userName+":"+password+"@"+address+"/"+databaseName)
	if err != nil {
		return err
	}
	return initDB()
}

func Close() {
	db.Close()
}

func initDB() error {
	rows, err := db.Query("show tables like 'users'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		rows.Close()
		return errors.New("users table not found")
	}
	rows.Close()

	rows, err = db.Query("show tables like 'projects'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table projects (id varchar(20) NOT NULL PRIMARY KEY, name varchar(32) NOT NULL UNIQUE, user_id varchar(32) NOT NULL)")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'member'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table member (user_id varchar(32) NOT NULL, project_id varchar(20) NOT NULL, PRIMARY KEY(user_id, project_id))")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'user_auth'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table user_auth (user_id varchar(32) NOT NULL PRIMARY KEY, auth varchar(20) NOT NULL)")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'posts'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table posts (id varchar(20) NOT NULL PRIMARY KEY, title varchar(32) NOT NULL, content text unicode NOT NULL, thumb_src varchar(32) NULL, user_id varchar(32) NOT NULL, created_at timestamp NOT NULL, updated_at timestamp NOT NULL, project_id varchar(20) NOT NULL, views int NOT NULL, has_deleted boolean NOT NULL, index(created_at))")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = db.Query("show tables like 'comments'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = db.Exec("create table comments (id varchar(20) NOT NULL PRIMARY KEY, content text unicode NOT NULL, user_id varchar(32) NOT NULL, post_id varchar(20) NOT NULL, created_at timestamp NOT NULL, has_deleted boolean NOT NULL, index(created_at))")
		if err != nil {
			return err
		}
	}
	rows.Close()

	return nil
}
