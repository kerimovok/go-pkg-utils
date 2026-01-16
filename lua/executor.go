package lua

import (
	"context"
	"fmt"
	"time"

	converter "github.com/kerimovok/go-lua-converter"
	lua "github.com/yuin/gopher-lua"
	"go.uber.org/zap"
)

// Script represents a script that can be executed.
type Script interface {
	GetID() string
	GetName() string
	GetVersion() string
	GetCode() string
}

// ExecutionResult represents the result of a script execution.
type ExecutionResult struct {
	ScriptID      string
	ScriptName    string
	ScriptVersion string
	Status        ExecutionStatus
	ErrorMessage  *string
	DurationMs    int64
	ExecutedAt    time.Time
}

// ExecutionStatus represents the status of an execution.
type ExecutionStatus string

const (
	ExecutionStatusSuccess ExecutionStatus = "success"
	ExecutionStatusFailure ExecutionStatus = "failure"
)

// HostFunctionRegistry allows registering custom host functions for Lua scripts.
type HostFunctionRegistry interface {
	RegisterFunctions(L *lua.LState, scriptID, scriptName, scriptVersion string)
}

// ExecutionRecorder allows recording execution results.
type ExecutionRecorder interface {
	RecordExecution(ctx context.Context, result ExecutionResult) error
}

// ExecutorConfig holds configuration for the executor.
type ExecutorConfig struct {
	Timeout       time.Duration
	Logger        *zap.Logger
	HostFunctions HostFunctionRegistry
	Recorder      ExecutionRecorder
	Sandbox       *SandboxConfig // Optional: if nil, DefaultSandboxConfig() is used
}

// Executor executes Lua scripts with timeout, error handling, and result recording.
type Executor struct {
	config ExecutorConfig
}

// NewExecutor creates a new Lua script executor.
func NewExecutor(config ExecutorConfig) *Executor {
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Second
	}

	return &Executor{
		config: config,
	}
}

// Execute executes a Lua script with the given payload.
// The payload is converted to a Lua table and passed to the script's "handle" function.
// Returns an ExecutionResult with the outcome.
func (e *Executor) Execute(ctx context.Context, script Script, payload map[string]interface{}) ExecutionResult {
	startTime := time.Now()
	var execErr error
	var errorMsg *string

	// Create a fresh sandboxed VM for this execution
	// Use provided sandbox config or default to strict sandboxing
	sandboxConfig := DefaultSandboxConfig()
	if e.config.Sandbox != nil {
		sandboxConfig = *e.config.Sandbox
	}
	L := NewVM(sandboxConfig)
	defer L.Close()

	// Register host functions if provided
	if e.config.HostFunctions != nil {
		e.config.HostFunctions.RegisterFunctions(L, script.GetID(), script.GetName(), script.GetVersion())
	}

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	// Set up timeout cancellation
	L.SetContext(execCtx)

	// Load and execute the script code
	if err := L.DoString(script.GetCode()); err != nil {
		execErr = err
		errStr := fmt.Sprintf("failed to load script: %v", err)
		errorMsg = &errStr
		if e.config.Logger != nil {
			e.config.Logger.Error("Failed to load script",
				zap.String("script_id", script.GetID()),
				zap.String("script_name", script.GetName()),
				zap.String("script_version", script.GetVersion()),
				zap.Error(err))
		}
	}

	if execErr == nil {
		// Get the handle function
		handleFn := L.GetGlobal("handle")
		if handleFn == lua.LNil {
			execErr = fmt.Errorf("script missing handle function")
			errStr := "script missing handle function"
			errorMsg = &errStr
			if e.config.Logger != nil {
				e.config.Logger.Error("Script missing handle function",
					zap.String("script_id", script.GetID()),
					zap.String("script_name", script.GetName()),
					zap.String("script_version", script.GetVersion()))
			}
		} else {
			// Convert payload to Lua table
			payloadTable := converter.MapToTable(L, payload)

			// Call the handle function
			if err := L.CallByParam(lua.P{
				Fn:      handleFn,
				NRet:    0,
				Protect: true,
			}, payloadTable); err != nil {
				execErr = err
				// Check if it was a timeout
				if execCtx.Err() == context.DeadlineExceeded {
					errStr := fmt.Sprintf("script execution timed out after %v", e.config.Timeout)
					errorMsg = &errStr
					if e.config.Logger != nil {
						e.config.Logger.Error("Script execution timed out",
							zap.String("script_id", script.GetID()),
							zap.String("script_name", script.GetName()),
							zap.String("script_version", script.GetVersion()),
							zap.Duration("timeout", e.config.Timeout))
					}
				} else {
					errStr := err.Error()
					errorMsg = &errStr
					if e.config.Logger != nil {
						e.config.Logger.Error("Script execution failed",
							zap.String("script_id", script.GetID()),
							zap.String("script_name", script.GetName()),
							zap.String("script_version", script.GetVersion()),
							zap.Error(err))
					}
				}
			}
		}
	}

	// Calculate duration
	duration := time.Since(startTime)
	durationMs := duration.Milliseconds()

	// Determine status
	status := ExecutionStatusSuccess
	if execErr != nil {
		status = ExecutionStatusFailure
	}

	// Log success
	if execErr == nil && e.config.Logger != nil {
		e.config.Logger.Info("Script executed successfully",
			zap.String("script_id", script.GetID()),
			zap.String("script_name", script.GetName()),
			zap.String("script_version", script.GetVersion()),
			zap.String("status", string(status)),
			zap.Int64("duration_ms", durationMs))
	}

	result := ExecutionResult{
		ScriptID:      script.GetID(),
		ScriptName:    script.GetName(),
		ScriptVersion: script.GetVersion(),
		Status:        status,
		ErrorMessage:  errorMsg,
		DurationMs:    durationMs,
		ExecutedAt:    startTime,
	}

	// Record execution result if recorder is provided
	if e.config.Recorder != nil {
		recordCtx, recordCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer recordCancel()

		if err := e.config.Recorder.RecordExecution(recordCtx, result); err != nil && e.config.Logger != nil {
			e.config.Logger.Error("Failed to record execution result",
				zap.String("script_id", script.GetID()),
				zap.String("script_name", script.GetName()),
				zap.String("script_version", script.GetVersion()),
				zap.Error(err))
		}
	}

	return result
}
