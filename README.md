# Resource-Aware Least Busy (RALB) - Load Balancing for Proxmox VMs (Golang Implementation)

## Overview

This repository is a part of my research __(the paper is still in draft, the complete paper will be released soon after paper trial!)__. The topic is **"Implementation of RALB in Proxmox Virtual Environment"**. 

This repo provides a Golang implementation of the **Resource-Aware Least Busy (RALB)** load balancing strategy, which is inspired by the work of Bouflous et al. (2023). The link: https://www.igi-global.com/article/resource-aware-least-busy-ralb-strategy-for-load-balancing-in-containerized-cloud-systems/328094

However, instead of focusing on containerized environments, this version targets virtual machines (VMs). The **RALB** strategy for Proxmox optimizes the allocation of incoming requests across VMs based on resource, ensuring that the least busy VM receives more traffic, resulting in better resource utilization.

## Features

- **Load Balancing for VMs**: Implements the RALB algorithm to balance load across VMs in Proxmox, considering both resource usage (CPU, memory and bandwidth).
- **Proxmox VE Integration**: Uses the Proxmox VE API to retrieve real-time VM resource usage information.
- **Golang Implementation**: Written in Go (Golang), optimized for performance.
- **HAProxy Integration**: Modifies the HAProxy configuration (`haproxy.cfg`) to adjust load balancing priority for VMs dynamically.
- **Operation**: Continuously monitors and updates VM load distribution in a loop, with adjustments every 1000ms (can be changed).

## Prerequisites

- Go version 1.18 or higher
- Proxmox VE API access and credentials
- HAProxy installed and configured

## Constraints

- Currently implemented in Proxmox VE
- Currently supported to be integrated with HAProxy