export default {
  wsConnected: state => state.ws.connection !== null,
  tokenAcquired: state => (tokenName) => {
    const matches = state.tokens.acquired.filter(a => a.value.token_name === tokenName);
    return matches.length > 0;
  },
};
