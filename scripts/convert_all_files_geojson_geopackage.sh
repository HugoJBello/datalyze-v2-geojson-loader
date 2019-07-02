
DIR="../data/raw_data"
cd ../data/raw_data
for FILE in $(find -name '*.json') # cycles through all files in directory (case-sensitive!)
do
    echo "converting file: $FILE..."
    FILENEW=`echo $FILE | sed "s/geojson.json/.gpkg/"` # replaces old filename
    docker run -v $(pwd):/data geodata/gdal ogr2ogr -f "GPKG" "$FILENEW" "$FILE"

    FILENEWCOMPRESSED=`echo $FILE | sed "s/geojson.json/geopackage.tar.xz/"` # replaces old filename
    tar -cJf $FILENEWCOMPRESSED "$FILENEW"


done
exit