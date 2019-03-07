package modules

type Discovery struct {
	id   string
	name string
}

type DiscoveryBuilder struct {
	*Discovery
}

func NewDiscoveryBuilder() *DiscoveryBuilder {
	return &DiscoveryBuilder{
		Discovery: &Discovery{},
	}
}

func (db *DiscoveryBuilder) WithID(id string) *DiscoveryBuilder {
	db.id = id
	return db
}

func (db *DiscoveryBuilder) WithName(name string) *DiscoveryBuilder {
	db.name = name
	return db
}
