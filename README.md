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
```
cd map
```
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

## Backend
```
$ cd backend
$ go build -o main
$ ./main <filepath> <category column head> <date time column head> <max level for nanocubes> <max number of points in nanocube>
```
An example of running the binary:
```
$ ./main ./parser/crime2020.csv PrimaryType Date 20 10000
```
It will automatically run on http://localhost:8080


### Note: In order to make the system work, you may run the backend first then frontend 

## CSV file format

It must have headers named Latitude,Longitude which represent the lat lng for the spatial points. It must also have two columns represent the category for each point and date time for each point. The date time format would be like:
10/30/2020 03:51:41 PM. A typical example for csv is backend/parser/crime2020.csv


