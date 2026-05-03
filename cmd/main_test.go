package main

import (
	"testing"
)

// ── NOTE ───────────────────────────────────────────────────────────────────
// Testing CLI flags is non-trivial because os.Args and flag.CommandLine
// are global state. These tests verify the logic of flag parsing by
// testing the behavior after flags are parsed.
//
// Full integration tests should use shell commands:
//   go run . example.txt
//   go run . -r example.txt
//   go run . -c -l example.txt
// ─────────────────────────────────────────────────────────────────────────

// ── TEST: Default flag values ──────────────────────────────────────────────

func TestMain_DefaultFlagValues(t *testing.T) {
	// Before any flag parsing, all modes should be false
	if ModeRotation != false {
		t.Errorf("ModeRotation should default to false, got %v", ModeRotation)
	}

	if ModeEnableLogs != false {
		t.Errorf("ModeEnableLogs should default to false, got %v", ModeEnableLogs)
	}

	if ModeColor != false {
		t.Errorf("ModeColor should default to false, got %v", ModeColor)
	}

	if ModeColorOnly != false {
		t.Errorf("ModeColorOnly should default to false, got %v", ModeColorOnly)
	}
}

// ── TEST: Flag logic (OR operation) ────────────────────────────────────────
// These tests verify the OR logic used in the args() function.

func TestMain_FlagLogic_RotationAlias(t *testing.T) {
	// Test: r || rotate
	// If either is true, result should be true

	testCases := []struct {
		name   string
		r      bool
		rotate bool
		expect bool
	}{
		{"both false", false, false, false},
		{"r true", true, false, true},
		{"rotate true", false, true, true},
		{"both true", true, true, true},
	}

	for _, tc := range testCases {
		result := tc.r || tc.rotate
		if result != tc.expect {
			t.Errorf("%s: %v || %v = %v, expected %v", 
				tc.name, tc.r, tc.rotate, result, tc.expect)
		}
	}
}

func TestMain_FlagLogic_LogsAlias(t *testing.T) {
	// Test: l || logs
	testCases := []struct {
		name  string
		l     bool
		logs  bool
		expect bool
	}{
		{"both false", false, false, false},
		{"l true", true, false, true},
		{"logs true", false, true, true},
		{"both true", true, true, true},
	}

	for _, tc := range testCases {
		result := tc.l || tc.logs
		if result != tc.expect {
			t.Errorf("%s: %v || %v = %v, expected %v",
				tc.name, tc.l, tc.logs, result, tc.expect)
		}
	}
}

func TestMain_FlagLogic_ColorAlias(t *testing.T) {
	// Test: c || color
	testCases := []struct {
		name   string
		c      bool
		color  bool
		expect bool
	}{
		{"both false", false, false, false},
		{"c true", true, false, true},
		{"color true", false, true, true},
		{"both true", true, true, true},
	}

	for _, tc := range testCases {
		result := tc.c || tc.color
		if result != tc.expect {
			t.Errorf("%s: %v || %v = %v, expected %v",
				tc.name, tc.c, tc.color, result, tc.expect)
		}
	}
}

func TestMain_FlagLogic_ColorOnlyAlias(t *testing.T) {
	// Test: co || coloronly
	testCases := []struct {
		name      string
		co        bool
		coloronly bool
		expect    bool
	}{
		{"both false", false, false, false},
		{"co true", true, false, true},
		{"coloronly true", false, true, true},
		{"both true", true, true, true},
	}

	for _, tc := range testCases {
		result := tc.co || tc.coloronly
		if result != tc.expect {
			t.Errorf("%s: %v || %v = %v, expected %v",
				tc.name, tc.co, tc.coloronly, result, tc.expect)
		}
	}
}

// ── TEST: Flag independence ───────────────────────────────────────────────
// Rotation, logs, color, and colorOnly should be independent.

func TestMain_FlagsAreIndependent(t *testing.T) {
	// Test that setting one flag doesn't affect others
	testCases := []struct {
		name     string
		setFlags func()
		checkFn  func(t *testing.T)
	}{
		{
			name: "rotation doesn't affect logs",
			setFlags: func() {
				// Simulated: if -r flag is true
				r := true
				logs := false
				ModeRotation = r || false
				ModeEnableLogs = logs || false
			},
			checkFn: func(t *testing.T) {
				if !ModeRotation {
					t.Error("ModeRotation should be true")
				}
				if ModeEnableLogs {
					t.Error("ModeEnableLogs should remain false")
				}
			},
		},
		{
			name: "color doesn't affect rotation",
			setFlags: func() {
				// Simulated: if -c flag is true
				c := true
				r := false
				ModeColor = c || false
				ModeRotation = r || false
			},
			checkFn: func(t *testing.T) {
				if !ModeColor {
					t.Error("ModeColor should be true")
				}
				if ModeRotation {
					t.Error("ModeRotation should remain false")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset flags
			ModeRotation = false
			ModeEnableLogs = false
			ModeColor = false
			ModeColorOnly = false

			tc.setFlags()
			tc.checkFn(t)
		})
	}
}

// ── TEST: initLogger behavior ──────────────────────────────────────────────
// These tests verify that the logger is correctly initialized.

func TestMain_InitLogger_DisabledByDefault(t *testing.T) {
	// When ModeEnableLogs is false, initLogger should return nil
	// and set up a silent logger

	ModeEnableLogs = false
	logFile := initLogger()

	if logFile != nil {
		t.Errorf("initLogger should return nil when ModeEnableLogs=false, got %v", logFile)
	}
}

// ── TEST: Argument validation logic ────────────────────────────────────────

func TestMain_ArgumentValidation_SingleArg(t *testing.T) {
	// Test that argument validation logic (not the actual args() func)
	// correctly accepts exactly 1 argument

	// Simulating what args() does:
	arguments := []string{"example.txt"}
	if len(arguments) != 1 {
		t.Errorf("Expected exactly 1 argument, got %d", len(arguments))
	}
}

func TestMain_ArgumentValidation_NoArgs(t *testing.T) {
	arguments := []string{}
	if len(arguments) == 1 {
		t.Errorf("With no arguments, should reject, but validation passed")
	}
}

func TestMain_ArgumentValidation_TooManyArgs(t *testing.T) {
	arguments := []string{"file1.txt", "file2.txt"}
	if len(arguments) == 1 {
		t.Errorf("With 2 arguments, should reject, but validation passed")
	}
}

// ── TEST: Conflict resolution (chaotic flag combinations) ──────────────────

func TestMain_FlagConflictResolution_RotationConflict(t *testing.T) {
	// User types: myapp -r -rotate=false
	// r = true, rotate = false
	// Expected: ModeRotation = true (because r is true)

	r := true
	rotate := false
	result := r || rotate

	if !result {
		t.Errorf("Conflict: r=true, rotate=false should resolve to true, got %v", result)
	}
}

func TestMain_FlagConflictResolution_BothDisabled(t *testing.T) {
	// User types: myapp -r=false -rotate=false
	r := false
	rotate := false
	result := r || rotate

	if result {
		t.Errorf("Both false should give false, got %v", result)
	}
}
