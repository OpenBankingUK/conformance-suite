/**
 * Taken from below and adjusted: https://github.com/danieldocki/slack-clone-vuejs-elixir-phoenix/blob/master/web/src/plugins/socket/index.js
 */
import { Socket as PhoenixSocket } from 'phoenix';

const getEndPoint = () => {
  const { host } = window.location;
  const isSecure = window.location.protocol === 'https:';
  const protocol = isSecure ? 'wss' : 'ws';
  const endPoint = `${protocol}://${host}/api/socket`;

  return endPoint;
};

const socketInstance = new PhoenixSocket(getEndPoint());

export const Socket = {
  connect(token, silent = false) {
    if (this.connClosed()) {
      socketInstance.params.token = token;
      socketInstance.connect();
      if (!silent) {
        // eslint-disable-next-line
        console.log('PhoenixSocket', 'Socket connected!');
      }

      socketInstance.onError((msg) => {
        // eslint-disable-next-line
        console.log('PhoenixSocket', 'there was an error with the connection!', 'msg', msg);
      });
      socketInstance.onClose((msg) => {
        // eslint-disable-next-line
        console.log('PhoenixSocket', 'the connection dropped', 'msg', msg);
      });

      return;
    }

    if (!this.connAvaiable()) {
      socketInstance.connect();
      if (!silent) {
        // eslint-disable-next-line
        console.log('PhoenixSocket', 'Socket reconnected!');
      }
    }
  },
  disconnect() {
    if (this.connClosed()) {
      return;
    }

    socketInstance.disconnect(() => {
      socketInstance.reconnectTimer.reset();
    });
  },
  connAvaiable() {
    return socketInstance && (socketInstance.connectionState() === 'open' ||
      socketInstance.connectionState() === 'connecting');
  },
  connClosed() {
    return socketInstance.connectionState() === 'closed';
  },
  findChannel(id, prefix = 'rooms') {
    return new Promise((resolve, reject) => {
      if (this.connClosed()) {
        const msg = 'NO_SOCKET_CONNECTION: No socket connection, please connect first';
        // eslint-disable-next-line
        console.error('PhoenixSocket', 'connClosed', 'msg', msg);

        reject(new Error(msg));
      } else {
        const topicName = `${prefix}:${id}`;

        let channel = socketInstance.channels.find(ch => ch.topic === topicName);
        if (!channel) {
          channel = socketInstance.channel(topicName, {});
        }

        if (channel.state === 'closed') {
          channel.join()
            .receive('ok', (response) => {
              resolve({ channel, response });
            })
            .receive('error', (msg) => {
              const err = new Error({ msg, err: `[Error] Joined ${channel.topic}`, channel });
              // eslint-disable-next-line
              console.error('PhoenixSocket', 'receive error', 'err', err, 'msg', msg);

              reject(err);
            })
            .receive('timeout', (msg) => {
              // eslint-disable-next-line
              console.error('PhoenixSocket', 'receive timeout', 'msg', msg);
            });
        } else {
          // eslint-disable-next-line
          console.log('PhoenixSocket', 'channel', channel);
          resolve({ channel });
        }
      }
    });
  },
  leaveChannel(id, prefix = 'rooms') {
    return new Promise((resolve, reject) => {
      if (this.connClosed()) {
        const msg = 'NO_SOCKET_CONNECTION: No socket connection, please connect first';
        // eslint-disable-next-line
        console.error('PhoenixSocket', 'leaveChannel', 'connClosed', 'msg', msg);
        reject(new Error(msg));
      } else {
        const topicName = `${prefix}:${id}`;

        const channel = socketInstance.channels.find(ch => ch.topic === topicName);
        if (channel.state === 'closed') {
          reject();
        } else {
          channel.leave()
            .receive('ok', () => {
              resolve({ channel });
            })
            .receive('error', (err) => {
              const msg = `[Error] Left ${channel.topic}`;
              reject(new Error({ err, msg, channel }));
            });
        }
      }
    });
  },
};

// receive connection and params by options
// https://vuejs.org/v2/guide/plugins.html
export default function install(vue) {
  Object.defineProperty(vue.prototype, '$socket', {
    get() {
      return socketInstance;
    },
  });
}
