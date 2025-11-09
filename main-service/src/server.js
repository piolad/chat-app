// server.js
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const winston = require('winston');
const { makeHandlers } = require('./handlers/browser-facade-handlers');
const { makeServers } = require('./lib/servers');

const protoPaths = [
  '../protos/BrowserFacade.proto',
  '../protos/auth.proto',
  '../protos/service.proto',
];

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.json(),
  transports: [new winston.transports.Console()],
});

process.on('uncaughtException', (err) => logger.error(`uncaught: ${err.stack || err}`));
process.on('unhandledRejection', (err) => logger.error(`unhandledRejection: ${err}`));

const packageDefinition = protoLoader.loadSync(protoPaths, {
  keepCase: false, longs: String, enums: String, defaults: true, oneofs: true,
});
const loadedProtos = grpc.loadPackageDefinition(packageDefinition);

const BrowserFacadeService = loadedProtos.BrowserFacade.BrowserFacade.service;

const services = makeServers(loadedProtos, logger, {
  authAddr: process.env.AUTH_ADDR || 'auth-service:50051',
  msgAddr : process.env.MSG_ADDR  || 'message-data-centre:50051',
  insecure: process.env.GRPC_INSECURE !== 'false',
});

const handlers = makeHandlers({ logger, services });

const server = new grpc.Server();
server.addService(BrowserFacadeService, handlers);

const BIND_ADDR = process.env.BIND_ADDR || '0.0.0.0:50050';
server.bindAsync(BIND_ADDR, grpc.ServerCredentials.createInsecure(), (err) => {
  if (err) { logger.error(err); process.exit(1); }
  logger.info(`Server running at ${BIND_ADDR}`);
});

// Graceful shutdown
process.on('SIGTERM', () => server.tryShutdown(() => process.exit(0)));
process.on('SIGINT',  () => server.tryShutdown(() => process.exit(0)));
