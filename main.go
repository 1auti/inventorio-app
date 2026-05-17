package main

import (
	"fmt"
	"net/http"
)

func main() {
	// HTTP HANDLER registra un handler para un ruta -> Seria equivalente como GetMapping
	// 1er parametro es la ruta | 2do parametro es el handler de la funcion
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
		// Response writer lo que escribis que va al cliente | Equivalente al return del metodo
		// Request es equivalente a Request body - Path variable , etc
	})

	// El Listen and Server lo que hace es arrancar el servidor en el puerto 8080
	fmt.Println("Servidor corriendo en el puerto 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor", err) // Si esta ocupado tira el error
	}

	/*
		var store Store = NewMemoryStore()

		// Cargamos productos iniciales
		store.Add(NewProducto(0, "Notebook", 10, 500.0))
		store.Add(NewProducto(0, "Mouse", 25, 29.99))
		store.Add(NewProducto(0, "Teclado", 15, 79.99))

		// --- UpdatePrice: caso exitoso ---
		fmt.Println("=== Actualizar precio ===")
		err := store.Update(1, 450.0)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			p, _ := store.Get(1) // _ descarta el error — sabemos que existe
			fmt.Println("Precio actualizado:", p)
		}

		// --- UpdatePrice: precio inválido ---
		err = store.Update(1, -100)
		if err != nil {
			fmt.Println("Error esperado:", err) // → precio inválido: -100.00
		}

		// --- Delete: caso exitoso ---
		fmt.Println("\n=== Eliminar producto ===")
		err = store.Delete(2)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Producto 2 eliminado")
		}

		// --- Delete: ID inexistente ---
		err = store.Delete(99)
		if err != nil {
			fmt.Println("Error esperado:", err) // → no se puede eliminar: producto 99
		}

		// --- Estado final del inventario ---
		fmt.Println("\n=== Inventario final ===")
		for _, p := range store.List() {
			fmt.Println(p)
		}

	*/
}
