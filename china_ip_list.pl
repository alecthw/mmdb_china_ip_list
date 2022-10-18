#!/usr/bin/perl

use utf8;
binmode(STDIN, ':encoding(utf8)');
binmode(STDOUT, ':encoding(utf8)');
binmode(STDERR, ':encoding(utf8)');

use MaxMind::DB::Writer::Tree;
use Path::Class;
use Data::Dumper;

# prepare continent and country data ---- Begin ----

my @mm_langs = ('de', 'ru', 'pt-BR', 'ja', 'en', 'fr', 'zh-CN', 'es');

my @mm_continent_codes = (6255146, 6255147, 6255148, 6255149, 6255150, 6255151, 6255152);
my %mm_continentid2name_map = (
    6255146 => 'AF',
    6255147 => 'AS',
    6255148 => 'EU',
    6255149 => 'NA',
    6255150 => 'SA',
    6255151 => 'OC',
    6255152 => 'AN',
);
my %mm_continent_map = (
    AF => {
        code => 'AF',
        names => {
            de => 'Afrika',
            ru => 'Африка',
            'pt-BR' => 'África',
            ja => 'アフリカ',
            en => 'Africa',
            fr => 'Afrique',
            'zh-CN' => '非洲',
            es => 'África',
        },
        geoname_id => 6255146,
    },
    AS => {
        code => 'AS',
        names => {
            de => 'Asien',
            ru => 'Азия',
            'pt-BR' => 'Ásia',
            ja => 'アジア',
            en => 'Asia',
            fr => 'Asie',
            'zh-CN' => '亚洲',
            es => 'Asia',
        },
        geoname_id => 6255147,
    },
    EU => {
        code => 'EU',
        names => {
            de => 'Europa',
            ru => 'Европа',
            'pt-BR' => 'Europa',
            ja => 'ヨーロッパ',
            en => 'Europe',
            fr => 'Europe',
            'zh-CN' => '欧洲',
            es => 'Europa',
        },
        geoname_id => 6255148,
    },
    NA => {
        code => 'NA',
        names => {
            de => 'Nordamerika',
            ru => 'Северная Америка',
            'pt-BR' => 'América do Norte',
            ja => '北アメリカ',
            en => 'North America',
            fr => 'Amérique du Nord',
            'zh-CN' => '北美洲',
            es => 'Norteamérica',
        },
        geoname_id => 6255149,
    },
    SA => {
        code => 'SA',
        names => {
            de => 'Südamerika',
            ru => 'Южная Америка',
            'pt-BR' => 'América do Sul',
            ja => '南アメリカ',
            en => 'South America',
            fr => 'Amérique du Sud',
            'zh-CN' => '南美洲',
            es => 'Sudamérica',
        },
        geoname_id => 6255150,
    },
    OC => {
        code => 'OC',
        names => {
            de => 'Ozeanien',
            ru => 'Океания',
            'pt-BR' => 'Oceania',
            ja => 'オセアニア',
            en => 'Oceania',
            fr => 'Océanie',
            'zh-CN' => '大洋洲',
            es => 'Oceanía',
        },
        geoname_id => 6255151,
    },
    AN => {
        code => 'AN',
        names => {
            de => 'Antarktis',
            ru => 'Антарктика',
            'pt-BR' => 'Antártica',
            ja => '南極大陸',
            en => 'Antarctica',
            fr => 'Antarctique',
            'zh-CN' => '南极洲',
            es => 'Antártida',
        },
        geoname_id => 6255152,
    },
);

my %mm_country_map = ();
my %mm_country2continent_map = ();

sub mm_country_map_insert{
    my @values = split(',', $_[0]);
    my $first_build = $_[1];

    my $geoname_id = $values[0];
    my $locale_code = $values[1];
    my $continent_code = $values[2];
    # my $continent_name = $values[3];
    my $country_iso_code = $values[4];
    my $country_name = $values[5];
    my $is_in_european_union = $values[6];

    if ( $geoname_id eq 'geoname_id' ) {
        return;
    }
    if ( grep { $_ eq $geoname_id } @mm_continent_codes ) {
        return;
    }

    $mm_country_map{$geoname_id}{'names'}{$locale_code} = $country_name;
    if ( $first_build > 0 ) {
        $mm_country_map{$geoname_id}{'iso_code'} = $country_iso_code;
        $mm_country_map{$geoname_id}{'geoname_id'} = $geoname_id;

        $mm_country2continent_map{$geoname_id} = $continent_code;
    }
    if ( $is_in_european_union > 0 ) {
        $mm_country_map{$geoname_id}{'is_in_european_union'} = 1;
    }
}

