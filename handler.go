package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// para serializar/deserializar JSON

// para manipular strings en las rutas

type Handler struct {
	store Store // Las Dependencias que inyectan en el constructor seria el equivalente al repositorio / service
}

// Creamos el handler que vamos a usar con la inteccionde dependencias
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes() {
	// Chain aplica Logger y Recovery a cada handler.
	// El orden importa: Logger es el más externo — ve el request primero.
	// Equivalente Spring: addInterceptor(logging).addInterceptor(recovery)
	http.HandleFunc("/products", Chain(
		h.handlerProductos,
		Logger,   // 1° — mide tiempo y loggea
		Recovery, // 2° — atrapa panics
	))

	http.HandleFunc("/products/", Chain(
		h.handlerProduct,
		Logger,
		Recovery,
	))

	// /health sin middleware — no necesita logging
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
}

func (h *Handler) handlerProductos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listProductos(w, r)
	case http.MethodPost:
		h.addProducto(w, r)
	default:
		http.Error(w, "Metodo no esta permitido", http.StatusMethodNotAllowed)
	}
}

// ESTA SERIA PARA METODOS /PRODUCTO/ID
func (h *Handler) handlerProduct(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: // GET /products/1
		h.getProducto(w, r)
	case http.MethodDelete: // DELETE /products/1
		h.deleteProducto(w, r)
	default:
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
	}
}

// Helpers

// Lo que hace es serializar a JSON y lo escribe en el response
// EQUIVALENTE RESPONSE BODY
func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json") // Le estamos diciendo que la respuesta es en JSON
	w.WriteHeader(status)                              // Le pasamos el status de la respuesta
	json.NewEncoder(w).Encode(data)                    // Escribimos en la respuesta el JSON
}

// Extraemos el ID de la URL
func extrearID(r *http.Request, prefix string) (int, error) {
	// Extraemos el ID de la URL
	idStr := strings.TrimPrefix(r.URL.Path, prefix)

	return strconv.Atoi(idStr) // Convertimos el id string a id int
}

// HANDLERS DE CADA ENDPONINT

func (h *Handler) listProductos(w http.ResponseWriter, r *http.Request) {
	productos := h.store.List()
	writeJson(w, http.StatusOK, productos)
}

func (h *Handler) addProducto(w http.ResponseWriter, r *http.Request) {
	var p Producto

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "ID INVALIDO", http.StatusBadRequest)
		return
	}

	h.store.Add(p)
	writeJson(w, http.StatusCreated, p)
}

func (h *Handler) getProducto(w http.ResponseWriter, r *http.Request) {
	id, err := extrearID(r, "/products/")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest) // 400
		return
	}

	p, err := h.store.Get(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // 404
		return
	}

	writeJson(w, http.StatusOK, p) // 200
}

// deleteProduct maneja DELETE /products/{id}
func (h *Handler) deleteProducto(w http.ResponseWriter, r *http.Request) {
	id, err := extrearID(r, "/products/")
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest) // 400
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // 404
		return
	}

	// 204 No Content — operación exitosa sin body de respuesta
	// Equivalente Spring: ResponseEntity.noContent().build()
	w.WriteHeader(http.StatusNoContent)
}
