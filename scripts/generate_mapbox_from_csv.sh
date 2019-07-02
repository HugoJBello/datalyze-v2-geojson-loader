#!/bin/bash

FILE_INPUT=$1
FILE_INPUT_NAME=aaa
FILE_OUTPUT=data/processed_data/$FILE_INPUT_NAME.mbtiles

#script='docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe  tippecanoe -z12 -f -o '$FILE_OUTPUT $FILE_INPUT

echo $script
eval $script 

output=$( go run . -generate-from-csv -path $FILE_INPUT -output data/processed_data -output-mapbox )
echo $output
eval $output
#docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 -f -o data/processed_data/example.mbtiles data/processed_data/example_geojson.json
