const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const winston = require('winston');
const util = require('util')


const logger = winston.createLogger({
  level: 'info',
  format: winston.format.simple(),
  //defaultMeta: { service: 'main-service' },
  transports: [
    new winston.transports.Console(),
  ],
});


//global error logger:
process.on('uncaughtException', (err) => {
  logger.error(`Uncaught Exception: ${err.message}`);
  //process.exit(1);
});

const protoPahts = [
  '../protos/BrowserFacade.proto',
  '../protos/auth.proto',
]


const packageDefinition = protoLoader.loadSync(protoPahts , {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});

const loadedProtos = grpc.loadPackageDefinition(packageDefinition);

const BrowserFacadeService = loadedProtos.BrowserFacade.BrowserFacade.service;

const client = new loadedProtos.auth.Auth('auth-service:50051', grpc.credentials.createInsecure());

// function definitions
function Login_toAuthService(username, password) {
  return new Promise((resolve, reject) => {

    let request = {
      username: username,
      password: password
    };

    client.login(request, (error, response) => {
      if (error) {
        logger.error(`Error from authentication service: ${error.message}`);
        reject(error); // Reject the Promise with the error
      } else {
        logger.info('Login Response:', response);
        logger.info(`Login status: ${response.status}`);
        resolve(response); // Resolve the Promise with the response
      }
    });
  });
}

function Register_toAuthService(firstname, lastname, birthdate, username, email, password) {
  return new Promise((resolve, reject) => {
    let request = {
      firstname: firstname,
      lastname: lastname,
      birthdate: birthdate,
      username: username,
      email: email,
      password: password
    };

    client.register(request, (error, response) => {
      if (error) {
        logger.error(`Error from authentication service: ${error.message}`);
        reject(error); // Reject the Promise with the error
      } else {
        logger.info('Register Response:', response);
        logger.info(`Register status: ${response.status}`);
        resolve(response); // Resolve the Promise with the response
      }
    });
  });
}

async function Login_fromBrowserFacade(call, callback) {
  logger.info(`Login request received from ${util.inspect(call.request, {depth: null})}`);

  try {
    const resp = await Login_toAuthService(call.request.username, call.request.password);
    logger.info(`B ${util.inspect(resp, {depth: null})}`);
    callback(null, {
      success: resp.status == 'Success',
      username: call.request.username,
      token: resp.token,
      message: resp.status
    });
  } catch (error) {
    logger.error(`Error occurred during login: ${error.message}`);
    callback(error, null); // Sending error response to the client
  }
}

async function Register_fromBrowserFacade(call, callback) {
  logger.info(`Register request received from ${util.inspect(call.request, {depth: null})}`);

  try {
    const resp = await Register_toAuthService(call.request.firstname, call.request.lastname, call.request.birthdate, call.request.username, call.request.email, call.request.password);
    logger.info(`B ${util.inspect(resp, {depth: null})}`);
    callback(null, {
      success: resp.status == 'Success',
      message: resp.status
    });
  } catch (error) {
    logger.error(`Error occurred during register: ${error.message}`);
    callback(error, null); // Sending error response to the client
  }
}

const server = new grpc.Server();
server.addService(BrowserFacadeService, { Login: Login_fromBrowserFacade, Register: Register_fromBrowserFacade });
server.bindAsync('0.0.0.0:50050', grpc.ServerCredentials.createInsecure(), () => {
  console.log('Server running at 0.0.0.0:50050');
});
