---
globs: ["**/*_test.go"]
---

# Go Testing Rules

## Test Structure

- Use table-driven tests for multiple test cases
- Use `t.Run()` for subtests with descriptive names
- Keep test functions focused on one behavior
- Use `t.Helper()` in test helper functions

## Naming

- Test functions: `TestFunctionName_Scenario`
- Example: `TestBuilder_Build_Success`, `TestBuilder_Build_Failure`
- Subtest names should describe the scenario

## Setup and Teardown

- Use `t.TempDir()` for temporary directories (auto-cleanup)
- Use `t.Cleanup()` for custom cleanup
- Avoid global state; inject dependencies

## Assertions

- Use clear error messages with `t.Errorf()`
- Include expected vs actual values: `got %v, want %v`
- Fail fast with `t.Fatalf()` for setup errors
- Use `t.Skip()` for conditional tests

## Mocking

- Define interfaces for external dependencies
- Create mock implementations in test files
- Use dependency injection to swap implementations

## Coverage

- Aim for high coverage on critical paths
- Don't chase 100% coverage at the expense of test quality
- Focus on testing behavior, not implementation details

## Example Test Structure

```go
func TestBuilder_Build(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    Result
        wantErr bool
    }{
        {
            name:    "successful build",
            input:   "valid",
            want:    Result{Success: true},
            wantErr: false,
        },
        {
            name:    "build failure",
            input:   "invalid",
            want:    Result{Success: false},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Build(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Build() = %v, want %v", got, tt.want)
            }
        })
    }
}
```
