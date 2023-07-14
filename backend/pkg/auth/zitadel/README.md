## RBAC VMClarity Proof of Concept

In order to enable RBAC in the whole project, it is required to have three things in place:
1. Identity Provier (IdP) instance, e.g. ZITADEL https://zitadel.com/, to be configured. This service can be either configured manually or automatically. 
For now, we shall keep it simple by having a preconfigured instance.
2. Protected API Middleware
3. Clearly defined clients. We want to be able to authenticate the following items: CLI components (CICD, VMClarity, CURL), UI Componenets (frontend, backend), and ScanOrchestrator. 
These clients need to be able to obtain access tokens based on their type. For example, ScanOrchestrator should be able to request a token which enables different access than e.g. CLI. 
This is the core of the problem.

   
### GOAL #1
Manually create a local Zitadel instance. Create an VMClarity API path that uses Zitadel protected middleware. 
Access route with CURL CLI.

1. Created a local Zitadel instance (manual)
```bash
docker-compose up -d
# console:  http://localhost:8080/ui/console
# username: zitadel-admin@zitadel.localhost
# password: Password1!
# client: 222868995949264899@project-demo
```
2. Wrapped API with zitadel middleware
3. `curl -i http://localhost:8888/api/assets` returns `auth header missing`
4. Configured Zitadel and created a service account JTW token in Zitadel project (manual)
5. `curl -i -H "Authorization: Bearer ${token}" http://localhost:8888/api/assets` returns `{"items":[]}`

### GOAL #2
Programmatically bootstrap a Zitadel instance. 

This is possible using native Zitadel Golang client connected with service account that has permissions 
to create organizations/projects/roles/users/groups.
In summary, this indicates that we only require an _API key with necessary permissions_ to fully configure and use Zitadel.
Moreover, this provides **dynamic management and configuration capabilities** that we can use to extend RBAC for our needs.
The only problem is:
- We need to have a deployed/existing Zitadel instance
- We always need a Service Account to interact with Zitadel instance

### GOAL #3
Finish bootstrapping Zitadel instance. Figure out how to distribute initially generated auth tokens.
