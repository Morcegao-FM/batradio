FROM node
WORKDIR /app
COPY server.js .
COPY package.json .
RUN npm install 
EXPOSE 9320
ENTRYPOINT [ "node", "server.js" ]