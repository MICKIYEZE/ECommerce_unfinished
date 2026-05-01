package http

import (
    "fmt"
    "net/http"

    "ecommerce/internal/domain/entity"
    productService "ecommerce/internal/service/product"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
)

type ProductHandler struct {
    service productService.ProductService
}

func NewProductHandler(service productService.ProductService) *ProductHandler {
    return &ProductHandler{service: service}
}

// ===============================
// PUBLIC ROUTES
// ===============================

// ListProducts godoc
// @Summary List all products
// @Description Returns all products with optional filtering (future use)
// @Tags Products
// @Produce json
// @Success 200 {array} entity.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
    filter := entity.ProductFilter{}

    products, err := h.service.List(filter)
    if err != nil {
        fmt.Println("LIST ERROR:", err)
        respondError(w, http.StatusInternalServerError, "failed to list products")
        return
    }

    respondJSON(w, http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Returns a single product by its UUID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid product ID")
        return
    }

    product, err := h.service.GetByID(id)
    if err != nil {
        respondError(w, http.StatusNotFound, "product not found")
        return
    }

    respondJSON(w, http.StatusOK, product)
}

// ===============================
// ADMIN ROUTES
// ===============================

// CreateProduct godoc
// @Summary Create a new product
// @Description Admin-only: Creates a new product
// @Tags Products
// @Accept json
// @Produce json
// @Param request body entity.Product true "Product Payload"
// @Success 201 {object} entity.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/products [post]
// @Security BearerAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    var req entity.Product

    if err := decodeJSON(r, &req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    if err := h.service.Create(&req); err != nil {
        respondError(w, http.StatusInternalServerError, "failed to create product")
        return
    }

    respondJSON(w, http.StatusCreated, req)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Admin-only: Updates product fields by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body entity.Product true "Updated Product Payload"
// @Success 200 {object} entity.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/products/{id} [put]
// @Security BearerAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid product ID")
        return
    }

    var req entity.Product
    if err := decodeJSON(r, &req); err != nil {
        respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    req.ID = id

    if err := h.service.Update(&req); err != nil {
        respondError(w, http.StatusInternalServerError, "failed to update product")
        return
    }

    respondJSON(w, http.StatusOK, req)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Admin-only: Deletes a product by ID
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/products/{id} [delete]
// @Security BearerAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        respondError(w, http.StatusBadRequest, "invalid product ID")
        return
    }

    if err := h.service.Delete(id); err != nil {
        respondError(w, http.StatusInternalServerError, "failed to delete product")
        return
    }

    respondJSON(w, http.StatusOK, map[string]string{"message": "product deleted"})
}
