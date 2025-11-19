package org

import "testing"

func TestShouldArchiveURL(t *testing.T) {
	testCases := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "regular website should be archived",
			url:      "https://example.com/article",
			expected: true,
		},
		{
			name:     "GitHub page should be archived",
			url:      "https://github.com/user/repo",
			expected: true,
		},
		{
			name:     "YouTube watch URL should not be archived",
			url:      "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
			expected: false,
		},
		{
			name:     "YouTube shorts URL should not be archived",
			url:      "https://www.youtube.com/shorts/dQw4w9WgXcQ",
			expected: false,
		},
		{
			name:     "YouTube short URL (youtu.be) should not be archived",
			url:      "https://youtu.be/dQw4w9WgXcQ",
			expected: false,
		},
		{
			name:     "Vimeo URL should not be archived",
			url:      "https://vimeo.com/123456789",
			expected: false,
		},
		{
			name:     "Twitch URL should not be archived",
			url:      "https://twitch.tv/streamer",
			expected: false,
		},
		{
			name:     "Dailymotion URL should not be archived",
			url:      "https://dailymotion.com/video/x123456",
			expected: false,
		},
		{
			name:     "blog post should be archived",
			url:      "https://blog.example.com/2024/01/post",
			expected: true,
		},
		{
			name:     "documentation should be archived",
			url:      "https://docs.example.com/guide",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := shouldArchiveURL(tc.url)
			if result != tc.expected {
				t.Errorf("shouldArchiveURL(%q) = %v, expected %v", tc.url, result, tc.expected)
			}
		})
	}
}

func TestSluggify(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{{
		input:    "",
		expected: "",
	}, {
		input:    "abcde",
		expected: "abcde",
	}, {
		input:    "abcde---",
		expected: "abcde",
	}, {
		input:    "a-b c--de",
		expected: "a-b-c-de",
	}, {
		input:    "a_bc__de",
		expected: "a-bc-de",
	}, {
		input:    "abcde$[)",
		expected: "abcde",
	}}
	for _, tc := range testCases {
		output := sluggify(tc.input)
		if output != tc.expected {
			t.Errorf("input \"%s\": expected %s, got %s", tc.input, tc.expected, output)
		}
	}
}
