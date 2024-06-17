# Forum Project

## Description

This project is a web forum built using SQLite for database management and Docker for containerization. Users can register, log in, create posts and comments, and like or dislike posts and comments.

## Requirements

- Go 1.16+
- Docker


## Installation

1. Clone the repository:
    ```sh
    git clone https://zero.academie.one/git/dkalykov/forum.git
    cd forum
    ```

2. Build and run the Docker containers:
    ```
    make run
    ```

    ## Usage

### Registration

To register as a new user, provide your email, username, and password. If the email or username is already taken, an error response will be returned.

### Login

To log in, provide your registered email and password. If the credentials are correct, you will be logged in and a session will be created using cookies. Each user can have only one active session at a time.

### Creating Posts and Comments

Only registered users can create posts and comments. When creating a post, you can associate one or more categories with it. All posts and comments are visible to both registered and non-registered users.

### Likes and Dislikes

Only registered users can like or dislike posts and comments. The number of likes and dislikes is visible to all users.

### Filtering

Users can filter posts by categories, created posts, and liked posts. The latter two filters are only available to registered users and refer to the logged-in user.


The project handles various types of errors, including:

- Website errors
- HTTP status codes
- Technical errors


### Author

*** Dana K ***