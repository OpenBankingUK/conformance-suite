export default {
  hasErrors: state => state.errors && state.errors.length > 0,
  errorMessages: state => state.errors.map(e => (e.message ? e.message : e)),
};
