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
        choose the time counting method (all, mean_epoch_time, profile)
  -csv_path string
        result csv file (if calc_mode = profile)
  -path string
        ***.log file path
```

### Necessary conditions for the launch
To run, you need to have a file with logs from MLPerf

### Launch example

#### Counting the total time
```bash
go build
./analyzer --path=input/file.log --calc_mode=all
```

Output:
```bash
33.868433
```

#### Calculating the average execution time of an epoch
```bash
go build
./analyzer --path=input/file.log --calc_mode=mean_epoch_time
```

Output:
```bash
1.252117
```

#### Profiling
```bash
go build
./analyzer --path=input/nvidia_unet_profile.log --csv_path=output/nvidia_unet_profile.csv --calc_mode=profile
```

Output:
```bash
check output/nvidia_unet_profile.csv
```

## Links
MLCommons:\
https://mlcommons.org

MLCommons MLPerf Training Github:\
https://github.com/mlcommons/training

MLCommons MLPerf Inference Github:\
https://github.com/mlcommons/inference
