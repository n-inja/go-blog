package utils

import (
	"database/sql"
	"errors"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB = nil
var RegexProjectName *regexp.Regexp

func open() error {
	address := ""
	if os.Getenv("DATABASE_ADDRESS") != "" {
		address = os.Getenv("DATABASE_ADDRESS")
	}
	userName := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")

	var err error
	DB, err = sql.Open("mysql", userName+":"+password+"@"+address+"/"+databaseName)
	if err != nil {
		return err
	}
	DB.SetMaxIdleConns(0)
	RegexProjectName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

	return initDB()
}

func Init() error {
	if DB == nil {
		err := open()
		if err != nil {
			return err
		}
	}
	return initDB()
}

func Open() error {
	if DB == nil {
		err := open()
		if err != nil {
			return err
		}
	}
	return nil
}

func Close() {
	DB.Close()
}

func initDB() error {
	rows, err := DB.Query("show tables like 'users'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		rows.Close()
		return errors.New("users table not found")
	}
	rows.Close()

	rows, err = DB.Query("show tables like 'projects'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = DB.Exec("create table projects (id varchar(20) NOT NULL PRIMARY KEY, name varchar(32) NOT NULL UNIQUE, display_name varchar(64) unicode NOT NULL, user_id varchar(32) NOT NULL, description text unicode NULL, foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = DB.Query("show tables like 'member'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = DB.Exec("create table member (user_id varchar(32) NOT NULL, project_id varchar(20) NOT NULL, PRIMARY KEY(user_id, project_id), index(user_id), index(project_id), foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = DB.Query("show tables like 'posts'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = DB.Exec("create table posts (id varchar(20) NOT NULL PRIMARY KEY, title text unicode NOT NULL, content longtext unicode NOT NULL, thumb_src varchar(64) NULL, user_id varchar(32) NOT NULL, number int NOT NULL, created_at timestamp NOT NULL, updated_at timestamp NOT NULL, project_id varchar(20) NOT NULL, views int NOT NULL, is_deleted boolean NOT NULL, index(user_id), index(created_at), index(is_deleted), index(project_id), unique(project_id, number), foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = DB.Query("show tables like 'comments'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = DB.Exec("create table comments (id varchar(20) NOT NULL PRIMARY KEY, content text unicode NOT NULL, user_id varchar(32) NOT NULL, post_id varchar(20) NOT NULL, created_at timestamp NOT NULL, is_deleted boolean NOT NULL, index(user_id), index(created_at), index(is_deleted), index(post_id), foreign key(post_id) references posts(id), foreign key(user_id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	rows, err = DB.Query("show tables like 'profiles'")
	if err != nil {
		return err
	}
	if !rows.Next() {
		_, err = DB.Exec("create table profiles (id varchar(32) NOT NULL PRIMARY KEY, description text unicode NULL, twitter_id varchar(32) NULL, github_id varchar(64) NULL, icon_src varchar(64) NULL, foreign key(id) references users(id)) engine=innodb")
		if err != nil {
			return err
		}
	}
	rows.Close()

	return nil
}

func Transact(txFunc func(*sql.Tx) error) (err error) {
	if DB == nil {
		open()
	}

	tx, err := DB.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}

func HasAuth(ID string) bool {
	if DB == nil {
		open()
	}

	var auth string
	err := DB.QueryRow("select auth from users where id = ?", ID).Scan(&auth)
	if err != nil {
		return false
	}
	if auth == "default" {
		go checkProfile(ID)
		return true
	} else {
		return false
	}
}

func HasCommentAuth(ID string) bool {
	if DB == nil {
		open()
	}

	rows, err := DB.Query("select id from users where id = ?", ID)
	if err != nil {
		return false
	}
	defer rows.Close()
	return rows.Next()
}

func checkProfile(ID string) {
	if DB == nil {
		open()
	}

	rows, err := DB.Query("select id from profiles where id = ?", ID)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		rows.Close()
		DB.Query("insert into profiles (id) value(?)", ID)
		return
	}
}
