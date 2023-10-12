package openalex

type Entity interface {
	GetType() string
	GetID() string
}
