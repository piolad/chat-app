const util = require('util')
const grpc = require('@grpc/grpc-js')
const { handler_wrapper } = require('../lib/handler_wrapper')



function makeHandlers ({logger, services}) {

  return {

    Login: handler_wrapper({
      name: 'Login',
      logger,
      handler: async (call) => {
        const {username, password} = call.request;
        const resp = await services.auth.Login({username, password});
        return {
          success: resp.status == 'Success',
          username: call.request.username,
          token: resp.token,
          message: resp.status,
        }
      }
    }),

    SendMessage: handler_wrapper({
      name: 'SendMessage',
      logger,
      handler: async (call) => {
        const { sender, receiver, message, timestamp } = call.request;
        const resp = await services.msg.SendMessage({ sender, receiver, message, timestamp });
        return { message: resp.message };
      }
    }),

    FetchLastXMessages: handler_wrapper({
      name: 'FetchLastXMessages',
      logger,
      handler: async (call) => {
        const {sender, receiver, startingPoint, count} = call.request;
        const resp = await services.msg.FetchLastXMessages({sender, receiver, startingPoint, count});
        return {
          messages: resp.messages,
          count: resp.count,
          hasMore: resp.hasMore
        }
      }
    }),

    FetchLastXConversations: handler_wrapper({
      name: 'FetchLastXConversations',
      logger,
      handler: async (call) => {
        const { conversationMember, count, start_index } = call.request;
        const resp = await services.msg.FetchLastXConversations({conversationMember, count, start_index});
        return {
          pairs: resp.pairs,
          count: resp.count,
          hasMore: resp.hasMore
        }
      }
    }),

  }  
}

module.exports = { makeHandlers };