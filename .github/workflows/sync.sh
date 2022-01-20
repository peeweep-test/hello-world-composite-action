#!/bin/bash
repo=repos/peeweep-test/test-action
for file in $(cat $repo); do
    repo=${repo#"repos"}
    echo sync $file to $repo/$line
    content=`cat $file | base64`
    echo curl -X PUT \
        -H "Accept: application/vnd.github.v3+json" \
        -H 'Authorization: Bearer ${{ secrets.GITHUB_TOKEN }}' \
        "https://api.github.com/repos/peeweep-test/test-action/contents/$line" \
        -d '{"message":"message","content":"'$content'"}'
done