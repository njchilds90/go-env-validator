# Contributing

Thank you for your interest in contributing to go-env-validator.

## Getting Started

1. Fork the repository and clone your fork.
2. Create a branch for your change: `git checkout -b your-feature-name`
3. Make your changes and add tests.
4. Run `go test -race ./...` and `go vet ./...` and ensure both pass.
5. Open a pull request against the `main` branch.

## Guidelines

- Every exported function must have a GoDoc comment and at least one example.
- New field kinds should include at least three test cases: valid, invalid, and default.
- The library must remain zero external dependencies. Do not add anything to `go.mod` beyond the Go standard library.
- Keep the public interface surface minimal. Prefer adding methods to existing types over introducing new top-level functions.
- All behavior must be deterministic. No randomness, no global mutable state.

## Reporting Issues

Open a GitHub issue with a minimal reproduction case. Include the Go version, operating system, and the exact error message.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
```

---

## 4. Release and Verification Instructions
```
RELEASE STEPS — GitHub Web UI Only

1. CREATE THE TAG AND RELEASE

   a. From your repository homepage, click "Releases" in the right sidebar.
   b. Click "Create a new release".
   c. In the "Choose a tag" dropdown, type:
         v1.0.0
      Then click "Create new tag: v1.0.0 on publish".
   d. Set the Release title to:
         v1.0.0 — Initial Release
   e. In the description box, paste:

         ## What's New

         Initial release of go-env-validator — declarative environment variable
         validation, type coercion, defaults, and machine-readable schema for Go.

         ### Features
         - Declare all config fields in one place with types, defaults, and descriptions
         - Validates and parses six types: string, integer, float, boolean, url, duration
         - Returns all validation errors at once — no fail-on-first behavior
         - Schema() outputs a JSON-safe machine-readable field description
         - ValidateMap for pure, testable validation without os.Getenv
         - Zero external dependencies, full context.Context support

   f. Make sure "Set as the latest release" is checked.
   g. Click "Publish release".

2. VERIFY ON pkg.go.dev

   Wait approximately 10 minutes, then visit:
   https://pkg.go.dev/github.com/njchilds90/go-env-validator

   If the page has not appeared yet, you can trigger indexing manually by visiting:
   https://pkg.go.dev/github.com/njchilds90/go-env-validator@v1.0.0

   pkg.go.dev discovers new modules automatically after a release tag is pushed.

3. SEMANTIC VERSIONING GUIDANCE

   v1.0.0  — current release (stable public API)
   v1.1.0  — add new Field options, new Kind values, new Result methods (backward compatible)
   v1.2.0  — add optional struct-tag-based declaration via reflection (backward compatible)
   v2.0.0  — only if the ValidateMap signature or Result accessor API must change
