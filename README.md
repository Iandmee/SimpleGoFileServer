# Simple Go File Server
## Get started
Clone the repository, and make a docker image
```shell
docker build -t go-file-server .
```
Now you should have a new docker image of this file server `go-file-server`
```shell
> docker images                             
REPOSITORY              TAG       IMAGE ID       CREATED          SIZE
go-file-server          latest    a384f3f922be   11 minutes ago   1.47GB
```

To run the docker image run the following command (specify `<port>` by your port number, For example: `8080`)
```shell
docker run -p <port>:8080 go-file-server
```
The server starts on all interfaces of your machine. To access it locally use ```http://127.0.0.1:<port>```
## Examples of the server usage
For example, I will use `8080` - as a forwarded port
### Delete file
Delete `/tmp/directory/file` file
```shell
curl -X DELETE http://127.0.0.1:8080/tmp/directory/file
```
If successful, the server returns
```shell
File "file" successfully deleted.
```

### Create/Replace file
```shell
curl -X POST -F "data=@/tmp/file" http://127.0.0.1:8080/tmp/directory/file
```
If successful, the server returns
```shell
File "file" successfully uploaded.
```

### Download file
```shell
wget http://127.0.0.1:8080/tmp/directory/file
```
If successful, the server returns the file's content
