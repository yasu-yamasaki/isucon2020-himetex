FROM nginx:1.19.2

RUN apt-get update \
    && apt-get -y upgrade \
    && apt-get install wget unzip --no-install-recommends -qy \
    && wget https://github.com/tkuchiki/alp/releases/download/v0.3.1/alp_linux_amd64.zip \
    && unzip alp_linux_amd64.zip \
    && install ./alp /usr/local/bin
