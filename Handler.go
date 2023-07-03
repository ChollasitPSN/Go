// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	_ "github.com/go-sql-driver/mysql"
// _ "golang.org/x/tools/go/analysis/passes/defers"
// _ "golang.org/x/tools/go/analysis/passes/nilfunc"
// )

// type Product struct {
// 	ID         int     `json: "ID"`
// 	Name       string  `json: "name"`
// 	Gender     string  `json: "gender"`
// 	Size       string  `json: "size"`
// 	Plaincolor string  `json: "plaincolor"`
// 	Pattern    string  `json: "pattern"`
// 	Figure     string  `json: "figure"`
// 	Price      float64 `json: "price"`
// }

// var Db *sql.DB
// var ProductList []Product

// const productPath = "products"
// const basePath = "/api"

// func init() {
// 	ProductJSON := `[
// 		{
// 			"ID":1,
// 			"name":"Hood",
// 			"Gender":"Male",
// 			"size":"S"
// 			"plaincolor":"White"
// 			"pattern":"Plan"
// 			"figure":"Fit"
// 			"price":"1500"
// 		},
// 		{
// 			"ID":2,
// 			"name":"Sweater",
// 			"Gender":"Female",
// 			"size":"XL"
// 			"plaincolor":"Blue"
// 			"pattern":"Dot"
// 			"figure":"Oversize"
// 			"price":"1000"
// 		},
// 		{
// 			"Id":3,
// 			"name":"Shirt",
// 			"Gender":"Male",
// 			"size":"XL"
// 			"plaincolor":"Pink"
// 			"pattern":"Plan"
// 			"figure":"Oversize"
// 			"price":"500"
// 		}
// 	]`
// 	err := json.Unmarshal([]byte(ProductJSON), &ProductList)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

// func getNextID() int {
// 	highestID := -1
// 	for _, product := range ProductList {
// 		if highestID < product.ID {
// 			highestID = product.ID
// 		}
// 	}
// 	return highestID + 1
// }

// func findID(ID int) (*Product, int) {
// 	for i, product := range ProductList {
// 		if product.ID == ID {
// 			return &product, i
// 		}
// 	}
// 	return nil, 0
// }

// func productHandler(w http.ResponseWriter, r *http.Request) {
// 	urlPathSegment := strings.Split(r.URL.Path, "product/")
// 	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
// 	if err != nil {
// 		log.Print(err)
// 		w.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	product, listItemIndex := findID(ID)
// 	if product == nil {
// 		http.Error(w, fmt.Sprintf("no product with ID %d", ID), http.StatusNotFound)
// 		return
// 	}
// 	switch r.Method {
// 	case http.MethodGet:
// 		productJSON, err := json.Marshal(product)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(productJSON)

// 	case http.MethodPut:
// 		var updateproduct Product
// 		byteBody, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		err = json.Unmarshal(byteBody, &updateproduct)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		if updateproduct.ID != ID {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		product = &updateproduct
// 		ProductList[listItemIndex] = *product
// 		w.WriteHeader(http.StatusOK)
// 		return
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 		return
// 	}

// }

// func productsHandler(w http.ResponseWriter, r *http.Request) {
// 	productJSON, err := json.Marshal(ProductList)
// 	switch r.Method {
// 	case http.MethodGet:
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write(productJSON)
// 	case http.MethodPost:
// 		var newProduct Product
// 		Bodybyte, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		err = json.Unmarshal(Bodybyte, &newProduct)
// 		if err != nil {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		if newProduct.ID != 0 {
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		newProduct.ID = getNextID()
// 		ProductList = append(ProductList, newProduct)
// 		w.WriteHeader(http.StatusCreated)
// 		return

// 	}

// }

// func middlewareHandler(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("before handler middle start")

// 		handler.ServeHTTP(w, r)
// 		fmt.Println("middle finish")
// 	})
// }

// func enableCorsMiddelware(handler http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Add("Access-Control-Allow-Origin", "*")
// 		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 		w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
// 		handler.ServeHTTP(w, r)
// 	})
// }

// func query(db *sql.DB) {

// 	var (
// 		id         int
// 		name       string
// 		plaincolor string
// 		pattern    string
// 		figure     string
// 		size       string
// 		price      float64
// 	)
// 	var inputID int
// 	fmt.Scan(&inputID)
// 	query := "SELECT id,name,plaincolor,pattern,figure,size,price FROM pruduct WHERE id = ?"
// 	if err := db.QueryRow(query, inputID).Scan(&id, &name, &plaincolor, &pattern, &figure, &size, &price); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(id, name, plaincolor, pattern, figure, size, price)

