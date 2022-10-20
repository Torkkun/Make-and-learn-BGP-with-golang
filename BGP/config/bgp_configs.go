package config

// Globalとneighborで分かれている恐らくGlobalはnoth,southでneighborはeast,west用かな？

type Neighbor struct {
	Config NeighborConfig
	State  NeighborState
}

type NeighborConfig struct {
}

type NeighborState struct {
	Message Messages
}

type Messages struct {
}
