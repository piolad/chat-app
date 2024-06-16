
    Situation: 
        User makes a search for the other user: 
        main service communicates with auth service to check if user with given username/email exist in database 

        If the user exist: 
        Option 1: 
            - User is not a friend: Show caption add to friends 
        Option 2: 
        User is a friend:
        
        When user clicks on the friend, main-service sends request to the auth-service, in order to retrive id and username of the selected friend 
        auth-service returns  and username  of the friend, along with the id and username of the user that requested it. Both id's along with the usernames of both users are then sent to message data center. 

        After receiving ids and usernames message data center must implement the funtion that sends a SELECT query that might look something like that: 

        SELECT * 
        FROM Messages 
        WHERE id_1 = username1 AND id_2 = username2
        ORDER BY date DESC 
        LIMIT 10;

        This SELECT will return 10 last messages between the users along with other information held in the table Messages that might look like a table below: 

        CREATE TABLE Messages (
            message_id PRIMARY KEY,                 <- 
            sender_id INT NOT NULL,                 <- could also be string, we can use username as the id 
            receiver_id INT NOT NULL,               <- could also be string, we can use username as the id
            message TEXT NOT NULL,                  
            date   (there is some data type to handle dates)                          
            if_read BOOLEAN DEFAULT FALSE,
        );

        "message" is encrypted, so message-data center must also have a function to retrive the original message with the help of id's.

        After original message, or 10 messages as proposed in the select query, is retirved, it goes back to the main service that sends it back to the browser facade that handles the display of messages. Table "Messages" has a column "sender_id", "receiver_id" and "date" that can help to display the messages properly by the browser facade

        This design creates a problem of concurency of transactions, when many users send messages at the same time, but i dont know how to solve it right now, 
        maybe by separating conversations between pair of users into the other table called "Chats"

        
3 tabele:
- conversationID
- user
- AES encrypted