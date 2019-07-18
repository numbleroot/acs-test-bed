# Test Bed for Evaluating Planet-Scale ACS on Public Clouds

Test bed for planet-scale ACS experiments on public clouds. Proposed and developed as part of my Master Thesis. Among others, utilized to evaluate our resilient mix-net proposal [zeno](https://github.com/numbleroot/zeno).

Folders `cmd` and `scripts` contain executables and scripts required to conduct experiments. Folder `results` contains all measurements gathered and presented in my Master Thesis for reproducibility. File `genplots.py` contains Python code to visualize the results.


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

Deploy the generated executable as a sidecar to the ACS under evaluation on each node of the deployment. Use the following to inspect available flags of the executable:
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
