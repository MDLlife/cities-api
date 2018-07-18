#!/usr/bin/env bash

rm -rf cities.db
rm -rf $PWD/data || true
mkdir $PWD/data || true
curl -O http://download.geonames.org/export/dump/cities1000.zip || true
unzip $PWD/cities1000.zip -d $PWD/data || true
rm $PWD/cities1000.zip || true
curl -O http://download.geonames.org/export/dump/alternateNames.zip || true
unzip $PWD/alternateNames.zip -d $PWD/data || true
rm $PWD/alternateNames.zip || true
rm $PWD/data/iso-languagecodes.txt || true
(cd $PWD/data && curl -O http://download.geonames.org/export/dump/countryInfo.txt) || true