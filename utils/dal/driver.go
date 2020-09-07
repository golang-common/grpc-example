package dal

import (
	"fmt"
	"github.com/crane/utils/encrypt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Database struct {
	DB *sqlx.DB
}

type DatabaseConfig struct {
	Driver      string
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	MaxOpenConn int
}

func connDB(cfg *DatabaseConfig) (*Database, error) {
	var DbInstance = Database{}
	passwd, err := encrypt.DesDeCrypt(cfg.Password)
	if err != nil {
		return nil, err
	}
	connStr := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%d  sslmode=disable",
		cfg.DBName, cfg.Host, cfg.User, passwd, cfg.Port)
	db, err := sqlx.Open(cfg.Driver, connStr)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	if err != nil || db == nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil || db == nil {
		return nil, err
	}
	log.Printf("DB/Connect: Connected to %s", cfg.Host)
	DbInstance.DB = db
	return &DbInstance, nil
}

func GetDB(cfg *DatabaseConfig) (*Database, error) {
	d, err := connDB(cfg)
	if err != nil {
		log.Printf("DB/Connect:%s", err.Error())
		return nil, err
	}
	err = d.DB.Ping()
	if err != nil {
		log.Printf("DB/Connect:%s", err)
		return nil, err
	}
	return d, nil
}
