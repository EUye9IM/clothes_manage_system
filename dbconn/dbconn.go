package dbconn

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

func Connect(db_ip string, db_port int, db_name string, db_user string, db_password string) error {
	var err error
	if db != nil {
		Close()
	}
	conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db_user, db_password, db_ip, db_port, db_name)
	db, err = sql.Open("mysql", conn_str)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		db.Close()
		return err
	}
	return nil
}

func Close() {
	if db != nil {
		db.Close()
	}
	db = nil
}

func Login(uname string, upasswd string) (bool, int, int) {
	var err error
	var uid, ugrant int
	var usalt, upw string
	err = db.QueryRow("SELECT u_id, u_salt, u_pw, u_grant FROM `user` WHERE u_name = ?", uname).Scan(&uid, &usalt, &upw, &ugrant)
	if err != nil {
		return false, 0, 0
	}
	passStr := getPassHex(usalt, upasswd)
	if passStr != upw {
		return false, 0, 0
	}
	return true, uid, ugrant
}
