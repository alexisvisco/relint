package lint030core

import (
	"context"

	assettypes "github.com/alexisvisco/relint/example/src/lint030assetserver/types" // want `LINT-030: package under root "lint030core" must not import sibling root "lint030assetserver" via "github\.com/alexisvisco/relint/example/src/lint030assetserver/types"`
	coretypes "github.com/alexisvisco/relint/example/src/lint030core/types"
	smarthubtypes "github.com/alexisvisco/relint/example/src/lint030smarthubserver/types" // want `LINT-030: package under root "lint030core" must not import sibling root "lint030smarthubserver" via "github\.com/alexisvisco/relint/example/src/lint030smarthubserver/types"`
)

var _ = context.Background
var _ = assettypes.AssetType{}
var _ = coretypes.CoreType{}
var _ = smarthubtypes.SmartHubType{}
