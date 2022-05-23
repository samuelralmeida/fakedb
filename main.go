package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"Fakedb/credentials"
	"Fakedb/fakedata/dbs/coisadb"
	"Fakedb/fakedata/dbs/tremdb"
	"Fakedb/script"
)

func main() {
	var source string
	var target string

	flag.StringVar(&source, "source", "all", "Specify which database. Options: all, trem, coisa. Default: all")
	flag.StringVar(&target, "target", "all", "Specify which dump destination. Options: all, dev, staging. Default: all")
	flag.Parse()

	sources := map[string]bool{"all": true, "trem": true, "coisa": true}
	targets := map[string]bool{"all": true, "dev": true, "staging": true}

	if _, ok := sources[source]; !ok {
		log.Fatalf("%s is a invalid source. Options: all, farme, indicame.", source)
	}

	if _, ok := targets[target]; !ok {
		log.Fatalf("%s is a invalid target. Options: all, dev, staging.", target)
	}

	credentialTargets := []credentials.ICredentials{}
	if target == "all" || target == "dev" {
		credentialTargets = append(credentialTargets, credentials.NewCredentials(credentials.Dev))
	}
	if target == "all" || target == "staging" {
		credentialTargets = append(credentialTargets, credentials.NewCredentials(credentials.Staging))
	}

	prodCredentials := credentials.NewCredentials(credentials.Prod)
	intermediateCredentials := credentials.NewCredentials(credentials.Intermediate)

	var wg sync.WaitGroup

	if source == "all" || source == "trem" {

		t := []credentials.ConnectionParams{}
		for _, ct := range credentialTargets {
			t = append(t, ct.GetCredentialsBySource(credentials.TremSource))
		}

		scriptParams := script.DatabaseParams{
			Origin:       prodCredentials.GetCredentialsBySource(credentials.TremSource),
			Intermediate: intermediateCredentials.GetCredentialsBySource(credentials.TremSource),
			Targets:      t,
			FakeData:     tremdb.NewTremDb(),
		}

		wg.Add(1)
		go func(dbParams script.DatabaseParams) {
			log.Println("START: trem dump")
			defer wg.Done()

			err := script.FakeDatabase(dbParams)
			if err != nil {
				err = fmt.Errorf("ERROR: não foi possivel fazer o dump do banco trem: %w", err)
				log.Println(err)
			}
		}(scriptParams)
	}

	if source == "all" || source == "coisa" {

		t := []credentials.ConnectionParams{}
		for _, ct := range credentialTargets {
			t = append(t, ct.GetCredentialsBySource(credentials.CoisaSource))
		}

		scriptParams := script.DatabaseParams{
			Origin:       prodCredentials.GetCredentialsBySource(credentials.CoisaSource),
			Intermediate: intermediateCredentials.GetCredentialsBySource(credentials.CoisaSource),
			Targets:      t,
			FakeData:     coisadb.NewCoisaDb(),
		}

		wg.Add(1)
		go func(dbParams script.DatabaseParams) {
			log.Println("START: coisa dump")
			defer wg.Done()

			err := script.FakeDatabase(dbParams)
			if err != nil {
				err = fmt.Errorf("ERROR: não foi possivel fazer o dump do banco coisa: %w", err)
				log.Println(err)
			}
		}(scriptParams)
	}

	wg.Wait()

}
