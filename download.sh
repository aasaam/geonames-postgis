#!/bin/bash

set -e

SCRIPT_DIR=`dirname "0"`
SCRIPT_PATH=`realpath $SCRIPT_DIR`

mkdir -p $SCRIPT_PATH/var/geonames
mkdir -p $SCRIPT_PATH/var/geonames_extract

# geonameid
wget -c -O "$SCRIPT_PATH/var/geonames/geonameid.zip" 'http://download.geonames.org/export/dump/cities500.zip'
rm -rf $SCRIPT_PATH/var/tmp
mkdir -p $SCRIPT_PATH/var/tmp
cp $SCRIPT_PATH/var/geonames/geonameid.zip $SCRIPT_PATH/var/tmp/geonameid.zip
cd $SCRIPT_PATH/var/tmp
7z x geonameid.zip
mv $SCRIPT_PATH/var/tmp/cities500.txt $SCRIPT_PATH/var/geonames_extract/geonameid.tsv

# countryInfo
wget -c -O "$SCRIPT_PATH/var/geonames/countryInfo.tsv" 'https://download.geonames.org/export/dump/countryInfo.txt'
cat "$SCRIPT_PATH/var/geonames/countryInfo.tsv" | grep -v '^#' > $SCRIPT_PATH/var/geonames_extract/countryInfo.tsv

# shapes_all_low
wget -c -O "$SCRIPT_PATH/var/geonames/shapes_all_low.zip" 'https://download.geonames.org/export/dump/shapes_all_low.zip'
rm -rf $SCRIPT_PATH/var/tmp
mkdir -p $SCRIPT_PATH/var/tmp
cp $SCRIPT_PATH/var/geonames/shapes_all_low.zip $SCRIPT_PATH/var/tmp/shapes_all_low.zip
cd $SCRIPT_PATH/var/tmp
7z x shapes_all_low.zip
cat $SCRIPT_PATH/var/tmp/shapes_all_low.txt | grep -E '^[0-9]+.*' > $SCRIPT_PATH/var/geonames_extract/shapes_all_low.tsv

rm -rf $SCRIPT_PATH/var/tmp
