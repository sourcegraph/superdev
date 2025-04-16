module.exports = {
  extends: ['react-app', 'react-app/jest'],
  // Force using a single plugin instance
  settings: {
    react: {
      version: 'detect'
    }
  }
};