package tokenizers

import (
	"fmt"
	"testing"
)

func TestNGramProcess(t *testing.T) {
	type tests struct {
		min   int
		max   int
		input string
		want  map[string]bool
	}

	testCases := []tests{
		// Single Words
		{
			min:   1,
			max:   1,
			input: "hello",
			want:  map[string]bool{"h": true, "e": true, "l": true, "o": true},
		},
		{
			min:   2,
			max:   2,
			input: "hello",
			want:  map[string]bool{"he": true, "el": true, "ll": true, "lo": true},
		},
		{
			min:   3,
			max:   3,
			input: "hello",
			want:  map[string]bool{"hel": true, "ell": true, "llo": true},
		},
		// Multiple Words
		{
			min:   1,
			max:   1,
			input: "hello world",
			want:  map[string]bool{" ": true, "h": true, "e": true, "l": true, "o": true, "w": true, "r": true, "d": true},
		},
		{
			min:   2,
			max:   2,
			input: "hello world",
			want:  map[string]bool{"he": true, "el": true, "ll": true, "lo": true, "o ": true, " w": true, "wo": true, "or": true, "rl": true, "ld": true},
		},
		// Real World Examples
		{
			min:   1,
			max:   10,
			input: "Can't swim",
			want:  map[string]bool{"an't swim": true, "n't ": true, "n't sw": true, "t ": true, "sw": true, "swi": true, "'t ": true, "an't swi": true, "swim": true, "n't swim": true, "n'": true, " swim": true, "ca": true, "t swim": true, " sw": true, "s": true, "wi": true, "im": true, "can't swi": true, "an't s": true, "i": true, "an't": true, "t sw": true, "w": true, "wim": true, "'t swim": true, "t s": true, "an'": true, "n't s": true, " ": true, "can't swim": true, "'t sw": true, "t": true, "can": true, "can't": true, "can't s": true, "an": true, "an't sw": true, "n": true, "'t s": true, "c": true, "can't ": true, "an't ": true, "n't": true, "n't swi": true, "'": true, "m": true, "can'": true, "can't sw": true, "a": true, "'t": true, "t swi": true, " swi": true, "'t swi": true, " s": true},
		},
	}

	outputMap := func(t *testing.T, got map[string]bool) {
		t.Helper()

		fmt.Println("map[string]bool{")
		for k, v := range got {
			fmt.Printf("%q: %t, ", k, v)
		}
		fmt.Println("}")
	}

	for _, tc := range testCases {
		// Create a new NGram tokenizer
		n := NewNGram(tc.min, tc.max)

		// Process the input
		err := n.Process(tc.input)
		if err != nil {
			t.Errorf("NGram.Process() failed: %v", err)
		}

		// Check the index and want are the same length
		// and that the values are the same
		if len(n.index) != len(tc.want) {
			t.Errorf("NGram.Process() failed count: got %d, want %d", len(n.index), len(tc.want))
		}

		// Check the values are the same
		for k, v := range tc.want {
			if n.index[k] != v {
				t.Errorf("NGram.Process() failed: got %+v, want %v", n.index, tc.want)
				outputMap(t, n.index)
			}
		}
	}
}

func TestNGramToSearch(t *testing.T) {
	type tests struct {
		min   int
		max   int
		input string
		want  []string
	}

	testCases := []tests{
		{
			min:   1,
			max:   1,
			input: "hello",
			want:  []string{"h"},
		},
		{
			min:   2,
			max:   2,
			input: "hello",
			want:  []string{"he"},
		},
		{
			min:   3,
			max:   3,
			input: "hello",
			want:  []string{"hel"},
		},
		{
			min:   1,
			max:   1,
			input: "hello world",
			want:  []string{"h"},
		},
		{
			min:   2,
			max:   2,
			input: "he",
			want:  []string{"he"},
		},
		{
			min:   1,
			max:   10,
			input: "Can't swim",
			want:  []string{"can't swim"},
		},
	}

	for _, tc := range testCases {
		// Create a new NGram tokenizer
		n := NewNGram(tc.min, tc.max)

		// Process the input
		val, err := n.ToSearch(tc.input)
		if err != nil {
			t.Errorf("NGram.Process() failed: %v", err)
		}

		// Check the values are the same
		if len(val) != len(tc.want) {
			t.Errorf("NGram.Process() failed: got %v, want %v", val, tc.want)
		}

		for i := range val {
			if val[i] != tc.want[i] {
				t.Errorf("NGram.Process() failed: got %v, want %v", val, tc.want)
			}
		}
	}
}

