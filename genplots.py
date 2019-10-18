#!/usr/bin/env python3

import sys
import os
import csv
import numpy as np
import matplotlib

from matplotlib import pyplot as plt
from matplotlib import patches as mpatches
from pylab import setp

matplotlib.use("pgf")
plt.rcParams.update({
    "font.family": "serif",
    "text.usetex": True,
    "pgf.rcfonts": False,
    "pgf.preamble": [
        "\\usepackage{units}",
        "\\usepackage{metalogo}",
        "\\usepackage{amsmath}",
        "\\usepackage{amssymb}",
        "\\usepackage{amsthm}",
        "\\usepackage{paratype}",
        "\\usepackage{FiraMono}",
    ]
})

# Load all measurement files.
metrics = {
    "01": {
        "zeno": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "zeno", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
        "pung": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "pung", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
        "vuvuzela": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "01_tc-off_proc-off", "vuvuzela", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
    },
    "02a": {
        "zeno": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "zeno", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
        "pung": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "pung", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
        "vuvuzela": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc1-on_proc-off", "vuvuzela", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
    },
    "02b": {
        "zeno": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "zeno", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
        "pung": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "pung", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
        "vuvuzela": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "TotalExpTime": os.path.join(sys.argv[1], "02_tc2-on_proc-off", "vuvuzela", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
    },
    "03": {
        "zeno": {
            "1000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-1000", "total-experiment-times_seconds.data"),
            },
            "2000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-2000", "total-experiment-times_seconds.data"),
            },
            "3000": {
                "Bandwidth": {
                    "Clients": {
                        "Avg": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_clients.data"),
                    },
                    "Servers": {
                        "Avg": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "traffic-volume_mebibytes_highest-at-end-of-time-window_servers.data"),
                    },
                },
                "Load": {
                    "Clients": {
                        "CPU": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_clients.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "memory_megabytes-used_all-values-in-time-window_clients.data"),
                    },
                    "Servers": {
                        "CPU": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "cpu_percentage-busy_all-values-in-time-window_servers.data"),
                        "Mem": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "memory_gigabytes-used_all-values-in-time-window_servers.data"),
                    },
                },
                "Latencies": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "transmission-latencies_seconds_all-values-in-time-window.data"),
                "MessagesPerMix": {
                    "Run01": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "run-01", "message-count-per-server_first-to-last-round.data"),
                    "Run02": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "run-02", "message-count-per-server_first-to-last-round.data"),
                    "Run03": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "run-03", "message-count-per-server_first-to-last-round.data"),
                },
                "TotalExpTime": os.path.join(sys.argv[1], "03_tc3-on_proc-on", "zeno", "clients-3000", "total-experiment-times_seconds.data"),
            },
        },
    },        
}


