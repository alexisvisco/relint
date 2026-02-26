package lint001

import "log/slog"

func Bad() {
	slog.Info("msg", "UserID", 1)          // want `LINT-001: slog key "UserID" must be lowercase_snake_case`
	slog.Info("msg", "RequestID", "abc")   // want `LINT-001: slog key "RequestID" must be lowercase_snake_case`
	slog.Error("msg", "HttpStatus", 500)   // want `LINT-001: slog key "HttpStatus" must be lowercase_snake_case`
	slog.Info("msg", "user_id", 1)         // ok
	slog.Info("msg", "error.message", "x") // ok
	slog.Info("msg", "request_id", "abc")  // ok
}
