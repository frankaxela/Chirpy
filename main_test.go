package main

import "testing"

func TestCleanChirp(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase profane word",
			input: "This is a kerfuffle opinion I need to share with the world",
			want:  "This is a **** opinion I need to share with the world",
		},
		{
			name:  "uppercase and mixed case",
			input: "KERFUFFLE and Sharbert and fornax",
			want:  "**** and **** and ****",
		},
		{
			name:  "punctuation is not replaced",
			input: "Sharbert! is fine but sharbert is not",
			want:  "Sharbert! is fine but **** is not",
		},
		{
			name:  "no profanity unchanged",
			input: "I hear Mastodon is better than Chirpy. sure I'll try it",
			want:  "I hear Mastodon is better than Chirpy. sure I'll try it",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := cleanChirp(c.input)
			if got != c.want {
				t.Errorf("cleanChirp(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}
