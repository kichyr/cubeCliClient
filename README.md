#  CubeClient(+ TestServer) [![Build Status](https://travis-ci.org/kichyr/cubeCliClient.svg?branch=master)](https://travis-ci.org/kichyr/cubeCliClient)
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