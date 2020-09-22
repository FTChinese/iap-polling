package config

import (
	"github.com/FTChinese/go-rest/connect"
	"github.com/spf13/viper"
	"log"
)

func GetConn(key string) (connect.Connect, error) {
	var conn connect.Connect
	err := viper.UnmarshalKey(key, &conn)
	if err != nil {
		return connect.Connect{}, err
	}

	return conn, nil
}

func MustDBConn(prod bool) connect.Connect {
	var conn connect.Connect
	var err error

	if prod {
		conn, err = GetConn("mysql.master")
	} else {
		conn, err = GetConn("mysql.dev")
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Using mysql server %s. Production: %t", conn.Host, prod)

	return conn
}

func MustHanqiConn() connect.Connect {
	conn, err := GetConn("email.hanqi")
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
