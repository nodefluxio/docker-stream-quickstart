# Plugin Installation

## With Docker Command

1. Install plugin API - Service to hit Dukcapil & Korlantas
```
sudo docker run -it -d -p 6014:6014 \
--restart unless-stopped \
--name search-plugin-api nodefluxio/vanilla-dashboard:1.2.3 ./search-plugin-api \
--seagate-base-url [IP for search face ] \
--polri-base-url [IP for search nik & plate] \
--polri-username [polri-username] \
--polri-password [polri-password] \
--fremisn-url [IP for fremis API]
```

2 Install plugin dashboard
```
sudo docker run -it -d -p 8081:80 \
--restart unless-stopped \
--env NODE_ENV=production \
--env ENV_DEST=/ \
--env SERVER_API=http://localhost:6014 \
--name search-plugin registry.gitlab.com/nodefluxio/vanilla-dashboard/search:1.2.3
```

3. Install Vanilla Dashboard with plugin, just add `--env` for **PLUGIN_HOST** and **PLUGIN_NAME** before `--name`
```
sudo docker run -it -d --net host \
--restart unless-stopped \
--env PLUGIN_HOST=http://localhost:8081 \
--env PLUGIN_NAME=Search \
--name vanilla-dashboard nodefluxio/vanilla-dashboard:1.2.3 ./vanend \
--visionaire-host localhost:4004 \
--website-host localhost
```

more info for command at **services: vanilla dashboard**: https://docs.nodeflux.io/visionaire-docker-stream/installation-guide/advanced-installation
---

## With Docker Compose

1. Download plugin image 
2. Edit or create docker-compose.yml with config as follows:
```
version: "3"
services:
  vanilla-dasboard:
    container_name: vanilla-dashboard
    image: registry.gitlab.com/nodefluxio/vanilla-dashboard:[tag version]
    command: ./vanend --visionaire-host [IP of visionaire docker stream host] --website-host [IP of dashboard website host]
    network_mode: "host"
    environment:
    - PLUGIN_HOST= [IP of plugin host]:[plugin host port]
    - PLUGIN_NAME= Search
  search-plugin:
    container_name: search-plugin
    image: registry.gitlab.com/nodefluxio/vanilla-dashboard/search:[tag version]
    ports:
    - [plugin host port]:80
    environment:
    - NODE_ENV=production
    - ENV_DEST=/
    - SERVER_API=[IP of server host]:6014
  search-plugin-api:
    container_name: search-plugin-api
    image: registry.gitlab.com/nodefluxio/vanilla-dashboard:[tag version]
    command: ./searchingpolri --seagate-base-url [IP for search face ] --polri-base-url [IP for search nik & plate] --fremisn-url [IP for fremis API]
    ports:
    - 6014:6014
```
3. save docker-compose.yml, and run command below to apply the configuration and run docker 
```
docker-compose up -d
```
---

if restart or update needed after running docker, you can use command below to stop docker
```
docker-compose down
```
---

### Docker compose file example
```
version: "3"
services:
  vanilla-dasboard:
    container_name: vanilla-dashboard
    image: registry.gitlab.com/nodefluxio/vanilla-dashboard:1.2.3
    command: ./vanend --visionaire-host 192.168.101.248:4004 --website-host 192.168.101.248
    network_mode: "host"
    environment:
    - PLUGIN_HOST= http://192.168.101.248:8081
    - PLUGIN_NAME= Search
  search-plugin:
    container_name: search-plugin
    image: registry.gitlab.com/nodefluxio/vanilla-dashboard/search:1.2.3
    ports:
    - 8081:80
    environment:
    - NODE_ENV=production
    - ENV_DEST=/
    - SERVER_API=http://192.168.101.248:6014
  search-plugin-api:
    container_name: search-plugin-api
    image: registry.gitlab.com/nodefluxio/vanilla-dashboard:1.2.3
    command: ./searchingpolri --seagate-base-url http://192.168.101.248:3003 --polri-base-url http://192.168.101.248:3003 --fremisn-url http://192.168.101.248:4005/v1/face
    ports:
    - 6014:6014
```
