# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server

# uses alpine base image
FROM  alpine

# installs ca-certificates package using run command
RUN apk add --no-cache ca-certificates

# website files into the container
COPY gateway gateway

# purposed for https
EXPOSE 443

# tell docker to exceute GO application
ENTRYPOINT [ "/gateway" ]

