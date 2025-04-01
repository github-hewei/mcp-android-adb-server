package vision_test

import (
	"path"
	"strings"
	"testing"

	"mcp-android-adb-server/vision"
)

// TestImage2Base64 tests the Image2Base64 function
func TestImage2Base64(t *testing.T) {
	testImagePath := path.Join("D:", "temp", "test.png")
	base64Str, err := vision.Image2Base64(testImagePath)

	if err != nil {
		t.Errorf("Image2Base64 failed: %v", err)
	}

	if base64Str == "" {
		t.Error("Image2Base64 returned empty string")
	}

	_, err = vision.Image2Base64("nonexistent.png")
	if err == nil {
		t.Error("Should return error for non-existent file")
	}
}

// TestNewModel tests the NewModel function
func TestNewModel(t *testing.T) {
	testCases := []struct {
		name    string
		apiKey  string
		model   string
		baseURL string
		wantNil bool
	}{
		{
			name:    "Normal initialization",
			apiKey:  "test-api-key",
			model:   "gpt-4-vision-preview",
			baseURL: "https://api.openai.com/v1",
			wantNil: false,
		},
		{
			name:    "Empty API Key",
			apiKey:  "",
			model:   "gpt-4-vision-preview",
			baseURL: "https://api.openai.com/v1",
			wantNil: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := vision.NewModel(tc.apiKey, tc.model, tc.baseURL)
			if (model == nil) != tc.wantNil {
				t.Errorf("NewModel() = %v, want nil: %v", model, tc.wantNil)
			}
		})
	}
}

// TestScreenshotDescription tests the ScreenshotDescription function
func TestScreenshotDescription(t *testing.T) {
	testImagePath := path.Join("D:", "temp", "test.png")
	model := vision.NewModel(
		"xxx",
		"qwen/qwen2.5-vl-72b-instruct:free",
		"https://openrouter.ai/api/v1/",
	)

	description, err := model.ScreenshotDescription(testImagePath)
	if err != nil {
		t.Errorf("ScreenshotDescription failed: %v", err)
	}

	if strings.TrimSpace(description) == "" {
		t.Error("ScreenshotDescription returned empty description")
	}
}
