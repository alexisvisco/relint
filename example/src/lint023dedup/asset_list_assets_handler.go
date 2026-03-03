package handler

type ListAssetsInput struct{} // want `LINT-023: type "ListAssetsInput" must be declared in route file "asset_list_handler\.go"`

type ListAssetsOutput struct{} // want `LINT-023: type "ListAssetsOutput" must be declared in route file "asset_list_handler\.go"`
