# Snippetbox

Snippetbox is a web application for creating, managing, and sharing text-based snippets. It allows users to sign up, log in, and post snippets with a specified expiration time.

## Features

*   **User Authentication:** Users can sign up for a new account, log in, and log out.
*   **Snippet Management:** Authenticated users can create new snippets with a title, content, and an expiration period (1, 7, or 365 days).
*   **View Snippets:** Users can view a list of the latest snippets on the homepage and can view individual snippets.
*   **Secure:** The application uses HTTPS to encrypt all traffic and has secure session management.
*   **Form Validation:** All forms have validation to ensure data integrity.
*   **Flash Messages:** The application provides feedback to the user through flash messages (e.g., "Snippet successfully created!").

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   [Go](https://golang.org/) (version 1.18 or newer)
*   [MySQL](https://www.mysql.com/) (any V)

### Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/snippetbox.git
    sudo apt install mysql-server
    sudo apt install golang-go
    
    cd pastebinclone51
    ```

2.  **Database Setup:** 3 DB total
    *   Connect to your MySQL instance and create a new database for the project.
        ```sql
        CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
        ```
    *   Create the necessary tables by running the SQL scripts in the `internal/models/testdata` directory. You will need to create the `snippets` and `users` tables.

3.  **Configuration:**
    *   Create a `config.json` file in the `cmd/web/` directory with your MySQL Data Source Name (DSN).
        ```json
        {
            "dsn": "your-username:your-password@/snippetbox?parseTime=true"
        }
        ```
    *   Replace `your-username` and `your-password` with your MySQL credentials.

4.  **TLS Certificates:**
    *   The application requires TLS certificates to run over HTTPS. You can generate self-signed certificates for local development.
    *   Create a `tls` directory in the root of the project.
    *   Generate the certificate and key files using a tool like OpenSSL:
        ```bash
        go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
        mv cert.pem key.pem ./tls/
        ```

### Running the Application

Once you have completed the setup steps, you can run the application with the following command:

```bash
go run ./cmd/web
```

The application will be available at `https://localhost:4000`.

## Technology Stack

*   **Backend:** [Go](https://golang.org/)
*   **Database:** [MySQL](https://www.mysql.com/)
*   **Routing:** [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
*   **Session Management:** [alexedwards/scs](https://github.com/alexedwards/scs)
*   **Templating:** Go's built-in `html/template` package

## Future Work

*   Automate infrastructure setup and deployment using [Terraform](https://www.terraform.io/). This includes provisioning the MySQL server, managing users, and handling TLS certificates.
