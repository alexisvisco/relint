## Linter vs Formatter Classification

### Formatter (auto-fixable)

These are purely stylistic and can be deterministically rewritten without semantic understanding.

---

**FMT-001 — Type declaration merging**
Multiple consecutive single `type` declarations MUST be merged into a single `type (...)` block.

**FMT-002 — File declaration order**
Declarations within a file MUST follow the order: `type`, `const`, `var`, `func`. A formatter can reorder top-level declaration groups.

**FMT-003 — Function body spacing**
Functions MUST NOT start or end with empty lines. Logical blocks MUST NOT be separated by more than one blank line within a function body.

**FMT-004 — Interface method spacing**
Interface method signatures MUST be separated by exactly one blank line.

**FMT-005 — Type block spec spacing**
Type specs inside `type (...)` blocks MUST be separated by exactly one blank line.

---

### Linter (requires static analysis, not auto-fixable)

All rules skip generated Go files (files marked with `// Code generated ... DO NOT EDIT.`).

---

**LINT-001 — Log key casing**
String-literal slog key arguments inspected by this rule MUST be in `lowercase_snake_case`. Keys using dot notation (e.g. `error.message`) are permitted. Keys in `PascalCase`, `camelCase`, or containing uppercase letters are flagged.

**LINT-002 — Log message casing**
The first argument (message string) passed to `slog.*` calls MUST start with a lowercase letter.

**LINT-003 — Log key dot notation for grouped keys**
Log keys that semantically belong to a group (e.g. error fields, user fields) MUST use dot notation.

This rule is configuration-driven via `-dot-notation` as comma-separated `key=dotted_key` pairs (for example: `error=error.message,userId=user.id`). Only configured keys are flagged.

**LINT-004 — Context as first parameter**
Any function that accepts a `context.Context` MUST have it as the first parameter. Functions with `context.Context` in any other position MUST be flagged.

**LINT-005 — Excessive function parameters**
Functions with more than 4 parameters MUST be flagged. The message SHOULD suggest introducing a `{Name}Params` struct.

**LINT-006 — Excessive return values**
Functions with more than 2 return values MUST be flagged. The message SHOULD suggest introducing a `{Name}Result` struct.

Exception: functions referenced by `fx.Provide(...)` are excluded from this rule.

**LINT-007 — Enum value prefix**
For any named type backed by a primitive (string, int, etc.) with associated `const` values, each constant MUST be prefixed with the type name. Constants that do not carry the type name as a prefix MUST be flagged.

Configurable exceptions are supported via `package.Type` values. Default exception: `environment.Environment`.

**LINT-008 — Package name underscore**
Package names MUST NOT contain underscores. Any `package` declaration with an underscore in the name MUST be flagged.

Package-name suffixes can be excluded from this check via configuration. Default excluded suffix: `_test`.

**LINT-009 — Package name plural**
Package names that are pluralized MUST NOT be used.

The rule detects plural names generically (for example names ending with `s`), with configurable package-name exceptions.
Default configured exception: `types`.

**LINT-010 — Interface location**
Only interfaces suffixed with `Service` or `Store` MUST be declared in a `types` package (i.e. a file whose package is `types`). `Service`/`Store` interface declarations found outside of a `types` package MUST be flagged. Other interfaces are allowed outside `types`.

**LINT-011 — Service interface suffix**
Interfaces whose names do not end with `Service` or `Store` and are located in a `types` package MUST be evaluated. Specifically, interfaces semantically acting as services MUST be suffixed `Service`, and those acting as stores MUST be suffixed `Store`. This rule is best enforced via a naming pattern: any interface in `types` that wraps data access methods and is not suffixed `Store` MUST be flagged; any interface wrapping business logic methods not suffixed `Service` MUST be flagged. In practice, enforce: all interfaces in `types/` MUST end with either `Service` or `Store`.

**LINT-012 — Store function return types**
In packages whose name contains `store`, methods on receivers suffixed `Store` MUST NOT return types from packages whose import path contains `core/model` (including pointers/slices of those types).

**LINT-013 — Store struct interface assertion**
In packages whose name contains `store`, every exported struct suffixed `Store` MUST have a compile-time assertion in `store.go` whose value side matches `(*{Name}Store)(nil)` (for example: `var _ types.AnyStore = (*UserStore)(nil)`).

