package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/idkwhyureadthis/project-practicum/restaurants/internal/service"
	"github.com/labstack/echo/v4"
)

// @Summary      Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "User credentials"
// @Success      200 {object} LoginRequest
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /login [post]
func (h *Handler) login(c echo.Context) error {
	var data LoginRequest

	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}
	tokens, err := h.s.LogIn(data.Login, data.Password)
	if errors.Is(err, service.ErrWrongData) {
		return Err2Json(err.Error(), c, http.StatusUnauthorized)
	} else if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.Value = tokens.Refresh
	cookie.Expires = time.Now().Add(time.Hour * 24)

	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens)
}

// @Summary      Verify token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        token body VerifyRequest true "Access token"
// @Success      200 {object} VerifyRequest
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /verify [post]
func (h *Handler) verify(c echo.Context) error {
	var data VerifyRequest
	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	role, err := h.s.Verify(data.Token, "access", c)
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"role": &role,
	})
}

// @Summary      Refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        refresh body GenerateRequest true "Refresh token"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Router       /refresh [post]
func (h *Handler) refresh(c echo.Context) error {
	var data GenerateRequest
	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	tokens, err := h.s.Generate(data.Token)
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusUnauthorized)
	}
	return c.JSON(http.StatusOK, tokens)
}

// @Summary      Add restaurant
// @Tags         Restaurants
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        restaurant body AddRestaurantRequest true "Restaurant data"
// @Success      201 {object} AddRestaurantRequest
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /restaurants [post]
func (h *Handler) addRestaurant(c echo.Context) error {
	var data AddRestaurantRequest
	role := c.Get("role")
	if role.(string) != "superadmin" {
		return Err2Json("only superadmin can add new restaurants", c, http.StatusUnauthorized)
	}

	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	restId, err := h.s.AddRestaurant(data.OpenTime, data.CloseTime, data.Name, data.Latitude, data.Longitude)

	if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"restaurant_id": restId.String(),
	})
}

// @Summary      Get restaurants
// @Tags         Restaurants
// @Produce      json
// @Success      200 {array} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /restaurants [get]
func (h *Handler) getRestaurants(c echo.Context) error {
	restaurants, err := h.s.GetRestaurants()
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, restaurants)
}

// @Summary      Create admin
// @Tags         Admins
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        admin body CreateAdminRequest true "Admin data"
// @Success      201 {object} CreateAdminRequest
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      409 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admins [post]
func (h *Handler) createAdmin(c echo.Context) error {
	var data CreateAdminRequest

	role := c.Get("role")
	if role.(string) != "superadmin" {
		return Err2Json("only superadmins are permitted to do that", c, http.StatusUnauthorized)
	}

	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	parsedUuid, err := uuid.Parse(data.RestaurantId)
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	adminID, err := h.s.CreateAdmin(parsedUuid, data.Login, data.Password)

	if err != nil && errors.Is(err, service.ErrLoginOccupied) {
		return Err2Json(err.Error(), c, http.StatusConflict)
	} else if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"created_id": adminID.String(),
		})
	}
}

// @Summary      Add item
// @Tags         Items
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        name formData string true "Item name"
// @Param        description formData string true "Description"
// @Param        sizes formData []string true "Sizes"
// @Param        prices formData []string true "Prices"
// @Param        images formData file true "Images"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /items [post]
func (h *Handler) addItem(c echo.Context) error {
	role := c.Get("role")
	if role.(string) != "superadmin" {
		return Err2Json("only superadmins are permitted to do that", c, http.StatusUnauthorized)
	}

	multiformData, err := c.MultipartForm()
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}
	if multiformData == nil {
		return Err2Json("no data for item provided", c, http.StatusBadRequest)
	}

	images := multiformData.File["images"]
	sizes := multiformData.Value["sizes"]
	prices := multiformData.Value["prices"]
	description := multiformData.Value["description"]
	name := multiformData.Value["name"]
	if len(description) != 1 || len(name) != 1 || len(images) == 0 || len(prices) == 0 || len(sizes) == 0 {
		return Err2Json("wrong item data provided", c, http.StatusBadRequest)
	}
	itemId, err := h.s.CreateItem(sizes, prices, name[0], description[0], images)

	if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"item_id": itemId.String(),
	})
}

// @Summary      Get items
// @Tags         Items
// @Produce      json
// @Param        restaurant_id query string true "Restaurant ID"
// @Success      200 {array} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /items [get]
func (h *Handler) getItems(c echo.Context) error {
	restaurantId := c.QueryParam("restaurant_id")

	parsedId, err := uuid.Parse(restaurantId)

	if err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	items, err := h.s.GetItems(parsedId)

	if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, *items)
}

// @Summary      Ban item
// @Tags         Items
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        item body ItemActionRequest true "Item to ban"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /items/ban [post]
func (h *Handler) banItem(c echo.Context) error {
	var data ItemActionRequest

	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	role := c.Get("role").(string)
	if role != "admin" {
		return Err2Json("only admins are permitted to do that", c, http.StatusUnauthorized)
	}

	itemId, err := uuid.Parse(data.ItemId)

	if err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	adminId := uuid.MustParse(c.Get("userID").(string))

	err = h.s.BanItem(itemId, adminId)
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "ban was successful",
	})
}

// @Summary      Unban item
// @Tags         Items
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        item body ItemActionRequest true "Item to unban"
// @Success      201 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /items/unban [post]
func (h *Handler) unbanItem(c echo.Context) error {
	var data ItemActionRequest

	if err := c.Bind(&data); err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	role := c.Get("role").(string)
	if role != "admin" {
		return Err2Json("only admins are permitted to do that", c, http.StatusUnauthorized)
	}

	itemId, err := uuid.Parse(data.ItemId)

	if err != nil {
		return Err2Json(err.Error(), c, http.StatusBadRequest)
	}

	adminId := uuid.MustParse(c.Get("userID").(string))

	err = h.s.UnbanItem(itemId, adminId)
	if err != nil {
		return Err2Json(err.Error(), c, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "ban was successful",
	})
}
