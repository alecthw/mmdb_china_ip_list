# mmdb_china_ip_list
Geoip MaxMind Database for china ip list!

git clone https://github.com/maxmind/MaxMind-DB-Writer-perl.git writer

cd writer

curl -LO http://xrl.us/cpanm

perl cpanm –installdeps .
./Build manifest

perl Build.PL
./Build install


{
    "database_type": "GeoLite2-Country",
    "binary_format_major_version": 2,
    "build_epoch": 1589304057,
    "ip_version": 6,
    "languages": [
        "de",
        "en",
        "es",
        "fr",
        "ja",
        "pt-BR",
        "ru",
        "zh-CN"
    ],
    "description": {
        "en": "GeoLite2 Country database"
    },
    "record_size": 24,
    "node_count": 616946,
    "binary_format_minor_version": 0
}



{
    "continent": {
        "code": "AS",
        "names": {
            "de": "Asien",
            "ru": "Азия",
            "pt-BR": "Ásia",
            "ja": "アジア",
            "en": "Asia",
            "fr": "Asie",
            "zh-CN": "亚洲",
            "es": "Asia"
        },
        "geoname_id": 6255147
    },
    "country": {
        "names": {
            "de": "China",
            "ru": "Китай",
            "pt-BR": "China",
            "ja": "中国",
            "en": "China",
            "fr": "Chine",
            "zh-CN": "中国",
            "es": "China"
        },
        "iso_code": "CN",
        "geoname_id": 1814991,
        "is_in_european_union": false,
    },
    "registered_country": {
        "names": {
            "de": "China",
            "ru": "Китай",
            "pt-BR": "China",
            "ja": "中国",
            "en": "China",
            "fr": "Chine",
            "zh-CN": "中国",
            "es": "China"
        },
        "iso_code": "CN",
        "geoname_id": 1814991
    },
    "represented_country": {
        "names": {
            "de": "China",
            "ru": "Китай",
            "pt-BR": "China",
            "ja": "中国",
            "en": "China",
            "fr": "Chine",
            "zh-CN": "中国",
            "es": "China"
        },
        "iso_code": "CN",
        "geoname_id": 1814991
    },
    "traits": {
        "is_anonymous_proxy": true,
        "is_satellite_provider": true
    }
}


https://raw.githubusercontent.com/17mon/china_ip_list/master/china_ip_list.txt
https://raw.githubusercontent.com/metowolf/iplist/master/data/country/CN.txt

https://github.com/maxmind/MaxMind-DB-Writer-perl


http://maxmind.github.io/MaxMind-DB/

https://blog.csdn.net/openex/article/details/53487465

https://www.jianshu.com/p/1b1a018ae729

https://codeload.github.com