package webscraping

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
)

func WebScraping(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/page/webScraping.page.html",
		"./ui/html/layout/base.layout.html",
		"./ui/html/partial/footer.partial.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере WebScrapping", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Внутренняя ошибка на сервере WebScrapping", 500)
		return
	}
}

// Обработчик!
func WebscrapingFormHandler(w http.ResponseWriter, r *http.Request) {
	//Проверяем, что метод запроса является POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Используйте r.FormValue для получения значений полей формы
	url := r.FormValue("url")                   //Адрес
	queryString := r.FormValue("queryString")   //Строка запроса
	textOrAtribut := r.FormValue("tip")         //Тип запроса текст или атрибуты
	valueAtribut := r.FormValue("valueAtribut") //Название атрибута
	resultTextArea := r.FormValue("result")     //Результат текстареа

	fmt.Println(url, queryString, textOrAtribut, valueAtribut, resultTextArea)

	parsing := webParsing(url, queryString, valueAtribut, textOrAtribut, resultTextArea) //функция парсинга страниц

	HttpResponse(w, parsing)

}

// HTTP-ответ
func HttpResponse(w http.ResponseWriter, parsing []string) {
	//Проверяем и записываем результат в тело HTTP-ответа
	if len(parsing) == 0 {
		fmt.Fprintln(w, "Если вы видете это сообщение, возможно имеет смысл проверить правильность 'Команды запроса', так как было возврашено 'Нулевой' результат!")
	}
	for i, v := range parsing {
		if v == "" {
			fmt.Fprintln(w, "Ты ошибся или пытаешься меня обмануть. Возможно ошибка в 'Команде запроса или в название атрибута'")
			return
		} else {
			fmt.Fprintf(w, "%d %s\n\n", i, v)
		}
	}
}

// Валидатор поля формы
func formFieldValidation(url, queryString, valueAtribut, textOrAtribut, resultTextArea string, resultPars []string) []string {
	//экранирования HTML-специальных символов
	url = template.HTMLEscapeString(url)
	queryString = template.HTMLEscapeString(queryString)
	valueAtribut = template.HTMLEscapeString(valueAtribut)
	textOrAtribut = template.HTMLEscapeString(textOrAtribut)
	resultTextArea = template.HTMLEscapeString(resultTextArea)
	//Проверочные Регулярки
	gapAndVoid := regexp.MustCompile(`\S`)                               //Регулярка на проверку пустоты и только пробела
	checkingFormFieldsForBrackets := regexp.MustCompile(`[<>(){}\[\]]+`) //Регулярка на проверку полей формы на скобки

	//Проверка на пустоту и только пробел поле Адреса
	matchUrl := gapAndVoid.MatchString(url)
	if !matchUrl {
		resultPars = append(resultPars, "Заполните поля адреса!")
		return resultPars
	}
	//Проверка на ввод запрещённых символов
	matchUrlCheckingFormFieldsForBrackets := checkingFormFieldsForBrackets.MatchString(url)
	if matchUrlCheckingFormFieldsForBrackets {
		resultPars = append(resultPars, "Ввод <>(){}[] запрещен")
		return resultPars
	}

	//Проверка на пустоту и только пробел поле Запроса
	matchQueryString := gapAndVoid.MatchString(queryString)
	if !matchQueryString {
		resultPars = append(resultPars, "Заполните поля запроса!")
		return resultPars
	}
	//Проверка на ввод запрещённых символов
	matchQueryStringCheckingFormFieldsForBrackets := checkingFormFieldsForBrackets.MatchString(queryString)
	if matchQueryStringCheckingFormFieldsForBrackets {
		resultPars = append(resultPars, "Ввод <>(){}[] запрещен")
		return resultPars
	}

	//Проверка на пустоту и только пробел поле название Атрибутов
	mathcValueAtribut := gapAndVoid.MatchString(valueAtribut)
	if textOrAtribut != "text" {
		if !mathcValueAtribut {
			resultPars = append(resultPars, "Заполните поля Название Атрибутов!")
			return resultPars
		}
	}

	return resultPars
}

// Парсим страницу
func webParsing(url, queryString, valueAtribut, textOrAtribut, resultTextArea string) []string {
	resultPars := []string{}                                                                                     //Сюда закидываем полученный результат
	createFile()                                                                                                 //Функция для создания файла                                                                                                //Функция создания файла
	validation := formFieldValidation(url, queryString, valueAtribut, textOrAtribut, resultTextArea, resultPars) //Функция проверки полей формы
	resultPars = append(resultPars, validation...)

	response, err := http.Get(url)

	//Проверяем что адрес рабочий
	if err != nil {
		resultPars = append(resultPars, "Нельзя спарсить это, возможно такой адрес не существует")
		return resultPars
	} else if response.StatusCode == 200 {

		document, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			fmt.Println(err)
		}

		if textOrAtribut == "text" { //выполниться если парсят текст
			document.Find(queryString).Each(func(i int, selection *goquery.Selection) {
				item := i
				title := selection.Text()
				resultPars = append(resultPars, title)
				saveOfFile(item, title) //функция для сохранения результата в файл
			})
		} else if textOrAtribut != "text" { //выполниться если парсят атрибуты
			document.Find(queryString).Each(func(i int, selection *goquery.Selection) {
				item := i
				title, _ := selection.Attr(valueAtribut)
				resultPars = append(resultPars, title)
				saveOfFile(item, title) //функция для сохранения результата в файл
			})
		} else {
			fmt.Println("Что-то пошло не так")
		}

	} else {
		resultPars = append(resultPars, "Нельзя спарсить это, проверьте Url")
		return resultPars
	}

	defer response.Body.Close()
	return resultPars
}

// Создать файл
func createFile() error {
	create, err := os.Create("parsingResult.txt")
	if err != nil {
		return err
	}
	defer create.Close()
	return nil
}

// Добавляем данные в файл
func saveOfFile(item int, title string) error {
	fmt.Println(item, title)

	create, err := os.OpenFile("parsingResult.txt", os.O_APPEND, 0777)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer create.Close()

	// Создаем кодировочный преобразователь для указанной кодировки (например, Windows-1251)
	encoder := charmap.Windows1251.NewEncoder()
	encoderTitle, err := encoder.String(title)
	if err != nil {
		fmt.Println(err)
	}

	_, err = create.WriteString(fmt.Sprintf("%s\n", encoderTitle))
	if err != nil {
		return err
	}
	return nil
}

func Download(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что запрос методом GET и путь равен "/webscraping/download"
	if r.Method == http.MethodGet && r.URL.Path == "/webscraping/download" {

		_, err := os.Stat("parsingResult.txt")
		if err != nil {
			fmt.Fprintf(w, "Зайгрузочный файл небыл создан. Сначала выполните запрос")
			return
		}

		// Устанавливаем заголовки для скачивания файла
		w.Header().Set("Content-Disposition", "attachment; filename=parsingResult.txt")
		w.Header().Set("Content-Type", "text/plain")

		// Открываем файл и отправляем его содержимое в ответ
		file, err := os.Open("parsingResult.txt")
		if err != nil {
			http.Error(w, "Не удалось прочитать файл", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		_, err = io.Copy(w, file)
		if err != nil {
			http.Error(w, "Не удалось отправить содержимое файла", http.StatusInternalServerError)
			return
		}

		// Закрываем файл
		err = file.Close()
		if err != nil {
			fmt.Println("Не удалось закрыть файл:", err)
		}

		// Удаляем файл после скачивания
		err = os.Remove("parsingResult.txt")
		if err != nil {
			fmt.Println("Не удалось удалить файл(:", err)
		}
	}
}
