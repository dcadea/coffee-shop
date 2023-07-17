# Coffee Shop Service

### Prerequisites
You should have one of docker distributions installed and running (eg. Docker, Rancher).

### Try it out
Run the following script included in this project
```shell
./run.sh
```
Execute request
```shell
curl -L -X POST 'http://localhost:8080/coffee' \
-H 'User-Id: 2' \
-H 'Membership-Type: americano_maniac' \
-H 'Content-Type: application/json' \
-d '{
    "coffee_type": "espresso"
}'
```
### Micro-spec :)
| Property        | Accepted values                             |
|-----------------|---------------------------------------------|
| User-Id         | any unassigned integer `uint`               |
| Membership-Type | `basic`, `coffee_lover`, `americano_maniac` |
| coffee_type     | `cappuccino`, `espresso`, `americano`       |
