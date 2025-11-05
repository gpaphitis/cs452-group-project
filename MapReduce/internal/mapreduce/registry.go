package mapreduce

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type Registry struct {
	mapFuncs    map[string]MapFunc
	reduceFuncs map[string]ReduceFunc
}

func NewRegistry() *Registry {
	return &Registry{
		mapFuncs:    make(map[string]MapFunc),
		reduceFuncs: make(map[string]ReduceFunc),
	}
}

func (registry *Registry) RegisterMap(name string, mapFunction MapFunc) {

	registry.mapFuncs[name] = mapFunction
}

func (registry *Registry) RegisterReduce(name string, reduceFunction ReduceFunc) {
	registry.reduceFuncs[name] = reduceFunction
}

func (registry *Registry) GetMap(name string) (MapFunc, bool) {
	fn, ok := registry.mapFuncs[name]
	return fn, ok
}

func (registry *Registry) GetReduce(name string) (ReduceFunc, bool) {
	fn, ok := registry.reduceFuncs[name]
	return fn, ok
}

func MapReverseIndex(filename string, contents string) []KeyValue {
	// Emit each (word, filename) ONCE per file to avoid duplicates.
	seenWords := make(map[string]struct{})

	// Tokenize: keep sequences of letters as words (Unicode-aware), lowercased.
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		w := strings.ToLower(b.String())
		seenWords[w] = struct{}{}
		b.Reset()
	}
	for _, r := range contents {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()

	out := make([]KeyValue, 0, len(seenWords))
	for w := range seenWords {
		out = append(out, KeyValue{Key: w, Value: filename})
	}
	return out
}

func ReduceReverseIndex(word string, values []string) string {
	// Deduplicate filenames, sort for stable output, and return:
	// "<count> file1,file2,..."
	uniq := make(map[string]struct{}, len(values))
	for _, f := range values {
		if f == "" {
			continue
		}
		uniq[f] = struct{}{}
	}
	files := make([]string, 0, len(uniq))
	for f := range uniq {
		files = append(files, f)
	}
	sort.Strings(files)
	return fmt.Sprintf("%d %s", len(files), strings.Join(files, ","))
}

// Optional: register these names for your framework to reference.
func (registry *Registry) PopulateRegistry() {
	registry.RegisterMap("inverseindex/map", MapReverseIndex)
	registry.RegisterReduce("inverseindex/reduce", ReduceReverseIndex)
}
