# MeXiCo

MeXiCo is an esoteric programming language, and a *compiler* for the language with the same name, compiling source code to [DNS MX records](https://en.m.wikipedia.org/wiki/MX_record). This also explains the name.

It's greatly inspired by [blinry](https://morr.cc)'s [legit](https://github.com/blinry/legit) project and [brainfuck](https://en.wikipedia.org/wiki/Brainfuck), and is written in [go](https://golang.org).
This project was born and written on a boring railroad trip from Hamburg to Düsseldorf.

## Machine specification

Like some other esoteric programming languages, a MeXiCo machine has a storage of unlimited size, also called *infinite tape*, and a first-in-last-out stack. It's possible to read from and write to the tape with a movable *head*. This makes MeXiCo turing-complete.

## Design

As stated above, the source code of a MeXiCo program completely resides in MX records. The compiler ensures that generated MX records are RFC-conform. This means, it is possible to deliver *MeXiCo* source code over a DNS server of your choice, and, using the [time-to-live](https://en.wikipedia.org/wiki/Time_to_live) value, cache source code on a DNS resolver of your choice.

Like in the old [BASIC](https://en.wikipedia.org/wiki/BASIC) days, source code *lines* are defined by a number at the beginning of a line, sitting in the [Priority](https://en.wikipedia.org/wiki/MX_record#MX_preference,_distance,_and_priority) value of the MX record. Lines which are not filled out are skipped. As a short explanation:

```
someprogram.esolang.mil		IN	MX	10 <line 10>
someprogram.esolang.mil		IN	MX	20 <line 20>
someprogram.esolang.mil		IN	MX	30 <line 30>
```

Due to the fact that the payload of a MX records needs to be a [FQDN](https://en.wikipedia.org/wiki/FQDN), every command is represented as a subdomain of the domain `mexico.invalid.`, which is obviously non-existent. If a command contains spaces, for example if they carry an argument (`push 5`), every space is replaced by a minus sign.

## Instructions

Below, you can find the instructions used in the MX records.

| Command | Involves | Consumes Stack | Pushes to Stack | Description |
| --- | --- | --- | --- | --- |
| `left` | Tape Head | 0 | 0 | Moves the tape head one cell to the left |
| `right` | Tape Head | 0 | 0 |  Moves the tape head one cell to the right |
| `pusht` | Tape, Stack | 0 | 1 | Reads the current cell value and pushes it on top of the stack |
| `push n` | Stack | 0 | 1 | Pushes the value `n` to the stack |
| `pop` | Tape, Stack | 1 | 0 | Pops top stack value to the current cell |
| `dup` | Stack | 1 | 2 | Duplicates the topmost stack value |
| `del` | Stack | 1 | 0 | Deletes the topmost stack value, ignoring its value |
| `eq` | Stack | 2 | 1 | Checks if `stack[0] == stack[1]`. Pushes `1` to the stack if equal, `0` otherwise | 
| `not` | Stack | 1 | 1 | Inverses `stack[0]` |
| `gt` | Stack | 2 | 1 | Checks if `stack[0] > stack[1]`. Pushes `1` to the stack if greater, `0` otherwise |
| `lt` | Stack | 2 | 1 | Checks if `stack[0] < stack[1]`. Pushes `1` to the stack if smaller, `0` otherwise |
| `add` | Stack | 2 | 1 | Calculates `stack[0] + stack[1]`, and pushes the result to the stack |
| `sub` | Stack | 2 | 1 | Calculates `stack[0] - stack[1]`, and pushes the result to the stack |
| `mult` | Stack | 2 | 1 | Calculates `stack[0] * stack[1]`, and pushes the result to the stack |
| `div` | Stack | 2 | 1 | Calculates `stack[0] / stack[1]`, and pushes the result to the stack |
| `mod` | Stack | 2 | 1 | Calculates `stack[0] % stack[1]`, and pushes the result to the stack |
| `read` | Stack | 0 | 1 | Reads a character from the user, and pushes its char value to the stack |
| `print` | Stack | 1 | 0 | Prints `stack[0]` as a character |
| `jmp` | Program Flow, Stack | 1 | 0 | Jumps to the line number specified by `stack[0]` |
| `jmpc` | Program Flow, Stack | 2 | 0 | Jumps to the line number specified by `stack[0]`, if `stack[1]` is not `0`. |

Please note that `stack[0]` refers to the topmost stack value, and `stack[i]` refers to the i-th stack value.

## Source code

The syntax of the source code is strongly aligned with the Instructions table above. However, there's a bit of *syntactical sugar* to make programming in this language enjoyable. Take a look into the `examples` directory of this repository to get a basic idea of it.

### Comments

Every line starting with `#`, `//` or `;` is ignored by the compiler.

### Labels

You can define labels like this:

```
// This program reads a character from the user, and subtracts 1 from it, until it is zero.
read

// Let's loop here
LOOP:
push 1
sub 
push 0
lt
push LOOP
jmpc

// We're done!
```

As you can see, labels can be defined with a `:` after the label name, and it can be used as a value for `push`. At compile time, it is replaced with the corresponding line number.

## Implementations

There is a reference implementation for the compiler, `mexico`, and a reference implementation for the interpreter, `mexigo`. Both can be found in this repository.

### Compiler "mexico"

Simply run `go get github.com/maride/mexico/mexico` to get the compiler.

The mexico compiler takes three arguments:

- `-input` to specify the source code file
- `-output` to specify the output path for the zonefile
- `-baseDomain`, the base domain to compile the source code for. This should be the domain you are planning to host the source code on.

For example. to compile the `Fibonacci.mxc` example for the domain `fibonacci.mxc.maride.cc`, you could use this command:

`./mexico --input ../examples/Fibonacci.mxc --output /srv/zones/fibonacci.mxc.maride.cc --baseDomain fibonacci.mxc.maride.cc`

If no problems occurred and the compiler didn't run into an issue, nothing is printed.

### Interpreter "mexigo"

Simply run `go get github.com/maride/mexico/mexigo` to get the interpreter.

The mexigo interpreter takes only one argument - the domain to execute:

`./mexigo fibonacci.mxc.maride.cc`

This will give you an output similar to this:

```
> $ ./mexigo fibonacci.mxc.maride.cc
mexigo - the reference interpreter for the mexico esolang!
See github.com/maride/mexico for further information.

2019/12/08 17:22:27 Resolving fibonacci.mxc.maride.cc for MX records
2019/12/08 17:22:27 Found 24 code lines, interpreting them...
'\x02' (2)
'\x03' (3)
'\x05' (5)
'\b' (8)
'\r' (13)
'\x15' (21)
'"' (34)
'7' (55)
'Y' (89)
'\u0090' (144)
'é' (233)
'Ź' (377)
'ɢ' (610)
'ϛ' (987)
'ؽ' (1597)
2019/12/08 17:22:27 Stack is currently 0 entries big
2019/12/08 17:22:27 Tape is currently 2 cells big
Cell 0: 987 (ϛ)
Cell 1: 1597 (ؽ)
2019/12/08 17:22:27 Found no commands after line 24. Stopping.
```

## Examples

You can find examples in the `examples` directory of this repository.

I currently host the `Fibonacci.mxc` on `fibonacci.mxc.maride.cc`, means you can run it with `mexigo fibonacci.mxc.maride.cc`!

I challenge you to write more examples. ;)

## License

I chose to release this project under the [CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/) license.
