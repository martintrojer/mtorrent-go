FROM steeve/libtorrent-go:linux-x64

RUN apt-get -y update
RUN apt-get -y install git tmux curl tree htop wget nano

RUN adduser --disabled-password --gecos '' martin
WORKDIR /home/martin
USER martin
ENV GOPATH /home/martin

RUN git clone https://github.com/steeve/libtorrent-go.git /home/martin/src/github.com/steeve/libtorrent-go
RUN cd src/github.com/steeve/libtorrent-go && make

RUN go get code.google.com/p/gcfg
RUN go get github.com/dustin/go-humanize

ENV PATH ${PATH}:/home/martin/bin
WORKDIR /home/martin/src/github.com/martintrojer/mtorrent-go