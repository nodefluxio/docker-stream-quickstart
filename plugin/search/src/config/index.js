const envVar = process.env;
const { NODE_ENV, SERVER_API } = envVar;

const config = {
  NODE_ENV,
  PUBLIC_URL: process.env.PUBLIC_URL,
  SERVER_API // server IP where this repo hosted
};

module.exports = config;
