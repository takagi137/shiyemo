package anonymous

import (
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("returns_non_empty_string", func(t *testing.T) {
		name := Generate()
		if name == "" {
			t.Fatal("expected non-empty anonymous name")
		}
	})

	t.Run("contains_adjective", func(t *testing.T) {
		name := Generate()
		found := false
		for _, adj := range adjectives {
			if strings.Contains(name, adj) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected name to contain an adjective, got '%s'", name)
		}
	})

	t.Run("contains_animal", func(t *testing.T) {
		name := Generate()
		found := false
		for _, animal := range animals {
			if strings.Contains(name, animal) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected name to contain an animal, got '%s'", name)
		}
	})

	t.Run("format_is_adj_animal", func(t *testing.T) {
		name := Generate()
		parts := strings.SplitN(name, "", 2)
		if len(parts) != 2 {
			t.Fatalf("expected name to have adj+animal format, got '%s'", name)
		}
	})

	t.Run("generates_different_names", func(t *testing.T) {
		names := make(map[string]bool)
		for i := 0; i < 100; i++ {
			name := Generate()
			names[name] = true
		}
		// With 15 adjectives x 15 animals = 225 possible combinations,
		// 100 generations should produce at least 5 different names
		if len(names) < 5 {
			t.Fatalf("expected at least 5 different names, got %d", len(names))
		}
	})
}

func TestNewGenerator(t *testing.T) {
	t.Run("creates_generator", func(t *testing.T) {
		gen := NewGenerator()
		if gen == nil {
			t.Fatal("expected non-nil generator")
		}
	})

	t.Run("generator_produces_valid_names", func(t *testing.T) {
		gen := NewGenerator()
		for i := 0; i < 50; i++ {
			name := gen.Generate()
			if name == "" {
				t.Fatal("expected non-empty name from generator instance")
			}
		}
	})
}

func TestGeneratorStructure(t *testing.T) {
	t.Run("adjectives_not_empty", func(t *testing.T) {
		if len(adjectives) == 0 {
			t.Fatal("expected adjectives to not be empty")
		}
	})

	t.Run("animals_not_empty", func(t *testing.T) {
		if len(animals) == 0 {
			t.Fatal("expected animals to not be empty")
		}
	})
}
