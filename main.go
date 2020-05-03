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

// Offer from the electric supplier
type ElectricOffer struct {
	CMPRate   float64
	EMERARate float64
	FixedTerm string
}

// Electric supplier
type ElectricSupplier struct {
	Company       string
	Offers        []ElectricOffer
	EarlyTermFee  string
	ContactNumber string
}

func main() {
	var currentCost float64
	var electricProvider string
	var lowestProvider string

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

	electricProvider = strings.ToUpper(electricProvider)

	lowestCost := currentCost

	fmt.Println("Maine electricity cost comparison tool.")

	// Grab all the rows
	rows := getHeadingRows(MAINE_SUPPLIER_URL)

	supplierList := buildSupplierList(rows)

	for _, supplier := range supplierList {
		fmt.Println("Supplier:", supplier.Company)
		fmt.Println("Offerings:")
			fmt.Println("CMP Rate\tEmera Rate\tTerm")
		for _, offer := range supplier.Offers {
			fmt.Printf("¢%.02f\t\t¢%.02f\t\t%s\r\n", offer.CMPRate, offer.EMERARate, offer.FixedTerm)

			switch electricProvider {
			case "CMP":
				if lowestCost > offer.CMPRate {
					lowestCost = offer.CMPRate
					lowestProvider = supplier.Company
				}
			case "EMERA":
				if lowestCost > offer.EMERARate {
					lowestCost = offer.EMERARate
					lowestProvider = supplier.Company
				}
			}
		}
		fmt.Println("Early Termination Fee:", supplier.EarlyTermFee)
		fmt.Println("Contact Number:", supplier.ContactNumber)
		fmt.Println("-----")
	}
	fmt.Println("Current Cost:", currentCost)
	fmt.Printf("Lowest Cost Found: %.02f from %s\r\n", lowestCost, lowestProvider)
	fmt.Println("If there's a provider offering cheaper electricity you can switch and save money.")
	fmt.Println("There's no change in equipment, lines, or anything. Just a billing change with your provider.")
	fmt.Println("For more details please visit: https://www.maine.gov/meopa/electricity/electricity-supply#CEPrates")
	if len(os.Args) < 2 {
		fmt.Println("Press enter to exit.")
		fmt.Scanf("%v\n", &electricProvider)
	}
}

// Build out a slice of electric suppliers from the HTML table
func buildSupplierList(rows [][]string) []ElectricSupplier {
	var supplierList []ElectricSupplier
	var supplier ElectricSupplier
	var offering ElectricOffer
	var lastSupplier = supplier

	// Loop for all the rows and hope they don't change the format often.
	for _, row := range rows {
		if len(row) == 6 {
			if lastSupplier.Company != "" {
				supplierList = append(supplierList, lastSupplier)
			}

			var supplier ElectricSupplier

			supplier.Company = row[0]
			supplier.EarlyTermFee = row[4]
			supplier.ContactNumber = row[5]

			for x := 1; x <= 2; x++ {
				if strings.IndexAny(row[1], " (%)") > -1 {
					rowSplit := strings.Split(row[1], " ")
					row[1] = rowSplit[0]
				}
			}

			offering.CMPRate, _ = strconv.ParseFloat(row[1], 64)
			offering.EMERARate, _ = strconv.ParseFloat(row[2], 64)
			offering.FixedTerm = row[3]

			supplier.Offers = append(supplier.Offers, offering)
			lastSupplier = supplier
		} else if len(row) == 3 {
			var offering ElectricOffer

			for x := 0; x <= 1; x++ {
				if strings.IndexAny(row[x], " (%)") > -1 {
					rowSplit := strings.Split(row[x], " ")
					row[x] = rowSplit[0]
				}
			}

			offering.CMPRate, _ = strconv.ParseFloat(row[0], 64)
			offering.EMERARate, _ = strconv.ParseFloat(row[1], 64)
			offering.FixedTerm = row[2]

			lastSupplier.Offers = append(lastSupplier.Offers, offering)

		}
	}

	return supplierList
}

// Get the headings and rows from the website
// Thanks to Salmoni for some of this logic.
// https://gist.github.com/salmoni/27aee5bb0d26536391aabe7f13a72494
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
