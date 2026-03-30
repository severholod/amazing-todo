package users_transport_http

import (
	core_logger "github.com/severholod/amazing-todo/internal/core/logger"
	core_http_response "github.com/severholod/amazing-todo/internal/core/transport/http/response"
	core_http_utils "github.com/severholod/amazing-todo/internal/core/transport/http/utils"
	"net/http"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'id' path value")
		return
	}

	userDomain, err := h.usersService.GetUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
