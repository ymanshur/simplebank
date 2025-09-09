# Development Best Practices

This document outlines the coding standards, testing guidelines, and development practices for the Simple Bank Service project.

## Coding Style & Naming Conventions

### Go Code Standards

- **Indentation**: Tabs (Go standard)
- **File naming**: snake_case for SQL files, camelCase for Go files
- **Function/variable naming**: camelCase for public, camelCase for private
- **Package naming**: Short, lowercase, single words when possible

### Database Conventions

- **Tables and columns**: snake_case naming
- **Foreign keys**: Follow pattern `table_name_id`
- **Indexes**: Descriptive names indicating purpose

### API Design

- **REST endpoints**: Follow RESTful conventions
- **gRPC services**: Use Protocol Buffers with HTTP Gateway
- **Error handling**: Consistent error response format
- **Validation**: Comprehensive input validation using Go validator

## Testing Guidelines

### Testing Framework

- **Primary**: Go's built-in testing package
- **Assertions**: Testify library for enhanced assertions
- **Mocking**: GoMock for database layer mocking
- **Coverage**: Aim for comprehensive test coverage

### Test Organization

- **Test files**: `*_test.go` alongside source files
- **Test naming**: `TestFunctionName_Scenario_ExpectedResult`
- **Test structure**: Arrange, Act, Assert pattern
- **Test data**: Use table-driven tests where appropriate

### Test Types

- **Unit tests**: Test individual functions and methods
- **Integration tests**: Test with real database connections
- **API tests**: Test HTTP and gRPC endpoints
- **Database tests**: Test SQL queries and transactions

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v ./api -run TestCreateAccount
```

## Commit & Pull Request Guidelines

### Commit Message Format

Follow conventional commit format: `type(scope): description`

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**

```markdown
feat(api): add user authentication endpoint
fix(db): resolve migration rollback issue
docs(readme): update installation instructions
test(transfer): add integration tests for money transfer
```

### Branch Naming

- **Feature branches**: `feature/description`
- **Bug fixes**: `fix/description`
- **Documentation**: `docs/description`
- **Automated updates**: `dependabot/*`

### Pull Request Process

1. **Create feature branch** from main/master
2. **Implement changes** following coding standards
3. **Write tests** for new functionality
4. **Update documentation** if needed
5. **Create PR** with descriptive title and description
6. **CI validation** - All tests must pass
7. **Code review** - At least one approval required
8. **Merge** after approval and CI success

## Build & Development Commands

### Essential Commands

```bash
# Start infrastructure
make postgres && make redis && make createdb && make migrateup

# Development workflow
make server         # Run the server
make test           # Run tests
make sqlc           # Generate SQL code
make proto          # Generate protobuf code
make mock           # Generate mocks

# Database operations
make migrateup      # Apply migrations
make migratedown    # Rollback migrations
make migratecreate name=<name>  # Create new migration

# Docker operations
make containers     # Run everything in Docker
make build          # Build Docker image
```

### Code Generation

- **SQLC**: Generates type-safe Go code from SQL
- **Protocol Buffers**: Generates gRPC and HTTP Gateway code
- **GoMock**: Generates mock interfaces for testing
- **Database docs**: Generates documentation from DBML

## Security Best Practices

### Authentication & Authorization

- **Token management**: Support both JWT and PASETO tokens
- **Token expiration**: Short-lived access tokens (15m), longer refresh tokens (24h)
- **Password security**: Bcrypt hashing with salt
- **Session management**: Store refresh tokens securely

### Data Protection

- **Input validation**: Validate all API inputs
- **SQL injection**: Use parameterized queries (SQLC handles this)
- **HTTPS**: Always use HTTPS in production
- **Environment variables**: Store secrets in environment variables

### Financial Data Integrity

- **ACID transactions**: All money operations in database transactions
- **Balance validation**: Prevent negative balances
- **Audit trail**: Complete transaction history
- **Currency validation**: Enforce valid currency codes

## Performance Considerations

### Database Optimization

- **Indexing**: Optimize indexes on frequently queried columns
- **Connection pooling**: Use database connection pools
- **Query optimization**: Monitor and optimize slow queries
- **Transaction scope**: Keep transactions as short as possible

### Application Performance

- **Background processing**: Use Redis queues for heavy tasks
- **Caching**: Implement caching where appropriate
- **Logging**: Structured logging with appropriate levels
- **Monitoring**: Health checks and metrics collection

## Error Handling

### Error Response Format

```json
{
    "error": "descriptive error message",
    "code": "ERROR_CODE",
    "details": {}
}
```

### Logging Standards

- **Structured logging**: Use Zerolog for JSON output
- **Log levels**: DEBUG, INFO, WARN, ERROR
- **Request tracing**: Include request IDs for tracing
- **Sensitive data**: Never log passwords or tokens

## Documentation Standards

### Code Documentation

- **Package comments**: Document package purpose
- **Function comments**: Document public functions
- **Complex logic**: Comment non-obvious code
- **API documentation**: Generate from Protocol Buffers

### Project Documentation

- **README**: Keep updated with setup instructions
- **BEST_PRACTICE**: This document is for development standards
- **API docs**: Auto-generated Swagger documentation
- **Database schema**: DBML documentation

## Continuous Integration

### GitHub Actions

- **Test pipeline**: Run full test suite on every PR
- **Database testing**: Use PostgreSQL service container
- **Code quality**: Lint and format checks
- **Security scanning**: Dependency vulnerability checks

### Quality Gates

- **All tests pass**: No failing tests allowed
- **Code coverage**: Maintain reasonable coverage levels
- **No security vulnerabilities**: Address security issues promptly
- **Code review**: Require peer review for all changes

---

*For questions about these practices or suggestions for improvements, please create an issue or discuss in pull requests.*
