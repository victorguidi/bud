package database

type IVectorDB interface {
	Initialize() error
	Save(embeddings interface{}) error
	Retrieve(embeddings interface{}) (VectorsTable, error)
}

type VectorsTable struct {
	DocName    string
	Text       string
	Created_at string
	Vector     []uint8
	Id         int
}
