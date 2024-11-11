## Contribution

### Project Setup

Before setting up the project, ensure you have the following installed:

- **[Node.js](https://nodejs.org/)** (includes npm)
- **[Golang](https://golang.org/dl/)**
- **[Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)**

### Setup Instructions

Follow these steps to set up and run the project locally:

1. **Install Node.js**

   Download and install Node.js from the [official website](https://nodejs.org/). This will also install npm, which is required for managing Node.js packages.

2. **Install Dependencies**

   Navigate to the project `/web` directory where the `package.json` file is located and run:

   ```bash
   npm install
   ```

3. **Start Docker Containers**

   In the project root directory (where your `docker-compose.yml` file is located), run:

   ```
   sudo docker-compose up
   ```
4. **Setup `.env` file**
    
   In folder `rate_shield` create a file named `.env` and add following content to it. You do not need to change these values if you are not working with something related to slack notification. Just copy and paste these 2 lines and it will work fine.


   ```
   SLACK_TOKEN=<Replace With Your Slack Token>
   SLACK_CHANNEL=<Replace With Your Slack Channel ID>
   ```
   
5. Run the Golang Application

   Open a new terminal window, navigate to the directory containing your `main.go` file, and run:

   ```
   go run main.go
   ```

6. Start the Frontend

   Open another terminal window, navigate to the frontend directory, and run:

   ```
   npm run dev
   ```

7. Access Application

   Open browser and go to `http://localhost:5173/`
