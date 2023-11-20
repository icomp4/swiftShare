******SwiftShare is a WIP file-sharing application that combines a backend written in Go with a frontend featuring user authentication, login, signup, and a main page for file management. The backend utilizes standard Go libraries and provides endpoints to support user management and authentication through JSON Web Tokens (JWT).******

****Features****

    User Authentication: Users can sign up, log in, and automatically log in on subsequent visits using JWT sessions.
    File Management: Access the main page for file sharing and management.
    Session Management: JWT manages user sessions securely.
    Email Authentication: Users can request an email code to perform changes such as updating their password, or deleting their account. 

****Examples****

 ![Demo](https://i.gyazo.com/117c0d1f536ccbd0206669426bec281e.gif)


****Tech Stack****

- [Golang Repository](https://github.com/golang/go)
- [HTML5 Wikipedia](https://en.wikipedia.org/wiki/HTML5)
- [Golang JWT Library](https://github.com/golang-jwt/jwt)
- [godotenv on GitHub](https://github.com/joho/godotenv)
- [pq on GitHub](https://github.com/lib/pq)
- [PostgreSQL Official Site](https://www.postgresql.org/)

****API Endpoints****

    POST /api/v1/signup: Create a new user account.
    POST /api/v1/login: Log in and obtain a JWT token.
    POST /api/v1/logout: Log out and invalidate the session.
    DELETE /api/v1/delete: Delete user account.
    POST /api/v1/request: Request email verification.
    POST /api/v1/password/update: Update user password.
    GET /api/v1/user: Retrieve user details.

****Folder Structure****
    
    controllers/: Handles database queries for an added layer of abstraction in /handlers.
    database/: Database-related functionalities.
    handlers/: Backend handlers for different endpoints.
    handlers/middleware/: Middleware functions for Auth, Logging, and Email functionality.
    handlers/validators/: Various validation such as password requirements and extracting JWT tokens.
    static/: Contains static files for the frontend.
    utils/: Various utilities for smoother development.
    
