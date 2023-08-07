# Keptn Lifecycle Controller - Python Runtime

## Build

```shell
docker build -t lifecycle-toolkit/python-runtime:${VERSION} .
```

## Usage

The Keptn `python-runtime` runner uses python3
and enables the following packages: requests, json, git, yaml

The Keptn Lifecycle Toolkit uses this runner to execute tasks defined as
[KeptnTaskDefinition](https://lifecycle.keptn.sh/docs/yaml-crd-ref/taskdefinition/)
resources.
for pre- and post-checks.

`KeptnTask`s can be tested locally with the runtime using the following commands.
Replace `${VERSION}` with the KLT version of your choice.
`SCRIPT` should refer to either a python file mounted locally in the container or to a url containing the file.

### mounting a python file

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

You can pass python command line arguments by specifying `CMD_ARGS`.
The following example will print the help of python3:

```shell
docker run -e "CMD_ARGS= -help" -it lifecycle-toolkit/python-runtime:${VERSION}
```

### Pass arguments to your python script

In this example we pass one argument (-i test.txt) to the script

```shell
docker run -v $(pwd)/samples/args.py:/args.py -e "SCRIPT=args.py -i test.txt"  -it lifecycle-toolkit/python-runtime:${VERSION}
```

### Use a script from url

We can call the hellopy.py script downloading it directly from github

```shell
docker run -e "SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/python-runtime/samples/hellopy.py" -it lifecycle-toolkit/python-runtime:${VERSION}
```

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />

