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
todo

## Codebase Organization
todo

## File List
todo

## Description
todo

## Limitations and Improvements
todo


