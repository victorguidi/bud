package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var (
	DBUSER     = os.Getenv("DBUSER")
	DBPASSWORD = os.Getenv("DBPASSWORD")
	DBHOST     = os.Getenv("DBHOST")
	DBPORT     = os.Getenv("DBPORT")
	connStr    = fmt.Sprintf("user=%s dbname=budvecs password=%s host=%s port=%s", DBUSER, DBPASSWORD, DBHOST, DBPORT)
)

type PostgresVectorDB struct {
	db *sql.DB
}

func New() *PostgresVectorDB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	return &PostgresVectorDB{
		db: db,
	}
}

// TODO: This needs to be dynamic for each Model used, since each has a certain ammount of dimmensions
func (p *PostgresVectorDB) Initialize() error {
	rows, err := p.db.Query("CREATE TABLE embeddings (id SERIAL PRIMARY KEY, vector pgvector(16));")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

func (p *PostgresVectorDB) Save(embeddings []float64) error {
	// Execute a query
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
