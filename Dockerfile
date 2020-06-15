FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY . /toaiapp
WORKDIR /toaiapp 

RUN go get ./app
# RUN go build
RUN go get github.com/cespare/reflex

# CMD if [${APP_ENV} = production]; \
#   then \
#   app; \
#   else \
#   reflex -c reflex.conf; \
#   fi
# EXPOSE 8080