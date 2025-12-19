FROM ubuntu
LABEL multi.label1="value1" multi.label2="value2" other="value3"
WORKDIR /yk
ADD ../../config/.env /yk/.env
ADD bin/amd64/myapp /yk/myapp
EXPOSE 8080
ENTRYPOINT ["/yk/myapp"]