sub build_country_map{
    my $dir = dir('./mindmax');

    # 标记第一次构建，后续的构建仅添加不同语言的名称
    my $first_build = 1;
    foreach $mm_lang (@mm_langs){
        my $file_name = 'GeoLite2-Country-Locations-';
        $file_name .= $mm_lang ;
        $file_name .= '.csv' ;

        # print '$file_name\n';

        my $file = $dir->file($file_name);
        my $content = $file->slurp();
        my $file_handle = $file->openr();
        binmode($file_handle, ':utf8');
        while( my $line = $file_handle->getline() ) {
            $line =~ s/\"//g;
            mm_country_map_insert($line, $first_build);
        }
        $first_build = 0;
    }

}

build_country_map;

# print Dumper(\%mm_continent_map);
# print Dumper(\%mm_country_map);

# prepare continent and country data ---- End ----


my %types = (
    continent => 'map',
    country => 'map',
    registered_country => 'map',
    represented_country => 'map',
    traits => 'map',
    code => 'utf8_string',
    names => 'map',
    geoname_id => 'uint32',
    iso_code => 'utf8_string',

    is_in_european_union => 'boolean',
    is_anonymous_proxy => 'boolean',
    is_satellite_provider => 'boolean',

    de => 'utf8_string',
    ru => 'utf8_string',
    'pt-BR' => 'utf8_string',
    ja => 'utf8_string',
    en => 'utf8_string',
    fr => 'utf8_string',
    'zh-CN' => 'utf8_string',
    es => 'utf8_string',
);

my $tree = MaxMind::DB::Writer::Tree->new(
    ip_version                  => 6,
    record_size                 => 24,
    database_type               => 'GeoLite2-Country',
    languages                   => ['de', 'ru', 'pt-BR', 'ja', 'en', 'fr', 'zh-CN', 'es'],
    description                 => { en => 'GeoLite2 Country database' },
    map_key_type_callback       => sub { $types{ $_[0] } },
);

sub insert_maxmind_ip{
    my $dir = dir('./mindmax');
    my $file = $dir->file($_[0]);
    my $content = $file->slurp();
    my $file_handle = $file->openr();
    binmode($file_handle, ':utf8');
    while( my $line = $file_handle->getline() ) {
        my @values = split(',', $line);

        my $network = $values[0];
        my $geoname_id = $values[1];
        my $registered_country_geoname_id = $values[2];
        my $represented_country_geoname_id = $values[3];
        my $is_anonymous_proxy = $values[4];
        my $is_satellite_provider = $values[5];

        if ( $network ne 'network' ) {
            my $data = {};

            if ( grep { $_ eq $geoname_id } @mm_continent_codes ) {
                $data -> {continent} = $mm_continent_map{$mm_continentid2name_map{$geoname_id}};
            } else {
                if ( exists($mm_country_map{$geoname_id} ) ) {
                    $data -> {continent} = $mm_continent_map{$mm_country2continent_map{$geoname_id}};
                    $data -> {country} = $mm_country_map{$geoname_id};
                }
            }

            if ( exists($mm_country_map{$registered_country_geoname_id} ) ) {
                $data -> {registered_country} = $mm_country_map{$registered_country_geoname_id};
            }

            if ( exists($mm_country_map{$represented_country_geoname_id} ) ) {
                $data -> {represented_country} = $mm_country_map{$represented_country_geoname_id};
            }

            if ( $is_anonymous_proxy > 0 ) {
                $data -> {traits} -> {is_anonymous_proxy} = 1;
            }

            if ( $is_satellite_provider > 0 ) {
                $data -> {traits} -> {is_satellite_provider} = 1;
            }

            # print '$network\n';
            # print Dumper($data);

            $tree->insert_network(
                $network,
                $data,
            );
        }
    }
}

sub insert_china_ip{
    my $dir = dir('.');
    my $file = $dir->file($_[0]);
    my $content = $file->slurp();
    my $file_handle = $file->openr();
    binmode($file_handle, ':utf8');
    while( my $line = $file_handle->getline() ) {
        my $data = {};
        $data -> {continent} = $mm_continent_map{'AS'};
        $data -> {country} = $mm_country_map{1814991};
        $data -> {registered_country} = $mm_country_map{1814991};
        # print Dumper($data);
        $tree->insert_network(
            $line,
            $data,
        );
    }
}

insert_maxmind_ip('GeoLite2-Country-Blocks-IPv4.csv');
insert_maxmind_ip('GeoLite2-Country-Blocks-IPv6.csv');
insert_china_ip('china_ip_list.txt');
insert_china_ip('CN.txt');

open my $fh, '>:raw', 'china_ip_list.mmdb';
$tree->write_tree($fh);
