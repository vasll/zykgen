# Zykgen
A fork of [zykgen](https://github.com/luc10/zykgen) with added featueres

## Install the utility
You can install it on your system with the `go get` command
```go
go get github.com/vasll/zykgen/cmd/zykgen
```

## Usage
```
zykgen --pass (-m|-n|-c) [-o <file>] [-l <length>] <serial>
zykgen --dump (-m|-n|-c) -o <file> [-l <length>] <RouterSerialRange>
```
- `--pass or --dump` The zykgen mode you want to use
- `(-m|-n|-c)` What letterlist you want to use for the key generation
- `[-o <file>]` Optional for the `--pass` mode but mandatory for the `--dump` mode
- `-l <length>` The length of the password to generate, by default 10. _For example the home&life routers use 16 chars._
- `<serial>` The serial of the router
- `<RouterSerialRange>` One of the three router serial ranges included i.e: `homelife`, `infostrada`, `tiscali`


## Examples
### zykgen --pass
`--pass` generates a password starting from a router's serial
```powershell
zykgen --pass -c -l 16 S182V30001171
```
Output
```
Serial: S182V30001171     
Password: M8TN4BPPLLT4NJ84
```

### zykgen --dump
`--dump` creates a file containing all the passwords between a serial router range, for example it can generate password from the home&life serial ranges which are `S182V00000000-S182V99999999` and `S192V00000000-S192V99999999`
```powershell
zykgen --dump -c -o test.txt -l 16 -t 4 homelife
```