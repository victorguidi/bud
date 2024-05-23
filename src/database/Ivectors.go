package database

type IVectorDB interface {
	Initialize() error
	Save(embeddings []float64) error
	Retrieve() error
}
