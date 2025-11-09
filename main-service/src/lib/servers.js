const grpc = require('@grpc/grpc-js');
const util = require('util');


const DEADLINE_MS = 5000;


function promisify(server, method){
    const fn = server[method].bind(server); // fn - new function server.method
    return(req, md = new grpc.Metadata(), timeoutMs = DEADLINE_MS) => 
        new Promise((resolve, reject) => {
            const deadline = new Date(Date.now() + timeoutMs);
            fn(req, md, { deadline }, (err, res) => (err ? reject(err) : resolve(res)));
        });
}


function makeServers(loadedProtos, logger, { authAddr, msgAddr, insecure = true } ){
    const creds = insecure ? grpc.credentials.createInsecure() : grpc.credentials.createSsl();
    const Auth = new loadedProtos.auth.Auth(authAddr, creds);
    const Message = new loadedProtos.service.MessageService(msgAddr, creds);

    return {
        auth: {
            Login: promisify(Auth, 'Login'),
            register: promisify(Auth, 'register'),
        },
        msg: {
            SendMessage: promisify(Message, 'SendMessage'),
            FetchLastXMessages: promisify(Message, 'FetchLastXMessages'),
            FetchLastXConversations: promisify(Message, 'FetchLastXConversations'),
        },
    };
}

module.exports = { makeServers };