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
	_, err := p.db.Exec("CREATE TABLE IF NOT EXISTS embeddings (id SERIAL PRIMARY KEY, embeddings vector(16));")
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresVectorDB) Save(embeddings interface{}) error {
	query := fmt.Sprintf("INSERT INTO embeddings (vector) VALUES('%s')", strings.Join(strings.Fields(fmt.Sprint(embeddings)), ","))
	log.Println(query)
	rows, err := p.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (p *PostgresVectorDB) Retrieve() error {
	// SELECT * FROM embeddings WHERE vector @@ (SELECT vector FROM embeddings WHERE id = 1) ORDER BY similarity(vector, (SELECT vector FROM embeddings WHERE id = 1)) LIMIT 5;
	return nil
}
