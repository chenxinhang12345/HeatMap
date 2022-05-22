import { useState, useMemo, useCallback, useRef, useEffect } from "react";
import {
  GoogleMap, Rectangle
} from "@react-google-maps/api";
import {convertNanoCubeBoundsToNSEW, getAllData, getTypes} from "../utils"
import InputLabel from '@mui/material/InputLabel';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import Select, { SelectChangeEvent } from '@mui/material/Select';

type LatLngLiteral = google.maps.LatLngLiteral;
type MapOptions = google.maps.MapOptions;

export default function Map() {
  const mapRef = useRef<GoogleMap>();
  const [zoom, setZoom] =useState <number | undefined>(10);
  const [bounds, setBounds] = useState(null);
  const [data, setData] = useState([]);
  const [types, setTypes] = useState({});
  const [currentType, setCurrentType] = useState('-1');
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
  const onBoundsChanged = () =>{
    setZoom(mapRef.current?.getZoom());
    setBounds(mapRef.current?.getBounds());
  }
  const onTypeChanged = (e :SelectChangeEvent) =>{
    console.log(e.target.value);
    setCurrentType(e.target.value);
  }
  useEffect (()=>{
    console.log('bounds', bounds);
    var values = [{h:0,j:0},{h:0,j:0}]
    if (bounds){
    console.log(values);
    console.log(values[1].h)
    console.log("bounds", bounds);
     values = Object.values(bounds);
     
    }
    const minLat = values[0].h;
    const maxLat = values[0].j;
    const minLng = values[1].h;
    const maxLng = values[1].j;
    console.log('currentType', currentType);
    getAllData(minLat,maxLat,minLng,maxLng,zoom,currentType).then(
      (res) =>{
        setData(res.data.data);
      }
    )
    
  },[bounds, currentType])

  useEffect (
    ()=>{
      console.log('types', types);
      getTypes().then(
        (res) =>{
          setTypes(res.data.data);
        }
      )
    },[]
  )
  return (
    <div className="container">
      <div className="controls">
        <h1>Crime Map</h1>
        <div>
        <FormControl variant="filled" sx={{ m: 1, minWidth: 120 }}>
        <InputLabel id="types">Types</InputLabel>
        <Select
          labelId="types"
          id="types"
          value={currentType}
          onChange={onTypeChanged}
          label="Type"
        >
          <MenuItem value={-1}>
            <em>All</em>
          </MenuItem>
          {
           Object.keys(types).map(
             (key, i) =>{
               return(
               <MenuItem value={types[key]}>{key}</MenuItem>
               )
             }
           )
          }
        </Select>
      </FormControl>
        </div>
      </div>
      <div className="map">
        <GoogleMap
          zoom={10}
          center={center}
          mapContainerClassName="map-container"
          options={options}
          onZoomChanged={onBoundsChanged}
          onDragEnd = {onBoundsChanged}
          onLoad={onLoad}
        >
          {
            data?.map((cube,index)=>{
              const {rectBounds, rectOpacity} = convertNanoCubeBoundsToNSEW(cube);
              const rectangleOptions = {
                strokeWeight: 0.1,
                fillColor: "#Fa240c",
                fillOpacity: rectOpacity
              }
              return (
                <Rectangle options = {rectangleOptions} bounds={rectBounds}/>
              )
            })
          }
            
        </GoogleMap>
      </div>
    </div>
  );
}

