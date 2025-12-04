package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdemeester/go-org-readwise/internal/readwise"
)

func TestGetUpdateAfterFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		wantTime    *time.Time
		wantErr     bool
	}{
		{
			name:        "valid timestamp",
			fileContent: "2024-05-11T10:04:01",
			wantTime:    parseTime(t, "2024-05-11T10:04:01"),
			wantErr:     false,
		},
		{
			name:        "empty file",
			fileContent: "",
			wantTime:    nil,
			wantErr:     true,
		},
		{
			name:        "invalid format",
			fileContent: "2024-05-11 10:04:01",
			wantTime:    nil,
			wantErr:     true,
		},
		{
			name:        "malformed timestamp",
			fileContent: "not-a-timestamp",
			wantTime:    nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			stateFile := filepath.Join(tmpDir, ".readwise-sync.state")

			// Write test content to file
			if err := os.WriteFile(stateFile, []byte(tt.fileContent), 0o644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			// Test the function
			got, err := getUpdateAfterFromFile(stateFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUpdateAfterFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.wantTime == nil && got != nil {
					t.Errorf("getUpdateAfterFromFile() = %v, want nil", got)
				} else if tt.wantTime != nil && got == nil {
					t.Errorf("getUpdateAfterFromFile() = nil, want %v", tt.wantTime)
				} else if tt.wantTime != nil && got != nil && !tt.wantTime.Equal(*got) {
					t.Errorf("getUpdateAfterFromFile() = %v, want %v", got, tt.wantTime)
				}
			}
		})
	}
}

func TestGetUpdateAfterFromFile_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, ".readwise-sync.state")

	// Don't create the file - test non-existent file behavior
	got, err := getUpdateAfterFromFile(stateFile)
	if err != nil {
		t.Errorf("getUpdateAfterFromFile() should not error on non-existent file, got: %v", err)
	}
	if got != nil {
		t.Errorf("getUpdateAfterFromFile() = %v, want nil for non-existent file", got)
	}
}

func TestWriteUpdateAfterToFile(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "current time",
			time: time.Now(),
		},
		{
			name: "specific timestamp",
			time: time.Date(2024, 5, 11, 10, 4, 1, 0, time.UTC),
		},
		{
			name: "zero time",
			time: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			stateFile := filepath.Join(tmpDir, ".readwise-sync.state")

			// Write the timestamp
			if err := writeUpdateAfterToFile(stateFile, tt.time); err != nil {
				t.Fatalf("writeUpdateAfterToFile() error = %v", err)
			}

			// Verify file exists
			if _, err := os.Stat(stateFile); os.IsNotExist(err) {
				t.Fatal("State file was not created")
			}

			// Read and verify the content
			content, err := os.ReadFile(stateFile)
			if err != nil {
				t.Fatalf("Failed to read state file: %v", err)
			}

			expectedContent := tt.time.Format(readwise.FormatUpdatedAfter)
			if string(content) != expectedContent {
				t.Errorf("State file content = %q, want %q", string(content), expectedContent)
			}
		})
	}
}

func TestStateFileRoundTrip(t *testing.T) {
	// This test verifies that writing a timestamp and then reading it back
	// produces the same value (within the precision of the format)
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, ".readwise-sync.state")

	originalTime := time.Now()

	// Write the timestamp
	if err := writeUpdateAfterToFile(stateFile, originalTime); err != nil {
		t.Fatalf("writeUpdateAfterToFile() error = %v", err)
	}

	// Read it back
	readTime, err := getUpdateAfterFromFile(stateFile)
	if err != nil {
		t.Fatalf("getUpdateAfterFromFile() error = %v", err)
	}

	if readTime == nil {
		t.Fatal("getUpdateAfterFromFile() returned nil")
	}

	// Format both times the same way to compare (the format loses sub-second precision)
	expectedStr := originalTime.Format(readwise.FormatUpdatedAfter)
	gotStr := readTime.Format(readwise.FormatUpdatedAfter)

	if gotStr != expectedStr {
		t.Errorf("Round trip failed: got %q, want %q", gotStr, expectedStr)
	}
}

func TestStateFilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	stateFile := filepath.Join(tmpDir, ".readwise-sync.state")

	if err := writeUpdateAfterToFile(stateFile, time.Now()); err != nil {
		t.Fatalf("writeUpdateAfterToFile() error = %v", err)
	}

	// Check file permissions
	info, err := os.Stat(stateFile)
	if err != nil {
		t.Fatalf("Failed to stat state file: %v", err)
	}

	// Verify permissions are 0644 (owner: rw, group: r, others: r)
	if info.Mode().Perm() != 0o644 {
		t.Errorf("State file permissions = %o, want 0644", info.Mode().Perm())
	}
}

// Helper function to parse time in tests
func parseTime(t *testing.T, s string) *time.Time {
	t.Helper()
	parsed, err := time.Parse(readwise.FormatUpdatedAfter, s)
	if err != nil {
		t.Fatalf("Failed to parse test time %q: %v", s, err)
	}
	return &parsed
}
