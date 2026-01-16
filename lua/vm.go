package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// SandboxConfig configures which libraries and functions are available in the Lua VM.
// By default, only safe libraries are enabled (base, table, string, math).
// Dangerous functions (dofile, loadfile, load, loadstring) are disabled by default.
type SandboxConfig struct {
	// Libraries to enable/disable
	EnableBase   bool // Basic functions (print, type, etc.) - default: true
	EnableTable  bool // Table manipulation - default: true
	EnableString bool // String manipulation - default: true
	EnableMath   bool // Math functions - default: true
	EnableOS     bool // OS functions (os.execute, os.exit, etc.) - default: false
	EnableIO     bool // IO functions (file operations) - default: false
	EnableDebug  bool // Debug functions - default: false

	// Functions to disable (even if their library is enabled)
	// These are dangerous functions that should typically be disabled
	DisableDofile     bool // Disable dofile() - default: true
	DisableLoadfile   bool // Disable loadfile() - default: true
	DisableLoad       bool // Disable load() - default: true
	DisableLoadstring bool // Disable loadstring() - default: true
}

// DefaultSandboxConfig returns a default sandbox configuration with strict security.
// Only safe libraries are enabled, and dangerous functions are disabled.
func DefaultSandboxConfig() SandboxConfig {
	return SandboxConfig{
		EnableBase:        true,
		EnableTable:       true,
		EnableString:      true,
		EnableMath:        true,
		EnableOS:          false,
		EnableIO:          false,
		EnableDebug:       false,
		DisableDofile:     true,
		DisableLoadfile:   true,
		DisableLoad:       true,
		DisableLoadstring: true,
	}
}

// NewSandboxedVM creates a new Lua VM with restricted libraries for security.
// Only safe libraries are loaded: base, table, string, and math.
// Dangerous functions like os, io, debug, loadfile, dofile are disabled.
// This function uses DefaultSandboxConfig() for backward compatibility.
func NewSandboxedVM() *lua.LState {
	return NewVM(DefaultSandboxConfig())
}

// NewVM creates a new Lua VM with the specified sandbox configuration.
// If config is nil, DefaultSandboxConfig() is used.
func NewVM(config SandboxConfig) *lua.LState {
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true, // Don't load any libs by default
	})

	// Open libraries based on configuration
	if config.EnableBase {
		lua.OpenBase(L)
	}
	if config.EnableTable {
		lua.OpenTable(L)
	}
	if config.EnableString {
		lua.OpenString(L)
	}
	if config.EnableMath {
		lua.OpenMath(L)
	}
	if config.EnableOS {
		lua.OpenOs(L)
	}
	if config.EnableIO {
		lua.OpenIo(L)
	}
	if config.EnableDebug {
		lua.OpenDebug(L)
	}

	// Remove potentially dangerous functions based on configuration
	if config.DisableDofile {
		L.SetGlobal("dofile", lua.LNil)
	}
	if config.DisableLoadfile {
		L.SetGlobal("loadfile", lua.LNil)
	}
	if config.DisableLoad {
		L.SetGlobal("load", lua.LNil)
	}
	if config.DisableLoadstring {
		L.SetGlobal("loadstring", lua.LNil)
	}

	return L
}
