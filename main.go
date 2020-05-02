package main

import (
	"fmt"      // Needed for printing
	"net/http" // Needed for http call
	"regexp"   // Needed for space trimming
	"strings"  // Needed for space trimming
	"strconv"  // Needed to convert string to float

	"github.com/PuerkitoBio/goquery" // Needed for HTML table parsing
)

const MAINE_SUPPLIER_URL = "https://www.maine.gov/meopa/electricity/electricity-supply"

func main() {
	fmt.Println("Maine electricity cost comparison tool.")

	electricProvider := "CMP"
	currentCost := 8.0
	lowestCost := currentCost

	// Grab all the rows
	rows := getHeadingRows(MAINE_SUPPLIER_URL)

	// Loop for all the rows and hope they don't change the format often.
	for _, row := range rows {
		if len(row) == 6 {
			fmt.Println("Provider:", row[0])
			fmt.Println("- CMP Rate:", row[1])
			fmt.Println("- Emra & BHE:", row[2])
			fmt.Println("- Fixed Term:", row[3])
			fmt.Println("- Early Termination Fee:", row[4])
			fmt.Println("- Contact number:", row[5])
			switch electricProvider {
			case "CMP":
				if strings.IndexAny(row[1], " (%)") > -1 {
					rowSplit := strings.Split(row[1], " ")
					row[1] = rowSplit[0]
				}
				rateCheck, err := strconv.ParseFloat(row[1], 64)
				handleErr(err)
				if lowestCost > rateCheck {
					lowestCost = rateCheck
				}
			case "Emra", "BHE":
				if strings.IndexAny(row[2], " (%)") > -1 {
					rowSplit := strings.Split(row[2], " ")
					row[2] = rowSplit[0]
				}
				rateCheck, err := strconv.ParseFloat(row[2], 64)
				handleErr(err)
				if lowestCost > rateCheck {
					lowestCost = rateCheck
				}
			}
		} else if len(row) == 3 {
			fmt.Println("- CMP Rate:", row[0])
			fmt.Println("- Emra & BHE:", row[1])
			fmt.Println("- Fixed Term:", row[2])
			fmt.Println("See above for early termination fee.")
			switch electricProvider {
			case "CMP":
				if strings.IndexAny(row[0], " (%)") > -1 {
					rowSplit := strings.Split(row[0], " ")
					row[0] = rowSplit[0]
				}
				rateCheck, err := strconv.ParseFloat(row[0], 64)
				handleErr(err)
				if lowestCost > rateCheck {
					lowestCost = rateCheck
				}
			case "Emra", "BHE":
				if strings.IndexAny(row[1], " (%)") > -1 {
					rowSplit := strings.Split(row[1], " ")
					row[1] = rowSplit[0]
				}
				rateCheck, err := strconv.ParseFloat(row[1], 64)
				handleErr(err)
				if lowestCost > rateCheck {
					lowestCost = rateCheck
				}
			}
		}
	}

	fmt.Println(currentCost, lowestCost)
	fmt.Println("If there's a provider offering cheaper electricity you can switch and save money.")
	fmt.Println("There's no change in equipment, lines, or anything. Just a billing change with your provider.")
	fmt.Println("For more details please visit: https://www.maine.gov/meopa/electricity/electricity-supply#CEPrates")
}

// Get the headings and rows from the website
func getHeadingRows(supplierURL string) ([][]string) {
	var row []string
	var rows [][]string

	spaceRemoval := regexp.MustCompile(`(\s\s)`)

	resp, err := http.Get(supplierURL)
	handleErr(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	handleErr(err)

	doc.Find("table").Each(func(index int, table *goquery.Selection) {
		table.Find("tr").Each(func(indextr int, tablerow *goquery.Selection) {
			tablerow.Find("td").Each(func(indexth int, cell *goquery.Selection) {
				row = append(row, spaceRemoval.ReplaceAllString(strings.TrimSpace(cell.Text()), " "))
			})
			rows = append(rows, row)
			row = nil
		})
	})

	return rows
}

// TODO: Actually implement error handling.
// For now, just panic and print the error.
func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}