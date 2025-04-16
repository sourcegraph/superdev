# Go Project Guidelines

## Project Structure
- Keep a flat structure for small projects
- Use `cmd/`, `pkg/`, and `internal/` directories as needed
- Place main application code in `main.go` or `cmd/`
- Store reusable packages in `pkg/`
- Hide implementation details in `internal/`

## Cobra Best Practices
- Define a root command that represents your application
- Create subcommands for different functionality
- Use `PersistentFlags` for flags shared across commands
- Use `Flags` for command-specific flags
- Implement proper help text for all commands and flags
- Use `PreRun` and `PostRun` hooks when appropriate

## Error Handling
- Return errors rather than using panic
- Use descriptive error messages
- Consider using custom error types for specific error cases
- Wrap errors with context using `fmt.Errorf("doing something: %w", err)`

## Testing
- Write tests for all public functions
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for high test coverage, especially for core logic
- Use `t.Parallel()` for tests that can run concurrently

## Code Style
- Follow standard Go formatting with `go fmt`
- Use `golint` and `go vet` regularly
- Keep functions small and focused
- Use meaningful variable names
- Document all exported functions, types, and packages

## Dependencies
- Use Go modules for dependency management
- Minimize external dependencies
- Pin dependencies to specific versions
- Regularly update dependencies

## Documentation
- Write clear godoc comments for exported functions and types
- Include examples in documentation when helpful
- Document command usage in CLI help text
- Create a README.md with installation and usage instructions
