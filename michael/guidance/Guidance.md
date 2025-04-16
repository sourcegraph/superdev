# Go Project Guidelines

## Testing
- Write tests for all public functions
- Use table-driven tests where appropriate
- Mock external dependencies
- Aim for high test coverage, especially for core logic
- Use `t.Parallel()` for tests that can run concurrently
