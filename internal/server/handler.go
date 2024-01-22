package server

import (
	"encoding/json"
	"github.com/go-playground/pure/v5"
	"net/http"
	"text/template"
)

type handler struct {
	s storer
}

func newHandler(s storer) *handler {
	return &handler{s}
}

func (h *handler) getOrder(w http.ResponseWriter, r *http.Request) {
	oid := pure.RequestVars(r).URLParam("id")

	o, found := h.s.Get(oid)
	if !found {
		http.Error(w, "order_id not found "+oid, http.StatusBadRequest)
		return
	}

	var orderData interface{}
	err := json.Unmarshal(o, &orderData)
	if err != nil {
		http.Error(w, "error decoding JSON", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.MarshalIndent(orderData, "", "  ")
	if err != nil {
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *handler) getIndex(w http.ResponseWriter, r *http.Request) {
	indexPage := `
	<h2>Список поступивших заказов{{ if not . }} пуст. {{ else }}:{{ end }}</h2>
	<ul>
		{{range .}}
			<li><a href="/order/{{.}}">{{.}}</a></li>
		{{end}}
	</ul>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	templ := template.Must(template.New("index").Parse(indexPage))
	templ.Execute(w, h.s.GetAll())
}
