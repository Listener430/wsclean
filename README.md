## Test harness for Boxes app websocket features
The code contains 2 handlers WSHandle and WSNotifications, also the code for OnlineHub and Notification clients. 

WsHandler:  It creates a NotificationClient and adds it to the NotificationPools map (it does not start any goroutines)
WsNotificationHandler:  It starts two goroutines for reading and writing to the WebSocket connection. It initialized OnlineClient and registers it to the OnlineHub.  

To send a notification the correct order is to hit the first end point (e.g. from postman or the FE)
ws://localhost:8081/ws/56fa9e31-8035-4c0c-b07d-c90d95a91d81 (request body is empty)

Then hit the second endpoint with the following 
ws://localhost:8081/ws/online/56fa9e31-8035-4c0c-b07d-c90d95a91d81

The second endpoint expects JSON in the body of the following format: 
{
  "userID": "9403be31-d7c7-423d-a529-1cb3f7f0cbcf",
  "actionUserID": "d1e2054f-1041-441b-b9fc-5a25af7175cf",
  "description": "You have a new follow request from 117",
  "status": true,
  "type": "FLW",
  "sourceId": "d1e2054f-1041-441b-b9fc-5a25af7175cf",
  "attributes": {
    "attribute1": "value1",
    "attribute2": "value2"
  }
}

There's also 3rd endpoint called follow, it tests notification storage (a local instance of postgress is required) 
http://localhost:8081/users/56fa9e31-8035-4c0c-b07d-c90d95a91d81/follows
The endpoint expects JSON like the one above in the body to send notifications and store it to the DB. 

Table creation in postgres
-- Table: public.notification

-- DROP TABLE IF EXISTS public.notification;

CREATE TABLE IF NOT EXISTS public.notification
(
    id uuid,
    _created timestamp without time zone NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
    _modified timestamp without time zone NOT NULL DEFAULT 'CURRENT_TIMESTAMP',
    _deleted boolean DEFAULT 'false',
    user_id uuid NOT NULL,
    description text COLLATE pg_catalog."default" NOT NULL,
    action_user_id uuid NOT NULL,
    status boolean DEFAULT 'false',
    type character varying(15) COLLATE pg_catalog."default" NOT NULL,
    source_id uuid NOT NULL,
    attributes jsonb NOT NULL DEFAULT '{}'::jsonb
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.notification
    OWNER to postgres;

## Start
go run main.go

By default the server starts at localhost, port 8081



