package main

import "fmt"

type Store interface {
	Add(p Producto)
	Get(id int) (Producto, error)
	List() []Producto
	Delete(id int) error
	Update(id int, nuevoPrecio float64) error
}

type error interface {
	Error() string
}

type NotFoundError struct {
	ID int
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("producto %d no encontrado", e.ID)
}

type MemoryStore struct {
	producto map[int]Producto
	nextID   int // es un autoincremental ID
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{ // EL & nos permite modificar la struct
		producto: make(map[int]Producto), // El make lo que hace es inicializar el MAP
		nextID:   1,
	}
}

func (s *MemoryStore) Add(p Producto) {
	p.ID = s.nextID
	s.producto[p.ID] = p
	s.nextID++
}

func (s *MemoryStore) Get(id int) (Producto, error) {
	p, ok := s.producto[id]
	if !ok {
		return Producto{}, fmt.Errorf("No se encontro el producto %d", id)
	}

	return p, nil // es como retornar null -> Estas diciendo que no hay errores
}

func (s *MemoryStore) List() []Producto {
	// Creamos un slice vacio con capacidad para todos los productos
	result := make([]Producto, 0, len(s.producto))

	for _, p := range s.producto {
		result = append(result, p)
	}

	return result
}

func (s *MemoryStore) Delete(id int) error {
	_, ok := s.producto[id]
	if !ok {
		return fmt.Errorf("NO se puede eliminar el producto %d", id)
	}

	delete(s.producto, id) // esta funcion es equivalente al map.remove en Java

	return nil // no hay errores
}

func (s *MemoryStore) Update(id int, nuevoPrecion float64) error {
	p, ok := s.producto[id]

	if !ok {
		return fmt.Errorf("No se encontro el producto %d", id)
	}

	if nuevoPrecion < 0 {
		fmt.Errorf("No se puede asignar un precio menor a 0  %.2f", nuevoPrecion)
	}

	p.Price = nuevoPrecion
	s.producto[id] = p

	return nil
}
