## Issues with connect topics

Connect automatically tries to create topics with replication-factor=3, so with only one broker you have to create the topic manually.

To do this you can start the broker container and exec into it with 
```
docker exec -it broker /bin/sh
```
and create each topic with kafka-topics
```
kafka-topics --create --topic connect-config --partitions 1 --replication-factor 1 --if-not-exists --bootstrap-server broker:9092
```
Connect might still complain that the topics are configure incorrectly (like cleanup.policy=delete instead of compact)
You can fix this with kafka-configs
```
kafka-configs --bootstrap-server broker:9092 --alter --entity-type topics --entity-name connect-offsets --add-config cleanup.policy=compact
```

## Issues with connect plugins

Had problems of adding connectors to the kafka-connect container.
Connectors have to be installed at start and cannot be added later.

Solution:
Use a custom build with Dockerfile 
```Dockerfile
FROM confluentinc/cp-kafka-connect-base:7.0.1
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:latest
```
and add the directory where these connectors will be installed by confluent-hub to the plugin_path
```yaml
CONNECT_PLUGIN_PATH: "/usr/share/java,/usr/share/confluent-hub-components"
```


## Connecting to Elastic

For dev make sure xpack.security.http.ssl=false in elasticsearch.yml

1. Connect kibana
    1. Copy the cert from you elastic container
    ```
    docker cp es-node01:/usr/share/elasticsearch/config/certs/http_ca.crt .
    ```
    2. Find the IP of your elastic container
    On the kibana page localhost:5601 enter configure manually and enter the elastic container IP
    3. It will ask you to authorize. Create a user with the elastic api
    ```
    curl --cacert http_ca.crt -u elastic -XPOST "http://localhost:9200/_security/user/[your user]" \
    -H 'Content-Type: application/json' -d '{"password": "[your password]", "roles": ["kibana_system"]}'
    ```
    and then authorize on the kibana page.
    4. Now you need to authorize again, this time as elastic
2. Create the connection
```
curl -X PUT -H "Content-Type: application/json" --data '
{
"connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
"type.name": "_doc",
"topics": "[your topic]",
"consumer.override.auto.offset.reset": "latest",
"key.ignore": "true",
"schema.ignore": "true",
"name": "[connection name]",
"value.converter.schemas.enable": "false",
"connection.url": "http://es-node01:9200",
"connection.username": "elastic",
"connection.password": "[your password]",
"value.converter": "org.apache.kafka.connect.json.JsonConverter",
"key.converter": "org.apache.kafka.connect.storage.StringConverter"
}' localhost:8083/connectors/[connection name]/config
```