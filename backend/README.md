### Build Image
`
docker build -t backend .
`

### Run backend
```
docker run -d \
--name backend \
-e SERVERPORT=replaceme \
-e INFURAADRESS=replaceme \
backend serve
```
