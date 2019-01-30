FROM scratch
EXPOSE 8080
ENTRYPOINT ["/cmg"]
COPY ./bin/ /