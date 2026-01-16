package lua

import (
	lua "github.com/yuin/gopher-lua"
)

// NewSandboxedVM creates a new Lua VM with restricted libraries for security.
// Only safe libraries are loaded: base, table, string, and math.
// Dangerous functions like os, io, debug, loadfile, dofile are disabled.
func NewSandboxedVM() *lua.LState {
	L := lua.NewState(lua.Options{
		SkipOpenLibs: true, // Don't load any libs by default
	})

	// Manually open only safe libraries
	lua.OpenBase(L)    // Basic functions (print, type, etc.)
	lua.OpenTable(L)   // Table manipulation
	lua.OpenString(L)  // String manipulation
	lua.OpenMath(L)    // Math functions
	// NOT loading: os, io, debug, loadfile, dofile (security risk)

	// Remove potentially dangerous functions from base
	L.SetGlobal("dofile", lua.LNil)
	L.SetGlobal("loadfile", lua.LNil)
	L.SetGlobal("load", lua.LNil)
	L.SetGlobal("loadstring", lua.LNil)

	return L
}
