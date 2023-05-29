package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"github.com/go-chi/chi"
	// "github.com/go-chi/chi/middleware"
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

type Inbox struct {
	Took int `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards Shards `json:"_shards"`
	Hits Hits `json:"hits"`	
}

type Shards struct {
	Total int `json:"total"`
	Successful int `json:"successful"`
	Skipped int `json:"skipped"`
	Failed int `json:"failed"`
}

type Hits struct {
	Total total `json:"total"`
	MaxScore float64 `json:"max_score"`
	HitsSon []HitsSon `json:"hits"`
}

type total struct {
	Value int `json:"value"`	
}

type HitsSon struct {
	Index string `json:"_index"`	
	TypeInfo string `json:"_type"`	
	Id string `json:"_id"`	
	Score float64 `json:"_score"`	
	Timestamp string `json:"@timestamp"`
	Source Source `json:"_source"`
}

type Source struct {
	Timestamp string `json:"@timestamp"`	
	ContentTransferEncoding string `json:"Content-Transfer-Encoding"`	
	ContentType string `json:"Content-Type"`	
	Date string `json:"Date"`	
	From string `json:"From"`	
	MessageID string `json:"Message-ID"`	
	MimeVersion string `json:"Mime-Version"`	
	Subject string `json:"Subject"`	
	To string `json:"To"`
	XFileName string `json:"X-FileName"`
	XFolder string `json:"X-Folder"`
	XFrom string `json:"X-From"`
	XOrigin string `json:"X-Origin"`
	XTo string `json:"X-To"`
	Xbcc string `json:"X-bcc"`
	Xcc string `json:"X-cc"`
	Body string `json:"body"`
}

type Request struct {
	Field string `json:"field"`
	Term string `json:"term"`
	Page string `json:"page"`
}

func injeccion(w http.ResponseWriter, r *http.Request) {

	root := "maildir"
	err := filepath.Walk(root, func (path string, info os.FileInfo, err error) error {
		// Ignorar errores que puedan ocurrir al acceder a la carpeta/archivo actual
		if err != nil {
			return err
		}
	
		// Si la ruta actual es una carpeta, imprimir el nombre de la carpeta y continuar
		if info.IsDir() {
			return nil
		}		
	
		// Si la ruta actual es un archivo, imprimir el nombre del archivo y continuar
		// fmt.Println("Archivo:", path)*
		infoPath := strings.Split(path, "\\")
		index := infoPath[2]
		person := infoPath[1]
		indexSave := chi.URLParam(r, "index")
		personSave := chi.URLParam(r, "person")


		fmt.Println("Directorio:", path)

		if index == indexSave && (personSave == "all" || personSave == person) {

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
		
			// Crear un escáner para leer el archivo línea por línea
			scanner := bufio.NewScanner(file)
		
			data := make(map[string]string)
			dataOld := [1]string{}
			// Iterar a través de cada línea del archivo
			bodyF := false
			for scanner.Scan() {
				if bodyF {
					data["body"] = data["body"] + scanner.Text()
				} else {
					if strings.Contains(scanner.Text(), ":") {
						infoCampo := strings.Split(scanner.Text(), ":")
						data[infoCampo[0]] = infoCampo[1]
						dataOld[0] = infoCampo[0]
						if infoCampo[0] == "X-FileName" {
							bodyF = true
						}					
					} else {
						data[dataOld[0]] = data[dataOld[0]] + scanner.Text()
					}
				}
			}			
			// Comprobar si ha ocurrido algún error durante la lectura del archivo
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				return nil
			}
		
			// Convertimos los registros a JSON
			jsonData, err := json.Marshal(data)
			if err != nil {
				return err
			}
			
			req, err := http.NewRequest("POST", "http://localhost:4080/api/Index_"+index+"/_doc", strings.NewReader(string(jsonData)))
			if err != nil {
				log.Fatal(err)
			}
			req.SetBasicAuth("admin", "Complexpass#123")
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
		
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			log.Println(resp.StatusCode)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(body))
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}	   
	
	w.Write([]byte("Insert finish"))
}

func search(w http.ResponseWriter, r *http.Request) {

	var data Request
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	index := chi.URLParam(r, "index")
	field := data.Field
	term := data.Term
	page := data.Page

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		fmt.Println("Error al convertir la cadena a entero:", err)
		return
	}

	max := pageNum * 5

	field = strings.TrimSpace(field)
	term = strings.TrimSpace(term)

	InfoInbox := get_search(index, field, term, strconv.Itoa(max))

	Data := []Source{}

	for _,source := range InfoInbox.Hits.HitsSon {
		Data = append(Data, source.Source)
	}
	
	// Convertir el array en JSON
	jsonBytes, err := json.Marshal(Data)
	if err != nil {
		fmt.Println("Error al convertir a JSON:", err)
		return
	}

	// Convertir los bytes en string para imprimirlos
	w.Write(jsonBytes)
}

func get_search(Index, Field, Term, Max string) Inbox {
	index := Index
	field := Field
	term := Term
	max := Max

	query := `{
		"_source": [],
		"from": 0,
        "max_results": `+max+`,
        "query": 
		{
			"field": "`+field+`",
            "term": "`+term+`"
        },
        "search_type": "match"        
    }`
    req, err := http.NewRequest("POST", "http://localhost:4080/api/Index_"+index+"/_search", strings.NewReader(query))
    if err != nil {
        log.Fatal(err)
    }
    req.SetBasicAuth("admin", "Complexpass#123")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

	// CREAMOS  EL JSON DEL RESULTADO

	var arreglo Inbox

	err = json.Unmarshal(body, &arreglo)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return arreglo
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	r := chi.NewRouter()

	// r.Use(middleware.Logger) // middleware para registrar las solicitudes HTTP
	// r.Use(middleware.Recoverer) // middleware para recuperarse de los errores de la aplicación
	r.Post("/insert-info/{person}/{index}", injeccion)
	r.Post("/search/{index}", search)
	r.Get("/", index)
	http.ListenAndServe(":3000", r)	
}
