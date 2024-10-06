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

async function FetchLastXMessages_fromBrowserFacade(call, callback, data) {
  const logger = data.logger;
  const FetchLastXMessages_fromDataCenter = data.fetchLastXMessages_fromDataCenter;
  logger.info(`FetchLastXMessages request received from ${util.inspect(call.request, {depth: null})}`);

  try {
    const { sender, receiver, startingPoint, count } = call.request;
    const resp = await FetchLastXMessages_fromDataCenter(sender, receiver, startingPoint, count);
    logger.info(`FetchLastXMessages Response: ${util.inspect(resp, { depth: null })}`);
    callback(null, {
      messages: resp.messages,
      count: resp.count,
      hasMore: resp.hasMore
    });
  } catch (error) {
    logger.error(`Error occurred during FetchLastXMessages: ${error.message}`);
    callback(error, null);
  }
}

async function FetchLastXConversations_fromBrowserFacade(call, callback, data) {
  const logger = data.logger;
  const FetchLastXConversations_fromDataCenter = data.fetchLastXConversations_fromDataCenter;
  logger.info(`FetchLastXConversations request received from ${util.inspect(call.request, {depth: null})}`);

  try {
    const request = {
      conversationMember: call.request.conversationMember,
      count: call.request.count,
      start_index: call.request.start_index
    };
    const resp = await FetchLastXConversations_fromDataCenter(request);
    logger.info(`FetchLastXConversations Response: ${util.inspect(resp, { depth: null })}`);
    callback(null, {
      pairs: resp.pairs,
      count: resp.count,
      hasMore: resp.hasMore
    });
  } catch (error) {
    logger.error(`Error occurred during FetchLastXConversations: ${error.message}`);
    callback(error, null);
  }
}

module.exports = {
  Login_fromBrowserFacade,
  SendMessage_fromBrowserFacade,
  FetchLastXMessages_fromBrowserFacade,
  FetchLastXConversations_fromBrowserFacade
};