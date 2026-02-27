package ecs

import "testing"

func TestGenreIDString(t *testing.T) {
	testCases := []struct {
		genre    GenreID
		expected string
	}{
		{GenreFantasy, "fantasy"},
		{GenreSciFi, "scifi"},
		{GenreHorror, "horror"},
		{GenreCyberpunk, "cyberpunk"},
		{GenrePostApoc, "postapoc"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.genre), func(t *testing.T) {
			if tc.genre.String() != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.genre.String())
			}
		})
	}
}

func TestGenreIDIsValid(t *testing.T) {
	validGenres := []GenreID{
		GenreFantasy,
		GenreSciFi,
		GenreHorror,
		GenreCyberpunk,
		GenrePostApoc,
	}

	for _, genre := range validGenres {
		t.Run(string(genre), func(t *testing.T) {
			if !genre.IsValid() {
				t.Errorf("Expected %s to be valid", genre)
			}
		})
	}
}

func TestGenreIDIsValidInvalid(t *testing.T) {
	invalidGenres := []GenreID{
		"invalid",
		"",
		"FANTASY",
		"sci-fi",
	}

	for _, genre := range invalidGenres {
		t.Run(string(genre), func(t *testing.T) {
			if genre.IsValid() {
				t.Errorf("Expected %s to be invalid", genre)
			}
		})
	}
}

func TestGetGenreName(t *testing.T) {
	testCases := []struct {
		genre        GenreID
		expectedName string
	}{
		{GenreFantasy, "Fantasy"},
		{GenreSciFi, "Science Fiction"},
		{GenreHorror, "Horror"},
		{GenreCyberpunk, "Cyberpunk"},
		{GenrePostApoc, "Post-Apocalyptic"},
		{"invalid", "Unknown"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.genre), func(t *testing.T) {
			name := tc.genre.GetGenreName()
			if name != tc.expectedName {
				t.Errorf("Expected name %s, got %s", tc.expectedName, name)
			}
		})
	}
}

func TestGetGenreDescription(t *testing.T) {
	testCases := []struct {
		genre GenreID
	}{
		{GenreFantasy},
		{GenreSciFi},
		{GenreHorror},
		{GenreCyberpunk},
		{GenrePostApoc},
	}

	for _, tc := range testCases {
		t.Run(string(tc.genre), func(t *testing.T) {
			desc := tc.genre.GetGenreDescription()
			if desc == "" {
				t.Errorf("Expected non-empty description for %s", tc.genre)
			}
			if desc == "Unknown world" {
				t.Errorf("Expected valid description for %s, got 'Unknown world'", tc.genre)
			}
		})
	}
}

func TestGetGenreDescriptionInvalid(t *testing.T) {
	invalidGenre := GenreID("invalid")
	desc := invalidGenre.GetGenreDescription()
	if desc != "Unknown world" {
		t.Errorf("Expected 'Unknown world' for invalid genre, got %s", desc)
	}
}

func TestAllGenres(t *testing.T) {
	genres := AllGenres()

	expectedCount := 5
	if len(genres) != expectedCount {
		t.Errorf("Expected %d genres, got %d", expectedCount, len(genres))
	}

	expectedGenres := map[GenreID]bool{
		GenreFantasy:   false,
		GenreSciFi:     false,
		GenreHorror:    false,
		GenreCyberpunk: false,
		GenrePostApoc:  false,
	}

	for _, genre := range genres {
		if _, exists := expectedGenres[genre]; !exists {
			t.Errorf("Unexpected genre in AllGenres: %s", genre)
		}
		expectedGenres[genre] = true
	}

	for genre, found := range expectedGenres {
		if !found {
			t.Errorf("Expected genre %s not found in AllGenres", genre)
		}
	}
}

func TestDefaultGenre(t *testing.T) {
	defaultGenre := DefaultGenre()

	if defaultGenre != GenreFantasy {
		t.Errorf("Expected default genre to be 'fantasy', got %s", defaultGenre)
	}

	if !defaultGenre.IsValid() {
		t.Error("Expected default genre to be valid")
	}
}

func TestGenreIDDeterminism(t *testing.T) {
	// Ensure genre constants are deterministic
	if GenreFantasy != "fantasy" {
		t.Errorf("Expected GenreFantasy to be 'fantasy', got %s", GenreFantasy)
	}

	if GenreSciFi != "scifi" {
		t.Errorf("Expected GenreSciFi to be 'scifi', got %s", GenreSciFi)
	}

	if GenreHorror != "horror" {
		t.Errorf("Expected GenreHorror to be 'horror', got %s", GenreHorror)
	}

	if GenreCyberpunk != "cyberpunk" {
		t.Errorf("Expected GenreCyberpunk to be 'cyberpunk', got %s", GenreCyberpunk)
	}

	if GenrePostApoc != "postapoc" {
		t.Errorf("Expected GenrePostApoc to be 'postapoc', got %s", GenrePostApoc)
	}
}
