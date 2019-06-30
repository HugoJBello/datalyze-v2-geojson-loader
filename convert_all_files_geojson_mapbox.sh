
DIR = "data/raw_data"
cd data/raw_data
for FILE in $(find -name '*geojson.json') # cycles through all files in directory (case-sensitive!)
do
    echo "converting file: $FILE..."
    FILENEW=`echo $FILE | sed "s/_geojson.json/.mbtiles/"` # replaces old filename
    docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 -f -o "$FILENEW" "$FILE"

    FILENEWCOMPRESSED=`echo $FILE | sed "s/_geojson.json/_mapbox.tar.xz/"` # replaces old filename
    tar -cJf $FILENEWCOMPRESSED "$FILENEW"


done
exit