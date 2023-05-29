package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/go-chi/chi"
)

func TestMyEndpoint(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/insert-info/{person}/{index}", injeccion) // Supongamos que esta es la ruta que quieres probar

	// Simular solicitud
	req, err := http.NewRequest("POST", "/insert-info/all/straw", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crear objeto de respuesta falso
	recorder := httptest.NewRecorder()

	// Ejecutar la prueba
	r.ServeHTTP(recorder, req)

	// Verificar la respuesta
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code esperado %v, pero obtuvo %v", http.StatusOK, status)
	}

	// Verificar el cuerpo de la respuesta
	expected := "Insert finish" // Supongamos que este es el cuerpo esperado
	if recorder.Body.String() != expected {
		t.Errorf("Cuerpo esperado %q, pero obtuvo %q", expected, recorder.Body.String())
	}
}

func TestMyEndpoint1(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/insert-info/{person}/{index}", injeccion) // Supongamos que esta es la ruta que quieres probar

	// Simular solicitud
	req, err := http.NewRequest("POST", "/insert-info/zufferli-j/all_documents", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crear objeto de respuesta falso
	recorder := httptest.NewRecorder()

	// Ejecutar la prueba
	r.ServeHTTP(recorder, req)

	// Verificar la respuesta
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code esperado %v, pero obtuvo %v", http.StatusOK, status)
	}

	// Verificar el cuerpo de la respuesta
	expected := "Insert finish" // Supongamos que este es el cuerpo esperado
	if recorder.Body.String() != expected {
		t.Errorf("Cuerpo esperado %q, pero obtuvo %q", expected, recorder.Body.String())
	}
}
func TestMyEndpoint2(t *testing.T) {
	r := chi.NewRouter()
	r.Post("/insert-info/{person}/{index}", injeccion) // Supongamos que esta es la ruta que quieres probar

	// Simular solicitud
	req, err := http.NewRequest("POST", "/insert-info/schoolcraft-d/tw_fuel_sales", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crear objeto de respuesta falso
	recorder := httptest.NewRecorder()

	// Ejecutar la prueba
	r.ServeHTTP(recorder, req)

	// Verificar la respuesta
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Status code esperado %v, pero obtuvo %v", http.StatusOK, status)
	}

	// Verificar el cuerpo de la respuesta
	expected := "Insert finish" // Supongamos que este es el cuerpo esperado
	if recorder.Body.String() != expected {
		t.Errorf("Cuerpo esperado %q, pero obtuvo %q", expected, recorder.Body.String())
	}
}
