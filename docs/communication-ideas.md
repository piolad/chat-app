ok so we have alice and bob, and they want to communicate with each other.

Each one of them generates a pair of asymeetric keys, a public and a private one.

When alice wants to send a message to bob, she generates symmetric enryption key , such as AES, to encrypt message content.This symmetric key is randomly generated and serves as a shared secret between Alice and Bob.

Alice encrypts the symmetric key with Bob's public key. By encrypting the symmetric key with Bob's public key, Alice ensures that only Bob, possessing the corresponding private key, can decrypt the symmetric key and access the encrypted message

Once Bob receives the encrypted message and decrypts the symmetric key with his private key, he can then decrypt the message content with the symmetric key.