// }

// func creatingTable(db *sql.DB) {
// 	query := `CREATE TABLE users (
//   		id INT NOT NULL AUTO_INCREMENT,
//   		username TEXT NOT NULL,
//   		password TEXT NOT NULL,
//   		cart TEXT NULL,
//   		PRIMARY KEY (id)
// 	);`

// 	if _, err := db.Exec(query); err != nil {
// 		log.Fatal(err)
// 	}
// }
// func Insert(db *sql.DB) {
// 	var username string
// 	var password string
// 	var cart string
// 	fmt.Scan(&username)
// 	fmt.Scan(&password)
// 	fmt.Scan(&cart)
// 	result, err := db.Exec(`INSERT INTO users (username, password, cart) VALUES(?,?,?)`, username, password, cart)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	id, err := result.LastInsertId()
// 	fmt.Println(id)

// }

// func Delete(db *sql.DB) {
// 	var deleteid int
// 	fmt.Scan(&deleteid)
// 	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, deleteid)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }

// func SetupDB() {
// 	var err error
// 	Db, err := sql.Open("mysql", "root:th30122544@tcp(127.0.0.1:3306)/productdb")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(Db)
// 	Db.SetConnMaxLifetime(time.Minute * 3)
// 	Db.SetMaxOpenConns(10)
// 	Db.SetMaxIdleConns(10)
// }

// func getProductList() ([]Product, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	results, err := Db.QueryContext(ctx, `SELECT
// 	id,
// 	name,
// 	plaincolor,
// 	pattern,
// 	figure,
// 	size,
// 	price
// 	FROM category`)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return nil, err
// 	}
// 	defer results.Close()
// 	products := make([]Product, 0)
// 	for results.Next() {
// 		var product Product
// 		results.Scan(
// 			&product.ID,
// 			&product.Name,
// 			&product.Plaincolor,
// 			&product.Pattern,
// 			&product.Figure,
// 			&product.Size,
// 			&product.Price)
// 		products = append(products, product)
// 	}
// 	return products, nil
// }

// func SetupRoutes(apiBasePath string) {
// 	productsHandler := http.HandlerFunc(handleProduct)
// 	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productPath), corsMiddleware(productsHandler))

// }

// func InsertProduct(product Product) (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	result, err := Db.ExecContext(ctx, `INSERT INTO category
// 	(id,
// 	name,
// 	plaincolor,
// 	pattern,
// 	figure,
// 	size,
// 	price)
// 	VALUE (?,?,?,?,?,?,?)`,
// 		product.ID,
// 		product.Name,
// 		product.Plaincolor,
// 		product.Pattern,
// 		product.Figure,
// 		product.Size,
// 		product.Price)
// 	if err != nil {
// 		log.Println(err.Error())
// 		return 0, err
// 	}
// 	insertID, err := result.LastInsertId()
// 	if err != nil {
// 		log.Println(err.Error())
// 		return 0, err
// 	}
// 	return int(insertID), nil
// }

// func handleProduct(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		productList, err := getProductList()
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		j, err := json.Marshal(productList)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		_, err = w.Write(j)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	case http.MethodPost:
// 		var product Product
// 		err := json.NewDecoder(r.Body).Decode(&product)
// 		if err != nil {
// 			log.Print(err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		ProductID, err := InsertProduct(product)
// 		if err != nil {
// 			log.Print(err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}
// 		w.WriteHeader(http.StatusCreated)
// 		w.Write([]byte(fmt.Sprintf(`{"productid":%d}`, ProductID)))
// 	case http.MethodOptions:
// 		return
// 	default:
// 		w.WriteHeader(http.StatusMethodNotAllowed)
// 	}
// }

// func test() {
// productItemHandler := http.HandlerFunc(productHandler)
// productListHandler := http.HandlerFunc(productsHandler)
// http.Handle("/product/", enableCorsMiddelware(productItemHandler))
// http.Handle("/product", enableCorsMiddelware(productListHandler))
// http.ListenAndServe(":5000", nil)
// db, err := sql.Open("mysql", "root:th30122544@tcp(127.0.0.1:3306)/productdb")
// if err != nil {
// 	fmt.Println("failed to connect")
// } else {
// 	fmt.Println("connect success")
// }
//Delete(db)
//Insert(db)
//creatingTable(db)
//query(db)
// SetupDB()
// SetupRoutes(basePath)
// log.Fatal(http.ListenAndServe(":5000", nil))
// }