package database

type IVectorDB interface {
	Initialize() error
	Save(embeddings interface{}) error
	Retrieve(embeddings interface{}) (VectorsTable, error)
}

type VectorsTable struct {
	Vector     string
	DocName    string
	Text       string
	Created_at string
	Id         int
}
