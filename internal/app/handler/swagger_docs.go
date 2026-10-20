package handler

// These small types exist only so swag can generate clear request schemas
// without depending on internal serializer type resolution.
type SwaggerUserAuthRequest struct {
	Name     string `json:"name" example:"lr4_creator"`
	Password string `json:"password" example:"123456"`
}

type SwaggerUserRegisterRequest struct {
	Name     string `json:"name" example:"lr4_creator"`
	Password string `json:"password" example:"123456"`
}

type SwaggerUserResponse struct {
	ID   uint   `json:"id" example:"3"`
	Name string `json:"name" example:"lr4_creator"`
}

type SwaggerUserAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int64  `json:"expires_in" example:"3599"`
	UserID      uint   `json:"user_id" example:"3"`
	Name        string `json:"name" example:"lr4_creator"`
	Role        string `json:"role" example:"creator"`
}

// @Summary Authenticate user
// @Description Checks creator credentials and returns a JWT access token.
// @Tags users
// @Accept json
// @Produce json
// @Param body body SwaggerUserAuthRequest true "Credentials"
// @Success 200 {object} SwaggerUserAuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/auth [post]
func swaggerAuth() {}

// @Summary Register user
// @Description Creates a creator account. Password is stored as a bcrypt hash.
// @Tags users
// @Accept json
// @Produce json
// @Param body body SwaggerUserRegisterRequest true "Registration data"
// @Success 201 {object} SwaggerUserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /users/register [post]
func swaggerRegister() {}

// @Summary Logout user
// @Description Revokes the current JWT by adding its jti to the Redis blacklist until token expiration.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/logout [post]
func swaggerLogout() {}

// @Summary Get electric car loads
// @Description Returns the service/load list. This endpoint is available without authentication.
// @Tags electric-car-loads
// @Produce json
// @Param category query string false "Category"
// @Param maxPowerW query int false "Maximum power draw in watts"
// @Param sort query string false "Sort order"
// @Success 200 {object} map[string]interface{}
// @Router /electric-car-loads [get]
func swaggerGetLoads() {}

// @Summary Get feed item
// @Description Returns a published feed item. It is available without authentication.
// @Tags electric-car-loads
// @Produce json
// @Param id query int false "Current item id"
// @Param next query bool false "Return next item"
// @Success 200 {object} map[string]interface{}
// @Router /electric-car-loads/feed [get]
func swaggerGetFeed() {}

// @Summary Get creator draft
// @Description Returns the authenticated creator's draft.
// @Tags electric-car-loads
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /electric-car-loads/draft [get]
func swaggerGetDraft() {}

// @Summary Create draft
// @Description Creates a new electric car load draft for the authenticated creator.
// @Tags electric-car-loads
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Name"
// @Param description formData string false "Description"
// @Param category formData string false "Category"
// @Param power_draw_w formData int false "Power draw in watts"
// @Param priority formData int false "Priority"
// @Param image formData file false "Image"
// @Param video formData file false "Video"
// @Success 201 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /electric-car-loads [post]
func swaggerCreateLoad() {}

// @Summary Publish/update load
// @Description Updates the creator's load and publishes it.
// @Tags electric-car-loads
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Load id"
// @Param name formData string false "Name"
// @Param description formData string false "Description"
// @Param category formData string false "Category"
// @Param power_draw_w formData int false "Power draw in watts"
// @Param priority formData int false "Priority"
// @Param image formData file false "Image"
// @Param video formData file false "Video"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /electric-car-loads/{id} [put]
func swaggerPublishLoad() {}

// @Summary Delete load
// @Description Deletes the authenticated creator's load.
// @Tags electric-car-loads
// @Produce json
// @Security BearerAuth
// @Param id path int true "Load id"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /electric-car-loads/{id} [delete]
func swaggerDeleteLoad() {}

// @Summary Like load
// @Description Adds/removes the authenticated creator's like.
// @Tags electric-car-loads
// @Produce json
// @Security BearerAuth
// @Param id path int true "Load id"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /electric-car-loads/{id}/like [post]
func swaggerLikeLoad() {}
