## Preface

After some thought I decided that it was more logical to have the buses themselves query the \
graphHopper api instead of having the request server send instructions

## Assumptions

The first assumption is that most of the passengers are already on board the bus. \
This means that we can handle destination points as drop-offs.\
(Otherwise we would have to make sure that we picked up a passenger before we could drop him off)

The second assumption is that it is not too much slower to always pickup a new passenger first before \
continuing dropping off passengers. \
This means that when the list of pending clients isn't empty we can stop what we were doing and first \
pick them up. \
This works together with the first assumption in that there will always be only one or two pending \
clients at most. \
A potential problem is that a bus could be stuck in an endless loop of picking up clients and not \
dropping any of them off. This can be solved by using a limited capacity. \
This shouldn't be a problem with a decent bus/client-rps ratio

## Routing Logic

When receiving a new request (client, dest) first find the path from the current location to the client.
```go
client, dest := newRequest
clients = append(clients, client)
destinations = append(destinations, dest)
points = append([]point{current}, clients...)
```
Then, each time we reach a client we remove them from the clients list. \
When the list of clients is empty, we find the optimal path to reach all the destinations
```go
points = append([]point{current}, destinations...)
```
And each time we reach a destination we remove it from the destinations list. \

Depending on the bus capacity it might be worth to keep a map of the clients and destinations for faster \
lookup when checking if we reached a point. But for smaller capacities it's not worth the memory allocation.