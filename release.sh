 ~/go/bin/versionbump $1 ./VERSION && \
    export COMMITSHA="$(git rev-parse --short HEAD)" && \
    export VERSION="$(cat ./VERSION)" && \
    git add ./VERSION && \
    git commit -m "$2" && \
    git checkout -b release/$VERSION && \
    git tag -a $VERSION -m 'Release $VERSION' && \
    git push origin $VERSION && \
    git push -u origin release/$VERSION && \
    git checkout main && \
    git merge release/$VERSION && \
    ##git commit -m 'Merge Release $VERSION into main' && \
    git push origin main
    ##&& \
    ##git checkout release/$VERSION && \
    ##git checkout development
    