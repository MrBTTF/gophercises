package phone

import (
	"database/sql"
)

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) FindPhone(number string) (*Phone, error) {
	var p Phone
	err := db.db.QueryRow("SELECT value FROM phone_numbers WHERE value=$1", number).Scan(&p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (db *DB) UpdatePhone(p *Phone) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(statement, p.Id, p.Number)
	return err
}

func (db *DB) DeletePhone(id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.db.Exec(statement, id)
	return err
}

func insertPhone(db *sql.DB, Phone string) (int, error) {
	var id int
	statement := "INSERT INTO phone_numbers(value) VALUES($1) RETURNING id"
	err := db.QueryRow(statement, Phone).Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

type Phone struct {
	Id     int
	Number string
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.Id, &p.Number); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return nil
}

func Reset(driverName, dataSourceName, dbName string) error {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL, 
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	return err
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT value FROM phone_numbers WHERE id=$1", id).Scan(&number)
	if err != nil {
		return "", err
	}
	return number, nil
}
