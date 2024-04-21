const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

const packageDefinition = protoLoader.loadSync('../protos/Greeter.proto', {});
const greeterProto = grpc.loadPackageDefinition(packageDefinition).Greeter;

function sayHello(call, callback) {
  callback(null, { message: 'Hello ' + call.request.name });
}

const server = new grpc.Server();
server.addService(greeterProto.Greeter.service, { SayHello: sayHello });
server.bindAsync('0.0.0.0:50050', grpc.ServerCredentials.createInsecure(), () => {
  console.log('Server running at 0.0.0.0:50050');
});
