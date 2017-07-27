WORKPATH=$(shell pwd)
OUTPATH=bin/
all:
	GOPATH=${WORKPATH} go build -o ${OUTPATH}cryptogen

clean:
	cd ${WORKPATH} && rm -rf ${OUTPATH}