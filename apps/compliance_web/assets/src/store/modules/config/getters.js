export default {
  getConfig: configState => configState.parsed,
  getPayload: configState => configState.payload.parsed,
};
