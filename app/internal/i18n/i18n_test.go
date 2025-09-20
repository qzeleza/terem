package i18n

import "testing"

func TestSetLanguageAndTranslate(t *testing.T) {
	t.Helper()

	if err := SetLanguage("en"); err != nil {
		t.Fatalf("set language en: %v", err)
	}
	t.Cleanup(func() { _ = SetLanguage("ru") })

	if got := T("menu.main.option.apps"); got != "Applications" {
		t.Fatalf("expected English translation, got %q", got)
	}
}

func TestFallbackToDefaultLanguage(t *testing.T) {
	t.Helper()

	ensureLoaded()

	mu.Lock()
	original := dictionaries["en"]["network.warn.invalid"]
	delete(dictionaries["en"], "network.warn.invalid")
	mu.Unlock()
	t.Cleanup(func() {
		mu.Lock()
		if original != "" {
			dictionaries["en"]["network.warn.invalid"] = original
		}
		mu.Unlock()
		_ = SetLanguage("ru")
	})

	if err := SetLanguage("en"); err != nil {
		t.Fatalf("set language en: %v", err)
	}

	mu.RLock()
	expected := dictionaries["ru"]["network.warn.invalid"]
	mu.RUnlock()

	if got := T("network.warn.invalid"); got != expected {
		t.Fatalf("expected fallback to %q, got %q", expected, got)
	}

	if missing := T("i18n.test.missing.key"); missing != "i18n.test.missing.key" {
		t.Fatalf("expected key for unknown translation, got %q", missing)
	}
}

func TestInvalidLanguage(t *testing.T) {
	if err := SetLanguage("zz"); err == nil {
		t.Fatal("expected error for unknown language")
	}
}
