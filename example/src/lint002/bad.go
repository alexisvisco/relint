package lint002

import "log/slog"

func Bad() {
	slog.Info("User created")  // want `LINT-002: slog message "User created" must start with a lowercase letter`
	slog.Error("Error found")  // want `LINT-002: slog message "Error found" must start with a lowercase letter`
	slog.Info("user created")  // ok
	slog.Error("error found")  // ok
	slog.Warn("Something bad") // want `LINT-002: slog message "Something bad" must start with a lowercase letter`
}
