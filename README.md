# promfmt

> This is in an early stage. It works, but may change in the future!

Creates a pattern for [Prometheus](https://prometheus.io) `.rules` files.

In short, files like this:

```
# this is my custom metric
custom_metric = sum(another_metric) BY (name)



# this will alert when blah blah blah blah
alert MyAlert
 IF custom_metric > 10 FOR 10m
  Labels {
	  call= "ghostbusters",
  }
  aNnotations { foo ="bar"}

# sometimes people do comment in blocks
# let's also make that work as expected...
alert Super_Alert
 if blah > 1
 for 1m
```

Will look like this:

```
# this is my custom metric
custom_metric{} = sum(another_metric) BY (name)

# this will alert when blah blah blah blah
ALERT MyAlert
  IF custom_metric > 10
  FOR 10m
  LABELS {
    call = "ghostbusters",
  }
  ANNOTATIONS {
    foo = "bar",
  }

# sometimes people do comment in blocks
# let's also make that work as expected...
ALERT Super_Alert
  IF blah > 1
  FOR 1m
```

## Usage

```console
# format and prints the output to stdout
promfmt FILE

# format and rewrites the file
promfmt --write FILE

# prints the diff
promfmt --diff FILE

# fails if not formatted
promfmt --fail FILE
```


## Installing

On macOS with homebrew:

```console
brew install caarlos0/tap/promfmt
```

Or, on Linux or macOS boxes without homebrew:

```console
curl -sfL https://git.io/promfmt | bash -s -- -b /usr/local/bin
```
