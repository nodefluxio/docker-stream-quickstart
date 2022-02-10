
# build stage - frontend
FROM node:12.13.0-alpine as build-frontend
ENV NODE_ENV=production
RUN \ 
  apk update && \
  apk add git python make gcc g++

ADD ./client /opt/app
WORKDIR /opt/app

RUN \
  yarn install && \
  npm run build

# build stage - backend
FROM golang:1.14-alpine as build-backend
ENV PATH="${PATH}:/usr/include/"

ARG ssh_prv_key
ARG ssh_pub_key
RUN echo $GOPATH
RUN echo $PATH

RUN \ 
  apk update && \
  apk add --no-cache git gcc build-base openssh-client curl

RUN mkdir -p /root/.ssh && \
  chmod 0700 /root/.ssh
RUN echo "$ssh_prv_key" > /root/.ssh/id_rsa && \
  echo "$ssh_pub_key" > /root/.ssh/id_rsa.pub && \
  chmod 600 /root/.ssh/id_rsa && \
  chmod 600 /root/.ssh/id_rsa.pub && \
  ssh-keyscan gitlab.com > /root/.ssh/known_hosts
RUN cat /root/.ssh/id_rsa.pub && cat /root/.ssh/id_rsa
RUN git config --global --add url."git@gitlab.com:".insteadOf "https://gitlab.com/"

# Compiling third party proto
ADD . /go/src/gitlab.com/nodefluxio/vanilla-dashboard
WORKDIR /go/src/gitlab.com/nodefluxio/vanilla-dashboard
ENV GENFILES=/go/src

# Compiling app
RUN go get -v ./cmd/vanend
# RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/swaggerui -i ./cmd/swagger-ui
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/vanend -i ./cmd/vanend
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/cescoordinator -i ./cmd/cescoordinator
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/cesagent -i ./cmd/cesagent
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/searchingpolri -i ./cmd/searchingpolri

# final stage
FROM frolvlad/alpine-glibc

RUN apk add tzdata nodejs-current yarn bash curl jq
WORKDIR /go/bin
RUN curl -fsSL -o /go/bin/dbmate https://github.com/amacneil/dbmate/releases/download/v1.12.0/dbmate-linux-amd64
RUN chmod +x /go/bin/dbmate
RUN cp /go/bin/dbmate /bin/
RUN yarn global add dotenv-to-json @ethical-jobs/dynamic-env
COPY --from=build-frontend /opt/app/build  /go/bin
COPY --from=build-frontend /opt/app/.env.build /go/bin/.env.example
COPY --from=build-backend /go/bin/vanend /go/bin/vanend
COPY --from=build-backend /go/bin/cescoordinator /go/bin/cescoordinator
COPY --from=build-backend /go/bin/cesagent /go/bin/cesagent
COPY --from=build-backend /go/bin/searchingpolri /go/bin/searchingpolri
COPY --from=build-backend /go/src/gitlab.com/nodefluxio/vanilla-dashboard/template /go/bin/template/
COPY --from=build-backend /go/src/gitlab.com/nodefluxio/vanilla-dashboard/script /go/bin/script/
COPY --from=build-backend /go/src/gitlab.com/nodefluxio/vanilla-dashboard/internal/infrastructure/db/psql/migrations /go/bin/internal/infrastructure/db/psql/migrations/
COPY --from=build-backend /go/src/gitlab.com/nodefluxio/vanilla-dashboard/internal/infrastructure/db/psql/migrations_ces_agent /go/bin/internal/infrastructure/db/psql/migrations_ces_agent/
COPY --from=build-backend /go/src/gitlab.com/nodefluxio/vanilla-dashboard/internal/infrastructure/db/psql/migrations_ces_coordinator /go/bin/internal/infrastructure/db/psql/migrations_ces_coordinator/


EXPOSE 80


