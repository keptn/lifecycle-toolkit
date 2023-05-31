# Keptn Lifecycle Controller - Function Runtime

## Build

```shell
docker build -t lifecycle-toolkit/python-runtime:${VERSION} .
```

## Usage

The Keptn python runtime uses python3, and enables the follwing packages: requests, json, git, yaml

The Keptn Lifecycle Toolkit uses this runtime to run [KeptnTask](https://lifecycle.keptn.sh/docs/tasks/write-tasks/)
for pre- and post-checks.

`KeptnTask`s can be tested locally with the runtime using the following command.
Replace `${VERSION}` with the KLT version of your choice.

```shell
docker run -v $(pwd)/samples/hellopy.py:/hellopy.py -e "SCRIPT=hellopy.py" -it lifecycle-toolkit/python-runtime:${VERSION}
```

Where the file in sample/hellopy.py contains python3 code:

```python3
import os

print("Hello, World!")
print(os.environ)
```

This should print in your shell, something like:

```shell
Hello, World!
environ({'HOSTNAME': 'myhost', 'PYTHON_VERSION': '3.9.16', 'PWD': '/', 'CMD_ARGS': '','SCRIPT': 'hellopy.py', ...})
```

### Pass command line arguments to the python command

You can pass python command line arguments by specifying CMD_ARGS. The following example will print the help of python3

```shell
docker run -e "CMD_ARGS= -help" -it lifecycle-toolkit/python-runtime:${VERSION}
```

### Pass arguments to your python script

In this example we pass one argument (-i test.txt) to the script

```shell
docker run -v $(pwd)/samples/args.py:/args.py -e "SCRIPT=args.py -i test.txt"  -it lifecycle-toolkit/python-runtime:${VERSION}
```

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
