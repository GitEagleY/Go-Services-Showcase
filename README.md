# Microservices Project

## Prerequisites

- Docker
- Docker Desktop (optional but recommended for monitoring service state)

If you are using Visual Studio Code, open the `Microservices.code-workspace` file.

## Building and Running

Follow these steps to set up and run the microservices project:

1. Switch to the "docker" folder.
2. Run the following command to build and run Docker images:
        make up_build
3. Start the frontend on port 8082 by running:
        make start
4. To stop the Docker images, use the following command:
        make down
5. To stop the frontend, run:
        make stop

## Notes

Please keep the following notes in mind while working on this project:

1. For proper functionality of the authentication service, you should execute the `users.sql` query in the PostgreSQL users database (details of connection in `docker-compose.yml`).

2. It is recommended to use the MailHog service for email-related tasks.

3. To check logs, it is recommended to use MongoDB Compass.

## Technologies Used

This project utilizes various technologies, including RPC, REST, RabbitMQ, and other popular technologies. Feel free to inspect the code for more details.

## License

This project is released under the GNU General Public License (GPL).

