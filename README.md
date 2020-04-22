# Test Bed for Planet-Scale ACS Experiments on Public Clouds

Test bed for planet-scale ACS experiments on public clouds. We developed this to evaluate our fault-tolerant mix-net
prototype [zeno](https://github.com/numbleroot/zeno) and to compare it against state-of-the-art competitors
[Vuvuzela](https://github.com/vuvuzela/vuvuzela) and [Pung](https://github.com/pung-project/pung). Please see
[numbleroot/acs-eval-2019](https://github.com/numbleroot/acs-eval-2019) for forks of Vuvuzela and Pung with minor
adjustments for evaluation purposes, deployed experiment configurations, and obtained measurements.

Folders `cmd` and `scripts` contain executables and scripts required to conduct the experiments. File `genplots.py`
contains Python code to visualize the results.


## Setup

Clone the repository and change into the newly created directory. We assume you have a working Go installation.

### Generate Configuration Files

Run:
```
$ make genconfigs
```

And use the following to inspect available flags of the created executable:
```
$ ./genconfigs -help
```

### Run Experiments

Run:
```
$ make runexperiments
```

And use the following to inspect available flags of the created executable:
```
$ ./runexperiments -help
```

### Run Collector Executable as Sidecar on Nodes

Run:
```
$ make collector
```

Deploy the generated executable as a sidecar to the ACS under evaluation on each deployed node. Use the following
to inspect available flags of the executable:
```
$ ./collector -help
```

### Perform Calculations Across Gathered Measurements

Run:
```
$ make calcstats
```

And use the following to inspect available flags of the created executable:
```
$ ./calcstats -help
```
