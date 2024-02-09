# Test task for Medods

This service is a part of authentication service. Implement two operations: singIn(get pair of access and refresh token) and refresh(accepts refresh token and returns new pair of tokens). Tokens are connected by generated sessionID.  

Access token - jwt, refresh token - random byte sequence.  

The main technologies are:  
- `chi-router` for routing;
- `mongo-go` as mongoDB driver;
- `golang-jwt` - to wotk with jwt tokens;
- `google/uuid` - to create and validate GUID(UUID);
- `godotenv` to work with environment variables;

Implement graceful shutdown.

To run the app You can clone this project, move to directory with this project and run commands:
```
docker build -t headgardener1/medods_app .

docker-compose -f docker-compose.yml up mongo_db

docker-compose -f docker-compose.yml up app

```
