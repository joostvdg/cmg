FROM scratch
EXPOSE 8080
ENTRYPOINT ["/cmg"]
CMD ["serve"]
COPY ./bin/ /
