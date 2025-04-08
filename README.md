# go-saas-template

My template for quickly setting up a simple Go lang SaaS project.

## Features

### Authentication System

The application provides a complete authentication system built on PocketBase, featuring:

- **User Registration**: Email and password-based account creation with validation
- **Login**: Secure authentication with JWT tokens
- **Password Reset**: Self-service flow for users who forget their passwords
- **Session Management**: Stateless authentication using cookies
- **Protected Routes**: Middleware for securing application routes

## Implementation Details

### Authentication Flow

The authentication system follows a stateless approach using JWT tokens:

1. **Registration**: Users create accounts with email/password
2. **Login**: Users authenticate and receive a token stored in a cookie
3. **Session Validation**: Requests to protected routes verify the token
4. **Logout**: Token is cleared from cookies

### Password Reset Process

The password reset functionality follows these steps:

1. User requests a password reset by providing their email
2. System generates a secure token tied to the user account
3. User receives reset instructions with a link containing the token
4. User sets a new password, which is verified against the token
5. Upon successful reset, user is redirected to login

### User Interface

- **Login Page**: Email/password fields with forgot password link
- **Registration Page**: Email and password creation with validation
- **Forgot Password**: Email submission form
- **Reset Password**: New password entry form
- **Home Dashboard**: Authenticated user view with logout functionality

## Technical Architecture

The application uses:

- **PocketBase**: For user management and authentication
- **Gorilla Mux**: For HTTP routing
- **HTML Templates**: For rendering user interfaces
- **Cookie-based Auth**: For maintaining authenticated state

## Development Setup

### Prerequisites

- Go 1.19+
- PocketBase

### Running the Application

```bash
go run cmd/server/main.go
```

The server will start at http://localhost:8080 by default.

## Security Considerations

- Passwords are securely hashed
- Authentication tokens have appropriate expiration
- Password reset tokens are single-use and time-limited
- Error messages are designed to prevent information leakage
