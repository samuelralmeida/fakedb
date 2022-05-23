package mysqldb

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"Fakedb/credentials"
)

func DumpDatabase(filename string, args credentials.ConnectionParams) error {
	outfile, err := os.Create(fmt.Sprintf("./%s", filename))
	if err != nil {
		return fmt.Errorf("mysqldump: error to open file: %w", err)
	}
	defer outfile.Close()

	mysqldumpArgs := []string{
		fmt.Sprintf("--host=%s", args.Host),
		fmt.Sprintf("--port=%s", args.Port),
		"--default-character-set=utf8",
		fmt.Sprintf("--user=%s", args.User),
		fmt.Sprintf("--password=%s", args.Password),
		args.Database,
		"--single-transaction",
	}

	cmd := exec.Command("mysqldump", mysqldumpArgs...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("mysqldump: error to start command: %w", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("mysqldump: error to wait command: %w", err)
	}

	return nil
}

func DropDatabase(args credentials.ConnectionParams) error {

	dropArgs := []string{
		fmt.Sprintf("--host=%s", args.Host),
		fmt.Sprintf("--port=%s", args.Port),
		fmt.Sprintf("--user=%s", args.User),
		fmt.Sprintf("--password=%s", args.Password),
		"-e",
		fmt.Sprintf("DROP DATABASE IF EXISTS %s", args.Database),
	}

	cmd := exec.Command("mysql", dropArgs...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("dropDatabase: error to run command: %w", err)
	}

	return nil
}

func CreateDatabase(args credentials.ConnectionParams) error {
	dropArgs := []string{
		fmt.Sprintf("--host=%s", args.Host),
		fmt.Sprintf("--port=%s", args.Port),
		fmt.Sprintf("--user=%s", args.User),
		fmt.Sprintf("--password=%s", args.Password),
		"-e",
		fmt.Sprintf("CREATE DATABASE %s", args.Database),
	}

	cmd := exec.Command("mysql", dropArgs...)
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("createDatabase: error to run command: %w", err)
	}

	return nil
}

func RestoreData(filename string, args credentials.ConnectionParams) error {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("./%s", filename))
	if err != nil {
		return fmt.Errorf("restoreData: error to read file: %w", err)
	}

	dropArgs := []string{
		fmt.Sprintf("--host=%s", args.Host),
		fmt.Sprintf("--port=%s", args.Port),
		fmt.Sprintf("--user=%s", args.User),
		fmt.Sprintf("--password=%s", args.Password),
		args.Database,
	}

	cmd := exec.Command("mysql", dropArgs...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("restoreData: error to get stdin pipe: %w", err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write(bytes)
	}()

	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("restoreData: error to run command: %w", err)
	}

	return nil
}
