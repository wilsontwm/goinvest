# Start from base image
FROM golang:alpine

# Create an /go/src/
RUN mkdir /goinvest

# Copy everything in root directory to /goinvest directory
ADD . /goinvest

# Specify that subsequent commands to be executed in /goinvest directory
WORKDIR /goinvest

# Build the application
RUN go build -o main .

# Expose necessary port
EXPOSE 3000

# Run the created binary executable after wait for mysql container to be up
CMD ["./wait-for.sh" , "mysql:3306" , "--timeout=300" , "--" , "/goinvest/main"]