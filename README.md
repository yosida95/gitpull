# gitpull

## About
This is a server program to execute `git pull` in the wake of TCP port knocking.

## How to use
```shell
$ make
$ ./gitpull --socket=":5000" --repository="origin" --refspec="master"
```

## LICENSE
This program is licensed under the [MIT LICENSE](http://yosida95.mit-license.org/)
