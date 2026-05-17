package main

import (
	"fmt"
	"net/http"
	"time" // para medir duración del request
)

// MiddlewareFunc es nuestro tipo para middlewares.
// Recibe un handler y retorna un handler envuelto.
// Equivalente Java: Function<HandlerMethod, HandlerMethod>
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// Logger registra método, ruta y duración de cada request.
// Equivalente Spring: HandlerInterceptor con preHandle + postHandle.
//
// La firma sigue el patrón: recibe el siguiente handler, retorna uno nuevo.
func Logger(next http.HandlerFunc) http.HandlerFunc {
	// Retornamos una función nueva que envuelve a "next"
	return func(w http.ResponseWriter, r *http.Request) {
		// ANTES del request — medimos el tiempo de inicio
		start := time.Now()

		// Llamamos al handler real — acá se procesa el request
		next(w, r)

		// DESPUÉS del request — calculamos cuánto tardó
		// time.Since(start) = tiempo transcurrido desde start
		// Equivalente Java: System.currentTimeMillis() - startTime
		duration := time.Since(start)

		fmt.Printf("[%s] %s %s — %v\n",
			time.Now().Format("15:04:05"), // hora actual HH:mm:ss
			r.Method,                      // GET, POST, DELETE...
			r.URL.Path,                    // /products, /products/1...
			duration,                      // cuánto tardó: 125µs, 1.2ms...
		)
	}
}

// Recovery atrapa panics y los convierte en respuestas 500.
// En Go un panic es como un error irrecuperable — si no lo atrapás, el servidor muere.
// Equivalente Spring: @ControllerAdvice con @ExceptionHandler
func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// defer se ejecuta cuando la función termina — pase lo que pase.
		// recover() atrapa un panic si hubo uno.
		// Equivalente Java: finally { } + catch(Throwable t)
		defer func() {
			if err := recover(); err != nil {
				// Si hubo panic, respondemos 500 en lugar de caer
				fmt.Printf("[PANIC] %v\n", err)
				http.Error(w,
					"error interno del servidor",
					http.StatusInternalServerError, // 500
				)
			}
		}()

		next(w, r) // si esto genera un panic, defer lo atrapa
	}
}

// Chain aplica múltiples middlewares en orden, de afuera hacia adentro.
// chain(handler, Logger, Recovery) → Logger(Recovery(handler))
// Equivalente Spring: cadena de interceptores en orden de registro.
//
// "...MiddlewareFunc" es variadic — acepta cualquier cantidad de middlewares.
// Equivalente Java: Object... args
func Chain(handler http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	// Aplicamos los middlewares en orden inverso
	// para que el primero de la lista sea el más externo
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
