package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

type chaveContexto string

const chaveCliente chaveContexto = "cliente"

func clienteDaRequisicao(r *http.Request) ApiClient {
	c, _ := r.Context().Value(chaveCliente).(ApiClient)
	return c
}

func limitarCorpo(maxBytes int64, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		next.ServeHTTP(w, r)
	})
}

func autenticar(clientes *ApiClientDAO, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chave := r.Header.Get("X-API-KEY")
		if chave == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		cliente, err := clientes.FindByKey(chave)
		if errors.Is(err, ErrClienteNaoEncontrado) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if err != nil {
			erroInterno(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), chaveCliente, cliente)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type LimitadorPorCliente struct {
	mu          sync.Mutex
	limitadores map[int]*rate.Limiter
	taxa        rate.Limit
	rajada      int
}

func NewLimitadorPorCliente(porSegundo float64, rajada int) *LimitadorPorCliente {
	return &LimitadorPorCliente{
		limitadores: make(map[int]*rate.Limiter),
		taxa:        rate.Limit(porSegundo),
		rajada:      rajada,
	}
}

func (l *LimitadorPorCliente) limitador(clienteID int) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	lim, existe := l.limitadores[clienteID]
	if !existe {
		lim = rate.NewLimiter(l.taxa, l.rajada)
		l.limitadores[clienteID] = lim
	}
	return lim
}

func (l *LimitadorPorCliente) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cliente := clienteDaRequisicao(r)
		if !l.limitador(cliente.ID).Allow() {
			log.Printf("limite de taxa atingido: cliente %q", cliente.Name)
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
