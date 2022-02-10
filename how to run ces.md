## How To Run Centralize Enrollment System
### Run Coodinator service  
Install postgre

```bash
sudo docker run \
--name=postgres \
-e POSTGRES_PASSWORD=test -p 5432:5432 \
-v /var/lib/postgresql/data:/var/lib/postgresql/data \
-e PGDATA=/var/lib/postgresql/data/pgdata \
-d postgres:12-alpine
```

This service is central in CES for organizing data from agent and this service must be run in 1 server only.  
Example docker run* :
```bash
$ docker run -it -d --net host --restart unless-stopped --name ces-coordinator registry.gitlab.com/nodefluxio/vanilla-dashboard:1.2.0 ./cescoordinator
```

*) change all ip address at your server and image version with version you needed


### Run Agent Service and Vanilla Dashboard
1. Agent Service 
   Example docker run* :
   ```bash
   $ docker run -it -d --net host --restart unless-stopped --name ces-agent registry.gitlab.com/nodefluxio/vanilla-dashboard:1.2.0 ./cesagent --agent-name <agent-name> --db-user postgres --db-password test --coor-url http://192.168.101.176:6012 --enrollment-vanilla-url http://192.168.101.248
   ```

   *) change all ip address at your server and image version with version you needed  
   *) agent-name is unique must be different at all server

2. Vanilla Dashboard
   Example docker run* :
   ```bash
   $ docker run -it -d --net host --restart unless-stopped --name vanilla-dashboard --env USE_CES=true  registry.gitlab.com/nodefluxio/vanilla-dashboard:1.2.0 ./vanend --visionaire-host 192.168.101.248:4004 --website-host 192.168.101.248 --db-host 192.168.101.248 --db-user postgres --db-password test --use-ces true
   ```

   *) change all ip address at your server and image version with version you needed  
