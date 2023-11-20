SwiftShare is a WIP file-sharing application that combines a backend written in Go with a frontend featuring user authentication, login, signup, and a main page for file management. The backend utilizes standard Go libraries and provides endpoints to support user management and authentication through JSON Web Tokens (JWT).
Features

    User Authentication: Users can sign up, log in, and automatically log in on subsequent visits using JWT sessions.
    File Management: Access the main page for file sharing and management.
    Session Management: JWT manages user sessions securely.
    Email Authentication: Users can request an email code to peform changes such as updating their password, or deleting their account. 

Tech Stack

    Backend: Implemented in Go, using standard libraries for handling HTTP requests and database management.
    Frontend: Developed using HTML templates and styled with Tailwind CSS.

API Endpoints

    POST /api/v1/signup: Create a new user account.
    POST /api/v1/login: Log in and obtain a JWT token.
    POST /api/v1/logout: Log out and invalidate the session.
    DELETE /api/v1/delete: Delete user account.
    POST /api/v1/request: Request email verification.
    POST /api/v1/password/update: Update user password.
    GET /api/v1/user: Retrieve user details.

Folder Structure
    
    controllers/: Handles database queries for an added layer of abstraction in /handlers.
    database/: Database-related functionalities.
    handlers/: Backend handlers for different endpoints.
    middleware/: Middleware functions for Auth, Logging, and Email functionality
    static/: Contains static files for the frontend.
    utils/: Various utilities for smoother development.
    

Contributing

Feel free to contribute by submitting bug reports, feature requests, or pull requests. Make sure to follow the contribution guidelines.
