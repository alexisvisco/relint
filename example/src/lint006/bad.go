package lint006

func Good() (int, error) { return 0, nil }

func Bad() (int, string, error) { return 0, "", nil } // want `LINT-006: function "Bad" has 3 return values, consider using a BadResult struct`

func AlsoBad() (int, string, bool, error) { return 0, "", false, nil } // want `LINT-006: function "AlsoBad" has 4 return values, consider using a AlsoBadResult struct`
