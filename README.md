# BusPool
A bus pooling system for optimized public transport

## Structure

Package bus simulates bus movement and reporting
produces data into an mqtt server and receives orders

Package client is a mock of expected client behavior.
At random locations and irregular intervals requests will be created to simulate real-world situations

Package giraffe handles rest api requests

Package rhino subscribes to an mqtt server and writes data into a kafka topic

Package gazelle subscribes to kafka and enriches the data for visualization

Package hippo subscribes to kafka and stores data in a db

Package lion processes client requests and sends orders to buses

## Milestones

1. Have buses produce to an mqtt server and have that data be written in kafka topic

2. Have the data be visualized (e.g. with elasticsearch and kibana)

3. Start accepting client requests and assigning them the nearest bus

4. Optimize the dispatch algorithm

