package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Data struct {
	Input  string
	Banner string
	Output string
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		// http.NotFound(writer, request)
		http.Error(writer, "Page Not Found", http.StatusNotFound)
		return
	}

	hometemp, err := template.ParseFiles("template/index.html")
	if err != nil {
		fmt.Fprint(writer, "home error")
	}

	hometemp.Execute(writer, Data{})
	fmt.Fprint(writer, "homepage")

}

func AsciiHandler(writer http.ResponseWriter, request *http.Request) {

	tempone, _ := template.ParseFiles("template/index.html")

	if request.Method != "POST" {
		http.Error(writer, "this method is not allowed", http.StatusNotFound)
		return
	}

	if request.Method == "POST" {
		temp, err := template.ParseFiles("template/index.html")
		if err != nil {
			fmt.Fprint(writer, "template error")
		}

		request.ParseForm()

		Word := request.FormValue("Input")
		Style := request.FormValue("banner")

		newascii := GenerateAscii(Word, Style)

		NewData := Data{
			Input:  Word,
			Banner: Style,
			Output: newascii,
		}

		temp.Execute(writer, NewData)
		return
	}

	// Placeholder := Data{
	// 	Input:  "",
	// 	Banner: "standard",
	// 	Output: "",
	// }

	tempone.Execute(writer, Data{})
	fmt.Fprint(writer, "asciipage")

	// http.ServeFile(writer, request, "index.html")
}

func GenerateAscii(text string, banner string) string {
	// opens file and writes error handling pattern
	file, err := os.Open(banner)
	if err != nil {
		fmt.Println("error opening file one", err)
		// return
	}
	defer file.Close() // closes the file

	scanner := bufio.NewScanner(file) // reads the file

	var content []string
	for scanner.Scan() {
		lines := scanner.Text()
		content = append(content, lines)
	}

	// error handling pattern for file reading
	err = scanner.Err()
	if err != nil {
		fmt.Println("error reading", err)
		// return
	}

	// splits "text" by newline

	var result strings.Builder

	newtext := strings.Split(text, "\r\n")

	for ind, val := range newtext {
		if val == "" {
			result.WriteString("\n")
			continue
		}

		// range loop to carry the split adjustment

		// outer loop prints the row
		for row := 0; row < 8; row++ {
			for i := 0; i < len(newtext[ind]); i++ { // inner loop prints the characters
				asciival := newtext[ind][i]
				start := int(asciival-32)*9 + 1 // logic to generate the ascii start
				char := content[start+row]
				result.WriteString(char)
			}
			result.WriteString("\n")
		}
	}
	return result.String()
}

func main() {

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/ascii", AsciiHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("server running on port: http://localhost:5000")
	http.ListenAndServe(":5000", nil)

}
