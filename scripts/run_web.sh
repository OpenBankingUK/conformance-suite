#!/usr/bin/env bash
echo -e '\033[92m  ---> Starting web file watcher ... \033[0m'
cd web && FORCE_COLOR=1 NODE_DISABLE_COLORS=0 yarn build-watch
