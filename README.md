# HeatMap
This is an interactive heatmap backend built primarily in Golang based on Nanocubes structure.

## Purpose
The purpose of this system is to create a heatmap for visualization of large spatial temporal points with their own categories, subject to minimize the memory usage for the backend.

## Deployment
<!-- ## Test
To test validity of the nanocubes function:
```
$ cd backend/parser
$ go test -run TestNanoCubeFromBigFile
```
It will load all the data points from the crime2020.csv which is all the spatial crime points in Chicago City, and form nanocubes data structure.  -->
```
git clone https://github.com/chenxinhang12345/HeatMap
```
### Backend
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
### Frontend
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

### Note: In order to make the system work, you may run the backend first then frontend 

### CSV file format

It must have headers named Latitude,Longitude which represent the lat lng for the spatial points. It must also have two columns represent the category for each point and date time for each point. The date time format would be like:
10/30/2020 03:51:41 PM. A typical example for csv is backend/parser/crime2020.csv

## Usage
After we open the frontend in http://localhost:3000, we can see a webpage like this.
<img width="1327" alt="image" src="https://user-images.githubusercontent.com/20518726/171042740-0b4394d8-5db6-477c-af4a-c955e9389697.png">
The selection field on the left with `THEFT` on it is the current crime type or category of spatial points. User can select on this field. The start time and end time of spatial points we wan to show can also be adjusted. Just need to adjust the time on `Start date time` field and `End date time` field, the heatmap will autimatically change by fetching the result from the backend.

## Requirements and Dependencies
### Dependencies
`go 1.15.2` or above
Other go dependencies are in go.mod file. Dependencies will be automatically installed after running `go run main.go`.
`npm 8.1.4` or above
Other npm dependencies will be automatically installed after running `npm install`.
### System requirements
#### MacOS
```
System Version: macOS 12.0.1
Kernel Version: Darwin 21.1.0
```
#### Linux
```
System version: Ubuntu 20.04
```
### Memory
At least 2GB memory.

## System Architecture
<img width="853" alt="image" src="https://user-images.githubusercontent.com/20518726/171046733-be0e5c36-e684-4c34-ae1b-5ba95a6769dc.png">
The system parse the data from a csv file and construct Nanocubes structure in Golang. It also start a backend server for the querys from the frontend. The frontend will query the data from the backend as the user interact with the map. 

## Codebase Organization
    .
    ├── README.md # readme markdown
    ├── backend # The whole backend system
    │   ├── go.mod # go mod file for golang dependency install
    │   ├── go.sum # dependency verification
    │   ├── main.go # entry of the backend system
    │   ├── nanocube # module for the Nanocubes in system architecture
    │   │   ├── nanocube.go # nanocubes data structure implementation
    │   │   └── nanocube_test.go # test of some of the functions
    │   ├── parser # module for doing parsing in system architecture
    │   │   ├── crime2020.csv # an example csv file used for demo
    │   │   ├── parser.go # parser implementation
    │   │   └── parser_test.go # end to end test for constructing nanocubes from a csv file
    │   ├── server # Golang server in system architecture
    │   │   ├── controllers # api implementation
    │   │   │   └── map.go # api methods for querys
    │   │   └── models # data models
    │   │       └── cube.go # struct and json definition for rectangles returned to frontend
    │   └── utils
    │       └── utils.go # golang utilities
    └── map
        ├── README.md
        ├── components
        │   └── map.tsx # map component(primary implementation in this file)
        ├── env.d.ts
        ├── next-env.d.ts
        ├── next.config.js
        ├── package-lock.json
        ├── package.json
        ├── pages
        │   ├── _app.tsx
        │   └── index.tsx # entry
        ├── public
        │   └── favicon.ico
        ├── styles
        │   └── globals.css
        ├── tsconfig.json
        ├── tsconfig.tsbuildinfo
        ├── utils
        │   └── index.tsx # utilities for Typescript
        └── yarn.lock


## File List
- [./backend/main.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/main.go) Entry of the backend system
- [./backend/nanocube/nanocube.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/nanocube/nanocube.go) Nanocubes data structure implementation
- [./backend/nanocube/nanocube_test.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/nanocube/nanocube_test.go) Tests of some of the functions in Nanocubes implementation
- [./backend/parser/parser.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/parser/parser.go) Parser implementation
- [./backend/parser/parser_test.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/parser/parser_test.go) End to end test for constructing nanocubes from a csv file
- [./backend/server/controllers/map.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/server/controllers/map.go) API methods for querys
- [./backend/server/models/cube.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/server/models/cubes.go) Struct and json definition for rectangles
- [./backend/server/utils/utils.go](https://github.com/chenxinhang12345/HeatMap/blob/main/backend/server/utils/utils.go) Golang utility functions for the backend system
- [./map/components/map.tsx](https://github.com/chenxinhang12345/HeatMap/blob/main/map/components/map.tsx) Map component(primary implementation in this file for frontend)
- [./map/utils/index.tsx](https://github.com/chenxinhang12345/HeatMap/blob/main/map/utils/index.tsx) Utility functions for Typescript

## Description
The system parse the data from a csv file and construct Nanocubes structure in Golang using algorithm described in [this paper](http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.696.7905&rep=rep1&type=pdf). It also start a backend server using web framework Gin for the querys from the frontend. The frontend will query the data from the backend as the user interact with the map. The api for quering the backend will be as follows:
http://127.0.0.1:8080/cubes/time?minLat=`minimum latitude`&maxLat=`maximum latitude`.92650702568175&minLng=`minimum longitutde`&maxLng=`maximum longitude`&zoom=`zooming level`&type=`type index`&startTime=`start time in seconds`&endTime=`end time in seconds`
Example: http://127.0.0.1:8080/cubes/time?minLat=41.7387437848875&maxLat=41.92650702568175&minLng=-87.79987031532278&maxLng=-87.63404542518606&zoom=12&type=-1&startTime=1500000000&endTime=1600000000

Example data response:
```json
{
  "data": [
    {
      "n": 41.923832594570314,
      "s": 41.920841584890624,
      "e": -87.79680524285156,
      "w": -87.79979625253125,
      "count": 1,
      "opacity": 0.15049712002086615
    },
    {
      "n": 41.91186855585156,
      "s": 41.90887754617187,
      "e": -87.79381423317187,
      "w": -87.79680524285156,
      "count": 1,
      "opacity": 0.15049712002086615
    }
]
}
```
Then the frontend will render these recatngles using their opacity values and their location properties as shown in figure previously presented.
## Limitations and Improvements
There are several limitations and imporvements for this system.
- For large amount of data, the constructing time is longer than we thought. One possible reason is that the garbage collection time is more than we expected. In order to decrease the number of GCs we could add a flag on summary to indicate wehther it is shared or not to prevent future clean up. This would cost a little extra memory usage but would decrease the number of GCs.
- In the frontend implementation, it only support single category selection. Users might want to check multiple categories. It would be better if we use checkboxs instead of selection.
- For initial location of the map, it should be adaptive according to the input csv file. The initial location might be in the center of the Nanocube region instead of hardcoding.


