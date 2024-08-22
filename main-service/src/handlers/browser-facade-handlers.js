const util = require('util')

async function Login_fromBrowserFacade(call, callback, data) {
    logger = data.logger
    Login_toAuthService = data.login_toAuthService
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

async function SendMessage_fromBrowserFacade(call, callback, data) {
  logger = data.logger
  SendMessage_toMessageDataCenter = data.sendMessage_toMessageDataCenter
  logger.info(`SendMessage request received from ${util.inspect(call.request, {depth: null})}`);

  try {
    const resp = await SendMessage_toMessageDataCenter(call.request.sender, call.request.receiver,  call.request.message, call.request.timestamp);
    logger.info(`B ${util.inspect(resp, {depth: null})}`);
    callback(null, {
      message: resp.message
    });
  } catch (error) {
    logger.error(`Error occurred during SendMessage: ${error.message}`);
    callback(error, null); // Sending error response to the client
  }
}

module.exports = {
    Login_fromBrowserFacade,
    SendMessage_fromBrowserFacade
}