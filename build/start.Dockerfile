FROM docker

WORKDIR /ci.io

COPY . .

RUN echo 'PS1="\[$(tput setaf 3)$(tput bold)[\]appname@\\h$:\\w]#\[$(tput sgr0) \]"' >> /root/.bashrc

EXPOSE 80

CMD . ./bin/build && /bin/sh
