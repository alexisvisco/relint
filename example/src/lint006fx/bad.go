package lint006fx

type fakeFX struct{}

func (fakeFX) Provide(...any) {}

var fx fakeFX

func init() {
	fx.Provide(LoadConfig)
}

func LoadConfig() (string, int, error) { return "", 0, nil } // ok - used in fx.Provide

func ParseConfig() (string, int, error) { return "", 0, nil } // want `LINT-006: function "ParseConfig" has 3 return values, consider using a ParseConfigResult struct`
