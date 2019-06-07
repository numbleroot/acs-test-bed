#!/usr/bin/env python3

import sys
import os
import csv
import numpy as np
import matplotlib

from matplotlib import pyplot as plt
from matplotlib import patches as mpatches
from matplotlib import ticker as ticker
from pylab import setp

matplotlib.use("pgf")
plt.rcParams.update({
    "font.family": "serif",
    "text.usetex": True,
    "pgf.rcfonts": False,
    "pgf.preamble": [
        "\\usepackage{units}",
        "\\usepackage{metalogo}",
        "\\usepackage{amsfonts}",
        "\\usepackage{amsmath}",
        "\\usepackage{amssymb}",
        "\\usepackage{amsthm}",
        "\\usepackage{mathtools}",
        "\\usepackage{paratype}",
        "\\usepackage{FiraMono}",
    ]
})

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
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-0500", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-0500", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-0500", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc-on_proc-off", "zeno", "clients-1000", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-0500", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc-off_proc-on", "zeno", "clients-1000", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-0500", "load_mem_avg_servers.data"),
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
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_mem_avg_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_cpu_lowest-to-highest_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "04_tc-on_proc-on", "zeno", "clients-1000", "load_mem_avg_servers.data"),
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

    bandwidthAvg = [set01_zeno_0500_Bandwidth_Clients_Avg, set02_zeno_0500_Bandwidth_Clients_Avg,
                    set03_zeno_0500_Bandwidth_Clients_Avg, set04_zeno_0500_Bandwidth_Clients_Avg,
                    set01_pung_0500_Bandwidth_Clients_Avg, set01_zeno_1000_Bandwidth_Clients_Avg,
                    set02_zeno_1000_Bandwidth_Clients_Avg, set03_zeno_1000_Bandwidth_Clients_Avg,
                    set04_zeno_1000_Bandwidth_Clients_Avg, set01_pung_1000_Bandwidth_Clients_Avg]

    # Draw plots.

    width = 1.0
    y_max = np.ceil((max(bandwidthAvg) + 10.0))

    _, ax = plt.subplots()

    # Draw all bars.

    ax.bar(1, set01_zeno_0500_Bandwidth_Clients_Avg, width, label='zeno (tc off, no failures)', edgecolor='black', color='gold', hatch='/')
    ax.bar(2, set02_zeno_0500_Bandwidth_Clients_Avg, width, label='zeno (tc on, no failures)', edgecolor='black', color='gold', hatch='x')
    ax.bar(3, set03_zeno_0500_Bandwidth_Clients_Avg, width, label='zeno (tc off, mix failure)', edgecolor='black', color='gold', hatch='o')
    ax.bar(4, set04_zeno_0500_Bandwidth_Clients_Avg, width, label='zeno (tc on, mix failure)', edgecolor='black', color='gold', hatch='+')
    ax.bar(5, set01_pung_0500_Bandwidth_Clients_Avg, width, label='pung (tc off, no failures)', edgecolor='black', color='steelblue', hatch='\\')
    ax.bar(7, set01_zeno_1000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='/')
    ax.bar(8, set02_zeno_1000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='x')
    ax.bar(9, set03_zeno_1000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='o')
    ax.bar(10, set04_zeno_1000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='+')
    ax.bar(11, set01_pung_1000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='\\')

    labels = ["%.2f" % avg for avg in bandwidthAvg]

    for bar, label in zip(ax.patches, labels):
        ax.text((bar.get_x() + (bar.get_width() / 2)), (bar.get_height() * 1.05), label, ha='center', va='bottom')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 12])
    ax.set_ylim([0, y_max])
    ax.set_xticks((3, 9))
    ax.set_xticklabels(('500 clients', '1000 clients'))

    # Add a legend.
    ax.legend(loc='upper left')

    plt.yscale('symlog')
    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Bandwidth usage (MiB) [log.]")

    plt.savefig(os.path.join(sys.argv[1], "bandwidth-usage_clients.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "bandwidth-usage_clients.pgf"), bbox_inches='tight')


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

    bandwidthAvg = [set01_zeno_0500_Bandwidth_Servers_AvgAll, set02_zeno_0500_Bandwidth_Servers_AvgAll,
                    set03_zeno_0500_Bandwidth_Servers_AvgAll, set04_zeno_0500_Bandwidth_Servers_AvgAll,
                    set01_pung_0500_Bandwidth_Servers_AvgAll, set01_zeno_1000_Bandwidth_Servers_AvgAll,
                    set02_zeno_1000_Bandwidth_Servers_AvgAll, set03_zeno_1000_Bandwidth_Servers_AvgAll,
                    set04_zeno_1000_Bandwidth_Servers_AvgAll, set01_pung_1000_Bandwidth_Servers_AvgAll]

    # Draw plots.

    width = 1.0
    barWidth = (1.0 / 12.0)
    y_max = np.ceil((max(bandwidthAvg) + 1000.0))

    _, ax = plt.subplots()

    # Draw all bars and corresponding average lines.

    ax.bar(1, set01_zeno_0500_Bandwidth_Servers_AvgAll, width, label='zeno (tc off, no failures)', edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_0500_Bandwidth_Servers_Avg, xmin=(0.5 * barWidth), xmax=(1.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(2, set02_zeno_0500_Bandwidth_Servers_AvgAll, width, label='zeno (tc on, no failures)', edgecolor='black', color='gold', hatch='x')
    plt.axhline(y=set02_zeno_0500_Bandwidth_Servers_Avg, xmin=(1.5 * barWidth), xmax=(2.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(3, set03_zeno_0500_Bandwidth_Servers_AvgAll, width, label='zeno (tc off, mix failure)', edgecolor='black', color='gold', hatch='o')
    plt.axhline(y=set03_zeno_0500_Bandwidth_Servers_Avg, xmin=(2.5 * barWidth), xmax=(3.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(4, set04_zeno_0500_Bandwidth_Servers_AvgAll, width, label='zeno (tc on, mix failure)', edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set04_zeno_0500_Bandwidth_Servers_Avg, xmin=(3.5 * barWidth), xmax=(4.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(5, set01_pung_0500_Bandwidth_Servers_AvgAll, width, label='pung (tc off, no failures)', edgecolor='black', color='steelblue', hatch='\\')

    ax.bar(7, set01_zeno_1000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_1000_Bandwidth_Servers_Avg, xmin=(6.5 * barWidth), xmax=(7.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(8, set02_zeno_1000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='x')
    plt.axhline(y=set02_zeno_1000_Bandwidth_Servers_Avg, xmin=(7.5 * barWidth), xmax=(8.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(9, set03_zeno_1000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='o')
    plt.axhline(y=set03_zeno_1000_Bandwidth_Servers_Avg, xmin=(8.5 * barWidth), xmax=(9.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(10, set04_zeno_1000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set04_zeno_1000_Bandwidth_Servers_Avg, xmin=(9.5 * barWidth), xmax=(10.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(11, set01_pung_1000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='steelblue', hatch='\\')
    
    labels = ["%.2f" % avg for avg in bandwidthAvg]

    for bar, label in zip(ax.patches, labels):
        ax.text((bar.get_x() + (bar.get_width() / 2)), (bar.get_height() + 200), label, ha='center', va='bottom')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
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

    plt.savefig(os.path.join(sys.argv[1], "bandwidth-usage_servers.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "bandwidth-usage_servers.pgf"), bbox_inches='tight')


def compileLoadCPUClients():

    global metrics

    # Ingest data.

    set01_zeno_0500_Load_CPU = []
    with open(metrics["01"]["zeno"]["0500"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_0500_Load_CPU.append(float(item))

    set01_pung_0500_Load_CPU = []
    with open(metrics["01"]["pung"]["0500"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_0500_Load_CPU.append(float(item))

    set01_zeno_1000_Load_CPU = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_CPU.append(float(item))

    set01_pung_1000_Load_CPU = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_CPU.append(float(item))

    set02_zeno_0500_Load_CPU = []
    with open(metrics["02"]["zeno"]["0500"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_0500_Load_CPU.append(float(item))

    set02_zeno_1000_Load_CPU = []
    with open(metrics["02"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_1000_Load_CPU.append(float(item))

    set03_zeno_0500_Load_CPU = []
    with open(metrics["03"]["zeno"]["0500"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_0500_Load_CPU.append(float(item))

    set03_zeno_1000_Load_CPU = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_CPU.append(float(item))

    set04_zeno_0500_Load_CPU = []
    with open(metrics["04"]["zeno"]["0500"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_0500_Load_CPU.append(float(item))

    set04_zeno_1000_Load_CPU = []
    with open(metrics["04"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_1000_Load_CPU.append(float(item))

    # Draw plots.

    width = 0.9

    _, ax = plt.subplots()

    set01_zeno01 = ax.boxplot(set01_zeno_0500_Load_CPU, positions=[1], widths=width, patch_artist=True, whis='range')
    set02_zeno01 = ax.boxplot(set02_zeno_0500_Load_CPU, positions=[2], widths=width, patch_artist=True, whis='range')
    set03_zeno01 = ax.boxplot(set03_zeno_0500_Load_CPU, positions=[3], widths=width, patch_artist=True, whis='range')
    set04_zeno01 = ax.boxplot(set04_zeno_0500_Load_CPU, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_pung01 = ax.boxplot(set01_pung_0500_Load_CPU, positions=[5], widths=width, patch_artist=True, whis='range')
    set01_zeno03 = ax.boxplot(set01_zeno_1000_Load_CPU, positions=[7], widths=width, patch_artist=True, whis='range')
    set02_zeno03 = ax.boxplot(set02_zeno_1000_Load_CPU, positions=[8], widths=width, patch_artist=True, whis='range')
    set03_zeno03 = ax.boxplot(set03_zeno_1000_Load_CPU, positions=[9], widths=width, patch_artist=True, whis='range')
    set04_zeno03 = ax.boxplot(set04_zeno_1000_Load_CPU, positions=[10], widths=width, patch_artist=True, whis='range')
    set01_pung03 = ax.boxplot(set01_pung_1000_Load_CPU, positions=[11], widths=width, patch_artist=True, whis='range')

    # Color boxplots.
    setp(set01_zeno01['boxes'], color='black')
    setp(set01_zeno01['boxes'], facecolor='gold')
    setp(set01_zeno01['boxes'], hatch='/')

    setp(set02_zeno01['boxes'], color='black')
    setp(set02_zeno01['boxes'], facecolor='gold')
    setp(set02_zeno01['boxes'], hatch='x')

    setp(set03_zeno01['boxes'], color='black')
    setp(set03_zeno01['boxes'], facecolor='gold')
    setp(set03_zeno01['boxes'], hatch='o')

    setp(set04_zeno01['boxes'], color='black')
    setp(set04_zeno01['boxes'], facecolor='gold')
    setp(set04_zeno01['boxes'], hatch='+')

    setp(set01_pung01['boxes'], color='black')
    setp(set01_pung01['boxes'], facecolor='steelblue')
    setp(set01_pung01['boxes'], hatch='\\')

    setp(set01_zeno03['boxes'], color='black')
    setp(set01_zeno03['boxes'], facecolor='gold')
    setp(set01_zeno03['boxes'], hatch='/')

    setp(set02_zeno03['boxes'], color='black')
    setp(set02_zeno03['boxes'], facecolor='gold')
    setp(set02_zeno03['boxes'], hatch='x')

    setp(set03_zeno03['boxes'], color='black')
    setp(set03_zeno03['boxes'], facecolor='gold')
    setp(set03_zeno03['boxes'], hatch='o')

    setp(set04_zeno03['boxes'], color='black')
    setp(set04_zeno03['boxes'], facecolor='gold')
    setp(set04_zeno03['boxes'], hatch='+')

    setp(set01_pung03['boxes'], color='black')
    setp(set01_pung03['boxes'], facecolor='steelblue')
    setp(set01_pung03['boxes'], hatch='\\')
    
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 12])
    ax.set_ylim([0.0, 50.0])
    ax.set_xticks((3.5, 9.5))
    ax.set_yticks([0, 10, 20, 30, 40, 50, 50])
    ax.set_xticklabels(('500 clients', '1000 clients'))

    # Add a legend.
    ax.legend([set01_zeno01['boxes'][0], set02_zeno01['boxes'][0], set03_zeno01['boxes'][0],
        set04_zeno01['boxes'][0], set01_pung01['boxes'][0]], ['zeno (tc off, no failures)',
        'zeno (tc on, no failures)', 'zeno (tc off, mix failure)', 'zeno (tc on, mix failure)',
        'pung (tc off, no failures)'], loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Busy CPU (percentage)")

    plt.savefig(os.path.join(sys.argv[1], "load-cpu_clients.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "load-cpu_clients.pgf"), bbox_inches='tight')


def compileLoadMemClients():

    global metrics

    # Ingest data.

    set01_zeno_0500_Load_Mem = 0.0
    with open(metrics["01"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set01_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set01_pung_0500_Load_Mem = 0.0
    with open(metrics["01"]["pung"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set01_pung_0500_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set01_zeno_1000_Load_Mem = 0.0
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set01_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set01_pung_1000_Load_Mem = 0.0
    with open(metrics["01"]["pung"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set01_pung_1000_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set02_zeno_0500_Load_Mem = 0.0
    with open(metrics["02"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set02_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set02_zeno_1000_Load_Mem = 0.0
    with open(metrics["02"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set02_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set03_zeno_0500_Load_Mem = 0.0
    with open(metrics["03"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set03_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set03_zeno_1000_Load_Mem = 0.0
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set03_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set04_zeno_0500_Load_Mem = 0.0
    with open(metrics["04"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set04_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000.0

    set04_zeno_1000_Load_Mem = 0.0
    with open(metrics["04"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        set04_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000.0

    # Draw plots.

    width = 0.9

    _, ax = plt.subplots()

    ax.bar(1, set01_zeno_0500_Load_Mem, width, label='zeno (tc off, no failures)', color='gold', hatch='/')
    ax.bar(2, set02_zeno_0500_Load_Mem, width, label='zeno (tc on, no failures)', color='gold', hatch='x')
    ax.bar(3, set03_zeno_0500_Load_Mem, width, label='zeno (tc off, mix failure)', color='gold', hatch='o')
    ax.bar(4, set04_zeno_0500_Load_Mem, width, label='zeno (tc on, mix failure)', color='gold', hatch='+')
    ax.bar(5, set01_pung_0500_Load_Mem, width, label='pung (tc off, no failures)', color='steelblue', hatch='\\')
    ax.bar(7, set01_zeno_1000_Load_Mem, width, color='gold', hatch='/')
    ax.bar(8, set02_zeno_1000_Load_Mem, width, color='gold', hatch='x')
    ax.bar(9, set03_zeno_1000_Load_Mem, width, color='gold', hatch='o')
    ax.bar(10, set04_zeno_1000_Load_Mem, width, color='gold', hatch='+')
    ax.bar(11, set01_pung_1000_Load_Mem, width, color='steelblue', hatch='\\')
    
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 12])
    ax.set_xticks((3.5, 9.5))
    ax.set_yticks([0, 50, 100, 150, 200, 250, 300, 350, 400, 450, 500])
    ax.set_xticklabels(('500 clients', '1000 clients'))

    # Add a legend.
    ax.legend(loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Used memory (MB)")

    plt.savefig(os.path.join(sys.argv[1], "load-mem_clients.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "load-mem_clients.pgf"), bbox_inches='tight')


def compileLoadCPUServers():

    global metrics

    # Ingest data.

    set01_zeno_0500_Load_CPU = []
    with open(metrics["01"]["zeno"]["0500"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_0500_Load_CPU.append(float(item))

    set01_pung_0500_Load_CPU = []
    with open(metrics["01"]["pung"]["0500"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_0500_Load_CPU.append(float(item))

    set01_zeno_1000_Load_CPU = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_CPU.append(float(item))

    set01_pung_1000_Load_CPU = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_CPU.append(float(item))

    set02_zeno_0500_Load_CPU = []
    with open(metrics["02"]["zeno"]["0500"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_0500_Load_CPU.append(float(item))

    set02_zeno_1000_Load_CPU = []
    with open(metrics["02"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_1000_Load_CPU.append(float(item))

    set03_zeno_0500_Load_CPU = []
    with open(metrics["03"]["zeno"]["0500"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_0500_Load_CPU.append(float(item))

    set03_zeno_1000_Load_CPU = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_CPU.append(float(item))

    set04_zeno_0500_Load_CPU = []
    with open(metrics["04"]["zeno"]["0500"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_0500_Load_CPU.append(float(item))

    set04_zeno_1000_Load_CPU = []
    with open(metrics["04"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_1000_Load_CPU.append(float(item))

    # Draw plots.

    width = 0.9

    _, ax = plt.subplots()

    set01_zeno01 = ax.boxplot(set01_zeno_0500_Load_CPU, positions=[1], widths=width, patch_artist=True, whis='range')
    set02_zeno01 = ax.boxplot(set02_zeno_0500_Load_CPU, positions=[2], widths=width, patch_artist=True, whis='range')
    set03_zeno01 = ax.boxplot(set03_zeno_0500_Load_CPU, positions=[3], widths=width, patch_artist=True, whis='range')
    set04_zeno01 = ax.boxplot(set04_zeno_0500_Load_CPU, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_pung01 = ax.boxplot(set01_pung_0500_Load_CPU, positions=[5], widths=width, patch_artist=True, whis='range')
    set01_zeno03 = ax.boxplot(set01_zeno_1000_Load_CPU, positions=[7], widths=width, patch_artist=True, whis='range')
    set02_zeno03 = ax.boxplot(set02_zeno_1000_Load_CPU, positions=[8], widths=width, patch_artist=True, whis='range')
    set03_zeno03 = ax.boxplot(set03_zeno_1000_Load_CPU, positions=[9], widths=width, patch_artist=True, whis='range')
    set04_zeno03 = ax.boxplot(set04_zeno_1000_Load_CPU, positions=[10], widths=width, patch_artist=True, whis='range')
    set01_pung03 = ax.boxplot(set01_pung_1000_Load_CPU, positions=[11], widths=width, patch_artist=True, whis='range')

    # Color boxplots.
    setp(set01_zeno01['boxes'], color='black')
    setp(set01_zeno01['boxes'], facecolor='gold')
    setp(set01_zeno01['boxes'], hatch='/')

    setp(set02_zeno01['boxes'], color='black')
    setp(set02_zeno01['boxes'], facecolor='gold')
    setp(set02_zeno01['boxes'], hatch='x')

    setp(set03_zeno01['boxes'], color='black')
    setp(set03_zeno01['boxes'], facecolor='gold')
    setp(set03_zeno01['boxes'], hatch='o')

    setp(set04_zeno01['boxes'], color='black')
    setp(set04_zeno01['boxes'], facecolor='gold')
    setp(set04_zeno01['boxes'], hatch='+')

    setp(set01_pung01['boxes'], color='black')
    setp(set01_pung01['boxes'], facecolor='steelblue')
    setp(set01_pung01['boxes'], hatch='\\')

    setp(set01_zeno03['boxes'], color='black')
    setp(set01_zeno03['boxes'], facecolor='gold')
    setp(set01_zeno03['boxes'], hatch='/')

    setp(set02_zeno03['boxes'], color='black')
    setp(set02_zeno03['boxes'], facecolor='gold')
    setp(set02_zeno03['boxes'], hatch='x')

    setp(set03_zeno03['boxes'], color='black')
    setp(set03_zeno03['boxes'], facecolor='gold')
    setp(set03_zeno03['boxes'], hatch='o')

    setp(set04_zeno03['boxes'], color='black')
    setp(set04_zeno03['boxes'], facecolor='gold')
    setp(set04_zeno03['boxes'], hatch='+')

    setp(set01_pung03['boxes'], color='black')
    setp(set01_pung03['boxes'], facecolor='steelblue')
    setp(set01_pung03['boxes'], hatch='\\')

    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 12])
    ax.set_ylim([0.0, 50.0])
    ax.set_xticks((3.5, 9.5))
    ax.set_yticks([0, 10, 20, 30, 40, 50])
    ax.set_xticklabels(('500 clients', '1000 clients'))

    # Add a legend.
    ax.legend([set01_zeno01['boxes'][0], set02_zeno01['boxes'][0], set03_zeno01['boxes'][0],
        set04_zeno01['boxes'][0], set01_pung01['boxes'][0]], ['zeno (tc off, no failures)',
        'zeno (tc on, no failures)', 'zeno (tc off, mix failure)', 'zeno (tc on, mix failure)',
        'pung (tc off, no failures)'], loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Busy CPU (percentage)")

    plt.savefig(os.path.join(sys.argv[1], "load-cpu_servers.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "load-cpu_servers.pgf"), bbox_inches='tight')


def compileLoadMemServers():

    global metrics

    # Ingest data.

    set01_zeno_0500_Load_Mem = 0.0
    with open(metrics["01"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set01_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set01_pung_0500_Load_Mem = 0.0
    with open(metrics["01"]["pung"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set01_pung_0500_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set01_zeno_1000_Load_Mem = 0.0
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set01_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set01_pung_1000_Load_Mem = 0.0
    with open(metrics["01"]["pung"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set01_pung_1000_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set02_zeno_0500_Load_Mem = 0.0
    with open(metrics["02"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set02_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set02_zeno_1000_Load_Mem = 0.0
    with open(metrics["02"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set02_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set03_zeno_0500_Load_Mem = 0.0
    with open(metrics["03"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set03_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set03_zeno_1000_Load_Mem = 0.0
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set03_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set04_zeno_0500_Load_Mem = 0.0
    with open(metrics["04"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set04_zeno_0500_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    set04_zeno_1000_Load_Mem = 0.0
    with open(metrics["04"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        set04_zeno_1000_Load_Mem = float(dataFile.read().strip()) / 1000000.0

    # Draw plots.

    width = 0.9

    _, ax = plt.subplots()

    ax.bar(1, set01_zeno_0500_Load_Mem, width, label='zeno (tc off, no failures)', color='gold', hatch='/')
    ax.bar(2, set02_zeno_0500_Load_Mem, width, label='zeno (tc on, no failures)', color='gold', hatch='x')
    ax.bar(3, set03_zeno_0500_Load_Mem, width, label='zeno (tc off, mix failure)', color='gold', hatch='o')
    ax.bar(4, set04_zeno_0500_Load_Mem, width, label='zeno (tc on, mix failure)', color='gold', hatch='+')
    ax.bar(5, set01_pung_0500_Load_Mem, width, label='pung (tc off, no failures)', color='steelblue', hatch='\\')
    ax.bar(7, set01_zeno_1000_Load_Mem, width, color='gold', hatch='/')
    ax.bar(8, set02_zeno_1000_Load_Mem, width, color='gold', hatch='x')
    ax.bar(9, set03_zeno_1000_Load_Mem, width, color='gold', hatch='o')
    ax.bar(10, set04_zeno_1000_Load_Mem, width, color='gold', hatch='+')
    ax.bar(11, set01_pung_1000_Load_Mem, width, color='steelblue', hatch='\\')

    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 12])
    ax.set_xticks((3.5, 9.5))
    ax.set_xticklabels(('500 clients', '1000 clients'))
    ax.set_yticks([0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10])

    # Add a legend.
    ax.legend(loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Used memory (GB)")

    plt.savefig(os.path.join(sys.argv[1], "load-mem_servers.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "load-mem_servers.pgf"), bbox_inches='tight')


def compileLatencies():

    global metrics

    # Ingest data.

    set01_zeno_0500_Latencies = []
    with open(metrics["01"]["zeno"]["0500"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_0500_Latencies.append(float(item))

    set01_zeno_1000_Latencies = []
    with open(metrics["01"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Latencies.append(float(item))

    set02_zeno_0500_Latencies = []
    with open(metrics["02"]["zeno"]["0500"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_0500_Latencies.append(float(item))

    set02_zeno_1000_Latencies = []
    with open(metrics["02"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_1000_Latencies.append(float(item))

    set03_zeno_0500_Latencies = []
    with open(metrics["03"]["zeno"]["0500"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_0500_Latencies.append(float(item))

    set03_zeno_1000_Latencies = []
    with open(metrics["03"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Latencies.append(float(item))

    set04_zeno_0500_Latencies = []
    with open(metrics["04"]["zeno"]["0500"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_0500_Latencies.append(float(item))

    set04_zeno_1000_Latencies = []
    with open(metrics["04"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_1000_Latencies.append(float(item))

    set01_pung_0500_Latencies = []
    with open(metrics["01"]["pung"]["0500"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_0500_Latencies.append(float(item))

    set01_pung_1000_Latencies = []
    with open(metrics["01"]["pung"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Latencies.append(float(item))

    # Prepare CDF arrays.

    set01_zeno_0500_Latencies = np.sort(set01_zeno_0500_Latencies)
    set01_zeno_0500_CDF = np.array(range(len(set01_zeno_0500_Latencies))) / float(len(set01_zeno_0500_Latencies))

    set01_zeno_1000_Latencies = np.sort(set01_zeno_1000_Latencies)
    set01_zeno_1000_CDF = np.array(range(len(set01_zeno_1000_Latencies))) / float(len(set01_zeno_1000_Latencies))

    set02_zeno_0500_Latencies = np.sort(set02_zeno_0500_Latencies)
    set02_zeno_0500_CDF = np.array(range(len(set02_zeno_0500_Latencies))) / float(len(set02_zeno_0500_Latencies))

    set02_zeno_1000_Latencies = np.sort(set02_zeno_1000_Latencies)
    set02_zeno_1000_CDF = np.array(range(len(set02_zeno_1000_Latencies))) / float(len(set02_zeno_1000_Latencies))

    set03_zeno_0500_Latencies = np.sort(set03_zeno_0500_Latencies)
    set03_zeno_0500_CDF = np.array(range(len(set03_zeno_0500_Latencies))) / float(len(set03_zeno_0500_Latencies))

    set03_zeno_1000_Latencies = np.sort(set03_zeno_1000_Latencies)
    set03_zeno_1000_CDF = np.array(range(len(set03_zeno_1000_Latencies))) / float(len(set03_zeno_1000_Latencies))

    set04_zeno_0500_Latencies = np.sort(set04_zeno_0500_Latencies)
    set04_zeno_0500_CDF = np.array(range(len(set04_zeno_0500_Latencies))) / float(len(set04_zeno_0500_Latencies))

    set04_zeno_1000_Latencies = np.sort(set04_zeno_1000_Latencies)
    set04_zeno_1000_CDF = np.array(range(len(set04_zeno_1000_Latencies))) / float(len(set04_zeno_1000_Latencies))

    set01_pung_0500_Latencies = np.sort(set01_pung_0500_Latencies)
    set01_pung_0500_CDF = np.array(range(len(set01_pung_0500_Latencies))) / float(len(set01_pung_0500_Latencies))

    set01_pung_1000_Latencies = np.sort(set01_pung_1000_Latencies)
    set01_pung_1000_CDF = np.array(range(len(set01_pung_1000_Latencies))) / float(len(set01_pung_1000_Latencies))

    # Draw plots.

    _, ax = plt.subplots(figsize=(9, 5))
    
    ax.plot(set01_pung_0500_Latencies, set01_pung_0500_CDF, label='pung, 500 clients (tc off, no failures)')
    ax.plot(set01_pung_1000_Latencies, set01_pung_1000_CDF, label='pung, 1000 clients (tc off, no failures)')
    ax.plot(set02_zeno_1000_Latencies, set02_zeno_1000_CDF, label='zeno, 1000 clients (tc on, no failures)')
    ax.plot(set02_zeno_0500_Latencies, set02_zeno_0500_CDF, label='zeno, 500 clients (tc on, no failures)')
    ax.plot(set01_zeno_0500_Latencies, set01_zeno_0500_CDF, label='zeno, 500 clients (tc off, no failures)')
    ax.plot(set04_zeno_1000_Latencies, set04_zeno_1000_CDF, label='zeno, 1000 clients (tc on, mix failure)')
    ax.plot(set01_zeno_1000_Latencies, set01_zeno_1000_CDF, label='zeno, 1000 clients (tc off, no failures)')
    ax.plot(set03_zeno_1000_Latencies, set03_zeno_1000_CDF, label='zeno, 1000 clients (tc off, mix failure)')
    ax.plot(set04_zeno_0500_Latencies, set04_zeno_0500_CDF, label='zeno, 500 clients (tc on, mix failure)')
    ax.plot(set03_zeno_0500_Latencies, set03_zeno_0500_CDF, label='zeno, 500 clients (tc off, mix failure)')

    ax.xaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_ylim([0.0, 1.0])
    ax.set_yticks((0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0))

    # Add a legend.
    ax.legend(loc='lower right')

    plt.tight_layout()
    plt.xlabel("End-to-end transmission latency (seconds)")
    plt.ylabel("Fraction of messages transmitted")

    plt.savefig(os.path.join(sys.argv[1], "msg-latencies.pdf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "msg-latencies.pgf"), bbox_inches='tight')


def compileMessagesPerMix():

    global metrics

    for setting in metrics:

        for numClients in {"0500", "1000"}:

            for run in {"Run01", "Run02", "Run03"}:

                outputFilePDF = os.path.join(os.path.dirname(metrics[setting]["zeno"][numClients]["MessagesPerMix"][run]), "msgs-per-mix_first-to-last-round.pdf")
                outputFilePGF = os.path.join(os.path.dirname(metrics[setting]["zeno"][numClients]["MessagesPerMix"][run]), "msgs-per-mix_first-to-last-round.pgf")

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
                    plt.plot(msgCounts, "-", label=labels[idx], markersize=2.0, color=np.random.rand(3,))

                boxOfPlot = ax.get_position()
                ax.set_position([boxOfPlot.x0, boxOfPlot.y0, (boxOfPlot.width * 0.8), boxOfPlot.height])
                ax.legend(loc='center left', bbox_to_anchor=(1, 0.5), fontsize='small')

                plt.grid()
                plt.tight_layout()

                plt.xlabel("Round number")
                plt.ylabel("Messages in all pools (count)")

                plt.savefig(outputFilePDF, bbox_inches='tight')
                plt.savefig(outputFilePGF, bbox_inches='tight')

# Create all figures.

# Build bandwidth figures.
# compileTrafficClients()
compileTrafficServers()

# Build load usage figures.
# compileLoadCPUClients()
# compileLoadCPUServers()
# compileLoadMemClients()
# compileLoadMemServers()

# Build message latencies figure.
# compileLatencies()

# Build figures describing the number of
# messages in each mix server over rounds.
# compileMessagesPerMix()
