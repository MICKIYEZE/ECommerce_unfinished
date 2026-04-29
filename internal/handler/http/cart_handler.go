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

// ---------------------------
// Add Item
// ---------------------------

type AddItemRequest struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

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

// ---------------------------
// Update Item
// ---------------------------

type UpdateItemRequest struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

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

// ---------------------------
// Remove Item
// ---------------------------

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

// ---------------------------
// Clear Cart
// ---------------------------

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(uuid.UUID)

    if err := h.cartService.ClearCart(userID); err != nil {
        respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{"message": "cart cleared"})
}

// ---------------------------
// View Cart
// ---------------------------

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
