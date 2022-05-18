# HeatMap
This is an interactive heatmap backend built primarily in Golang based on Nanocubes structure.

## Test
To test validity of the nanocubes function:
```
$ cd backend/parser
$ go test -run TestNanoCubeFromBigFile
```
It will load all the data points from the crime2020.csv which is all the spatial crime points in Chicago City, and form nanocubes data structure. 

## Frontend
Put .env.local in map folder 

```
NEXT_PUBLIC_GOOGLE_MAPS_API_KEY="your google map api key"
NEXT_PUBLIC_SERVER_HOST="http://<server_ip>:<server_port>"
```

Then run 
```
$ npm install
$ npm run dev
```
