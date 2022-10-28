ARGS="$*"
COMMITS=$(echo "$ARGS" | cut -f 3- -d" " | awk -v RS='[\n ]' '{print}')

PREVIOUS_RELEASE=$1

NEXT_RELEASE=$2



git fetch
git checkout $PREVIOUS_RELEASE
git checkout -b "release-${NEXT_RELEASE}"
git push -u origin HEAD

git checkout -b "release-${NEXT_RELEASE}-bump"
if [ ! -z "$COMMITS" ]
then
    git cherry-pick $COMMITS
fi
sed "s/\"version\": \".*/\"version\": \"${NEXT_RELEASE}\",/g" plugin.json | sed "s/release_notes_url\": \(.*\)\/v.*/release_notes_url\": \1\/v${NEXT_RELEASE}\",/g" > temp.json
mv temp.json plugin.json
make apply && git add server/manifest.go && git add webapp/src/manifest.*

git add .
git commit -m "bump version ${NEXT_RELEASE}"
git push -u origin HEAD
hub browse
