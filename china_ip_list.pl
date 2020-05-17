use MaxMind::DB::Writer::Tree;

use Path::Class;
use utf8;

my %types = (
    continent => 'map',
    country => 'map',
    registered_country => 'map',
    code => 'utf8_string',
    names => 'map',
    geoname_id => 'uint32',
    iso_code => 'utf8_string',

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
    ip_version            => 6,
    record_size           => 24,
    database_type         => 'GeoLite2-Country',
    languages             => ['en'],
    description           => { en => 'GeoLite2 Country database' },
    map_key_type_callback => sub { $types{ $_[0] } },
);

sub tree_insert_network{
    $tree->insert_network(
        $_[0],
        {
            continent => {
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
            country => {
                names => {
                    de => 'China',
                    ru => 'Китай',
                    'pt-BR' => 'China',
                    ja => '中国',
                    en => 'China',
                    fr => 'China',
                    'zh-CN' => '中国',
                    es => 'China',
                },
                iso_code => 'CN',
                geoname_id => 1814991,
            },
            registered_country => {
                names => {
                    de => 'China',
                    ru => 'Китай',
                    'pt-BR' => 'China',
                    ja => '中国',
                    en => 'China',
                    fr => 'China',
                    'zh-CN' => '中国',
                    es => 'China',
                },
                iso_code => 'CN',
                geoname_id => 1814991,
            },
        },
    );
}

sub build_tree{
    my $dir = dir(".");
    my $file = $dir->file($_[0]);
    my $content = $file->slurp();
    my $file_handle = $file->openr();
    binmode($file_handle, ":utf8");
    while( my $line = $file_handle->getline() ) {
        tree_insert_network($line);
    }
}
build_tree('china_ip_list.txt');

open my $fh, '>:raw', 'china_ip_list.mmdb';
$tree->write_tree($fh);