**LINT-014 — Service struct interface assertion**
In packages whose name contains `service`, every exported struct suffixed `Service` MUST have a compile-time assertion in `service.go` whose value side matches `(*{Name}Service)(nil)` (for example: `var _ types.AnyService = (*UserService)(nil)`).

**LINT-015 — One exported function per store/service file**
Files in packages whose name contains `store`, `service`, or `handler` (excluding `store.go`, `service.go`, and `fx_module.go`) are checked based on exported methods whose receiver name ends with `Store`, `Service`, or `Handler`.

If a file contains more than one such exported layer method, it is flagged. Exported non-method functions are ignored by this rule.

**LINT-016 — Middleware naming: Inject***
In `handler` packages, any function named `Inject{Name}` or `inject{Name}` (with non-empty `{Name}`) MUST be declared in `inject_{name}.go`. Violations are flagged.

**LINT-017 — Middleware naming: Require***
In `handler` packages, any function named `Require{Name}` or `require{Name}` (with non-empty `{Name}`) MUST be declared in `require_{name}.go`. Violations are flagged.

**LINT-018 — Middleware naming outside handler**
Outside `handler` packages, exported functions with middleware signature `func(http.Handler) http.Handler` MUST be named `Middleware`. Non-matching names are flagged.

**LINT-019 — fx_module.go presence**
Packages whose names end with `store`, `service`, or `handler` MUST contain an `fx_module.go` file. Absence is flagged.

If present, `fx_module.go` MUST declare an exported `FxModule` variable.

**LINT-020 — Error variable location (types package)**
In `types` packages only, error variables prefixed with `Err` MUST be declared in `errors.go`. `Err*` variables declared in other files within `types` MUST be flagged. Non-`types` packages are excluded from this rule.

**LINT-021 — RecordNotFound as typed error**
In packages whose name contains `store`, direct `return` expressions of these known not-found sentinels are flagged:
- `sql.ErrNoRows`
- `pgx.ErrNoRows`
- `gorm.ErrRecordNotFound`

**LINT-022 — Handler route file naming**
In `handler` packages, exported methods on receivers `*{Name}Handler` MUST be located in files named `{name}_{route}_handler.go` (where `{route}` is the method name in snake_case). Deviations are flagged.

Special case: when route name equals handler base name (for example `TenantHandler.Tenant`), the valid file is `{name}_handler.go` (for example `tenant_handler.go`).

**LINT-023 — Route Input/Output type location**
In `handler` packages, types suffixed `Input` or `Output` are treated as route types. If a matching handler method `{Route}` exists, the type MUST be declared either:
- in the corresponding route file (`{name}_{route}_handler.go`, or `{name}_handler.go` when route equals handler base name), or
- in the shared handler file (`{name}.go`).

Otherwise it is ignored by this rule.

**LINT-024 — Shared body type naming**
In `handler` packages, for files not ending in `_handler.go`, type names containing `Body` MUST match `{Name}BodyInput` or `{Name}BodyOutput`. Non-matching names are flagged.

**LINT-025 — Handler struct file location**
In `handler` packages, struct types suffixed `Handler` MUST be declared in `{name}.go` (for example `TenantHandler` in `tenant.go`).

**LINT-026 — Body-only helper struct naming**
In `handler` packages, struct types that are referenced only by body structs (`*BodyInput`/`*BodyOutput`) MUST:
- start with the parent body prefix (parent name without `Input`/`Output`), and
- end with the corresponding parent suffix (`Input` or `Output`).

**LINT-027 — No json tags in model structs**
In `model` packages, struct fields MUST NOT declare `json` tags. Fields with `json` tags are flagged.

**LINT-028 — Exported model fields require gorm tag**
In `model` packages, exported struct fields MUST declare a `gorm` tag attribute.

**LINT-029 — Relation field pointer shape**
In `model` packages, relation fields identified by `gorm` tag attributes `foreignKey`, `many2many`, or `polymorphicType` MUST be either:
- a pointer (`*Type`), or
- a slice of pointers (`[]*Type`).
