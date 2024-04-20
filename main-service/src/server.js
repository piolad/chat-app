const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const packageDefinition = protoLoader.loadSync('greet.proto', {});
const greeterProto = grpc.loadPackageDefinition(packageDefinition).greeter;

function sayHello(call, callback) {
  callback(null, { message: 'Hello ' + call.request.name });
}

const server = new grpc.Server();
server.addService(greeterProto.Greeter.service, { sayHello: sayHello });
server.bindAsync('0.0.0.0:50051', grpc.ServerCredentials.createInsecure(), () => {

  
  console.log('Server running at http://0.0.0.0:50051');
});
