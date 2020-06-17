# DEPLOYMENT

```console
bash build.sh 
```

Done!

## Test

### Create user with command

```console
docker exec -it toaiapp_app_1 app --configFile=config/config.yaml auth add $(username) $(password)
```

### Authorize
Goto: http://localhost:8082/oauth/login
=> Login with username and password for session


### Use Oauth2 

Goto: http://localhost:8084
=>Token is received

### Todo:
Add endpoint to test token