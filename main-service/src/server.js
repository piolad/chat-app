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
  '../protos/browser-facade.proto',
  '../protos/auth.proto',
]

const packageDefinition1 = protoLoader.loadSync('../protos/Greeter.proto', {});
const greeterProto = grpc.loadPackageDefinition(packageDefinition1).Greeter;

const packageDefinition = protoLoader.loadSync(protoPahts , {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true
});
const loadedProtos = grpc.loadPackageDefinition(packageDefinition);

const BrowserFacadeService = loadedProtos.browserfacade.BrowserFacade.service;


const client = new loadedProtos.auth.Auth('auth-service:50051', grpc.credentials.createInsecure());

function sayHello(call, callback) {
  callback(null, { message: 'Hello ' + call.request.name });
}
// function definitions
function Login_toAuthService(username, password) {
  return new Promise((resolve, reject) => {

  let request = {
    username: username,
    password: password
  };

  client.login(request, (error, response) => {
    if (!error) {
      console.log('Login Response:', response);
    } else {
      console.error('Error:', error.message);
    }
    console.log('Login status:', response.status);

    logger.info (`A ${(util.inspect(response, {depth: null}))}`)


    resolve(response);
  });
});

}

async function Login_fromBrowserFacade(call, callback) {

  logger.info (`Login request received from ${(util.inspect(call.request, {depth: null}))}`);
  //logger.info(`Login request received from ${call.request}`);
  const resp  = await Login_toAuthService(call.request.username, call.request.password)
  
  logger.info (`B ${(util.inspect(resp, {depth: null}))}`)
    callback(null, { success: resp.status == 'Success',
      username: call.request.username,
      token: '<TEMP TOKEN>',
       message: resp.status });


}

const server = new grpc.Server();

server.addService(BrowserFacadeService, { Login: Login_fromBrowserFacade });
server.addService(greeterProto.Greeter.service, { SayHello: sayHello });
//server.addService(AuthServiceService, { Login: Login_toAuthService });

server.bindAsync('0.0.0.0:50051', grpc.ServerCredentials.createInsecure(), () => {
  
  console.log('Server running at http://0.0.0.0:50051');

});
