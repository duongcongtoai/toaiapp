# DEPLOYMENT

```console
bash build.sh 
```

Done!

## Test

### Create user with command

```console
docker exec -it ($app_container_id) app --configFile=config/config.yaml auth add $(username) $(password)
```

### Get access token

```console
curl --location --request POST 'localhost:8082/api/v1/user/login' \
--form 'username=($username)' \
--form 'password=($password)'
```

### Do something with access token

```console
curl --location --request GET 'localhost:8082/api/v1/user/' \
--header 'Authorization: Bearer ($received_token)' 
```
