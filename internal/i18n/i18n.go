package i18n

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localesFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	currentLang string
)

// Initialize sets up the i18n system with language detection
// Must be called before any translations are used
func Initialize(configLang string) error {
	// Create bundle with default language (Swedish)
	bundle = i18n.NewBundle(language.Swedish)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load Swedish translations from embedded files
	svData, err := localesFS.ReadFile("locales/active.sv.toml")
	if err != nil {
		return fmt.Errorf("failed to read embedded Swedish translations: %w", err)
	}
	if _, err := bundle.ParseMessageFileBytes(svData, "active.sv.toml"); err != nil {
		return fmt.Errorf("failed to parse Swedish translations: %w", err)
	}

	// Load English translations from embedded files
	enData, err := localesFS.ReadFile("locales/active.en.toml")
	if err != nil {
		return fmt.Errorf("failed to read embedded English translations: %w", err)
	}
	if _, err := bundle.ParseMessageFileBytes(enData, "active.en.toml"); err != nil {
		return fmt.Errorf("failed to parse English translations: %w", err)
	}

	// Detect language using cascading strategy
	lang := detectLanguage(configLang)
	currentLang = lang
	localizer = i18n.NewLocalizer(bundle, lang)

	return nil
}

// detectLanguage implements the cascading language detection strategy:
// 1. WORKLOG_LANG env var
// 2. Config file language setting (passed as parameter)
// 3. System LANG env var
// 4. Default to Swedish ("sv")
func detectLanguage(configLang string) string {
	// 1. Check WORKLOG_LANG
	if lang := os.Getenv("WORKLOG_LANG"); lang != "" {
		return normalizeLanguage(lang)
	}

	// 2. Check config file setting
	if configLang != "" {
		return normalizeLanguage(configLang)
	}

	// 3. Check system LANG
	if systemLang := os.Getenv("LANG"); systemLang != "" {
		return normalizeLanguage(systemLang)
	}

	// 4. Default to Swedish
	return "sv"
}

// normalizeLanguage converts various language formats to ISO 639-1 codes
// Examples: "sv_SE.UTF-8" -> "sv", "en_US" -> "en", "svenska" -> "sv"
func normalizeLanguage(lang string) string {
	lang = strings.ToLower(lang)

	// Extract language code from locale format (e.g., "sv_SE.UTF-8" -> "sv")
	if idx := strings.Index(lang, "_"); idx != -1 {
		lang = lang[:idx]
	}
	if idx := strings.Index(lang, "."); idx != -1 {
		lang = lang[:idx]
	}

	// Map full language names to codes
	switch lang {
	case "svenska", "swedish":
		return "sv"
	case "english", "engelska":
		return "en"
	}

	// Return as-is if already a valid code
	if lang == "sv" || lang == "en" {
		return lang
	}

	// Default to Swedish if unrecognized
	return "sv"
}

// T translates a message by its ID with optional template data
func T(messageID string, templateData ...interface{}) string {
	if localizer == nil {
		// Fallback if i18n not initialized
		return messageID
	}

	var data interface{}
	if len(templateData) > 0 {
		data = templateData[0]
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		// Return message ID if translation not found
		return messageID
	}
	return msg
}

// Tf translates a message with template data (convenience function)
// Example: Tf("cost_display", map[string]interface{}{"Cost": 5000.50})
func Tf(messageID string, templateData map[string]interface{}) string {
	return T(messageID, templateData)
}

// GetCurrentLanguage returns the currently active language code
func GetCurrentLanguage() string {
	if currentLang == "" {
		return "sv"
	}
	return currentLang
}
