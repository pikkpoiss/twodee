set -e  # Aborts if any step fails.

ROOT=`git rev-parse --show-toplevel`

function red {
  echo -e "\033[1;31m$1\033[0m $2"
}

function green {
  echo -e "\033[1;32m$1\033[0m $2"
}

function yellow {
  echo -e "\033[1;33m$1\033[0m $2"
}

function file_exists {
  if [ -n "$FORCE" ]; then
    echo "skipping check for $1"
    return 1
  fi
  if [ -e "$1" ]; then
    echo "$1 exists"
    return 0
  fi
  echo "$1 does not exist"
  return 1
}
