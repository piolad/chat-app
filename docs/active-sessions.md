## gRPC Service Methods:

### CreateSession
- **Request:** `user_data`
  - *Description:* Request to create a session of a user with the given data.
- **Response:** `session_data`
  - *Description:* Contains the session ID associated with the created session.

### UpdateSession
- **Request:** `user_data`
  - *Description:* Request to update a user session with the provided data.
- **Response:** `session_data`
  - *Description:* Contains the session ID or other relevant data after the update operation.

### DeleteSession
- **Request:** `session_data`
  - *Description:* Request to delete a user session based on the provided session ID.
- **Response:** `user_data`
  - *Description:* Contains the data of the removed user associated with the deleted session.

### GetData
- **Request:** `session_data`
  - *Description:* Request to retrieve the data of a user with the given session ID.
- **Response:** `userdata`
  - *Description:* Contains the user data associated with the provided session ID.

## gRPC Message Definitions:

### user_data:
- **Fields:**
  - `username`: string
  - `email`: string
  - `password`: string
  - `name`: string
  - `surename`: string
  - `location`: string
  - `birth_date`: string

### session_data:
- **Fields:**
  - `session_id`: string