def compileTrafficClients():

    global metrics

    # Ingest and prepare data.

    set02a_zeno_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_zeno_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set02a_vuvuzela_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set02a_pung_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["pung"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_pung_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    bandwidthAvg = [set02a_zeno_1000_Bandwidth_Clients_Avg, set02a_vuvuzela_1000_Bandwidth_Clients_Avg, set02a_pung_1000_Bandwidth_Clients_Avg]

    # Draw plots.

    width = 1.0
    y_max = np.ceil((max(bandwidthAvg) + 10.0))

    _, ax = plt.subplots()

    # Draw all bars.

    ax.bar(1, set02a_zeno_1000_Bandwidth_Clients_Avg, width, label='zeno', edgecolor='black', color='gold', hatch='/')
    ax.bar(2, set02a_vuvuzela_1000_Bandwidth_Clients_Avg, width, label='vuvuzela', edgecolor='black', color='olive', hatch='+')
    ax.bar(3, set02a_pung_1000_Bandwidth_Clients_Avg, width, label='pung', edgecolor='black', color='steelblue', hatch='x')

    labels = ["%.2f" % avg for avg in bandwidthAvg]

    for bar, label in zip(ax.patches, labels):
        ax.text((bar.get_x() + (bar.get_width() / 2)), (bar.get_height() * 1.075), label, ha='center', va='bottom')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 4])
    ax.set_ylim([0, y_max])
    ax.set_xticks([2])
    ax.set_xticklabels(['1,000 clients'])

    # Add a legend.
    ax.legend(loc='upper left')

    plt.yscale('symlog')
    plt.tight_layout()

    ax.set_title("Average Highest Traffic Volume on Clients (high delay, no failures)")
    plt.xlabel("Number of clients")
    plt.ylabel("Traffic volume (MiB) [log.]")

    plt.savefig(os.path.join(sys.argv[1], "traffic-volume_clients.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "traffic-volume_clients.pdf"), bbox_inches='tight')


def compileTrafficServers():

    global metrics

    # Ingest and prepare data.

    set02a_zeno_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_zeno_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02a_zeno_1000_Bandwidth_Servers_AvgAll = set02a_zeno_1000_Bandwidth_Servers_Avg * 21.0

    set02a_vuvuzela_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02a_vuvuzela_1000_Bandwidth_Servers_AvgAll = set02a_vuvuzela_1000_Bandwidth_Servers_Avg * 4.0

    set02a_pung_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["pung"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_pung_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())

    bandwidthAvg = [set02a_zeno_1000_Bandwidth_Servers_AvgAll,
                    set02a_vuvuzela_1000_Bandwidth_Servers_AvgAll,
                    set02a_pung_1000_Bandwidth_Servers_Avg]

    # Draw plots.

    width = 1.0
    barWidth = (1.0 / 4.0)
    y_max = np.ceil((max(bandwidthAvg) + 5000.0))

    _, ax = plt.subplots()

    # Draw all bars and corresponding average lines.

    ax.bar(1, set02a_zeno_1000_Bandwidth_Servers_AvgAll, width, label='zeno', edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set02a_zeno_1000_Bandwidth_Servers_Avg, xmin=(0.5 * barWidth), xmax=(1.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(2, set02a_vuvuzela_1000_Bandwidth_Servers_AvgAll, width, label='vuvuzela', edgecolor='black', color='olive', hatch='+')
    plt.axhline(y=set02a_vuvuzela_1000_Bandwidth_Servers_Avg, xmin=(1.5 * barWidth), xmax=(2.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(3, set02a_pung_1000_Bandwidth_Servers_Avg, width, label='pung', edgecolor='black', color='steelblue', hatch='x')
    
    labels = ['{:,.0f}'.format(avg) for avg in bandwidthAvg]

    for bar, label in zip(ax.patches, labels):
        ax.text((bar.get_x() + (bar.get_width() / 2)), (bar.get_height() + 200), label, ha='center', va='bottom')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 4])
    ax.set_ylim([0, y_max])
    ax.set_xticks([2])
    ax.set_xticklabels(['1,000 clients'])
    ax.get_yaxis().set_major_formatter(matplotlib.ticker.FuncFormatter(lambda x, p: format(int(x), ',')))

    # Add a legend.
    ax.legend(loc='upper left')

    ax.set_title("Average Highest Traffic Volume on Servers (medium delay, no failures)")
    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Traffic volume (MiB)")

    plt.savefig(os.path.join(sys.argv[1], "traffic-volume_servers.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "traffic-volume_servers.pdf"), bbox_inches='tight')


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
    ax.set_xticks((3, 9))
    ax.set_yticks([0, 10, 20, 30, 40, 50, 50])
    ax.set_xticklabels(('500 clients', '1,000 clients'))

    # Add a legend.
    ax.legend([set01_zeno01['boxes'][0], set02_zeno01['boxes'][0], set03_zeno01['boxes'][0],
        set04_zeno01['boxes'][0], set01_pung01['boxes'][0]], ['zeno (tc off, no failures)',
        'zeno (tc on, no failures)', 'zeno (tc off, mix failure)', 'zeno (tc on, mix failure)',
        'pung (tc off, no failures)'], loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Busy CPU (percentage)")

    plt.savefig(os.path.join(sys.argv[1], "cpu-busy_clients.pgf"), bbox_inches='tight')


def compileLoadMemClients():

    global metrics

    # Ingest data.

    set01_zeno_0500_Load_Mem = []
    with open(metrics["01"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_0500_Load_Mem.append(float(item))

    set01_pung_0500_Load_Mem = []
    with open(metrics["01"]["pung"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_0500_Load_Mem.append(float(item))

    set01_zeno_1000_Load_Mem = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_Mem.append(float(item))

    set01_pung_1000_Load_Mem = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_Mem.append(float(item))

    set02_zeno_0500_Load_Mem = []
    with open(metrics["02"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_0500_Load_Mem.append(float(item))

    set02_zeno_1000_Load_Mem = []
    with open(metrics["02"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_1000_Load_Mem.append(float(item))

    set03_zeno_0500_Load_Mem = []
    with open(metrics["03"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_0500_Load_Mem.append(float(item))

    set03_zeno_1000_Load_Mem = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_Mem.append(float(item))

    set04_zeno_0500_Load_Mem = []
    with open(metrics["04"]["zeno"]["0500"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_0500_Load_Mem.append(float(item))

    set04_zeno_1000_Load_Mem = []
    with open(metrics["04"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_1000_Load_Mem.append(float(item))
    
    # Draw plots.

    width = 0.9

    _, ax = plt.subplots()

    set01_zeno01 = ax.boxplot(set01_zeno_0500_Load_Mem, positions=[1], widths=width, patch_artist=True, whis='range')
    set02_zeno01 = ax.boxplot(set02_zeno_0500_Load_Mem, positions=[2], widths=width, patch_artist=True, whis='range')
    set03_zeno01 = ax.boxplot(set03_zeno_0500_Load_Mem, positions=[3], widths=width, patch_artist=True, whis='range')
    set04_zeno01 = ax.boxplot(set04_zeno_0500_Load_Mem, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_pung01 = ax.boxplot(set01_pung_0500_Load_Mem, positions=[5], widths=width, patch_artist=True, whis='range')
    set01_zeno03 = ax.boxplot(set01_zeno_1000_Load_Mem, positions=[7], widths=width, patch_artist=True, whis='range')
    set02_zeno03 = ax.boxplot(set02_zeno_1000_Load_Mem, positions=[8], widths=width, patch_artist=True, whis='range')
    set03_zeno03 = ax.boxplot(set03_zeno_1000_Load_Mem, positions=[9], widths=width, patch_artist=True, whis='range')
    set04_zeno03 = ax.boxplot(set04_zeno_1000_Load_Mem, positions=[10], widths=width, patch_artist=True, whis='range')
    set01_pung03 = ax.boxplot(set01_pung_1000_Load_Mem, positions=[11], widths=width, patch_artist=True, whis='range')

    # Log values for text mention.

    print("\nClients:\n")
    for scenData in [set01_zeno01, set02_zeno01, set03_zeno01, set04_zeno01, set01_pung01, set01_zeno03, set02_zeno03, set03_zeno03, set04_zeno03, set01_pung03]:

        for whis in scenData['whiskers']:
            print("whis=", whis.get_ydata()[1])
        
        for med in scenData['medians']:
            print(" med=", med.get_ydata()[1])
        
        print("")
    print("")
    
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
    ax.set_xticks((3, 9))
    ax.set_ylim([0, 600])
    ax.set_yticks([0, 100, 200, 300, 400, 500, 600, 700])
    ax.set_xticklabels(('500 clients', '1,000 clients'))

    # Add a legend.
    ax.legend([set01_zeno01['boxes'][0], set02_zeno01['boxes'][0], set03_zeno01['boxes'][0],
        set04_zeno01['boxes'][0], set01_pung01['boxes'][0]], ['zeno (tc off, no failures)',
        'zeno (tc on, no failures)', 'zeno (tc off, mix failure)', 'zeno (tc on, mix failure)',
        'pung (tc off, no failures)'], loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Used memory (MB)")

    plt.savefig(os.path.join(sys.argv[1], "memory-used_clients.pgf"), bbox_inches='tight')


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
    ax.set_xticks((3, 9))
    ax.set_yticks([0, 10, 20, 30, 40, 50])
    ax.set_xticklabels(('500 clients', '1,000 clients'))

    # Add a legend.
    ax.legend([set01_zeno01['boxes'][0], set02_zeno01['boxes'][0], set03_zeno01['boxes'][0],
        set04_zeno01['boxes'][0], set01_pung01['boxes'][0]], ['zeno (tc off, no failures)',
        'zeno (tc on, no failures)', 'zeno (tc off, mix failure)', 'zeno (tc on, mix failure)',
        'pung (tc off, no failures)'], loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Busy CPU (percentage)")

    plt.savefig(os.path.join(sys.argv[1], "cpu-busy_servers.pgf"), bbox_inches='tight')


def compileLoadMemServers():

    global metrics

    # Ingest data.

    set01_zeno_0500_Load_Mem = []
    with open(metrics["01"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_0500_Load_Mem.append(float(item))

    set01_pung_0500_Load_Mem = []
    with open(metrics["01"]["pung"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_0500_Load_Mem.append(float(item))

    set01_zeno_1000_Load_Mem = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_Mem.append(float(item))

    set01_pung_1000_Load_Mem = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_Mem.append(float(item))

    set02_zeno_0500_Load_Mem = []
    with open(metrics["02"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_0500_Load_Mem.append(float(item))

    set02_zeno_1000_Load_Mem = []
    with open(metrics["02"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02_zeno_1000_Load_Mem.append(float(item))

    set03_zeno_0500_Load_Mem = []
    with open(metrics["03"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_0500_Load_Mem.append(float(item))

    set03_zeno_1000_Load_Mem = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_Mem.append(float(item))

    set04_zeno_0500_Load_Mem = []
    with open(metrics["04"]["zeno"]["0500"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_0500_Load_Mem.append(float(item))

    set04_zeno_1000_Load_Mem = []
    with open(metrics["04"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set04_zeno_1000_Load_Mem.append(float(item))

    # Draw plots.

    width = 0.9

    _, ax = plt.subplots()

    set01_zeno01 = ax.boxplot(set01_zeno_0500_Load_Mem, positions=[1], widths=width, patch_artist=True, whis='range')
    set02_zeno01 = ax.boxplot(set02_zeno_0500_Load_Mem, positions=[2], widths=width, patch_artist=True, whis='range')
    set03_zeno01 = ax.boxplot(set03_zeno_0500_Load_Mem, positions=[3], widths=width, patch_artist=True, whis='range')
    set04_zeno01 = ax.boxplot(set04_zeno_0500_Load_Mem, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_pung01 = ax.boxplot(set01_pung_0500_Load_Mem, positions=[5], widths=width, patch_artist=True, whis='range')
    set01_zeno03 = ax.boxplot(set01_zeno_1000_Load_Mem, positions=[7], widths=width, patch_artist=True, whis='range')
    set02_zeno03 = ax.boxplot(set02_zeno_1000_Load_Mem, positions=[8], widths=width, patch_artist=True, whis='range')
    set03_zeno03 = ax.boxplot(set03_zeno_1000_Load_Mem, positions=[9], widths=width, patch_artist=True, whis='range')
    set04_zeno03 = ax.boxplot(set04_zeno_1000_Load_Mem, positions=[10], widths=width, patch_artist=True, whis='range')
    set01_pung03 = ax.boxplot(set01_pung_1000_Load_Mem, positions=[11], widths=width, patch_artist=True, whis='range')

    # Log values for text mention.
    
    print("Servers:\n")
    for scenData in [set01_zeno01, set02_zeno01, set03_zeno01, set04_zeno01, set01_pung01, set01_zeno03, set02_zeno03, set03_zeno03, set04_zeno03, set01_pung03]:

        for whis in scenData['whiskers']:
            print("whis=", whis.get_ydata()[1])
        
        for med in scenData['medians']:
            print(" med=", med.get_ydata()[1])
        
        print("")
    
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
    ax.set_xticks((3, 9))
    ax.set_xticklabels(('500 clients', '1,000 clients'))
    ax.set_yticks([0, 2, 4, 6, 8, 10, 12, 14, 16, 18])

    # Add a legend.
    ax.legend([set01_zeno01['boxes'][0], set02_zeno01['boxes'][0], set03_zeno01['boxes'][0],
        set04_zeno01['boxes'][0], set01_pung01['boxes'][0]], ['zeno (tc off, no failures)',
        'zeno (tc on, no failures)', 'zeno (tc off, mix failure)', 'zeno (tc on, mix failure)',
        'pung (tc off, no failures)'], loc='upper left')

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Used memory (GB)")

    plt.savefig(os.path.join(sys.argv[1], "memory-used_servers.pgf"), bbox_inches='tight')


def compileLatencies():

    global metrics

    # Ingest data.

    x_max = 0.0

    set02a_zeno_1000_Latencies = []
    with open(metrics["02a"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) > x_max:
                    x_max = float(item)
                set02a_zeno_1000_Latencies.append(float(item))
    
    set02a_vuvuzela_1000_Latencies = []
    with open(metrics["02a"]["vuvuzela"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) > x_max:
                    x_max = float(item)
                set02a_vuvuzela_1000_Latencies.append(float(item))

    set02a_pung_1000_Latencies = []
    with open(metrics["02a"]["pung"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) > x_max:
                    x_max = float(item)
                set02a_pung_1000_Latencies.append(float(item))
    

    print("x_max=", x_max)
    x_max = np.ceil(x_max) + 10.0
    print("x_max=", x_max)

    # Prepare CDF arrays.

    set02a_zeno_1000_Latencies = np.sort(set02a_zeno_1000_Latencies)
    set02a_zeno_1000_CDF = np.array(range(len(set02a_zeno_1000_Latencies))) / float(len(set02a_zeno_1000_Latencies))

    set02a_vuvuzela_1000_Latencies = np.sort(set02a_vuvuzela_1000_Latencies)
    set02a_vuvuzela_1000_CDF = np.array(range(len(set02a_vuvuzela_1000_Latencies))) / float(len(set02a_vuvuzela_1000_Latencies))

    set02a_pung_1000_Latencies = np.sort(set02a_pung_1000_Latencies)
    set02a_pung_1000_CDF = np.array(range(len(set02a_pung_1000_Latencies))) / float(len(set02a_pung_1000_Latencies))

    # Draw plots.

    _, ax = plt.subplots()
    
    ax.plot(set02a_pung_1000_Latencies, set02a_pung_1000_CDF, label='pung')
    ax.plot(set02a_zeno_1000_Latencies, set02a_zeno_1000_CDF, label='zeno')
    ax.plot(set02a_vuvuzela_1000_Latencies, set02a_vuvuzela_1000_CDF, label='vuvuzela')

    ax.xaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0.0, x_max])
    ax.set_ylim([0.0, 1.0])
    ax.set_yticks((0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0))

    # Add a legend.
    ax.legend(loc='lower right')

    ax.set_title("CDFs of End-to-End Transmission Latencies on Clients (delay, no failures)")
    plt.tight_layout()
    plt.xlabel("End-to-end transmission latency (seconds)")
    plt.ylabel("Fraction of messages transmitted")

    plt.savefig(os.path.join(sys.argv[1], "e2e-transmission-latencies.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "e2e-transmission-latencies.pdf"), bbox_inches='tight')


def compileMessagesPerMix():

    global metrics

    for setting in metrics:

        for numClients in {"0500", "1000"}:

            for run in {"Run01", "Run02", "Run03"}:

                outputFile = os.path.join(os.path.dirname(metrics[setting]["zeno"][numClients]["MessagesPerMix"][run]), "message-count-per-server_first-to-last-round.pgf")

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

                ax.get_yaxis().set_major_formatter(matplotlib.ticker.FuncFormatter(lambda x, p: format(int(x), ',')))

                plt.grid()
                plt.tight_layout()

                plt.xlabel("Round number")
                plt.ylabel("Message count")

                plt.savefig(outputFile, bbox_inches='tight')

# Create all figures.

# Build bandwidth figures.
compileTrafficClients()
compileTrafficServers()

# Build load usage figures.
# compileLoadCPUClients()
# compileLoadCPUServers()
# compileLoadMemClients()
# compileLoadMemServers()

# Build message latencies figure.
compileLatencies()

# Build figures describing the number of
# messages in each mix server over rounds.
# compileMessagesPerMix()
