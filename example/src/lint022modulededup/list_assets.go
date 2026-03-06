package assethandler

import "net/http"

type AssetHandler struct{}

func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) { // want `LINT-022: route handler "ListAssets" on "AssetHandler" must be in file "list\.go"`
}
