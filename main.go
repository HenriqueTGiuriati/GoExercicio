package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

//City is a struct
type City struct {
	name   string
	number []int
}

//CityAverage is a struct
type CityAverage struct {
	Cidade string  `json:"cidade"`
	Idade  float64 `json:"idade"`
}

//Averages is a struct
type Averages struct {
	Medias []CityAverage `json:"medias"`
}

var cities []City

func main() {
	fmt.Println("Bem vindo ao Exercício 1.")

	arquivo, err := os.Open("input.csv")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	reader := bufio.NewReader(arquivo)

	for {
		line, err := reader.ReadString('\n')

		slice := strings.Split(line, ",")

		if err == io.EOF {
			break
		}

		if checkIfCityAlreadyInSlice(convertName(slice[1])) {
			for i := range cities {
				if cities[i].name == convertName(slice[1]) {
					cities[i].number = append(cities[i].number, convertStringToInt(slice[2]))
				}
			}

		} else {
			var number []int

			n := convertStringToInt(slice[2])

			city := City{name: convertName(slice[1]), number: append(number, n)}
			cities = append(cities, city)
		}

	}

	average := Averages{}
	for i := range cities {
		var total int
		count := 0
		for j := range cities[i].number {
			total = total + cities[i].number[j]
			count++
		}

		av := float64(total) / float64(count)
		average.Medias = append(average.Medias, CityAverage{Cidade: cities[i].name, Idade: math.Round(av*100) / 100})
	}

	buffer, _ := json.MarshalIndent(average, "", "  ")

	jsn := string(buffer)
	fmt.Println(jsn)

	resp, err := http.Post("https://zeit-endpoint.brmaeji.now.sh/api/avg", "application/json", bytes.NewBuffer(buffer))

	if err != nil {
		fmt.Println("Ocorreu um erro ao enviar a request:", err)
	} else {
		fmt.Println("Resposta da request:", resp)
	}
}

func convertStringToInt(number string) int {
	number = strings.Trim(number, "\t \n")
	n, err := strconv.Atoi(number)

	if err != nil {
		fmt.Println(err)
	}

	return n
}

func checkIfCityAlreadyInSlice(name string) bool {
	if cities != nil {
		for _, city := range cities {
			if city.name == strings.ToLower(strings.Trim(name, "\t \n\x00")) {
				return true
			}
		}
	}
	return false
}

func convertName(name string) string {
	b := make([]byte, len(name))

	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}

	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	_, _, err := t.Transform(b, []byte(name), true)
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	return strings.Trim(strings.ToLower(string(b)), "\t \n\x00")
}
