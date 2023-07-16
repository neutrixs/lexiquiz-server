package db

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/exp/slices"
)

type field struct {
	Name string
	Type string
} 

type dbStructureType struct {
	User_auth []field
	Discord_login []field
}

var dbStructure = dbStructureType{
	User_auth: []field{
		{"state", "tinytext"},
		{"login_provider", "tinytext"},
		{"refresh_token", "mediumtext"},
		{"email", "mediumtext"},
		{"timestamp", "bigint"},
	},
	Discord_login: []field{
		{"state", "tinytext"},
		{"scopes", "mediumtext"},
	},
}

func createTable(name string, fields []field) {
	var keyValPair []string
	db := GetDB()

	for _, data := range fields {
		combined := data.Name + " " + data.Type
		keyValPair = append(keyValPair, combined)
	}

	_, err := db.Query(fmt.Sprintf("CREATE TABLE %s (%s)", name, strings.Join(keyValPair, ",")))
	if err != nil {
		log.Println(err)
	}
}

func dbCheck() {
	v := reflect.ValueOf(dbStructure)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name := strings.ToLower(t.Field(i).Name)
		var value []field

		if val, ok := v.Field(i).Interface().([]field); ok {
			value = val
		}

		db := GetDB()
		rows, err := db.Query(fmt.Sprintf("DESCRIBE %s", name))
		if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1146 {
			createTable(name, value)
		}

		var remoteValue []field
		for rows.Next() {
			var data field
			err := rows.Scan(&data.Name, &data.Type, &sql.NullString{}, &sql.NullString{}, &sql.NullString{}, &sql.NullString{})
			if err != nil {
				log.Println(err)
			}

			remoteValue = append(remoteValue, data)
		}

		for _, currentValue := range value {
			remoteIndex := slices.IndexFunc(remoteValue, func(search field) bool {return search.Name == currentValue.Name})
			if remoteIndex == -1 {
				fieldType := currentValue.Name + " " + currentValue.Type
				db.Query(fmt.Sprintf("ALTER TABLE %s ADD %s", name, fieldType))
				continue
			}

			currentRV := remoteValue[remoteIndex]
			if currentRV.Type != currentValue.Type {
				fieldType := currentValue.Name + " " + currentValue.Type
				db.Query(fmt.Sprintf("ALTER TABLE %s MODIFY %s", name, fieldType))
			}
		}
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}