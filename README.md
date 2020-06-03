#  CubeClient(+ TestServer) [![Build Status](https://travis-ci.org/kichyr/cubeCliClient.svg?branch=master)](https://travis-ci.org/kichyr/cubeCliClient)
This repo contains implementation of CubeClient and Cube test server that simulate oath2 communication with authentication server in purpose of token scope check.

Quick start:
```
$ make start-test-server &
$ make build
$ ./cubeclient localhost 8091 test1 read
```

Run test in docker:
```
$ make test
```

Run tests locally:
```
$ pip3 install -r ./test/requirements.txt
$ make test-local
```