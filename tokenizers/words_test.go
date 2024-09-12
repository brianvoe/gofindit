package tokenizers

import (
	"reflect"
	"testing"
)

func TestWordsProcess(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty",
			text:    "",
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "simple",
			text:    "Hello, World!",
			want:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "complex",
			text:    "Hello, World! How are you?",
			want:    []string{"hello", "world", "how", "are", "you"},
			wantErr: false,
		},
		{
			name:    "apostrophe",
			text:    "It's not about what's been lost, but what's yet to be found.",
			want:    []string{"its", "not", "about", "whats", "been", "lost", "but", "whats", "yet", "to", "be", "found"},
			wantErr: false,
		},
		{
			name:    "accented",
			text:    "Héllö, Wörld! Hów áré ÿöü?",
			want:    []string{"hello", "world", "how", "are", "you"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := NewWords()
			if err := words.Process(tt.text); (err != nil) != tt.wantErr {
				t.Errorf("Words.Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(words.words, tt.want) {
				t.Errorf("NewStandard() = %v, want %v", words.words, tt.want)
			}
		})
	}

}

func TestWordsToSearch(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		want    []string
		wantErr bool
	}{
		{
			name:    "empty",
			text:    "",
			want:    []string{},
			wantErr: true,
		},
		{
			name:    "simple",
			text:    "Hello, World!",
			want:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "complex",
			text:    "Hello, World! How are you?",
			want:    []string{"hello", "world", "how", "are", "you"},
			wantErr: false,
		},
		{
			name:    "apostrophe",
			text:    "It's not about what's been lost, but what's yet to be found.",
			want:    []string{"its", "not", "about", "whats", "been", "lost", "but", "whats", "yet", "to", "be", "found"},
			wantErr: false,
		},
		{
			name:    "accented",
			text:    "Héllö, Wörld! Hów áré ÿöü?",
			want:    []string{"hello", "world", "how", "are", "you"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := NewWords()
			got, err := words.ToSearch(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Words.Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStandard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordsSearch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		search   []string
		found    bool
		foundErr bool
	}{
		{
			name:     "simple",
			text:     "Hello, World!",
			search:   []string{"hello"},
			found:    true,
			foundErr: false,
		},
		{
			name:     "complex",
			text:     "Hello, World! How are you?",
			search:   []string{"hello", "world"},
			found:    true,
			foundErr: false,
		},
		{
			name:     "apostrophe",
			text:     "It's not about what's been lost, but what's yet to be found.",
			search:   []string{"its", "not", "about", "whats", "been", "lost", "but", "whats", "yet", "to", "be", "found"},
			found:    true,
			foundErr: false,
		},
		{
			name:     "accented",
			text:     "Héllö, Wörld! Hów áré ÿöü?",
			search:   []string{"hello", "world", "how", "are", "you"},
			found:    true,
			foundErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words := NewWords()
			if err := words.Process(tt.text); err != nil {
				t.Errorf("Words.Process() error = %v", err)
				return
			}

			got, err := words.Search(tt.search)
			if (err != nil) != tt.foundErr {
				t.Errorf("Words.Search() error = %v, foundErr %v", err, tt.foundErr)
				return
			}

			if !tt.foundErr && got != tt.found {
				t.Errorf("Words.Search() = %v, found %v", got, tt.found)
			}
		})
	}
}
