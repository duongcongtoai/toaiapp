# DEPLOYMENT

bash build.sh
Done!

## Test

docker exec -it ($app_container_id) app --configFile=config/config.yaml auth add $(username) $(password)

Test:
curl --location --request POST 'localhost:8082/api/v1/user/login' \
--form 'username=($username)' \
--form 'password=($password)'

curl --location --request GET 'localhost:8082/api/v1/user/' \
--header 'Authorization: Bearer ($received_token)' \
--form 'username=toai' \
--form 'password=toai'