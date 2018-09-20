for app in apps/*; do
  pushd $app
  echo running tests for: $app
  if [ "$app" = "apps/log_consumer" ]; then
    mix test --no-start
  else
    mix test
  fi;
  popd
done
