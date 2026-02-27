// Package ecs provides the core Entity Component System framework for VANIA.
package ecs

// GenreID represents a game genre/theme identifier
type GenreID string

// Genre constants define the five supported game themes
const (
	GenreFantasy   GenreID = "fantasy"
	GenreSciFi     GenreID = "scifi"
	GenreHorror    GenreID = "horror"
	GenreCyberpunk GenreID = "cyberpunk"
	GenrePostApoc  GenreID = "postapoc"
)

// String returns the string representation of the genre
func (g GenreID) String() string {
	return string(g)
}

// IsValid checks if a genre ID is valid
func (g GenreID) IsValid() bool {
	switch g {
	case GenreFantasy, GenreSciFi, GenreHorror, GenreCyberpunk, GenrePostApoc:
		return true
	default:
		return false
	}
}

// GetGenreName returns a human-readable name for the genre
func (g GenreID) GetGenreName() string {
	switch g {
	case GenreFantasy:
		return "Fantasy"
	case GenreSciFi:
		return "Science Fiction"
	case GenreHorror:
		return "Horror"
	case GenreCyberpunk:
		return "Cyberpunk"
	case GenrePostApoc:
		return "Post-Apocalyptic"
	default:
		return "Unknown"
	}
}

// GetGenreDescription returns a brief description of the genre's world concept
func (g GenreID) GetGenreDescription() string {
	switch g {
	case GenreFantasy:
		return "Enchanted castle with magical barriers and vine-covered doorways"
	case GenreSciFi:
		return "Derelict space hulk with zero-G sections and hull-breach bulkheads"
	case GenreHorror:
		return "Haunted mansion with creaking floors and spirit seals"
	case GenreCyberpunk:
		return "Megastructure with data-stream platforms and hacking-gated doors"
	case GenrePostApoc:
		return "Collapsed bunker with rubble platforms and sealed blast doors"
	default:
		return "Unknown world"
	}
}

// AllGenres returns a slice of all valid genre IDs
func AllGenres() []GenreID {
	return []GenreID{
		GenreFantasy,
		GenreSciFi,
		GenreHorror,
		GenreCyberpunk,
		GenrePostApoc,
	}
}

// DefaultGenre returns the default genre (fantasy)
func DefaultGenre() GenreID {
	return GenreFantasy
}
