# browser-facade
Browser-facade microservice is an intermediary buffer service taking REST requests from user's browser and relaying the adequate gRPC requests further, to the main-service service. It serves as a presentation layer for the user, also sending HTML and CSS files to the user's browser.

## Incomming traffic (user&#8594;browser-facade) - REST

### login (POST)
Used to log user in. 
##### params:

```cpp
    string username
    string email
    string password
```
##### returns:
1.  `200 OK` (successful login)
    ```cpp
    string sessionId
    (some time format?) TTL 
    ```
1.  `401 Unauthorized` (failed login)
    ```cpp
    (some enum type) reason
    string longerReason
    ```


### register (POST)
Create new user.
```cpp
    string username
    string email
    string password
    string name
    string surname
    (some time format?) birthdate
```
##### returns:
1.  `200 OK` (successful register)
    ```cpp
    string sessionId
    (some time format?) TTL 
    ```
1.  `401 Unauthorized` (failed login)
    ```cpp
    (some enum type) reason
    string longerReason
    ```
 

### username_free (GET) 
Check if username is free (such username is not used in the application).
##### params:
##### &nbsp;
```cpp
        string username
```

##### returns:
1.  `200 OK` (username free)
    ```cpp
    string username
    ```
1.  `409 Conflict` (such user already exits)
    ```cpp
    string username
    ```
    

### email_free (GET) 
Check if email is free (such username is not used in the application).
##### params:
##### &nbsp;
```cpp
        string email
```

##### returns:
1.  `200 OK` (email free)
    ```cpp
    string email
    ```
1.  `409 Conflict` (such email already exits)
    ```cpp
    string username
    ```

### session_ttl (POST) 
Get current session's Time To Live
##### params:
##### &nbsp;
```cpp
        string sessionId
```

##### returns:
1.  `200 OK` ()
    ```cpp
    (some time format?) deadline
    ```

### send_message (POST) 
Get current session's Time To Live
##### params:
##### &nbsp;
```cpp
        string sessionId
        string revceiver
        string message
```

##### returns:
1.  `200 OK` ()
    ```cpp
    string messageId
    ```



> Following is the proposal of the mechanism for fetching messages and alerting about the new ones.
> Keep in mind that the messages should include the ones send from the current user.

### new_messages (POST) 
Check if there are any new messages since the given timestamp.
##### params:
##### &nbsp;
```cpp
        string sessionId
        (some time format?) since
```

##### returns:
1.  `200 OK` (no new messages)
    ```cpp
    string username
    ```
1.  `200 OK` (there are new messages)
    ```cpp
    string username
    int messageCount
    string[] senders
    ```

### get_messages__for_sender_on_span (POST) 
Get all messages withinin given 2 boundaries.
##### params (version A - time):
##### &nbsp;
```cpp
        string sessionId
        string sender
        (some time format?) since
        (some time format?) to
```

##### params (version B - messageId):
##### &nbsp;
```cpp
        string sessionId
        string sender
        string messageSince
        string messageTo
```

##### returns:
1.  `200 OK` ()
    ```cpp
    (some message type)[] messages
    ```


## Outgoing traffic (browser-facade&#8594;main-service) - gRPC
> _to be implemented_