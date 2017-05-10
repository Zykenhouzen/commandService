FROM scratch
EXPOSE 8080
COPY ./command-service /command-service
CMD ["./command-service"]
