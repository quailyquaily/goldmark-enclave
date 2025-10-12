package object

import (
	"strings"
	"testing"

	"github.com/quailyquaily/goldmark-enclave/core"
)

func TestGetPodbeanHtml_LightAndDark(t *testing.T) {
	tests := []struct {
		name         string
		theme        string
		expectedSkin string
	}{
		{name: "light skin", theme: "light", expectedSkin: "f6f6f6"},
		{name: "dark skin", theme: "dark", expectedSkin: "1b1b1b"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			enc := &core.Enclave{
				ObjectID:       "s9x5a-196f966-pb",
				Theme:          tc.theme,
				IframeDisabled: false,
			}
			html, err := GetPodbeanHtml(enc)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Basic structural checks
			if !strings.Contains(html, "<iframe") || !strings.Contains(html, "data-name=\"pb-iframe-player\"") {
				t.Fatalf("unexpected iframe markup: %s", html)
			}

			// Verify src contains id and the expected skin
			expectedSrc := "https://www.podbean.com/player-v2/?from=embed&i=s9x5a-196f966-pb"
			if !strings.Contains(html, expectedSrc) {
				t.Errorf("expected src to contain %q, got: %s", expectedSrc, html)
			}
			if !strings.Contains(html, "skin="+tc.expectedSkin) {
				t.Errorf("expected skin %q, got: %s", tc.expectedSkin, html)
			}
		})
	}
}
