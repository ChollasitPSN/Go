package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

const productPath = "products"
const basepath = "/api"

type Product struct {
	ID         int     `json: "id"`
	Gender     string  `json: "gender"`
	Plaincolor string  `json: "plaincolor"`
	Pattern    string  `json: "pattern"`
	Figure     string  `json: "figure"`
	Size       string  `json: "size"`
	Price      float64 `json: "price"`
}

func SetupDB() {
	var err error
	Db, err = sql.Open("mysql", "root:th30122544@tcp(127.0.0.1:3306)/productdb")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(Db)
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	results, err := Db.QueryContext(ctx, `SELECT
	id,
	gender,
	plaincolor,
	pattern,
	figure,
	size,
	price
	FROM category`)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(
			&product.ID,
			&product.Gender,
			&product.Plaincolor,
			&product.Pattern,
			&product.Figure,
			&product.Size,
			&product.Price)
		products = append(products, product)
	}
	return products, nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := Db.ExecContext(ctx, `INSERT INTO category
	(id,
	gender,
	plaincolor,
	pattern,
	figure,
	size,
	price
	)
	VALUE (?,?,?,?,?,?,?)`,
		product.ID,
		product.Gender,
		product.Plaincolor,
		product.Pattern,
		product.Figure,
		product.Size,
		product.Price)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	return int(insertID), nil
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(productList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var product Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ProductID, err := insertProduct(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"productid":%d}`, ProductID)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleProduct(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, fmt.Sprintf("%s/", productPath))
	if len(urlPathSegments[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		product, err := getProduct(productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if product == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(product)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		err := removeProduct(productID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getProduct(productid int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	row := Db.QueryRowContext(ctx, `SELECT
	id,
	gender,
	plaincolor,
	pattern,
	figure,
	size,
	price
	FROM category
	WHERE id = ?`, productid)
	product := &Product{}
	err := row.Scan(
		&product.ID,
		&product.Gender,
		&product.Plaincolor,
		&product.Pattern,
		&product.Figure,
		&product.Size,
		&product.Price,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return product, nil
}

func removeProduct(productID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := Db.ExecContext(ctx, `DELETE FROM category where id = ?`, productID)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
		handler.ServeHTTP(w, r)
	})
}

func SetupRoutes(apiBasePath string) {
	productHandler := http.HandlerFunc(handleProduct)
	productsHandler := http.HandlerFunc(handleProducts)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productPath), corsMiddleware(productHandler))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productPath), corsMiddleware(productsHandler))
}

func main() {
	SetupDB()
	SetupRoutes(basepath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
