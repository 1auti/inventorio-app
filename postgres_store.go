package main

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx" // _   <--- ese simbolo sirve para decrile al compilador que sirve ese import
	// No usamos ninguna función de él directamente —
	// pero su init() registra el driver "pgx" en database/sql
)

type PostgresStore struct {
	db *sqlx.DB // Esto es la conexion de la base de datos,maneja los pool de conexiones internamente
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("Error al abriendo la db: %w", err) // Cumple la funcion de encapsular el error %w
	}

	// Verifica que abre la conexion real y que la base de datos responda
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Error conectando a db: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Add(p Producto) {
	query := `
	INSERT INTO products (name, quantity, price)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := s.db.QueryRow(query, p.Name, p.Quantity, p.Price).Scan(&p.ID)
	if err != nil {
		fmt.Errorf("Error al ejecutar la query %w", err)
	}
}

func (s *PostgresStore) List() []Producto {
	var products []Producto // Mapea muchs filas directamente en el slice | equivalente al findAll

	query := `SELECT id, name, quantity price FROM products ORDER BY id`
	err := s.db.Select(&products, query)
	if err != nil {
		fmt.Println("Error al ejecutar la query", err)
		return []Producto{}
	}

	return products
}

func (s *PostgresStore) Get(id int) (Producto, error) {
	var p Producto

	// sqlx.Get es el helper estrella de sqlx —
	// ejecuta la query y mapea el resultado directamente al struct.
	// Equivalente JPA: findById(id) — pero con SQL explícito.
	//
	// Para que funcione, los nombres de columna SQL deben coincidir
	// con los tags `db:` del struct Producto
	query := `SELECT id, name, quantity, price FROM products WHERE id = $1`
	err := s.db.Get(&p, query, id)
	if err != nil {
		return Producto{}, fmt.Errorf("producto %d no encontrado", id)
	}

	return p, nil
}

func (s *PostgresStore) Delete(id int) error {
	// Exec ejecuta queries que no retornan filas (INSERT, UPDATE, DELETE).
	// Equivalente JPA: deleteById(id)
	query := `DELETE FROM products WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error eliminando producto %d: %w", id, err)
	}

	// RowsAffected nos dice cuántas filas se modificaron.
	// Si es 0, el ID no existía — equivalente a Optional.isEmpty() en JPA.
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("producto %d no encontrado", id)
	}

	return nil
}

func (s *PostgresStore) Update(id int, newPrice float64) error {
	if newPrice <= 0 {
		return fmt.Errorf("precio inválido: %.2f", newPrice)
	}

	query := `UPDATE products SET price = $1 WHERE id = $2`
	result, err := s.db.Exec(query, newPrice, id)
	if err != nil {
		return fmt.Errorf("error actualizando precio: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("producto %d no encontrado", id)
	}

	return nil
}
