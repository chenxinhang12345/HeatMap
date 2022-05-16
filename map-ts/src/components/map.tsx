import { useState, useMemo, useCallback, useRef } from "react";
import {
  GoogleMap, Rectangle
} from "@react-google-maps/api";
import "../styles/globals.css";
import { on } from "events";

type LatLngLiteral = google.maps.LatLngLiteral;
type MapOptions = google.maps.MapOptions;
const bounds = {
  north: 42.022585817,
  south: 41.0226,
  east: -87.9345,
  west: -88.025
}
const rectangleOptions = {
  strokeWeight: 0.1,
  fillColor: "#fa240c"
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
  // @ts-ignore
  const onLoad = useCallback((map) => (mapRef.current = map), []);
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
          onLoad = {onLoad}
        >
            <Rectangle
              options = {rectangleOptions}
              bounds={bounds}
            />
        </GoogleMap>
      </div>
    </div>
  );
}