package lege

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

// ParseOptions are options passed to a parser
type ParseOptions struct {
	BoundaryOptions []BoundaryOption
}

// BoundaryOption are boundaries to use when collecting strings
type BoundaryOption struct {
	Starts []string
	Ends   []string
}

func (options *ParseOptions) maxStartLength() (max int) {
	for _, s := range options.getAllStarts() {
		if l := utf8.RuneCountInString(s); l > max {
			max = l
		}
	}
	return max
}

func (options *ParseOptions) maxEndLength() (max int) {
	for _, s := range options.getAllEnds() {
		if l := utf8.RuneCountInString(s); l > max {
			max = l
		}
	}
	return max
}

func (options *ParseOptions) getAllStarts() []string {
	starts := make([]string, 0)
	for _, boundary := range options.BoundaryOptions {
		for _, start := range boundary.Starts {
			starts = append(starts, start)
		}
	}
	return starts
}

func (options *ParseOptions) getAllEnds() []string {
	ends := make([]string, 0)
	for _, boundary := range options.BoundaryOptions {
		for _, end := range boundary.Ends {
			ends = append(ends, end)
		}
	}
	return ends
}

func (options *ParseOptions) getCorrespondingBoundary(start string) *BoundaryOption {
	for _, boundary := range options.BoundaryOptions {
		for _, s := range boundary.Starts {
			if s == start {
				return &boundary
			}
		}
	}
	return nil
}

func (options *ParseOptions) mustGetCorrespondingBoundary(start string) *BoundaryOption {
	for _, boundary := range options.BoundaryOptions {
		for _, s := range boundary.Starts {
			if s == start {
				return &boundary
			}
		}
	}
	panic(fmt.Sprintf("boundary not found for start: %s", start))
}

// Validate checks the parse options and returns an error if they are invalid
func (options *ParseOptions) Validate() error {

	allBoundaries := options.BoundaryOptions
	allStarts := options.getAllStarts()
	allEnds := options.getAllEnds()

	if len(allBoundaries) == 0 {
		return errors.New("must supply at least one boundary")
	}

	if len(allStarts) == 0 {
		return errors.New("must supply at least one start string")
	}

	if len(allEnds) == 0 {
		return errors.New("must supply at least one end string")
	}

	for _, start := range allStarts {
		if boundary := options.getCorrespondingBoundary(start); boundary == nil {
			return fmt.Errorf("start boundary %q must have a corresponding end boundary", start)
		}
	}

	return nil
}

// Parser is used to parse a source for collections to extract
type Parser struct {
	options *ParseOptions
}

// NewParser creates a *Parser
func NewParser(options *ParseOptions) (*Parser, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
	}
	parser := &Parser{options: options}
	return parser, nil
}

// newWindow produces a window for a parser, from the boundary options
func (parser *Parser) newWindow() []rune {
	maxStartLen := parser.options.maxStartLength()
	maxEndLen := parser.options.maxEndLength()
	windowSize := 0
	if maxStartLen > maxEndLen {
		windowSize = maxStartLen
	} else {
		windowSize = maxEndLen
	}
	return make([]rune, windowSize)
}

// windowMatchesString checks if the runes in the window are equivalent to a string
func (parser *Parser) windowMatchesString(window []rune, compareTo string) (bool, string) {
	var winString string
	runesInOption := utf8.RuneCountInString(compareTo)
	if runesInOption < len(window) {
		winString = string(window[len(window)-runesInOption:])
	} else {
		winString = string(window)
	}
	return compareTo == winString, winString
}

// Parse takes a reader
func (parser *Parser) Parse(reader io.Reader) (Collections, error) {
	r := bufio.NewReader(reader)
	window := parser.newWindow()
	windowSize := len(window)
	index := 0
	lineCounter := 1
	positionCounter := 1
	collections := make(Collections, 0)
	collecting := false

	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				if collecting { // if we're still collecting at the EOF, drop the last collection
					collections = collections[:len(collections)-1]
				}
				break
			} else {
				return nil, err
			}
		}

		// fmt.Printf("%q : %q : %v : %d : %d\n", string(window), c, collecting, lineCounter, positionCounter)

		if index < windowSize { // the window needs to be initially populated
			window[index] = c
			index++
			positionCounter++
			continue
		}

		if !collecting { // if we're not collecting, we're looking for a start match
			for _, startOption := range parser.options.getAllStarts() { // find a match with any of the possible starts
				if match, _ := parser.windowMatchesString(window, startOption); match { // if the window matches a start option
					collecting = true // go into collecting mode
					boundary := parser.options.mustGetCorrespondingBoundary(startOption)
					collections = append(collections, &Collection{
						runes:    []rune{c},
						Boundary: *boundary,
						StartLocation: Location{
							Line: lineCounter,
							Pos:  positionCounter,
						},
					}) // create a new collection, starting with this rune
					break
				}
			}
		} else { // if we're collecting, we're looking for an end match and storing runes along the way
			currentCollection := collections.getLast()
			for _, endOption := range currentCollection.Boundary.Ends {
				if match, _ := parser.windowMatchesString(window, endOption); match { // if the window matches an end option
					collecting = false // leave collecting mode
					// if we're stopping collection, since the window trails the current index, we need to reslice the current collection to take off
					// the runes we just matched
					runeCount := utf8.RuneCountInString(endOption)
					currentCollection.trimRightRunes(runeCount)
					currentCollection.EndLocation = Location{
						Line: lineCounter,
						Pos:  positionCounter - runeCount - 1,
					}
					break
				}
			}
			if collecting {
				currentCollection.addRune(c)
			}
		}

		// shift the window by one rune
		for i := range window {
			if i == len(window)-1 { // if we're at the last spot in the window
				window[i] = c // assign it to the current rune
			} else { // otherwise, assign the current spot in the window to what's in the next spot
				window[i] = window[i+1]
			}
		}
		index++
		positionCounter++

		if string(c) == "\n" {
			lineCounter++
			positionCounter = 1
		}
	}

	return collections, nil
}
