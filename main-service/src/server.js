const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const packageDefinition = protoLoader.loadSync('../protos/browser-facade.proto', {});

const greeterProto = grpc.loadPackageDefinition(packageDefinition).browserfacade;

function Login_fromBrowserFacade(call, callback) {
  
  console.log(call.request);
  callback(null, { message: 'Hello ' + call.request.name });  

}


function sayHello(call, callback) {
  callback(null, { message: 'Hello ' + call.request.name });
}

const server = new grpc.Server();
server.addService(greeterProto.BrowserFacade.service, { Login: Login_fromBrowserFacade });
server.bindAsync('0.0.0.0:50051', grpc.ServerCredentials.createInsecure(), () => {
  
  console.log('Server running at http://0.0.0.0:50051');
});