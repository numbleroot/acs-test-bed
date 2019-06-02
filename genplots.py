#!/usr/bin/env python3

import sys
import os
import csv
import numpy as np
import matplotlib

from matplotlib import pyplot as plt
from matplotlib import patches as mpatches
from pylab import setp

# matplotlib.rcParams['font.size'] = 10

# Load all measurement files.
metrics = {
    "01": {
        "zeno": {
            "0500": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
        },
        "pung": {
            "0500": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "msg-latencies_lowest-to-highest.data"),
            },
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "msg-latencies_lowest-to-highest.data"),
            },
        },
    },
    "02": {
        "zeno": {
            "0500": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
        },
    },
    "03": {
        "zeno": {
            "0500": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
        },
    },
    "04": {
        "zeno": {
            "0500": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "bandwidth_highest_avg_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "bandwidth_highest_avg_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_cpu_lowest-to-highest_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_mem_lowest-to-highest_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_mem_lowest-to-highest_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "msg-latencies_lowest-to-highest.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "run-01", "msgs-per-mix_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "run-02", "msgs-per-mix_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "run-03", "msgs-per-mix_first-to-last-round.data"),
                },
            },
        },
    },
}


def compileTrafficClients():

    global metrics

    # Ingest and prepare data.

    set01_zeno_0500_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["01"]["zeno"]["0500"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_zeno_0500_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set01_zeno_0500_Bandwidth_Clients_Avg = set01_zeno_0500_Bandwidth_Clients_AvgAll / 500.0

    set01_zeno_1000_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["01"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_zeno_1000_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set01_zeno_1000_Bandwidth_Clients_Avg = set01_zeno_1000_Bandwidth_Clients_AvgAll / 1000.0

    set01_pung_0500_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["01"]["pung"]["0500"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_pung_0500_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set01_pung_0500_Bandwidth_Clients_Avg = set01_pung_0500_Bandwidth_Clients_AvgAll / 500.0

    set01_pung_1000_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["01"]["pung"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_pung_1000_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set01_pung_1000_Bandwidth_Clients_Avg = set01_pung_1000_Bandwidth_Clients_AvgAll / 1000.0

    set02_zeno_0500_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["02"]["zeno"]["0500"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02_zeno_0500_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set02_zeno_0500_Bandwidth_Clients_Avg = set02_zeno_0500_Bandwidth_Clients_AvgAll / 500.0

    set02_zeno_1000_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["02"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02_zeno_1000_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set02_zeno_1000_Bandwidth_Clients_Avg = set02_zeno_1000_Bandwidth_Clients_AvgAll / 1000.0

    set03_zeno_0500_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["03"]["zeno"]["0500"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set03_zeno_0500_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set03_zeno_0500_Bandwidth_Clients_Avg = set03_zeno_0500_Bandwidth_Clients_AvgAll / 500.0

    set03_zeno_1000_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["03"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set03_zeno_1000_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set03_zeno_1000_Bandwidth_Clients_Avg = set03_zeno_1000_Bandwidth_Clients_AvgAll / 1000.0

    set04_zeno_0500_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["04"]["zeno"]["0500"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set04_zeno_0500_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set04_zeno_0500_Bandwidth_Clients_Avg = set04_zeno_0500_Bandwidth_Clients_AvgAll / 500.0

    set04_zeno_1000_Bandwidth_Clients_AvgAll = 0.0
    with open(metrics["04"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set04_zeno_1000_Bandwidth_Clients_AvgAll = float(
            dataFile.read().strip())
    set04_zeno_1000_Bandwidth_Clients_Avg = set04_zeno_1000_Bandwidth_Clients_AvgAll / 1000.0

    bandwidthAvg = [set01_zeno_0500_Bandwidth_Clients_Avg, set01_zeno_1000_Bandwidth_Clients_Avg, set01_pung_0500_Bandwidth_Clients_Avg, set01_pung_1000_Bandwidth_Clients_Avg, set02_zeno_0500_Bandwidth_Clients_Avg,
                    set02_zeno_1000_Bandwidth_Clients_Avg, set03_zeno_0500_Bandwidth_Clients_Avg, set03_zeno_1000_Bandwidth_Clients_Avg, set04_zeno_0500_Bandwidth_Clients_Avg, set04_zeno_1000_Bandwidth_Clients_Avg]

    # Draw plots.

    width = 1.0
    barWidth = (1.0 / 12.0)
    y_max = np.ceil((max(bandwidthAvg) + 1.0))

    _, ax = plt.subplots()

    # Draw all bars.

    ax.bar(1, set01_zeno_0500_Bandwidth_Clients_Avg, width,
           label='zeno (tc off, no failures)', edgecolor='black', color='gold', hatch='/')
    ax.bar(2, set02_zeno_0500_Bandwidth_Clients_Avg, width,
           label='zeno (tc on, no failures)', edgecolor='black', color='gold', hatch='x')
    ax.bar(3, set03_zeno_0500_Bandwidth_Clients_Avg, width,
           label='zeno (tc off, mix failure)', edgecolor='black', color='gold', hatch='o')
    ax.bar(4, set04_zeno_0500_Bandwidth_Clients_Avg, width,
           label='zeno (tc on, mix failure)', edgecolor='black', color='gold', hatch='+')
    ax.bar(5, set01_pung_0500_Bandwidth_Clients_Avg, width,
           label='pung (tc off, no failures)', edgecolor='black', color='steelblue', hatch='\\')
    ax.bar(7, set01_zeno_1000_Bandwidth_Clients_Avg, width,
           edgecolor='black', color='gold', hatch='/')
    ax.bar(8, set02_zeno_1000_Bandwidth_Clients_Avg, width,
           edgecolor='black', color='gold', hatch='x')
    ax.bar(9, set03_zeno_1000_Bandwidth_Clients_Avg, width,
           edgecolor='black', color='gold', hatch='o')
    ax.bar(10, set04_zeno_1000_Bandwidth_Clients_Avg, width,
           edgecolor='black', color='gold', hatch='+')
    ax.bar(11, set01_pung_1000_Bandwidth_Clients_Avg, width,
           edgecolor='black', color='steelblue', hatch='\\')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major',
                  color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 12])
    ax.set_ylim([0, y_max])
    ax.set_xticks((3, 9))
    ax.set_xticklabels(('500 clients', '1000 clients'))

    # Add a legend.
    ax.legend(loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Bandwidth usage (MiB)")

    plt.savefig(os.path.join(
        sys.argv[1], "bandwidth-usage_clients.png"), bbox_inches='tight', dpi=400)


def compileTrafficServers():

    global metrics

    # Ingest and prepare data.

    set01_zeno_0500_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["01"]["zeno"]["0500"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_zeno_0500_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set01_zeno_0500_Bandwidth_Servers_Avg = set01_zeno_0500_Bandwidth_Servers_AvgAll / 21.0

    set01_zeno_1000_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["01"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_zeno_1000_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set01_zeno_1000_Bandwidth_Servers_Avg = set01_zeno_1000_Bandwidth_Servers_AvgAll / 21.0

    set01_pung_0500_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["01"]["pung"]["0500"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_pung_0500_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())

    set01_pung_1000_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["01"]["pung"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_pung_1000_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())

    set02_zeno_0500_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["02"]["zeno"]["0500"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02_zeno_0500_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set02_zeno_0500_Bandwidth_Servers_Avg = set02_zeno_0500_Bandwidth_Servers_AvgAll / 21.0

    set02_zeno_1000_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["02"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02_zeno_1000_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set02_zeno_1000_Bandwidth_Servers_Avg = set02_zeno_1000_Bandwidth_Servers_AvgAll / 21.0

    set03_zeno_0500_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["03"]["zeno"]["0500"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set03_zeno_0500_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set03_zeno_0500_Bandwidth_Servers_Avg = set03_zeno_0500_Bandwidth_Servers_AvgAll / 21.0

    set03_zeno_1000_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["03"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set03_zeno_1000_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set03_zeno_1000_Bandwidth_Servers_Avg = set03_zeno_1000_Bandwidth_Servers_AvgAll / 21.0

    set04_zeno_0500_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["04"]["zeno"]["0500"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set04_zeno_0500_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set04_zeno_0500_Bandwidth_Servers_Avg = set04_zeno_0500_Bandwidth_Servers_AvgAll / 21.0

    set04_zeno_1000_Bandwidth_Servers_AvgAll = 0.0
    with open(metrics["04"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set04_zeno_1000_Bandwidth_Servers_AvgAll = float(
            dataFile.read().strip())
    set04_zeno_1000_Bandwidth_Servers_Avg = set04_zeno_1000_Bandwidth_Servers_AvgAll / 21.0

    bandwidthAvg = [set01_zeno_0500_Bandwidth_Servers_AvgAll, set01_zeno_1000_Bandwidth_Servers_AvgAll, set01_pung_0500_Bandwidth_Servers_AvgAll, set01_pung_1000_Bandwidth_Servers_AvgAll, set02_zeno_0500_Bandwidth_Servers_AvgAll,
                    set02_zeno_1000_Bandwidth_Servers_AvgAll, set03_zeno_0500_Bandwidth_Servers_AvgAll, set03_zeno_1000_Bandwidth_Servers_AvgAll, set04_zeno_0500_Bandwidth_Servers_AvgAll, set04_zeno_1000_Bandwidth_Servers_AvgAll]

    # Draw plots.

    width = 1.0
    barWidth = (1.0 / 12.0)
    y_max = np.ceil((max(bandwidthAvg) + 1000.0))

    _, ax = plt.subplots()

    # Draw all bars and corresponding average lines.

    ax.bar(1, set01_zeno_0500_Bandwidth_Servers_AvgAll, width,
           label='zeno (tc off, no failures)', edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_0500_Bandwidth_Servers_Avg,
                xmin=(0.5 * barWidth), xmax=(1.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(2, set02_zeno_0500_Bandwidth_Servers_AvgAll, width,
           label='zeno (tc on, no failures)', edgecolor='black', color='gold', hatch='x')
    plt.axhline(y=set02_zeno_0500_Bandwidth_Servers_Avg,
                xmin=(1.5 * barWidth), xmax=(2.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(3, set03_zeno_0500_Bandwidth_Servers_AvgAll, width,
           label='zeno (tc off, mix failure)', edgecolor='black', color='gold', hatch='o')
    plt.axhline(y=set03_zeno_0500_Bandwidth_Servers_Avg,
                xmin=(2.5 * barWidth), xmax=(3.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(4, set04_zeno_0500_Bandwidth_Servers_AvgAll, width,
           label='zeno (tc on, mix failure)', edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set04_zeno_0500_Bandwidth_Servers_Avg,
                xmin=(3.5 * barWidth), xmax=(4.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(5, set01_pung_0500_Bandwidth_Servers_AvgAll, width,
           label='pung (tc off, no failures)', edgecolor='black', color='steelblue', hatch='\\')

    ax.bar(7, set01_zeno_1000_Bandwidth_Servers_AvgAll,
           width, edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_1000_Bandwidth_Servers_Avg,
                xmin=(6.5 * barWidth), xmax=(7.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(8, set02_zeno_1000_Bandwidth_Servers_AvgAll,
           width, edgecolor='black', color='gold', hatch='x')
    plt.axhline(y=set02_zeno_1000_Bandwidth_Servers_Avg,
                xmin=(7.5 * barWidth), xmax=(8.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(9, set03_zeno_1000_Bandwidth_Servers_AvgAll,
           width, edgecolor='black', color='gold', hatch='o')
    plt.axhline(y=set03_zeno_1000_Bandwidth_Servers_Avg,
                xmin=(8.5 * barWidth), xmax=(9.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(10, set04_zeno_1000_Bandwidth_Servers_AvgAll,
           width, edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set04_zeno_1000_Bandwidth_Servers_Avg,
                xmin=(9.5 * barWidth), xmax=(10.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(11, set01_pung_1000_Bandwidth_Servers_AvgAll, width,
           edgecolor='black', color='steelblue', hatch='\\')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major',
                  color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 12])
    ax.set_ylim([0, y_max])
    ax.set_xticks((3, 9))
    ax.set_xticklabels(('500 clients', '1000 clients'))

    # Add a legend.
    ax.legend(loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Bandwidth usage (MiB)")

    plt.savefig(os.path.join(
        sys.argv[1], "bandwidth-usage_servers.png"), bbox_inches='tight', dpi=400)


def compileLoad(path, bothSystems, isClientsFigure):

    global metrics

    # Ingest data.

    zeno0500LoadCPU = []
    with open(zeno0500LoadCPUPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            zeno0500LoadCPU.append(list(map(float, row)))

    zeno0500LoadMem = []
    with open(zeno0500LoadMemPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            zeno0500LoadMem.append(list(map(float, row)))

    zeno1000LoadCPU = []
    with open(zeno1000LoadCPUPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            zeno1000LoadCPU.append(list(map(float, row)))

    zeno1000LoadMem = []
    with open(zeno1000LoadMemPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            zeno1000LoadMem.append(list(map(float, row)))

    pung0500LoadCPU = []
    with open(pung0500LoadCPUPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            pung0500LoadCPU.append(list(map(float, row)))

    pung0500LoadMem = []
    with open(pung0500LoadMemPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            pung0500LoadMem.append(list(map(float, row)))

    pung1000LoadCPU = []
    with open(pung1000LoadCPUPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            pung1000LoadCPU.append(list(map(float, row)))

    pung1000LoadMem = []
    with open(pung1000LoadMemPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            pung1000LoadMem.append(list(map(float, row)))

    # Draw plots.

    _, ax = plt.subplots()

    zeno01 = ax.boxplot(zeno0500LoadCPU, positions=[
                        2], widths=0.85, patch_artist=True, whis='range')
    zeno02 = ax.boxplot(zeno0500LoadMem, positions=[
                        3], widths=0.85, patch_artist=True, whis='range')
    pung01 = ax.boxplot(pung0500LoadCPU, positions=[
                        5], widths=0.85, patch_artist=True, whis='range')
    pung02 = ax.boxplot(pung0500LoadMem, positions=[
                        6], widths=0.85, patch_artist=True, whis='range')

    zeno03 = ax.boxplot(zeno1000LoadCPU, positions=[
                        9], widths=0.85, patch_artist=True, whis='range')
    zeno04 = ax.boxplot(zeno1000LoadMem, positions=[
                        10], widths=0.85, patch_artist=True, whis='range')
    pung03 = ax.boxplot(pung1000LoadCPU, positions=[
                        12], widths=0.85, patch_artist=True, whis='range')
    pung04 = ax.boxplot(pung1000LoadMem, positions=[
                        13], widths=0.85, patch_artist=True, whis='range')

    # Color boxes.
    setp(zeno01['boxes'], color='gold')
    setp(zeno02['boxes'], color='gold')
    setp(zeno03['boxes'], color='gold')
    setp(zeno04['boxes'], color='gold')

    setp(pung01['boxes'], color='steelblue')
    setp(pung02['boxes'], color='steelblue')
    setp(pung03['boxes'], color='steelblue')
    setp(pung04['boxes'], color='steelblue')

    ax.yaxis.grid(True, linestyle='-', which='major',
                  color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 15])
    ax.set_ylim([0.0, 100.0])
    ax.set_xticks((4, 11))
    ax.set_xticklabels(('500 clients', '1000 clients'))

    zeno_patch = mpatches.Patch(color='gold', label='zeno')
    pung_patch = mpatches.Patch(color='steelblue', label='pung')

    plt.legend(handles=[zeno_patch, pung_patch])

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Load (percentage)")

    if isClientsFigure:
        plt.savefig(os.path.join(
            sys.argv[1], "01_tc-off_proc-off_load_clients.png"), bbox_inches='tight', dpi=400)
    else:
        plt.savefig(os.path.join(
            sys.argv[1], "01_tc-off_proc-off_load_servers.png"), bbox_inches='tight', dpi=400)


def compileLatencies(path, bothSystems):

    global metrics

    # Ingest data.

    zeno0500Latencies = []
    with open(zeno0500LatenciesPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            zeno0500Latencies.append(list(map(float, row)))

    pung0500Latencies = []
    with open(pung0500LatenciesPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            pung0500Latencies.append(list(map(float, row)))

    zeno1000Latencies = []
    with open(zeno1000LatenciesPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            zeno1000Latencies.append(list(map(float, row)))

    pung1000Latencies = []
    with open(pung1000LatenciesPath, newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            pung1000Latencies.append(list(map(float, row)))

    # Draw plots.

    _, ax = plt.subplots()

    n_bins = 1000
    ax.hist(zeno0500Latencies, n_bins, density=True, histtype='step',
            cumulative=True, label='zeno (500 clients)')
    ax.hist(pung0500Latencies, n_bins, density=True, histtype='step',
            cumulative=True, label='pung (500 clients)')
    ax.hist(zeno1000Latencies, n_bins, density=True, histtype='step',
            cumulative=True, label='zeno (1000 clients)')
    ax.hist(pung1000Latencies, n_bins, density=True, histtype='step',
            cumulative=True, label='pung (1000 clients)')

    ax.legend(loc='lower right')
    ax.set_axisbelow(True)

    plt.grid()
    plt.tight_layout()

    plt.xlabel("End-to-end transmission latency (seconds)")
    plt.ylabel("Fraction of messages transmitted")

    plt.savefig(os.path.join(
        sys.argv[1], "01_tc-off_proc-off_msg-latencies.png"), bbox_inches='tight', dpi=400)


def compileMessagesPerMix():

    global metrics

    for setting in metrics:

        for numClients in {"0500", "1000"}:

            for run in {"Run01", "Run02", "Run03"}:

                outputFile = os.path.join(os.path.dirname(
                    metrics[setting]["zeno"][numClients]["MessagesPerMix"][run]), "msgs-per-mix_first-to-last-round.png")

                labels = []
                data = []
                with open(metrics[setting]["zeno"][numClients]["MessagesPerMix"][run], newline='') as dataFile:
                    reader = csv.reader(dataFile, delimiter=',')
                    for idx, row in enumerate(reader):
                        if idx == 0:
                            labels = row
                        else:
                            data.append(list(map(int, row)))

                flat_data = [count for mix in data for count in mix]
                y_max = np.ceil(max(flat_data) * 1.07)

                x_max = len(data[0])
                for msgCounts in data:
                    if len(msgCounts) > x_max:
                        x_max = len(msgCounts)

                _, ax = plt.subplots()
                ax.set_xlim([0, x_max])
                ax.set_ylim([0, y_max])

                for idx, msgCounts in enumerate(data):
                    plt.plot(
                        msgCounts, "-", label=labels[idx], markersize=2.0, color=np.random.rand(3,))

                boxOfPlot = ax.get_position()
                ax.set_position([boxOfPlot.x0, boxOfPlot.y0,
                                 (boxOfPlot.width * 0.8), boxOfPlot.height])
                ax.legend(loc='center left', bbox_to_anchor=(
                    1, 0.5), fontsize='small')

                plt.grid()
                plt.tight_layout()

                plt.xlabel("Round number")
                plt.ylabel("Messages in all pools (count)")

                plt.savefig(outputFile, bbox_inches='tight', dpi=400)

# Create all figures.


# Build bandwidth figures.
compileTrafficClients()
# compileTrafficServers()

# Build load usage figures.
# compileLoad()

# Build message latencies figure.
# compileLatencies()

# Build figures describing the number of
# messages in each mix server over rounds.
# compileMessagesPerMix()
