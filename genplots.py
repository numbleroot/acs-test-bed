#!/usr/bin/env python3

import sys
import os
import csv
import numpy as np
import matplotlib
from matplotlib import pyplot as plt

# matplotlib.rcParams['text.usetex'] = True
matplotlib.rcParams['font.size'] = 12


def compileTraffic(outgoing, dataFile, outputFile):

    data = []
    with open(dataFile, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            data.append(list(map(int, row)))

    _, ax = plt.subplots()
    # ax.set_xlim([x_min, x_max])
    # ax.set_ylim([y_min, y_max])
    # ax.set_aspect(1, share=True)
    ax.boxplot(data)

    plt.grid()
    plt.tight_layout()

    plt.xlabel("Experiment time (seconds)")

    if outgoing:
        plt.ylabel("Outgoing traffic (bytes)")
    else:
        plt.ylabel("Incoming traffic (bytes)")

    plt.savefig(outputFile, bbox_inches = 'tight', dpi = 200)


def compileMemory(dataFile, outputFile):

    data = []
    with open(dataFile, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            data.append(list(map(float, row)))

    _, ax = plt.subplots()
    # ax.set_xlim([x_min, x_max])
    ax.set_ylim([0.0, 100.0])
    # ax.set_aspect(1, share=True)
    ax.boxplot(data)

    plt.grid()
    plt.tight_layout()

    plt.xlabel("Experiment time (seconds)")
    plt.ylabel("Used memory (percentage)")

    plt.savefig(outputFile, bbox_inches = 'tight', dpi = 200)


def compileLoad(dataFile, outputFile):

    data = []
    with open(dataFile, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            data.append(list(map(float, row)))

    _, ax = plt.subplots()
    # ax.set_xlim([x_min, x_max])
    ax.set_ylim([0.0, 100.0])
    # ax.set_aspect(1, share=True)
    ax.boxplot(data)

    plt.grid()
    plt.tight_layout()

    plt.xlabel("Experiment time (seconds)")
    plt.ylabel("CPU(s) utilization (percentage)")

    plt.savefig(outputFile, bbox_inches = 'tight', dpi = 200)


def compileLatencies(dataFile, outputFile):
    
    data = []
    with open(dataFile, newline='') as csvfile:
        reader = csv.reader(csvfile, delimiter=',')
        for row in reader:
            data.append(list(map(float, row)))
    
    flat_data = [latency for message in data for latency in message]
    y_max = np.ceil(max(flat_data))

    _, ax = plt.subplots()
    ax.set_ylim([0, y_max])
    # ax.set_aspect(1, share=True)
    ax.boxplot(data)

    plt.grid()
    plt.tight_layout()

    plt.xlabel("Message (sequence number)")
    plt.ylabel("End-to-end latency (seconds)")

    plt.savefig(outputFile, bbox_inches = 'tight', dpi = 200)


clientMetricsPath = os.path.join(sys.argv[1], "clients")
mixMetricsPath = os.path.join(sys.argv[1], "mixes")

clientSentBytesDataFile = os.path.join(clientMetricsPath, "sent-bytes_per_second.boxplot")
clientSentBytesOutputFile = os.path.join(clientMetricsPath, "sent-bytes_per_second.png")
mixSentBytesDataFile = os.path.join(mixMetricsPath, "sent-bytes_per_second.boxplot")
mixSentBytesOutputFile = os.path.join(mixMetricsPath, "sent-bytes_per_second.png")

clientRecvdBytesDataFile = os.path.join(clientMetricsPath, "recvd-bytes_per_second.boxplot")
clientRecvdBytesOutputFile = os.path.join(clientMetricsPath, "recvd-bytes_per_second.png")
mixRecvdBytesDataFile = os.path.join(mixMetricsPath, "recvd-bytes_per_second.boxplot")
mixRecvdBytesOutputFile = os.path.join(mixMetricsPath, "recvd-bytes_per_second.png")

clientMemoryDataFile = os.path.join(clientMetricsPath, "memory_per_second.boxplot")
clientMemoryOutputFile = os.path.join(clientMetricsPath, "memory_per_second.png")
mixMemoryDataFile = os.path.join(mixMetricsPath, "memory_per_second.boxplot")
mixMemoryOutputFile = os.path.join(mixMetricsPath, "memory_per_second.png")

clientLoadDataFile = os.path.join(clientMetricsPath, "load_per_second.boxplot")
clientLoadOutputFile = os.path.join(clientMetricsPath, "load_per_second.png")
mixLoadDataFile = os.path.join(mixMetricsPath, "load_per_second.boxplot")
mixLoadOutputFile = os.path.join(mixMetricsPath, "load_per_second.png")

latenciesDataFile = os.path.join(clientMetricsPath, "latency_per_message.boxplot")
latenciesOutputFile = os.path.join(clientMetricsPath, "latency_per_message.png")


# Create figures for outgoing metrics.
# compileTraffic(True, clientSentBytesDataFile, clientSentBytesOutputFile)
# compileTraffic(True, mixSentBytesDataFile, mixSentBytesOutputFile)

# Create figures for incoming metrics.
# compileTraffic(False, clientRecvdBytesDataFile, clientRecvdBytesOutputFile)
# compileTraffic(False, mixRecvdBytesDataFile, mixRecvdBytesOutputFile)

# Create figures for memory metrics.
# compileMemory(clientMemoryDataFile, clientMemoryOutputFile)
# compileMemory(mixMemoryDataFile, mixMemoryOutputFile)

# Create figures for load metrics.
# compileLoad(clientLoadDataFile, clientLoadOutputFile)
# compileLoad(mixLoadDataFile, mixLoadOutputFile)

# Clients-only: create figures for message latency metrics.
compileLatencies(latenciesDataFile, latenciesOutputFile)
