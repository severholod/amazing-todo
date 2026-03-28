package users_transport_http

import (
	"fmt"
	"github.com/severholod/amazing-todo/internal/core/domain"
	core_errors "github.com/severholod/amazing-todo/internal/core/errors"
	core_logger "github.com/severholod/amazing-todo/internal/core/logger"
	core_http_request "github.com/severholod/amazing-todo/internal/core/transport/http/request"
	core_http_response "github.com/severholod/amazing-todo/internal/core/transport/http/response"
	core_http_types "github.com/severholod/amazing-todo/internal/core/transport/http/types"
	core_http_utils "github.com/severholod/amazing-todo/internal/core/transport/http/utils"
	"net/http"
	"strings"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get 'id' path value")
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func userPatchFromRequest(req PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    req.FullName.ToDomain(),
		PhoneNumber: req.PhoneNumber.ToDomain(),
	}
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("'full_name' cannot be nil.")
		}
		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf(
				"invalid `FullName` length: %d: %w`",
				fullNameLen,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value == nil {
			return fmt.Errorf("'phone_number' cannot be nil.")
		}

		phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf(
				"invalid `PhoneNumber` length: %d: %w",
				phoneNumberLen,
				core_errors.ErrInvalidArgument,
			)
		}
		if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
			return fmt.Errorf("invalid 'PhoneNumber' prefix: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
