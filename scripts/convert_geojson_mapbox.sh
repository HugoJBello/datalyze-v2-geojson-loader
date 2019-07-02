#!/bin/bash

FILE_INPUT=$2
FILE_OUTPUT=$1

docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 -f -o $FILE_OUTPUT $FILE_INPUT


#docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 -f -o data/raw_data/polygon_cusecs_geojson.mbtiles data/raw_data/polygon_cusecs_geojson_crs.json

