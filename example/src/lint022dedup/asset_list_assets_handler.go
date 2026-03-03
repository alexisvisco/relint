package handler

type AssetHandler struct{}

func (h *AssetHandler) ListAssets() {} // want `LINT-022: route handler "ListAssets" on "AssetHandler" must be in file "asset_list_handler\.go"`
