# Builds slate md source into static html files
# For example source files see the "source" directory in the slate github repo
# docker run -v /absolute/path/to/source:/source -v /absolute/path/to/build:/build -v /absolute/path/to/manifest:/manifest pwittrock/brodocs 

FROM node:7.2
MAINTAINER Phillip Wittrock <pwittroc@google.com>

RUN apt-get update && apt-get install -y npm git && apt-get clean && rm -rf /var/lib/apt/lists/*

RUN echo "v1.4"
RUN git clone --depth=1 https://github.com/Birdrock/brodocs.git
WORKDIR brodocs
RUN npm install
RUN node brodoc.js

COPY runbrodocs.sh .

CMD ["./runbrodocs.sh"]
