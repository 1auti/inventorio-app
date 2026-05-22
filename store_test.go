package main

import "testing"

// TestMemoryStoreAdd usa table-driven tests —
// el patrón más idiomático de Go para múltiples casos.
// Equivalente Java: @ParameterizedTest con @MethodSource
func TestMemoryStoreAdd(t *testing.T) {
	// Definimos una tabla de casos como slice de structs anónimos.
	// Cada struct es un caso de test con sus datos y resultado esperado.
	tests := []struct {
		name    string // nombre del caso — aparece en el output
		product Producto
		wantID  int // ID esperado después del Add
	}{
		{
			name:    "primer producto obtiene ID 1",
			product: NewProducto(0, "Notebook", 10, 500.0),
			wantID:  1,
		},
		{
			name:    "segundo producto obtiene ID 2",
			product: NewProducto(0, "Mouse", 5, 29.99),
			wantID:  2,
		},
	}

	// Un store nuevo por test — estado limpio siempre
	store := NewMemoryStore()

	for _, tt := range tests {
		// t.Run crea un subtest con el nombre del caso.
		// Equivalente Java: el nombre del @ParameterizedTest
		t.Run(tt.name, func(t *testing.T) {
			store.Add(tt.product)

			// Buscamos el producto que acabamos de agregar
			got, err := store.Get(tt.wantID)
			if err != nil {
				t.Fatalf("no se encontró el producto con ID %d: %v",
					tt.wantID, err)
				// t.Fatalf es como t.Errorf pero DETIENE el test inmediatamente
				// Equivalente Java: fail() o assumeTrue()
			}

			if got.ID != tt.wantID {
				t.Errorf("esperaba ID=%d, obtuve ID=%d", tt.wantID, got.ID)
			}
		})
	}
}

func TestMemoryStoreGet(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool // true si esperamos un error
	}{
		{name: "producto existente", id: 1, wantErr: false},
		{name: "producto inexistente", id: 99, wantErr: true},
		{name: "ID negativo", id: -1, wantErr: true},
	}

	// Setup — store con un producto
	store := NewMemoryStore()
	store.Add(NewProducto(0, "Notebook", 10, 500.0))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.Get(tt.id)

			// Verificamos si el error coincide con lo esperado
			if tt.wantErr && err == nil {
				t.Error("esperaba un error pero no hubo ninguno")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("no esperaba error pero obtuve: %v", err)
			}
		})
	}
}

func TestMemoryStoreDelete(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{name: "eliminar existente", id: 1, wantErr: false},
		{name: "eliminar inexistente", id: 99, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup dentro del subtest — store fresco para cada caso
			store := NewMemoryStore()
			store.Add(NewProducto(0, "Notebook", 10, 500.0))

			err := store.Delete(tt.id)

			if tt.wantErr && err == nil {
				t.Error("esperaba error pero no hubo ninguno")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("no esperaba error pero obtuve: %v", err)
			}
		})
	}
}
