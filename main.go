package main

import (
	"flag"     // Needed for CLI parsing
	"fmt"      // Needed for printing
	"net/http" // Needed for http call
	"os"       // Needed to do len(os.Args)
	"regexp"   // Needed for space trimming
	"strconv"  // Needed to convert string to float
	"strings"  // Needed for space trimming

	"github.com/PuerkitoBio/goquery" // Needed for HTML table parsing
)

const MAINE_SUPPLIER_URL = "https://www.maine.gov/meopa/electricity/electricity-supply"

func main() {
	var currentCost float64
	var electricProvider string

	// Check to see if the user supplied args, or wants to take input.
	if len(os.Args) > 1 {
		flag.Float64Var(&currentCost, "current", 10.0, "Current cost of your electricity in ¢/kWh")
		flag.StringVar(&electricProvider, "provider", "CMP", "Current electric provider, defaults to CMP")
		flag.Parse()
	} else {
		fmt.Printf("Enter your current price for electricity (¢/kWh): ")
		_, err := fmt.Scanf("%f\n", &currentCost)
		handleErr(err)

		fmt.Printf("Enter your current provider. (CMP or EMERA): ")
		_, err = fmt.Scanf("%s\n", &electricProvider)
		electricProvider = strings.ToUpper(electricProvider)
		handleErr(err)
	}

	lowestCost := currentCost

	fmt.Println("Maine electricity cost comparison tool.")

	// Grab all the rows
	rows := getHeadingRows(MAINE_SUPPLIER_URL)

	// Loop for all the rows and hope they don't change the format often.
	for _, row := range rows {
		if len(row) == 6 {
			fmt.Println("Provider:", row[0])
			fmt.Println("- CMP Rate:", row[1])
			fmt.Println("- Emera & BHE:", row[2])
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
			case "EMERA":
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
			fmt.Println("- Emera & BHE:", row[1])
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
			case "EMERA":
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

	fmt.Println("Current Cost:", currentCost)
	fmt.Println("Lowest Cost Found:", lowestCost)
	fmt.Println("If there's a provider offering cheaper electricity you can switch and save money.")
	fmt.Println("There's no change in equipment, lines, or anything. Just a billing change with your provider.")
	fmt.Println("For more details please visit: https://www.maine.gov/meopa/electricity/electricity-supply#CEPrates")
	if len(flag.Args()) < 2 {
		fmt.Println("Press enter to exit.")
		fmt.Scanf("%v\n", &electricProvider)
	}
}

// Get the headings and rows from the website
func getHeadingRows(supplierURL string) [][]string {
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
