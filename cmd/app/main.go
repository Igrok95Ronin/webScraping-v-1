package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"webScraping/internal/calculator"
	"webScraping/internal/home"
	"webScraping/internal/passwordGenerator"
	"webScraping/internal/todolist"
	"webScraping/internal/webscraping"
)

func main() {
	addr := flag.String("addr", ":4000", "Сетевой адрес HTTP") // Создаем новый флаг командной строки, значение по умолчанию: ":4000"
	flag.Parse()                                               // Мы вызываем функцию flag.Parse() для извлечения флага из командной строки

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)                  //logs
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile) //logs

	mux := http.NewServeMux()
	mux.HandleFunc("/", home.Home)
	//Calculator
	mux.HandleFunc("/calculator", calculator.Calculator)
	mux.HandleFunc("/formHandler", calculator.FormHandler)
	mux.HandleFunc("/deleteEntry", calculator.DeleteEntry)
	//ToDoList
	mux.HandleFunc("/todolist", todolist.ToDoList)
	mux.HandleFunc("/formHandlerToDoList", todolist.FormHandlerToDoList)
	mux.HandleFunc("/deleteEntryToDoList", todolist.DeleteEntryToDoList)
	mux.HandleFunc("/editPost", todolist.EditPost)

	//Password Generator
	mux.HandleFunc("/passwordgenerator", passwordGenerator.PasswordGenerator)

	//Web Scraping
	mux.HandleFunc("/webscraping", webscraping.WebScraping)
	mux.HandleFunc("/webscrapingFormHandler", webscraping.WebscrapingFormHandler)
	mux.HandleFunc("/webscraping/download", webscraping.Download)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{ //logs
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Запуск сервера на %s", *addr) // Значение, возвращаемое функцией flag.String(), является указателем на значение go run ./cmd/web -addr=":9999"
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}
