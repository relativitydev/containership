#!/bin/sh
cd /github/workspace/

while getopts "c:o:t:" opt; do
  case ${opt} in
    c )
      if [ -z "${OPTARG}" ]; then
        echo "::error ::git-chlog path is not set using flag '-c <configuration directory>'"
        exit 1
      fi
      config=$OPTARG
      ;;
    o )
      if [ ! -z ${OPTARG} ]; then
        output="${OPTARG}"
      fi
      ;;
    t )
      tag="${OPTARG}"
      ;;
  esac
done
shift $((OPTIND -1))


if [ -f "${config}/config.yml" ] && [ -f "${config}/CHANGELOG.tpl.md" ]; then
  echo "::debug ::git-chlog: -c '${config}'"
  echo "::debug ::git-chlog: -o '${output}'"
  echo "::debug ::git-chlog: -t '${tag}'"
  echo "::info ::git-chlog executing command: /usr/local/bin/git-chglog --config "${config}/config.yml" ${tag}"

  lasttag=$(git describe --abbrev=0 --tags `git rev-list --tags --skip=1 --max-count=1`)
  changelog=$(/usr/local/bin/git-chglog --config "${config}/config.yml" ${tag})

  echo "----------------------------------------------------------"
  echo "${changelog}"
  echo "----------------------------------------------------------"

  echo "::debug ::git-chlog: -o '$HOME/${output}'"
  if [[ ! -z "$output" ]]; then
    echo "::info ::git-chlog -o options is set. writing changelog to ${output}"
    echo "${changelog}" > "$HOME/${output}"
  fi

  changelog="${changelog//'%'/'%25'}"
  changelog="${changelog//$'\n'/'%0A'}"
  changelog="${changelog//$'\r'/'%0D'}"
  echo "::set-output name=changelog::${changelog}"
  echo "::set-output name=filepath::/home/runner/work/_temp/_github_home/${output}"

else
  echo "::warning ::git-chlog configuration was not found, skipping changelog generation."
fi
