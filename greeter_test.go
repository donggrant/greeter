package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/googleapis/gax-go/v2"
)

// mockTranslationClient implements a fake translation client for testing
type mockTranslationClient struct {
	translate.TranslationClient
	// Mock responses for different languages
	translations map[string]string
	// For testing error conditions
	shouldError bool
}

func (m *mockTranslationClient) TranslateText(_ context.Context, req *translatepb.TranslateTextRequest, _ ...gax.CallOption) (*translatepb.TranslateTextResponse, error) {
	if m.shouldError {
		return nil, errors.New("mock translation error")
	}

	text := req.GetContents()[0]
	lang := req.GetTargetLanguageCode()

	// Simple mock translations
	if m.translations == nil {
		m.translations = map[string]string{
			"es": "¡" + text + "!",
			"fr": text + " !",
			"ja": text + "！",
			"de": text + "!",
		}
	}

	// Return error for invalid language code
	if _, ok := m.translations[lang]; !ok {
		return nil, errors.New("unsupported language code")
	}

	translation := m.translations[lang]
	return &translatepb.TranslateTextResponse{
		Translations: []*translatepb.Translation{
			{TranslatedText: translation},
		},
	}, nil
}

func (m *mockTranslationClient) Close() error {
	return nil
}

// Helper function to create a test greeter
func newTestGreeter(recipient string, client TranslationClient) *Greeter {
	tmpfile, _ := os.CreateTemp("", "translation_cache_*.json")

	return &Greeter{
		recipient: recipient,
		language:  "en",
		ctx:       context.Background(),
		client:    client,
		cache: &TranslationCache{
			Translations: make(map[string]map[string]string),
		},
		cacheFile: tmpfile.Name(),
		projectID: "test-project",
		stats:     &Stats{}, // Initialize stats
	}
}

// TestGreeterBasicFunctionality tests the core greeting functionality
func TestGreeterBasicFunctionality(t *testing.T) {
	g := newTestGreeter("Test", &mockTranslationClient{})
	defer os.Remove(g.cacheFile)

	// Test English greeting (no translation needed)
	greeting, err := g.Greet()
	if err != nil {
		t.Errorf("English greeting failed: %v", err)
	}
	if greeting == "" {
		t.Error("English greeting was empty")
	}

	// Test Spanish translation
	g.SetLanguage("es")
	greeting, err = g.Greet()
	if err != nil {
		t.Errorf("Spanish greeting failed: %v", err)
	}
	if greeting == "" {
		t.Error("Spanish greeting was empty")
	}

	// Test cache functionality
	g.SetLanguage("es")
	greeting2, err := g.Greet()
	if err != nil {
		t.Errorf("Cached Spanish greeting failed: %v", err)
	}
	if greeting != greeting2 {
		t.Errorf("Cache not working: got different translations for same text")
	}
	if g.stats.CacheHits != 2 { // One for English, one for cached Spanish
		t.Errorf("Expected 2 cache hits, got %d", g.stats.CacheHits)
	}
}

// TestTimeBasedGreeting tests that different times produce different greetings
func TestTimeBasedGreeting(t *testing.T) {
	tests := []struct {
		hour     int
		expected string
	}{
		{8, "Good morning"},
		{13, "Good afternoon"},
		{19, "Good evening"},
		{23, "Good night"},
	}

	originalNow := timeNow
	defer func() { timeNow = originalNow }()

	g := newTestGreeter("Test", &mockTranslationClient{})
	defer os.Remove(g.cacheFile)

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			mockTime := time.Date(2024, 1, 1, tt.hour, 0, 0, 0, time.UTC)
			timeNow = func() time.Time {
				return mockTime
			}
			greeting := g.getTimeBasedGreeting()
			if greeting != tt.expected+", Test!" {
				t.Errorf("At hour %d, got %q, want %q", tt.hour, greeting, tt.expected+", Test!")
			}
		})
	}
}

