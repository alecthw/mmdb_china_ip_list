package main

import (
	"bufio"
	"encoding/csv"
	"strings"

	"log"

	"github.com/maxmind/mmdbwriter/mmdbtype"

	"os"
	"path/filepath"
	"strconv"
)

var (
	chinaIpList    []string
	chunzhenIpList []string

	clangIpV4List []string
	clangIpV6List []string

	cloudCNListV4List []string
	cloudCNListV6List []string

	chinaOperatorIpV4List []string
	chinaOperatorIpV6List []string

	countryChina = mmdbtype.Map{
		"iso_code":             mmdbtype.String("CN"),
		"geoname_id":           mmdbtype.Uint32(1814991),
		"is_in_european_union": mmdbtype.Bool(false),
		"names": mmdbtype.Map{
			"de":    mmdbtype.String("China"),
			"en":    mmdbtype.String("China"),
			"es":    mmdbtype.String("China"),
			"fr":    mmdbtype.String("Chine"),
			"ja":    mmdbtype.String("中国"),
			"pt-BR": mmdbtype.String("China"),
			"ru":    mmdbtype.String("Китай"),
			"zh-CN": mmdbtype.String("中国"),
		},
	}

	cnData = mmdbtype.Map{
		"continent": mmdbtype.Map{
			"code":       mmdbtype.String("AS"),
			"geoname_id": mmdbtype.Uint32(6255147),
			"names": mmdbtype.Map{
				"de":    mmdbtype.String("Asien"),
				"en":    mmdbtype.String("Asia"),
				"es":    mmdbtype.String("Asia"),
				"fr":    mmdbtype.String("Asie"),
				"ja":    mmdbtype.String("アジア"),
				"pt-BR": mmdbtype.String("Ásia"),
				"ru":    mmdbtype.String("Азия"),
				"zh-CN": mmdbtype.String("亚洲"),
			},
		},
		"registered_country": countryChina,
		"country":            countryChina,
	}

	liteCountryMap map[uint64]mmdbtype.Map
)

func init() {
	chunzhenIpList = readFileToStringArray("chunzhen_cn.txt")
	chinaIpList = readFileToStringArray("china_ip_list.txt")
	clangIpV4List = readFileToStringArray("all_cn.txt")
	clangIpV6List = readFileToStringArray("all_cn_ipv6.txt")
	chinaOperatorIpV4List = readFileToStringArray("china_operator_ipv4.txt")
	chinaOperatorIpV6List = readFileToStringArray("china_operator_ipv6.txt")

	initLiteCountryMap()
	initCloudCN()

	log.Printf("chunzhenIpList size, %d\n", len(chunzhenIpList))
	log.Printf("chinaIpList size, %d\n", len(chinaIpList))
	log.Printf("clangIpV4List size, %d\n", len(clangIpV4List))
	log.Printf("clangIpV6List size, %d\n", len(clangIpV6List))
	log.Printf("chinaOperatorIpV4List size, %d\n", len(chinaOperatorIpV4List))
	log.Printf("chinaOperatorIpV6List size, %d\n", len(chinaOperatorIpV6List))
}

func readFileToStringArray(filePath string) []string {
	var strList []string
	fh, err := os.Open(filepath.Join(workDir, filePath))
	if err != nil {
		log.Printf("fail to open %s\n", err)
		return strList
	}
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		strList = append(strList, scanner.Text())
	}

	return strList
}

func initLiteCountryMap() {
	csvFile, err := os.Open(filepath.Join(workDir, "mindmax", "GeoLite2-Country-Locations-zh-CN.csv"))
	if err != nil {
		log.Fatalf("fail to open %s\n", err)
	}
	reader := csv.NewReader(csvFile)

	countryLocationsZhCn, err := reader.ReadAll()
	if err != nil {
		log.Printf("fail to read csv %s\n", err)
		return
	}

	liteCountryMap = make(map[uint64]mmdbtype.Map)

	for index, value := range countryLocationsZhCn {
		if index == 0 {
			continue
		}

		if len(value[4]) == 0 {
			continue
		}

		geoNameId, _ := strconv.ParseUint(value[0], 10, 32)

		liteRecord := mmdbtype.Map{
			"country": mmdbtype.Map{
				"geoname_id": mmdbtype.Uint32(geoNameId),
				"iso_code":   mmdbtype.String(value[4]),
			},
		}

		liteCountryMap[geoNameId] = liteRecord
	}

	//content, _ := json.Marshal(liteCountryMap)
	//log.Printf("%v\n", string(content))
}

func initCloudCN() {
	fh, err := os.Open(filepath.Join(workDir, "cloud_cn.list"))
	if err != nil {
		log.Printf("fail to open %s\n", err)
		return
	}
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		clashRule := scanner.Text()

		if strings.HasPrefix(clashRule, "IP-CIDR,") {
			tmp := strings.Replace(clashRule, "IP-CIDR,", "", -1)
			ipCidr := strings.Replace(tmp, ",no-resolve", "", -1)
			cloudCNListV4List = append(cloudCNListV4List, ipCidr)
		} else if strings.HasPrefix(clashRule, "IP-CIDR6,") {
			tmp := strings.Replace(clashRule, "IP-CIDR6,", "", -1)
			ipCidr := strings.Replace(tmp, ",no-resolve", "", -1)
			cloudCNListV6List = append(cloudCNListV6List, ipCidr)
		}
	}
}
