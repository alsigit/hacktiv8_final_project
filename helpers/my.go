package helpers

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"os"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func BcryptEncode(str string) (bool, string) {
	password := []byte(str)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error occured : ", err.Error())
		return false, ""
	} else {
		return true, string(hashedPassword)
	}
}

func ValidMailAddress(address string) bool {
	_, err := mail.ParseAddress(address)
	return err == nil
}

func mapBytesToString(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			m[k] = string(b)
		}
	}
	return m
}

func DatabaseQueryRows(db *sqlx.DB, query string, args ...interface{}) []map[string]interface{} {

	var datarows []map[string]interface{}
	rows, err := db.Queryx(query, args...)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			// fmt.Println(err)
			datarows = append(datarows, mapBytesToString(results))
		}
	}
	return datarows
}

func DatabaseQuerySingleRow(db *sqlx.DB, query string, args ...interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	rows, err := db.Queryx(query, args...)

	if err != nil {
		fmt.Println("Query Error", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			results := make(map[string]interface{})
			err = rows.MapScan(results)
			return mapBytesToString(results)
		}
	}

	return result
}

func B64Decode(b64string string, filepath string) error {
	var tr_err error
	dec, err := base64.StdEncoding.DecodeString(b64string)
	if err != nil {
		fmt.Println("A")
		tr_err = err
	}

	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println("B", filepath)
		tr_err = err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		tr_err = err
		fmt.Println("C")
	}

	if err := f.Sync(); err != nil {
		tr_err = err
		fmt.Println("D")
	}
	return tr_err
}