func TestNGramSearch(t *testing.T) {
	type tests struct {
		min    int
		max    int
		input  string
		search string
		match  bool
	}

	testCases := []tests{
		{
			min:    1,
			max:    10,
			input:  "hello world",
			search: "hello",
			match:  true,
		},
		{
			min:    1,
			max:    10,
			input:  "hello world",
			search: "world",
			match:  true,
		},
		{
			min:    1,
			max:    10,
			input:  "hello world",
			search: "lo wo",
			match:  true,
		},
		{
			min:    1,
			max:    10,
			input:  "Billy went to the store",
			search: "went to",
			match:  true,
		},
		{
			min:    1,
			max:    10,
			input:  "I can't help but feel like I'm missing something",
			search: "can't help",
			match:  true,
		},
	}

	for _, tc := range testCases {
		// Create a new NGram tokenizer
		n := NewNGram(tc.min, tc.max)

		// Process the input
		err := n.Process(tc.input)
		if err != nil {
			t.Errorf("NGram.Process() failed: %v", err)
		}

		// Search the input
		match, err := n.Search([]string{tc.search})
		if err != nil {
			t.Errorf("NGram.Search() failed: %v", err)
		}

		if match != tc.match {
			t.Errorf("NGram.Search() failed: got %v, want %v", match, tc.match)
		}
	}
}

// BenchmarkNGramProcess benchmarks the Process method of the NGram struct
func BenchmarkNGramProcessSmall(b *testing.B) {
	n := NewNGram(1, 10)

	for i := 0; i < b.N; i++ {
		n.Process("hello world")
	}
}

func BenchmarkNGramProcessMedium(b *testing.B) {
	n := NewNGram(1, 10)

	for i := 0; i < b.N; i++ {
		n.Process("hello world, how are you doing today?")
	}
}

func BenchmarkNGramProcessLarge(b *testing.B) {
	n := NewNGram(1, 10)

	for i := 0; i < b.N; i++ {
		n.Process("hello world, how are you doing today? I'm doing well, thank you for asking.")
	}
}

// BenchmarkNGramToSearch benchmarks the ToSearch method of the NGram struct
func BenchmarkNGramToSearchSmall(b *testing.B) {
	n := NewNGram(1, 10)

	for i := 0; i < b.N; i++ {
		n.ToSearch("hello world")
	}
}

func BenchmarkNGramToSearchMedium(b *testing.B) {
	n := NewNGram(1, 10)

	for i := 0; i < b.N; i++ {
		n.ToSearch("hello world, how are you doing today?")
	}
}

func BenchmarkNGramToSearchLarge(b *testing.B) {
	n := NewNGram(1, 10)

	for i := 0; i < b.N; i++ {
		n.ToSearch("hello world, how are you doing today? I'm doing well, thank you for asking.")
	}
}

// BenchmarkNGramSearch benchmarks the Search method of the NGram struct
func BenchmarkNGramSearchSmall(b *testing.B) {
	n := NewNGram(1, 10)
	n.Process("hello world")

	for i := 0; i < b.N; i++ {
		n.Search([]string{"hello"})
	}
}

func BenchmarkNGramSearchMedium(b *testing.B) {
	n := NewNGram(1, 10)
	n.Process("hello world, how are you doing today?")

	for i := 0; i < b.N; i++ {
		n.Search([]string{"hello"})
	}
}

func BenchmarkNGramSearchLarge(b *testing.B) {
	n := NewNGram(1, 10)
	n.Process("hello world, how are you doing today? I'm doing well, thank you for asking.")

	for i := 0; i < b.N; i++ {
		n.Search([]string{"hello"})
	}
}
