package users_transport_http

import (
	core_logger "github.com/severholod/amazing-todo/internal/core/logger"
	core_http_response "github.com/severholod/amazing-todo/internal/core/transport/http/response"
	core_http_utils "github.com/severholod/amazing-todo/internal/core/transport/http/utils"
	"net/http"
)

func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'id' path value")
		return
	}

	if err := h.usersService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")
		return
	}

	responseHandler.NoContentResponse()
}
