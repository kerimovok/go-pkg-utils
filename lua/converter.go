package lua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// ConvertToLua converts a Go value to a Lua value.
// Supports: string, numbers (int, int32, int64, float32, float64), bool,
// map[string]interface{}, []interface{}, and nil.
// Unknown types are converted to string representation.
func ConvertToLua(L *lua.LState, v interface{}) lua.LValue {
	if v == nil {
		return lua.LNil
	}

	switch val := v.(type) {
	case string:
		return lua.LString(val)
	case float64:
		return lua.LNumber(val)
	case float32:
		return lua.LNumber(val)
	case int:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case int32:
		return lua.LNumber(val)
	case bool:
		return lua.LBool(val)
	case map[string]interface{}:
		return MapToLuaTable(L, val)
	case []interface{}:
		return SliceToLuaTable(L, val)
	default:
		// For unknown types, convert to string
		return lua.LString(fmt.Sprintf("%v", val))
	}
}

// MapToLuaTable converts a Go map[string]interface{} to a Lua table.
func MapToLuaTable(L *lua.LState, m map[string]interface{}) *lua.LTable {
	table := L.NewTable()

	for k, v := range m {
		table.RawSetString(k, ConvertToLua(L, v))
	}

	return table
}

// SliceToLuaTable converts a Go []interface{} to a Lua table.
// Lua arrays are 1-indexed, so the slice is converted accordingly.
func SliceToLuaTable(L *lua.LState, arr []interface{}) *lua.LTable {
	table := L.NewTable()

	for i, item := range arr {
		table.RawSetInt(i+1, ConvertToLua(L, item)) // Lua arrays are 1-indexed
	}

	return table
}
