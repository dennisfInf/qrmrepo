### Build Image
docker build -t relay .

### Prepare database
docker run -d \
--name relay \
-e SERVERHOST=replaceme \
-e SERVERPORT=replaceme \
-e DBHOST=replaceme \
-e DBPORT=replaceme \
-e DBNAME=replaceme \
-e DBUSER=replaceme \
-e DBPASSWORD=replaceme \
relay setup

### Run server
docker run -d \
--name relay \
-e SERVERHOST=replaceme \
-e SERVERPORT=replaceme \
-e DBHOST=replaceme \
-e DBPORT=replaceme \
-e DBNAME=replaceme \
-e DBUSER=replaceme \
-e DBPASSWORD=replaceme \
relay serve