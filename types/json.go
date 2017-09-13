package types

// JSONFileParser reads json from a file
type JSONFileParser interface {
	Read() error
}

// JSONFile represents a json file
type JSONFile struct {
	Path string
}
