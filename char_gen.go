package spg

import (
	"fmt"
	"strings"
)

// Character types for Character and Separator generation
const ( // character types
	CTUpper      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CTLower      = "abcdefghijklmnopqrstuvwxyz"
	CTDigits     = "0123456789"
	CTAmbiguous  = "0O1Il5S"
	CTSymbols    = "!#%)*+,-.:=>?@]^_}~"
	CTWhiteSpace = " \t"
)

// CTFlag is the type for the be
type CTFlag uint32

// Character type flags
const (
	Uppers CTFlag = 1 << iota
	Lowers
	Digits
	Symbols
	Ambiguous
	WhiteSpace

	Letters = Uppers | Lowers
)

// charTypesByFlag
var charTypeByFlag = map[CTFlag]string{
	Uppers:     CTUpper,
	Lowers:     CTLower,
	Digits:     CTDigits,
	Symbols:    CTSymbols,
	Ambiguous:  CTAmbiguous,
	WhiteSpace: CTWhiteSpace,
}

/*** Character type passwords ***/

// Generate a password using the character generator. The attributes contain
// all of the details needed for generating the password
func (r CharRecipe) Generate() (*Password, error) {

	if r.Length < 1 {
		return nil, fmt.Errorf("don't ask for passwords of length %d", r.Length)
	}

	p := &Password{}
	chars := r.buildCharacterList()

	toks := make([]Token, r.Length)
	for i := 0; i < r.Length; i++ {
		c := chars[Int31n(uint32(len(chars)))]
		toks[i] = Token{c, AtomTokenType}
	}
	p.Tokens = toks
	p.Entropy = r.Entropy()
	return p, nil
}

// buildCharacterList constructs the "alphabet" that is all and only those
// characters (actually strings of length 1) that are all and only those
// characters from which the password will be build. It also ensures that
// there are no duplicates
func (r CharRecipe) buildCharacterList() []string {

	ab := r.IncludeExtra
	exclude := r.ExcludeExtra
	for f, ct := range charTypeByFlag {
		if r.Allow&f == f {
			ab += ct
		}
		// Treat Require as Allow for now
		if r.Require&f == f {
			ab += ct
		}
		if r.Exclude&f == f {
			exclude += ct
		}
	}

	alphabet := subtractString(ab, exclude)
	return strings.Split(alphabet, "")
}

// Entropy returns the entropy of a character password given the generator attributes
func (r CharRecipe) Entropy() float32 {
	size := len(r.buildCharacterList())
	return float32(entropySimple(r.Length, size))
}

// CharInclusion holds the inclusion/exclusion value for some character class
type CharInclusion int

// CI{Included,Required,Excluded,Unstated} indicate how some class of characters (such as digts)
// are to be included (or not) in the generated password
const (
	CIUnstated = iota // Not included by this statement, but not excluded either
	CIAllow           // Allowed in the generated password
	CIRequire         // At least one of these must be in each generated password
	CIExclude         // None of these may appear in a generated password
)

// CharRecipe are generator attributes relevent for character list generation
type CharRecipe struct {
	Length       int    // Length of generated password in characters
	Allow        CTFlag // Flags for which character types to allow
	Require      CTFlag // Flags for which character types to require
	Exclude      CTFlag // Flags for which character types to exclude
	ExcludeExtra string // Specific characters caller may want excluded
	IncludeExtra string // Specific characters caller may want excluded (this is where to put emojis. Please don't)
}

// We need a way to map certain field names to the alphabets they correspond to
// I got worried about keeping this in sync with CharRecipe, so there's a test
// for that.
var fieldNamesAlphabets = map[string]string{
	"Uppers":    CTUpper,
	"Lowers":    CTLower,
	"Digits":    CTDigits,
	"Symbols":   CTSymbols,
	"Ambiguous": CTAmbiguous,
}

// NewCharRecipe creates CharRecipe with reasonable defaults and Length length
// more structure
func NewCharRecipe(length int) *CharRecipe {

	r := new(CharRecipe)
	r.Length = length

	r.Allow = Letters | Digits | Symbols
	r.Exclude = Ambiguous

	return r
}