// TestSetLanguage tests language switching
func TestSetLanguage(t *testing.T) {
	g := newTestGreeter("Test", &mockTranslationClient{})
	defer os.Remove(g.cacheFile)

	languages := []Language{"en", "es", "fr", "ja"}
	for _, lang := range languages {
		g.SetLanguage(lang)
		if g.language != lang {
			t.Errorf("SetLanguage(%q) failed: got %q, want %q", lang, g.language, lang)
		}
	}
}

// TestCachePersistence tests that translations are properly cached and loaded
func TestCachePersistence(t *testing.T) {
	// Create a temporary cache file
	tmpfile, err := os.CreateTemp("", "translation_cache_*.json")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// First greeter makes a translation
	g1 := newTestGreeter("Test", &mockTranslationClient{})
	g1.cacheFile = tmpfile.Name()
	defer os.Remove(g1.cacheFile)

	g1.SetLanguage("es")
	greeting1, err := g1.Greet()
	if err != nil {
		t.Fatalf("Initial translation failed: %v", err)
	}

	// Second greeter should load from cache
	g2 := newTestGreeter("Test", &mockTranslationClient{})
	g2.cacheFile = tmpfile.Name()
	defer os.Remove(g2.cacheFile)

	// Load cache from file
	if data, err := os.ReadFile(tmpfile.Name()); err == nil {
		if err := json.Unmarshal(data, &g2.cache.Translations); err != nil {
			t.Fatalf("Failed to load cache: %v", err)
		}
	}

	g2.SetLanguage("es")
	greeting2, err := g2.Greet()
	if err != nil {
		t.Fatalf("Second translation failed: %v", err)
	}

	if greeting1 != greeting2 {
		t.Errorf("Cache persistence failed: got different translations")
	}
	if g2.stats.APICalls != 0 {
		t.Errorf("Expected no API calls for cached translation, got %d", g2.stats.APICalls)
	}
}

// TestErrorHandling tests various error conditions
func TestErrorHandling(t *testing.T) {
	g := newTestGreeter("Test", &mockTranslationClient{shouldError: true})
	defer os.Remove(g.cacheFile)

	g.SetLanguage("es")
	_, err := g.Greet()
	if err == nil {
		t.Error("Expected error from translation service, got none")
	}
}

// TestInvalidLanguageCode tests handling of unsupported language codes
func TestInvalidLanguageCode(t *testing.T) {
	g := newTestGreeter("Test", &mockTranslationClient{})
	defer os.Remove(g.cacheFile)

	g.SetLanguage("xx") // Invalid language code
	_, err := g.Greet()
	if err == nil {
		t.Error("Expected error for invalid language code, got none")
	}
}

// TestCostEstimation tests the cost estimation functionality
func TestCostEstimation(t *testing.T) {
	g := newTestGreeter("Test", &mockTranslationClient{})
	defer os.Remove(g.cacheFile)

	g.SetLanguage("es")
	_, err := g.Greet()
	if err != nil {
		t.Fatalf("Greeting failed: %v", err)
	}

	// Cost should be calculated based on character count
	expectedChars := len("Good morning, Test!") // Base greeting length
	if g.stats.CharsSent != expectedChars {
		t.Errorf("Expected %d characters sent, got %d", expectedChars, g.stats.CharsSent)
	}

	// Check cost estimation (assuming $0.00002 per character)
	expectedCost := float64(expectedChars) * 0.00002
	if g.stats.CostEstimate != expectedCost {
		t.Errorf("Expected cost estimate of %f, got %f", expectedCost, g.stats.CostEstimate)
	}
}

// TestMultipleLanguageTranslations tests translations across multiple languages
func TestMultipleLanguageTranslations(t *testing.T) {
	g := newTestGreeter("Test", &mockTranslationClient{})
	defer os.Remove(g.cacheFile)

	languages := []Language{"es", "fr", "ja", "de"}
	for _, lang := range languages {
		g.SetLanguage(lang)
		greeting, err := g.Greet()
		if err != nil {
			t.Errorf("Translation failed for %s: %v", lang, err)
		}
		if greeting == "" {
			t.Errorf("Empty greeting for language %s", lang)
		}
	}

	// Check stats after multiple translations
	if g.stats.APICalls != len(languages) {
		t.Errorf("Expected %d API calls, got %d", len(languages), g.stats.APICalls)
	}
}
