import { useState, useMemo, useCallback, useRef } from "react";
import {
  GoogleMap, Rectangle
} from "@react-google-maps/api";
import {convertNanoCubeBoundsToNSEW} from "../utils"
type LatLngLiteral = google.maps.LatLngLiteral;
type MapOptions = google.maps.MapOptions;
const bounds = {
  north: 42.022585817,
  south: 41.0226,
  east: -87.8345,
  west: -88.925
}


const chicagoBounds = {
  north: 42.022585817,
  south: 42.022585817-0.424,
  east: -87.9345+0.424,
  west: -87.9345
}

const ncBounds20 = {
  "bounds":
  {"lng":-87.77115154101563,
  "lat":41.97298620617969,
  "width":0.0007286721738281321,
  "height":0.0007286721738281321},
  "count":1
}
const ncBounds1 = {"bounds":{"lng":-87.905227221,"lat":42.022535914,"width":0.37308015300000363,"height":0.37308015300000363},"count":1000}
const ncBounds5 = {"bounds":{"lng":-87.905227221,"lat":41.9992184044375,"width":0.023317509562500227,"height":0.023317509562500227},"count":5}
const ncBounds10 = {"bounds":{"lng":-87.59991358016602,"lat":41.755113226205076,"width":0.0007286721738281321,"height":0.0007286721738281321},"count":1}
const ncBounds2 = {"bounds":{"lng":-87.905227221,"lat":42.022535914,"width":0.18654007650000182,"height":0.18654007650000182},"count":218}
const rectangleOptions = {
  strokeWeight: 0.1,
  fillColor: "#Fa240c"
}
export default function Map() {
  const mapRef = useRef<GoogleMap>();
  const center = useMemo<LatLngLiteral>(
    () => ({ lat: 41.8337329, lng: -87.7319639 }),
    []
  );
  const options = useMemo<MapOptions>(
    () => ({
      disableDefaultUI: true,
      clickableIcons: false,
    }),
    []
  );
  const onLoad = useCallback((map) => (mapRef.current = map), []);

  useEffect(() => {
    // Update the document title using the browser API
  });
  return (
    <div className="container">
      <div className="controls">
        <h1>Map</h1>
      </div>
      <div className="map">
        <GoogleMap
          zoom={10}
          center={center}
          mapContainerClassName="map-container"
          options={options}
          onLoad={onLoad}
        >
          <Rectangle
              options = {rectangleOptions}
              bounds={convertNanoCubeBoundsToNSEW(ncBounds1)}
            />
            <Rectangle
              options = {rectangleOptions}
              bounds={convertNanoCubeBoundsToNSEW(ncBounds2)}
            />
            <Rectangle
              options = {rectangleOptions}
              bounds={convertNanoCubeBoundsToNSEW(ncBounds5)}
            />
            
        </GoogleMap>
      </div>
    </div>
  );
}

