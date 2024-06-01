ok so we have alice and bob, and they want to communicate with each other.

Each one of them generates a pair of asymeetric keys, a public and a private one.

When alice wants to send a message to bob, she generates symmetric enryption key , such as AES, to encrypt message content.This symmetric key is randomly generated and serves as a shared secret between Alice and Bob.

Alice encrypts the symmetric key with Bob's public key. By encrypting the symmetric key with Bob's public key, Alice ensures that only Bob, possessing the corresponding private key, can decrypt the symmetric key and access the encrypted message

Once Bob receives the encrypted message and decrypts the symmetric key with his private key, he can then decrypt the message content with the symmetric key.

Conversation:
- messages all
- last timetamp
- sender
- receiver
- AES key encrypted for each person(shoud be dynamic, one for each person)
- iv for decryption, initialization vector


database - message-data-centre:
- message(one message encryoted witn AES key)
- timestamp
- status (read, unread, deleted, etc.)

database - user-data-centre: MORE
- public key
- private key

main service should be doing this
it can be also done in auth-service but it will be more complex