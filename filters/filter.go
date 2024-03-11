package filters

import (
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"runtime"
	"strings"
)

type FilterFunc func([]string) ([]string, error)

var DefaultFilters = []FilterFunc{
	Lowercase,
}

func FuncsID(filters ...FilterFunc) string {
	// if no filters, return empty string
	if len(filters) == 0 {
		return ""
	}

	identifierStr := ""
	for _, filter := range filters {
		funcPtr := reflect.ValueOf(filter).Pointer()
		funcName := runtime.FuncForPC(funcPtr).Name()

		identifierStr += funcName + ";"
	}

	// Use SHA-256 and then truncate.
	hasher := sha256.New()
	hasher.Write([]byte(identifierStr))
	fullHash := hasher.Sum(nil)
	// Truncate to the first 12 characters.
	truncatedHash := hex.EncodeToString(fullHash)[:12]

	return truncatedHash
}

// Lowercase converts all tokens to lowercase
func Lowercase(tokens []string) ([]string, error) {
	out := make([]string, len(tokens))
	for i, token := range tokens {
		out[i] = strings.ToLower(token)
	}
	return out, nil
}

var stopwords = map[string]struct{}{
	"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {}, "all": {},
	"am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "arent": {}, "as": {}, "at": {},
	"be": {}, "because": {}, "been": {}, "before": {}, "being": {}, "below": {}, "between": {},
	"both": {}, "but": {}, "by": {}, "cant": {}, "cannot": {}, "could": {}, "couldnt": {},
	"did": {}, "didnt": {}, "do": {}, "does": {}, "doesnt": {}, "doing": {}, "dont": {},
	"down": {}, "during": {}, "each": {}, "few": {}, "for": {}, "from": {}, "further": {},
	"had": {}, "hadnt": {}, "has": {}, "hasnt": {}, "have": {}, "havent": {}, "having": {},
	"he": {}, "hed": {}, "hell": {}, "hes": {}, "her": {}, "here": {}, "heres": {},
	"hers": {}, "herself": {}, "him": {}, "himself": {}, "his": {}, "how": {}, "hows": {},
	"i": {}, "id": {}, "ill": {}, "im": {}, "ive": {}, "if": {}, "in": {}, "into": {},
	"is": {}, "isnt": {}, "it": {}, "its": {}, "itself": {}, "lets": {},
	"me": {}, "more": {}, "most": {}, "mustnt": {}, "my": {}, "myself": {}, "no": {},
	"nor": {}, "not": {}, "of": {}, "off": {}, "on": {}, "once": {}, "only": {}, "or": {},
	"other": {}, "ought": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {}, "over": {},
	"own": {}, "same": {}, "shant": {}, "she": {}, "shed": {}, "shell": {}, "shes": {},
	"should": {}, "shouldnt": {}, "so": {}, "some": {}, "such": {}, "than": {}, "that": {},
	"thats": {}, "the": {}, "their": {}, "theirs": {}, "them": {}, "themselves": {}, "then": {},
	"there": {}, "theres": {}, "these": {}, "they": {}, "theyd": {}, "theyll": {}, "theyre": {},
	"theyve": {}, "this": {}, "those": {}, "through": {}, "to": {}, "too": {}, "under": {},
	"until": {}, "up": {}, "very": {}, "was": {}, "wasnt": {}, "we": {}, "wed": {}, "well": {},
	"were": {}, "weve": {}, "werent": {}, "what": {}, "whats": {}, "when": {},
	"whens": {}, "where": {}, "wheres": {}, "which": {}, "while": {}, "who": {}, "whos": {},
	"whom": {}, "why": {}, "whys": {}, "with": {}, "wont": {}, "would": {}, "wouldnt": {},
	"you": {}, "youd": {}, "youll": {}, "youre": {}, "youve": {}, "your": {}, "yours": {},
	"yourself": {}, "yourselves": {},
}

func RemoveStopwords(tokens []string) ([]string, error) {
	out := make([]string, 0, len(tokens))
	for _, token := range tokens {
		if _, ok := stopwords[token]; !ok {
			out = append(out, token)
		}
	}
	return out, nil
}
