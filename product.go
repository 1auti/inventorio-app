package main

import (
	"fmt"
)

// Si el nombre de la estructura esta en MAYUS -> PUBLICO  | MINUSCULA --> PRIVADO ( Solo visible adentro del paquete)
type Producto struct {
	ID       int
	Name     string
	Quantity int
	Price    float64
}

// El New Producto cumple la funcion del constructor, GO no tiene constructores reales pero por convencion se usa de esta manera
func NewProducto(id int, name string, quantity int, price float64) Producto {
	return Producto{
		ID:       id,
		Name:     name,
		Quantity: quantity,
		Price:    price,
	}
}

// String() convierte en objeto en un texto legible
func (p Producto) String() string {
	return fmt.Sprintf("{%d} - %s - stock : %d - precio $%.2f", p.ID, p.Name, p.Quantity, p.Price)
}

// El metodo p es el que vinculo el metood con la clase Producto
func (p Producto) tieneStock() bool {
	return p.Quantity > 0
}

func (p Producto) valorStock() float64 {
	return float64(p.Quantity) * p.Price // el float64 convierte de int a float |GO NO HACE CONVERSACIONES IMPLICITAS
}

// Aplicamos el puntero apuntando al Original
func (p *Producto) aplicaDescuento(porcentaje float64) {
	p.Price = p.Price * (1 - porcentaje/100)
}

/*  Al aplicar el ( p * Producto ) modificamos al Original en cambio cuando p Producto modificamos una copia nunca se modifica EL ORIGINAL  */

func (p Producto) cambiarValorSinPuntero(valor float64) float64 {
	p.Price = valor
	return p.Price
}

func (p *Producto) cambiarValorConPuntero() float64 {
	p.Price = 1000
	return p.Price
}
