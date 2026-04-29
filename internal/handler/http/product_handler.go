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

// -----------------------------
// PUBLIC ROUTES
// -----------------------------

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

// -----------------------------
// ADMIN ROUTES
// -----------------------------

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
