package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
	"path/filepath"
	"sync/atomic"
)

var (
	db1 *geoip2.Reader
	db2 *geoip2.Reader
)

func hosts(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	return ips, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func loadMmdb(mmdbFile1 string, mmdbFile2 string) {
	var err error

	db1, err = geoip2.Open(filepath.Join(workDir, mmdbFile1))
	if err != nil {
		log.Fatal(err)
	}

	db2, err = geoip2.Open(filepath.Join(workDir, mmdbFile2))
	if err != nil {
		log.Fatal(err)
	}
}

func testIp(dbNum int, ip string) string {
	var db *geoip2.Reader

	if dbNum == 1 {
		db = db1
	} else {
		db = db2
	}

	ipNet := net.ParseIP(ip)
	record, err := db.Country(ipNet)
	if err != nil {
		log.Fatal(err)
	}

	return record.Country.IsoCode
}

func testAllIPDiff(mmdbFile1 string, mmdbFile2 string) {
	loadMmdb(mmdbFile1, mmdbFile2)

	var res []string
	var count uint64 = 0
	var process int = 0

	total := len(chinaIpList)

	for _, ips := range chinaIpList {
		process++
		log.Printf("process: %v/%v, current: %v", process, total, ips)

		hosts, _ := hosts(ips)
		for _, ip := range hosts {
			iso1 := testIp(1, ip)
			iso2 := testIp(2, ip)
			if iso1 != iso2 {
				atomic.AddUint64(&count, 1)
				log.Printf("ip: %v, before: %v, after: %v\n", ip, iso1, iso2)
				res = append(res, fmt.Sprintf("ip: %v, before: %v, after: %v\n", ip, iso1, iso2))
			}
		}
	}

	log.Printf("Different country from ipip cn count: %v\n", count)
	log.Print(res)

	defer db1.Close()
	defer db2.Close()
}

func testSingleIp(ipStr string, mmdbFile string) {
	db, err := geoip2.Open(filepath.Join(workDir, mmdbFile))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(ipStr)
	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("ISO country code: %v\n", record.Country.IsoCode)
}
