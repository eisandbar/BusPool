# Setting up the buses

Setting up the base containers is as easy as running

```
docker compose --profile base up
```
It sets up MQTT and Kafka for communication.\
It starts a GraphHopper server for pathing.\
It starts the bus, rhino and lion go services.\
\
If you want to start simulating client requests run
```
docker compose up client
```

# Setting up kafka connector and elastic

## Starting the containers
To setup the kafka connector we first want to use a custom password for elastic
```
export EP=your_password
```
and then we can start the containers
```
docker compose --profile connect-elastic up
```
After the containers set up you will need to go to the url given by kibana.\
Here instead of using an enrollment token click 'Configure manually'.\
When prompted to enter an address enter 
```
http://es-node01:9200
```
This is the address by which kibana can always find the elastic container.\

## Connecting the data
To view our data on our map we need to do 2 things.\
First we will create our index.\
Open the console in kibana and run
```
PUT bus-positions-elastic
{
  "mappings": {
    "dynamic_templates": [
      {
        "dates": {
          "mapping": {
            "format": "epoch_millis",
            "type": "date"
          },
          "match": "*time"
        }
      },
      {
        "locations": {
          "match": "*location",
          "mapping": {
            "type": "geo_point"
          }
        }
      }
    ],
    "properties": {
      "id": {
        "type": "long"
      },
      "location": {
        "type": "geo_point"
      },
      "time": {
        "type": "date"
      }
    }
  }
}
```

What this does is it maps fields in our kafka data to known types like date and geo_point.\
\
Next in a terminal send this request, remembering to change connection.password to your password
```
curl -X PUT -H "Content-Type: application/json" --data '
{
  "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
  "type.name": "_doc",
  "topics": "bus-positions-elastic",
  "consumer.override.auto.offset.reset": "latest",
  "key.ignore": "true",
  "schema.ignore": "true",
  "name": "bus-positions-elastic",
  "value.converter.schemas.enable": "false",
  "connection.url": "http://es-node01:9200",
  "connection.username": "elastic",
  "connection.password": "your_password",
  "value.converter": "org.apache.kafka.connect.json.JsonConverter",
  "key.converter": "org.apache.kafka.connect.storage.StringConverter"
}' localhost:8083/connectors/bus-positions-elastic/config

```
This creates a connection between the kafka topic and your index in elastic.\
Now if you refresh your index it should contain the data from kafka.\

## Creating a map

To create a map we first need a data view.\
Go to Kibana / Data Views and create a view that includes the index we created.\
Be sure to include the timestamp field time.\
\
Now go to Maps and click Create map.\
Click Add layer, Documents and select the view you just created.\
Save the layer and the map.\
If you don't see anything try clicking on the calendar and your time window if your data is a bit old.\

## Security

This doc uses elastic without security for simplicity, but adding security is pretty simple and is explained [here](https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html#docker-compose-file)
