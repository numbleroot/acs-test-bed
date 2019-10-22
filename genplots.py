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

    # 1000 clients.

    set01_zeno_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_zeno_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set01_vuvuzela_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["vuvuzela"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_vuvuzela_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set01_pung_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["pung"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_pung_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set02a_zeno_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_zeno_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set02a_vuvuzela_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set02a_pung_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["pung"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_pung_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set02b_zeno_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_zeno_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set02b_vuvuzela_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["vuvuzela"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_vuvuzela_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())

    set02b_pung_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["pung"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_pung_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set03_zeno_1000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["03"]["zeno"]["1000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set03_zeno_1000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    # 2000 clients.

    set01_zeno_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["zeno"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_zeno_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set01_vuvuzela_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["vuvuzela"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_vuvuzela_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set01_pung_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["pung"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_pung_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set02a_zeno_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["zeno"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_zeno_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02a_vuvuzela_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02a_pung_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["pung"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_pung_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set02b_zeno_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["zeno"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_zeno_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02b_vuvuzela_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["vuvuzela"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_vuvuzela_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02b_pung_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["pung"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_pung_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set03_zeno_2000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["03"]["zeno"]["2000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set03_zeno_2000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    # 3000 clients.

    set01_zeno_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["zeno"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_zeno_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set01_vuvuzela_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["vuvuzela"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_vuvuzela_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set01_pung_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["01"]["pung"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set01_pung_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set02a_zeno_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["zeno"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_zeno_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02a_vuvuzela_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02a_pung_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02a"]["pung"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02a_pung_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set02b_zeno_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["zeno"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_zeno_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02b_vuvuzela_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["vuvuzela"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_vuvuzela_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())
    
    set02b_pung_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["02b"]["pung"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set02b_pung_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    set03_zeno_3000_Bandwidth_Clients_Avg = 0.0
    with open(metrics["03"]["zeno"]["3000"]["Bandwidth"]["Clients"]["Avg"], 'r') as dataFile:
        set03_zeno_3000_Bandwidth_Clients_Avg = float(dataFile.read().strip())


    bandwidthAvg = [
        set01_zeno_1000_Bandwidth_Clients_Avg,
        set02a_zeno_1000_Bandwidth_Clients_Avg,
        set02b_zeno_1000_Bandwidth_Clients_Avg,
        set03_zeno_1000_Bandwidth_Clients_Avg,
        set01_vuvuzela_1000_Bandwidth_Clients_Avg,
        set02a_vuvuzela_1000_Bandwidth_Clients_Avg,
        set02b_vuvuzela_1000_Bandwidth_Clients_Avg,
        set01_pung_1000_Bandwidth_Clients_Avg,
        set02a_pung_1000_Bandwidth_Clients_Avg,
        set02b_pung_1000_Bandwidth_Clients_Avg,

        set01_zeno_2000_Bandwidth_Clients_Avg,
        set02a_zeno_2000_Bandwidth_Clients_Avg,
        set02b_zeno_2000_Bandwidth_Clients_Avg,
        set03_zeno_2000_Bandwidth_Clients_Avg,
        set01_vuvuzela_2000_Bandwidth_Clients_Avg,
        set02a_vuvuzela_2000_Bandwidth_Clients_Avg,
        set02b_vuvuzela_2000_Bandwidth_Clients_Avg,
        set01_pung_2000_Bandwidth_Clients_Avg,
        set02a_pung_2000_Bandwidth_Clients_Avg,
        set02b_pung_2000_Bandwidth_Clients_Avg,

        set01_zeno_3000_Bandwidth_Clients_Avg,
        set02a_zeno_3000_Bandwidth_Clients_Avg,
        set02b_zeno_3000_Bandwidth_Clients_Avg,
        set03_zeno_3000_Bandwidth_Clients_Avg,
        set01_vuvuzela_3000_Bandwidth_Clients_Avg,
        set02a_vuvuzela_3000_Bandwidth_Clients_Avg,
        set02b_vuvuzela_3000_Bandwidth_Clients_Avg,
        set01_pung_3000_Bandwidth_Clients_Avg,
        set02a_pung_3000_Bandwidth_Clients_Avg,
        set02b_pung_3000_Bandwidth_Clients_Avg
    ]

    # Draw plots.

    width = 1.0
    y_max = np.ceil((max(bandwidthAvg) + 10.0))

    _, ax = plt.subplots(figsize=(14, 5))

    # Draw all bars.

    ax.bar(1, set01_zeno_1000_Bandwidth_Clients_Avg, width, label='zeno (no impediments)', edgecolor='black', color='gold', hatch='/')
    ax.bar(2, set02a_zeno_1000_Bandwidth_Clients_Avg, width, label='zeno (high delay, no failures)', edgecolor='black', color='gold', hatch='//')
    ax.bar(3, set02b_zeno_1000_Bandwidth_Clients_Avg, width, label='zeno (high loss, no failures)', edgecolor='black', color='gold', hatch='+')
    ax.bar(4, set03_zeno_1000_Bandwidth_Clients_Avg, width, label='zeno (high network troubles, failures)', edgecolor='black', color='gold', hatch='.')
    ax.bar(5, set01_vuvuzela_1000_Bandwidth_Clients_Avg, width, label='vuvuzela (no impediments)', edgecolor='black', color='darkseagreen', hatch='\\')
    ax.bar(6, set02a_vuvuzela_1000_Bandwidth_Clients_Avg, width, label='vuvuzela (high delay, no failures)', edgecolor='black', color='darkseagreen', hatch='\\\\')
    ax.bar(7, set02b_vuvuzela_1000_Bandwidth_Clients_Avg, width, label='vuvuzela (high loss, no failures)', edgecolor='black', color='darkseagreen', hatch='-')
    ax.bar(8, set01_pung_1000_Bandwidth_Clients_Avg, width, label='pung (no impediments)', edgecolor='black', color='steelblue', hatch='x')
    ax.bar(9, set02a_pung_1000_Bandwidth_Clients_Avg, width, label='pung (high delay, no failures)', edgecolor='black', color='steelblue', hatch='o')
    ax.bar(10, set02b_pung_1000_Bandwidth_Clients_Avg, width, label='pung (high loss, no failures)', edgecolor='black', color='steelblue', hatch='*')

    ax.bar(12, set01_zeno_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='/')
    ax.bar(13, set02a_zeno_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='//')
    ax.bar(14, set02b_zeno_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='+')
    ax.bar(15, set03_zeno_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='.')
    ax.bar(16, set01_vuvuzela_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='darkseagreen', hatch='\\')
    ax.bar(17, set02a_vuvuzela_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='darkseagreen', hatch='\\\\')
    ax.bar(18, set02b_vuvuzela_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='darkseagreen', hatch='-')
    ax.bar(19, set01_pung_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='x')
    ax.bar(20, set02a_pung_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='o')
    ax.bar(21, set02b_pung_2000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='*')

    ax.bar(23, set01_zeno_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='/')
    ax.bar(24, set02a_zeno_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='//')
    ax.bar(25, set02b_zeno_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='+')
    ax.bar(26, set03_zeno_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='gold', hatch='.')
    ax.bar(27, set01_vuvuzela_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='darkseagreen', hatch='\\')
    ax.bar(28, set02a_vuvuzela_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='darkseagreen', hatch='\\\\')
    ax.bar(29, set02b_vuvuzela_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='darkseagreen', hatch='-')
    ax.bar(30, set01_pung_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='x')
    ax.bar(31, set02a_pung_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='o')
    ax.bar(32, set02b_pung_3000_Bandwidth_Clients_Avg, width, edgecolor='black', color='steelblue', hatch='*')

    labels = ["%.2f" % avg for avg in bandwidthAvg]

    for bar, label in zip(ax.patches, labels):
        ax.text((bar.get_x() + (bar.get_width() / 2)), ((bar.get_height() * 1.07) + 0.05), label, ha='center', va='bottom')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 33])
    ax.set_ylim([0, y_max])
    ax.set_xticks((5.5, 16.5, 27.5))
    ax.set_xticklabels(('1,000 clients', '2,000 clients', '3,000 clients'))

    # Add a legend.
    ax.legend(loc='upper left')

    plt.yscale('symlog')
    plt.tight_layout()

    ax.set_title("Average Highest Traffic Volume on Clients")
    plt.xlabel("Number of clients")
    plt.ylabel("Traffic volume (MiB) [log.]")

    # plt.savefig(os.path.join(sys.argv[1], "traffic-volume_clients.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "traffic-volume_clients.pdf"), bbox_inches='tight')



def compileTrafficServers():

    global metrics

    # Ingest and prepare data.

    # 1000 clients.

    set01_zeno_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_zeno_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set01_zeno_1000_Bandwidth_Servers_AvgAll = set01_zeno_1000_Bandwidth_Servers_Avg * 21.0

    set01_vuvuzela_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["vuvuzela"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_vuvuzela_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set01_vuvuzela_1000_Bandwidth_Servers_AvgAll = set01_vuvuzela_1000_Bandwidth_Servers_Avg * 4.0

    set01_pung_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["pung"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_pung_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


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


    set02b_zeno_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_zeno_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02b_zeno_1000_Bandwidth_Servers_AvgAll = set02b_zeno_1000_Bandwidth_Servers_Avg * 21.0

    set02b_vuvuzela_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["vuvuzela"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_vuvuzela_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02b_vuvuzela_1000_Bandwidth_Servers_AvgAll = set02b_vuvuzela_1000_Bandwidth_Servers_Avg * 4.0

    set02b_pung_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["pung"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_pung_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set03_zeno_1000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["03"]["zeno"]["1000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set03_zeno_1000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set03_zeno_1000_Bandwidth_Servers_AvgAll = set03_zeno_1000_Bandwidth_Servers_Avg * 21.0


    # 2000 clients.

    set01_zeno_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["zeno"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_zeno_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set01_zeno_2000_Bandwidth_Servers_AvgAll = set01_zeno_2000_Bandwidth_Servers_Avg * 21.0
    
    set01_vuvuzela_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["vuvuzela"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_vuvuzela_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set01_vuvuzela_2000_Bandwidth_Servers_AvgAll = set01_vuvuzela_2000_Bandwidth_Servers_Avg * 4.0
    
    set01_pung_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["pung"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_pung_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set02a_zeno_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["zeno"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_zeno_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02a_zeno_2000_Bandwidth_Servers_AvgAll = set02a_zeno_2000_Bandwidth_Servers_Avg * 21.0
    
    set02a_vuvuzela_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02a_vuvuzela_2000_Bandwidth_Servers_AvgAll = set02a_vuvuzela_2000_Bandwidth_Servers_Avg * 4.0
    
    set02a_pung_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["pung"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_pung_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set02b_zeno_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["zeno"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_zeno_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02b_zeno_2000_Bandwidth_Servers_AvgAll = set02b_zeno_2000_Bandwidth_Servers_Avg * 21.0
    
    set02b_vuvuzela_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["vuvuzela"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_vuvuzela_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02b_vuvuzela_2000_Bandwidth_Servers_AvgAll = set02b_vuvuzela_2000_Bandwidth_Servers_Avg * 4.0
    
    set02b_pung_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["pung"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_pung_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set03_zeno_2000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["03"]["zeno"]["2000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set03_zeno_2000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set03_zeno_2000_Bandwidth_Servers_AvgAll = set03_zeno_2000_Bandwidth_Servers_Avg * 21.0


    # 3000 clients.

    set01_zeno_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["zeno"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_zeno_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set01_zeno_3000_Bandwidth_Servers_AvgAll = set01_zeno_3000_Bandwidth_Servers_Avg * 21.0
    
    set01_vuvuzela_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["vuvuzela"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_vuvuzela_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set01_vuvuzela_3000_Bandwidth_Servers_AvgAll = set01_vuvuzela_3000_Bandwidth_Servers_Avg * 4.0
    
    set01_pung_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["01"]["pung"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set01_pung_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set02a_zeno_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["zeno"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_zeno_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02a_zeno_3000_Bandwidth_Servers_AvgAll = set02a_zeno_3000_Bandwidth_Servers_Avg * 21.0
    
    set02a_vuvuzela_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["vuvuzela"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_vuvuzela_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02a_vuvuzela_3000_Bandwidth_Servers_AvgAll = set02a_vuvuzela_3000_Bandwidth_Servers_Avg * 4.0
    
    set02a_pung_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02a"]["pung"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02a_pung_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set02b_zeno_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["zeno"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_zeno_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02b_zeno_3000_Bandwidth_Servers_AvgAll = set02b_zeno_3000_Bandwidth_Servers_Avg * 21.0
    
    set02b_vuvuzela_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["vuvuzela"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_vuvuzela_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set02b_vuvuzela_3000_Bandwidth_Servers_AvgAll = set02b_vuvuzela_3000_Bandwidth_Servers_Avg * 4.0
    
    set02b_pung_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["02b"]["pung"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set02b_pung_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())


    set03_zeno_3000_Bandwidth_Servers_Avg = 0.0
    with open(metrics["03"]["zeno"]["3000"]["Bandwidth"]["Servers"]["Avg"], 'r') as dataFile:
        set03_zeno_3000_Bandwidth_Servers_Avg = float(dataFile.read().strip())
    set03_zeno_3000_Bandwidth_Servers_AvgAll = set03_zeno_3000_Bandwidth_Servers_Avg * 21.0


    bandwidthAvg = [
        set01_zeno_1000_Bandwidth_Servers_AvgAll,
        set02a_zeno_1000_Bandwidth_Servers_AvgAll,
        set02b_zeno_1000_Bandwidth_Servers_AvgAll,
        set03_zeno_1000_Bandwidth_Servers_AvgAll,
        set01_vuvuzela_1000_Bandwidth_Servers_AvgAll,
        set02a_vuvuzela_1000_Bandwidth_Servers_AvgAll,
        set02b_vuvuzela_1000_Bandwidth_Servers_AvgAll,
        set01_pung_1000_Bandwidth_Servers_Avg,
        set02a_pung_1000_Bandwidth_Servers_Avg,
        set02b_pung_1000_Bandwidth_Servers_Avg,

        set01_zeno_2000_Bandwidth_Servers_AvgAll,
        set02a_zeno_2000_Bandwidth_Servers_AvgAll,
        set02b_zeno_2000_Bandwidth_Servers_AvgAll,
        set03_zeno_2000_Bandwidth_Servers_AvgAll,
        set01_vuvuzela_2000_Bandwidth_Servers_AvgAll,
        set02a_vuvuzela_2000_Bandwidth_Servers_AvgAll,
        set02b_vuvuzela_2000_Bandwidth_Servers_AvgAll,
        set01_pung_2000_Bandwidth_Servers_Avg,
        set02a_pung_2000_Bandwidth_Servers_Avg,
        set02b_pung_2000_Bandwidth_Servers_Avg,

        set01_zeno_3000_Bandwidth_Servers_AvgAll,
        set02a_zeno_3000_Bandwidth_Servers_AvgAll,
        set02b_zeno_3000_Bandwidth_Servers_AvgAll,
        set03_zeno_3000_Bandwidth_Servers_AvgAll,
        set01_vuvuzela_3000_Bandwidth_Servers_AvgAll,
        set02a_vuvuzela_3000_Bandwidth_Servers_AvgAll,
        set02b_vuvuzela_3000_Bandwidth_Servers_AvgAll,
        set01_pung_3000_Bandwidth_Servers_Avg,
        set02a_pung_3000_Bandwidth_Servers_Avg,
        set02b_pung_3000_Bandwidth_Servers_Avg
    ]

    # Draw plots.

    width = 1.0
    barWidth = (1.0 / 33.0)
    y_max = np.ceil((max(bandwidthAvg) + 5000.0))

    _, ax = plt.subplots(figsize=(14, 5))

    # Draw all bars and corresponding average lines.

    ax.bar(1, set01_zeno_1000_Bandwidth_Servers_AvgAll, width, label='zeno (no impediments)', edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_1000_Bandwidth_Servers_Avg, xmin=(0.5 * barWidth), xmax=(1.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(2, set02a_zeno_1000_Bandwidth_Servers_AvgAll, width, label='zeno (high delay, no failures)', edgecolor='black', color='gold', hatch='//')
    plt.axhline(y=set02a_zeno_1000_Bandwidth_Servers_Avg, xmin=(1.5 * barWidth), xmax=(2.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(3, set02b_zeno_1000_Bandwidth_Servers_AvgAll, width, label='zeno (high loss, no failures)', edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set02b_zeno_1000_Bandwidth_Servers_Avg, xmin=(2.5 * barWidth), xmax=(3.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(4, set03_zeno_1000_Bandwidth_Servers_AvgAll, width, label='zeno (high network troubles, failures)', edgecolor='black', color='gold', hatch='.')
    plt.axhline(y=set03_zeno_1000_Bandwidth_Servers_Avg, xmin=(3.5 * barWidth), xmax=(4.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    
    ax.bar(5, set01_vuvuzela_1000_Bandwidth_Servers_AvgAll, width, label='vuvuzela (no impediments)', edgecolor='black', color='darkseagreen', hatch='\\')
    plt.axhline(y=set01_vuvuzela_1000_Bandwidth_Servers_Avg, xmin=(4.5 * barWidth), xmax=(5.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(6, set02a_vuvuzela_1000_Bandwidth_Servers_AvgAll, width, label='vuvuzela (high delay, no failures)', edgecolor='black', color='darkseagreen', hatch='\\\\')
    plt.axhline(y=set02a_vuvuzela_1000_Bandwidth_Servers_Avg, xmin=(5.5 * barWidth), xmax=(6.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(7, set02b_vuvuzela_1000_Bandwidth_Servers_AvgAll, width, label='vuvuzela (high loss, no failures)', edgecolor='black', color='darkseagreen', hatch='-')
    plt.axhline(y=set02b_vuvuzela_1000_Bandwidth_Servers_Avg, xmin=(6.5 * barWidth), xmax=(7.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(8, set01_pung_1000_Bandwidth_Servers_Avg, width, label='pung (no impediments)', edgecolor='black', color='steelblue', hatch='x')
    ax.bar(9, set02a_pung_1000_Bandwidth_Servers_Avg, width, label='pung (high delay, no failures)', edgecolor='black', color='steelblue', hatch='o')
    ax.bar(10, set02b_pung_1000_Bandwidth_Servers_Avg, width, label='pung (high loss, no failures)', edgecolor='black', color='steelblue', hatch='*')

    ax.bar(12, set01_zeno_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_2000_Bandwidth_Servers_Avg, xmin=(11.5 * barWidth), xmax=(12.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(13, set02a_zeno_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='//')
    plt.axhline(y=set02a_zeno_2000_Bandwidth_Servers_Avg, xmin=(12.5 * barWidth), xmax=(13.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(14, set02b_zeno_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set02b_zeno_2000_Bandwidth_Servers_Avg, xmin=(13.5 * barWidth), xmax=(14.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(15, set03_zeno_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='.')
    plt.axhline(y=set03_zeno_2000_Bandwidth_Servers_Avg, xmin=(14.5 * barWidth), xmax=(15.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(16, set01_vuvuzela_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='darkseagreen', hatch='\\')
    plt.axhline(y=set01_vuvuzela_2000_Bandwidth_Servers_Avg, xmin=(15.5 * barWidth), xmax=(16.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(17, set02a_vuvuzela_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='darkseagreen', hatch='\\\\')
    plt.axhline(y=set02a_vuvuzela_2000_Bandwidth_Servers_Avg, xmin=(16.5 * barWidth), xmax=(17.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(18, set02b_vuvuzela_2000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='darkseagreen', hatch='-')
    plt.axhline(y=set02b_vuvuzela_2000_Bandwidth_Servers_Avg, xmin=(17.5 * barWidth), xmax=(18.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(19, set01_pung_2000_Bandwidth_Servers_Avg, width, edgecolor='black', color='steelblue', hatch='x')
    ax.bar(20, set02a_pung_2000_Bandwidth_Servers_Avg, width, edgecolor='black', color='steelblue', hatch='o')
    ax.bar(21, set02b_pung_2000_Bandwidth_Servers_Avg, width, edgecolor='black', color='steelblue', hatch='*')

    ax.bar(23, set01_zeno_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='/')
    plt.axhline(y=set01_zeno_3000_Bandwidth_Servers_Avg, xmin=(22.5 * barWidth), xmax=(23.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(24, set02a_zeno_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='//')
    plt.axhline(y=set02a_zeno_3000_Bandwidth_Servers_Avg, xmin=(23.5 * barWidth), xmax=(24.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(25, set02b_zeno_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='+')
    plt.axhline(y=set02b_zeno_3000_Bandwidth_Servers_Avg, xmin=(24.5 * barWidth), xmax=(25.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(26, set03_zeno_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='gold', hatch='.')
    plt.axhline(y=set03_zeno_3000_Bandwidth_Servers_Avg, xmin=(25.5 * barWidth), xmax=(26.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(27, set01_vuvuzela_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='darkseagreen', hatch='\\')
    plt.axhline(y=set01_vuvuzela_3000_Bandwidth_Servers_Avg, xmin=(26.5 * barWidth), xmax=(27.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(28, set02a_vuvuzela_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='darkseagreen', hatch='\\\\')
    plt.axhline(y=set02a_vuvuzela_3000_Bandwidth_Servers_Avg, xmin=(27.5 * barWidth), xmax=(28.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')
    ax.bar(29, set02b_vuvuzela_3000_Bandwidth_Servers_AvgAll, width, edgecolor='black', color='darkseagreen', hatch='-')
    plt.axhline(y=set02b_vuvuzela_3000_Bandwidth_Servers_Avg, xmin=(28.5 * barWidth), xmax=(29.5 * barWidth), linewidth=1.5, linestyle='--', color='crimson')

    ax.bar(30, set01_pung_3000_Bandwidth_Servers_Avg, width, edgecolor='black', color='steelblue', hatch='x')
    ax.bar(31, set02a_pung_3000_Bandwidth_Servers_Avg, width, edgecolor='black', color='steelblue', hatch='o')
    ax.bar(32, set02b_pung_3000_Bandwidth_Servers_Avg, width, edgecolor='black', color='steelblue', hatch='*')
    
    labels = ['{:,.0f}'.format(avg) for avg in bandwidthAvg]

    for bar, label in zip(ax.patches, labels):
        ax.text((bar.get_x() + (bar.get_width() / 2)), (bar.get_height() + 200), label, ha='center', va='bottom')

    # Show a light horizontal grid.
    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    # Limit x and y axes and configure ticks and labels.
    ax.set_xlim([0, 33])
    ax.set_ylim([0, y_max])
    ax.set_xticks((5.5, 16.5, 27.5))
    ax.set_xticklabels(('1,000 clients', '2,000 clients', '3,000 clients'))
    ax.get_yaxis().set_major_formatter(matplotlib.ticker.FuncFormatter(lambda x, p: format(int(x), ',')))

    ax.legend(loc='upper left')
    ax.set_title("Average Highest Traffic Volume on Servers")

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Traffic volume (MiB)")

    # plt.savefig(os.path.join(sys.argv[1], "traffic-volume_servers.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "traffic-volume_servers.pdf"), bbox_inches='tight')


def compileLoadCPUClients():

    global metrics

    # Ingest data.

    # 1000 clients.

    set01_zeno_1000_Load_CPU = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_CPU.append(float(item))
    
    set01_vuvuzela_1000_Load_CPU = []
    with open(metrics["01"]["vuvuzela"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_1000_Load_CPU.append(float(item))
    
    set01_pung_1000_Load_CPU = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_CPU.append(float(item))


    set02a_zeno_1000_Load_CPU = []
    with open(metrics["02a"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_1000_Load_CPU.append(float(item))
    
    set02a_vuvuzela_1000_Load_CPU = []
    with open(metrics["02a"]["vuvuzela"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_1000_Load_CPU.append(float(item))
    
    set02a_pung_1000_Load_CPU = []
    with open(metrics["02a"]["pung"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_1000_Load_CPU.append(float(item))


    set02b_zeno_1000_Load_CPU = []
    with open(metrics["02b"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_1000_Load_CPU.append(float(item))
    
    set02b_vuvuzela_1000_Load_CPU = []
    with open(metrics["02b"]["vuvuzela"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_1000_Load_CPU.append(float(item))
    
    set02b_pung_1000_Load_CPU = []
    with open(metrics["02b"]["pung"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_1000_Load_CPU.append(float(item))


    set03_zeno_1000_Load_CPU = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_CPU.append(float(item))


    # 2000 clients.

    set01_zeno_2000_Load_CPU = []
    with open(metrics["01"]["zeno"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_2000_Load_CPU.append(float(item))
    
    set01_vuvuzela_2000_Load_CPU = []
    with open(metrics["01"]["vuvuzela"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_2000_Load_CPU.append(float(item))
    
    set01_pung_2000_Load_CPU = []
    with open(metrics["01"]["pung"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_2000_Load_CPU.append(float(item))


    set02a_zeno_2000_Load_CPU = []
    with open(metrics["02a"]["zeno"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_2000_Load_CPU.append(float(item))
    
    set02a_vuvuzela_2000_Load_CPU = []
    with open(metrics["02a"]["vuvuzela"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_2000_Load_CPU.append(float(item))
    
    set02a_pung_2000_Load_CPU = []
    with open(metrics["02a"]["pung"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_2000_Load_CPU.append(float(item))


    set02b_zeno_2000_Load_CPU = []
    with open(metrics["02b"]["zeno"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_2000_Load_CPU.append(float(item))
    
    set02b_vuvuzela_2000_Load_CPU = []
    with open(metrics["02b"]["vuvuzela"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_2000_Load_CPU.append(float(item))
    
    set02b_pung_2000_Load_CPU = []
    with open(metrics["02b"]["pung"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_2000_Load_CPU.append(float(item))


    set03_zeno_2000_Load_CPU = []
    with open(metrics["03"]["zeno"]["2000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_2000_Load_CPU.append(float(item))


    # 3000 clients.

    set01_zeno_3000_Load_CPU = []
    with open(metrics["01"]["zeno"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_3000_Load_CPU.append(float(item))
    
    set01_vuvuzela_3000_Load_CPU = []
    with open(metrics["01"]["vuvuzela"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_3000_Load_CPU.append(float(item))
    
    set01_pung_3000_Load_CPU = []
    with open(metrics["01"]["pung"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_3000_Load_CPU.append(float(item))


    set02a_zeno_3000_Load_CPU = []
    with open(metrics["02a"]["zeno"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_3000_Load_CPU.append(float(item))
    
    set02a_vuvuzela_3000_Load_CPU = []
    with open(metrics["02a"]["vuvuzela"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_3000_Load_CPU.append(float(item))
    
    set02a_pung_3000_Load_CPU = []
    with open(metrics["02a"]["pung"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_3000_Load_CPU.append(float(item))


    set02b_zeno_3000_Load_CPU = []
    with open(metrics["02b"]["zeno"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_3000_Load_CPU.append(float(item))
    
    set02b_vuvuzela_3000_Load_CPU = []
    with open(metrics["02b"]["vuvuzela"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_3000_Load_CPU.append(float(item))
    
    set02b_pung_3000_Load_CPU = []
    with open(metrics["02b"]["pung"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_3000_Load_CPU.append(float(item))


    set03_zeno_3000_Load_CPU = []
    with open(metrics["03"]["zeno"]["3000"]["Load"]["Clients"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_3000_Load_CPU.append(float(item))


    # Draw plots.

    width = 0.9

    _, ax = plt.subplots(figsize=(14, 5))

    set01_zeno1 = ax.boxplot(set01_zeno_1000_Load_CPU, positions=[1], widths=width, patch_artist=True, whis='range')
    set02a_zeno1 = ax.boxplot(set02a_zeno_1000_Load_CPU, positions=[2], widths=width, patch_artist=True, whis='range')
    set02b_zeno1 = ax.boxplot(set02b_zeno_1000_Load_CPU, positions=[3], widths=width, patch_artist=True, whis='range')
    set03_zeno1 = ax.boxplot(set03_zeno_1000_Load_CPU, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela1 = ax.boxplot(set01_vuvuzela_1000_Load_CPU, positions=[5], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela1 = ax.boxplot(set02a_vuvuzela_1000_Load_CPU, positions=[6], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela1 = ax.boxplot(set02b_vuvuzela_1000_Load_CPU, positions=[7], widths=width, patch_artist=True, whis='range')
    set01_pung1 = ax.boxplot(set01_pung_1000_Load_CPU, positions=[8], widths=width, patch_artist=True, whis='range')
    set02a_pung1 = ax.boxplot(set02a_pung_1000_Load_CPU, positions=[9], widths=width, patch_artist=True, whis='range')
    set02b_pung1 = ax.boxplot(set02b_pung_1000_Load_CPU, positions=[10], widths=width, patch_artist=True, whis='range')

    set01_zeno2 = ax.boxplot(set01_zeno_2000_Load_CPU, positions=[12], widths=width, patch_artist=True, whis='range')
    set02a_zeno2 = ax.boxplot(set02a_zeno_2000_Load_CPU, positions=[13], widths=width, patch_artist=True, whis='range')
    set02b_zeno2 = ax.boxplot(set02b_zeno_2000_Load_CPU, positions=[14], widths=width, patch_artist=True, whis='range')
    set03_zeno2 = ax.boxplot(set03_zeno_2000_Load_CPU, positions=[15], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela2 = ax.boxplot(set01_vuvuzela_2000_Load_CPU, positions=[16], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela2 = ax.boxplot(set02a_vuvuzela_2000_Load_CPU, positions=[17], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela2 = ax.boxplot(set02b_vuvuzela_2000_Load_CPU, positions=[18], widths=width, patch_artist=True, whis='range')
    set01_pung2 = ax.boxplot(set01_pung_2000_Load_CPU, positions=[19], widths=width, patch_artist=True, whis='range')
    set02a_pung2 = ax.boxplot(set02a_pung_2000_Load_CPU, positions=[20], widths=width, patch_artist=True, whis='range')
    set02b_pung2 = ax.boxplot(set02b_pung_2000_Load_CPU, positions=[21], widths=width, patch_artist=True, whis='range')

    set01_zeno3 = ax.boxplot(set01_zeno_3000_Load_CPU, positions=[23], widths=width, patch_artist=True, whis='range')
    set02a_zeno3 = ax.boxplot(set02a_zeno_3000_Load_CPU, positions=[24], widths=width, patch_artist=True, whis='range')
    set02b_zeno3 = ax.boxplot(set02b_zeno_3000_Load_CPU, positions=[25], widths=width, patch_artist=True, whis='range')
    set03_zeno3 = ax.boxplot(set03_zeno_3000_Load_CPU, positions=[26], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela3 = ax.boxplot(set01_vuvuzela_3000_Load_CPU, positions=[27], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela3 = ax.boxplot(set02a_vuvuzela_3000_Load_CPU, positions=[28], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela3 = ax.boxplot(set02b_vuvuzela_3000_Load_CPU, positions=[29], widths=width, patch_artist=True, whis='range')
    set01_pung3 = ax.boxplot(set01_pung_3000_Load_CPU, positions=[30], widths=width, patch_artist=True, whis='range')
    set02a_pung3 = ax.boxplot(set02a_pung_3000_Load_CPU, positions=[31], widths=width, patch_artist=True, whis='range')
    set02b_pung3 = ax.boxplot(set02b_pung_3000_Load_CPU, positions=[32], widths=width, patch_artist=True, whis='range')


    # Color boxplots.

    setp(set01_zeno1['boxes'], color='black'); setp(set01_zeno2['boxes'], color='black'); setp(set01_zeno3['boxes'], color='black')
    setp(set01_zeno1['boxes'], facecolor='gold'); setp(set01_zeno2['boxes'], facecolor='gold'); setp(set01_zeno3['boxes'], facecolor='gold')
    setp(set01_zeno1['boxes'], hatch='/'); setp(set01_zeno2['boxes'], hatch='/'); setp(set01_zeno3['boxes'], hatch='/')

    setp(set02a_zeno1['boxes'], color='black'); setp(set02a_zeno2['boxes'], color='black'); setp(set02a_zeno3['boxes'], color='black')
    setp(set02a_zeno1['boxes'], facecolor='gold'); setp(set02a_zeno2['boxes'], facecolor='gold'); setp(set02a_zeno3['boxes'], facecolor='gold')
    setp(set02a_zeno1['boxes'], hatch='//'); setp(set02a_zeno2['boxes'], hatch='//'); setp(set02a_zeno3['boxes'], hatch='//')

    setp(set02b_zeno1['boxes'], color='black'); setp(set02b_zeno2['boxes'], color='black'); setp(set02b_zeno3['boxes'], color='black')
    setp(set02b_zeno1['boxes'], facecolor='gold'); setp(set02b_zeno2['boxes'], facecolor='gold'); setp(set02b_zeno3['boxes'], facecolor='gold')
    setp(set02b_zeno1['boxes'], hatch='+'); setp(set02b_zeno2['boxes'], hatch='+'); setp(set02b_zeno3['boxes'], hatch='+')

    setp(set03_zeno1['boxes'], color='black'); setp(set03_zeno2['boxes'], color='black'); setp(set03_zeno3['boxes'], color='black')
    setp(set03_zeno1['boxes'], facecolor='gold'); setp(set03_zeno2['boxes'], facecolor='gold'); setp(set03_zeno3['boxes'], facecolor='gold')
    setp(set03_zeno1['boxes'], hatch='.'); setp(set03_zeno2['boxes'], hatch='.'); setp(set03_zeno3['boxes'], hatch='.')

    setp(set01_vuvuzela1['boxes'], color='black'); setp(set01_vuvuzela2['boxes'], color='black'); setp(set01_vuvuzela3['boxes'], color='black')
    setp(set01_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set01_vuvuzela1['boxes'], hatch='\\'); setp(set01_vuvuzela2['boxes'], hatch='\\'); setp(set01_vuvuzela3['boxes'], hatch='\\')

    setp(set02a_vuvuzela1['boxes'], color='black'); setp(set02a_vuvuzela2['boxes'], color='black'); setp(set02a_vuvuzela3['boxes'], color='black')
    setp(set02a_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02a_vuvuzela1['boxes'], hatch='\\\\'); setp(set02a_vuvuzela2['boxes'], hatch='\\\\'); setp(set02a_vuvuzela3['boxes'], hatch='\\\\')

    setp(set02b_vuvuzela1['boxes'], color='black'); setp(set02b_vuvuzela2['boxes'], color='black'); setp(set02b_vuvuzela3['boxes'], color='black')
    setp(set02b_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02b_vuvuzela1['boxes'], hatch='-'); setp(set02b_vuvuzela2['boxes'], hatch='-'); setp(set02b_vuvuzela3['boxes'], hatch='-')

    setp(set01_pung1['boxes'], color='black'); setp(set01_pung2['boxes'], color='black'); setp(set01_pung3['boxes'], color='black')
    setp(set01_pung1['boxes'], facecolor='steelblue'); setp(set01_pung2['boxes'], facecolor='steelblue'); setp(set01_pung3['boxes'], facecolor='steelblue')
    setp(set01_pung1['boxes'], hatch='x'); setp(set01_pung2['boxes'], hatch='x'); setp(set01_pung3['boxes'], hatch='x')

    setp(set02a_pung1['boxes'], color='black'); setp(set02a_pung2['boxes'], color='black'); setp(set02a_pung3['boxes'], color='black')
    setp(set02a_pung1['boxes'], facecolor='steelblue'); setp(set02a_pung2['boxes'], facecolor='steelblue'); setp(set02a_pung3['boxes'], facecolor='steelblue')
    setp(set02a_pung1['boxes'], hatch='o'); setp(set02a_pung2['boxes'], hatch='o'); setp(set02a_pung3['boxes'], hatch='o')

    setp(set02b_pung1['boxes'], color='black'); setp(set02b_pung2['boxes'], color='black'); setp(set02b_pung3['boxes'], color='black')
    setp(set02b_pung1['boxes'], facecolor='steelblue'); setp(set02b_pung2['boxes'], facecolor='steelblue'); setp(set02b_pung3['boxes'], facecolor='steelblue')
    setp(set02b_pung1['boxes'], hatch='*'); setp(set02b_pung2['boxes'], hatch='*'); setp(set02b_pung3['boxes'], hatch='*')


    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 33])
    ax.set_xticks((5.5, 16.5, 27.5))
    ax.set_xticklabels(('1,000 clients', '2,000 clients', '3,000 clients'))
    ax.set_ylim([0.0, 20.0])
    ax.set_yticks([0, 5, 10, 15, 20])

    # Add a legend.
    ax.legend([
        set01_zeno1['boxes'][0],
        set02a_zeno1['boxes'][0],
        set02b_zeno1['boxes'][0],
        set03_zeno1['boxes'][0],
        set01_vuvuzela1['boxes'][0],
        set02a_vuvuzela1['boxes'][0],
        set02b_vuvuzela1['boxes'][0],
        set01_pung1['boxes'][0],
        set02a_pung1['boxes'][0],
        set02b_pung1['boxes'][0]
    ], [
        'zeno (no impediments)',
        'zeno (high delay, no failures)',
        'zeno (high loss, no failures)',
        'zeno (high network troubles, failures)',
        'vuvuzela (no impediments)',
        'vuvuzela (high delay, no failures)',
        'vuvuzela (high loss, no failures)',
        'pung (no impediments)',
        'pung (high delay, no failures)',
        'pung (high loss, no failures)'
    ],
    loc='upper center',
    ncol=3)
    ax.set_title("Computational Load per Physical Client (10 Logical Clients)")

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Busy CPU (percentage)")

    # plt.savefig(os.path.join(sys.argv[1], "cpu-busy_clients.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "cpu-busy_clients.pdf"), bbox_inches='tight')


def compileLoadMemClients():

    global metrics

    # Ingest data.

    # 1000 clients.

    set01_zeno_1000_Load_Mem = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_Mem.append(float(item))
    
    set01_vuvuzela_1000_Load_Mem = []
    with open(metrics["01"]["vuvuzela"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_1000_Load_Mem.append(float(item))
    
    set01_pung_1000_Load_Mem = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_Mem.append(float(item))


    set02a_zeno_1000_Load_Mem = []
    with open(metrics["02a"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_1000_Load_Mem.append(float(item))

    set02a_vuvuzela_1000_Load_Mem = []
    with open(metrics["02a"]["vuvuzela"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_1000_Load_Mem.append(float(item))

    set02a_pung_1000_Load_Mem = []
    with open(metrics["02a"]["pung"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_1000_Load_Mem.append(float(item))


    set02b_zeno_1000_Load_Mem = []
    with open(metrics["02b"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_1000_Load_Mem.append(float(item))
    
    set02b_vuvuzela_1000_Load_Mem = []
    with open(metrics["02b"]["vuvuzela"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_1000_Load_Mem.append(float(item))
    
    set02b_pung_1000_Load_Mem = []
    with open(metrics["02b"]["pung"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_1000_Load_Mem.append(float(item))


    set03_zeno_1000_Load_Mem = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_Mem.append(float(item))


    # 2000 clients.

    set01_zeno_2000_Load_Mem = []
    with open(metrics["01"]["zeno"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_2000_Load_Mem.append(float(item))
    
    set01_vuvuzela_2000_Load_Mem = []
    with open(metrics["01"]["vuvuzela"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_2000_Load_Mem.append(float(item))
    
    set01_pung_2000_Load_Mem = []
    with open(metrics["01"]["pung"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_2000_Load_Mem.append(float(item))


    set02a_zeno_2000_Load_Mem = []
    with open(metrics["02a"]["zeno"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_2000_Load_Mem.append(float(item))
    
    set02a_vuvuzela_2000_Load_Mem = []
    with open(metrics["02a"]["vuvuzela"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_2000_Load_Mem.append(float(item))
    
    set02a_pung_2000_Load_Mem = []
    with open(metrics["02a"]["pung"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_2000_Load_Mem.append(float(item))


    set02b_zeno_2000_Load_Mem = []
    with open(metrics["02b"]["zeno"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_2000_Load_Mem.append(float(item))
    
    set02b_vuvuzela_2000_Load_Mem = []
    with open(metrics["02b"]["vuvuzela"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_2000_Load_Mem.append(float(item))
    
    set02b_pung_2000_Load_Mem = []
    with open(metrics["02b"]["pung"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_2000_Load_Mem.append(float(item))


    set03_zeno_2000_Load_Mem = []
    with open(metrics["03"]["zeno"]["2000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_2000_Load_Mem.append(float(item))


    # 3000 clients.

    set01_zeno_3000_Load_Mem = []
    with open(metrics["01"]["zeno"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_3000_Load_Mem.append(float(item))
    
    set01_vuvuzela_3000_Load_Mem = []
    with open(metrics["01"]["vuvuzela"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_3000_Load_Mem.append(float(item))
    
    set01_pung_3000_Load_Mem = []
    with open(metrics["01"]["pung"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_3000_Load_Mem.append(float(item))


    set02a_zeno_3000_Load_Mem = []
    with open(metrics["02a"]["zeno"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_3000_Load_Mem.append(float(item))
    
    set02a_vuvuzela_3000_Load_Mem = []
    with open(metrics["02a"]["vuvuzela"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_3000_Load_Mem.append(float(item))
    
    set02a_pung_3000_Load_Mem = []
    with open(metrics["02a"]["pung"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_3000_Load_Mem.append(float(item))


    set02b_zeno_3000_Load_Mem = []
    with open(metrics["02b"]["zeno"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_3000_Load_Mem.append(float(item))
    
    set02b_vuvuzela_3000_Load_Mem = []
    with open(metrics["02b"]["vuvuzela"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_3000_Load_Mem.append(float(item))
    
    set02b_pung_3000_Load_Mem = []
    with open(metrics["02b"]["pung"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_3000_Load_Mem.append(float(item))


    set03_zeno_3000_Load_Mem = []
    with open(metrics["03"]["zeno"]["3000"]["Load"]["Clients"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_3000_Load_Mem.append(float(item))


    # Draw plots.

    width = 0.9

    _, ax = plt.subplots(figsize=(14, 5))

    set01_zeno1 = ax.boxplot(set01_zeno_1000_Load_Mem, positions=[1], widths=width, patch_artist=True, whis='range')
    set02a_zeno1 = ax.boxplot(set02a_zeno_1000_Load_Mem, positions=[2], widths=width, patch_artist=True, whis='range')
    set02b_zeno1 = ax.boxplot(set02b_zeno_1000_Load_Mem, positions=[3], widths=width, patch_artist=True, whis='range')
    set03_zeno1 = ax.boxplot(set03_zeno_1000_Load_Mem, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela1 = ax.boxplot(set01_vuvuzela_1000_Load_Mem, positions=[5], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela1 = ax.boxplot(set02a_vuvuzela_1000_Load_Mem, positions=[6], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela1 = ax.boxplot(set02b_vuvuzela_1000_Load_Mem, positions=[7], widths=width, patch_artist=True, whis='range')
    set01_pung1 = ax.boxplot(set01_pung_1000_Load_Mem, positions=[8], widths=width, patch_artist=True, whis='range')
    set02a_pung1 = ax.boxplot(set02a_pung_1000_Load_Mem, positions=[9], widths=width, patch_artist=True, whis='range')
    set02b_pung1 = ax.boxplot(set02b_pung_1000_Load_Mem, positions=[10], widths=width, patch_artist=True, whis='range')

    set01_zeno2 = ax.boxplot(set01_zeno_2000_Load_Mem, positions=[12], widths=width, patch_artist=True, whis='range')
    set02a_zeno2 = ax.boxplot(set02a_zeno_2000_Load_Mem, positions=[13], widths=width, patch_artist=True, whis='range')
    set02b_zeno2 = ax.boxplot(set02b_zeno_2000_Load_Mem, positions=[14], widths=width, patch_artist=True, whis='range')
    set03_zeno2 = ax.boxplot(set03_zeno_2000_Load_Mem, positions=[15], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela2 = ax.boxplot(set01_vuvuzela_2000_Load_Mem, positions=[16], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela2 = ax.boxplot(set02a_vuvuzela_2000_Load_Mem, positions=[17], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela2 = ax.boxplot(set02b_vuvuzela_2000_Load_Mem, positions=[18], widths=width, patch_artist=True, whis='range')
    set01_pung2 = ax.boxplot(set01_pung_2000_Load_Mem, positions=[19], widths=width, patch_artist=True, whis='range')
    set02a_pung2 = ax.boxplot(set02a_pung_2000_Load_Mem, positions=[20], widths=width, patch_artist=True, whis='range')
    set02b_pung2 = ax.boxplot(set02b_pung_2000_Load_Mem, positions=[21], widths=width, patch_artist=True, whis='range')

    set01_zeno3 = ax.boxplot(set01_zeno_3000_Load_Mem, positions=[23], widths=width, patch_artist=True, whis='range')
    set02a_zeno3 = ax.boxplot(set02a_zeno_3000_Load_Mem, positions=[24], widths=width, patch_artist=True, whis='range')
    set02b_zeno3 = ax.boxplot(set02b_zeno_3000_Load_Mem, positions=[25], widths=width, patch_artist=True, whis='range')
    set03_zeno3 = ax.boxplot(set03_zeno_3000_Load_Mem, positions=[26], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela3 = ax.boxplot(set01_vuvuzela_3000_Load_Mem, positions=[27], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela3 = ax.boxplot(set02a_vuvuzela_3000_Load_Mem, positions=[28], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela3 = ax.boxplot(set02b_vuvuzela_3000_Load_Mem, positions=[29], widths=width, patch_artist=True, whis='range')
    set01_pung3 = ax.boxplot(set01_pung_3000_Load_Mem, positions=[30], widths=width, patch_artist=True, whis='range')
    set02a_pung3 = ax.boxplot(set02a_pung_3000_Load_Mem, positions=[31], widths=width, patch_artist=True, whis='range')
    set02b_pung3 = ax.boxplot(set02b_pung_3000_Load_Mem, positions=[32], widths=width, patch_artist=True, whis='range')


    # Log values for text mention.

    print("\nClients:\n")
    for scenData in [set01_zeno1, set02a_zeno1, set02b_zeno1, set03_zeno1,
                     set01_vuvuzela1, set02a_vuvuzela1, set02b_vuvuzela1,
                     set01_pung1, set02a_pung1, set02b_pung1,
                     set01_zeno2, set02a_zeno2, set02b_zeno2, set03_zeno2,
                     set01_vuvuzela2, set02a_vuvuzela2, set02b_vuvuzela2,
                     set01_pung2, set02a_pung2, set02b_pung2,
                     set01_zeno3, set02a_zeno3, set02b_zeno3, set03_zeno3,
                     set01_vuvuzela3, set02a_vuvuzela3, set02b_vuvuzela3,
                     set01_pung3, set02a_pung3, set02b_pung3]:

        for whis in scenData['whiskers']:
            print("whis=", whis.get_ydata()[1])
        
        for med in scenData['medians']:
            print(" med=", med.get_ydata()[1])
        
        print("")
    print("")


    # Color boxplots.

    setp(set01_zeno1['boxes'], color='black'); setp(set01_zeno2['boxes'], color='black'); setp(set01_zeno3['boxes'], color='black')
    setp(set01_zeno1['boxes'], facecolor='gold'); setp(set01_zeno2['boxes'], facecolor='gold'); setp(set01_zeno3['boxes'], facecolor='gold')
    setp(set01_zeno1['boxes'], hatch='/'); setp(set01_zeno2['boxes'], hatch='/'); setp(set01_zeno3['boxes'], hatch='/')

    setp(set02a_zeno1['boxes'], color='black'); setp(set02a_zeno2['boxes'], color='black'); setp(set02a_zeno3['boxes'], color='black')
    setp(set02a_zeno1['boxes'], facecolor='gold'); setp(set02a_zeno2['boxes'], facecolor='gold'); setp(set02a_zeno3['boxes'], facecolor='gold')
    setp(set02a_zeno1['boxes'], hatch='//'); setp(set02a_zeno2['boxes'], hatch='//'); setp(set02a_zeno3['boxes'], hatch='//')

    setp(set02b_zeno1['boxes'], color='black'); setp(set02b_zeno2['boxes'], color='black'); setp(set02b_zeno3['boxes'], color='black')
    setp(set02b_zeno1['boxes'], facecolor='gold'); setp(set02b_zeno2['boxes'], facecolor='gold'); setp(set02b_zeno3['boxes'], facecolor='gold')
    setp(set02b_zeno1['boxes'], hatch='+'); setp(set02b_zeno2['boxes'], hatch='+'); setp(set02b_zeno3['boxes'], hatch='+')

    setp(set03_zeno1['boxes'], color='black'); setp(set03_zeno2['boxes'], color='black'); setp(set03_zeno3['boxes'], color='black')
    setp(set03_zeno1['boxes'], facecolor='gold'); setp(set03_zeno2['boxes'], facecolor='gold'); setp(set03_zeno3['boxes'], facecolor='gold')
    setp(set03_zeno1['boxes'], hatch='.'); setp(set03_zeno2['boxes'], hatch='.'); setp(set03_zeno3['boxes'], hatch='.')

    setp(set01_vuvuzela1['boxes'], color='black'); setp(set01_vuvuzela2['boxes'], color='black'); setp(set01_vuvuzela3['boxes'], color='black')
    setp(set01_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set01_vuvuzela1['boxes'], hatch='\\'); setp(set01_vuvuzela2['boxes'], hatch='\\'); setp(set01_vuvuzela3['boxes'], hatch='\\')

    setp(set02a_vuvuzela1['boxes'], color='black'); setp(set02a_vuvuzela2['boxes'], color='black'); setp(set02a_vuvuzela3['boxes'], color='black')
    setp(set02a_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02a_vuvuzela1['boxes'], hatch='\\\\'); setp(set02a_vuvuzela2['boxes'], hatch='\\\\'); setp(set02a_vuvuzela3['boxes'], hatch='\\\\')

    setp(set02b_vuvuzela1['boxes'], color='black'); setp(set02b_vuvuzela2['boxes'], color='black'); setp(set02b_vuvuzela3['boxes'], color='black')
    setp(set02b_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02b_vuvuzela1['boxes'], hatch='-'); setp(set02b_vuvuzela2['boxes'], hatch='-'); setp(set02b_vuvuzela3['boxes'], hatch='-')

    setp(set01_pung1['boxes'], color='black'); setp(set01_pung2['boxes'], color='black'); setp(set01_pung3['boxes'], color='black')
    setp(set01_pung1['boxes'], facecolor='steelblue'); setp(set01_pung2['boxes'], facecolor='steelblue'); setp(set01_pung3['boxes'], facecolor='steelblue')
    setp(set01_pung1['boxes'], hatch='x'); setp(set01_pung2['boxes'], hatch='x'); setp(set01_pung3['boxes'], hatch='x')

    setp(set02a_pung1['boxes'], color='black'); setp(set02a_pung2['boxes'], color='black'); setp(set02a_pung3['boxes'], color='black')
    setp(set02a_pung1['boxes'], facecolor='steelblue'); setp(set02a_pung2['boxes'], facecolor='steelblue'); setp(set02a_pung3['boxes'], facecolor='steelblue')
    setp(set02a_pung1['boxes'], hatch='o'); setp(set02a_pung2['boxes'], hatch='o'); setp(set02a_pung3['boxes'], hatch='o')

    setp(set02b_pung1['boxes'], color='black'); setp(set02b_pung2['boxes'], color='black'); setp(set02b_pung3['boxes'], color='black')
    setp(set02b_pung1['boxes'], facecolor='steelblue'); setp(set02b_pung2['boxes'], facecolor='steelblue'); setp(set02b_pung3['boxes'], facecolor='steelblue')
    setp(set02b_pung1['boxes'], hatch='*'); setp(set02b_pung2['boxes'], hatch='*'); setp(set02b_pung3['boxes'], hatch='*')


    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 33])
    ax.set_xticks((5.5, 16.5, 27.5))
    ax.set_xticklabels(('1,000 clients', '2,000 clients', '3,000 clients'))
    ax.set_ylim([500.0, 3000.0])
    ax.set_yticks([500, 1000, 1500, 2000, 2500, 3000])

    # Add a legend.
    ax.legend([
        set01_zeno1['boxes'][0],
        set02a_zeno1['boxes'][0],
        set02b_zeno1['boxes'][0],
        set03_zeno1['boxes'][0],
        set01_vuvuzela1['boxes'][0],
        set02a_vuvuzela1['boxes'][0],
        set02b_vuvuzela1['boxes'][0],
        set01_pung1['boxes'][0],
        set02a_pung1['boxes'][0],
        set02b_pung1['boxes'][0]
    ], [
        'zeno (no impediments)',
        'zeno (high delay, no failures)',
        'zeno (high loss, no failures)',
        'zeno (high network troubles, failures)',
        'vuvuzela (no impediments)',
        'vuvuzela (high delay, no failures)',
        'vuvuzela (high loss, no failures)',
        'pung (no impediments)',
        'pung (high delay, no failures)',
        'pung (high loss, no failures)'
    ],
    loc='upper center',
    ncol=3)
    ax.set_title("Memory Load per Physical Client (10 Logical Clients)")

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Used memory (MB)")

    # plt.savefig(os.path.join(sys.argv[1], "memory-used_clients.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "memory-used_clients.pdf"), bbox_inches='tight')


def compileLoadCPUServers():

    global metrics

    # Ingest data.

    # 1000 clients.

    set01_zeno_1000_Load_CPU = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_CPU.append(float(item))
    
    set01_vuvuzela_1000_Load_CPU = []
    with open(metrics["01"]["vuvuzela"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_1000_Load_CPU.append(float(item))
    
    set01_pung_1000_Load_CPU = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_CPU.append(float(item))


    set02a_zeno_1000_Load_CPU = []
    with open(metrics["02a"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_1000_Load_CPU.append(float(item))
    
    set02a_vuvuzela_1000_Load_CPU = []
    with open(metrics["02a"]["vuvuzela"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_1000_Load_CPU.append(float(item))
    
    set02a_pung_1000_Load_CPU = []
    with open(metrics["02a"]["pung"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_1000_Load_CPU.append(float(item))


    set02b_zeno_1000_Load_CPU = []
    with open(metrics["02b"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_1000_Load_CPU.append(float(item))
    
    set02b_vuvuzela_1000_Load_CPU = []
    with open(metrics["02b"]["vuvuzela"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_1000_Load_CPU.append(float(item))
    
    set02b_pung_1000_Load_CPU = []
    with open(metrics["02b"]["pung"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_1000_Load_CPU.append(float(item))


    set03_zeno_1000_Load_CPU = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_CPU.append(float(item))


    # 2000 clients.

    set01_zeno_2000_Load_CPU = []
    with open(metrics["01"]["zeno"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_2000_Load_CPU.append(float(item))
    
    set01_vuvuzela_2000_Load_CPU = []
    with open(metrics["01"]["vuvuzela"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_2000_Load_CPU.append(float(item))
    
    set01_pung_2000_Load_CPU = []
    with open(metrics["01"]["pung"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_2000_Load_CPU.append(float(item))


    set02a_zeno_2000_Load_CPU = []
    with open(metrics["02a"]["zeno"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_2000_Load_CPU.append(float(item))
    
    set02a_vuvuzela_2000_Load_CPU = []
    with open(metrics["02a"]["vuvuzela"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_2000_Load_CPU.append(float(item))
    
    set02a_pung_2000_Load_CPU = []
    with open(metrics["02a"]["pung"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_2000_Load_CPU.append(float(item))


    set02b_zeno_2000_Load_CPU = []
    with open(metrics["02b"]["zeno"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_2000_Load_CPU.append(float(item))
    
    set02b_vuvuzela_2000_Load_CPU = []
    with open(metrics["02b"]["vuvuzela"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_2000_Load_CPU.append(float(item))
    
    set02b_pung_2000_Load_CPU = []
    with open(metrics["02b"]["pung"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_2000_Load_CPU.append(float(item))


    set03_zeno_2000_Load_CPU = []
    with open(metrics["03"]["zeno"]["2000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_2000_Load_CPU.append(float(item))


    # 3000 clients.

    set01_zeno_3000_Load_CPU = []
    with open(metrics["01"]["zeno"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_3000_Load_CPU.append(float(item))
    
    set01_vuvuzela_3000_Load_CPU = []
    with open(metrics["01"]["vuvuzela"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_3000_Load_CPU.append(float(item))
    
    set01_pung_3000_Load_CPU = []
    with open(metrics["01"]["pung"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_3000_Load_CPU.append(float(item))


    set02a_zeno_3000_Load_CPU = []
    with open(metrics["02a"]["zeno"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_3000_Load_CPU.append(float(item))
    
    set02a_vuvuzela_3000_Load_CPU = []
    with open(metrics["02a"]["vuvuzela"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_3000_Load_CPU.append(float(item))
    
    set02a_pung_3000_Load_CPU = []
    with open(metrics["02a"]["pung"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_3000_Load_CPU.append(float(item))


    set02b_zeno_3000_Load_CPU = []
    with open(metrics["02b"]["zeno"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_3000_Load_CPU.append(float(item))
    
    set02b_vuvuzela_3000_Load_CPU = []
    with open(metrics["02b"]["vuvuzela"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_3000_Load_CPU.append(float(item))
    
    set02b_pung_3000_Load_CPU = []
    with open(metrics["02b"]["pung"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_3000_Load_CPU.append(float(item))


    set03_zeno_3000_Load_CPU = []
    with open(metrics["03"]["zeno"]["3000"]["Load"]["Servers"]["CPU"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_3000_Load_CPU.append(float(item))

    # Draw plots.

    width = 0.9

    _, ax = plt.subplots(figsize=(14, 5))

    set01_zeno1 = ax.boxplot(set01_zeno_1000_Load_CPU, positions=[1], widths=width, patch_artist=True, whis='range')
    set02a_zeno1 = ax.boxplot(set02a_zeno_1000_Load_CPU, positions=[2], widths=width, patch_artist=True, whis='range')
    set02b_zeno1 = ax.boxplot(set02b_zeno_1000_Load_CPU, positions=[3], widths=width, patch_artist=True, whis='range')
    set03_zeno1 = ax.boxplot(set03_zeno_1000_Load_CPU, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela1 = ax.boxplot(set01_vuvuzela_1000_Load_CPU, positions=[5], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela1 = ax.boxplot(set02a_vuvuzela_1000_Load_CPU, positions=[6], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela1 = ax.boxplot(set02b_vuvuzela_1000_Load_CPU, positions=[7], widths=width, patch_artist=True, whis='range')
    set01_pung1 = ax.boxplot(set01_pung_1000_Load_CPU, positions=[8], widths=width, patch_artist=True, whis='range')
    set02a_pung1 = ax.boxplot(set02a_pung_1000_Load_CPU, positions=[9], widths=width, patch_artist=True, whis='range')
    set02b_pung1 = ax.boxplot(set02b_pung_1000_Load_CPU, positions=[10], widths=width, patch_artist=True, whis='range')

    set01_zeno2 = ax.boxplot(set01_zeno_2000_Load_CPU, positions=[12], widths=width, patch_artist=True, whis='range')
    set02a_zeno2 = ax.boxplot(set02a_zeno_2000_Load_CPU, positions=[13], widths=width, patch_artist=True, whis='range')
    set02b_zeno2 = ax.boxplot(set02b_zeno_2000_Load_CPU, positions=[14], widths=width, patch_artist=True, whis='range')
    set03_zeno2 = ax.boxplot(set03_zeno_2000_Load_CPU, positions=[15], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela2 = ax.boxplot(set01_vuvuzela_2000_Load_CPU, positions=[16], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela2 = ax.boxplot(set02a_vuvuzela_2000_Load_CPU, positions=[17], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela2 = ax.boxplot(set02b_vuvuzela_2000_Load_CPU, positions=[18], widths=width, patch_artist=True, whis='range')
    set01_pung2 = ax.boxplot(set01_pung_2000_Load_CPU, positions=[19], widths=width, patch_artist=True, whis='range')
    set02a_pung2 = ax.boxplot(set02a_pung_2000_Load_CPU, positions=[20], widths=width, patch_artist=True, whis='range')
    set02b_pung2 = ax.boxplot(set02b_pung_2000_Load_CPU, positions=[21], widths=width, patch_artist=True, whis='range')

    set01_zeno3 = ax.boxplot(set01_zeno_3000_Load_CPU, positions=[23], widths=width, patch_artist=True, whis='range')
    set02a_zeno3 = ax.boxplot(set02a_zeno_3000_Load_CPU, positions=[24], widths=width, patch_artist=True, whis='range')
    set02b_zeno3 = ax.boxplot(set02b_zeno_3000_Load_CPU, positions=[25], widths=width, patch_artist=True, whis='range')
    set03_zeno3 = ax.boxplot(set03_zeno_3000_Load_CPU, positions=[26], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela3 = ax.boxplot(set01_vuvuzela_3000_Load_CPU, positions=[27], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela3 = ax.boxplot(set02a_vuvuzela_3000_Load_CPU, positions=[28], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela3 = ax.boxplot(set02b_vuvuzela_3000_Load_CPU, positions=[29], widths=width, patch_artist=True, whis='range')
    set01_pung3 = ax.boxplot(set01_pung_3000_Load_CPU, positions=[30], widths=width, patch_artist=True, whis='range')
    set02a_pung3 = ax.boxplot(set02a_pung_3000_Load_CPU, positions=[31], widths=width, patch_artist=True, whis='range')
    set02b_pung3 = ax.boxplot(set02b_pung_3000_Load_CPU, positions=[32], widths=width, patch_artist=True, whis='range')


    # Color boxplots.

    setp(set01_zeno1['boxes'], color='black'); setp(set01_zeno2['boxes'], color='black'); setp(set01_zeno3['boxes'], color='black')
    setp(set01_zeno1['boxes'], facecolor='gold'); setp(set01_zeno2['boxes'], facecolor='gold'); setp(set01_zeno3['boxes'], facecolor='gold')
    setp(set01_zeno1['boxes'], hatch='/'); setp(set01_zeno2['boxes'], hatch='/'); setp(set01_zeno3['boxes'], hatch='/')

    setp(set02a_zeno1['boxes'], color='black'); setp(set02a_zeno2['boxes'], color='black'); setp(set02a_zeno3['boxes'], color='black')
    setp(set02a_zeno1['boxes'], facecolor='gold'); setp(set02a_zeno2['boxes'], facecolor='gold'); setp(set02a_zeno3['boxes'], facecolor='gold')
    setp(set02a_zeno1['boxes'], hatch='//'); setp(set02a_zeno2['boxes'], hatch='//'); setp(set02a_zeno3['boxes'], hatch='//')

    setp(set02b_zeno1['boxes'], color='black'); setp(set02b_zeno2['boxes'], color='black'); setp(set02b_zeno3['boxes'], color='black')
    setp(set02b_zeno1['boxes'], facecolor='gold'); setp(set02b_zeno2['boxes'], facecolor='gold'); setp(set02b_zeno3['boxes'], facecolor='gold')
    setp(set02b_zeno1['boxes'], hatch='+'); setp(set02b_zeno2['boxes'], hatch='+'); setp(set02b_zeno3['boxes'], hatch='+')

    setp(set03_zeno1['boxes'], color='black'); setp(set03_zeno2['boxes'], color='black'); setp(set03_zeno3['boxes'], color='black')
    setp(set03_zeno1['boxes'], facecolor='gold'); setp(set03_zeno2['boxes'], facecolor='gold'); setp(set03_zeno3['boxes'], facecolor='gold')
    setp(set03_zeno1['boxes'], hatch='.'); setp(set03_zeno2['boxes'], hatch='.'); setp(set03_zeno3['boxes'], hatch='.')

    setp(set01_vuvuzela1['boxes'], color='black'); setp(set01_vuvuzela2['boxes'], color='black'); setp(set01_vuvuzela3['boxes'], color='black')
    setp(set01_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set01_vuvuzela1['boxes'], hatch='\\'); setp(set01_vuvuzela2['boxes'], hatch='\\'); setp(set01_vuvuzela3['boxes'], hatch='\\')

    setp(set02a_vuvuzela1['boxes'], color='black'); setp(set02a_vuvuzela2['boxes'], color='black'); setp(set02a_vuvuzela3['boxes'], color='black')
    setp(set02a_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02a_vuvuzela1['boxes'], hatch='\\\\'); setp(set02a_vuvuzela2['boxes'], hatch='\\\\'); setp(set02a_vuvuzela3['boxes'], hatch='\\\\')

    setp(set02b_vuvuzela1['boxes'], color='black'); setp(set02b_vuvuzela2['boxes'], color='black'); setp(set02b_vuvuzela3['boxes'], color='black')
    setp(set02b_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02b_vuvuzela1['boxes'], hatch='-'); setp(set02b_vuvuzela2['boxes'], hatch='-'); setp(set02b_vuvuzela3['boxes'], hatch='-')

    setp(set01_pung1['boxes'], color='black'); setp(set01_pung2['boxes'], color='black'); setp(set01_pung3['boxes'], color='black')
    setp(set01_pung1['boxes'], facecolor='steelblue'); setp(set01_pung2['boxes'], facecolor='steelblue'); setp(set01_pung3['boxes'], facecolor='steelblue')
    setp(set01_pung1['boxes'], hatch='x'); setp(set01_pung2['boxes'], hatch='x'); setp(set01_pung3['boxes'], hatch='x')

    setp(set02a_pung1['boxes'], color='black'); setp(set02a_pung2['boxes'], color='black'); setp(set02a_pung3['boxes'], color='black')
    setp(set02a_pung1['boxes'], facecolor='steelblue'); setp(set02a_pung2['boxes'], facecolor='steelblue'); setp(set02a_pung3['boxes'], facecolor='steelblue')
    setp(set02a_pung1['boxes'], hatch='o'); setp(set02a_pung2['boxes'], hatch='o'); setp(set02a_pung3['boxes'], hatch='o')

    setp(set02b_pung1['boxes'], color='black'); setp(set02b_pung2['boxes'], color='black'); setp(set02b_pung3['boxes'], color='black')
    setp(set02b_pung1['boxes'], facecolor='steelblue'); setp(set02b_pung2['boxes'], facecolor='steelblue'); setp(set02b_pung3['boxes'], facecolor='steelblue')
    setp(set02b_pung1['boxes'], hatch='*'); setp(set02b_pung2['boxes'], hatch='*'); setp(set02b_pung3['boxes'], hatch='*')


    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 33])
    ax.set_xticks((5.5, 16.5, 27.5))
    ax.set_xticklabels(('1,000 clients', '2,000 clients', '3,000 clients'))
    ax.set_ylim([0.0, 120.0])
    ax.set_yticks([0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100])

    # Add a legend.
    ax.legend([
        set01_zeno1['boxes'][0],
        set02a_zeno1['boxes'][0],
        set02b_zeno1['boxes'][0],
        set03_zeno1['boxes'][0],
        set01_vuvuzela1['boxes'][0],
        set02a_vuvuzela1['boxes'][0],
        set02b_vuvuzela1['boxes'][0],
        set01_pung1['boxes'][0],
        set02a_pung1['boxes'][0],
        set02b_pung1['boxes'][0]
    ], [
        'zeno (no impediments)',
        'zeno (high delay, no failures)',
        'zeno (high loss, no failures)',
        'zeno (high network troubles, failures)',
        'vuvuzela (no impediments)',
        'vuvuzela (high delay, no failures)',
        'vuvuzela (high loss, no failures)',
        'pung (no impediments)',
        'pung (high delay, no failures)',
        'pung (high loss, no failures)'
    ],
    loc='upper center',
    ncol=3)
    ax.set_title("Computational Load per Server")

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Busy CPU (percentage)")

    # plt.savefig(os.path.join(sys.argv[1], "cpu-busy_servers.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "cpu-busy_servers.pdf"), bbox_inches='tight')


def compileLoadMemServers():

    global metrics

    # Ingest data.

    # 1000 clients.

    set01_zeno_1000_Load_Mem = []
    with open(metrics["01"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_1000_Load_Mem.append(float(item))
    
    set01_vuvuzela_1000_Load_Mem = []
    with open(metrics["01"]["vuvuzela"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_1000_Load_Mem.append(float(item))
    
    set01_pung_1000_Load_Mem = []
    with open(metrics["01"]["pung"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_1000_Load_Mem.append(float(item))


    set02a_zeno_1000_Load_Mem = []
    with open(metrics["02a"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_1000_Load_Mem.append(float(item))

    set02a_vuvuzela_1000_Load_Mem = []
    with open(metrics["02a"]["vuvuzela"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_1000_Load_Mem.append(float(item))

    set02a_pung_1000_Load_Mem = []
    with open(metrics["02a"]["pung"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_1000_Load_Mem.append(float(item))


    set02b_zeno_1000_Load_Mem = []
    with open(metrics["02b"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_1000_Load_Mem.append(float(item))
    
    set02b_vuvuzela_1000_Load_Mem = []
    with open(metrics["02b"]["vuvuzela"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_1000_Load_Mem.append(float(item))
    
    set02b_pung_1000_Load_Mem = []
    with open(metrics["02b"]["pung"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_1000_Load_Mem.append(float(item))


    set03_zeno_1000_Load_Mem = []
    with open(metrics["03"]["zeno"]["1000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_1000_Load_Mem.append(float(item))


    # 2000 clients.

    set01_zeno_2000_Load_Mem = []
    with open(metrics["01"]["zeno"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_2000_Load_Mem.append(float(item))
    
    set01_vuvuzela_2000_Load_Mem = []
    with open(metrics["01"]["vuvuzela"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_2000_Load_Mem.append(float(item))
    
    set01_pung_2000_Load_Mem = []
    with open(metrics["01"]["pung"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_2000_Load_Mem.append(float(item))


    set02a_zeno_2000_Load_Mem = []
    with open(metrics["02a"]["zeno"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_2000_Load_Mem.append(float(item))
    
    set02a_vuvuzela_2000_Load_Mem = []
    with open(metrics["02a"]["vuvuzela"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_2000_Load_Mem.append(float(item))
    
    set02a_pung_2000_Load_Mem = []
    with open(metrics["02a"]["pung"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_2000_Load_Mem.append(float(item))


    set02b_zeno_2000_Load_Mem = []
    with open(metrics["02b"]["zeno"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_2000_Load_Mem.append(float(item))
    
    set02b_vuvuzela_2000_Load_Mem = []
    with open(metrics["02b"]["vuvuzela"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_2000_Load_Mem.append(float(item))
    
    set02b_pung_2000_Load_Mem = []
    with open(metrics["02b"]["pung"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_2000_Load_Mem.append(float(item))


    set03_zeno_2000_Load_Mem = []
    with open(metrics["03"]["zeno"]["2000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_2000_Load_Mem.append(float(item))


    # 3000 clients.

    set01_zeno_3000_Load_Mem = []
    with open(metrics["01"]["zeno"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_zeno_3000_Load_Mem.append(float(item))
    
    set01_vuvuzela_3000_Load_Mem = []
    with open(metrics["01"]["vuvuzela"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_vuvuzela_3000_Load_Mem.append(float(item))
    
    set01_pung_3000_Load_Mem = []
    with open(metrics["01"]["pung"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set01_pung_3000_Load_Mem.append(float(item))


    set02a_zeno_3000_Load_Mem = []
    with open(metrics["02a"]["zeno"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_zeno_3000_Load_Mem.append(float(item))
    
    set02a_vuvuzela_3000_Load_Mem = []
    with open(metrics["02a"]["vuvuzela"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_vuvuzela_3000_Load_Mem.append(float(item))
    
    set02a_pung_3000_Load_Mem = []
    with open(metrics["02a"]["pung"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02a_pung_3000_Load_Mem.append(float(item))


    set02b_zeno_3000_Load_Mem = []
    with open(metrics["02b"]["zeno"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_zeno_3000_Load_Mem.append(float(item))
    
    set02b_vuvuzela_3000_Load_Mem = []
    with open(metrics["02b"]["vuvuzela"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_vuvuzela_3000_Load_Mem.append(float(item))
    
    set02b_pung_3000_Load_Mem = []
    with open(metrics["02b"]["pung"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set02b_pung_3000_Load_Mem.append(float(item))


    set03_zeno_3000_Load_Mem = []
    with open(metrics["03"]["zeno"]["3000"]["Load"]["Servers"]["Mem"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                set03_zeno_3000_Load_Mem.append(float(item))


    # Draw plots.

    width = 0.9

    _, ax = plt.subplots(figsize=(14, 5))

    set01_zeno1 = ax.boxplot(set01_zeno_1000_Load_Mem, positions=[1], widths=width, patch_artist=True, whis='range')
    set02a_zeno1 = ax.boxplot(set02a_zeno_1000_Load_Mem, positions=[2], widths=width, patch_artist=True, whis='range')
    set02b_zeno1 = ax.boxplot(set02b_zeno_1000_Load_Mem, positions=[3], widths=width, patch_artist=True, whis='range')
    set03_zeno1 = ax.boxplot(set03_zeno_1000_Load_Mem, positions=[4], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela1 = ax.boxplot(set01_vuvuzela_1000_Load_Mem, positions=[5], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela1 = ax.boxplot(set02a_vuvuzela_1000_Load_Mem, positions=[6], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela1 = ax.boxplot(set02b_vuvuzela_1000_Load_Mem, positions=[7], widths=width, patch_artist=True, whis='range')
    set01_pung1 = ax.boxplot(set01_pung_1000_Load_Mem, positions=[8], widths=width, patch_artist=True, whis='range')
    set02a_pung1 = ax.boxplot(set02a_pung_1000_Load_Mem, positions=[9], widths=width, patch_artist=True, whis='range')
    set02b_pung1 = ax.boxplot(set02b_pung_1000_Load_Mem, positions=[10], widths=width, patch_artist=True, whis='range')

    set01_zeno2 = ax.boxplot(set01_zeno_2000_Load_Mem, positions=[12], widths=width, patch_artist=True, whis='range')
    set02a_zeno2 = ax.boxplot(set02a_zeno_2000_Load_Mem, positions=[13], widths=width, patch_artist=True, whis='range')
    set02b_zeno2 = ax.boxplot(set02b_zeno_2000_Load_Mem, positions=[14], widths=width, patch_artist=True, whis='range')
    set03_zeno2 = ax.boxplot(set03_zeno_2000_Load_Mem, positions=[15], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela2 = ax.boxplot(set01_vuvuzela_2000_Load_Mem, positions=[16], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela2 = ax.boxplot(set02a_vuvuzela_2000_Load_Mem, positions=[17], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela2 = ax.boxplot(set02b_vuvuzela_2000_Load_Mem, positions=[18], widths=width, patch_artist=True, whis='range')
    set01_pung2 = ax.boxplot(set01_pung_2000_Load_Mem, positions=[19], widths=width, patch_artist=True, whis='range')
    set02a_pung2 = ax.boxplot(set02a_pung_2000_Load_Mem, positions=[20], widths=width, patch_artist=True, whis='range')
    set02b_pung2 = ax.boxplot(set02b_pung_2000_Load_Mem, positions=[21], widths=width, patch_artist=True, whis='range')

    set01_zeno3 = ax.boxplot(set01_zeno_3000_Load_Mem, positions=[23], widths=width, patch_artist=True, whis='range')
    set02a_zeno3 = ax.boxplot(set02a_zeno_3000_Load_Mem, positions=[24], widths=width, patch_artist=True, whis='range')
    set02b_zeno3 = ax.boxplot(set02b_zeno_3000_Load_Mem, positions=[25], widths=width, patch_artist=True, whis='range')
    set03_zeno3 = ax.boxplot(set03_zeno_3000_Load_Mem, positions=[26], widths=width, patch_artist=True, whis='range')
    set01_vuvuzela3 = ax.boxplot(set01_vuvuzela_3000_Load_Mem, positions=[27], widths=width, patch_artist=True, whis='range')
    set02a_vuvuzela3 = ax.boxplot(set02a_vuvuzela_3000_Load_Mem, positions=[28], widths=width, patch_artist=True, whis='range')
    set02b_vuvuzela3 = ax.boxplot(set02b_vuvuzela_3000_Load_Mem, positions=[29], widths=width, patch_artist=True, whis='range')
    set01_pung3 = ax.boxplot(set01_pung_3000_Load_Mem, positions=[30], widths=width, patch_artist=True, whis='range')
    set02a_pung3 = ax.boxplot(set02a_pung_3000_Load_Mem, positions=[31], widths=width, patch_artist=True, whis='range')
    set02b_pung3 = ax.boxplot(set02b_pung_3000_Load_Mem, positions=[32], widths=width, patch_artist=True, whis='range')


    # Log values for text mention.
    
    print("Servers:\n")
    for scenData in [set01_zeno1, set02a_zeno1, set02b_zeno1, set03_zeno1,
                     set01_vuvuzela1, set02a_vuvuzela1, set02b_vuvuzela1,
                     set01_pung1, set02a_pung1, set02b_pung1,
                     set01_zeno2, set02a_zeno2, set02b_zeno2, set03_zeno2,
                     set01_vuvuzela2, set02a_vuvuzela2, set02b_vuvuzela2,
                     set01_pung2, set02a_pung2, set02b_pung2,
                     set01_zeno3, set02a_zeno3, set02b_zeno3, set03_zeno3,
                     set01_vuvuzela3, set02a_vuvuzela3, set02b_vuvuzela3,
                     set01_pung3, set02a_pung3, set02b_pung3]:

        for whis in scenData['whiskers']:
            print("whis=", whis.get_ydata()[1])
        
        for med in scenData['medians']:
            print(" med=", med.get_ydata()[1])
        
        print("")


    # Color boxplots.

    setp(set01_zeno1['boxes'], color='black'); setp(set01_zeno2['boxes'], color='black'); setp(set01_zeno3['boxes'], color='black')
    setp(set01_zeno1['boxes'], facecolor='gold'); setp(set01_zeno2['boxes'], facecolor='gold'); setp(set01_zeno3['boxes'], facecolor='gold')
    setp(set01_zeno1['boxes'], hatch='/'); setp(set01_zeno2['boxes'], hatch='/'); setp(set01_zeno3['boxes'], hatch='/')

    setp(set02a_zeno1['boxes'], color='black'); setp(set02a_zeno2['boxes'], color='black'); setp(set02a_zeno3['boxes'], color='black')
    setp(set02a_zeno1['boxes'], facecolor='gold'); setp(set02a_zeno2['boxes'], facecolor='gold'); setp(set02a_zeno3['boxes'], facecolor='gold')
    setp(set02a_zeno1['boxes'], hatch='//'); setp(set02a_zeno2['boxes'], hatch='//'); setp(set02a_zeno3['boxes'], hatch='//')

    setp(set02b_zeno1['boxes'], color='black'); setp(set02b_zeno2['boxes'], color='black'); setp(set02b_zeno3['boxes'], color='black')
    setp(set02b_zeno1['boxes'], facecolor='gold'); setp(set02b_zeno2['boxes'], facecolor='gold'); setp(set02b_zeno3['boxes'], facecolor='gold')
    setp(set02b_zeno1['boxes'], hatch='+'); setp(set02b_zeno2['boxes'], hatch='+'); setp(set02b_zeno3['boxes'], hatch='+')

    setp(set03_zeno1['boxes'], color='black'); setp(set03_zeno2['boxes'], color='black'); setp(set03_zeno3['boxes'], color='black')
    setp(set03_zeno1['boxes'], facecolor='gold'); setp(set03_zeno2['boxes'], facecolor='gold'); setp(set03_zeno3['boxes'], facecolor='gold')
    setp(set03_zeno1['boxes'], hatch='.'); setp(set03_zeno2['boxes'], hatch='.'); setp(set03_zeno3['boxes'], hatch='.')

    setp(set01_vuvuzela1['boxes'], color='black'); setp(set01_vuvuzela2['boxes'], color='black'); setp(set01_vuvuzela3['boxes'], color='black')
    setp(set01_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set01_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set01_vuvuzela1['boxes'], hatch='\\'); setp(set01_vuvuzela2['boxes'], hatch='\\'); setp(set01_vuvuzela3['boxes'], hatch='\\')

    setp(set02a_vuvuzela1['boxes'], color='black'); setp(set02a_vuvuzela2['boxes'], color='black'); setp(set02a_vuvuzela3['boxes'], color='black')
    setp(set02a_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02a_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02a_vuvuzela1['boxes'], hatch='\\\\'); setp(set02a_vuvuzela2['boxes'], hatch='\\\\'); setp(set02a_vuvuzela3['boxes'], hatch='\\\\')

    setp(set02b_vuvuzela1['boxes'], color='black'); setp(set02b_vuvuzela2['boxes'], color='black'); setp(set02b_vuvuzela3['boxes'], color='black')
    setp(set02b_vuvuzela1['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela2['boxes'], facecolor='darkseagreen'); setp(set02b_vuvuzela3['boxes'], facecolor='darkseagreen')
    setp(set02b_vuvuzela1['boxes'], hatch='-'); setp(set02b_vuvuzela2['boxes'], hatch='-'); setp(set02b_vuvuzela3['boxes'], hatch='-')

    setp(set01_pung1['boxes'], color='black'); setp(set01_pung2['boxes'], color='black'); setp(set01_pung3['boxes'], color='black')
    setp(set01_pung1['boxes'], facecolor='steelblue'); setp(set01_pung2['boxes'], facecolor='steelblue'); setp(set01_pung3['boxes'], facecolor='steelblue')
    setp(set01_pung1['boxes'], hatch='x'); setp(set01_pung2['boxes'], hatch='x'); setp(set01_pung3['boxes'], hatch='x')

    setp(set02a_pung1['boxes'], color='black'); setp(set02a_pung2['boxes'], color='black'); setp(set02a_pung3['boxes'], color='black')
    setp(set02a_pung1['boxes'], facecolor='steelblue'); setp(set02a_pung2['boxes'], facecolor='steelblue'); setp(set02a_pung3['boxes'], facecolor='steelblue')
    setp(set02a_pung1['boxes'], hatch='o'); setp(set02a_pung2['boxes'], hatch='o'); setp(set02a_pung3['boxes'], hatch='o')

    setp(set02b_pung1['boxes'], color='black'); setp(set02b_pung2['boxes'], color='black'); setp(set02b_pung3['boxes'], color='black')
    setp(set02b_pung1['boxes'], facecolor='steelblue'); setp(set02b_pung2['boxes'], facecolor='steelblue'); setp(set02b_pung3['boxes'], facecolor='steelblue')
    setp(set02b_pung1['boxes'], hatch='*'); setp(set02b_pung2['boxes'], hatch='*'); setp(set02b_pung3['boxes'], hatch='*')


    ax.yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    ax.set_axisbelow(True)

    ax.set_xlim([0, 33])
    ax.set_xticks((5.5, 16.5, 27.5))
    ax.set_xticklabels(('1,000 clients', '2,000 clients', '3,000 clients'))
    ax.set_ylim([0.0, 50.0])
    ax.set_yticks([0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50])

    # Add a legend.
    ax.legend([
        set01_zeno1['boxes'][0],
        set02a_zeno1['boxes'][0],
        set02b_zeno1['boxes'][0],
        set03_zeno1['boxes'][0],
        set01_vuvuzela1['boxes'][0],
        set02a_vuvuzela1['boxes'][0],
        set02b_vuvuzela1['boxes'][0],
        set01_pung1['boxes'][0],
        set02a_pung1['boxes'][0],
        set02b_pung1['boxes'][0]
    ], [
        'zeno (no impediments)',
        'zeno (high delay, no failures)',
        'zeno (high loss, no failures)',
        'zeno (high network troubles, failures)',
        'vuvuzela (no impediments)',
        'vuvuzela (high delay, no failures)',
        'vuvuzela (high loss, no failures)',
        'pung (no impediments)',
        'pung (high delay, no failures)',
        'pung (high loss, no failures)'
    ],
    loc='upper center',
    ncol=3)
    ax.set_title("Memory Load per Server")

    plt.tight_layout()
    plt.xlabel("Number of clients")
    plt.ylabel("Used memory (GB)")

    # plt.savefig(os.path.join(sys.argv[1], "memory-used_servers.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "memory-used_servers.pdf"), bbox_inches='tight')


def compileLatencies():

    global metrics

    x_min = 100.0
    x_max = 0.0

    # Ingest data.

    # 1000 clients.

    set01_zeno_1000_Latencies = []
    with open(metrics["01"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_zeno_1000_Latencies.append(float(item))
    
    set01_vuvuzela_1000_Latencies = []
    with open(metrics["01"]["vuvuzela"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_vuvuzela_1000_Latencies.append(float(item))
    
    set01_pung_1000_Latencies = []
    with open(metrics["01"]["pung"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_pung_1000_Latencies.append(float(item))
    
    set02a_zeno_1000_Latencies = []
    with open(metrics["02a"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_zeno_1000_Latencies.append(float(item))
    
    set02a_vuvuzela_1000_Latencies = []
    with open(metrics["02a"]["vuvuzela"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_vuvuzela_1000_Latencies.append(float(item))
    
    set02a_pung_1000_Latencies = []
    with open(metrics["02a"]["pung"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_pung_1000_Latencies.append(float(item))
    
    set02b_zeno_1000_Latencies = []
    with open(metrics["02b"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_zeno_1000_Latencies.append(float(item))
    
    set02b_vuvuzela_1000_Latencies = []
    with open(metrics["02b"]["vuvuzela"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_vuvuzela_1000_Latencies.append(float(item))
    
    set02b_pung_1000_Latencies = []
    with open(metrics["02b"]["pung"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_pung_1000_Latencies.append(float(item))
    
    set03_zeno_1000_Latencies = []
    with open(metrics["03"]["zeno"]["1000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set03_zeno_1000_Latencies.append(float(item))
    

    # 2000 clients.

    set01_zeno_2000_Latencies = []
    with open(metrics["01"]["zeno"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_zeno_2000_Latencies.append(float(item))
    
    set01_vuvuzela_2000_Latencies = []
    with open(metrics["01"]["vuvuzela"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_vuvuzela_2000_Latencies.append(float(item))
    
    set01_pung_2000_Latencies = []
    with open(metrics["01"]["pung"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_pung_2000_Latencies.append(float(item))
    
    set02a_zeno_2000_Latencies = []
    with open(metrics["02a"]["zeno"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_zeno_2000_Latencies.append(float(item))
    
    set02a_vuvuzela_2000_Latencies = []
    with open(metrics["02a"]["vuvuzela"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_vuvuzela_2000_Latencies.append(float(item))
    
    set02a_pung_2000_Latencies = []
    with open(metrics["02a"]["pung"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_pung_2000_Latencies.append(float(item))
    
    set02b_zeno_2000_Latencies = []
    with open(metrics["02b"]["zeno"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_zeno_2000_Latencies.append(float(item))
    
    set02b_vuvuzela_2000_Latencies = []
    with open(metrics["02b"]["vuvuzela"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_vuvuzela_2000_Latencies.append(float(item))
    
    set02b_pung_2000_Latencies = []
    with open(metrics["02b"]["pung"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_pung_2000_Latencies.append(float(item))
    
    set03_zeno_2000_Latencies = []
    with open(metrics["03"]["zeno"]["2000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set03_zeno_2000_Latencies.append(float(item))
    

    # 3000 clients.

    set01_zeno_3000_Latencies = []
    with open(metrics["01"]["zeno"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_zeno_3000_Latencies.append(float(item))
    
    set01_vuvuzela_3000_Latencies = []
    with open(metrics["01"]["vuvuzela"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_vuvuzela_3000_Latencies.append(float(item))
    
    set01_pung_3000_Latencies = []
    with open(metrics["01"]["pung"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set01_pung_3000_Latencies.append(float(item))
    
    set02a_zeno_3000_Latencies = []
    with open(metrics["02a"]["zeno"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_zeno_3000_Latencies.append(float(item))
    
    set02a_vuvuzela_3000_Latencies = []
    with open(metrics["02a"]["vuvuzela"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_vuvuzela_3000_Latencies.append(float(item))
    
    set02a_pung_3000_Latencies = []
    with open(metrics["02a"]["pung"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02a_pung_3000_Latencies.append(float(item))
    
    set02b_zeno_3000_Latencies = []
    with open(metrics["02b"]["zeno"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_zeno_3000_Latencies.append(float(item))
    
    set02b_vuvuzela_3000_Latencies = []
    with open(metrics["02b"]["vuvuzela"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_vuvuzela_3000_Latencies.append(float(item))
    
    set02b_pung_3000_Latencies = []
    with open(metrics["02b"]["pung"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set02b_pung_3000_Latencies.append(float(item))
    
    set03_zeno_3000_Latencies = []
    with open(metrics["03"]["zeno"]["3000"]["Latencies"], newline='') as dataFile:
        reader = csv.reader(dataFile, delimiter=',')
        for row in reader:
            for item in row:
                if float(item) < x_min:
                    x_min = float(item)
                if float(item) > x_max:
                    x_max = float(item)
                set03_zeno_3000_Latencies.append(float(item))


    # Prepare CDF arrays.

    # 1000 clients.

    set01_zeno_1000_Latencies = np.sort(set01_zeno_1000_Latencies)
    set01_zeno_1000_CDF = np.array(range(len(set01_zeno_1000_Latencies))) / float(len(set01_zeno_1000_Latencies))
    set01_vuvuzela_1000_Latencies = np.sort(set01_vuvuzela_1000_Latencies)
    set01_vuvuzela_1000_CDF = np.array(range(len(set01_vuvuzela_1000_Latencies))) / float(len(set01_vuvuzela_1000_Latencies))
    set01_pung_1000_Latencies = np.sort(set01_pung_1000_Latencies)
    set01_pung_1000_CDF = np.array(range(len(set01_pung_1000_Latencies))) / float(len(set01_pung_1000_Latencies))

    set02a_zeno_1000_Latencies = np.sort(set02a_zeno_1000_Latencies)
    set02a_zeno_1000_CDF = np.array(range(len(set02a_zeno_1000_Latencies))) / float(len(set02a_zeno_1000_Latencies))
    set02a_vuvuzela_1000_Latencies = np.sort(set02a_vuvuzela_1000_Latencies)
    set02a_vuvuzela_1000_CDF = np.array(range(len(set02a_vuvuzela_1000_Latencies))) / float(len(set02a_vuvuzela_1000_Latencies))
    set02a_pung_1000_Latencies = np.sort(set02a_pung_1000_Latencies)
    set02a_pung_1000_CDF = np.array(range(len(set02a_pung_1000_Latencies))) / float(len(set02a_pung_1000_Latencies))

    set02b_zeno_1000_Latencies = np.sort(set02b_zeno_1000_Latencies)
    set02b_zeno_1000_CDF = np.array(range(len(set02b_zeno_1000_Latencies))) / float(len(set02b_zeno_1000_Latencies))
    set02b_vuvuzela_1000_Latencies = np.sort(set02b_vuvuzela_1000_Latencies)
    set02b_vuvuzela_1000_CDF = np.array(range(len(set02b_vuvuzela_1000_Latencies))) / float(len(set02b_vuvuzela_1000_Latencies))
    set02b_pung_1000_Latencies = np.sort(set02b_pung_1000_Latencies)
    set02b_pung_1000_CDF = np.array(range(len(set02b_pung_1000_Latencies))) / float(len(set02b_pung_1000_Latencies))

    set03_zeno_1000_Latencies = np.sort(set03_zeno_1000_Latencies)
    set03_zeno_1000_CDF = np.array(range(len(set03_zeno_1000_Latencies))) / float(len(set03_zeno_1000_Latencies))

    # 2000 clients.

    set01_zeno_2000_Latencies = np.sort(set01_zeno_2000_Latencies)
    set01_zeno_2000_CDF = np.array(range(len(set01_zeno_2000_Latencies))) / float(len(set01_zeno_2000_Latencies))
    set01_vuvuzela_2000_Latencies = np.sort(set01_vuvuzela_2000_Latencies)
    set01_vuvuzela_2000_CDF = np.array(range(len(set01_vuvuzela_2000_Latencies))) / float(len(set01_vuvuzela_2000_Latencies))
    set01_pung_2000_Latencies = np.sort(set01_pung_2000_Latencies)
    set01_pung_2000_CDF = np.array(range(len(set01_pung_2000_Latencies))) / float(len(set01_pung_2000_Latencies))

    set02a_zeno_2000_Latencies = np.sort(set02a_zeno_2000_Latencies)
    set02a_zeno_2000_CDF = np.array(range(len(set02a_zeno_2000_Latencies))) / float(len(set02a_zeno_2000_Latencies))
    set02a_vuvuzela_2000_Latencies = np.sort(set02a_vuvuzela_2000_Latencies)
    set02a_vuvuzela_2000_CDF = np.array(range(len(set02a_vuvuzela_2000_Latencies))) / float(len(set02a_vuvuzela_2000_Latencies))
    set02a_pung_2000_Latencies = np.sort(set02a_pung_2000_Latencies)
    set02a_pung_2000_CDF = np.array(range(len(set02a_pung_2000_Latencies))) / float(len(set02a_pung_2000_Latencies))

    set02b_zeno_2000_Latencies = np.sort(set02b_zeno_2000_Latencies)
    set02b_zeno_2000_CDF = np.array(range(len(set02b_zeno_2000_Latencies))) / float(len(set02b_zeno_2000_Latencies))
    set02b_vuvuzela_2000_Latencies = np.sort(set02b_vuvuzela_2000_Latencies)
    set02b_vuvuzela_2000_CDF = np.array(range(len(set02b_vuvuzela_2000_Latencies))) / float(len(set02b_vuvuzela_2000_Latencies))
    set02b_pung_2000_Latencies = np.sort(set02b_pung_2000_Latencies)
    set02b_pung_2000_CDF = np.array(range(len(set02b_pung_2000_Latencies))) / float(len(set02b_pung_2000_Latencies))

    set03_zeno_2000_Latencies = np.sort(set03_zeno_2000_Latencies)
    set03_zeno_2000_CDF = np.array(range(len(set03_zeno_2000_Latencies))) / float(len(set03_zeno_2000_Latencies))

    # 3000 clients.

    set01_zeno_3000_Latencies = np.sort(set01_zeno_3000_Latencies)
    set01_zeno_3000_CDF = np.array(range(len(set01_zeno_3000_Latencies))) / float(len(set01_zeno_3000_Latencies))
    set01_vuvuzela_3000_Latencies = np.sort(set01_vuvuzela_3000_Latencies)
    set01_vuvuzela_3000_CDF = np.array(range(len(set01_vuvuzela_3000_Latencies))) / float(len(set01_vuvuzela_3000_Latencies))
    set01_pung_3000_Latencies = np.sort(set01_pung_3000_Latencies)
    set01_pung_3000_CDF = np.array(range(len(set01_pung_3000_Latencies))) / float(len(set01_pung_3000_Latencies))

    set02a_zeno_3000_Latencies = np.sort(set02a_zeno_3000_Latencies)
    set02a_zeno_3000_CDF = np.array(range(len(set02a_zeno_3000_Latencies))) / float(len(set02a_zeno_3000_Latencies))
    set02a_vuvuzela_3000_Latencies = np.sort(set02a_vuvuzela_3000_Latencies)
    set02a_vuvuzela_3000_CDF = np.array(range(len(set02a_vuvuzela_3000_Latencies))) / float(len(set02a_vuvuzela_3000_Latencies))
    set02a_pung_3000_Latencies = np.sort(set02a_pung_3000_Latencies)
    set02a_pung_3000_CDF = np.array(range(len(set02a_pung_3000_Latencies))) / float(len(set02a_pung_3000_Latencies))

    set02b_zeno_3000_Latencies = np.sort(set02b_zeno_3000_Latencies)
    set02b_zeno_3000_CDF = np.array(range(len(set02b_zeno_3000_Latencies))) / float(len(set02b_zeno_3000_Latencies))
    set02b_vuvuzela_3000_Latencies = np.sort(set02b_vuvuzela_3000_Latencies)
    set02b_vuvuzela_3000_CDF = np.array(range(len(set02b_vuvuzela_3000_Latencies))) / float(len(set02b_vuvuzela_3000_Latencies))
    set02b_pung_3000_Latencies = np.sort(set02b_pung_3000_Latencies)
    set02b_pung_3000_CDF = np.array(range(len(set02b_pung_3000_Latencies))) / float(len(set02b_pung_3000_Latencies))

    set03_zeno_3000_Latencies = np.sort(set03_zeno_3000_Latencies)
    set03_zeno_3000_CDF = np.array(range(len(set03_zeno_3000_Latencies))) / float(len(set03_zeno_3000_Latencies))


    # Draw plots.

    fig, axes = plt.subplots(3, figsize=(10, 12))

    # 1000 clients.

    pung1, = axes[0].plot(set01_pung_1000_Latencies, set01_pung_1000_CDF, label='Pung (no impediments)', color='steelblue')
    pung2, = axes[0].plot(set02a_pung_1000_Latencies, set02a_pung_1000_CDF, label='Pung (high delay, no failures)', color='dodgerblue')
    pung3, = axes[0].plot(set02b_pung_1000_Latencies, set02b_pung_1000_CDF, label='Pung (high loss, no failures)', color='skyblue')

    zeno1, = axes[0].plot(set01_zeno_1000_Latencies, set01_zeno_1000_CDF, label='zeno (no impediments)', color='gold')
    zeno2, = axes[0].plot(set02a_zeno_1000_Latencies, set02a_zeno_1000_CDF, label='zeno (high delay, no failures)', color='khaki')
    zeno3, = axes[0].plot(set02b_zeno_1000_Latencies, set02b_zeno_1000_CDF, label='zeno (high loss, no failures)', color='sandybrown')
    zeno4, = axes[0].plot(set03_zeno_1000_Latencies, set03_zeno_1000_CDF, label='zeno (high network troubles, failures)', color='orange')

    vuvuzela1, = axes[0].plot(set01_vuvuzela_1000_Latencies, set01_vuvuzela_1000_CDF, label='Vuvuzela (no impediments)', color='darkseagreen')
    vuvuzela2, = axes[0].plot(set02a_vuvuzela_1000_Latencies, set02a_vuvuzela_1000_CDF, label='Vuvuzela (high delay, no failures)', color='limegreen')
    vuvuzela3, = axes[0].plot(set02b_vuvuzela_1000_Latencies, set02b_vuvuzela_1000_CDF, label='Vuvuzela (high loss, no failures)', color='olive')

    # 2000 clients.

    axes[1].plot(set01_pung_2000_Latencies, set01_pung_2000_CDF, color='steelblue')
    axes[1].plot(set02a_pung_2000_Latencies, set02a_pung_2000_CDF, color='dodgerblue')
    axes[1].plot(set02b_pung_2000_Latencies, set02b_pung_2000_CDF, color='skyblue')

    axes[1].plot(set01_zeno_2000_Latencies, set01_zeno_2000_CDF, color='gold')
    axes[1].plot(set02a_zeno_2000_Latencies, set02a_zeno_2000_CDF, color='khaki')
    axes[1].plot(set02b_zeno_2000_Latencies, set02b_zeno_2000_CDF, color='sandybrown')
    axes[1].plot(set03_zeno_2000_Latencies, set03_zeno_2000_CDF, color='orange')

    axes[1].plot(set01_vuvuzela_2000_Latencies, set01_vuvuzela_2000_CDF, color='darkseagreen')
    axes[1].plot(set02a_vuvuzela_2000_Latencies, set02a_vuvuzela_2000_CDF, color='limegreen')
    axes[1].plot(set02b_vuvuzela_2000_Latencies, set02b_vuvuzela_2000_CDF, color='olive')

    # 3000 clients.

    axes[2].plot(set01_pung_3000_Latencies, set01_pung_3000_CDF, color='steelblue')
    axes[2].plot(set02a_pung_3000_Latencies, set02a_pung_3000_CDF, color='dodgerblue')
    axes[2].plot(set02b_pung_3000_Latencies, set02b_pung_3000_CDF, color='skyblue')

    axes[2].plot(set01_zeno_3000_Latencies, set01_zeno_3000_CDF, color='gold')
    axes[2].plot(set02a_zeno_3000_Latencies, set02a_zeno_3000_CDF, color='khaki')
    axes[2].plot(set02b_zeno_3000_Latencies, set02b_zeno_3000_CDF, color='sandybrown')
    axes[2].plot(set03_zeno_3000_Latencies, set03_zeno_3000_CDF, color='orange')

    axes[2].plot(set01_vuvuzela_3000_Latencies, set01_vuvuzela_3000_CDF, color='darkseagreen')
    axes[2].plot(set02a_vuvuzela_3000_Latencies, set02a_vuvuzela_3000_CDF, color='limegreen')
    axes[2].plot(set02b_vuvuzela_3000_Latencies, set02b_vuvuzela_3000_CDF, color='olive')


    axes[0].xaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    axes[0].yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    axes[0].set_axisbelow(True)

    axes[1].xaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    axes[1].yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    axes[1].set_axisbelow(True)

    axes[2].xaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    axes[2].yaxis.grid(True, linestyle='-', which='major', color='lightgrey', alpha=0.5)
    axes[2].set_axisbelow(True)

    axes[0].set_xlim([x_min, x_max])
    axes[0].set_ylim([0.0, 1.0])
    axes[0].set_yticks((0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0))

    axes[1].set_xlim([x_min, x_max])
    axes[1].set_ylim([0.0, 1.0])
    axes[1].set_yticks((0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0))

    axes[2].set_xlim([x_min, x_max])
    axes[2].set_ylim([0.0, 1.0])
    axes[2].set_yticks((0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0))

    axes[0].set_title("1,000 clients")
    axes[1].set_title("2,000 clients")
    axes[2].set_title("3,000 clients")
    fig.suptitle(t="CDFs of End-to-End Transmission Latencies on Clients", y=1.02, fontsize=16)

    axes[0].set_xlabel("End-to-end transmission latency (seconds)")
    axes[0].set_ylabel("Fraction of messages transmitted")
    axes[1].set_xlabel("End-to-end transmission latency (seconds)")
    axes[1].set_ylabel("Fraction of messages transmitted")
    axes[2].set_xlabel("End-to-end transmission latency (seconds)")
    axes[2].set_ylabel("Fraction of messages transmitted")

    # Add legends.
    box = axes[2].get_position()
    pung_legend = fig.legend(handles=[pung1, pung2, pung3], loc="upper left", bbox_to_anchor=((box.x0 - 0.071), 0.), bbox_transform=plt.gcf().transFigure)
    zeno_legend = fig.legend(handles=[zeno1, zeno2, zeno3, zeno4], loc="upper center", bbox_to_anchor=((box.width / 2.0) + box.x0, 0.), bbox_transform=plt.gcf().transFigure)
    fig.gca().add_artist(pung_legend)
    fig.gca().add_artist(zeno_legend)
    fig.legend(handles=[vuvuzela1, vuvuzela2, vuvuzela3], loc="upper right", bbox_to_anchor=((box.x1 + 0.091), 0.), bbox_transform=plt.gcf().transFigure)

    fig.set_tight_layout(True)
    plt.tight_layout(pad=0.4, w_pad=0.5, h_pad=1.0)

    # plt.savefig(os.path.join(sys.argv[1], "e2e-transmission-latencies.pgf"), bbox_inches='tight')
    plt.savefig(os.path.join(sys.argv[1], "e2e-transmission-latencies.pdf"), bbox_inches='tight')



# Create all figures.

# Build bandwidth figures.
compileTrafficClients()
compileTrafficServers()

# Build load usage figures.
compileLoadCPUClients()
compileLoadCPUServers()
compileLoadMemClients()
compileLoadMemServers()

# Build message latencies figure.
compileLatencies()
