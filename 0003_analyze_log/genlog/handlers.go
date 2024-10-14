package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

var store []Product

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func init() {
	store = make([]Product, 2)
	product := Product{
		ID:   11,
		Name: "一加 12",
	}
	product2 := Product{
		ID:   22,
		Name: "IPhone 17",
	}

	store = append(store, product, product2)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	pid, _ := strconv.Atoi(r.PathValue("pid"))

	var goods Product

	for i := 0; i < len(store); i++ {
		if store[i].ID == pid {
			goods = store[i]
		}
	}

	slog.Info(fmt.Sprintf("%s - %s", r.Method, r.RequestURI), slog.Int("product", pid))

	err := json.NewEncoder(w).Encode(&goods)
	if err != nil {
		slog.Error("写入product错误", "err", err.Error())
	}
}
