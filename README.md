# Cards service


## Create deck
#### (shuffle=true,false; cards=AS,KH,2C,QH)

`curl --request POST \
  --url 'http://localhost:8080/decks/create?shuffle=false&cards=AS%2CKH%2C2C%2CQH'`
  
## Open Deck

  `curl --request GET \
  --url http://localhost:8080/decks/{deck-id}`
   
## Draw cards from deck
#### (Require deck version in payload to keep this handler idempotent. starts from 0.)

 ` curl --request PATCH \
  --url http://localhost:8080/decks/{deck-id} \
  --header 'Content-Type: application/json' \
  --data '{
	"draw": 2,
	"version": 0
}'`
