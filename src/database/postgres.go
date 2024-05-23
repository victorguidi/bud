package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var connStr string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	connStr = fmt.Sprintf("user=%s dbname=budvecs password=%s host=%s port=%s sslmode=disable",
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWORD"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"))
}

type PostgresVectorDB struct {
	db *sql.DB
}

// FIX: Vector DB in postgres not working
func New() *PostgresVectorDB {
	log.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	return &PostgresVectorDB{
		db: db,
	}
}

// TODO: This needs to be dynamic for each Model used, since each has a certain ammount of dimmensions
func (p *PostgresVectorDB) Initialize() error {
	_, err := p.db.Exec("CREATE TABLE IF NOT EXISTS embeddings (id SERIAL PRIMARY KEY, docName varchar(255), text TEXT, embeddings vector(384), created_at TIMESTAMP DEFAULT now());")
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresVectorDB) Save(docName, content string, embeddings interface{}) error {
	query := fmt.Sprintf("INSERT INTO embeddings (docName, text, embeddings) VALUES('%s', '%s', '%s'::vector)", docName, content, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(embeddings)), ","), "{}"))
	_, err := p.db.Query(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresVectorDB) Retrieve(embeddings interface{}) (*VectorsTable, error) {
	// Format embeddings as a string for the query
	embeddingStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(embeddings)), ","), "{}")

	// Correct query formation
	query := fmt.Sprintf("SELECT * FROM embeddings ORDER BY embeddings <-> '%s'::vector LIMIT 5;", embeddingStr)

	// Execute query
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vectorTable VectorsTable

	// Process the result set
	for rows.Next() {
		// Scan each row into your variables
		if err := rows.Scan(&vectorTable.Id, &vectorTable.DocName, &vectorTable.Text, &vectorTable.Vector, &vectorTable.Created_at); err != nil {
			return nil, err
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &vectorTable, nil
}
