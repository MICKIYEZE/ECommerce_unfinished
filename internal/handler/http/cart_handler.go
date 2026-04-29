package http

import (
    "net/http"

    "github.com/gin-gonic/gin"
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
    ProductID string `json:"product_id" binding:"required"`
    Quantity  int    `json:"quantity" binding:"required"`
}

func (h *CartHandler) AddItem(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    var req AddItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    productID, err := uuid.Parse(req.ProductID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
        return
    }

    err = h.cartService.AddItem(userID, productID, req.Quantity)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "item added"})
}

// ---------------------------
// Update Item
// ---------------------------

type UpdateItemRequest struct {
    ProductID string `json:"product_id" binding:"required"`
    Quantity  int    `json:"quantity" binding:"required"`
}

func (h *CartHandler) UpdateItem(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    var req UpdateItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    productID, err := uuid.Parse(req.ProductID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
        return
    }

    err = h.cartService.UpdateItem(userID, productID, req.Quantity)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "item updated"})
}

// ---------------------------
// Remove Item
// ---------------------------

func (h *CartHandler) RemoveItem(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    productIDStr := c.Param("productID")
    productID, err := uuid.Parse(productIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
        return
    }

    err = h.cartService.RemoveItem(userID, productID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}

// ---------------------------
// Clear Cart
// ---------------------------

func (h *CartHandler) ClearCart(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    err := h.cartService.ClearCart(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "cart cleared"})
}

// ---------------------------
// View Cart
// ---------------------------

func (h *CartHandler) GetCart(c *gin.Context) {
    userID := c.MustGet("userID").(uuid.UUID)

    cart, items, err := h.cartService.GetCart(userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "cart":  cart,
        "items": items,
    })
}
