### Deploy to heroku

To deploy to heroku for the first time:

```sh
npm install -g heroku-cli
```

To verify your CLI installation use the heroku --version command.

```sh
heroku --version
```

Setup application.

```sh
heroku login

heroku create --region eu <newname>

heroku addons:create redistogo # or any other redis add-on
heroku addons:create mongolab:sandbox

heroku config:set DEBUG=error,log
heroku config:set OB_PROVISIONED=false
heroku config:set OB_DIRECTORY_HOST=http://ob-directory.example.com
heroku config:set SIGNING_KEY='xxx'
heroku config:set OB_ISSUING_CA='xxx'
heroku config:set TRANSPORT_CERT='xxx'
heroku config:set TRANSPORT_KEY='xxx'

git push heroku master
```

## Configuration of ASPSP Authorisation Servers

### Adding and Updating ASPSP authorisation servers

```sh
heroku run npm run updateAuthServersAndOpenIds
```

### Listing available ASPSP authorisation servers

```sh
heroku run npm run listAuthServers
```

### Adding Client Credentials for ASPSP Authorisation Servers

```sh
heroku run npm run saveCreds authServerId=123 clientId=456 clientSecret=789
```

#### Setting client credentials for running against Reference Mock Server

##### Remotely on Heroku

```sh
heroku run npm run saveCreds authServerId=aaaj4NmBD8lQxmLh2O clientId=spoofClientId clientSecret=spoofClientSecret

heroku run npm run saveCreds authServerId=bbbX7tUB4fPIYB0k1m clientId=spoofClientId clientSecret=spoofClientSecret

heroku run npm run saveCreds authServerId=cccbN8iAsMh74sOXhk clientId=spoofClientId clientSecret=spoofClientSecret
```
