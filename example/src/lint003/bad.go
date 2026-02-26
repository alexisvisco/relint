package lint003

import "log/slog"

func Bad() {
	slog.Info("msg", "error", "some error")      // want `LINT-003: slog key "error" should use dot notation, e.g. "error.message"`
	slog.Info("msg", "userId", "123")            // want `LINT-003: slog key "userId" should use dot notation, e.g. "user.id"`
	slog.Info("msg", "userID", "123")            // want `LINT-003: slog key "userID" should use dot notation, e.g. "user.id"`
	slog.Info("msg", "sessionId", "abc")         // want `LINT-003: slog key "sessionId" should use dot notation, e.g. "session.id"`
	slog.Info("msg", "sessionID", "abc")         // want `LINT-003: slog key "sessionID" should use dot notation, e.g. "session.id"`
	slog.Info("msg", "error.message", "some err") // ok
	slog.Info("msg", "user.id", "123")            // ok
	slog.Info("msg", "session.id", "abc")         // ok
}
