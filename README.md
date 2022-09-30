# BusPool
A bus pooling system for optimized public transport

## Structure

### Bus
Simulates real life bus activity using go routines.\
Buses publish their location and other data to an mqtt topic.\
Buses receive instructions through individual mqtt topics.\
Buses move along their routes and use the GraphHopper API to find an optimal route\
when receiving new instructions.

### Client
Simulates client requests.\
At a given rate sends requests to the dispatch server.\
Requests consist of a client location and destination which are chosen randomly from\
all the bus stops in the Berlin area.

### Rhino
A connector program that reads data from mqtt toics and publishes it to Kafka topics.\

### Lion
Dispatch API.\
Processes requests and sends instruction to buses via mqtt.\
Keeps track of bus locations by reading from the Kafka topic.

## Milestones

1. ~~Have buses produce to an mqtt server and have that data be written in kafka topic~~

2. ~~Have the data be visualized (e.g. with elasticsearch and kibana)~~

3. ~~Start accepting client requests and assigning them the nearest bus~~

4. As kibana can't be used client side, write a UI for a client to view bus activity and send requests.

