package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
)

/*
Script Name: Restore Mongo DB Backup BSON Files in One Click
Github URI: https://github.com/riyasop/go-scripts.git
Author: Riyas OP
Author URI: https://riyas.dev
Description: To Restore Multiple BSON Mongo DB Collection Backup Files in One Click
Version: 1.0
Tags: mongodb, bson, restore mongo db, restore mongodb, restore mongo db , restore bson files, restore mongodb collections
*/
var wg sync.WaitGroup

func execute(file string) {
	fname := path.Base(file)
	cmd := exec.Command("mongorestore", "--uri", os.Args[1], "-d", os.Args[2], "-c", fname[:len(fname)-5], file)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
	}

	// start the command after having set up the pipe
	if err := cmd.Start(); err != nil {
	}

	// read command's stdout line by line
	in := bufio.NewScanner(stdout)

	for in.Scan() {
		log.Printf(in.Text()) // write each line to your log, or anything you need
	}
	if err := in.Err(); err != nil {
		log.Printf("error: %s", err)
	}
	wg.Done()
}
func main() {

	if len(os.Args) < 4 || os.Args[1] == "--help" {
		fmt.Println(`
--help for Help

Run this Command with following Arguments
1. Connection string
2. Database Name
3. Directory Path

---------Run Before Build----------
go run mongo_restore_bson_go.go <connection string> <database> <directory full path>
Example:-
go run mongo_restore_bson_go.go "mongodb+srv://username:password@xxxxx.mongodb.net" test /home/abcd/mongobackup

--To Build--
go build mongo_restore_bson_go.go

---------Run After Build----------
./mongo_restore_bson_go <connection string> <database> <directory full path>
Example:-
./mongo_restore_bson_go "mongodb+srv://username:password@xxxxxxx.mongodb.net" test /home/abcd/mongobackup

		`)
		return
	}
	err := filepath.Walk(os.Args[3], func(fullpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if fullpath[len(fullpath)-4:] == "bson" {
			wg.Add(1)
			go execute(fullpath)
		}
		return nil
	})
	wg.Wait()
	if err != nil {
		panic(err)
	}

}
