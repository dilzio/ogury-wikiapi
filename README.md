# gt_mtc_takehome

docker build -t mtc-api .
docker run -p 8080:8080 -it --rm --name mtc-api mtc-api
docker run mtc-api go test ./...
http://localhost:8080/mostviewedday/Albert_Einstein/2015/07
http://localhost:8080/mostviewed/20210101/20210401
