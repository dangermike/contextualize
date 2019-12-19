# Contextualize -- golang edition

## Usage

```text
contextualize -- prefix lines with a previously-found expression

Usage: contextualize [-h, --help] expression [filename]

parameters:
  -h, --help    Display this message
   expression   RE2-style regex. Matched values will be used as the prefix
   filename     Optional filename or glob. If not provided, stdin will be read

The expression is only searched in the first 65535 characters of a line
```

Syntax for RE2 can be found [here](https://github.com/google/re2/wiki/Syntax)

### Example

```shell
contextualize "^(?:func ([^\\(]+)|\\}$)" *.go
```

or

```shell
cat main.go | contextualize "^(?:func ([^\\(]+)|\\}$)"
```


## Installation

This was written with Go modules, but since it uses nothing but standard library packages, that doesn't really matter. You can install it via

`go install`

## Why? What's wrong with you?

The [previous version](https://github.com/dangermike/contextualize-rb) of this utility was a tiny, little Ruby script. It took advantage of Ruby's [ARGF](https://ruby-doc.org/core-2.5.0/ARGF.html) to handle filenames or stdin. Ruby also provides a whole bunch of convenience methods to keep the code short. However, the Ruby version had its deficiencies:

* No support for globs (e.g. "*.txt") or multiple files
* Long lines could eat up memory
* Regular expressions ignored capturing groups
* Insufficiently rad

This version handles all that and more. The tool is comically overbuilt, but it was fun.
