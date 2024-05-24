package database

import (
	"database/sql"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type SqlDB struct {
	db *sql.DB
}

type TableDirs struct {
	Id         string
	Dir        string
	Created_at string
	Updated_at string
}

func NewSqlDB() *SqlDB {
	db, err := sql.Open("sqlite3", filepath.Join("data", "userdata.db"))
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &SqlDB{
		db: db,
	}
}

func (d *SqlDB) Init() error {
	if _, err := d.db.Exec(`
    CREATE TABLE IF NOT EXISTS dirs (id VARCHAR(128) PRIMARY KEY, dir VARCHAR(100), created_at DATETIME DEFAULT CURRENT_TIMESTAMP, updated_at DATETIME DEFAULT CURRENT_TIMESTAMP);
    `); err != nil {
		return err
	}
	return nil
}

func (d *SqlDB) InsertDirs(path string) error {
	dirId, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	query := "INSERT INTO dirs (id, dir) VALUES($1, $2)"
	_, err = d.db.Exec(query, dirId, path)
	if err != nil {
		return err
	}
	return nil
}

func (d *SqlDB) SelectDirs() ([]TableDirs, error) {
	query := "SELECT * FROM dirs"
	// query := "SELECT * FROM dirs WHERE created_at >= datetime(CURRENT_TIMESTAMP, '-30 seconds') OR updated_at >= datetime(CURRENT_TIMESTAMP, '-30 seconds');"
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := make([]TableDirs, 0)
	for rows.Next() {
		var resp TableDirs
		err := rows.Scan(&resp.Id, &resp.Dir, &resp.Created_at, &resp.Updated_at)
		if err != nil {
			return nil, err
		}
		response = append(response, resp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}
