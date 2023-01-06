FROM golang:1.18.2-alpine
RUN mkdir -p /home/app
# copies from host to container
COPY . /home/app 
#run the entrypoint command
CMD ["/home/app/main"]