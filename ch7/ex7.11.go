package ch7

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%3.2f", d) }

type database map[string]dollars

//func (d database) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
//	//TODO implement me
//	switch request.URL.Path {
//	case "/list":
//		for item, price := range d {
//			fmt.Fprintf(writer, "%s: %v\n", item, price)
//		}
//	case "/price":
//		item := request.URL.Query().Get("item")
//		price, ok := d[item]
//		if !ok {
//			writer.WriteHeader(http.StatusNotFound)
//			fmt.Fprintf(writer, "no such item: %s\n", item)
//			return
//		}
//		fmt.Fprintf(writer, "%s: %v\n", item, price)
//	default:
//		writer.WriteHeader(http.StatusNotFound)
//		fmt.Fprintf(writer, "no such page: %s\n", request.URL)
//	}
//}

func (d database) list(w http.ResponseWriter, r *http.Request) {
	var itemList = template.Must(template.New("database").Parse(`
	<h1>database</h1>
	<table>
	<tr style='text-align: left'>
		<th>item</th>
		<th>price</th>
	</tr>
	{{range $item, $price := .}}
	<tr>
		<td>{{$item}}</td>
		<td>{{$price}}</td>
	</tr>
	{{end}}
	</table>
`))
	if err := itemList.Execute(w, d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	/*	for item, price := range d {
		fmt.Fprintf(w, "%s: %v\n", item, price)
	}*/
}

func (d database) price(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price, ok := d[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %s\n", item)
		return
	}
	fmt.Fprintf(w, "%s: %v", item, price)
}

func (d database) update(w http.ResponseWriter, r *http.Request) {
	item := r.URL.Query().Get("item")
	price := r.URL.Query().Get("price")
	fmt.Printf("%s: %v\n", item, price)

	if _, ok := d[item]; !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "no such item: %s\n", item)
		return
	} else {
		floatPrice, _ := strconv.ParseFloat(price, 32)
		d[item] = dollars(floatPrice)
		fmt.Fprintf(w, "update %s: %v\n", item, d[item])
	}
}

func main() {
	db := database{"shoes": 50, "socks": 5}
	//mux := http.NewServeMux()
	http.Handle("/list", http.HandlerFunc(db.list))
	http.Handle("/price", http.HandlerFunc(db.price))
	http.Handle("/update", http.HandlerFunc(db.update))
	log.Fatal(http.ListenAndServe("localhost: 8000", nil))
}
