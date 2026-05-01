package http

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"

    "ecommerce/internal/service/cart"
)

type CartHandler struct {
    cartService cart.CartService
}

func NewCartHandler(cartService cart.CartService) *CartHandler {
    return &CartHandler{cartService: cartService}
}

// ===============================
// Add Item
// ===============================

// AddItemRequest represents the request body for adding an item to the cart
type AddItemRequest struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

// AddItem godoc
// @Summary Add item to cart
// @Description Adds a product to the authenticated user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body AddItemRequest true "Add Item Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /cart/add [post]
// @Security BearerAuth
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uuid.UUID)

    var req AddItemRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    productID, err := uuid.Parse(req.ProductID)
    if err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product ID"})
        return
    }

    if err := h.cartService.AddItem(userID, productID, req.Quantity); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{"message": "item added"})
}

// ===============================
// Update Item
// ===============================

type UpdateItemRequest struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

// UpdateItem godoc
// @Summary Update cart item quantity
// @Description Updates the quantity of a product in the user's cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body UpdateItemRequest true "Update Item Payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /cart/update [put]
// @Security BearerAuth
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uuid.UUID)

    var req UpdateItemRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    productID, err := uuid.Parse(req.ProductID)
    if err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product ID"})
        return
    }

    if err := h.cartService.UpdateItem(userID, productID, req.Quantity); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{"message": "item updated"})
}

// ===============================
// Remove Item
// ===============================

// RemoveItem godoc
// @Summary Remove item from cart
// @Description Removes a product from the authenticated user's cart
// @Tags Cart
// @Produce json
// @Param productID path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /cart/remove/{productID} [delete]
// @Security BearerAuth
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uuid.UUID)

    productIDStr := chi.URLParam(r, "productID")
    productID, err := uuid.Parse(productIDStr)
    if err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product ID"})
        return
    }

    if err := h.cartService.RemoveItem(userID, productID); err != nil {
        respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{"message": "item removed"})
}

// ===============================
// Clear Cart
// ===============================

// ClearCart godoc
// @Summary Clear cart
// @Description Removes all items from the authenticated user's cart
// @Tags Cart
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cart/clear [delete]
// @Security BearerAuth
func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uuid.UUID)

    if err := h.cartService.ClearCart(userID); err != nil {
        respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{"message": "cart cleared"})
}

// ===============================
// View Cart
// ===============================

// GetCart godoc
// @Summary Get user cart
// @Description Returns the user's cart and all items inside it
// @Tags Cart
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /cart [get]
// @Security BearerAuth
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uuid.UUID)

    cart, items, err := h.cartService.GetCart(userID)
    if err != nil {
        respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]interface{}{
        "cart":  cart,
        "items": items,
    })
}
