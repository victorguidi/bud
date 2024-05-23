package database

import "github.com/pgvector/pgvector-go"

type IVectorDB interface {
	Initialize() error
	Save(embeddings interface{}) error
	Retrieve(embeddings interface{}) (VectorsTable, error)
}

type VectorsTable struct {
	DocName    string
	Text       string
	Created_at string
	Vector     pgvector.Vector
	Id         int
}
