default: test

build:
	go build

getdumpfiles:
	mkdir data
	curl -O http://download.geonames.org/export/dump/cities1000.zip
	unzip cities1000.zip -d data
	rm cities1000.zip
	curl -O http://download.geonames.org/export/dump/alternateNames.zip
	unzip alternateNames.zip -d data
	rm alternateNames.zip
	rm data/iso-languagecodes.txt
	(cd data && curl -O http://download.geonames.org/export/dump/countryInfo.txt)

prepare: getdumpfiles

dockerbuild:
	docker build -t cities .

dockerrun: dockerbuild
	docker run -t -p 80:8080 --name cities --rm cities

dockerrm:
	docker stop cities || true
	docker rm cities || true

dockerrerun: dockerrm dockerrun

test:
	godep go vet ./...
	godep go test ./... -cover
