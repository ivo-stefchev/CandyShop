package store

import "io"
import "log"
import "net/http"
import "regexp"
import "bytes"
import "strconv"
import "strings"
import "fmt"
import "sort"
import "errors"
import "encoding/json"

const delimiter = "::DELIMITER::"

func Do() {
	resp, err := http.Get("https://candystore.zimpler.net/")
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	tds, err := getCustomersTds(body)
	if err != nil {
		log.Fatal(err)
	}
	favourite := getTotalAndFavourite(tds)
	marshaled, _ := json.MarshalIndent(favourite, "", "  ")

	fmt.Printf("%v\n", string(marshaled))
}

type TotalAndFavorite struct {
	// name of person
	Name             string `json:"name"`
	// name of favourite snack
	Candy            string `json:"favouriteSnack"`
	// quantity of consumed favourite snack
	TopCandyQuantity int    `json:"-"`
	// total consumed snacks
	Total            int    `json:"totalSnacks"`
}

var reCustomersTable = regexp.MustCompile(`<table id="top.customers" class="top.customers details">.*<tbody>(.*)</tbody>\s*</table>`)
var reTds = regexp.MustCompile(`<td>([\w\s\p{L}]+)</td>`)

// extract <td>s from customers table
func getCustomersTds(inputHtml []byte) ([]string, error) {
	// remove all newlines as regexp does not work in multiline
	noNewLines := bytes.ReplaceAll(inputHtml, []byte("\n"), []byte(""))

	tableCustomers := reCustomersTable.FindSubmatch(noNewLines)
	if len(tableCustomers) < 1 {
		return nil, errors.New("Did not find top customers table")
	}

	tds := reTds.FindAllSubmatch(tableCustomers[1], -1)

	result := make([]string, 0, len(tds))
	for _, v := range tds {
		result = append(result, string(v[1]))
	}

	return result, nil
}

// len(input) should be multiple of 3
func getTotalAndFavourite(input []string) []TotalAndFavorite {
	if len(input) % 3 != 0 {
		log.Fatal("Number of input should be multiple of 3.")
	}

	uniquePersonAndSnack := make(map[string]int)
	for i := 0; i < len(input) / 3; i++ {
		y := i * 3
		numberStr := input[y + 2]
		quantity, err := strconv.Atoi(numberStr)
		if err != nil {
			log.Fatalf("Could not convert to number from string '%v'.", numberStr)
		}
		name := input[y]
		candy := input[y + 1]
		key := name + delimiter + candy
		uniquePersonAndSnack[key] += quantity
	}

	personToTotalAndFavorite := make(map[string]TotalAndFavorite)

	for key, quantity := range uniquePersonAndSnack {
		split := strings.Split(key, delimiter)
		name := split[0]
		candy := split[1]
		current := personToTotalAndFavorite[name]
		current.Name = name
		current.Total += quantity
		if quantity > current.TopCandyQuantity {
			current.Candy = candy
			current.TopCandyQuantity = quantity
		}
		personToTotalAndFavorite[name] = current
	}

	totalAndFavorite := make([]TotalAndFavorite, 0, len(personToTotalAndFavorite))
	for _, p := range personToTotalAndFavorite {
		totalAndFavorite = append(totalAndFavorite, p)
	}

	sort.Slice(totalAndFavorite, func(i, j int) bool {
		return totalAndFavorite[i].Total > totalAndFavorite[j].Total
	})

	return totalAndFavorite
}
