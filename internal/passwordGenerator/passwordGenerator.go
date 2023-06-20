package passwordGenerator

import (
	"log"
	"net/http"
	"text/template"
)

func PasswordGenerator(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/page/passwordGenerator.page.html",
		"./ui/html/layout/base.layout.html",
		"./ui/html/partial/footer.partial.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере2", 500)
		return
	}
}
