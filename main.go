package main

import (
	"bufio"
	"log"
	"net"
	"os"

	mmdb "github.com/maxmind/mmdbwriter"
	mmtype "github.com/maxmind/mmdbwriter/mmdbtype"
)

var cnRecord = mmtype.Map{
	"country": mmtype.Map{
		"iso_code": mmtype.String("CN"),
	},
}

func insertRecords(db *mmdb.Tree, name string) (err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()
	scanner := bufio.NewScanner(f)
	var ipnet *net.IPNet
	for scanner.Scan() {
		_, ipnet, err = net.ParseCIDR(scanner.Text())
		if err != nil {
			return
		}
		err = db.Insert(ipnet, cnRecord)
		if err != nil {
			return
		}
	}
	return scanner.Err()
}

func writeFile(db *mmdb.Tree, name string) (err error) {
	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer func() {
		cerr := f.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = db.WriteTo(f)
	return
}

func conv(from string, to string) (err error) {
	db, err := mmdb.New(mmdb.Options{
		DatabaseType: "GeoIP2-Country",
		RecordSize:   24,
	})
	if err != nil {
		return
	}

	err = insertRecords(db, from)
	if err != nil {
		return
	}
	err = writeFile(db, to)
	return
}

func main() {
	err := conv(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
}
