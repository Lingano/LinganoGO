# LinganoGO

LinganoGO is the backend service for the Lingano platform, written in Go. It provides a robust GraphQL API to power the Lingano ecosystem.

## Tech stack
- **GraphQL**: implemented with Go and gqlgen
- **PostgreSQL**: database
- **[Ent](https://entgo.io/docs/)**: ORM used to access the database

## Features

- **High Performance**: Built with Go for speed and efficiency
- **GraphQL API**: Provides flexible querying capabilities with gqlgen
- **PostgreSQL Database**: Uses Ent ORM for efficient data access
- **Scalable Architecture**: Designed to handle high traffic and large datasets
- **Secure**: Implements best practices for authentication and data protection

## Getting Started

Follow these instructions to set up and run the LinganoGO backend on your local machine or server.

### Prerequisites

-   [Go](https://golang.org/) (version 1.18 or higher)
-   [Docker](https://www.docker.com/) (optional, for containerized deployment)
-   A running instance of PostgreSQL (configured in `config/db.go`)

### Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/LinganoGO.git
    cd LinganoGO
    ```
2. Install the dependencies:
    ```bash
    go mod download
    ```
3. Configure the environment variables:
   Copy the `.env.example` file to `.env` and update the values as needed.
4. Run the application:
    ```bash
    go run cmd/server/main.go
    ```
5. Access the GraphQL Playground:
   Open `http://localhost:8080/` in your web browser to interact with the GraphQL API.

## Contributing

We welcome contributions to LinganoGO! Please follow these steps to contribute:

1. Fork the repository on GitHub.
2. Create a new branch for your feature or bug fix:
    ```bash
    git checkout -b my-feature-branch
    ```
3. Make your changes and commit them:
    ```bash
    git commit -m "Add my feature"
    ```
4. Push your changes to your forked repository:
    ```bash
    git push origin my-feature-branch
    ```
5. Submit a pull request on GitHub.

Please ensure that your code adheres to the existing coding standards and includes appropriate tests.

## License

LinganoGO is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.