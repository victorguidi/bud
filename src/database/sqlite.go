package database

import (
	"database/sql"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type SqlDB[T any] struct {
	db *sql.DB
}

type Engine struct{}

type ISqlDB[T any] interface {
	Init() error
	Get(id string, obj *T) error
	GetAll(objs *[]T) error
	Insert(obj *T) error
	Update(id string, obj *T) error
	DeleteOne(id string) error
	Clean() error
	Run(query string) error
}

type TableDirs struct {
	Id         string
	Dir        string
	Created_at string
	Updated_at string
}

func NewSqlDB[T any]() ISqlDB[T] {
	db, err := sql.Open("sqlite3", filepath.Join("data", "userdata.db"))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return SqlDB[T]{
		db: db,
	}
}

// TODO: REGISTER A NEW WORKER ON THE TABLE AUTOMATIC
func (d SqlDB[T]) Init() error {
	if _, err := d.db.Exec(`
    CREATE TABLE IF NOT EXISTS dirs (
    id VARCHAR(128) PRIMARY KEY,
    dir VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);

    CREATE TABLE IF NOT EXISTS usersconfig (
    id VARCHAR(128) PRIMARY KEY, 
    username VARCHAR(256),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);

    CREATE TABLE IF NOT EXISTS workers (
    id VARCHAR(128) PRIMARY KEY,
    worker VARCHAR(128),
    ollamaModel VARCHAR(128),
    embeddingModel VARCHAR(128),
    tokens INTEGER,
    workerstate INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);

    `); err != nil {
		return err
	}
	return nil
}

// func (d *SqlDB) InsertDirs(path string) error {
// 	dirId, err := uuid.NewUUID()
// 	if err != nil {
// 		return err
// 	}
//
// 	query := "INSERT INTO dirs (id, dir) VALUES($1, $2)"
// 	_, err = d.db.Exec(query, dirId, path)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (d *SqlDB) SelectDirs() ([]TableDirs, error) {
// 	query := "SELECT * FROM dirs"
// 	// query := "SELECT * FROM dirs WHERE created_at >= datetime(CURRENT_TIMESTAMP, '-30 seconds') OR updated_at >= datetime(CURRENT_TIMESTAMP, '-30 seconds');"
// 	rows, err := d.db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	response := make([]TableDirs, 0)
// 	for rows.Next() {
// 		var resp TableDirs
// 		err := rows.Scan(&resp.Id, &resp.Dir, &resp.Created_at, &resp.Updated_at)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response = append(response, resp)
// 	}
//
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return response, nil
// }
//
// func (d *SqlDB) SelectDir(dirname string) (TableDirs, error) {
// 	return TableDirs{}, nil
// }

func (d SqlDB[T]) Run(query string) error {
	_, err := d.db.Exec(query)
	return err
}

func (d SqlDB[T]) Get(id string, obj *T) error {
	query := "SELECT * FROM table WHERE id = ?"
	return d.querySingle(query, obj, id)
}

func (d SqlDB[T]) GetAll(objs *[]T) error {
	query := "SELECT * FROM table"
	return d.queryMultiple(query, objs)
}

func (d SqlDB[T]) Insert(obj *T) error {
	// Example insert query
	query := "INSERT INTO table (columns) VALUES (values)"
	// You will need to fill the values and columns accordingly
	_, err := d.db.Exec(query)
	return err
}

func (d SqlDB[T]) Update(id string, obj *T) error {
	// Example update query
	query := "UPDATE table SET columns = values WHERE id = ?"
	// You will need to fill the columns and values accordingly
	_, err := d.db.Exec(query, id)
	return err
}

func (d SqlDB[T]) DeleteOne(id string) error {
	query := "DELETE FROM table WHERE id = ?"
	_, err := d.db.Exec(query, id)
	return err
}

func (d SqlDB[T]) Clean() error {
	query := "DELETE FROM table"
	_, err := d.db.Exec(query)
	return err
}

func (d SqlDB[T]) querySingle(query string, obj *T, args ...interface{}) error {
	row := d.db.QueryRow(query, args...)
	return row.Scan(obj)
}

func (d SqlDB[T]) queryMultiple(query string, objs *[]T, args ...interface{}) error {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var obj T
		err := rows.Scan(&obj)
		if err != nil {
			return err
		}
		*objs = append(*objs, obj)
	}
	return rows.Err()
}
