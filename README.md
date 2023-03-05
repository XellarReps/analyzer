# Analyzer

## Description
MLPerf Log Analyzer

### Compilation and launch
In the root of the project, run:

```bash
go build
./analyzer *args*
```

### Command line arguments
The module supports different startup modes and support for command line arguments has been added for this.
To find out information about all the arguments, you need to type in the command line after compilation:
```bash
./analyzer --help
```

Detailed description of each argument:
```text
Usage of ./analyzer:
  -calc_mode string
        choose the time counting method (all, mean_epoch_time)
  -path string
        ***.log file path
```

### Necessary conditions for the launch
To run, you need to have a file with logs from MLPerf

### Launch example
```bash
go build
./analyzer --path=input/file.log --calc_mode=all
```

Output:
```bash
33.868433
```

```bash
go build
./analyzer --path=input/file.log --calc_mode=mean_epoch_time
```

Output:
```bash
1.252117
```
