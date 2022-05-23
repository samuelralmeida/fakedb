package script

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"Fakedb/credentials"
	"Fakedb/fakedata"
	"Fakedb/mysqldb"
)

type DatabaseParams struct {
	Origin       credentials.ConnectionParams
	Intermediate credentials.ConnectionParams
	Targets      []credentials.ConnectionParams
	FakeData     fakedata.IFakeData
}

func FakeDatabase(dp DatabaseParams) error {
	var (
		err                 error
		sqlOriginFile       = fmt.Sprintf("%s_origin.sql", dp.Origin.Database)
		sqlIntermediateFile = fmt.Sprintf("%s_intermediate.sql", dp.Intermediate.Database)
	)

	log.Printf("dump database - %s", dp.Origin.Database)
	err = mysqldb.DumpDatabase(sqlOriginFile, dp.Origin)
	if err != nil {
		return err
	}

	log.Printf("drop database if exists - %s", dp.Intermediate.Database)
	err = mysqldb.DropDatabase(dp.Intermediate)
	if err != nil {
		return err
	}

	log.Printf("create database - %s", dp.Intermediate.Database)
	err = mysqldb.CreateDatabase(dp.Intermediate)
	if err != nil {
		return err
	}

	log.Printf("populate database - %s", dp.Intermediate.Database)
	err = mysqldb.RestoreData(sqlOriginFile, dp.Intermediate)
	if err != nil {
		return err
	}

	log.Printf("fake data - %s", dp.Intermediate.Database)
	err = FakeData(dp.FakeData, dp.Intermediate)
	if err != nil {
		return err
	}

	log.Printf("dump database - %s", dp.Intermediate.Database)
	err = mysqldb.DumpDatabase(sqlIntermediateFile, dp.Intermediate)
	if err != nil {
		return err
	}

	err = restoreTargets(sqlIntermediateFile, dp.Targets)
	if err != nil {
		return err
	}

	log.Printf("drop database - %s", dp.Intermediate.Database)
	err = mysqldb.DropDatabase(dp.Intermediate)
	if err != nil {
		return err
	}

	log.Println("remove temp files")
	err = removeFiles([]string{sqlOriginFile, sqlIntermediateFile})
	if err != nil {
		return err
	}

	return nil
}

func restoreTargets(filename string, targets []credentials.ConnectionParams) error {
	var (
		wg   sync.WaitGroup
		errs = []string{}
	)

	for _, target := range targets {
		t := target

		wg.Add(1)
		go func(credential credentials.ConnectionParams) {
			defer wg.Done()

			var err error

			log.Printf("drop database if exists - %s", credential.Database)
			err = mysqldb.DropDatabase(credential)
			if err != nil {
				err = fmt.Errorf("target - %s: %w", credential.Database, err)
				errs = append(errs, err.Error())
				return
			}

			log.Printf("create database - %s", credential.Database)
			err = mysqldb.CreateDatabase(credential)
			if err != nil {
				err = fmt.Errorf("target - %s: %w", credential.Database, err)
				errs = append(errs, err.Error())
				return
			}

			log.Printf("restore database - %s", credential.Database)
			err = mysqldb.RestoreData(filename, credential)
			if err != nil {
				err = fmt.Errorf("target - %s: %w", credential.Database, err)
				errs = append(errs, err.Error())
				return
			}

			log.Printf("DONE: restore database - %s", credential.Database)

		}(t)
	}

	wg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "&&"))
	}
	return nil
}

func removeFiles(filenames []string) error {
	cmd := exec.Command("rm", filenames...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("removeFiles: error to run command: %w", err)
	}

	return nil
}
