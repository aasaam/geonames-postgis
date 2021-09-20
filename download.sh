#!/bin/bash

set -e

SCRIPT_DIR=`dirname "0"`
SCRIPT_PATH=`realpath $SCRIPT_DIR`

mkdir -p $SCRIPT_PATH/var/geonames
mkdir -p $SCRIPT_PATH/var/geonames_extract

# allCountries
wget -c -O "$SCRIPT_PATH/var/geonames/allCountries.zip" 'https://download.geonames.org/export/dump/allCountries.zip'
rm -rf $SCRIPT_PATH/var/tmp
mkdir -p $SCRIPT_PATH/var/tmp
cp $SCRIPT_PATH/var/geonames/allCountries.zip $SCRIPT_PATH/var/tmp/allCountries.zip
cd $SCRIPT_PATH/var/tmp
7z x allCountries.zip
mv $SCRIPT_PATH/var/tmp/allCountries.txt $SCRIPT_PATH/var/geonames_extract/allCountries.tsv

# alternateNamesV2
wget -c -O "$SCRIPT_PATH/var/geonames/alternateNamesV2.zip" 'https://download.geonames.org/export/dump/alternateNamesV2.zip'
rm -rf $SCRIPT_PATH/var/tmp
mkdir -p $SCRIPT_PATH/var/tmp
cp $SCRIPT_PATH/var/geonames/alternateNamesV2.zip $SCRIPT_PATH/var/tmp/alternateNamesV2.zip
cd $SCRIPT_PATH/var/tmp
7z x alternateNamesV2.zip
mv $SCRIPT_PATH/var/tmp/alternateNamesV2.txt $SCRIPT_PATH/var/geonames_extract/alternateNamesV2.tsv
mv $SCRIPT_PATH/var/tmp/iso-languagecodes.txt $SCRIPT_PATH/var/geonames_extract/iso-languagecodes.tsv

# countryInfo
wget -c -O "$SCRIPT_PATH/var/geonames/countryInfo.tsv" 'https://download.geonames.org/export/dump/countryInfo.txt'
cat "$SCRIPT_PATH/var/geonames/countryInfo.tsv" | grep -v '^#' > $SCRIPT_PATH/var/geonames_extract/countryInfo.tsv

# hierarchy
wget -c -O "$SCRIPT_PATH/var/geonames/hierarchy.zip" 'https://download.geonames.org/export/dump/hierarchy.zip'
rm -rf $SCRIPT_PATH/var/tmp
mkdir -p $SCRIPT_PATH/var/tmp
cp $SCRIPT_PATH/var/geonames/hierarchy.zip $SCRIPT_PATH/var/tmp/hierarchy.zip
cd $SCRIPT_PATH/var/tmp
7z x hierarchy.zip
mv $SCRIPT_PATH/var/tmp/hierarchy.txt $SCRIPT_PATH/var/geonames_extract/hierarchy.tsv

# shapes_all_low
wget -c -O "$SCRIPT_PATH/var/geonames/shapes_all_low.zip" 'https://download.geonames.org/export/dump/shapes_all_low.zip'
rm -rf $SCRIPT_PATH/var/tmp
mkdir -p $SCRIPT_PATH/var/tmp
cp $SCRIPT_PATH/var/geonames/shapes_all_low.zip $SCRIPT_PATH/var/tmp/shapes_all_low.zip
cd $SCRIPT_PATH/var/tmp
7z x shapes_all_low.zip
cat $SCRIPT_PATH/var/tmp/shapes_all_low.txt | grep -E '^[0-9]+.*' > $SCRIPT_PATH/var/geonames_extract/shapes_all_low.tsv

rm -rf $SCRIPT_PATH/var/tmp
