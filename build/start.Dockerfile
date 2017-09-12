FROM docker

WORKDIR /ci.io

COPY . .
#RUN echo 'PS1="\[$(tput setaf 3)$(tput bold)[\]appname@\\h$:\\w]#\[$(tput sgr0) \]"' >> /root/.bashrc
#RUN echo 'export PS1="\[\033[38;5;195m\]\[\033[48;5;244m\]\u\[$(tput sgr0)\]\[\033[38;5;15m\]\[\033[48;5;-1m\] [\H] > { \v:\l } > \w\n\\$\[$(tput sgr0)\]"' >> /root/.bashrc
#RUN . /root/.bashrc

EXPOSE 80

CMD . ./bin/build && /bin/sh
