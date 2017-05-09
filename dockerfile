FROM scratch

COPY /app /app/

CMD [“/app”]
EXPOSE 8080
