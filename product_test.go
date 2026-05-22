package main

import "testing"

// TODOS LOS TEST EMPEIZAN CON Test
func TestNewProducto(t *testing.T) {
	// Creamos un producto con nuestro constructor
	p := NewProducto(0, "Notebook", 10, 500.0)

	// t.Errorf marca el test como fallido pero CONTINÚA ejecutando.
	// Equivalente Java: assertEquals() — pero sin lanzar excepción.
	if p.Name != "Notebook" {
		t.Errorf("esperaba Name='Notebook', obtuve '%s'", p.Name)
	}

	if p.Price != 500.0 {
		t.Errorf("esperaba Price=500.0, obtuve %.2f", p.Price)
	}

	if p.Quantity != 10 {
		t.Errorf("esperaba Quantity=10, obtuve %d", p.Quantity)
	}
}

func TestProductHasStock(t *testing.T) {
	// Caso: tiene stock
	p := NewProducto(0, "Notebook", 10, 500.0)
	if !p.tieneStock() {
		t.Error("esperaba HasStock()=true con quantity=10")
		// t.Error es igual a t.Errorf pero sin formato
	}

	// Caso: sin stock
	sinStock := NewProducto(0, "Agotado", 0, 100.0)
	if sinStock.tieneStock() {
		t.Error("esperaba HasStock()=false con quantity=0")
	}
}

func TestProductTotalValue(t *testing.T) {
	p := NewProducto(0, "Notebook", 10, 500.0)

	// 10 unidades * $500 = $5000
	expected := 5000.0
	got := p.valorStock()

	if got != expected {
		t.Errorf("esperaba TotalValue()=%.2f, obtuve %.2f", expected, got)
	}
}

func TestApplyDiscount(t *testing.T) {
	p := NewProducto(0, "Notebook", 10, 500.0)

	// Aplicamos 10% de descuento → $450
	p.aplicaDescuento(10)

	expected := 450.0
	if p.Price != expected {
		t.Errorf("esperaba Price=%.2f después del descuento, obtuve %.2f",
			expected, p.Price)
	}
}
