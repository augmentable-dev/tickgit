package lege

// Location represents location in an input string
type Location struct {
	Line int
	Pos  int
}

// Collection represents a string that has been "plucked" from a source
type Collection struct {
	runes         []rune
	Boundary      BoundaryOption
	StartLocation Location
	EndLocation   Location
}

// Collections is a list of *Collection
type Collections []*Collection

func (collections Collections) getLast() *Collection {
	return collections[len(collections)-1]
}

// Strings returns each collection as a string, in a list of strings
func (collections Collections) Strings() (s []string) {
	for _, collection := range collections {
		s = append(s, collection.String())
	}
	return s
}

func (collection *Collection) addRune(r rune) {
	collection.runes = append(collection.runes, r)
}

func (collection *Collection) trimRightRunes(num int) {
	if num <= len(collection.runes) {
		collection.runes = collection.runes[:len(collection.runes)-num]
	}
}

func (collection *Collection) String() string {
	return string(collection.runes)
}
