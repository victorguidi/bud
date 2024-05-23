package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	pgvector "github.com/pgvector/pgvector-go"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var connStr string

// NOTE: all-minilm = 384 dimensions
// NOTE: mxbai-embed-large = 1024 demensions

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
	_, err := p.db.Exec(`
    CREATE TABLE IF NOT EXISTS embeddings (id SERIAL PRIMARY KEY, docName varchar(255) UNIQUE, text TEXT, embeddings vector(1024), created_at TIMESTAMP DEFAULT now());
    `)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresVectorDB) Save(docName, content string, embeddings []float32) error {
	// embeddingStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(embeddings)), ","), "{}")
	query := `INSERT INTO embeddings (docName, text, embeddings)
    VALUES ($1, $2, $3::vector)
    ON CONFLICT (docName)
    DO UPDATE SET
        text = EXCLUDED.text,
        embeddings = EXCLUDED.embeddings;`

	_, err := p.db.Exec(query, docName, content, pgvector.NewVector(embeddings))
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresVectorDB) Retrieve(embeddings []float32) (*VectorsTable, error) {
	// embeddingStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(embeddings)), ","), "{}")
	query := "SELECT *, embeddings <=> $1::vector AS similarity FROM embeddings ORDER BY similarity LIMIT 1;"

	rows, err := p.db.Query(query, pgvector.NewVector(embeddings))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vectorTable []VectorsTable

	for rows.Next() {
		var vTable VectorsTable
		var similarity float32
		if err := rows.Scan(&vTable.Id, &vTable.DocName, &vTable.Text, &vTable.Vector, &vTable.Created_at, &similarity); err != nil {
			return nil, err
		}

		log.Printf("\n================\nSIMILARITY: %f\n================\n", similarity)

		if similarity < 0.5 {
			vectorTable = append(vectorTable, vTable)
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(vectorTable) > 0 {
		return &vectorTable[0], nil
	}

	return &VectorsTable{}, nil
}
