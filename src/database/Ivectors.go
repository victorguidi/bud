package database

import "github.com/pgvector/pgvector-go"

type IVectorDB interface {
	Initialize() error
	Save(docName, content string, embeddings []float32) error
	Retrieve(embeddings []float32) (*VectorsTable, error)
}

type VectorsTable struct {
	DocName    string
	Text       string
	Created_at string
	Vector     pgvector.Vector
	Id         int
}
