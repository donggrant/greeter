package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	translatepb "cloud.google.com/go/translate/apiv3/translatepb"
	"github.com/googleapis/gax-go/v2"
)

// TranslationClient interface for mocking in tests
type TranslationClient interface {
	TranslateText(context.Context, *translatepb.TranslateTextRequest, ...gax.CallOption) (*translatepb.TranslateTextResponse, error)
	Close() error
}

// TranslationCache represents the cache structure
type TranslationCache struct {
	Translations map[string]map[string]string // map[sourceText]map[targetLang]translatedText
	mu           sync.RWMutex
}

// Stats tracks translation statistics
type Stats struct {
	APICalls     int     `json:"apiCalls"`     // Number of API calls made
	CharsSent    int     `json:"charsSent"`    // Number of characters sent to API
	CostEstimate float64 `json:"costEstimate"` // Estimated cost in USD
	CacheHits    int     `json:"cacheHits"`    // Number of cache hits
}

// Language represents ISO 639-1 language codes
type Language string

// timeNow allows overriding time.Now in tests
var timeNow = time.Now

// Greeter manages greetings in different languages
type Greeter struct {
	recipient string
	language  Language
	client    TranslationClient
	ctx       context.Context
	cache     *TranslationCache
	cacheFile string
	projectID string
	stats     *Stats
}

// NewGreeter creates a new Greeter with default English language
func NewGreeter(recipient string) (*Greeter, error) {
	ctx := context.Background()

	// Get project ID from environment
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT_ID environment variable is not set")
	}

	// Check for credentials
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is not set. Please set it to the path of your service account key")
	}

	client, err := translate.NewTranslationClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create translate client: %v", err)
	}

	cache := &TranslationCache{
		Translations: make(map[string]map[string]string),
	}

	// Try to load existing cache
	cacheFile := "translation_cache.json"
	if data, err := os.ReadFile(cacheFile); err == nil {
		if err := json.Unmarshal(data, &cache.Translations); err != nil {
			log.Printf("Warning: Could not load cache: %v", err)
		}
	}

	return &Greeter{
		recipient: recipient,
		language:  "en",
		client:    client,
		ctx:       ctx,
		cache:     cache,
		cacheFile: cacheFile,
		projectID: projectID,
		stats:     &Stats{}, // Initialize with zero values
	}, nil
}

// saveCache saves the translation cache to disk
func (g *Greeter) saveCache() error {
	g.cache.mu.RLock()
	defer g.cache.mu.RUnlock()

	data, err := json.MarshalIndent(g.cache.Translations, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %v", err)
	}

	if err := os.WriteFile(g.cacheFile, data, 0644); err != nil {
		return fmt.Errorf("failed to save cache: %v", err)
	}

	return nil
}

// SetLanguage changes the greeting language
func (g *Greeter) SetLanguage(lang Language) {
	g.language = lang
}

// getTimeBasedGreeting returns appropriate greeting based on time of day
func (g *Greeter) getTimeBasedGreeting() string {
	hour := timeNow().Hour()
	var greeting string

	switch {
	case hour < 12:
		greeting = "Good morning"
	case hour < 17:
		greeting = "Good afternoon"
	case hour < 22:
		greeting = "Good evening"
	default:
		greeting = "Good night"
	}

	return fmt.Sprintf("%s, %s!", greeting, g.recipient)
}

// translateGreeting translates the greeting to the target language
func (g *Greeter) translateGreeting(text string) (string, error) {
	targetLang := string(g.language)

	// Check cache first
	g.cache.mu.RLock()
	if langCache, ok := g.cache.Translations[text]; ok {
		if translation, ok := langCache[targetLang]; ok {
			g.stats.CacheHits++
			g.cache.mu.RUnlock()
			return translation, nil
		}
	}
	g.cache.mu.RUnlock()

	// If not in cache, translate using API
	g.stats.APICalls++
	g.stats.CharsSent += len(text)
	g.stats.CostEstimate += float64(len(text)) * 0.00002 // $0.00002 per character

	log.Printf("Translating text to %s", targetLang)

	req := &translatepb.TranslateTextRequest{
		Contents:           []string{text},
		TargetLanguageCode: targetLang,
		SourceLanguageCode: "en",
		MimeType:           "text/plain",
		Parent:             fmt.Sprintf("projects/%s", g.projectID),
	}

	resp, err := g.client.TranslateText(g.ctx, req)
	if err != nil {
		return "", fmt.Errorf("translation failed: %v", err)
	}

	if len(resp.GetTranslations()) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	translation := resp.GetTranslations()[0].GetTranslatedText()

	// Add to cache
	g.cache.mu.Lock()
	if g.cache.Translations[text] == nil {
		g.cache.Translations[text] = make(map[string]string)
	}
	g.cache.Translations[text][targetLang] = translation
	g.cache.mu.Unlock()

	// Save cache to disk
	if err := g.saveCache(); err != nil {
		log.Printf("Warning: Failed to save cache: %v", err)
	}

	return translation, nil
}

// Greet returns a greeting in the current language
func (g *Greeter) Greet() (string, error) {
	greeting := g.getTimeBasedGreeting()

	if string(g.language) != "en" {
		translated, err := g.translateGreeting(greeting)
		if err != nil {
			return "", fmt.Errorf("translation failed: %v", err)
		}
		greeting = translated
	} else {
		// Count English as a cache hit since we're using our built-in English templates
		g.stats.CacheHits++
	}

	return greeting, nil
}

// Close cleans up resources used by the Greeter
func (g *Greeter) Close() error {
	return g.client.Close()
}

// RunCLI runs the greeter in command-line mode
func RunCLI() {
	log.SetFlags(log.Ltime) // Only show time in logs

	// Validate command line arguments
	if len(os.Args) != 3 {
		log.Fatal("Usage: greeter <recipient> <language-code>")
	}
	recipient := os.Args[1]
	lang := Language(os.Args[2])

	// Create a new greeter
	greeter, err := NewGreeter(recipient)
	if err != nil {
		log.Fatalf("Failed to create greeter: %v", err)
	}
	defer greeter.Close()

	// Set language and get greeting
	greeter.SetLanguage(lang)
	greeting, err := greeter.Greet()
	if err != nil {
		log.Fatalf("Error greeting in %s: %v", lang, err)
	}
	fmt.Println(greeting)

	// Only print stats if we made an API call
	if greeter.stats.APICalls > 0 {
		fmt.Printf("\nTranslation Statistics:\n")
		fmt.Printf("From cache: %v\n", greeter.stats.APICalls == 0)
		fmt.Printf("Characters Translated: %d\n", greeter.stats.CharsSent)
		fmt.Printf("Estimated Cost: $%.5f\n", greeter.stats.CostEstimate)
	}
}
