# TheDeeps Testing Guidelines

## Overview

TheDeeps employs comprehensive testing practices to ensure code quality, reliability, and maintainability. This document outlines the testing approach, tools, and best practices used throughout the project.

## Testing Philosophy

- **Test-Driven Development**: Write tests before implementing features when possible
- **Comprehensive Coverage**: Aim for high test coverage across all components
- **Edge Case Testing**: Explicitly test boundary conditions and error scenarios
- **Maintainable Tests**: Write clear, concise tests that are easy to understand and maintain
- **Continuous Integration**: Run tests automatically on code changes

## Testing Tools and Frameworks

### Go Testing

- **Standard Library**: Use the Go standard `testing` package as the foundation
- **Testify**: Leverage the `github.com/stretchr/testify` package for assertions and mocks
  - `assert`: For test assertions with clear failure messages
  - `require`: For assertions that should terminate the test on failure
  - `mock`: For creating mock objects
- **Ginkgo**: Use Ginkgo for BDD-style tests when appropriate
- **Gomega**: Pair with Ginkgo for expressive assertions

### JavaScript/TypeScript Testing

- **Jest**: Primary testing framework for client-side code
- **React Testing Library**: For testing React components
- **Mock Service Worker**: For mocking API requests in frontend tests

## Test Types

### Unit Tests

- Test individual functions, methods, and components in isolation
- Mock dependencies to focus on the unit being tested
- Keep tests small, focused, and fast

### Integration Tests

- Test interactions between components
- Verify that components work together correctly
- Focus on API boundaries and data flow

### End-to-End Tests

- Test complete user workflows
- Verify system behavior from the user's perspective
- Cover critical paths through the application

## Test Structure

### Go Tests

- **Table-Driven Tests**: Use table-driven tests for testing multiple scenarios
  ```go
  tests := []struct {
      name     string
      input    string
      expected string
      wantErr  bool
  }{
      {"Valid Input", "test", "result", false},
      {"Invalid Input", "", "", true},
  }
  
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          result, err := functionUnderTest(tt.input)
          if tt.wantErr {
              assert.Error(t, err)
          } else {
              assert.NoError(t, err)
              assert.Equal(t, tt.expected, result)
          }
      })
  }
  ```

- **Test Naming**: Use descriptive names that indicate what is being tested
  - Format: `Test<FunctionName>_<Scenario>`
  - Example: `TestCreateCharacter_ValidInput`

- **Test Organization**: Group related tests together in the same file
  - Place test files in the same package as the code being tested
  - Use `_test.go` suffix for test files

### JavaScript/TypeScript Tests

- **Test Naming**: Use descriptive names that indicate what is being tested
  - Format: `should <expected behavior> when <condition>`
  - Example: `should display error message when input is invalid`

- **Test Organization**: Group related tests using `describe` blocks
  ```javascript
  describe('Component', () => {
    describe('when initialized', () => {
      it('should render correctly', () => {
        // Test code
      });
    });
  });
  ```

## Test Coverage

- **Coverage Goals**: Aim for at least 80% code coverage
- **Coverage Reports**: Generate and review coverage reports regularly
- **Critical Path Coverage**: Ensure 100% coverage for critical business logic
- **Make Targets**: Use the provided make targets to run tests with coverage
  - `make server-test`: Run all server tests
  - `make server-test-coverage`: Run server tests with coverage report
  - `make server-test-ginkgo`: Run server tests using Ginkgo
  - `make client-test`: Run all client tests
  - `make client-test-coverage`: Run client tests with coverage report

## Mocking

### Go Mocking

- Use interfaces to make code testable
- Create mock implementations of interfaces for testing
- Use the testify mock package for complex mocking scenarios

### JavaScript/TypeScript Mocking

- Use Jest's mocking capabilities for functions and modules
- Use Mock Service Worker for API mocking
- Create test doubles (stubs, spies, mocks) as needed

## Test Data

- **Test Fixtures**: Create reusable test fixtures for common test scenarios
- **Random Data**: Use deterministic random data with fixed seeds for reproducibility
- **Edge Cases**: Include edge cases in test data (empty values, maximum values, etc.)

## Best Practices

1. **Independent Tests**: Each test should be independent and not rely on the state from other tests
2. **Fast Tests**: Tests should run quickly to encourage frequent testing
3. **Deterministic Tests**: Tests should produce the same results each time they run
4. **Clear Assertions**: Make assertions clear and specific
5. **Test Behavior, Not Implementation**: Focus on testing what the code does, not how it does it
6. **Avoid Test Duplication**: Don't duplicate test logic; use helpers and fixtures
7. **Test Error Handling**: Explicitly test error conditions and edge cases
8. **Clean Up**: Clean up resources after tests (close files, connections, etc.)

## Continuous Integration

- Tests run automatically on pull requests
- All tests must pass before merging
- Coverage reports are generated and reviewed
- Performance benchmarks are run to detect regressions

## Debugging Tests

- Use `t.Log` or `fmt.Printf` for debugging output
- Run specific tests with the `-run` flag
- Use the `-v` flag for verbose output
- Use breakpoints in your IDE for step-through debugging