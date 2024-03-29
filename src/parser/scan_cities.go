package parser

import (
	"bufio"
	"github.com/boltdb/bolt"
	"cities-api/src/ds"
	"os"
	"strconv"
	"strings"
)

func scanCities(
	db *bolt.DB, filename string, minPopulation int,
) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	citiesCount := 0

	err = db.Batch(func(tx *bolt.Tx) error {
		citiesBucket := tx.Bucket(ds.CitiesBucketName)
		cityNamesBucket := tx.Bucket(ds.CityNamesBucketName)

		for scanner.Scan() {
			cityData := strings.Split(scanner.Text(), "\t")
			cityBytes, err := prepareCityBytes(cityData)
			if err != nil {
				return err
			}

			population, _ := strconv.ParseInt(cityData[14], 0, 64)
			if population > int64(minPopulation) {
				citiesBucket.Put([]byte(cityData[0]), cityBytes)

				addCityToIndex(
					cityNamesBucket, cityData[0], cityData[1], "", uint32(population),
				)

				citiesCount++
			}
		}

		return err
	})

	return citiesCount, err
}
