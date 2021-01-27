[![CircleCI](https://circleci.com/gh/giantswarm/azure-scheduled-events.svg?style=shield)](https://circleci.com/gh/giantswarm/azure-scheduled-events)

# Azure Scheduled Events

This app is meant to be run on node pool instances for azure >14.0.0 clusters.
It listens to [Azure Scheduled Events](https://docs.microsoft.com/en-us/azure/virtual-machines/linux/scheduled-events) and
it automatically drains the node before it gets terminated.
