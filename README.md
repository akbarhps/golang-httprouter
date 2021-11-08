# Golang HTTP Router

Sumber Tutorial:
[Udemy](https://www.udemy.com/course/pemrograman-go-lang-pemula-sampai-mahir/learn/lecture/27069708#overview) |
[Slide](https://docs.google.com/presentation/d/1RaRkNSeaXopQXODvOANJiNyQ8ECg8QYJE8RvVrrjedg/edit#slide=id.p)


## Pengenalan HTTP Router
---

- HttpRouter merupakan salah satu OpenSource Library yang populer untuk Http Handler di Go-Lang
- HttpRouter terkenal dengan kecepatannya dan juga sangat minimalis
- Hal ini dikarenakan HttpRouter hanya memiliki fitur untuk routing saja, tidak memiliki fitur apapun selain itu
- https://github.com/julienschmidt/httprouter 


### Menambahkan HTTP Router ke Project

```bash
go get github.com/julienschmidt/httprouter

go get github.com/stretchr/testify
```


### Kode: go.mod

```go
module golang-httprouter

go 1.17

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.1.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)
```


## Router
---

- Inti dari library HttpRouter adalah struct Router
- Router ini merupakan implementasi dari http.Handler, sehingga kita bisa dengan mudah menambahkan ke dalam http.Server
- Untuk membuat Router, kita bisa menggunakan function httprouter.New(), yang akan mengembalikan Router pointer


### Kode: Router

```go
router := httprouter.New()

server := http.Server{
    Addr:    ":8080",
    Handler: router,
}

err := server.ListenAndServe()
if err != nil {
    panic(err)
}
```


### HTTP Method

- Router mirip dengan `ServeMux`, dimana kita bisa menambahkan route ke dalam Router
- Kelebihan dibandingkan dengan `ServeMux` adalah, pada Router, kita bisa menentukan HTTP Method yang ingin kita gunakan, misal GET, POST, PUT, dan lain-lain
- Cara menambahkan route ke dalam Router adalah gunakan function yang sama dengan HTTP Method nya, misal `router.GET()`, `router.POST()`, dan lain-lain


### `httprouter.Handler`

- Saat kita menggunakan `ServeMux`, ketika menambah route, kita bisa menambahkan http.Handler
- Berbeda dengan Router, pada Router kita tidak menggunakan `http.Handler` lagi, melainkan menggunakan type `httprouter.Handle`
- Perbedaan dengan `http.Handler` adalah, pada httprouter.Handle, terdapat parameter ke tiga yaitu `Params`, yang akan kita bahas nanti di chapter tersendiri


```go
// Handle is a function that can be registered to a route to handle HTTP requests.
// Like http.HandlerFunc, but has a third parameter for the values of wildcards (variables).
type Handle func(http.ResponseWriter, *http.Request, Params)
```


### Kode: Route

```go
func TestRouter(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
        fmt.Fprint(writer, "Hello HTTP Router")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	t.Log(string(responseBody))
}
```


### Params

- `httprouter.Handle` memiliki parameter yang ketiga, yaitu Params. Untuk apa kegunaan Params?
- Params merupakan tempat untuk menyimpan parameter yang dikirim dari client
- Namun Params ini bukan query parameter, melainkan parameter di URL
- Kadang kita butuh membuat URL yang tidak fix, alias bisa berubah-ubah, misal /products/1, /products/2, dan seterusnya
- `ServeMux` tidak mendukung hal tersebut, namun Router mendukung hal tersebut
- Parameter yang dinamis yang terdapat di URL, secara otomatis dikumpulkan di Params
- Namun, agar Router tahu, kita harus memberi tahu ketika menambahkan Route, dibagian mana kita akan buat URL path nya menjadi dinamis


### Kode: Params

```go
func TestRouterParams(t *testing.T) {
	router := httprouter.New()
	router.GET("/products/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		t.Log(params)
		fmt.Fprintf(writer, "You requested a product with id %s", params.ByName("id"))
	})

	productId := "1"
	request := httptest.NewRequest(http.MethodGet, "/products/"+productId, nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	responseBody, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200")
	assert.Equal(t, "You requested a product with id "+productId, string(responseBody), "response body should be 'You requested a product with id 1'")
	t.Log(string(responseBody))
}
```


## Router Pattern
---

- Sekarang kita sudah tahu bahwa dengan menggunakan Router, kita bisa menambah params di URL
- Sekarang pertanyaannya, bagaimana pattern (pola) pembuatan parameter nya?


### Named Parameter

- Named parameter adalah pola pembuatan parameter dengan menggunakan nama
- Setiap nama parameter harus diawali dengan : (titik dua), lalu diikuti dengan nama parameter
- Contoh, jika kita memiliki pattern seperti ini :

| Pattern             | `/user/:user` |
| ------------------- | ------------- |
| `/user/eko`         | match         |
| `/user/you`         | match         |
| `/user/eko/profile` | not match     |
| `/user/`            | not match     |


### Catch All Parameter

- Selain named parameter, ada juga yang bernama catch all parameter, yaitu menangkap semua parameter
- Catch all parameter harus diawali dengan * (bintang), lalu diikuti dengan nama parameter
- Catch all parameter harus berada di posisi akhir URL

| Pattern                | `src/*filepath` |
| ---------------------- | --------------- |
| `/src/`                | not match       |
| `/src/somefile`        | match           |
| `/src/subdir/somefile` | match           |


### Kode: Named Parameter

```go
func TestRouterPatternNamedParameter(t *testing.T) {
	router := httprouter.New()
	router.GET("/products/:id/items/:itemId", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		t.Log(params)

		id := params.ByName("id")
		itemId := params.ByName("itemId")
		fmt.Fprintf(writer, "product id: %s, item id: %s", id, itemId)
	})

	request := httptest.NewRequest("GET", "/products/1/items/2", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "product id: 1, item id: 2", string(responseBody))

	t.Log(string(responseBody))
}
```


### Kode: Catch All Parameter

```go
func TestRouterPatternCatchAllParameter(t *testing.T) {
	router := httprouter.New()
	router.GET("/images/*image", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		t.Log(params)

		image := params.ByName("image")
		fmt.Fprintf(writer, "image: %s", image)
	})

	request := httptest.NewRequest("GET", "/images/resources/img.jpg", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "image: /resources/img.jpg", string(responseBody))

	t.Log(string(responseBody))
}
```


## ServeFile
---

- Pada materi Go-Lang Web, kita sudah pernah membahas tentang Serve File
- Pada Router pun, mendukung serve static file menggunakan function `ServeFiles(Path, FileSystem)`
- Dimana pada Path, kita harus menggunakan Catch All Parameter
- Sedangkan pada FileSystem kita bisa melakukan manual load dari folder atau menggunakan golang embed, seperti yang pernah kita bahas di materi Go-Lang Web


### Kode: ServeFile

```go
func TestRouterServeFile(t *testing.T) {
	router := httprouter.New()

	// create a dir destination so we don't need
	// to specify the file folder, i.e resources
	dir, _ := fs.Sub(resources, "resources")

	// *filepath is hardcoded in ServeFile
	// so its must be named *filepath
	router.ServeFiles("/files/*filepath", http.FS(dir))

	request := httptest.NewRequest(http.MethodGet, "/files/test.txt", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "test", string(responseBody))

	t.Log(string(responseBody))
}
```


## Panic Handler
---

- Apa yang terjadi jika terjadi panic pada logic Handler yang kita buat?
- Secara otomatis akan terjadi error, dan web akan berhenti mengembalikan response
- Kadang saat terjadi panic, kita ingin melakukan sesuatu, misal memberitahu jika terjadi kesalahan di web, atau bahkan mengirim informasi log kesalahan yang terjadi
- Sebelumnya, seperti yang sudah kita bahas di materi Go-Lang Web, jika kita ingin menangani panic, kita harus membuat Middleware khusus secara manual
- Namun di Router, sudah disediakan untuk menangani panic, caranya dengan menggunakan attribute 

```go 
PanicHandler : func(http.ResponseWriter, *http.Request, interface{})
```


### Kode: Panic Handler

```go
func TestRouterPanicHandler(t *testing.T) {
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		t.Logf("Panic : %s", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal Server Error\nError : %s", err)
	}
	router.GET("/panic", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		panic("Oops, something went wrong")
	})

	request := httptest.NewRequest("GET", "/panic", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	responseBody, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, "Internal Server Error\nError : Oops, something went wrong", string(responseBody))

	t.Log(string(responseBody))
}
```


## Not Found Handler
---

- Selain panic handler, Router juga memiliki not found handler
- Not found handler adalah handler yang dieksekusi ketika client mencoba melakukan request URL yang memang tidak terdapat di Router
- Secara default, jika tidak ada route tidak ditemukan, Router akan melanjutkan request ke `http.NotFound`, namun kita bisa mengubah nya
- Caranya dengan mengubah `router.NotFound = http.Handler`


### Kode: Not Found Handler

```go
func TestRouterNotFound(t *testing.T) {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Not Found")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	result := response.Result()
	assert.Equal(t, http.StatusNotFound, result.StatusCode)

	responseBody, _ := ioutil.ReadAll(result.Body)
	assert.Equal(t, "404 Not Found", string(responseBody))

	t.Log(string(responseBody))
}
```


## Method Not Allowed Handler
---

- Saat menggunakan ServeMux, kita tidak bisa menentukan HTTP Method apa yang digunakan untuk Handler
- Namun pada Router, kita bisa menentukan HTTP Method yang ingin kita gunakan, lantas apa yang terjadi jika client tidak mengirim HTTP Method sesuai dengan yang kita tentukan? 
- Maka akan terjadi error Method Not Allowed
- Secara default, jika terjadi error seperti ini, maka Router akan memanggil function http.Error
- Jika kita ingin mengubahnya, kita bisa gunakan router.MethodNotAllowed = http.Handler


### Kode: Method Not Allowed Handler

```go
func TestRouterMethodNotAllowed(t *testing.T) {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		fmt.Fprint(w, "Method not allowed")
	})
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	request := httptest.NewRequest(http.MethodPost, "/", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	result := response.Result()
	assert.Equal(t, http.StatusTeapot, result.StatusCode)

	resultBody, _ := ioutil.ReadAll(result.Body)
	assert.Equal(t, "Method not allowed", string(resultBody))

	t.Log(string(resultBody))
}
```


## Middleware
---

- HttpRouter hanyalah library untuk http router saja, tidak ada fitur lain selain router
- Dan karena Router merupakan implementasi dari http.Handler, jadi untuk middleware, kita bisa membuat sendiri, seperti yang sudah kita bahas pada course Go-Lang Web


### Kode: Log Middleware

```go
type Logger struct {
	Handler http.Handler
}

func (log *Logger) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Before handler execution")
	log.Handler.ServeHTTP(w, r)
	fmt.Println("After handler execution")
}
```


### Kode: Test Log Middleware

```go
func TestRouterLogMiddleware(t *testing.T) {
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Hello world!")
	})

	logger := &Logger{
		Handler: router,
	}

	request := httptest.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	logger.ServerHTTP(response, request)
	result := response.Result()
	assert.Equal(t, http.StatusOK, result.StatusCode)

	resultBody, err := ioutil.ReadAll(result.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Hello world!", string(resultBody))

	t.Log(string(resultBody))
}
```