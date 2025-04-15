package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	/*
		!!! !!! !!!
		обратите внимание - в задании обязательно нужен отчет
		делать его лучше в самом начале, когда вы видите уже узкие места, но еще не оптимизировалм их
		так же обратите внимание на команду в параметром -http
		перечитайте еще раз задание
		!!! !!! !!!
	*/
	// открытие файла
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// полное чтение.
	fileContents, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Создаёт регулярное выражение для поиска символа "@"
	r := regexp.MustCompile("@")
	seenBrowsers := make(map[string]struct{})
	uniqueBrowsers := 0
	foundUsers := ""

	// Создание буфера для чтения
	buf := make([]byte, 4096)
	var data []byte
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		data = append(data, buf[:n]...)
		if err == io.EOF {
			break
		}
	}

	// парсинг текста
	lines := strings.Split(string(fileContents), "\n")

	for i, line := range lines {
		var user map[string]interface{}
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			continue
		}

		browsers, ok := user["browsers"].([]interface{})
		if !ok {
			continue
		}

		isAndroid := false
		isMSIE := false

		for _, browserRaw := range browsers {
			browser, ok := browserRaw.(string)
			if !ok {
				continue
			}

			// Проверка на Android и MSIE
			if strings.Contains(browser, "Android") {
				isAndroid = true
			}
			if strings.Contains(browser, "MSIE") {
				isMSIE = true
			}

			// Используем карту для хранения уникальных браузеров
			if isAndroid || isMSIE {
				if _, exists := seenBrowsers[browser]; !exists {
					seenBrowsers[browser] = struct{}{}
					uniqueBrowsers++
				}
			}
		}
		// Если пользователь использует и Android, и MSIE, добавляем его в список
		if isAndroid && isMSIE {
			email := r.ReplaceAllString(user["email"].(string), " [at] ")
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user["name"], email)
		}
	}

	// fmt.Fprintln(out, "found users:\n"+foundUsers)
	// fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}
