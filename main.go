package main

import (
	"encoding/csv"
	"flag"
	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	writer *mmdbwriter.Tree

	isNew   bool
	workDir string
	out     string
)

func init() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	flag.BoolVar(&isNew, "new", false, "Do not re using official database and create new mmdb. Default false.")
	flag.StringVar(&workDir, "dir", dir, "The directory which contains china_ip_list.txt and CN.txt. Default executable file directory.")
	flag.StringVar(&out, "out", "china_ip_list.mmdb", "The output mmdb file name. Default china_ip_list.mmdb.")
	flag.Parse()
}

func parseIpNets(strIps []string) []*net.IPNet {
	var ipNets = make([]*net.IPNet, 0, 50)
	for _, strIp := range strIps {
		_, ipNet, err := net.ParseCIDR(strIp)
		if err != nil || ipNet == nil {
			log.Printf("%s fail to parse to CIDR\n", strIp)
			continue
		}
		ipNets = append(ipNets, ipNet)
	}
	return ipNets
}

func insertIps(strIps []string, data mmdbtype.DataType) {
	ipList := parseIpNets(strIps)
	for _, ip := range ipList {
		err := writer.Insert(ip, data)
		if err != nil {
			log.Fatalf("fail to insert to writer %v\n", err)
		}
	}
}

func testResult(mmdbFile string) {
	log.Print("Sleep 5 minutes before test ips.")
	time.Sleep(time.Duration(5) * time.Second)

	//testAllIPDiff("GeoLite2-Country.mmdb", out)
	testSingleIp("1.4.9.249", mmdbFile)
	testSingleIp("2400:bc40::1", mmdbFile)
}

func buildAll() {
	log.Print("Start build all.")

	var err error
	if isNew {
		writer, err = mmdbwriter.New(
			mmdbwriter.Options{
				DatabaseType: "GeoIP2-Country",
				RecordSize:   24,
			},
		)
	} else {
		writer, err = mmdbwriter.Load(filepath.Join(workDir, "GeoLite2-Country.mmdb"), mmdbwriter.Options{})
	}

	if err != nil {
		log.Fatal(err)
	}

	insertIps(clangIpV4List, cnData)
	insertIps(chunzhenIpList, cnData)
	insertIps(chinaIpList, cnData)
	insertIps(clangIpV6List, cnData)

	fh, err := os.Create(filepath.Join(workDir, out))
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("End build all.")

	testResult(out)
}

func buildLite() {
	log.Print("Start build lite.")

	ipv4Csv, err := os.Open(filepath.Join(workDir, "mindmax", "GeoLite2-Country-Blocks-IPv4.csv"))
	if err != nil {
		log.Fatalf("fail to open %s\n", err)
	}
	reader := csv.NewReader(ipv4Csv)

	ipv4s, err := reader.ReadAll()
	if err != nil {
		log.Printf("fail to read csv %s\n", err)
		return
	}

	writer, err = mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "GeoIP2-Country",
			RecordSize:   24,
		},
	)

	for index, value := range ipv4s {
		if index == 0 {
			continue
		}

		_, ipNet, err := net.ParseCIDR(value[0])
		if err != nil || ipNet == nil {
			log.Printf("%s fail to parse to CIDR\n", value[0])
			continue
		}

		geoNameId, _ := strconv.ParseUint(value[1], 10, 32)

		if geoNameId == 1814991 {
			continue
		}

		err = writer.Insert(ipNet, liteCountryMap[geoNameId])
		if err != nil {
			log.Fatalf("fail to insert to writer %v\n", err)
		}
	}

	insertIps(clangIpV4List, liteCountryMap[1814991])
	insertIps(chunzhenIpList, liteCountryMap[1814991])
	insertIps(chinaIpList, liteCountryMap[1814991])
	insertIps(clangIpV6List, liteCountryMap[1814991])

	fh, err := os.Create(filepath.Join(workDir, "lite_"+out))
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("End build lite.")

	testResult("lite_" + out)
}

func main() {
	buildAll()
	buildLite()
}
