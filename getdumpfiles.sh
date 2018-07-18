#!/usr/bin/env bash

rm -rf data || true
mkdir data || true
curl -O http://download.geonames.org/export/dump/cities1000.zip || true
unzip cities1000.zip -d data || true
rm cities1000.zip || true
curl -O http://download.geonames.org/export/dump/alternateNames.zip || true
unzip alternateNames.zip -d data || true
rm alternateNames.zip || true
rm data/iso-languagecodes.txt || true
(cd data && curl -O http://download.geonames.org/export/dump/countryInfo.txt) || true