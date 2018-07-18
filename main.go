package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"cities-api/src/cache"
	"cities-api/src/ds"
	"cities-api/src/parser"
	"cities-api/src/config"
	"os"
	"strconv"
	"github.com/Sirupsen/logrus"
	"strings"
	"cities-api/src/server"
	"os/exec"
	"time"
)

var (
	log *logrus.Entry = logrus.WithField("package", "server")

	Port = os.Getenv("PORT")

	CORSOriginsSTR = os.Getenv("CORS_ORIGINS")
	CORSOrigins    = []string{"http://localhost"}

	TimeoutSTR = os.Getenv("TIMEOUT")
	Timeout    = 5

	LocalesSTR = os.Getenv("LOCALES")
	Locales    = []string{"ru", "uk", "be", "en", "de"}

	CitiesFile = os.Getenv("CITIES_FILE")

	CountriesFile = os.Getenv("COUNTRIES_FILE")

	MinPopulationSTR = os.Getenv("MIN_POPULATION")
	MinPopulation    = 2000

	AlternateNamesFile = os.Getenv("ALTERNATE_NAMES_FILE")
)

func init() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if Port == "" {
		Port = "8080"
	}

	if CORSOriginsSTR != "" {
		CORSOrigins = []string{CORSOriginsSTR}
	}

	if TimeoutSTR != "" {
		integers, err := strconv.Atoi(TimeoutSTR)
		if err != nil {
			log.WithError(err).Warn("Going to use default Timeout")
		} else {
			Timeout = integers
		}
	}

	if LocalesSTR != "" {
		Locales = strings.Split(LocalesSTR, ",")
	}

	if CitiesFile == "" {
		CitiesFile = dir+"/data/cities1000.txt"
	}

	if CountriesFile == "" {
		CountriesFile = dir+"/data/countryInfo.txt"
	}

	if MinPopulationSTR != "" {
		integers, err := strconv.Atoi(MinPopulationSTR)
		if err != nil {
			log.WithError(err).Warn("Going to use default MinPopulation")
		} else {
			MinPopulation = integers
		}
	}

	if AlternateNamesFile == "" {
		AlternateNamesFile = dir+"/data/alternateNames.txt"
	}

}

func main() {

	out, err := exec.Command("/bin/sh", "getdumpfiles.sh").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("We managed to download all files and prepare for the DB ??  %s\n", out)

	time.Sleep(time.Second * 15)



	log.Info("* Booting cities service...")



	options := &config.Options{
		Port:               Port,
		Timeout:            Timeout,
		CORSOrigins:        CORSOrigins,
		Locales:            Locales,
		MinPopulation:      MinPopulation,
		CountriesFile:      CountriesFile,
		CitiesFile:         CitiesFile,
		AlternateNamesFile: AlternateNamesFile,
	}

	log.Info("* Loading configuration...", options)

	log.Info("* Connecting to the database...")
	db, err := bolt.Open("cities.db", 0600, nil)
	if err != nil {
		panic(fmt.Sprintf("[DB] Couldn't connect to the db: %v", err))
	}

	c := cache.New()
	parsingDone := make(chan bool, 1)

	if ds.GetAppStatus(db).IsIndexed() {
		log.Info("[PARSER] Skipping, already done")
		parsingDone <- true
	} else {
		go parser.Scan(
			db, parsingDone, Locales, MinPopulation,
			CountriesFile, CitiesFile,
			AlternateNamesFile,
		)
	}

	<-parsingDone
	log.Info("[CACHE] Warming up...")
	server.WarmUpSearchCache(db, c, Locales, 5)
	log.Info("[CACHE] Warming up done")

	log.Infof("* Listening on port %s\n\n", Port)
	log.Fatal(server.Server(db, options, c).ListenAndServe())
}
