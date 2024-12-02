package main

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/psmccarty/currency-converter/client"
)

var (
	base        string
	newCurrency string
	value       string
)

func main() {

	c := client.NewConverter(os.Getenv("APP_ID"))
	currenciesMap, err := c.GetCurrencies()
	if err != nil {
		log.Fatal(err)
	}

	currenciesSlice := slices.Collect(maps.Keys(currenciesMap))
	slices.Sort(currenciesSlice)

	options := make([]huh.Option[string], 0, len(currenciesSlice))
	for _, currency := range currenciesSlice {
		options = append(options, huh.NewOption(fmt.Sprintf("%s: %s", currency, currenciesMap[currency]), currency))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What is your base currency").
				Options(
					options...,
				).
				Height(8).
				Value(&base),

			huh.NewSelect[string]().
				Title("What do you want to convert into").
				Options(
					options...,
				).
				Height(8).
				Value(&newCurrency),

			huh.NewInput().
				Title("How much to convert").
				Value(&value).
				Validate(func(str string) error {
					val, err := strconv.ParseFloat(str, 64)
					if err != nil {
						return err
					}
					if val < 0 {
						return errors.New("must be non negative number")
					}
					return nil
				}),
		),
	)
	form.Run()

	amount, err := c.Convert(base, newCurrency, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%s %s converts to %f %s\n\n", value, base, amount, newCurrency)
}
