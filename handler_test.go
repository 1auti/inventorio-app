package main

import (
	"bytes" // para crear un body de request
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest" // librería estándar para testear HTTP
	"testing"
)

// MockStore es nuestra implementación de Store para tests.
// No necesita Mockito — es un struct normal que implementa la interfaz.
// Equivalente Java: @Mock Store store — pero sin framework.
type MockStore struct {
	// Guardamos los productos en memoria — igual que MemoryStore
	productos map[int]Producto
	nextID    int

	// Flags para simular errores — controlamos el comportamiento desde el test
	failGet    bool
	failAdd    bool
	failDelete bool
}

func NewMockStore() *MockStore {
	return &MockStore{
		productos: make(map[int]Producto),
		nextID:    1,
	}
}

// Implementamos todos los métodos de Store — igual que MemoryStore
// pero con la capacidad de simular errores
func (m *MockStore) Add(p Producto) {
	if m.failAdd {
		return
	}
	p.ID = m.nextID
	m.productos[p.ID] = p
	m.nextID++
}

func (m *MockStore) Get(id int) (Producto, error) {
	if m.failGet {
		return Producto{}, fmt.Errorf("error simulado en Get")
	}
	p, ok := m.productos[id]
	if !ok {
		return Producto{}, fmt.Errorf("producto %d no encontrado", id)
	}
	return p, nil
}

func (m *MockStore) List() []Producto {
	result := make([]Producto, 0, len(m.productos))
	for _, p := range m.productos {
		result = append(result, p)
	}
	return result
}

func (m *MockStore) Delete(id int) error {
	if m.failDelete {
		return fmt.Errorf("error simulado en Delete")
	}
	_, ok := m.productos[id]
	if !ok {
		return fmt.Errorf("producto %d no encontrado", id)
	}
	delete(m.productos, id)
	return nil
}

func (m *MockStore) UpdatePrice(id int, price float64) error {
	p, ok := m.productos[id]
	if !ok {
		return fmt.Errorf("producto %d no encontrado", id)
	}
	p.Price = price
	m.productos[id] = p
	return nil
}

// --- Tests de los handlers ---

// TestHandlerListProductos testea GET /products
func TestHandlerListProductos(t *testing.T) {
	// Setup — store con datos
	store := NewMockStore()
	store.Add(NewProduct(0, "Notebook", 10, 500.0))
	store.Add(NewProduct(0, "Mouse", 25, 29.99))

	handler := NewHandler(store)

	// httptest.NewRequest crea un request de prueba sin red real.
	// Equivalente Spring: MockMvcRequestBuilders.get("/products")
	req := httptest.NewRequest(http.MethodGet, "/products", nil)

	// httptest.NewRecorder captura la respuesta sin red real.
	// Equivalente Spring: MockMvcResult
	rec := httptest.NewRecorder()

	// Llamamos al handler directamente — sin levantar servidor
	handler.handlerProductos(rec, req)

	// Verificamos el status code
	if rec.Code != http.StatusOK {
		t.Errorf("esperaba status 200, obtuve %d", rec.Code)
	}

	// Verificamos el Content-Type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("esperaba Content-Type=application/json, obtuve %s", contentType)
	}

	// Decodificamos el body de la respuesta
	var productos []Producto
	err := json.NewDecoder(rec.Body).Decode(&productos)
	if err != nil {
		t.Fatalf("error decodificando respuesta: %v", err)
	}

	if len(productos) != 2 {
		t.Errorf("esperaba 2 productos, obtuve %d", len(productos))
	}
}

// TestHandlerGetProducto testea GET /products/{id}
func TestHandlerGetProducto(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{name: "producto existente", url: "/products/1", wantStatus: 200},
		{name: "producto inexistente", url: "/products/99", wantStatus: 404},
		{name: "ID inválido", url: "/products/abc", wantStatus: 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewMockStore()
			store.Add(NewProducto(0, "Notebook", 10, 500.0))

			handler := NewHandler(store)

			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.handlerProduct(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("esperaba status %d, obtuve %d",
					tt.wantStatus, rec.Code)
			}
		})
	}
}

// TestHandlerAddProducto testea POST /products
func TestHandlerAddProducto(t *testing.T) {
	store := NewMockStore()
	handler := NewHandler(store)

	// Creamos el body del request como JSON
	// Equivalente Spring: MockMvcRequestBuilders.post().content(json)
	body := Producto{Name: "Teclado", Quantity: 15, Price: 79.99}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewReader(bodyBytes), // el body del request
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.handlerProductos(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("esperaba status 201, obtuve %d", rec.Code)
	}

	// Verificamos que se guardó en el store
	if len(store.List()) != 1 {
		t.Error("esperaba 1 producto en el store después del Add")
	}
}
