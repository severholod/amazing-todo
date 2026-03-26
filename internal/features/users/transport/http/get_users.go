package users_transport_http

import (
	"fmt"
	core_logger "github.com/severholod/amazing-todo/internal/core/logger"
	core_http_response "github.com/severholod/amazing-todo/internal/core/transport/http/response"
	core_http_utils "github.com/severholod/amazing-todo/internal/core/transport/http/utils"
	"net/http"
)

type GetUsersResponse []UserDTOResponse

const (
	QUERY_PARAM_LIMIT  = "limit"
	QUERY_PARAM_OFFSET = "offset"
)

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
	limit, offset, err := getQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query parameters",
		)
	}
	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
		return
	}

	response := GetUsersResponse(usersDTOFromDomain(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getQueryParams(r *http.Request) (*int, *int, error) {
	queryParamLimit, err := core_http_utils.GetIntQueryParam(r, QUERY_PARAM_LIMIT)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}
	queryParamOffset, err := core_http_utils.GetIntQueryParam(r, QUERY_PARAM_OFFSET)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return queryParamLimit, queryParamOffset, nil
}
