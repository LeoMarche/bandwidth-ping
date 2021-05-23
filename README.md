# bandwidth-ping

This repository contains an experiment I made when I realized that bandwidth had influence on latency.
This isn't really accurate but it's funny!

## Principle

If we model the time packet spend in cables as a fixed propagation time (due to length and speed in cable) ans a throughput time (due to size of packet and bandwidth), we realize that we can deduce a bandwidth from a few RTT measures using ping, that's really fun !

## Using it

1. Clone this repository
2. go build in the root directory
3. ```./bandwidth-ping <target> <number of samples> <number of ping samples>```
