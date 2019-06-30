#docker run -v $(pwd):/data geodata/gdal ogr2ogr -f GeoJSON -t_srs EPSG:4326 data/raw_data/polygon_cusecs_geojson_crs.json data/raw_data/polygon_cusecs_geojson.json

docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 -f -o data/raw_data/polygon_cusecs_geojson.mbtiles data/raw_data/polygon_cusecs_geojson_crs.json


# docker run -it -v $(pwd):/data fuzzytolerance/tippecanoe tippecanoe -z12 --drop-densest-as-needed --force-feature-limit -f -o data/raw_data/polygon_cusecs_geojson.mbtiles data/raw_data/polygon_cusecs_geojson_crs.json
