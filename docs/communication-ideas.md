# Idea 1: conversations

When a user starts messaging another user, a conversion is created. the list of conversations with keys will be kept in the auth-service.
When a user logs in, the auth-service will fetch the list of conversations for that user and push them to the active-sessions.

Keys to a conversation are symmetric for users taking part in the conversation. The key is used to decrypt the conversation from the database.

### Conversation
Conversation is a list of messages between two users.
The conversation keys are kept in a separate table in the auth-service.


### Message
Message is a text message sent by a user to another user. Every message has a unique id, timestamp and data - the actual message.

### Keys to conversation
A key to a conversation is the same for all users in the conversation. After login, all keys of a given user are stored within the active-sessions database in the form of a dictionary pairs "conversation_key": "conversation_id".
```json
session = {
    "user_id": "user_id",
    "username": "username",
    "location": "location",
    "timestamp": "timestamp",
    "user_key": "user_key",
    "conversations": {
        "conversation_key": "conversation_id",
        "conversation_key": "conversation_id",
        "conversation_key": "conversation_id",
        ...
    }
}
```


This key is used to decrypt the conversation, accessing all previous messages.

### Needed implementation:

#### Auth-service
- Managing table of conversations and additional keys
- Saving keys to conversations on users request


### Pros and cons

#### Pros
- Easy in implementation (but requires chanes in auth-service)
- Probably easey for multiple users in a conversation

#### Cons
- Requires changes in auth-service
- Not as secure
- Not as fast (lots of keys, encryption/decryption)