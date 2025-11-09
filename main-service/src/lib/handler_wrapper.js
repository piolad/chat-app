const grpc = require('@grpc/grpc-js');

function handler_wrapper({ name, logger, handler }) {
  return async (call, callback) => {
    const rid = call.metadata?.get('x-request-id')?.[0] || Math.random().toString(36).slice(2);
    try {
      logger.info(`[${name}] start rid=${rid} req=${JSON.stringify(call.request)}`);
      const res = await handler(call, { rid });
      logger.info(`[${name}] ok rid=${rid}`);
      callback(null, res);
    } catch (err) {
      const code = err.code && Number.isInteger(err.code) ? err.code : grpc.status.UNKNOWN;
      logger.error(`[${name}] fail rid=${rid} code=${code} err=${err.message}`);
      callback({ code, details: err.message });
    }
  };
}

module.exports = { handler_wrapper };
