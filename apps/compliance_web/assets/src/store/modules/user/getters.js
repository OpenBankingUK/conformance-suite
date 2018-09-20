export default {
  isSignedIn: state => state.signedIn,
  getAccessToken: state => state.profile && state.profile.access_token,
  isLoading: state => state.loading,
};
