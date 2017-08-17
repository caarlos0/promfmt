# promfmt

> This is in an early stage. It works, but may change in the future!

Creates a pattern for [Prometheus](https://prometheus.io) `.rules` files.

In short, files like this:

```
custom_metric = sum(another_metric) BY (name)

alert MyAlert
 IF custom_metric > 10
  FOR 10m
  Labels {
	  call= "ghostbusters",
  }
  aNnotations {
	  foo ="bar",
  }
```

Will look like this:

```
custom_metric{} = sum(another_metric) BY (name)

ALERT MyAlert
	IF custom_metric > 10
	FOR 10m
	LABELS {
		call = "ghostbusters",
	}
	ANNOTATIONS {
		foo = "bar",
	}
```

