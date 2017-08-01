WORKPATH=$(shell pwd)
OUTPATH=bin/

all:
	GOPATH=${WORKPATH} go install beego_main

clean:
	cd ${WORKPATH} && rm -rf ${OUTPATH}