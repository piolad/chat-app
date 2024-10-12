const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const winston = require('winston');
const util = require('util')

const BrowserFacadeHandlers = require('./handlers/browser-facade-handlers')


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
  '../protos/service.proto',
]

const packageDefinition = protoLoader.loadSync(protoPahts, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});

const loadedProtos = grpc.loadPackageDefinition(packageDefinition);

// console.log(util.inspect(loadedProtos, {depth: null}));

const BrowserFacadeService = loadedProtos.BrowserFacade.BrowserFacade.service;
const AuthServiceClient = new loadedProtos.auth.Auth('auth-service:50051', grpc.credentials.createInsecure());
const MessageDataCenterClient = new loadedProtos.service.MessageService('message-data-centre:50051', grpc.credentials.createInsecure());

// function definitions
function Login_toAuthService(username, password) {
  return new Promise((resolve, reject) => {

    let request = {
      username: username,
      password: password
    };

    AuthServiceClient.login(request, (error, response) => {
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

    AuthServiceClient.register(request, (error, response) => {
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

async function SendMessage_toMessageDataCenter(sender, receiver, message, timestamp) {
  return new Promise((resolve, reject) => {
    let request = {
      sender: sender,
      receiver: receiver,
      message: message,
      timestamp: timestamp
    };

    MessageDataCenterClient.sendMessage(request, (error, response) => {
      if (error) {
        logger.error(`Error from message-data-center: ${error.message}`);
        reject(error); // Reject the Promise with the error
      } else {
        logger.info('SendMessage Response:', util.inspect(response, {depth: null}));
        resolve(response); // Resolve the Promise with the response
      }
    });
  });
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

async function FetchLastXMessages_toMessageDataCenter(sender, receiver, startingPoint, count) {
  return new Promise((resolve, reject) => {
    const request = {
      sender: sender,
      receiver: receiver,
      startingPoint: startingPoint,
      count: count
    };

    logger.info(`Sending FetchLastXMessages request to message-data-center: ${util.inspect(request, { depth: null })}`);

    MessageDataCenterClient.fetchLastXMessages(request, (error, response) => {
      if (error) {
        logger.error(`Error from message-data-center (FetchLastXMessages): ${error.message}`);
        reject(error);
      } else {
        logger.info('FetchLastXMessages Response:', util.inspect(response, { depth: null }));
        resolve(response);
      }
    });
  });
}

async function FetchLastXConversations_toMessageDataCenter(conversationMember, count, start_index) {
  return new Promise((resolve, reject) => {
    const request = {
      conversationMember: conversationMember,
      count: count,
      startIndex: start_index
    };
    logger.info(`Sending FetchLastXConversations request to message-data-center: ${util.inspect(request, { depth: null })}`);
    MessageDataCenterClient.fetchLastXConversations(request, (error, response) => {
      if (error) {
        logger.error(`Error from message-data-center (FetchLastXConversations): ${error.message}`);
        reject(error);
      } else {
        logger.info('FetchLastXConversations Response:', util.inspect(response, { depth: null }));
        resolve(response);
      }
    });
  });
}

const server = new grpc.Server();
server.addService(BrowserFacadeService, {
  Login: (call, callback) => {
    BrowserFacadeHandlers.Login_fromBrowserFacade(call, callback, {
      logger: logger,
      login_toAuthService: Login_toAuthService
    });
  },
  Register: Register_fromBrowserFacade,
  SendMessage: (call, callback) => {
    BrowserFacadeHandlers.SendMessage_fromBrowserFacade(call, callback, {
      logger: logger,
      sendMessage_toMessageDataCenter: SendMessage_toMessageDataCenter
    });
  },
  FetchLastXMessages: (call, callback) => {
    BrowserFacadeHandlers.FetchLastXMessages_fromBrowserFacade(call, callback, {
      logger: logger,
      fetchLastXMessages_fromDataCenter: FetchLastXMessages_toMessageDataCenter
    });
  },
  FetchLastXConversations: (call, callback) => {
    BrowserFacadeHandlers.FetchLastXConversations_fromBrowserFacade(call, callback, {
      logger: logger,
      fetchLastXConversations_fromDataCenter: FetchLastXConversations_toMessageDataCenter
    });
  }
});

server.bindAsync('0.0.0.0:50050', grpc.ServerCredentials.createInsecure(), () => {
  console.log('Server running at 0.0.0.0:50050');
});
