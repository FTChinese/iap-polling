package config

import (
	"github.com/FTChinese/go-rest/connect"
	"github.com/spf13/viper"
	"log"
)

func SetupViper() error {
	viper.SetConfigName("api")
	viper.AddConfigPath("$HOME/config")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}

func MustSetupViper() {
	if err := SetupViper(); err != nil {
		panic(err)
	}
}

func GetConn(key string) (connect.Connect, error) {
	var conn connect.Connect
	err := viper.UnmarshalKey(key, &conn)
	if err != nil {
		return connect.Connect{}, err
	}

	return conn, nil
}

func MustGetDBConn(prod bool) connect.Connect {
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

func MustRedisAddr(prod bool) string {
	var addr string
	if prod {
		log.Print("Using production redis")
		addr = viper.GetString("redis.production")
	} else {
		log.Print("Using development redis")
		addr = viper.GetString("redis.development")
	}

	if addr == "" {
		log.Fatal("Redis address not found")
	}

	return addr
}

func MustGetHanqiConn() connect.Connect {
	conn, err := GetConn("email.hanqi")
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func MustIAPSecret() string {
	pw := viper.GetString("apple.receipt_password")
	if pw == "" {
		panic("empty receipt verification password")
	}

	return pw
}
