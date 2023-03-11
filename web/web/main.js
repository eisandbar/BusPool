import './style.css';
import {Map, View} from 'ol';
import {Tile as TileLayer, Vector as VectorLayer} from 'ol/layer';
import * as olProj from 'ol/proj';
import VectorSource from 'ol/source/Vector';
import OSM from 'ol/source/OSM';
import Feature from 'ol/Feature';
import Point from 'ol/geom/Point';
import {
    Circle as CircleStyle,
    Fill,
    Stroke,
    Style,
  } from 'ol/style';
import {io} from 'socket.io-client'

const style = new Style({
    image: new CircleStyle({
        radius: 7,
        fill: new Fill({color: 'blue'}),
        stroke: new Stroke({
            color: 'black',
            width: 2,
        })
    })
})

let userStyle = new Style({
    image: new CircleStyle({
        radius: 7,
        fill: new Fill({color: 'red'}),
        stroke: new Stroke({
            color: 'black',
            width: 2,
        })
    })
})

var features = []
for (let i= 0; i<200; i++) {
    features[i] = new Feature({
        geometry: new Point(olProj.fromLonLat([13.415679931640627, 52.51371369804256]))
    })
}

const map = new Map({
  target: 'map',
  layers: [
    new TileLayer({
      source: new OSM()
    }),
    new VectorLayer({
        source: new VectorSource({
            features: features
        }),
        style: style,
    }),
    new VectorLayer({
        source: new VectorSource({
            features: [features[0]]
        }),
        style: userStyle,
    }),
  ],
  view: new View({
    center: olProj.fromLonLat([13.415679931640627, 52.51371369804256]),
    zoom: 12
  })
});

const socket = io("localhost:3000")

socket.on("update location", message => {
    console.log(message)
    for (let i = 0; i < message.length; i++) {
        if (message[i]) {
            features[i].getGeometry().setCoordinates(olProj.fromLonLat([message[i].Lon, message[i].Lat]))
        }
    }
})