export const gapi = (googleId, isSignedIn) => ({
  load: (a, b) => b(),
  auth2: {
    init: () => ({
      currentUser: {
        get: () => ({
          getId: () => googleId || 'GOOGLEID',
          isSignedIn: () => isSignedIn || true
        })
      },
      signOut: () => {}
    })
  }
});
