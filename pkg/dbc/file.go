package dbc

// File is a parsed DBC source file.
type File struct {
	// Name of the file.
	Name string
	// Data contains the raw file data.
	Data []byte
	// Defs in the file.
	Defs []Def
}
