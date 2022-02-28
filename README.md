# Test Bed for Planet-Scale ACS Experiments on Public Clouds

Related publication: ["Strong Anonymity is not Enough: Introducing Fault Tolerance to Planet-Scale Anonymous Communication Systems"](https://dl.acm.org/doi/10.1145/3465481.3469189).

Test bed for planet-scale ACS experiments on public clouds. We developed this to evaluate our
fault-tolerant mixnet proof-of-concept [FTMix](https://github.com/numbleroot/zeno) (formerly *zeno*)
and compare it against state-of-the-art competitors [Vuvuzela](https://github.com/vuvuzela/vuvuzela)
and [Pung](https://github.com/pung-project/pung). Please see
[numbleroot/acs-eval-2019](https://github.com/numbleroot/acs-eval-2019) for forks of Vuvuzela
and Pung with minor adjustments for evaluation purposes, deployed experiment configurations,
and obtained measurements.

Folders `cmd` and `scripts` contain executables and scripts required to conduct the experiments.

**Mind:** At the time of our experiments, [our proof-of-concept fault-tolerant mixnet was still
called *zeno*](https://github.com/numbleroot/zeno#note-on-name-and-scope-of-repository). We renamed
it to *FTMix* (**f**ault-**t**olerant **mix**net) due to scope change and to make its purpose
immediately clear through its name. In order not to create inconsistencies in the data sets, however,
we have not replaced 'zeno' with 'FTMix' in any of the instrumentation code files. Please keep that
in mind when you look at these files.


## Setup

Clone the repository and change into the newly created directory. We assume you have a working
Go installation.

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

Deploy the generated executable as a sidecar to the ACS under evaluation on each deployed
node. Use the following to inspect available flags of the executable:
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
