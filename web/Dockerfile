FROM node:23-alpine

WORKDIR /app

COPY . .

RUN npm install

ENV VITE_RATE_SHIELD_BACKEND_BASE_URL=${VITE_RATE_SHIELD_BACKEND_BASE_URL}
ENV PORT=${PORT}

EXPOSE $PORT

CMD [ "npm", "run", "dev" ]